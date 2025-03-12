package traefik

import (
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/pointer"
)

func CreateTraefikDeployment() *appsv1.Deployment {
	return &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "traefik",
			Namespace: "traefik",
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: pointer.Int32(1),
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": "traefik",
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": "traefik",
					},
				},
				Spec: corev1.PodSpec{
					ServiceAccountName: "traefik",
					Containers: []corev1.Container{
						{
							Name:  "traefik",
							Image: "traefik:v2.9",
							Args: []string{
								"--providers.kubernetescrd",
								"--entrypoints.web.address=:80",
							},
							Ports: []corev1.ContainerPort{
								{
									ContainerPort: 80,
									Name:          "web",
								},
							},
						},
					},
				},
			},
		},
	}
}