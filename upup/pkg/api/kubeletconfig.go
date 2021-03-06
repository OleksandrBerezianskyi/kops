package api

import "k8s.io/kops/upup/pkg/fi/utils"

// NodeLabels are defined in the InstanceGroup, but set flags on the kubelet config.
// We have a conflict here: on the one hand we want an easy to use abstract specification
// for the cluster, on the other hand we don't want two fields that do the same thing.
// So we make the logic for combining a KubeletConfig part of our core logic.
// NodeLabels are set on the instanceGroup.  We might allow specification of them on the kubelet
// config as well, but for now the precedence is not fully specified.
// (Today, NodeLabels on the InstanceGroup are merged in to NodeLabels on the KubeletConfig in the Cluster).
// In future, we will likely deprecate KubeletConfig in the Cluster, and move it into componentconfig,
// once that is part of core k8s.

// BuildKubeletConfigSpec returns the kubeletconfig for the specified instanceGroup
func BuildKubeletConfigSpec(cluster *Cluster, instanceGroup *InstanceGroup) (*KubeletConfigSpec, error) {
	// Merge KubeletConfig for NodeLabels
	c := &KubeletConfigSpec{}
	if instanceGroup.Spec.Role == InstanceGroupRoleMaster {
		utils.JsonMergeStruct(c, cluster.Spec.MasterKubelet)
	} else {
		utils.JsonMergeStruct(c, cluster.Spec.Kubelet)
	}

	for k, v := range instanceGroup.Spec.NodeLabels {
		if c.NodeLabels == nil {
			c.NodeLabels = make(map[string]string)
		}
		c.NodeLabels[k] = v
	}

	return c, nil
}
