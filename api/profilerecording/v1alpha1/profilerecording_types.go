/*
Copyright 2020 The Kubernetes Authors.

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

package v1alpha1

import (
	"fmt"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"sigs.k8s.io/security-profiles-operator/internal/pkg/config"
)

type ProfileRecordingKind string

const (
	ProfileRecordingKindSeccompProfile ProfileRecordingKind = "SeccompProfile"
	ProfileRecordingKindSelinuxProfile ProfileRecordingKind = "SelinuxProfile"
)

type ProfileRecorder string

const (
	ProfileRecorderHook ProfileRecorder = "hook"
	ProfileRecorderLogs ProfileRecorder = "logs"
	ProfileRecorderBpf  ProfileRecorder = "bpf"
)

// ProfileRecordingSpec defines the desired state of ProfileRecording.
type ProfileRecordingSpec struct {
	// Kind of object to be recorded.
	// +kubebuilder:validation:Enum=SeccompProfile;SelinuxProfile
	Kind ProfileRecordingKind `json:"kind"`

	// Recorder to be used.
	// +kubebuilder:validation:Enum=bpf;hook;logs
	Recorder ProfileRecorder `json:"recorder"`

	// PodSelector selects the pods to record. This field follows standard
	// label selector semantics. An empty podSelector matches all pods in this
	// namespace.
	PodSelector metav1.LabelSelector `json:"podSelector"`

	// Containers is a set of containers to record. This allows to select
	// only specific containers to record instead of all containers present
	// in the pod.
	// +optional
	Containers []string `json:"containers,omitempty"`
}

// ProfileRecordingStatus contains status of the ProfileRecording.
type ProfileRecordingStatus struct {
	ActiveWorkloads []string `json:"activeWorkloads,omitempty"`
}

// +kubebuilder:object:root=true

// ProfileRecording is the Schema for the profilerecordings API.
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="PodSelector",type=string,priority=10,JSONPath=`.spec.podSelector`
type ProfileRecording struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ProfileRecordingSpec   `json:"spec,omitempty"`
	Status ProfileRecordingStatus `json:"status,omitempty"`
}

func (pr *ProfileRecording) CtrAnnotation(replica, ctrName string) (key, value string, err error) {
	ctrReplicaName := ctrName
	if replica != "" {
		ctrReplicaName += replica
	}

	switch pr.Spec.Kind {
	case ProfileRecordingKindSeccompProfile:
		return pr.ctrAnnotationSeccomp(ctrReplicaName, ctrName)
	case ProfileRecordingKindSelinuxProfile:
		return pr.ctrAnnotationSelinux(ctrReplicaName, ctrName)
	}

	return "", "", fmt.Errorf(
		"invalid kind: %s", pr.Spec.Kind,
	)
}

func (pr *ProfileRecording) IsKindSupported() bool {
	switch pr.Spec.Kind {
	case ProfileRecordingKindSelinuxProfile, ProfileRecordingKindSeccompProfile:
		return true
	}
	return false
}

func (pr *ProfileRecording) ctrAnnotationSeccomp(ctrReplicaName, ctrName string) (key, value string, err error) {
	var annotationPrefix string

	switch pr.Spec.Recorder {
	case ProfileRecorderHook:
		annotationPrefix = config.SeccompProfileRecordHookAnnotationKey
		value = fmt.Sprintf(
			"of:%s/%s-%s-%d.json",
			config.ProfileRecordingOutputPath,
			pr.GetName(),
			ctrReplicaName,
			time.Now().Unix(),
		)

	case ProfileRecorderLogs:
		annotationPrefix = config.SeccompProfileRecordLogsAnnotationKey
		value = fmt.Sprintf(
			"%s-%s-%d",
			pr.GetName(),
			ctrReplicaName,
			time.Now().Unix(),
		)

	case ProfileRecorderBpf:
		annotationPrefix = config.SeccompProfileRecordBpfAnnotationKey
		value = fmt.Sprintf(
			"%s-%s-%d",
			pr.GetName(),
			ctrReplicaName,
			time.Now().Unix(),
		)

	default:
		return "", "", fmt.Errorf(
			"invalid recorder: %s", pr.Spec.Recorder,
		)
	}

	key = annotationPrefix + ctrName
	return key, value, err
}

func (pr *ProfileRecording) ctrAnnotationSelinux(ctrReplicaName, ctrName string) (key, value string, err error) {
	var annotationPrefix string

	switch pr.Spec.Recorder {
	case ProfileRecorderLogs:
		annotationPrefix = config.SelinuxProfileRecordLogsAnnotationKey
		value = fmt.Sprintf(
			"%s-%s-%d",
			pr.GetName(),
			ctrReplicaName,
			time.Now().Unix(),
		)

	case ProfileRecorderHook:
	case ProfileRecorderBpf:
	default:
		return "", "", fmt.Errorf(
			"invalid recorder: %s, only %s is supported", pr.Spec.Recorder, ProfileRecorderLogs,
		)
	}

	key = annotationPrefix + ctrName
	return
}

// +kubebuilder:object:root=true

// ProfileRecordingList contains a list of ProfileRecording.
type ProfileRecordingList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ProfileRecording `json:"items"`
}

func init() { //nolint:gochecknoinits // required to init the scheme
	SchemeBuilder.Register(&ProfileRecording{}, &ProfileRecordingList{})
}
