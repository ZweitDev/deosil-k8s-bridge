package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"path/filepath"

	"deosil-k8s-bridge/lib/k8s/traefik"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

func main() {
	var kubeconfigPath *string
	if home := homedir.HomeDir(); home != "" {
		kubeconfigPath = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfigPath = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()
	// Load kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfigPath)
	if err != nil {
		log.Fatalf("Failed to load kubeconfig: %v", err)
	}

	// Create clientsets
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatalf("Failed to create clientset: %v", err)
	}
	dynamicClient, err := dynamic.NewForConfig(config)
	if err != nil {
		log.Fatalf("Failed to create dynamic client: %v", err)
	}

	// Create Traefik resources
	createResource(clientset.CoreV1().ServiceAccounts("traefik"), traefik.CreateTraefikServiceAccount())
	createResource(clientset.RbacV1().ClusterRoles(), traefik.CreateTraefikClusterRole())
	createResource(clientset.RbacV1().ClusterRoleBindings(), traefik.CreateTraefikClusterRoleBinding())
	createResource(clientset.AppsV1().Deployments("traefik"), traefik.CreateTraefikDeployment())
	createResource(clientset.CoreV1().Services("traefik"), traefik.CreateTraefikService())

	// Create IngressRoute using dynamic client
	ingressRouteGVR := traefik.IngressRouteGVR()
	ingressRoute := traefik.CreateMinioIngressRoute()
	_, err = dynamicClient.Resource(ingressRouteGVR).Namespace("minio").Create(
		context.TODO(),
		ingressRoute,
		v1.CreateOptions{},
	)
	if err != nil {
		log.Printf("Failed to create IngressRoute: %v", err)
	}

	fmt.Println("Traefik and Minio IngressRoute created successfully!")
}

// Helper function to create Kubernetes resources
func createResource[T any](client interface {
	Create(context.Context, *T, v1.CreateOptions) (*T, error)
}, obj *T) {
	_, err := client.Create(context.TODO(), obj, v1.CreateOptions{})
	if err != nil {
		log.Printf("Failed to create resource: %v", err)
	}
}