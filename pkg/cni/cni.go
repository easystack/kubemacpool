package cni

import (
	"github.com/go-logr/logr"
	utilerrors "k8s.io/apimachinery/pkg/util/errors"
	kubevirt "kubevirt.io/client-go/api/v1"
)

var mutateVirtualMachineFns = make([]MutateVirtualMachineFn, 0)

// TODO: MutatePodFn if need

// If there is no need to mutate VirtualMachine, return nil.
type MutateVirtualMachineFn func(virtualMachine *kubevirt.VirtualMachine, parentLogger logr.Logger) (*kubevirt.VirtualMachine, error)

func registerMutateVirtualMachineFn(fn MutateVirtualMachineFn) {
	mutateVirtualMachineFns = append(mutateVirtualMachineFns, fn)
}

// TODO: MutatePodFactory if need

type MutateVirtualMachineFactory struct {}

func (f *MutateVirtualMachineFactory) Run(virtualMachine *kubevirt.VirtualMachine, parentLogger logr.Logger) error {
	errList := make([]error, 0)
	for _, fn := range mutateVirtualMachineFns {
		copyVM := virtualMachine.DeepCopy()
		vm, err := fn(copyVM, parentLogger)
		if err != nil {
			errList = append(errList, err)
			continue
		}
		if vm != nil {
			vm.DeepCopyInto(virtualMachine)
		}
	}
	return utilerrors.NewAggregate(errList)
}
