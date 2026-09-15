package controller

import (
	"context"
	"testing"
	"time"

	v1alpha1 "github.com/KandakatlaChandramouli/kubebridge/api/v1alpha1"
	"github.com/KandakatlaChandramouli/kubebridge/internal/connectors"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
)

type eventTestConnector struct {
	name          string
	connectErr    error
	listResources []connectors.Resource
	listErr       error
}

func (c *eventTestConnector) Name() string {
	return c.name
}

func (c *eventTestConnector) Connect(context.Context, connectors.Config) error {
	return c.connectErr
}

func (c *eventTestConnector) Disconnect(context.Context) error {
	return nil
}

func (c *eventTestConnector) ListResources(context.Context) ([]connectors.Resource, error) {
	return c.listResources, c.listErr
}

func newEventTestScheme(t *testing.T) *runtime.Scheme {
	t.Helper()

	s := runtime.NewScheme()
	require.NoError(t, corev1.AddToScheme(s))
	require.NoError(t, v1alpha1.AddToScheme(s))

	return s
}

func newEventTestReconciler(
	t *testing.T,
	connector *v1alpha1.Connector,
	registry *connectors.Registry,
	recorder record.EventRecorder,
) *ConnectorReconciler {
	t.Helper()

	s := newEventTestScheme(t)

	var objects []client.Object
	if connector != nil {
		objects = append(objects, connector)
	}

	builder := fake.NewClientBuilder().
		WithScheme(s).
		WithObjects(objects...)

	if connector != nil {
		builder = builder.WithStatusSubresource(connector)
	}

	kubeClient := builder.Build()

	return &ConnectorReconciler{
		Client:   kubeClient,
		Scheme:   s,
		Registry: registry,
		Recorder: recorder,
	}
}

func reconcileEventTest(
	t *testing.T,
	reconciler *ConnectorReconciler,
	connector *v1alpha1.Connector,
) {
	t.Helper()

	_, err := reconciler.Reconcile(
		context.Background(),
		ctrl.Request{
			NamespacedName: types.NamespacedName{
				Name:      connector.Name,
				Namespace: connector.Namespace,
			},
		},
	)

	require.NoError(t, err)
}

func waitForEvent(t *testing.T, recorder *record.FakeRecorder) string {
	t.Helper()

	select {
	case event := <-recorder.Events:
		return event
	case <-time.After(time.Second):
		t.Fatal("timed out waiting for Kubernetes event")
		return ""
	}
}

func requireNoEvent(t *testing.T, recorder *record.FakeRecorder) {
	t.Helper()

	select {
	case event := <-recorder.Events:
		t.Fatalf("unexpected event: %s", event)
	case <-time.After(100 * time.Millisecond):
	}
}

func TestConnectorReconcilerEmitsUnknownConnectorEvent(t *testing.T) {
	connector := &v1alpha1.Connector{}
	connector.Name = "unknown"
	connector.Namespace = "default"
	connector.Spec.Type = "does-not-exist"

	recorder := record.NewFakeRecorder(10)
	reconciler := newEventTestReconciler(
		t,
		connector,
		connectors.NewRegistry(),
		recorder,
	)

	reconcileEventTest(t, reconciler, connector)

	event := waitForEvent(t, recorder)

	require.Contains(t, event, "UnknownConnectorType")
	require.Contains(t, event, "Warning")
}

func TestConnectorReconcilerEmitsNoEventWhenObjectDoesNotExist(t *testing.T) {
	recorder := record.NewFakeRecorder(10)
	reconciler := newEventTestReconciler(
		t,
		nil,
		connectors.NewRegistry(),
		recorder,
	)

	_, err := reconciler.Reconcile(
		context.Background(),
		ctrl.Request{
			NamespacedName: types.NamespacedName{
				Name:      "missing",
				Namespace: "default",
			},
		},
	)

	require.NoError(t, err)
	requireNoEvent(t, recorder)
}

func TestConnectorReconcilerEmitsConfigurationErrorEvent(t *testing.T) {
	connector := &v1alpha1.Connector{}
	connector.Name = "configuration-error"
	connector.Namespace = "default"
	connector.Spec.Type = "fake"
	connector.Spec.SecretRef = &v1alpha1.SecretKeySelector{
		Name: "missing-secret",
		Key:  "token",
	}

	registry := connectors.NewRegistry()
	registry.Register(&eventTestConnector{name: "fake"})

	recorder := record.NewFakeRecorder(10)
	reconciler := newEventTestReconciler(
		t,
		connector,
		registry,
		recorder,
	)

	reconcileEventTest(t, reconciler, connector)

	event := waitForEvent(t, recorder)

	require.Contains(t, event, "ConfigurationError")
	require.Contains(t, event, "Warning")
}

func TestConnectorReconcilerEmitsConnectionErrorEvent(t *testing.T) {
	connector := &v1alpha1.Connector{}
	connector.Name = "connection-error"
	connector.Namespace = "default"
	connector.Spec.Type = "fake"

	registry := connectors.NewRegistry()
	registry.Register(&eventTestConnector{
		name:       "fake",
		connectErr: context.Canceled,
	})

	recorder := record.NewFakeRecorder(10)
	reconciler := newEventTestReconciler(
		t,
		connector,
		registry,
		recorder,
	)

	reconcileEventTest(t, reconciler, connector)

	event := waitForEvent(t, recorder)

	require.Contains(t, event, "ConnectionError")
	require.Contains(t, event, "Warning")
}

func TestConnectorReconcilerEmitsResourceListErrorEvent(t *testing.T) {
	connector := &v1alpha1.Connector{}
	connector.Name = "resource-list-error"
	connector.Namespace = "default"
	connector.Spec.Type = "fake"

	registry := connectors.NewRegistry()
	registry.Register(&eventTestConnector{
		name:    "fake",
		listErr: context.Canceled,
	})

	recorder := record.NewFakeRecorder(10)
	reconciler := newEventTestReconciler(
		t,
		connector,
		registry,
		recorder,
	)

	reconcileEventTest(t, reconciler, connector)

	event := waitForEvent(t, recorder)

	require.Contains(t, event, "ResourceListError")
	require.Contains(t, event, "Warning")
}

func TestConnectorReconcilerEmitsConnectorReadyEvent(t *testing.T) {
	connector := &v1alpha1.Connector{}
	connector.Name = "ready"
	connector.Namespace = "default"
	connector.Spec.Type = "fake"

	registry := connectors.NewRegistry()
	registry.Register(&eventTestConnector{
		name: "fake",
		listResources: []connectors.Resource{
			{
				ID:   "resource-1",
				Name: "Resource One",
				Type: "test_resource",
			},
			{
				ID:   "resource-2",
				Name: "Resource Two",
				Type: "test_resource",
			},
		},
	})

	recorder := record.NewFakeRecorder(10)
	reconciler := newEventTestReconciler(
		t,
		connector,
		registry,
		recorder,
	)

	reconcileEventTest(t, reconciler, connector)

	event := waitForEvent(t, recorder)

	require.Contains(t, event, "ConnectorReady")
	require.Contains(t, event, "Normal")
	require.Contains(t, event, "Successfully discovered 2 resources")
}
