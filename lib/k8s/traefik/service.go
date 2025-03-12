package traefik

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

func CreateTraefikService() *corev1.Service {
	return &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "traefik",
			Namespace: "traefik",
		},
		Spec: corev1.ServiceSpec{
			Type: corev1.ServiceTypeLoadBalancer,
			Selector: map[string]string{
				"app": "traefik",
			},
			Ports: []corev1.ServicePort{
				{
					Name:       "web",
					Port:       80,
					TargetPort: intstr.FromInt(80),
				},
			},
		},
	}
}