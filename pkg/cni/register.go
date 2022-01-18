package cni

import (
	kubeovn "github.com/k8snetworkplumbingwg/kubemacpool/pkg/cni/kube-ovn"
	"github.com/k8snetworkplumbingwg/kubemacpool/pkg/cni/util"
)

func init() {
	if util.LookupEnvAsBool(kubeovn.ENV) {
		registerMutateVirtualMachineFn(kubeovn.MutateVirtualMachineFitKubeOVNFn)
	}
}
