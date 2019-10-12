package k8s

import (
	"fmt"
	"github.com/jhonnli/container-orchestration-api/api/k8s"
	"github.com/jhonnli/logs"
	coreV1 "k8s.io/api/core/v1"
	v12 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/typed/core/v1"
)

func NewNodeService() k8s.NodeInterface {
	podService := NewPodService()
	return &nodeService{
		podService: podService,
	}
}

type nodeService struct {
	podService k8s.PodInterface
}

func (ns *nodeService) getClient(env string) v1.NodeInterface {
	return K8sClient.getClientset(env).CoreV1().Nodes()
}

func (ns *nodeService) List(env string) (*coreV1.NodeList, error) {
	return ns.getClient(env).List(v12.ListOptions{})
}

func (ns *nodeService) Get(env, nodeName string) (*coreV1.Node, error) {
	return ns.getClient(env).Get(nodeName, v12.GetOptions{})
}

func (ns *nodeService) Drain(cluster, nodeName string) error {
	node, err := ns.getClient(cluster).Get(nodeName, v12.GetOptions{})
	if err != nil {
		logs.Error(err)
		return err
	}
	if !node.Spec.Unschedulable {
		err := ns.cordonOpt(cluster, nodeName, true)
		if err != nil {
			return err
		}
	}
	podList, err := ns.podService.ListByNode(cluster, nodeName)
	if err != nil {
		logs.Error(err)
		return err
	}
	podMap := ns.splitPodNameByNamespace(podList.Items)
	for namespace, podSlice := range podMap {
		ns.podService.EvictionPods(cluster, namespace, podSlice)
	}
	return err
}

func (ns nodeService) splitPodNameByNamespace(data []coreV1.Pod) map[string][]string {
	result := make(map[string][]string)
	for _, item := range data {
		ns := item.Namespace
		podListOfNS, ok := result[ns]
		if !ok {
			podListOfNS = make([]string, 0)
		}
		podListOfNS = append(podListOfNS, item.Name)
		result[ns] = podListOfNS
	}
	return result
}

func (ns *nodeService) Cordon(cluster, nodeName string) error {
	return ns.cordonOpt(cluster, nodeName, true)
}

func (ns *nodeService) Uncordon(cluster, nodeName string) error {
	return ns.cordonOpt(cluster, nodeName, false)
}

func (ns *nodeService) cordonOpt(cluster, nodeName string, unschedulable bool) error {
	node, err := ns.getClient(cluster).Get(nodeName, v12.GetOptions{})
	if err != nil {
		return err
	}
	node.Spec.Unschedulable = unschedulable
	data, err := ns.getClient(cluster).Update(node)
	fmt.Println(data)
	if err != nil {
		return err
	}
	return nil
}

func (ns *nodeService) AddNodeLabels(env, nodeName string, nodeLabels map[string]string) error {
	node, err := ns.getClient(env).Get(nodeName, v12.GetOptions{})
	if err != nil {
		return err
	}
	labels := node.Labels
	for k, v := range nodeLabels {
		labels[k] = v
	}
	node.Labels = labels
	_, err = ns.getClient(env).Update(node)
	return err
}

func (ns *nodeService) DeleteNodeLabels(env, nodeName string, nodeLabels map[string]string) error {
	node, err := ns.getClient(env).Get(nodeName, v12.GetOptions{})
	if err != nil {
		return err
	}
	labels := node.Labels
	for k := range nodeLabels {
		delete(labels, k)
	}
	node.Labels = labels
	_, err = ns.getClient(env).Update(node)
	return err
}

func (ns *nodeService) ListNodeLabels(env, labelSelector string) (*coreV1.NodeList, error) {
	return ns.getClient(env).List(v12.ListOptions{LabelSelector: labelSelector})
}
