package k8s

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"github.com/anthhub/forwarder"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"os"
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
func GetSgPodsNames(client *kubernetes.Clientset, nameSpace, sgGroupName, nodeName, row, rowFilePath string) ([]string, error) {
	pods, err := client.CoreV1().Pods(nameSpace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	podsNames := make([]string, 0)

	for _, pod := range pods.Items {
		found, err := checkIfNameIsInGroupOrInNode(pod.Name, sgGroupName, nodeName, row, rowFilePath)
		if err != nil {
			return nil, err
		}
		if found {
			podsNames = append(podsNames, pod.Name)
		}
	}
	if len(podsNames) > 0 {
		return podsNames, nil
	} else {
		return nil, fmt.Errorf("No pods ware found in node %s for group %s: ", nodeName, sgGroupName)
	}
}

// GetSgDeploymentsNames This func will return a slice of all the deployments names related to group or to a node or to both.
func GetSgDeploymentsNames(client *kubernetes.Clientset, nameSpace, sgGroupName, nodeName, row, rowFilePath string) ([]string, error) {
	deployments, err := client.AppsV1().Deployments(nameSpace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	deploymentsNames := make([]string, 0)

	for _, deployment := range deployments.Items {
		found, err := checkIfNameIsInGroupOrInNode(deployment.Name, sgGroupName, nodeName, row, rowFilePath)
		if err != nil {
			return nil, err
		}
		if found {
			deploymentsNames = append(deploymentsNames, deployment.Name)
		}
	}
	if len(deploymentsNames) > 0 {
		return deploymentsNames, nil
	} else {
		return nil, fmt.Errorf("No pods ware found in node %s for group %s: ", nodeName, sgGroupName)
	}
}

func GetNamesSplit(nodeName, sgGroupName string) ([]string, []string, error) {
	nodeNameToCheck := strings.Split(nodeName, "-")
	if sgGroupName != "" && len(nodeNameToCheck) != 2 {
		return nil, nil, fmt.Errorf("node name=%s is invalid", nodeName)
	}
	sgGroupNameToCheck := strings.Split(sgGroupName, "-")
	if sgGroupName != "" && len(sgGroupNameToCheck) != 2 {
		return nil, nil, fmt.Errorf("sg group name=%s is invalid", sgGroupName)
	}
	return nodeNameToCheck, sgGroupNameToCheck, nil
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

// GetPodNameForDeployment this func will return the pod name for a given deployment name.
func GetPodNameForDeployment(client *kubernetes.Clientset, nameSpace, deploymentName string) (string, error) {
	pods, err := client.CoreV1().Pods(nameSpace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return "", err
	}
	deploymentNameToCheck := strings.Split(deploymentName, "-")
	for _, pod := range pods.Items {
		podNameToCheck := strings.Split(pod.Name, "-")
		if deploymentNameToCheck[0] == podNameToCheck[0] && deploymentNameToCheck[1] == podNameToCheck[1] &&
			deploymentNameToCheck[2] == podNameToCheck[2] && deploymentNameToCheck[3] == podNameToCheck[4] {
			return pod.Name, nil
		}
	}
	return "", errors.New("No pod was find to deployment: " + deploymentName)
}

func getSgNameInRowFromRowFile(rowFilePath, row string) (map[string][]string, error) {
	if rowFilePath == "" {
		return nil, errors.New("rowFilePath cannot be empty")
	}
	f, err := os.OpenFile(rowFilePath, os.O_RDONLY, os.ModePerm)
	if err != nil {
		return nil, fmt.Errorf("open file error: %v", err)
	}
	defer f.Close()
	rd := bufio.NewScanner(f)
	sgInRow := map[string][]string{}
	for rd.Scan() {
		line := rd.Text()
		fmt.Println(line) // GET the line string
		lineToCheck := strings.Split(line, " ")
		if lineToCheck[0] == row {
			sgInRow[(lineToCheck[1][0 : len(lineToCheck[1])-1])] = strings.Split(lineToCheck[1][0:len(lineToCheck[1])-1], "-")
		}
	}
	return sgInRow, nil
}

func checkIfNameIsInGroupOrInNode(nameToCheck, sgGroupName, nodeName, row, rowFilePath string) (bool, error) {
	nodeNameToCheck := strings.Split(nodeName, "-")
	if nodeName != "" && len(nodeNameToCheck) != 2 {
		return false, fmt.Errorf("node name=%s is invalid node name shuold look like host-1", nodeName)
	}
	sgGroupNameToCheck := strings.Split(sgGroupName, "-")
	if sgGroupName != "" && len(sgGroupNameToCheck) != 2 {
		return false, fmt.Errorf("sg group name=%s is invalid group name should look like sg-1", sgGroupName)
	}
	nameToCheckSplit := strings.Split(nameToCheck, "-")
	if nameToCheck != "" && len(nameToCheckSplit) < 2 {
		return false, fmt.Errorf("name to check=%s is invalid", nameToCheck)
	}

	var sgNameIsRow map[string][]string
	var err error
	if row != "" {
		if sgNameIsRow, err = getSgNameInRowFromRowFile(rowFilePath, row); err != nil {
			return false, err
		}
	}

	if sgGroupName != "" {
		if nameToCheckSplit[0] == sgGroupNameToCheck[0] && nameToCheckSplit[1] == sgGroupNameToCheck[1] {
			if nodeName != "" {
				if nameToCheckSplit[2] == nodeNameToCheck[0] && nameToCheckSplit[3] == nodeNameToCheck[1] {
					return true, nil
				}
			} else {
				return true, nil
			}
		}
	} else if nodeName != "" {
		if nameToCheckSplit[2] == nodeNameToCheck[0] && nameToCheckSplit[3] == nodeNameToCheck[1] {
			if row != "" {
				if _, found := sgNameIsRow[strings.Join([]string{nameToCheckSplit[0], nameToCheckSplit[1]}, "-")]; found {
					return true, nil
				}
			} else {
				return true, nil
			}
		}
	} else {
		return false, nil
	}
	return false, nil
}
