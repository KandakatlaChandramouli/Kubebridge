package controller

import (
	"context"
	"fmt"
	"time"

	"github.com/go-logr/logr"
	"github.com/prometheus/client_golang/prometheus"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	"github.com/KandakatlaChandramouli/kubebridge/api/v1alpha1"
	"github.com/KandakatlaChandramouli/kubebridge/internal/connectors"
)

var (
	reconciliationTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "kubebridge_connector_reconciliation_total",
			Help: "Total number of Connector reconciliations.",
		},
		[]string{"connector_type", "result"},
	)

	reconciliationDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: "kubebridge_connector_reconciliation_duration_seconds",
			Help: "Duration of Connector reconciliations in seconds.",
		},
		[]string{"connector_type"},
	)

	resourcesSynced = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "kubebridge_connector_resources_synced",
			Help: "Current number of resources discovered by each Connector.",
		},
		[]string{"namespace", "name", "connector_type"},
	)

	connectorStatus = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "kubebridge_connector_status",
			Help: "Current Connector status: 1 for Ready, 0 for Failed or Pending.",
		},
		[]string{"namespace", "name", "connector_type", "phase"},
	)

	secretEventsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "kubebridge_secret_events_total",
			Help: "Total Secret events processed by the Connector watch.",
		},
		[]string{"event"},
	)
)

func init() {
	prometheus.MustRegister(reconciliationTotal)
	prometheus.MustRegister(reconciliationDuration)
	prometheus.MustRegister(resourcesSynced)
	prometheus.MustRegister(connectorStatus)
	prometheus.MustRegister(secretEventsTotal)
}

type ConnectorReconciler struct {
	client.Client
	Scheme   *runtime.Scheme
	Log      logr.Logger
	Registry *connectors.Registry
	Recorder record.EventRecorder
}

func (r *ConnectorReconciler) emitEvent(
	object runtime.Object,
	eventType string,
	reason string,
	message string,
) {
	if r.Recorder != nil {
		r.Recorder.Event(object, eventType, reason, message)
	}
}

func (r *ConnectorReconciler) Reconcile(
	ctx context.Context,
	req ctrl.Request,
) (ctrl.Result, error) {
	start := time.Now()

	connector := &v1alpha1.Connector{}
	if err := r.Get(ctx, req.NamespacedName, connector); err != nil {
		if apierrors.IsNotFound(err) {
			reconciliationTotal.WithLabelValues("unknown", "not_found").Inc()
			return ctrl.Result{}, nil
		}

		reconciliationTotal.WithLabelValues("unknown", "error").Inc()
		return ctrl.Result{}, err
	}

	connectorType := connector.Spec.Type
	log := r.Log.WithValues(
		"connector",
		req.NamespacedName.String(),
		"connectorType",
		connectorType,
	)

	resultLabel := "success"

	defer func() {
		reconciliationDuration.WithLabelValues(connectorType).Observe(time.Since(start).Seconds())
	}()

	defer func() {
		if resultLabel == "success" {
			reconciliationTotal.WithLabelValues(connectorType, "success").Inc()
		}
	}()

	log.Info(
		"starting connector reconciliation",
		"namespace",
		connector.Namespace,
		"name",
		connector.Name,
	)

	if r.Registry == nil {
		resultLabel = "error"
		err := fmt.Errorf("connector registry is nil")
		r.emitEvent(connector, corev1.EventTypeWarning, "ReconcileError", err.Error())

		log.Error(err, "connector reconciliation failed")
		_ = r.updateStatus(ctx, connector, "Failed", 0, err.Error())
		r.recordStatusMetrics(connector, "Failed", 0)
		reconciliationTotal.WithLabelValues(connectorType, "error").Inc()
		return ctrl.Result{RequeueAfter: 5 * time.Minute}, nil
	}

	connectorImpl, ok := r.Registry.Get(connector.Spec.Type)
	if !ok {
		resultLabel = "error"
		err := fmt.Errorf("unknown connector type: %s", connector.Spec.Type)
		r.emitEvent(connector, corev1.EventTypeWarning, "UnknownConnectorType", err.Error())

		log.Error(err, "connector reconciliation failed")
		_ = r.updateStatus(ctx, connector, "Failed", 0, err.Error())
		r.recordStatusMetrics(connector, "Failed", 0)
		reconciliationTotal.WithLabelValues(connectorType, "error").Inc()
		return ctrl.Result{RequeueAfter: 5 * time.Minute}, nil
	}

	config, err := r.resolveConnectorConfig(ctx, connector)
	if err != nil {
		resultLabel = "error"
		log.Error(err, "failed to resolve Connector configuration")
		r.emitEvent(connector, corev1.EventTypeWarning, "ConfigurationError", err.Error())

		_ = r.updateStatus(ctx, connector, "Failed", 0, err.Error())
		r.recordStatusMetrics(connector, "Failed", 0)
		reconciliationTotal.WithLabelValues(connectorType, "error").Inc()
		return ctrl.Result{RequeueAfter: 5 * time.Minute}, nil
	}

	log.Info("connecting to external connector")
	if err := connectorImpl.Connect(ctx, config); err != nil {
		resultLabel = "error"
		log.Error(err, "connector connection failed")
		r.emitEvent(connector, corev1.EventTypeWarning, "ConnectionError", err.Error())

		_ = r.updateStatus(ctx, connector, "Failed", 0, err.Error())
		r.recordStatusMetrics(connector, "Failed", 0)
		reconciliationTotal.WithLabelValues(connectorType, "error").Inc()
		return ctrl.Result{RequeueAfter: 5 * time.Minute}, nil
	}

	defer func() {
		if err := connectorImpl.Disconnect(ctx); err != nil {
			log.Error(err, "connector disconnect failed")
		}
	}()

	resources, err := connectorImpl.ListResources(ctx)
	if err != nil {
		resultLabel = "error"
		log.Error(err, "failed to list connector resources")
		r.emitEvent(connector, corev1.EventTypeWarning, "ResourceListError", err.Error())

		_ = r.updateStatus(ctx, connector, "Failed", 0, err.Error())
		r.recordStatusMetrics(connector, "Failed", 0)
		reconciliationTotal.WithLabelValues(connectorType, "error").Inc()
		return ctrl.Result{RequeueAfter: 5 * time.Minute}, nil
	}

	log.Info(
		"connector reconciliation completed",
		"resourcesSynced",
		len(resources),
	)

	if err := r.updateStatus(
		ctx,
		connector,
		"Ready",
		int32(len(resources)),
		fmt.Sprintf("successfully discovered %d resources", len(resources)),
	); err != nil {
		resultLabel = "error"
		log.Error(err, "failed to update Connector status")

		r.emitEvent(
			connector,
			corev1.EventTypeWarning,
			"StatusUpdateFailed",
			err.Error(),
		)
		reconciliationTotal.WithLabelValues(connectorType, "error").Inc()
		return ctrl.Result{}, err
	}

	r.recordStatusMetrics(connector, "Ready", int32(len(resources)))

	r.emitEvent(
		connector,
		corev1.EventTypeNormal,
		"ConnectorReady",
		fmt.Sprintf("Successfully discovered %d resources", len(resources)),
	)

	return ctrl.Result{RequeueAfter: 5 * time.Minute}, nil
}

