// Copyright (c) Unikraft GmbH
// SPDX-License-Identifier: MPL-2.0

package mock

import (
	"context"

	"github.com/stretchr/testify/mock"

	"unikraft.com/cloud/sdk/platform"
)

// PlatformClient is a mock implementation with methods used by resources
type PlatformClient struct {
	mock.Mock
}

// Instance methods

func (m *PlatformClient) CreateInstance(ctx context.Context, req platform.CreateInstanceRequest, ropts ...platform.RequestOption) (*platform.Response[platform.CreateInstanceResponseData], error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*platform.Response[platform.CreateInstanceResponseData]), args.Error(1)
}

func (m *PlatformClient) GetInstanceByUUID(ctx context.Context, uuid string, details bool, ropts ...platform.RequestOption) (*platform.Response[platform.GetInstancesResponseData], error) {
	args := m.Called(ctx, uuid, details)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*platform.Response[platform.GetInstancesResponseData]), args.Error(1)
}

func (m *PlatformClient) GetInstances(ctx context.Context, request []platform.NameOrUUID, details bool, ropts ...platform.RequestOption) (*platform.Response[platform.GetInstancesResponseData], error) {
	args := m.Called(ctx, request, details)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*platform.Response[platform.GetInstancesResponseData]), args.Error(1)
}

func (m *PlatformClient) DeleteInstanceByUUID(ctx context.Context, uuid string, ropts ...platform.RequestOption) (*platform.Response[platform.DeleteInstancesResponseData], error) {
	args := m.Called(ctx, uuid)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*platform.Response[platform.DeleteInstancesResponseData]), args.Error(1)
}

// Certificate methods

func (m *PlatformClient) CreateCertificate(ctx context.Context, req platform.CreateCertificateRequest, ropts ...platform.RequestOption) (*platform.Response[platform.CreateCertificateResponseData], error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*platform.Response[platform.CreateCertificateResponseData]), args.Error(1)
}

func (m *PlatformClient) GetCertificateByUUID(ctx context.Context, uuid string, ropts ...platform.RequestOption) (*platform.Response[platform.GetCertificatesResponseData], error) {
	args := m.Called(ctx, uuid)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*platform.Response[platform.GetCertificatesResponseData]), args.Error(1)
}

func (m *PlatformClient) DeleteCertificateByUUID(ctx context.Context, uuid string, ropts ...platform.RequestOption) (*platform.Response[platform.DeleteCertificatesResponseData], error) {
	args := m.Called(ctx, uuid)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*platform.Response[platform.DeleteCertificatesResponseData]), args.Error(1)
}

// Volume methods

func (m *PlatformClient) CreateVolume(ctx context.Context, req platform.CreateVolumeRequest, ropts ...platform.RequestOption) (*platform.Response[platform.CreateVolumeResponseData], error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*platform.Response[platform.CreateVolumeResponseData]), args.Error(1)
}

func (m *PlatformClient) GetVolumeByUUID(ctx context.Context, uuid string, details bool, ropts ...platform.RequestOption) (*platform.Response[platform.GetVolumesResponseData], error) {
	args := m.Called(ctx, uuid, details)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*platform.Response[platform.GetVolumesResponseData]), args.Error(1)
}

func (m *PlatformClient) DeleteVolumeByUUID(ctx context.Context, uuid string, ropts ...platform.RequestOption) (*platform.Response[platform.DeleteVolumesResponseData], error) {
	args := m.Called(ctx, uuid)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*platform.Response[platform.DeleteVolumesResponseData]), args.Error(1)
}

func (m *PlatformClient) UpdateVolumeByUUID(ctx context.Context, uuid string, request platform.UpdateVolumeByUUIDRequestBody, ropts ...platform.RequestOption) (*platform.Response[platform.UpdateVolumesResponseData], error) {
	args := m.Called(ctx, uuid, request)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*platform.Response[platform.UpdateVolumesResponseData]), args.Error(1)
}
