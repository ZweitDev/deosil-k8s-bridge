package minio

import (
	"context"
	"log"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/utils/pointer"
)

func CreatePVC(clientset *kubernetes.Clientset){
	pvcName := "minio-pvc"	
	namespace := "minio"
	storageSize := "10Gi"
	storageClass := "standard"

	pvc := &corev1.PersistentVolumeClaim{
		ObjectMeta: metav1.ObjectMeta{
			Name: pvcName,
			Namespace: namespace,
		},
		Spec: corev1.PersistentVolumeClaimSpec{
			AccessModes: []corev1.PersistentVolumeAccessMode {
				corev1.ReadWriteOnce,
			},
			Resources: corev1.VolumeResourceRequirements{
				Requests: corev1.ResourceList{
					corev1.ResourceStorage: resource.MustParse(storageSize),
				},
			},
			StorageClassName: pointer.String(storageClass),
		},
	}

	createdPVC, err := clientset.CoreV1().PersistentVolumeClaims(namespace).Create(
		context.TODO(),
		pvc,
		metav1.CreateOptions{},
	)
	if err != nil {
		//panic(fmt.Sprintf("Failed to create PVC: %v", err))
		log.Printf("Failed to create PVC: %v", err)
		log.Printf("Perhaps the PVC already exists?")
	} else {
		log.Printf("Successfully created PVC %q in namespace %q", createdPVC.Name, createdPVC.Namespace)
	}

}