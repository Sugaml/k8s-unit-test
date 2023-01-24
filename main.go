package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func uppercasePodLabel(clientset kubernetes.Interface, namespace, podName, labelKey string) (string, error) {
	pod, err := clientset.CoreV1().Pods(namespace).Get(
		context.Background(),
		podName,
		v1.GetOptions{},
	)
	if err != nil {
		return "", err
	}

	labelValue, ok := pod.ObjectMeta.Labels[labelKey]
	if !ok {
		return "", fmt.Errorf("no label with key %s for pod %s/%s", labelKey, namespace, podName)
	}
	return strings.ToUpper(labelValue), nil
}

func main() {
	fmt.Println("Get Kubernetes pod label to uppercase")

	userHomeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Printf("error getting user home dir: %v\n", err)
		os.Exit(1)
	}
	kubeConfigPath := filepath.Join(userHomeDir, ".kube", "config")
	fmt.Printf("Using kubeconfig: %s\n", kubeConfigPath)

	kubeConfig, err := clientcmd.BuildConfigFromFlags("", kubeConfigPath)
	if err != nil {
		fmt.Printf("error getting Kubernetes config: %v\n", err)
		os.Exit(1)
	}

	clientset, err := kubernetes.NewForConfig(kubeConfig)
	if err != nil {
		fmt.Printf("error getting Kubernetes clientset: %v\n", err)
		os.Exit(1)
	}

	labelValue, err := uppercasePodLabel(clientset, "01cloud-staging", "01cloud-staging-api-c58f64574-trnbw", "app")
	if err != nil {
		fmt.Printf("error getting pod label: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Pod label value: %s\n", labelValue)
}
