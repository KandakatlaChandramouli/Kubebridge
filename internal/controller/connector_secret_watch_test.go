package controller

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"

	"github.com/KandakatlaChandramouli/kubebridge/api/v1alpha1"
)

func TestSecretToConnectorRequests(t *testing.T) {
	scheme := runtime.NewScheme()

	require.NoError(t, v1alpha1.AddToScheme(scheme))
	require.NoError(t, corev1.AddToScheme(scheme))

	matching := &v1alpha1.Connector{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "matching",
			Namespace: "kubebridge-system",
		},
		Spec: v1alpha1.ConnectorSpec{
			Type: "github",
			SecretRef: &v1alpha1.SecretKeySelector{
				Name: "github-credentials",
				Key:  "token",
			},
		},
	}

	differentSecret := &v1alpha1.Connector{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "different-secret",
			Namespace: "kubebridge-system",
		},
		Spec: v1alpha1.ConnectorSpec{
			Type: "github",
			SecretRef: &v1alpha1.SecretKeySelector{
				Name: "other-credentials",
				Key:  "token",
			},
		},
	}

	differentNamespace := &v1alpha1.Connector{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "different-namespace",
			Namespace: "default",
		},
		Spec: v1alpha1.ConnectorSpec{
			Type: "github",
			SecretRef: &v1alpha1.SecretKeySelector{
				Name: "github-credentials",
				Key:  "token",
			},
		},
	}

	noSecretRef := &v1alpha1.Connector{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "no-secret-ref",
			Namespace: "kubebridge-system",
		},
		Spec: v1alpha1.ConnectorSpec{
			Type: "github",
		},
	}

	reconciler := &ConnectorReconciler{
		Client: fake.NewClientBuilder().
			WithScheme(scheme).
			WithObjects(
				matching,
				differentSecret,
				differentNamespace,
				noSecretRef,
			).
			Build(),
	}

	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "github-credentials",
			Namespace: "kubebridge-system",
		},
	}

	requests := reconciler.secretToConnectorRequests(
		context.Background(),
		secret,
	)

	require.Len(t, requests, 1)
	require.Equal(
		t,
		client.ObjectKey{
			Name:      "matching",
			Namespace: "kubebridge-system",
		},
		requests[0].NamespacedName,
	)
}

func TestSecretToConnectorRequestsNoMatches(t *testing.T) {
	scheme := runtime.NewScheme()

	require.NoError(t, v1alpha1.AddToScheme(scheme))
	require.NoError(t, corev1.AddToScheme(scheme))

	connectorResource := &v1alpha1.Connector{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "connector",
			Namespace: "kubebridge-system",
		},
		Spec: v1alpha1.ConnectorSpec{
			Type: "github",
			SecretRef: &v1alpha1.SecretKeySelector{
				Name: "different-secret",
				Key:  "token",
			},
		},
	}

	reconciler := &ConnectorReconciler{
		Client: fake.NewClientBuilder().
			WithScheme(scheme).
			WithObjects(connectorResource).
			Build(),
	}

	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "github-credentials",
			Namespace: "kubebridge-system",
		},
	}

	requests := reconciler.secretToConnectorRequests(
		context.Background(),
		secret,
	)

	require.Empty(t, requests)
}
