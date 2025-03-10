package minio

import (
	"context"
	"fmt"
	"log"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func CreateNamespace(clientset *kubernetes.Clientset) {
	namespaceName := "minio"

	namespace := &corev1.Namespace {
		ObjectMeta: metav1.ObjectMeta{
			Name: namespaceName,
		},
	}

	createdNamespace, err := clientset.CoreV1().Namespaces().Create(
		context.TODO(),
		namespace,
		metav1.CreateOptions{},
	)

	if err != nil {
		if err.Error() == fmt.Sprintf(`namespaces "%s" already exists`, namespaceName) {
			log.Printf("Namespace %q already exists", namespaceName)
			return
		}
		// We could panic here, but for the purposes of demonstration, lets just log
	}
	
	log.Printf("Successfullly created namespace %q", createdNamespace.Name)
}