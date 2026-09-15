package controller

import (
	"context"
	"testing"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/utils/ptr"

	"github.com/go-logr/logr"
	"github.com/stretchr/testify/require"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"

	"github.com/KandakatlaChandramouli/kubebridge/api/v1alpha1"
	"github.com/KandakatlaChandramouli/kubebridge/internal/connectors"
)

type e2eConnector struct {
	config connectors.Config
}

func (c *e2eConnector) Name() string {
	return "e2e"
}

func (c *e2eConnector) Connect(_ context.Context, config connectors.Config) error {
	c.config = config
	return nil
}

func (c *e2eConnector) Disconnect(_ context.Context) error {
	return nil
}

func (c *e2eConnector) ListResources(_ context.Context) ([]connectors.Resource, error) {
	return []connectors.Resource{
		{
			ID:   c.config["token"],
			Name: "secret-backed-resource",
			Type: "e2e",
		},
	}, nil
}

func TestSecretRotationTriggersReconciliationFlow(t *testing.T) {
	s := runtime.NewScheme()
	require.NoError(t, scheme.AddToScheme(s))
	require.NoError(t, v1alpha1.AddToScheme(s))

	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "connector-secret",
			Namespace: "default",
		},
		Data: map[string][]byte{
			"token": []byte("initial-token"),
		},
	}

	connector := &v1alpha1.Connector{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "github-connector",
			Namespace: "default",
		},
		Spec: v1alpha1.ConnectorSpec{
			Type: "e2e",
			SecretRef: &v1alpha1.SecretKeySelector{
				Name: "connector-secret",
				Key:  "token",
			},
		},
	}

	cl := fake.NewClientBuilder().
		WithScheme(s).
		WithStatusSubresource(&v1alpha1.Connector{}).
		WithObjects(secret, connector).
		Build()

	mock := &e2eConnector{}
	registry := connectors.NewRegistry()
	registry.Register(mock)

	r := &ConnectorReconciler{
		Client:   cl,
		Scheme:   s,
		Log:      logr.Discard(),
		Registry: registry,
	}

	firstRequest := ctrl.Request{
		NamespacedName: types.NamespacedName{
			Name:      "github-connector",
			Namespace: "default",
		},
	}

	_, err := r.Reconcile(context.Background(), firstRequest)
	require.NoError(t, err)
	require.Equal(t, "initial-token", mock.config["token"])

	updatedSecret := &corev1.Secret{}
	require.NoError(t, cl.Get(
		context.Background(),
		client.ObjectKey{
			Name:      "connector-secret",
			Namespace: "default",
		},
		updatedSecret,
	))

	updatedSecret.Data["token"] = []byte("rotated-token")
	require.NoError(t, cl.Update(context.Background(), updatedSecret))

	secretRequests := r.secretToConnectorRequests(
		context.Background(),
		updatedSecret,
	)

	require.Len(t, secretRequests, 1)
	require.Equal(t, firstRequest.NamespacedName, secretRequests[0].NamespacedName)

	_, err = r.Reconcile(context.Background(), secretRequests[0])
	require.NoError(t, err)

	require.Equal(t, "rotated-token", mock.config["token"])

	storedConnector := &v1alpha1.Connector{}
	require.NoError(t, cl.Get(
		context.Background(),
		client.ObjectKey{
			Name:      "github-connector",
			Namespace: "default",
		},
		storedConnector,
	))

	require.Equal(t, "Ready", storedConnector.Status.Phase)
	require.Equal(t, int32(1), storedConnector.Status.ResourcesSynced)
	require.Equal(t, ptr.To("successfully discovered 1 resources"), ptr.To(storedConnector.Status.Message))
}
