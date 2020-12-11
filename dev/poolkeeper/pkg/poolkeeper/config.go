// Copyright (c) 2020 TypeFox GmbH. All rights reserved.
// Licensed under the Gitpod Enterprise Source Code License,
// See License.enterprise.txt in the project root folder.

package poolkeeper

import "github.com/gitpod-io/gitpod/common-go/util"

// Config is the poolkeeper config
type Config struct {

	// Tasks contains all tasks poolkeeper should execute on the NodePools
	Tasks []*Task `json:"tasks"`
}

// Task is an action that PoolKeeper should perform regularly
type Task struct {
	// Interval configures how often we check and adjust the configured config in the cluster
	Interval util.Duration `json:"interval"`

	// PatchDeployment patches all deployments that run pods on this NodePool
	PatchDeployment *PatchDeploymentTask `json:"patchDeployment,omitempty"`
}

// TODO Needs a lease: https://carlosbecker.com/posts/k8s-leader-election
// // WipeDirectory struct
// type WipeDirectory struct {
// 	// Name identifies the NodePool
// 	NodeSelector string `json:"nodeSelector"`
// }
