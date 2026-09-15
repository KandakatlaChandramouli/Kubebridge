package controller

import (
	"context"
	"testing"

	"github.com/prometheus/client_golang/prometheus/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	v1alpha1 "github.com/KandakatlaChandramouli/kubebridge/api/v1alpha1"
	"github.com/KandakatlaChandramouli/kubebridge/internal/connectors"
)

func TestRecordStatusMetricsReady(t *testing.T) {
	connector := &v1alpha1.Connector{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "metrics-ready",
			Namespace: "metrics-test",
		},
		Spec: v1alpha1.ConnectorSpec{
			Type: "github",
		},
	}

	reconciler := &ConnectorReconciler{}
	reconciler.recordStatusMetrics(connector, "Ready", 3)

	assert.Equal(
		t,
		float64(3),
		testutil.ToFloat64(
			resourcesSynced.WithLabelValues(
				"metrics-test",
				"metrics-ready",
				"github",
			),
		),
	)

	assert.Equal(
		t,
		float64(1),
		testutil.ToFloat64(
			connectorStatus.WithLabelValues(
				"metrics-test",
				"metrics-ready",
				"github",
				"Ready",
			),
		),
	)

	require.Equal(
		t,
		float64(0),
		testutil.ToFloat64(
			connectorStatus.WithLabelValues(
				"metrics-test",
				"metrics-ready",
				"github",
				"Failed",
			),
		),
	)
}

func TestRecordStatusMetricsFailedReplacesReady(t *testing.T) {
	connector := &v1alpha1.Connector{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "metrics-transition",
			Namespace: "metrics-test",
		},
		Spec: v1alpha1.ConnectorSpec{
			Type: "github",
		},
	}

	reconciler := &ConnectorReconciler{}

	reconciler.recordStatusMetrics(connector, "Ready", 5)
	reconciler.recordStatusMetrics(connector, "Failed", 0)

	assert.Equal(
		t,
		float64(0),
		testutil.ToFloat64(
			connectorStatus.WithLabelValues(
				"metrics-test",
				"metrics-transition",
				"github",
				"Ready",
			),
		),
	)

	assert.Equal(
		t,
		float64(0),
		testutil.ToFloat64(
			connectorStatus.WithLabelValues(
				"metrics-test",
				"metrics-transition",
				"github",
				"Failed",
			),
		),
	)

	assert.Equal(
		t,
		float64(0),
		testutil.ToFloat64(
			resourcesSynced.WithLabelValues(
				"metrics-test",
				"metrics-transition",
				"github",
			),
		),
	)
}

func TestRecordStatusMetricsDoesNotDeleteOtherConnector(t *testing.T) {
	first := &v1alpha1.Connector{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "first",
			Namespace: "metrics-test",
		},
		Spec: v1alpha1.ConnectorSpec{
			Type: "github",
		},
	}

	second := &v1alpha1.Connector{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "second",
			Namespace: "metrics-test",
		},
		Spec: v1alpha1.ConnectorSpec{
			Type: "github",
		},
	}

	reconciler := &ConnectorReconciler{}

	reconciler.recordStatusMetrics(first, "Ready", 4)
	reconciler.recordStatusMetrics(second, "Failed", 0)

	assert.Equal(
		t,
		float64(1),
		testutil.ToFloat64(
			connectorStatus.WithLabelValues(
				"metrics-test",
				"first",
				"github",
				"Ready",
			),
		),
	)

	assert.Equal(
		t,
		float64(0),
		testutil.ToFloat64(
			connectorStatus.WithLabelValues(
				"metrics-test",
				"second",
				"github",
				"Failed",
			),
		),
	)

	assert.Equal(
		t,
		float64(4),
		testutil.ToFloat64(
			resourcesSynced.WithLabelValues(
				"metrics-test",
				"first",
				"github",
			),
		),
	)
}

func TestReconcileMetrics(t *testing.T) {
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

	connector := &v1alpha1.Connector{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "metrics-reconcile",
			Namespace: "default",
		},
		Spec: v1alpha1.ConnectorSpec{
			Type: "github",
		},
	}

	reconciler := newTestReconciler(t, connector, registry)

	beforeSuccess := testutil.ToFloat64(
		reconciliationTotal.WithLabelValues("github", "success"),
	)

	_, err := reconciler.Reconcile(
		context.Background(),
		reconcileRequest("metrics-reconcile", "default"),
	)
	require.NoError(t, err)

	afterSuccess := testutil.ToFloat64(
		reconciliationTotal.WithLabelValues("github", "success"),
	)

	assert.Equal(t, float64(1), afterSuccess-beforeSuccess)

	assert.Greater(
		t,
		testutil.ToFloat64(
			resourcesSynced.WithLabelValues(
				"default",
				"metrics-reconcile",
				"github",
			),
		),
		float64(0),
	)
}
