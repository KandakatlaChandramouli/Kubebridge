package controller

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"

	"github.com/KandakatlaChandramouli/kubebridge/api/v1alpha1"
)

func TestResolveConnectorConfigFromSecret(t *testing.T) {
	scheme := runtime.NewScheme()

	require.NoError(t, v1alpha1.AddToScheme(scheme))
	require.NoError(t, corev1.AddToScheme(scheme))

	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "github-credentials",
			Namespace: "kubebridge-system",
		},
		Data: map[string][]byte{
			"token": []byte("test-token"),
		},
	}

	connectorResource := &v1alpha1.Connector{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "github-connector",
			Namespace: "kubebridge-system",
		},
		Spec: v1alpha1.ConnectorSpec{
			Type: "github",
			Config: map[string]string{
				"organization": "kubernetes",
			},
			SecretRef: &v1alpha1.SecretKeySelector{
				Name: "github-credentials",
				Key:  "token",
			},
		},
	}

	reconciler := &ConnectorReconciler{
		Client: fake.NewClientBuilder().
			WithScheme(scheme).
			WithObjects(secret, connectorResource).
			Build(),
	}

	config, err := reconciler.resolveConnectorConfig(
		context.Background(),
		connectorResource,
	)

	require.NoError(t, err)
	require.Equal(t, "kubernetes", config["organization"])
	require.Equal(t, "test-token", config["token"])
}

func TestResolveConnectorConfigFailsWhenSecretKeyMissing(t *testing.T) {
	scheme := runtime.NewScheme()

	require.NoError(t, v1alpha1.AddToScheme(scheme))
	require.NoError(t, corev1.AddToScheme(scheme))

	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "github-credentials",
			Namespace: "kubebridge-system",
		},
		Data: map[string][]byte{
			"wrong-key": []byte("test-token"),
		},
	}

	connectorResource := &v1alpha1.Connector{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "github-connector",
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

	reconciler := &ConnectorReconciler{
		Client: fake.NewClientBuilder().
			WithScheme(scheme).
			WithObjects(secret, connectorResource).
			Build(),
	}

	_, err := reconciler.resolveConnectorConfig(
		context.Background(),
		connectorResource,
	)

	require.Error(t, err)
	require.Contains(t, err.Error(), "does not contain key")
}

func TestSecretLookupUsesConnectorNamespace(t *testing.T) {
	scheme := runtime.NewScheme()

	require.NoError(t, v1alpha1.AddToScheme(scheme))
	require.NoError(t, corev1.AddToScheme(scheme))

	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "github-credentials",
			Namespace: "other-namespace",
		},
		Data: map[string][]byte{
			"token": []byte("test-token"),
		},
	}

	connectorResource := &v1alpha1.Connector{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "github-connector",
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

	reconciler := &ConnectorReconciler{
		Client: fake.NewClientBuilder().
			WithScheme(scheme).
			WithObjects(secret, connectorResource).
			Build(),
	}

	_, err := reconciler.resolveConnectorConfig(
		context.Background(),
		connectorResource,
	)

	require.Error(t, err)
	require.Contains(t, err.Error(), "unable to read secret")
}

func TestSecretReferenceDeepCopy(t *testing.T) {
	original := &v1alpha1.Connector{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test",
			Namespace: "default",
		},
		Spec: v1alpha1.ConnectorSpec{
			Type: "github",
			SecretRef: &v1alpha1.SecretKeySelector{
				Name: "credentials",
				Key:  "token",
			},
		},
	}

	copy := original.DeepCopy()

	require.NotSame(t, original, copy)
	require.NotSame(t, original.Spec.SecretRef, copy.Spec.SecretRef)
	require.Equal(t, original.Spec.SecretRef, copy.Spec.SecretRef)

	copy.Spec.SecretRef.Name = "changed"

	require.Equal(t, "credentials", original.Spec.SecretRef.Name)
	require.Equal(t, types.NamespacedName{
		Name:      "test",
		Namespace: "default",
	}, types.NamespacedName{
		Name:      original.Name,
		Namespace: original.Namespace,
	})
}