func (r *ConnectorReconciler) resolveConnectorConfig(
	ctx context.Context,
	connector *v1alpha1.Connector,
) (connectors.Config, error) {
	config := connectors.Config{}

	for key, value := range connector.Spec.Config {
		config[key] = value
	}

	if connector.Spec.SecretRef == nil {
		return config, nil
	}

	if connector.Spec.SecretRef.Name == "" {
		return nil, fmt.Errorf("secretRef.name must not be empty")
	}

	if connector.Spec.SecretRef.Key == "" {
		return nil, fmt.Errorf("secretRef.key must not be empty")
	}

	secret := &corev1.Secret{}
	key := client.ObjectKey{
		Name:      connector.Spec.SecretRef.Name,
		Namespace: connector.Namespace,
	}

	if err := r.Get(ctx, key, secret); err != nil {
		return nil, fmt.Errorf("unable to read secret %s: %w", key.String(), err)
	}

	value, ok := secret.Data[connector.Spec.SecretRef.Key]
	if !ok {
		return nil, fmt.Errorf(
			"Secret %s does not contain key %q",
			key.String(),
			connector.Spec.SecretRef.Key,
		)
	}

	if len(value) == 0 {
		return nil, fmt.Errorf(
			"Secret %s key %q is empty",
			key.String(),
			connector.Spec.SecretRef.Key,
		)
	}

	config["token"] = string(value)
	return config, nil
}

func (r *ConnectorReconciler) updateStatus(
	ctx context.Context,
	connector *v1alpha1.Connector,
	phase string,
	resourceCount int32,
	message string,
) error {
	connector.Status.Phase = phase
	connector.Status.ResourcesSynced = resourceCount
	connector.Status.Message = message
	now := metav1.Now()
	connector.Status.LastSyncTime = &now

	return r.Status().Update(ctx, connector)
}

func (r *ConnectorReconciler) recordStatusMetrics(
	connector *v1alpha1.Connector,
	phase string,
	resourceCount int32,
) {
	resourcesSynced.WithLabelValues(
		connector.Namespace,
		connector.Name,
		connector.Spec.Type,
	).Set(float64(resourceCount))

	for _, existingPhase := range []string{"Ready", "Failed", "Pending"} {
		connectorStatus.DeleteLabelValues(
			connector.Namespace,
			connector.Name,
			connector.Spec.Type,
			existingPhase,
		)
	}

	value := float64(0)
	if phase == "Ready" {
		value = 1
	}

	connectorStatus.WithLabelValues(
		connector.Namespace,
		connector.Name,
		connector.Spec.Type,
		phase,
	).Set(value)
}

func (r *ConnectorReconciler) secretToConnectorRequests(
	ctx context.Context,
	obj client.Object,
) []reconcile.Request {
	secret, ok := obj.(*corev1.Secret)
	if !ok {
		return nil
	}

	secretEventsTotal.WithLabelValues(eventType(secret)).Inc()

	connectorsList := &v1alpha1.ConnectorList{}
	if err := r.List(ctx, connectorsList); err != nil {
		r.Log.Error(err, "failed to list Connectors for Secret event")
		return nil
	}

	requests := make([]reconcile.Request, 0)
	for i := range connectorsList.Items {
		item := &connectorsList.Items[i]

		if item.Namespace != secret.Namespace ||
			item.Spec.SecretRef == nil ||
			item.Spec.SecretRef.Name != secret.Name {
			continue
		}

		requests = append(requests, reconcile.Request{
			NamespacedName: types.NamespacedName{
				Name:      item.Name,
				Namespace: item.Namespace,
			},
		})
	}

	return requests
}

func eventType(secret *corev1.Secret) string {
	switch secret.ResourceVersion {
	case "":
		return "unknown"
	default:
		return "updated"
	}
}

func (r *ConnectorReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&v1alpha1.Connector{}).
		Watches(
			&corev1.Secret{},
			handler.EnqueueRequestsFromMapFunc(r.secretToConnectorRequests),
		).
		Complete(r)
}

var _ client.Object = (*v1alpha1.Connector)(nil)
