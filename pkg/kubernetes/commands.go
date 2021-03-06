// Copyright (c) 2017 Intel Corporation
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package kubernetes

import (
	"github.com/pkg/errors"
	"k8s.io/client-go/1.5/kubernetes"
	"k8s.io/client-go/1.5/pkg/api"
	"k8s.io/client-go/1.5/pkg/api/v1"
	"k8s.io/client-go/1.5/rest"
)

func getReadyNodes(k8sAPIAddress string) ([]v1.Node, error) {
	kubectlConfig := &rest.Config{
		Host:     k8sAPIAddress,
		Username: "",
		Password: "",
	}

	k8sClientset, err := kubernetes.NewForConfig(kubectlConfig)
	if err != nil {
		return nil, errors.Wrapf(err, "could not create new Kubernetes client on %q", k8sAPIAddress)
	}

	nodes, err := k8sClientset.Core().Nodes().List(api.ListOptions{})
	if err != nil {
		return nil, errors.Wrapf(err, "could not obtain Kubernetes node list on %q", k8sAPIAddress)
	}

	var readyNodes []v1.Node
	for _, node := range nodes.Items {
		for _, condition := range node.Status.Conditions {
			if condition.Type == v1.NodeReady && condition.Status == v1.ConditionTrue {
				readyNodes = append(readyNodes, node)
			}
		}
	}

	return readyNodes, nil
}
