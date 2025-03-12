package traefik

import (
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

func CreateMinioIngressRoute() *unstructured.Unstructured {
	return &unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": "traefik.containo.us/v1alpha1",
			"kind":       "IngressRoute",
			"metadata": map[string]interface{}{
				"name":      "minio-ingressroute",
				"namespace": "minio",
			},
			"spec": map[string]interface{}{
				"entryPoints": []string{"web"},
				"routes": []map[string]interface{}{
					{
						"match": "Host(`minio.local`)",
						"kind":  "Rule",
						"services": []map[string]interface{}{
							{
								"name": "minio-service",
								"port": 9000,
							},
						},
					},
				},
			},
		},
	}
}

// Define the GroupVersionResource for IngressRoute
func IngressRouteGVR() schema.GroupVersionResource {
	return schema.GroupVersionResource{
		Group:    "traefik.containo.us",
		Version:  "v1alpha1",
		Resource: "ingressroutes",
	}
}