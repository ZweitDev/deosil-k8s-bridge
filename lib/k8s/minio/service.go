package minio

import (
	"context"
	"fmt"
	"log"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/kubernetes"
)

func CreateMinioService(clientset *kubernetes.Clientset) {
	namespace := "minio"
	serviceName := "minio-service"

	service := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta {
			Name: serviceName,
			Namespace: namespace,
			Labels: map[string]string{
				"app": "minio",
			},
		},
		Spec: corev1.ServiceSpec {
			Selector: map[string]string{
				"app": "minio",
			},
			Ports: []corev1.ServicePort{
				{
					Name: "http",
					Port: 9000,
					TargetPort: intstr.FromInt32(9000),
					Protocol: corev1.ProtocolTCP,
				},
			},
			Type: corev1.ServiceTypeLoadBalancer,
		},
	}

	createdService, err := clientset.CoreV1().Services(namespace).Create(
		context.TODO(),
		service,
		metav1.CreateOptions{},
	)

	if err != nil {
		log.Printf("Failed to create Minio Service: %v", err)
		log.Printf("Perhaps the Minio Service already exists?")
	} else {
		log.Printf("Successfully created Minio Service %q in namespace %q\n", createdService.Name, createdService.Namespace)
	}
}

func GetMinioServiceEndpoint(clientset *kubernetes.Clientset) (string, error) {
	namespace := "minio"
	serviceName := "minio-service"

	service, err := clientset.CoreV1().Services(namespace).Get(
		context.TODO(),
		serviceName,
		metav1.GetOptions{},
	)

	if err != nil {
		return "", fmt.Errorf("failed to get Minio Service: %v", err)
	}

	clusterIP := service.Spec.ClusterIP
	if clusterIP == "" {
		return "", fmt.Errorf("ClusterIP not found for Minio Service")
	}

	port := service.Spec.Ports[0].Port

	// Return the endpoint in the format "IP:Port"
	return fmt.Sprintf("%s:%d", clusterIP, port), nil
}
