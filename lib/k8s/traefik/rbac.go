package traefik

import (
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func CreateTraefikClusterRole() *rbacv1.ClusterRole {
	return &rbacv1.ClusterRole{
		ObjectMeta: metav1.ObjectMeta{
			Name: "traefik",
		},
		Rules: []rbacv1.PolicyRule{
			{
				APIGroups: []string{""},
				Resources: []string{"services", "endpoints", "secrets"},
				Verbs:     []string{"get", "list", "watch"},
			},
			{
				APIGroups: []string{"networking.k8s.io"},
				Resources: []string{"ingresses"},
				Verbs:     []string{"get", "list", "watch"},
			},
			{
				APIGroups: []string{"traefik.containo.us"},
				Resources: []string{"ingressroutes", "middlewares"},
				Verbs:     []string{"get", "list", "watch"},
			},
		},
	}
}

func CreateTraefikClusterRoleBinding() *rbacv1.ClusterRoleBinding {
	return &rbacv1.ClusterRoleBinding{
		ObjectMeta: metav1.ObjectMeta{
			Name: "traefik",
		},
		RoleRef: rbacv1.RoleRef{
			APIGroup: "rbac.authorization.k8s.io",
			Kind:     "ClusterRole",
			Name:     "traefik",
		},
		Subjects: []rbacv1.Subject{
			{
				Kind:      "ServiceAccount",
				Name:      "traefik",
				Namespace: "traefik",
			},
		},
	}
}