package aws

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// This types are copied from Provisioner.

// InfrastructureConfig infrastructure configuration resource
type InfrastructureConfig struct {
	metav1.TypeMeta

	// EnableECRAccess specifies whether the IAM role policy for the worker nodes shall contain
	// permissions to access the ECR.
	// default: true
	EnableECRAccess *bool `json:"enableECRAccess,omitempty"`

	// Networks is the AWS specific network configuration (VPC, subnets, etc.)
	Networks Networks `json:"networks"`
}

// Networks holds information about the Kubernetes and infrastructure networks.
type Networks struct {
	// VPC indicates whether to use an existing VPC or create a new one.
	VPC VPC `json:"vpc"`
	// Zones belonging to the same region
	Zones []Zone `json:"zones"`
}

// Zone describes the properties of a zone
type Zone struct {
	// Name is the name for this zone.
	Name string `json:"name"`
	// Internal is the private subnet range to create (used for internal load balancers).
	Internal string `json:"internal"`
	// Public is the public subnet range to create (used for bastion and load balancers).
	Public string `json:"public"`
	// Workers isis the workers subnet range to create  (used for the VMs).
	Workers string `json:"workers"`
}

// VPC contains information about the AWS VPC and some related resources.
type VPC struct {
	// ID is the VPC id.
	ID *string `json:"id,omitempty"`
	// CIDR is the VPC CIDR.
	CIDR *string `json:"cidr,omitempty"`
}

// This types are copied from https://github.com/gardener/gardener-extensions/blob/master/controllers/provider-aws/pkg/apis/aws/types_controlplane.go

// ControlPlaneConfig contains configuration settings for the control plane.
type ControlPlaneConfig struct {
	metav1.TypeMeta

	// CloudControllerManager contains configuration settings for the cloud-controller-manager.
	CloudControllerManager *CloudControllerManagerConfig `json:"cloudControllerManager,omitempty"`
}

// CloudControllerManagerConfig contains configuration settings for the cloud-controller-manager.
type CloudControllerManagerConfig struct {
	// FeatureGates contains information about enabled feature gates.
	FeatureGates map[string]bool
}

// WorkerConfig contains configuration settings for the worker nodes.
type WorkerConfig struct {
	metav1.TypeMeta `json:",inline"`

	// InstanceMetadataOptions contains configuration for controlling access to the metadata API.
	InstanceMetadataOptions *InstanceMetadataOptions `json:"instanceMetadataOptions,omitempty"`
}

// HTTPTokensValue is a constant for HTTPTokens values.
type HTTPTokensValue string

const (
	// HTTPTokensRequired is a constant for requiring the use of tokens to access IMDS. Effectively disables access via
	// the IMDSv1 endpoints.
	HTTPTokensRequired HTTPTokensValue = "required"
	// HTTPTokensOptional that makes the use of tokens for IMDS optional. Effectively allows access via both IMDSv1 and
	// IMDSv2 endpoints.
	HTTPTokensOptional HTTPTokensValue = "optional"
)

// InstanceMetadataOptions contains configuration for controlling access to the metadata API.
type InstanceMetadataOptions struct {
	// HTTPTokens enforces the use of metadata v2 API.
	HTTPTokens *HTTPTokensValue `json:"httpTokens,omitempty"`
	// HTTPPutResponseHopLimit is the response hop limit for instance metadata requests.
	// Valid values are between 1 and 64.
	HTTPPutResponseHopLimit *int64 `json:"httpPutResponseHopLimit,omitempty"`
}