package k8sclient

import (
	"context"
	"fmt"

	discovery "github.com/gkarthiks/k8s-discovery"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func InitClient() (*discovery.K8s, error) {

	// https://medium.com/swlh/clientset-module-for-in-cluster-and-out-cluster-3f0d80af79ed
	k8sClients, err := discovery.NewK8s()
	if err != nil {
		return nil, fmt.Errorf("InitClient: error getting Kubernetes clientset: %w", err)
	}

	return k8sClients, nil
}

func GetNamespaces(k8sClients discovery.K8s) ([]string, error) {

	namespaces, err := k8sClients.Clientset.CoreV1().Namespaces().List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("GetNamespaces: %w", err)
	}

	var currentNamespaces []string

	for _, ns := range namespaces.Items {
		currentNamespaces = append(currentNamespaces, ns.Name)
	}

	return currentNamespaces, nil
}

func GetPods(k8sClients discovery.K8s, namespace string) (map[string]string, error) {

	// Check namespace exists first
	_, err := k8sClients.Clientset.CoreV1().Namespaces().Get(context.Background(), namespace, metav1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("GetPods: %w", err)
	}

	pods, err := k8sClients.Clientset.CoreV1().Pods(namespace).List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("GetPods: failed to get Pods from K8s API: %w", err)
	}

	currentPods := make(map[string]string)

	for _, pod := range pods.Items {
		currentPods[pod.Name] = string(pod.Status.Phase)
	}

	return currentPods, nil
}

func CreateNamespace(k8sClients discovery.K8s, name string) error {

	ns := &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
	}

	_, err := k8sClients.Clientset.CoreV1().Namespaces().Create(context.Background(), ns, metav1.CreateOptions{})
	if err != nil {
		return fmt.Errorf("CreateNamespaces: %w", err)
	}

	return nil
}
