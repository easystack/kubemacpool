package kubeovn

import (
	"fmt"
	"github.com/go-logr/logr"
	kubevirt "kubevirt.io/client-go/api/v1"
	"strings"
)

const ENV = "ENABLE_KUBE_OVN"

const CNI_NAME = "kube-ovn"

const MAC_ADDRESS_ANNOTATION_KEY = "ovn.kubernetes.io/mac_address"

func MutateVirtualMachineFitKubeOVNFn(virtualMachine *kubevirt.VirtualMachine, parentLogger logr.Logger) (*kubevirt.VirtualMachine, error) {
	logger := parentLogger.WithName("mutateVirtualMachinesForKubeOVN")

	defaultNetworkName := getDefaultMultusName(getVirtualMachineNetworks(virtualMachine))
	if defaultNetworkName == "" {
		return nil, nil
	}
	for _, iface := range getVirtualMachineInterfaces(virtualMachine) {
		if iface.Name != defaultNetworkName {
			continue
		}
		templateAnnotations := getVirtualMachineTemplateAnnotations(virtualMachine)
		if templateAnnotations == nil {
			templateAnnotations = map[string]string{}
		}
		templateAnnotations[MAC_ADDRESS_ANNOTATION_KEY] = iface.MacAddress
		virtualMachine.Spec.Template.ObjectMeta.Annotations = templateAnnotations
		break
	}

	logger.Info(fmt.Sprintf("successfully added annotation for %s", CNI_NAME))

	return virtualMachine, nil
}

func getDefaultMultusName(networks []kubevirt.Network) string {
	for _, network := range networks {
		if network.Multus != nil && network.Multus.Default && getNetworkName(network.Multus.NetworkName) == CNI_NAME {
			return network.Name
		}
	}
	return ""
}

func getNetworkName(fullNetworkName string) string {
	if strings.Contains(fullNetworkName, "/") {
		return strings.Split(fullNetworkName, "/")[1]
	}
	return fullNetworkName
}

func getVirtualMachineInterfaces(virtualMachine *kubevirt.VirtualMachine) []kubevirt.Interface {
	return virtualMachine.Spec.Template.Spec.Domain.Devices.Interfaces
}

func getVirtualMachineNetworks(virtualMachine *kubevirt.VirtualMachine) []kubevirt.Network {
	return virtualMachine.Spec.Template.Spec.Networks
}

func getVirtualMachineTemplateAnnotations(virtualMachine *kubevirt.VirtualMachine) map[string]string {
	return virtualMachine.Spec.Template.ObjectMeta.GetAnnotations()
}
