package k8sclient

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func InitClient() (*kubernetes.Clientset, error) {

	clientSet := &kubernetes.Clientset{}

	userHomeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("InitClient: error getting user home dir: %w", err)
	}

	kubeConfigPath := filepath.Join(userHomeDir, ".kube", "config")

	kubeConfig, err := clientcmd.BuildConfigFromFlags("", kubeConfigPath)
	if err != nil {
		return nil, fmt.Errorf("InitClient: error getting Kubernetes config: %w", err)
	}

	clientSet, err = kubernetes.NewForConfig(kubeConfig)
	if err != nil {
		return nil, fmt.Errorf("InitClient: error getting Kubernetes clientset: %w", err)
	}

	return clientSet, nil
}

func GetNamespaces(clientSet kubernetes.Clientset) ([]string, error) {

	namespaces, err := clientSet.CoreV1().Namespaces().List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("GetNamespaces: %w", err)
	}

	var currentNamespaces []string

	for _, ns := range namespaces.Items {
		currentNamespaces = append(currentNamespaces, ns.Name)
	}

	return currentNamespaces, nil
}

func GetPods(clientSet kubernetes.Clientset, namespace string) (map[string]string, error) {

	// Check namespace exists first
	_, err := clientSet.CoreV1().Namespaces().Get(context.Background(), namespace, metav1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("GetPods: %w", err)
	}

	pods, err := clientSet.CoreV1().Pods(namespace).List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("GetPods: failed to get Pods from K8s API: %w", err)
	}

	currentPods := make(map[string]string)

	for _, pod := range pods.Items {
		currentPods[pod.Name] = string(pod.Status.Phase)
	}

	return currentPods, nil
}

func CreateNamespace(clientSet kubernetes.Clientset, name string) error {

	ns := &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
	}

	_, err := clientSet.CoreV1().Namespaces().Create(context.Background(), ns, metav1.CreateOptions{})
	if err != nil {
		return fmt.Errorf("CreateNamespaces: %w", err)
	}

	return nil
}
