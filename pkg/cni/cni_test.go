package cni

import (
	kubeovn "github.com/k8snetworkplumbingwg/kubemacpool/pkg/cni/kube-ovn"
	kubevirt "kubevirt.io/client-go/api/v1"
	ctrl "sigs.k8s.io/controller-runtime"
	"testing"
)

func makeVirtualMachineWithDefaultPodNetwork() *kubevirt.VirtualMachine {
	return &kubevirt.VirtualMachine{
		Spec: kubevirt.VirtualMachineSpec{
			Template: &kubevirt.VirtualMachineInstanceTemplateSpec{
				Spec: kubevirt.VirtualMachineInstanceSpec{
					Domain: kubevirt.DomainSpec{
						Devices: kubevirt.Devices{
							Interfaces: []kubevirt.Interface{
								{
									Name: "net",
									InterfaceBindingMethod: kubevirt.InterfaceBindingMethod{
										Masquerade: &kubevirt.InterfaceMasquerade{},
									},
								},
							},
						},
					},
					Networks: []kubevirt.Network{
						{
							Name: "net",
							NetworkSource: kubevirt.NetworkSource{
								Pod: &kubevirt.PodNetwork{},
							},
						},
					},
				},
			},
		},
	}
}

func makeVirtualMachineWithMultus() *kubevirt.VirtualMachine {
	return &kubevirt.VirtualMachine{
		Spec: kubevirt.VirtualMachineSpec{
			Template: &kubevirt.VirtualMachineInstanceTemplateSpec{
				Spec: kubevirt.VirtualMachineInstanceSpec{
					Domain: kubevirt.DomainSpec{
						Devices: kubevirt.Devices{
							Interfaces: []kubevirt.Interface{
								{
									Name: "net",
									MacAddress: "00:00:00:63:D2:F9",
									InterfaceBindingMethod: kubevirt.InterfaceBindingMethod{
										Bridge: &kubevirt.InterfaceBridge{},
									},
								},
							},
						},
					},
					Networks: []kubevirt.Network{
						{
							Name: "net",
							NetworkSource: kubevirt.NetworkSource{
								Multus: &kubevirt.MultusNetwork{
									Default: true,
									NetworkName: "secure-container/kube-ovn",
								},
							},
						},
					},
				},
			},
		},
	}
}

func TestDefaultPodNetWork(t *testing.T) {
	vm := makeVirtualMachineWithDefaultPodNetwork()
	log := ctrl.Log.WithName("TestFitCNI")

	factory := MutateVirtualMachineFactory{}
	if err := factory.Run(vm, log); err != nil {
		t.Errorf("Test failed: %v", err)
	}

	templateAnnotations := vm.Spec.Template.ObjectMeta.GetAnnotations()
	if _, ok := templateAnnotations[kubeovn.MAC_ADDRESS_ANNOTATION_KEY]; ok {
		t.Errorf("Test failed")
	}
}

func TestFitCNI(t *testing.T) {
	vm := makeVirtualMachineWithMultus()
	log := ctrl.Log.WithName("TestFitCNI")

	factory := MutateVirtualMachineFactory{}
	if err := factory.Run(vm, log); err != nil {
		t.Errorf("Test failed: %v", err)
	}

	templateAnnotations := vm.Spec.Template.ObjectMeta.GetAnnotations()
	if _, ok := templateAnnotations[kubeovn.MAC_ADDRESS_ANNOTATION_KEY]; !ok {
		t.Errorf("Test failed")
	}
}
