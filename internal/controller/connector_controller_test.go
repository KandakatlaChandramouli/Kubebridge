package controller

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"

	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"

	v1alpha1 "github.com/KandakatlaChandramouli/kubebridge/api/v1alpha1"
	"github.com/KandakatlaChandramouli/kubebridge/internal/connectors"
)

type mockConnector struct {
	name         string
	connectErr   error
	listErr      error
	resources    []connectors.Resource
	connectCalls int
	listCalls    int
}

func (m *mockConnector) Name() string {
	return m.name
}

func (m *mockConnector) Connect(
	ctx context.Context,
	config connectors.Config,
) error {
	m.connectCalls++
	return m.connectErr
}

func (m *mockConnector) Disconnect(ctx context.Context) error {
	return nil
}

func (m *mockConnector) ListResources(
	ctx context.Context,
) ([]connectors.Resource, error) {
	m.listCalls++

	if m.listErr != nil {
		return nil, m.listErr
	}

	return m.resources, nil
}

func testScheme(t *testing.T) *runtime.Scheme {
	t.Helper()

	scheme := runtime.NewScheme()

	err := v1alpha1.AddToScheme(scheme)
	require.NoError(t, err)

	return scheme
}

func newTestReconciler(
	t *testing.T,
	connectorResource *v1alpha1.Connector,
	registry *connectors.Registry,
) *ConnectorReconciler {
	t.Helper()

	scheme := testScheme(t)

	builder := fake.NewClientBuilder().
		WithScheme(scheme).
		WithStatusSubresource(&v1alpha1.Connector{})

	if connectorResource != nil {
		builder = builder.WithObjects(connectorResource)
	}

	return &ConnectorReconciler{
		Client:   builder.Build(),
		Scheme:   scheme,
		Registry: registry,
	}
}

func reconcileRequest(name, namespace string) ctrl.Request {
	return ctrl.Request{
		NamespacedName: types.NamespacedName{
			Name:      name,
			Namespace: namespace,
		},
	}
}

func TestReconcileUnknownConnector(t *testing.T) {
	registry := connectors.NewRegistry()

	connectorResource := &v1alpha1.Connector{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "unknown",
			Namespace: "default",
		},
		Spec: v1alpha1.ConnectorSpec{
			Type: "unknown",
		},
	}

	reconciler := newTestReconciler(t, connectorResource, registry)

	result, err := reconciler.Reconcile(
		context.Background(),
		reconcileRequest("unknown", "default"),
	)

	require.NoError(t, err)
	assert.Equal(t, 5*time.Minute, result.RequeueAfter)

	var updated v1alpha1.Connector

	err = reconciler.Get(
		context.Background(),
		types.NamespacedName{
			Name:      "unknown",
			Namespace: "default",
		},
		&updated,
	)

	require.NoError(t, err)
	assert.Equal(t, "Failed", updated.Status.Phase)
	assert.Contains(t, updated.Status.Message, "connector type")
}

func TestReconcileSuccessfulConnector(t *testing.T) {
	registry := connectors.NewRegistry()

	mock := &mockConnector{
		name: "github",
		resources: []connectors.Resource{
			{
				ID:   "1",
				Name: "repo-one",
				Type: "github_repository",
			},
			{
				ID:   "2",
				Name: "repo-two",
				Type: "github_repository",
			},
		},
	}

	registry.Register(mock)

	connectorResource := &v1alpha1.Connector{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "github-connector",
			Namespace: "default",
		},
		Spec: v1alpha1.ConnectorSpec{
			Type: "github",
			Config: map[string]string{
				"organization": "kubernetes",
			},
		},
	}

	reconciler := newTestReconciler(t, connectorResource, registry)

	result, err := reconciler.Reconcile(
		context.Background(),
		reconcileRequest("github-connector", "default"),
	)

	require.NoError(t, err)
	assert.Equal(t, 5*time.Minute, result.RequeueAfter)

	assert.Equal(t, 1, mock.connectCalls)
	assert.Equal(t, 1, mock.listCalls)

	var updated v1alpha1.Connector

	err = reconciler.Get(
		context.Background(),
		types.NamespacedName{
			Name:      "github-connector",
			Namespace: "default",
		},
		&updated,
	)

	require.NoError(t, err)
	assert.Equal(t, "Ready", updated.Status.Phase)
	assert.Equal(t, int32(2), updated.Status.ResourcesSynced)
	assert.NotNil(t, updated.Status.LastSyncTime)
	assert.Equal(
		t,
		"successfully discovered 2 resources",
		updated.Status.Message,
	)
}

func TestReconcileConnectorConnectionFailure(t *testing.T) {
	registry := connectors.NewRegistry()

	mock := &mockConnector{
		name:       "github",
		connectErr: errors.New("connection failed"),
	}

	registry.Register(mock)

	connectorResource := &v1alpha1.Connector{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "github-connector",
			Namespace: "default",
		},
		Spec: v1alpha1.ConnectorSpec{
			Type: "github",
		},
	}

	reconciler := newTestReconciler(t, connectorResource, registry)

	result, err := reconciler.Reconcile(
		context.Background(),
		reconcileRequest("github-connector", "default"),
	)

	require.NoError(t, err)
	assert.Equal(t, 5*time.Minute, result.RequeueAfter)

	assert.Equal(t, 1, mock.connectCalls)
	assert.Equal(t, 0, mock.listCalls)

	var updated v1alpha1.Connector

	err = reconciler.Get(
		context.Background(),
		types.NamespacedName{
			Name:      "github-connector",
			Namespace: "default",
		},
		&updated,
	)

	require.NoError(t, err)
	assert.Equal(t, "Failed", updated.Status.Phase)
	assert.Contains(t, updated.Status.Message, "connection failed")
}

func TestReconcileNotFound(t *testing.T) {
	registry := connectors.NewRegistry()

	reconciler := newTestReconciler(t, nil, registry)

	result, err := reconciler.Reconcile(
		context.Background(),
		reconcileRequest("missing", "default"),
	)

	require.NoError(t, err)
	assert.Equal(t, time.Duration(0), result.RequeueAfter)
}
