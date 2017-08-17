// Code generated by protoc-gen-gogo.
// source: alecton.proto
// DO NOT EDIT!

/*
Package api is a generated protocol buffer package.

It is generated from these files:
	alecton.proto

It has these top-level messages:
	Cluster
	Image
	Target
	Application
	GetImageRequest
	ListImagesRequest
	ListImagesResponse
	ListApplicationsRequest
	ListApplicationsResponse
	GetApplicationRequest
	DeployRequest
	DeployResponse
	RollbackRequest
	RollbackResponse
	ListReleasesRequest
	ListReleasesResponse
*/
package api

import github_com_mwitkow_go_proto_validators "github.com/mwitkow/go-proto-validators"
import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"
import _ "google.golang.org/genproto/googleapis/api/annotations"
import _ "k8s.io/helm/pkg/proto/hapi/release"
import _ "github.com/mwitkow/go-proto-validators"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

func (this *Cluster) Validate() error {
	return nil
}
func (this *Image) Validate() error {
	return nil
}
func (this *Target) Validate() error {
	// Validation of proto3 map<> fields is unsupported.
	return nil
}
func (this *Application) Validate() error {
	for _, item := range this.Targets {
		if item != nil {
			if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(item); err != nil {
				return github_com_mwitkow_go_proto_validators.FieldError("Targets", err)
			}
		}
	}
	return nil
}
func (this *GetImageRequest) Validate() error {
	return nil
}
func (this *ListImagesRequest) Validate() error {
	return nil
}
func (this *ListImagesResponse) Validate() error {
	for _, item := range this.Images {
		if item != nil {
			if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(item); err != nil {
				return github_com_mwitkow_go_proto_validators.FieldError("Images", err)
			}
		}
	}
	return nil
}
func (this *ListApplicationsRequest) Validate() error {
	return nil
}
func (this *ListApplicationsResponse) Validate() error {
	for _, item := range this.Applications {
		if item != nil {
			if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(item); err != nil {
				return github_com_mwitkow_go_proto_validators.FieldError("Applications", err)
			}
		}
	}
	return nil
}
func (this *GetApplicationRequest) Validate() error {
	return nil
}
func (this *DeployRequest) Validate() error {
	return nil
}
func (this *DeployResponse) Validate() error {
	if this.Release != nil {
		if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(this.Release); err != nil {
			return github_com_mwitkow_go_proto_validators.FieldError("Release", err)
		}
	}
	return nil
}
func (this *RollbackRequest) Validate() error {
	return nil
}
func (this *RollbackResponse) Validate() error {
	if this.Release != nil {
		if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(this.Release); err != nil {
			return github_com_mwitkow_go_proto_validators.FieldError("Release", err)
		}
	}
	return nil
}
func (this *ListReleasesRequest) Validate() error {
	return nil
}
func (this *ListReleasesResponse) Validate() error {
	for _, item := range this.Releases {
		if item != nil {
			if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(item); err != nil {
				return github_com_mwitkow_go_proto_validators.FieldError("Releases", err)
			}
		}
	}
	return nil
}
