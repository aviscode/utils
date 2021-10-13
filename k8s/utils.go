package k8s

import (
	"context"
	"errors"
	"fmt"
	"github.com/anthhub/forwarder"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"strings"
)

// CreatServicePortForwarder this func will create a port forwarding to a given service.
func CreatServicePortForwarder(localPort, remotePort int, nameSpace, serviceName, kubeconfig string) (*forwarder.Result, error) {
	options := []*forwarder.Option{{LocalPort: localPort, RemotePort: remotePort, Namespace: nameSpace, ServiceName: serviceName, Source: "svc/" + serviceName}}
	return forwarder.WithForwarders(context.Background(), options, kubeconfig)
}

// GetServiceNameForTenant this func will return the service name for a given tenant name.
func GetServiceNameForTenant(client *kubernetes.Clientset, nameSpace, tenantDomainName string) (string, error) {
	services, err := client.CoreV1().Services(nameSpace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return "", err
	}
	for _, svc := range services.Items {
		svcTenantName := strings.Split(svc.Name, "-")
		if svcTenantName[0] == tenantDomainName && svcTenantName[1] == "sg" {
			return svc.Name, nil
		}
	}
	return "", errors.New("No service was find to: " + tenantDomainName)
}

// GetSgPodsNames This func will return a slice of all the pod names related to group or to a node or to both.
func GetSgPodsNames(client *kubernetes.Clientset, sgGroupNum, nameSpace, nodeName, row string) ([]string, error) {
	pods, err := client.CoreV1().Pods(nameSpace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	podsNames := make([]string, 0)
	nodeNameToCheck := strings.Split(nodeName, "-")
	for _, pod := range pods.Items {
		podNameToCheck := strings.Split(pod.Name, "-")
		if sgGroupNum != "" {
			if podNameToCheck[0] == "sg" && podNameToCheck[1] == sgGroupNum {
				if nodeName != "" && len(nodeNameToCheck) > 2 {
					if podNameToCheck[2] == nodeNameToCheck[0] && podNameToCheck[3] == nodeNameToCheck[1] {
						podsNames = append(podsNames, pod.Name)
					}
				} else {
					podsNames = append(podsNames, pod.Name)
				}
			}
		} else {
			if nodeName != "" && len(nodeNameToCheck) == 2 {
				if podNameToCheck[2] == nodeNameToCheck[0] && podNameToCheck[3] == nodeNameToCheck[1] {
					podsNames = append(podsNames, pod.Name)
				}
			} else {
				podsNames = append(podsNames, pod.Name)
			}
		}
	}
	if len(podsNames) > 0 {
		return podsNames, nil
	} else {
		return nil, fmt.Errorf("No pods ware found in node %s for group %s: ", nodeName, sgGroupNum)
	}
}

// GetSecretContent this func returns the secret content of a given secret name.
func GetSecretContent(client *kubernetes.Clientset, nameSpace, secretName string) (map[string][]byte, error) {
	secret, err := client.CoreV1().Secrets(nameSpace).Get(context.TODO(), secretName, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	return secret.Data, nil
}

// ConfigK8sClient this func creates a k8s client from kubeconfig file.
func ConfigK8sClient(kubeconfig string) (*kubernetes.Clientset, error) {
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		return nil, err
	}
	if client, err := kubernetes.NewForConfig(config); err != nil {
		return nil, err
	} else {
		return client, err
	}
}
