// Copyright (c) 2020 TypeFox GmbH. All rights reserved.
// Licensed under the Gitpod Enterprise Source Code License,
// See License.enterprise.txt in the project root folder.

package poolkeeper

import (
	"time"

	// corev1 "k8s.io/api/core/v1"
	// metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	// types "k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes"

	log "github.com/sirupsen/logrus"
)

// PoolKeeper is the entity responsiple to perform the configures actions per NodePool
type PoolKeeper struct {
	Clientset *kubernetes.Clientset
	Config    *Config

	stop chan struct{}
	done chan struct{}
}

// NewPoolKeeper creates a new PoolKeeper instance
func NewPoolKeeper(clientset *kubernetes.Clientset, config *Config) *PoolKeeper {
	return &PoolKeeper{
		Clientset: clientset,
		Config:    config,
	}
}

// Start starts the PoolKeeper and is meant to be run in a goroutine
func (pk *PoolKeeper) Start() {
	defer func() {
		pk.done <- struct{}{}
	}()

	ticker := time.NewTicker(time.Duration(pk.Config.Interval))

	for {
		pk.performTasks()

		select {
		case <-pk.stop:
			return
		case <-ticker.C:
			continue
		}
	}
}

func (pk *PoolKeeper) performTasks() {
	if len(pk.Config.Tasks) == 0 {
		log.Info("no nodepools configured, doing nothing.")
		return
	}

	for _, task := range pk.Config.Tasks {
		if task.PatchDeployment != nil {
			task.PatchDeployment.run(pk.Clientset)
		}
	}
}

// Stop stops PoolKeeper and waits until is has done so
func (pk *PoolKeeper) Stop() {
	pk.stop <- struct{}{}
	<-pk.done
}

// keepOut
func keepOut() {

}
