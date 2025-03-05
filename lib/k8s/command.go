package k8s

import (
	"context"
	"encoding/json"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func GetPods(clientset *kubernetes.Clientset, namespace string) ([]string, error) {
	pods, err := clientset.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	var podNames []string
	
	for _, pod := range pods.Items {
		podData, err := json.Marshal(pod)
		if err != nil {
			break	
		}

		data := string(podData)

		podNames = append(podNames, data)
	}

	return podNames, nil
}

func GetNodes(clientset *kubernetes.Clientset) ([]string, error) {
	nodes, err := clientset.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	var nodeNames []string

	for _, node := range nodes.Items {
		nodeNames = append(nodeNames, node.Name)
	}

	return nodeNames, nil

}