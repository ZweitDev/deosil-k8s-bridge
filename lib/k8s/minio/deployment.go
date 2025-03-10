package minio

import (
	"context"
	"log"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/utils/pointer"
)

func CreateMinioDeployment(clientset *kubernetes.Clientset) {
	namespace := "minio"
	deploymentName := "minio"
	replicas := int32(1)
	minioImage := "minio/minio:latest"
	pvcName := "minio-pvc"

	// Define the Minio StatefulSet
	statefulSet := &appsv1.StatefulSet{
		ObjectMeta: metav1.ObjectMeta{
			Name: deploymentName,
			Namespace: namespace,
			Labels: map[string]string{
				"app": "minio",
			},
		},
		Spec: appsv1.StatefulSetSpec{
			Replicas: pointer.Int32(replicas),
			Selector: &metav1.LabelSelector {
				MatchLabels: map[string]string{
					"app": "minio",
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta {
					Labels: map[string]string{
						"app": "minio",
					},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name: "minio",
							Image: minioImage,
							Args: []string{
								"server",
								"/data",
							},
							Ports: []corev1.ContainerPort{
								{
									ContainerPort: 9000,
								},
							},
							Env: []corev1.EnvVar {
								{
									Name: "MINIO_ROOT_USER",
									Value: "minioadmin",
								},
								{
									Name: "MINIO_ROOT_PASSWORD",
									Value: "minioadmin",
								},
							},
							VolumeMounts: []corev1.VolumeMount{
								{
									Name: "minio-data",
									MountPath: "/data",
								},
							},
						},
					},
					Volumes: []corev1.Volume{
						{
							Name: "minio-data",
							VolumeSource: corev1.VolumeSource {
								PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
									ClaimName: pvcName,
								},
							},
						},
					},
				},
			},
		},
	}

	// Actually create the Stateful Set now...
	createdStatefulSet, err := clientset.AppsV1().StatefulSets(namespace).Create(
		context.TODO(),
		statefulSet,
		metav1.CreateOptions{},
	)

	if err != nil {
		log.Printf("Failed to create Minio StatefulSet: %v", err)
		log.Printf("Perhaps the Minio StatefulSet already exists?")
	} else {
		log.Printf("Successfully created Minio StatefulSet %q in namespace %q\n", createdStatefulSet.Name, createdStatefulSet.Namespace)
	}
}