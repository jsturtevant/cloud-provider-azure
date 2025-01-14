/*
Copyright 2021 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package options

import (
	"fmt"
	"strings"

	"github.com/spf13/pflag"

	nodeipamconfig "sigs.k8s.io/cloud-provider-azure/pkg/nodeipam/config"
)

const (
	// DefaultNodeMaskCIDRIPv4 is default mask size for IPv4 node cidr
	DefaultNodeMaskCIDRIPv4 = 24
	// DefaultNodeMaskCIDRIPv6 is default mask size for IPv6 node cidr
	DefaultNodeMaskCIDRIPv6 = 64
	// DefaultNodeCIDRMaskSize is the default mask size for node cidr
	DefaultNodeCIDRMaskSize = 24
)

// NodeIPAMControllerOptions holds the NodeIpamController options.
type NodeIPAMControllerOptions struct {
	*nodeipamconfig.NodeIPAMControllerConfiguration
}

// AddFlags adds flags related to NodeIpamController for controller manager to the specified FlagSet.
func (o *NodeIPAMControllerOptions) AddFlags(fs *pflag.FlagSet) {
	if o == nil {
		return
	}
	fs.StringVar(&o.ServiceCIDR, "service-cluster-ip-range", "", "CIDR Range for Services in cluster. Requires --allocate-node-cidrs to be true")
	fs.Int32Var(&o.NodeCIDRMaskSize, "node-cidr-mask-size", DefaultNodeCIDRMaskSize, "Mask size for node cidr in cluster. Default is 24 for IPv4 and 64 for IPv6.")
	fs.Int32Var(&o.NodeCIDRMaskSizeIPv4, "node-cidr-mask-size-ipv4", 0, "Mask size for IPv4 node cidr in dual-stack cluster. Default is 24.")
	fs.Int32Var(&o.NodeCIDRMaskSizeIPv6, "node-cidr-mask-size-ipv6", 0, "Mask size for IPv6 node cidr in dual-stack cluster. Default is 64.")
}

// ApplyTo fills up NodeIpamController config with options.
func (o *NodeIPAMControllerOptions) ApplyTo(cfg *nodeipamconfig.NodeIPAMControllerConfiguration) error {
	if o == nil {
		return nil
	}

	// split the cidrs list and assign primary and secondary
	serviceCIDRList := strings.Split(o.ServiceCIDR, ",")
	if len(serviceCIDRList) > 0 {
		cfg.ServiceCIDR = serviceCIDRList[0]
	}
	if len(serviceCIDRList) > 1 {
		cfg.SecondaryServiceCIDR = serviceCIDRList[1]
	}

	cfg.NodeCIDRMaskSize = o.NodeCIDRMaskSize
	cfg.NodeCIDRMaskSizeIPv4 = o.NodeCIDRMaskSizeIPv4
	cfg.NodeCIDRMaskSizeIPv6 = o.NodeCIDRMaskSizeIPv6

	return nil
}

// Validate checks validation of NodeIPAMControllerOptions.
func (o *NodeIPAMControllerOptions) Validate() []error {
	if o == nil {
		return nil
	}
	errs := make([]error, 0)

	serviceCIDRList := strings.Split(o.ServiceCIDR, ",")
	if len(serviceCIDRList) > 2 {
		errs = append(errs, fmt.Errorf("--service-cluster-ip-range can not contain more than two entries"))
	}

	return errs
}

func defaultNodeIPAMControllerOptions() *NodeIPAMControllerOptions {
	return &NodeIPAMControllerOptions{
		&nodeipamconfig.NodeIPAMControllerConfiguration{
			ServiceCIDR:          "",
			NodeCIDRMaskSize:     DefaultNodeCIDRMaskSize,
			NodeCIDRMaskSizeIPv4: 0,
			NodeCIDRMaskSizeIPv6: 0,
		},
	}
}
