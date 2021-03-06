/*
Licensed to the Apache Software Foundation (ASF) under one or more
contributor license agreements.  See the NOTICE file distributed with
this work for additional information regarding copyright ownership.
The ASF licenses this file to You under the Apache License, Version 2.0
(the "License"); you may not use this file except in compliance with
the License.  You may obtain a copy of the License at

   http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package build

import (
	"context"

	corev1 "k8s.io/api/core/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"

	"github.com/apache/camel-k/pkg/apis/camel/v1alpha1"
)

// NewMonitorPodAction creates a new monitor action for scheduled pod
func NewMonitorPodAction() Action {
	return &monitorPodAction{}
}

type monitorPodAction struct {
	baseAction
}

// Name returns a common name of the action
func (action *monitorPodAction) Name() string {
	return "monitor-pod"
}

// CanHandle tells whether this action can handle the build
func (action *monitorPodAction) CanHandle(build *v1alpha1.Build) bool {
	return (build.Status.Phase == v1alpha1.BuildPhasePending ||
		build.Status.Phase == v1alpha1.BuildPhaseRunning) &&
		build.Spec.Platform.Build.BuildStrategy == v1alpha1.IntegrationPlatformBuildStrategyPod
}

// Handle handles the builds
func (action *monitorPodAction) Handle(ctx context.Context, build *v1alpha1.Build) (*v1alpha1.Build, error) {
	// Get the build pod
	pod := &corev1.Pod{}
	err := action.client.Get(ctx, types.NamespacedName{Namespace: build.Namespace, Name: buildPodName(build.Spec.Meta)}, pod)
	if err != nil {
		if k8serrors.IsNotFound(err) {
			// Let's reschedule the build
			build.Status.Phase = v1alpha1.BuildPhaseScheduling
		} else {
			return nil, err
		}
	}

	var buildPhase v1alpha1.BuildPhase

	switch pod.Status.Phase {
	case corev1.PodSucceeded:
		buildPhase = v1alpha1.BuildPhaseSucceeded
	case corev1.PodFailed:
		buildPhase = v1alpha1.BuildPhaseFailed
	}

	if build.Status.Phase == buildPhase {
		// Status is already up-to-date
		return nil, nil
	}

	build.Status.Phase = buildPhase

	return build, nil
}
