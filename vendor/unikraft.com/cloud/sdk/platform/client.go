// This file is auto-generated. DO NOT EDIT.
// SPDX-License-Identifier: BSD-3-Clause
// Copyright (c) 2025, Unikraft GmbH.
// Licensed under the BSD-3-Clause License (the "License").
// You may not use this file except in compliance with the License.

package platform

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"unikraft.com/cloud/sdk/pkg/httpclient"
)

type Client interface {
	// Create an autoscale configuration for the specified service group given
	// its UUID.
	//
	// @param `uuid`
	// 	The UUID of the service to create a configuration for.
	// 	Mutually exclusive with name.
	//
	// @param `request`
	// 	The request body for this operation.
	//
	// @param `ropts`
	// 	Optional request modifiers.
	//
	// Performs: POST /v1/services/{uuid}/autoscale
	//
	// See: https://unikraft.com/docs/api/platform/v1/autoscale#create-autoscale-configuration-by-service-group-uuid
	CreateAutoscaleConfigurationByServiceGroupUUID(ctx context.Context, uuid string, request CreateAutoscaleConfigurationByServiceGroupUUIDRequest, ropts ...RequestOption) (*Response[CreateAutoscaleConfigurationsResponseData], error)
	// Add a new autoscale policy to an autoscale configuration given a service
	// group UUID.
	//
	// @param `uuid`
	//
	// @param `request`
	// 	The request body for this operation.
	//
	// @param `ropts`
	// 	Optional request modifiers.
	//
	// Performs: POST /v1/services/{uuid}/autoscale/policies
	//
	// See: https://unikraft.com/docs/api/platform/v1/autoscale#create-autoscale-configuration-policy
	CreateAutoscaleConfigurationPolicy(ctx context.Context, uuid string, request CreateAutoscaleConfigurationPolicyRequest, ropts ...RequestOption) (*Response[CreateAutoscaleConfigurationPolicyResponseData], error)
	// Create one or more autoscale configurations for the specified service groups
	// given their UUIDs or names.
	//
	// @param `request`
	// 	The request body for this operation.
	//
	// @param `ropts`
	// 	Optional request modifiers.
	//
	// Performs: POST /v1/services/autoscale
	//
	// See: https://unikraft.com/docs/api/platform/v1/autoscale#create-autoscale-configurations
	CreateAutoscaleConfigurations(ctx context.Context, request []CreateAutoscaleConfigurationsRequestConfiguration, ropts ...RequestOption) (*Response[CreateAutoscaleConfigurationsResponseData], error)
	// Delete one or more autoscale policies for a given service group.
	//
	// @param `uuid`
	//
	// @param `request`
	// 	The request body for this operation.
	//
	// @param `ropts`
	// 	Optional request modifiers.
	//
	// Performs: DELETE /v1/services/{uuid}/autoscale/policies
	//
	// See: https://unikraft.com/docs/api/platform/v1/autoscale#delete-autoscale-configuration-policies
	DeleteAutoscaleConfigurationPolicies(ctx context.Context, uuid string, request DeletePolicyRequest, ropts ...RequestOption) (*Response[DeleteAutoscaleConfigurationPolicyResponseData], error)
	// Delete an autoscale policy by name given the service group UUID.
	//
	// @param `uuid`
	// 	The UUID of the service group.
	//
	// @param `name`
	// 	The name of the policy to delete.
	//
	// @param `ropts`
	// 	Optional request modifiers.
	//
	// Performs: DELETE /v1/services/{uuid}/autoscale/policies/{name}
	//
	// See: https://unikraft.com/docs/api/platform/v1/autoscale#delete-autoscale-configuration-policy-by-name
	DeleteAutoscaleConfigurationPolicyByName(ctx context.Context, uuid string, name string, ropts ...RequestOption) (*Response[DeleteAutoscaleConfigurationPolicyResponseData], error)
	// Delete autoscale configuration for a given set of service groups given
	// their UUIDs or names.
	//
	// @param `request`
	// 	The request body for this operation.
	//
	// @param `ropts`
	// 	Optional request modifiers.
	//
	// Performs: DELETE /v1/services/autoscale
	//
	// See: https://unikraft.com/docs/api/platform/v1/autoscale#delete-autoscale-configurations
	DeleteAutoscaleConfigurations(ctx context.Context, request []NameOrUUID, ropts ...RequestOption) (*Response[DeleteAutoscaleConfigurationsResponseData], error)
	// Delete the autoscale configuration for the service group given its UUID.
	//
	// Unikraft Cloud will immediately drain all connections from all instances
	// that have been created by autoscale and delete the instances afterwards.
	// The draining phase is allowed to take at most `cooldown_time_ms`
	// milliseconds after which remaining connections are forcefully closed.  The
	// master instance is never deleted.  However, deleting the autoscale
	// configuration causes the master instance to start if it is stopped.
	//
	// @param `uuid`
	// 	The UUID of the service group.
	//
	// @param `ropts`
	// 	Optional request modifiers.
	//
	// Performs: DELETE /v1/services/{uuid}/autoscale
	//
	// See: https://unikraft.com/docs/api/platform/v1/autoscale#delete-autoscale-configurations-by-service-group-uuid
	DeleteAutoscaleConfigurationsByServiceGroupUUID(ctx context.Context, uuid string, ropts ...RequestOption) (*Response[DeleteAutoscaleConfigurationsResponseData], error)
	// List the autoscale policies for a given service group given its UUID.
	//
	// @param `uuid`
	//
	// @param `request`
	// 	The request body for this operation.
	//
	// @param `ropts`
	// 	Optional request modifiers.
	//
	// Performs: GET /v1/services/{uuid}/autoscale/policies
	//
	// See: https://unikraft.com/docs/api/platform/v1/autoscale#get-autoscale-configuration-policies
	GetAutoscaleConfigurationPolicies(ctx context.Context, uuid string, request GetAutoscaleConfigurationPolicyRequest, ropts ...RequestOption) (*Response[GetAutoscaleConfigurationPolicyResponseData], error)
	// Return the current state and configuration of an autoscale policy given
	// the service group UUID and the name of the policy.
	//
	// @param `uuid`
	// 	The UUID of the service group.
	//
	// @param `name`
	// 	The name of the policy to get.
	//
	// @param `ropts`
	// 	Optional request modifiers.
	//
	// Performs: GET /v1/services/{uuid}/autoscale/policies/{name}
	//
	// See: https://unikraft.com/docs/api/platform/v1/autoscale#get-autoscale-configuration-policy-by-name
	GetAutoscaleConfigurationPolicyByName(ctx context.Context, uuid string, name string, ropts ...RequestOption) (*Response[GetAutoscaleConfigurationPolicyResponseData], error)
	// Return the current states and configurations of autoscale configurations
	// for a given set of service groups given their UUIDs or names.
	//
	// @param `request`
	// 	The request body for this operation.
	//
	// @param `ropts`
	// 	Optional request modifiers.
	//
	// Performs: GET /v1/services/autoscale
	//
	// See: https://unikraft.com/docs/api/platform/v1/autoscale#get-autoscale-configurations
	GetAutoscaleConfigurations(ctx context.Context, request []NameOrUUID, ropts ...RequestOption) (*Response[GetAutoscaleConfigurationsResponseData], error)
	// Return the current states and configurations of autoscale configurations
	// given a service group UUID.
	//
	// @param `uuid`
	// 	The UUID of the service group.
	//
	// @param `ropts`
	// 	Optional request modifiers.
	//
	// Performs: GET /v1/services/{uuid}/autoscale
	//
	// See: https://unikraft.com/docs/api/platform/v1/autoscale#get-autoscale-configurations-by-service-group-uuid
	GetAutoscaleConfigurationsByServiceGroupUUID(ctx context.Context, uuid string, ropts ...RequestOption) (*Response[GetAutoscaleConfigurationsResponseData], error)
	// Upload a new certificate with the given configuration.
	//
	// @param `request`
	// 	The request body for this operation.
	//
	// @param `ropts`
	// 	Optional request modifiers.
	//
	// Performs: POST /v1/certificates
	//
	// See: https://unikraft.com/docs/api/platform/v1/certificates#create-certificate
	CreateCertificate(ctx context.Context, request CreateCertificateRequest, ropts ...RequestOption) (*Response[CreateCertificateResponseData], error)
	// Delete a specified certificate by its UUID.  After this call the UUID of
	// the certificate are no longer valid.
	//
	// @param `uuid`
	// 	The UUID of the certificate to delete.
	//
	// @param `ropts`
	// 	Optional request modifiers.
	//
	// Performs: DELETE /v1/certificates/{uuid}
	//
	// See: https://unikraft.com/docs/api/platform/v1/certificates#delete-certificate-by-uuid
	DeleteCertificateByUUID(ctx context.Context, uuid string, ropts ...RequestOption) (*Response[DeleteCertificatesResponseData], error)
	// Delete the specified certificate(s).  After this call the name of the
	// certificate(s) are no longer valid.
	//
	// @param `request`
	// 	The request body for this operation.
	//
	// @param `ropts`
	// 	Optional request modifiers.
	//
	// Performs: DELETE /v1/certificates
	//
	// See: https://unikraft.com/docs/api/platform/v1/certificates#delete-certificates
	DeleteCertificates(ctx context.Context, request []NameOrUUID, ropts ...RequestOption) (*Response[DeleteCertificatesResponseData], error)
	// Get a specified certificate by its UUID.
	//
	// @param `uuid`
	// 	The UUID of the certificate.
	//
	// @param `ropts`
	// 	Optional request modifiers.
	//
	// Performs: GET /v1/certificates/{uuid}
	//
	// See: https://unikraft.com/docs/api/platform/v1/certificates#get-certificate-by-uuid
	GetCertificateByUUID(ctx context.Context, uuid string, ropts ...RequestOption) (*Response[GetCertificatesResponseData], error)
	// Get one or many certificates with their current status and configuration.
	// It's possible to filter this list by name or UUID.
	//
	// @param `request`
	// 	The request body for this operation.
	//
	// @param `details`
	// 	Whether to include details about the certificate in the response.  By
	// 	default this is set to true, meaning that all information about the
	// 	certificate will be included in the response.  If set to false, only the
	// 	basic information about the certificate will be included, such as its name
	// 	and UUID.
	//
	// @param `ropts`
	// 	Optional request modifiers.
	//
	// Performs: GET /v1/certificates
	//
	// See: https://unikraft.com/docs/api/platform/v1/certificates#get-certificates
	GetCertificates(ctx context.Context, request []NameOrUUID, details bool, ropts ...RequestOption) (*Response[GetCertificatesResponseData], error)
	// Retrieve an image by its digest.
	//
	// @param `digest`
	//
	// @param `ropts`
	// 	Optional request modifiers.
	//
	// Performs: GET /v1/images/digest/{digest}
	//
	// See: https://unikraft.com/docs/api/platform/v1/images#get-image-by-digest
	GetImageByDigest(ctx context.Context, digest string, ropts ...RequestOption) (*Response[GetImageResponseData], error)
	// Retrieve an image by its tag.
	//
	// @param `tag`
	//
	// @param `ropts`
	// 	Optional request modifiers.
	//
	// Performs: GET /v1/images/tag/{tag}
	//
	// See: https://unikraft.com/docs/api/platform/v1/images#get-image-by-tag
	GetImageByTag(ctx context.Context, tag string, ropts ...RequestOption) (*Response[GetImageResponseData], error)
	// Create an instance in Unikraft Cloud.
	//
	// @param `request`
	// 	The request body for this operation.
	//
	// @param `ropts`
	// 	Optional request modifiers.
	//
	// Performs: POST /v1/instances
	//
	// See: https://unikraft.com/docs/api/platform/v1/instances#create-instance
	CreateInstance(ctx context.Context, request CreateInstanceRequest, ropts ...RequestOption) (*Response[CreateInstanceResponseData], error)
	// Convert one or more existing instances by their UUID(s) or name(s) into
	// template instances that can be used to create new instances.
	//
	// The existing instances must be in the `stopped` state and not have existing
	// snapshots.
	//
	// @param `request`
	// 	The request body for this operation.
	//
	// @param `ropts`
	// 	Optional request modifiers.
	//
	// Performs: POST /v1/instances/templates
	//
	// See: https://unikraft.com/docs/api/platform/v1/instances#create-template-instances
	CreateTemplateInstances(ctx context.Context, request CreateTemplateInstancesRequest, ropts ...RequestOption) (*Response[CreateTemplateInstancesResponseData], error)
	// Delete a specified instance by its UUID.  After this call the UUID of the
	// instance is no longer valid.  If the instance is currently running,
	// it is force-stopped.
	//
	// @param `uuid`
	// 	The UUID of the instance to delete.
	//
	// @param `ropts`
	// 	Optional request modifiers.
	//
	// Performs: DELETE /v1/instances/{uuid}
	//
	// See: https://unikraft.com/docs/api/platform/v1/instances#delete-instance-by-uuid
	DeleteInstanceByUUID(ctx context.Context, uuid string, ropts ...RequestOption) (*Response[DeleteInstancesResponseData], error)
	// Delete the specified instance(s) by ID(s) (name or UUID).  After this call
	// the IDs of the instances are no longer valid.  If the instances are
	// currently running, they are force-stopped.
	//
	// @param `request`
	// 	The request body for this operation.
	//
	// @param `ropts`
	// 	Optional request modifiers.
	//
	// Performs: DELETE /v1/instances
	//
	// See: https://unikraft.com/docs/api/platform/v1/instances#delete-instances
	DeleteInstances(ctx context.Context, request []NameOrUUID, ropts ...RequestOption) (*Response[DeleteInstancesResponseData], error)
	// Delete a specified template instance by its UUID.  After this call the UUID
	// of the template instance is no longer valid.
	//
	// @param `uuid`
	// 	The UUID of the template instance to delete.
	//
	// @param `ropts`
	// 	Optional request modifiers.
	//
	// Performs: DELETE /v1/instances/templates/{uuid}
	//
	// See: https://unikraft.com/docs/api/platform/v1/instances#delete-template-instance-by-uuid
	DeleteTemplateInstanceByUUID(ctx context.Context, uuid string, ropts ...RequestOption) (*Response[DeleteTemplateInstancesResponseData], error)
	// Delete the specified template instance(s) by ID(s) (name or UUID).  After
	// this call the IDs of the template instances are no longer valid.
	//
	// @param `request`
	// 	The request body for this operation.
	//
	// @param `ropts`
	// 	Optional request modifiers.
	//
	// Performs: DELETE /v1/instances/templates
	//
	// See: https://unikraft.com/docs/api/platform/v1/instances#delete-template-instances
	DeleteTemplateInstances(ctx context.Context, request []NameOrUUID, ropts ...RequestOption) (*Response[DeleteTemplateInstancesResponseData], error)
	// Get a single instance by its UUID.
	//
	// @param `uuid`
	// 	The UUID of the instance to get.
	//
	// @param `details`
	// 	Whether to include details about the instance in the response.  By default
	// 	this is set to true, meaning that all information about the instance will
	// 	be included in the response.  If set to false, only the basic information
	// 	about the instance will be included, such as its name and UUID.
	//
	// @param `ropts`
	// 	Optional request modifiers.
	//
	// Performs: GET /v1/instances/{uuid}
	//
	// See: https://unikraft.com/docs/api/platform/v1/instances#get-instance-by-uuid
	GetInstanceByUUID(ctx context.Context, uuid string, details bool, ropts ...RequestOption) (*Response[GetInstancesResponseData], error)
	// Retrieve the logs of one or more instances by ID(s) (name or UUID).
	//
	// @param `request`
	// 	The request body for this operation.
	//
	// @param `ropts`
	// 	Optional request modifiers.
	//
	// Performs: GET /v1/instances/log
	//
	// See: https://unikraft.com/docs/api/platform/v1/instances#get-instance-logs
	GetInstanceLogs(ctx context.Context, request []GetInstancesLogsRequestItem, ropts ...RequestOption) (*Response[GetInstancesLogsResponseData], error)
	// Retrieve the logs of an instance by its UUID.
	//
	// @param `uuid`
	// 	The UUID of the instance to retrieve logs for.
	//
	// @param `request`
	// 	The request body for this operation.
	//
	// @param `ropts`
	// 	Optional request modifiers.
	//
	// Performs: GET /v1/instances/{uuid}/log
	//
	// See: https://unikraft.com/docs/api/platform/v1/instances#get-instance-logs-by-uuid
	GetInstanceLogsByUUID(ctx context.Context, uuid string, request GetInstanceLogsByUUIDRequestBody, ropts ...RequestOption) (*Response[GetInstancesLogsResponseData], error)
	// Get the metrics of one or more instances by their ID(s) (name or UUID).
	//
	// @param `request`
	// 	The request body for this operation.
	//
	// @param `ropts`
	// 	Optional request modifiers.
	//
	// Performs: GET /v1/instances/metrics
	//
	// See: https://unikraft.com/docs/api/platform/v1/instances#get-instance-metrics
	GetInstanceMetrics(ctx context.Context, request []NameOrUUID, ropts ...RequestOption) (*Response[GetInstancesMetricsResponseData], error)
	// Get the metrics of an instance by its UUID.
	//
	// @param `uuid`
	// 	The UUID of the instance to retrieve metrics for.
	//
	// @param `ropts`
	// 	Optional request modifiers.
	//
	// Performs: GET /v1/instances/{uuid}/metrics
	//
	// See: https://unikraft.com/docs/api/platform/v1/instances#get-instance-metrics-by-uuid
	GetInstanceMetricsByUUID(ctx context.Context, uuid string, ropts ...RequestOption) (*Response[GetInstancesMetricsResponseData], error)
	// Get one or many instances with their current status and configuration.
	// It's possible to filter this list by ID(s) (name or UUID).
	//
	// @param `request`
	// 	The request body for this operation.
	//
	// @param `details`
	// 	Whether to include details about the instance in the response.  By default
	// 	this is set to true, meaning that all information about the instance will
	// 	be included in the response.  If set to false, only the basic information
	// 	about the instance will be included, such as its name and UUID.
	//
	// @param `ropts`
	// 	Optional request modifiers.
	//
	// Performs: GET /v1/instances
	//
	// See: https://unikraft.com/docs/api/platform/v1/instances#get-instances
	GetInstances(ctx context.Context, request []NameOrUUID, details bool, ropts ...RequestOption) (*Response[GetInstancesResponseData], error)
	// Get a single template instance by its UUID.
	//
	// @param `uuid`
	// 	The UUID of the template instance to retrieve.
	//
	// @param `details`
	// 	Whether to include details about the templates in the response.  By default
	// 	this is set to true, meaning that all information about the templates will
	// 	be included in the response.  If set to false, only the basic information
	// 	about the templates will be included, such as their name and UUID.
	//
	// @param `ropts`
	// 	Optional request modifiers.
	//
	// Performs: GET /v1/instances/templates/{uuid}
	//
	// See: https://unikraft.com/docs/api/platform/v1/instances#get-template-instance-by-uuid
	GetTemplateInstanceByUUID(ctx context.Context, uuid string, details bool, ropts ...RequestOption) (*Response[GetTemplateInstancesResponseData], error)
	// Get one or more template instances by their UUID(s) or name(s).
	//
	// @param `details`
	// 	Whether to include details about the templates in the response.  By default
	// 	this is set to true, meaning that all information about the templates will
	// 	be included in the response.  If set to false, only the basic information
	// 	about the templates will be included, such as their name and UUID.
	//
	// @param `fromUuid`
	// 	If set, the listing starts from (but does not include) the template with
	// 	the given UUID.  This is useful for pagination.
	//
	// @param `count`
	// 	The maximum number of template instances to return.  This is useful for
	// 	pagination.  If not set, all the template instances matching filters will
	// 	be returned.  When filtering by IDs, this should not be set.
	//
	// @param `tags`
	// 	A list of tags to filter the template instances by.
	//
	// @param `ropts`
	// 	Optional request modifiers.
	//
	// Performs: GET /v1/instances/templates
	//
	// See: https://unikraft.com/docs/api/platform/v1/instances#get-template-instances
	GetTemplateInstances(ctx context.Context, details bool, fromUuid string, count int32, tags []string, ropts ...RequestOption) (*Response[GetTemplateInstancesResponseData], error)
	// Start a previously stopped instance by its UUID or do nothing if the
	// instance is already running.
	//
	// @param `uuid`
	// 	The UUID of the instance to start.
	//
	// @param `ropts`
	// 	Optional request modifiers.
	//
	// Performs: PUT /v1/instances/{uuid}/start
	//
	// See: https://unikraft.com/docs/api/platform/v1/instances#start-instance-by-uuid
	StartInstanceByUUID(ctx context.Context, uuid string, ropts ...RequestOption) (*Response[StartInstancesResponseData], error)
	// Start previously stopped instances by ID(s) (name or UUID) or do
	// nothing if the instances are already running.
	//
	// @param `request`
	// 	The request body for this operation.
	//
	// @param `ropts`
	// 	Optional request modifiers.
	//
	// Performs: PUT /v1/instances/start
	//
	// See: https://unikraft.com/docs/api/platform/v1/instances#start-instances
	StartInstances(ctx context.Context, request []NameOrUUID, ropts ...RequestOption) (*Response[StartInstancesResponseData], error)
	// Stop a running instance by its UUID or do nothing if the instance is
	// already stopped.
	//
	// @param `uuid`
	// 	The UUID of the instance to stop.
	//
	// @param `force`
	// 	Whether to immediately force stop the instance.
	//
	// @param `drainTimeoutMs`
	// 	Timeout for draining connections in milliseconds.  The instance does not
	// 	receive new connections in the draining phase.  The instance is stopped
	// 	when the last connection has been closed or the timeout expired.  The
	// 	maximum timeout may vary.  Use -1 for the largest possible value.
	//
	// 	Note: This endpoint does not block.  Use the wait endpoint for the instance
	// 	to reach the stopped state.
	//
	// @param `ropts`
	// 	Optional request modifiers.
	//
	// Performs: PUT /v1/instances/{uuid}/stop
	//
	// See: https://unikraft.com/docs/api/platform/v1/instances#stop-instance-by-uuid
	StopInstanceByUUID(ctx context.Context, uuid string, force bool, drainTimeoutMs int32, ropts ...RequestOption) (*Response[StopInstancesResponseData], error)
	// Stop one or more running instance by ID(s) (name or UUID) or do
	// nothing if the instances are already stopped.
	//
	// @param `request`
	// 	The request body for this operation.
	//
	// @param `ropts`
	// 	Optional request modifiers.
	//
	// Performs: PUT /v1/instances/stop
	//
	// See: https://unikraft.com/docs/api/platform/v1/instances#stop-instances
	StopInstances(ctx context.Context, request []StopInstancesRequestItem, ropts ...RequestOption) (*Response[StopInstancesResponseData], error)
	// Update (modify) an instance by its UUID.  The instance must be in a stopped
	// state for most update operations.
	//
	// @param `uuid`
	// 	The UUID of the instance to update.
	//
	// @param `request`
	// 	The request body for this operation.
	//
	// @param `ropts`
	// 	Optional request modifiers.
	//
	// Performs: PATCH /v1/instances/{uuid}
	//
	// See: https://unikraft.com/docs/api/platform/v1/instances#update-instance-by-uuid
	UpdateInstanceByUUID(ctx context.Context, uuid string, request UpdateInstanceByUUIDRequestBody, ropts ...RequestOption) (*Response[UpdateInstancesResponseData], error)
	// Update (modify) one or more instances by ID(s) (name or UUID).  The
	// instances must be in a stopped state for most update operations.
	//
	// @param `request`
	// 	The request body for this operation.
	//
	// @param `ropts`
	// 	Optional request modifiers.
	//
	// Performs: PATCH /v1/instances
	//
	// See: https://unikraft.com/docs/api/platform/v1/instances#update-instances
	UpdateInstances(ctx context.Context, request []UpdateInstancesRequestItem, ropts ...RequestOption) (*Response[UpdateInstancesResponseData], error)
	// Update (modify) a template instance by its UUID.
	//
	// @param `uuid`
	// 	The UUID of the template instance to update.
	//
	// @param `request`
	// 	The request body for this operation.
	//
	// @param `ropts`
	// 	Optional request modifiers.
	//
	// Performs: PATCH /v1/instances/templates/{uuid}
	//
	// See: https://unikraft.com/docs/api/platform/v1/instances#update-template-instance-by-uuid
	UpdateTemplateInstanceByUUID(ctx context.Context, uuid string, request UpdateTemplateInstanceByUUIDRequestBody, ropts ...RequestOption) (*Response[UpdateTemplateInstancesResponseData], error)
	// Update (modify) one or more template instances by ID(s) (name or UUID).
	//
	// @param `request`
	// 	The request body for this operation.
	//
	// @param `ropts`
	// 	Optional request modifiers.
	//
	// Performs: PATCH /v1/instances/templates
	//
	// See: https://unikraft.com/docs/api/platform/v1/instances#update-template-instances
	UpdateTemplateInstances(ctx context.Context, request []UpdateTemplateInstancesRequestItem, ropts ...RequestOption) (*Response[UpdateTemplateInstancesResponseData], error)
	// Wait for an instance to reach a certain state, by its UUID.
	//
	// If the instance is already in the desired state, the request will return
	// immediately.  If the instance is not in the desired state, the request will
	// block until the instance reaches the desired state or the timeout is
	// reached.  If the timeout is reached, the request will fail with an error.
	// If the timeout is -1, the request will block indefinitely until the
	// instance reaches the desired state.
	//
	// @param `uuid`
	// 	The UUID of the instance to wait for.
	//
	// @param `request`
	// 	The request body for this operation.
	//
	// @param `ropts`
	// 	Optional request modifiers.
	//
	// Performs: GET /v1/instances/{uuid}/wait
	//
	// See: https://unikraft.com/docs/api/platform/v1/instances#wait-instance-by-uuid
	WaitInstanceByUUID(ctx context.Context, uuid string, request WaitInstanceByUUIDRequestBody, ropts ...RequestOption) (*Response[WaitInstancesResponseData], error)
	// Wait for one or more instances to reach certain states by ID(s)
	// (name or UUID).
	//
	// If the instances are already in the desired states, the request will return
	// immediately.  If the instances are not in the desired state, the request will
	// block until the instances reach the desired state or the timeout is
	// reached.  If the timeout is reached, the request will fail with an error.
	// If the timeout is -1, the request will block indefinitely until the
	// instances reach the desired states.
	//
	// @param `request`
	// 	The request body for this operation.
	//
	// @param `ropts`
	// 	Optional request modifiers.
	//
	// Performs: GET /v1/instances/wait
	//
	// See: https://unikraft.com/docs/api/platform/v1/instances#wait-instances
	WaitInstances(ctx context.Context, request []WaitInstancesRequestItem, ropts ...RequestOption) (*Response[WaitInstancesResponseData], error)
	// Return the status of a full-system health check of the node.
	//
	// @param `ropts`
	// 	Optional request modifiers.
	//
	// Performs: GET /v1/healthz
	//
	// See: https://unikraft.com/docs/api/platform/v1/node#healthz
	Healthz(ctx context.Context, ropts ...RequestOption) (*Response[HealthzResponseData], error)
	// Create a new service with the given configuration.
	//
	// Note that the service properties like published ports can only be defined
	// during creation.  They cannot be changed later.  Each port in a service can
	// specify a list of handlers that determine how traffic arriving at the port
	// is handled. See Connection Handlers for a complete overview.
	//
	// @param `request`
	// 	The request body for this operation.
	//
	// @param `ropts`
	// 	Optional request modifiers.
	//
	// Performs: POST /v1/services
	//
	// See: https://unikraft.com/docs/api/platform/v1/service-groups#create-service-group
	CreateServiceGroup(ctx context.Context, request CreateServiceGroupRequest, ropts ...RequestOption) (*Response[CreateServiceGroupResponseData], error)
	// Delete a specified service group by its UUID.  After this call the UUID of
	// the service group is no longer valid.
	//
	// @param `uuid`
	// 	The UUID of the service group to delete.
	//
	// @param `ropts`
	// 	Optional request modifiers.
	//
	// Performs: DELETE /v1/services/{uuid}
	//
	// See: https://unikraft.com/docs/api/platform/v1/service-groups#delete-service-group-by-uuid
	DeleteServiceGroupByUUID(ctx context.Context, uuid string, ropts ...RequestOption) (*Response[DeleteServiceGroupsResponseData], error)
	// Delete the specified service group(s).  After this call the name of the
	// service group(s) are no longer valid.
	//
	// @param `request`
	// 	The request body for this operation.
	//
	// @param `ropts`
	// 	Optional request modifiers.
	//
	// Performs: DELETE /v1/services
	//
	// See: https://unikraft.com/docs/api/platform/v1/service-groups#delete-service-groups
	DeleteServiceGroups(ctx context.Context, request []NameOrUUID, ropts ...RequestOption) (*Response[DeleteServiceGroupsResponseData], error)
	// Get a specified service group by its UUID.
	//
	// @param `uuid`
	// 	The UUID of the service group to retrieve.
	//
	// @param `details`
	// 	Whether to include details about the service group in the response.  By
	// 	default this is set to true, meaning that all information about the service
	// 	group will be included in the response.  If set to false, only the basic
	// 	information about the service group will be included, such as its name and
	// 	UUID.
	//
	// @param `ropts`
	// 	Optional request modifiers.
	//
	// Performs: GET /v1/services/{uuid}
	//
	// See: https://unikraft.com/docs/api/platform/v1/service-groups#get-service-group-by-uuid
	GetServiceGroupByUUID(ctx context.Context, uuid string, details bool, ropts ...RequestOption) (*Response[GetServiceGroupsResponseData], error)
	// Get one or many service groups with their current status and configuration.
	// It's possible to filter this list by name or UUID.
	//
	// @param `request`
	// 	The request body for this operation.
	//
	// @param `details`
	// 	Whether to include details about the service group in the response.  By
	// 	default this is set to true, meaning that all information about the service
	// 	group will be included in the response.  If set to false, only the basic
	// 	information about the service group will be included, such as its name and
	// 	UUID.
	//
	// @param `ropts`
	// 	Optional request modifiers.
	//
	// Performs: GET /v1/services
	//
	// See: https://unikraft.com/docs/api/platform/v1/service-groups#get-service-groups
	GetServiceGroups(ctx context.Context, request []NameOrUUID, details bool, ropts ...RequestOption) (*Response[GetServiceGroupsResponseData], error)
	// Update a service group by its UUID.
	//
	// @param `uuid`
	// 	The UUID of the service group to update.
	//
	// @param `request`
	// 	The request body for this operation.
	//
	// @param `ropts`
	// 	Optional request modifiers.
	//
	// Performs: PATCH /v1/services/{uuid}
	//
	// See: https://unikraft.com/docs/api/platform/v1/service-groups#update-service-group-by-uuid
	UpdateServiceGroupByUUID(ctx context.Context, uuid string, request UpdateServiceGroupByUUIDRequestBody, ropts ...RequestOption) (*Response[UpdateServiceGroupsResponseData], error)
	// Update one or more service groups.
	//
	// @param `request`
	// 	The request body for this operation.
	//
	// @param `ropts`
	// 	Optional request modifiers.
	//
	// Performs: PATCH /v1/services
	//
	// See: https://unikraft.com/docs/api/platform/v1/service-groups#update-service-groups
	UpdateServiceGroups(ctx context.Context, request []UpdateServiceGroupsRequestItem, ropts ...RequestOption) (*Response[UpdateServiceGroupsResponseData], error)
	// Create new user accounts. This will return 409 Conflict when any of the
	// requested users already existed on the target.
	//
	// @param `ropts`
	// 	Optional request modifiers.
	//
	// Performs: POST /v1/users
	//
	// See: https://unikraft.com/docs/api/platform/v1/users#add-users
	AddUsers(ctx context.Context, ropts ...RequestOption) (*Response[AddUsersResponseData], error)
	// List quota usage and limits of your user account.
	// Limits are hard limits that cannot be exceeded.
	//
	// @param `ropts`
	// 	Optional request modifiers.
	//
	// Performs: GET /v1/users/quotas
	//
	// See: https://unikraft.com/docs/api/platform/v1/users#get-user
	GetUser(ctx context.Context, ropts ...RequestOption) (*Response[QuotasResponseData], error)
	// List quota usage and limits of a user account by UUID.
	// Limits are hard limits that cannot be exceeded.
	//
	// @param `uuid`
	// 	The UUID of the user to retrieve quotas for.
	//
	// @param `ropts`
	// 	Optional request modifiers.
	//
	// Performs: GET /v1/users/{uuid}/quotas
	//
	// See: https://unikraft.com/docs/api/platform/v1/users#get-user-by-uuid
	GetUserByUUID(ctx context.Context, uuid string, ropts ...RequestOption) (*Response[QuotasResponseData], error)
	// Attach a volume by UUID to an instance so that the volume is mounted when
	// the instance starts.  The volume needs to be in `available` state and the
	// instance must be in `stopped` state.
	//
	// @param `uuid`
	// 	The UUID of the volume to attach.
	//
	// @param `request`
	// 	The request body for this operation.
	//
	// @param `ropts`
	// 	Optional request modifiers.
	//
	// Performs: PUT /v1/volumes/{uuid}/attach
	//
	// See: https://unikraft.com/docs/api/platform/v1/volumes#attach-volume-by-uuid
	AttachVolumeByUUID(ctx context.Context, uuid string, request AttachVolumeByUUIDRequestBody, ropts ...RequestOption) (*Response[AttachVolumesResponseData], error)
	// Attach one or more volumes specified by ID(s) (name or UUID) to instances
	// so that the volumes are mounted when the instances start.  The volumes need
	// to be in `available` state and the instances must be in `stopped` state.
	//
	// @param `request`
	// 	The request body for this operation.
	//
	// @param `ropts`
	// 	Optional request modifiers.
	//
	// Performs: PUT /v1/volumes/attach
	//
	// See: https://unikraft.com/docs/api/platform/v1/volumes#attach-volumes
	AttachVolumes(ctx context.Context, request AttachVolumesRequest, ropts ...RequestOption) (*Response[AttachVolumesResponseData], error)
	// Create a volume given the specified configuration parameters.
	// The volume is automatically initialized with an empty file system.
	// After initialization, the volume is in the `available` state and can be
	// attached to an instance with the `PUT /v1/volumes/attach` endpoint.
	// Note that, the size of a volume cannot be changed after creation.
	//
	// @param `request`
	// 	The request body for this operation.
	//
	// @param `ropts`
	// 	Optional request modifiers.
	//
	// Performs: POST /v1/volumes
	//
	// See: https://unikraft.com/docs/api/platform/v1/volumes#create-volume
	CreateVolume(ctx context.Context, request CreateVolumeRequest, ropts ...RequestOption) (*Response[CreateVolumeResponseData], error)
	// Delete the specified volume by its UUID.  If the volume is still attached
	// to an instance, the operation fails.  After this call, the IDs associated
	// with the volume are no longer valid.
	//
	// @param `uuid`
	// 	The UUID of the volume to delete.
	//
	// @param `ropts`
	// 	Optional request modifiers.
	//
	// Performs: DELETE /v1/volumes/{uuid}
	//
	// See: https://unikraft.com/docs/api/platform/v1/volumes#delete-volume-by-uuid
	DeleteVolumeByUUID(ctx context.Context, uuid string, ropts ...RequestOption) (*Response[DeleteVolumesResponseData], error)
	// Delete one or more volumes by their UUID(s) or name(s).  If the volumes are
	// still attached to an instance, the operation fails.  After this call, the
	// IDs associated with the volumes are no longer valid.
	//
	// @param `request`
	// 	The request body for this operation.
	//
	// @param `ropts`
	// 	Optional request modifiers.
	//
	// Performs: DELETE /v1/volumes
	//
	// See: https://unikraft.com/docs/api/platform/v1/volumes#delete-volumes
	DeleteVolumes(ctx context.Context, request []NameOrUUID, ropts ...RequestOption) (*Response[DeleteVolumesResponseData], error)
	// Detaches a volume by UUID from instances.  If no particular instance is
	// specified the volume is detached from all instances.  The instances from
	// which to detach must not have the volume mounted.  The API returns an error
	// for each instance from which it was unable to detach the volume.  If the
	// volume has been created together with an instance, detaching the volume
	// will make it persistent (i.e., it survives the deletion of the instance).
	//
	// @param `uuid`
	// 	The UUID of the volume to detach.
	//
	// @param `request`
	// 	The request body for this operation.
	//
	// @param `ropts`
	// 	Optional request modifiers.
	//
	// Performs: PUT /v1/volumes/{uuid}/detach
	//
	// See: https://unikraft.com/docs/api/platform/v1/volumes#detach-volume-by-uuid
	DetachVolumeByUUID(ctx context.Context, uuid string, request DetachVolumeByUUIDRequestBody, ropts ...RequestOption) (*Response[DetachVolumesResponseData], error)
	// Detach volumes specified by ID(s) (name or UUID) from instances.  If no
	// particular instance is specified the volume is detached from all instances.
	// The instances from which to detach must not have the volumes mounted.  The
	// API returns an error for each instance from which it was unable to detach
	// the volume.  If the volume has been created together with an instance,
	// detaching the volume will make it persistent (i.e., it survives the
	// deletion of the instance).
	//
	// @param `request`
	// 	The request body for this operation.
	//
	// @param `ropts`
	// 	Optional request modifiers.
	//
	// Performs: PUT /v1/volumes/detach
	//
	// See: https://unikraft.com/docs/api/platform/v1/volumes#detach-volumes
	DetachVolumes(ctx context.Context, request DetachVolumesRequest, ropts ...RequestOption) (*Response[DetachVolumesResponseData], error)
	// Return the current status and the configuration of a particular volume by
	// its UUID.
	//
	// @param `uuid`
	// 	The UUID of the volume to retrieve.
	//
	// @param `details`
	// 	Whether to include details about the volume in the response.  By
	// 	default this is set to true, meaning that all information about the
	// 	volume will be included in the response.  If set to false, only the
	// 	basic information about the volume will be included, such as its name
	// 	and UUID.
	//
	// @param `ropts`
	// 	Optional request modifiers.
	//
	// Performs: GET /v1/volumes/{uuid}
	//
	// See: https://unikraft.com/docs/api/platform/v1/volumes#get-volume-by-uuid
	GetVolumeByUUID(ctx context.Context, uuid string, details bool, ropts ...RequestOption) (*Response[GetVolumesResponseData], error)
	// Return the current status and the configuration of one or more volumes
	// specified by either UUID(s) or name(s).  If no identifier is provided,
	// all volumes are returned.
	//
	// @param `request`
	// 	The request body for this operation.
	//
	// @param `details`
	// 	Whether to include details about the volume in the response.  By
	// 	default this is set to true, meaning that all information about the
	// 	volume will be included in the response.  If set to false, only the
	// 	basic information about the volume will be included, such as its name
	// 	and UUID.
	//
	// @param `ropts`
	// 	Optional request modifiers.
	//
	// Performs: GET /v1/volumes
	//
	// See: https://unikraft.com/docs/api/platform/v1/volumes#get-volumes
	GetVolumes(ctx context.Context, request []NameOrUUID, details bool, ropts ...RequestOption) (*Response[GetVolumesResponseData], error)
	// Update the specified volume by its UUID.
	//
	// @param `uuid`
	// 	The UUID of the service group to update.
	//
	// @param `request`
	// 	The request body for this operation.
	//
	// @param `ropts`
	// 	Optional request modifiers.
	//
	// Performs: PATCH /v1/volumes/{uuid}
	//
	// See: https://unikraft.com/docs/api/platform/v1/volumes#update-volume-by-uuid
	UpdateVolumeByUUID(ctx context.Context, uuid string, request UpdateVolumeByUUIDRequestBody, ropts ...RequestOption) (*Response[UpdateVolumesResponseData], error)
	// Update one or more volumes specified by either UUID(s) or name(s).
	//
	// @param `request`
	// 	The request body for this operation.
	//
	// @param `ropts`
	// 	Optional request modifiers.
	//
	// Performs: PATCH /v1/volumes
	//
	// See: https://unikraft.com/docs/api/platform/v1/volumes#update-volumes
	UpdateVolumes(ctx context.Context, request []UpdateVolumesRequestItem, ropts ...RequestOption) (*Response[UpdateVolumesResponseData], error)
	// WithMetro sets the metro to use when connecting to the API.
	WithMetro(string) Client
	// WithTimeout sets the timeout when making the request.
	WithTimeout(time.Duration) Client
	// WithHTTPClient overwrites the base HTTP client.
	WithHTTPClient(httpclient.HTTPClient) Client
}

// NewClient creates a new client for the API.
func NewClient(copts ...ClientOption) Client {
	options := ClientOptions{}

	for _, opt := range copts {
		opt(&options)
	}

	if options.Token() == "" {
		options.SetToken(os.Getenv("UKC_TOKEN"))
	}

	if options.Token() == "" {
		options.SetToken(os.Getenv("UNIKRAFT_CLOUD_TOKEN"))
	}

	if options.Token() == "" {
		options.SetToken(os.Getenv("KRAFTCLOUD_TOKEN"))
	}

	if options.DefaultMetro() == "" {
		options.SetDefaultMetro(DefaultMetro)
	}

	if options.AllowInsecure() && options.HTTPClient() == nil {
		options.SetHTTPClient(httpclient.NewInsecureHTTPClient())
	}

	if options.HTTPClient() == nil {
		options.SetHTTPClient(httpclient.NewHTTPClient())
	}

	return &client{
		request: &Request{
			copts: &options,
		},
	}
}

type client struct {
	request *Request
}

// WithMetro sets the metro to use when connecting to the API.
func (c *client) WithMetro(m string) Client {
	ccpy := c.clone()
	ccpy.request = c.request.WithMetro(m)
	return ccpy
}

// WithHTTPClient overwrites the base HTTP client.
func (c *client) WithHTTPClient(hc httpclient.HTTPClient) Client {
	ccpy := c.clone()
	ccpy.request = c.request.WithHTTPClient(hc)
	return ccpy
}

// WithTimeout sets the timeout when making a request.
func (c *client) WithTimeout(to time.Duration) Client {
	ccpy := c.clone()
	ccpy.request = c.request.WithTimeout(to)
	return ccpy
}

// clone returns a shallow copy of c.
func (c *client) clone() *client {
	ccpy := *c
	return &ccpy
}

func (c *client) CreateAutoscaleConfigurationByServiceGroupUUID(ctx context.Context, uuid string, request CreateAutoscaleConfigurationByServiceGroupUUIDRequest, ropts ...RequestOption) (*Response[CreateAutoscaleConfigurationsResponseData], error) {
	requestPath := "/v1/services/{uuid}/autoscale"
	requestPath = strings.ReplaceAll(requestPath, "{uuid}", url.PathEscape(uuid))

	body, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("error marshalling request body: %w", err)
	}

	resp := &Response[CreateAutoscaleConfigurationsResponseData]{}
	if err := doRequest[CreateAutoscaleConfigurationsResponseData](ctx, c.request, http.MethodPost, requestPath, nil, bytes.NewReader(body), resp, ropts...); err != nil {
		return nil, fmt.Errorf("performing the request: %w", err)
	}
	return resp, nil
}

func (c *client) CreateAutoscaleConfigurationPolicy(ctx context.Context, uuid string, request CreateAutoscaleConfigurationPolicyRequest, ropts ...RequestOption) (*Response[CreateAutoscaleConfigurationPolicyResponseData], error) {
	requestPath := "/v1/services/{uuid}/autoscale/policies"
	requestPath = strings.ReplaceAll(requestPath, "{uuid}", url.PathEscape(uuid))

	body, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("error marshalling request body: %w", err)
	}

	resp := &Response[CreateAutoscaleConfigurationPolicyResponseData]{}
	if err := doRequest[CreateAutoscaleConfigurationPolicyResponseData](ctx, c.request, http.MethodPost, requestPath, nil, bytes.NewReader(body), resp, ropts...); err != nil {
		return nil, fmt.Errorf("performing the request: %w", err)
	}
	return resp, nil
}

func (c *client) CreateAutoscaleConfigurations(ctx context.Context, request []CreateAutoscaleConfigurationsRequestConfiguration, ropts ...RequestOption) (*Response[CreateAutoscaleConfigurationsResponseData], error) {
	requestPath := "/v1/services/autoscale"

	var body []byte
	var err error
	if request != nil {
		body, err = json.Marshal(request)
		if err != nil {
			return nil, fmt.Errorf("error marshalling request body: %w", err)
		}
	}

	resp := &Response[CreateAutoscaleConfigurationsResponseData]{}
	if err := doRequest[CreateAutoscaleConfigurationsResponseData](ctx, c.request, http.MethodPost, requestPath, nil, bytes.NewReader(body), resp, ropts...); err != nil {
		return nil, fmt.Errorf("performing the request: %w", err)
	}
	return resp, nil
}

func (c *client) DeleteAutoscaleConfigurationPolicies(ctx context.Context, uuid string, request DeletePolicyRequest, ropts ...RequestOption) (*Response[DeleteAutoscaleConfigurationPolicyResponseData], error) {
	requestPath := "/v1/services/{uuid}/autoscale/policies"
	requestPath = strings.ReplaceAll(requestPath, "{uuid}", url.PathEscape(uuid))

	body, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("error marshalling request body: %w", err)
	}

	resp := &Response[DeleteAutoscaleConfigurationPolicyResponseData]{}
	if err := doRequest[DeleteAutoscaleConfigurationPolicyResponseData](ctx, c.request, http.MethodDelete, requestPath, nil, bytes.NewReader(body), resp, ropts...); err != nil {
		return nil, fmt.Errorf("performing the request: %w", err)
	}
	return resp, nil
}

func (c *client) DeleteAutoscaleConfigurationPolicyByName(ctx context.Context, uuid string, name string, ropts ...RequestOption) (*Response[DeleteAutoscaleConfigurationPolicyResponseData], error) {
	requestPath := "/v1/services/{uuid}/autoscale/policies/{name}"
	requestPath = strings.ReplaceAll(requestPath, "{uuid}", url.PathEscape(uuid))
	requestPath = strings.ReplaceAll(requestPath, "{name}", url.PathEscape(name))

	resp := &Response[DeleteAutoscaleConfigurationPolicyResponseData]{}
	if err := doRequest[DeleteAutoscaleConfigurationPolicyResponseData](ctx, c.request, http.MethodDelete, requestPath, nil, nil, resp, ropts...); err != nil {
		return nil, fmt.Errorf("performing the request: %w", err)
	}
	return resp, nil
}

func (c *client) DeleteAutoscaleConfigurations(ctx context.Context, request []NameOrUUID, ropts ...RequestOption) (*Response[DeleteAutoscaleConfigurationsResponseData], error) {
	requestPath := "/v1/services/autoscale"

	var body []byte
	var err error
	if request != nil {
		body, err = json.Marshal(request)
		if err != nil {
			return nil, fmt.Errorf("error marshalling request body: %w", err)
		}
	}

	resp := &Response[DeleteAutoscaleConfigurationsResponseData]{}
	if err := doRequest[DeleteAutoscaleConfigurationsResponseData](ctx, c.request, http.MethodDelete, requestPath, nil, bytes.NewReader(body), resp, ropts...); err != nil {
		return nil, fmt.Errorf("performing the request: %w", err)
	}
	return resp, nil
}

func (c *client) DeleteAutoscaleConfigurationsByServiceGroupUUID(ctx context.Context, uuid string, ropts ...RequestOption) (*Response[DeleteAutoscaleConfigurationsResponseData], error) {
	requestPath := "/v1/services/{uuid}/autoscale"
	requestPath = strings.ReplaceAll(requestPath, "{uuid}", url.PathEscape(uuid))

	resp := &Response[DeleteAutoscaleConfigurationsResponseData]{}
	if err := doRequest[DeleteAutoscaleConfigurationsResponseData](ctx, c.request, http.MethodDelete, requestPath, nil, nil, resp, ropts...); err != nil {
		return nil, fmt.Errorf("performing the request: %w", err)
	}
	return resp, nil
}

func (c *client) GetAutoscaleConfigurationPolicies(ctx context.Context, uuid string, request GetAutoscaleConfigurationPolicyRequest, ropts ...RequestOption) (*Response[GetAutoscaleConfigurationPolicyResponseData], error) {
	requestPath := "/v1/services/{uuid}/autoscale/policies"
	requestPath = strings.ReplaceAll(requestPath, "{uuid}", url.PathEscape(uuid))

	body, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("error marshalling request body: %w", err)
	}

	resp := &Response[GetAutoscaleConfigurationPolicyResponseData]{}
	if err := doRequest[GetAutoscaleConfigurationPolicyResponseData](ctx, c.request, http.MethodGet, requestPath, nil, bytes.NewReader(body), resp, ropts...); err != nil {
		return nil, fmt.Errorf("performing the request: %w", err)
	}
	return resp, nil
}

func (c *client) GetAutoscaleConfigurationPolicyByName(ctx context.Context, uuid string, name string, ropts ...RequestOption) (*Response[GetAutoscaleConfigurationPolicyResponseData], error) {
	requestPath := "/v1/services/{uuid}/autoscale/policies/{name}"
	requestPath = strings.ReplaceAll(requestPath, "{uuid}", url.PathEscape(uuid))
	requestPath = strings.ReplaceAll(requestPath, "{name}", url.PathEscape(name))

	resp := &Response[GetAutoscaleConfigurationPolicyResponseData]{}
	if err := doRequest[GetAutoscaleConfigurationPolicyResponseData](ctx, c.request, http.MethodGet, requestPath, nil, nil, resp, ropts...); err != nil {
		return nil, fmt.Errorf("performing the request: %w", err)
	}
	return resp, nil
}

func (c *client) GetAutoscaleConfigurations(ctx context.Context, request []NameOrUUID, ropts ...RequestOption) (*Response[GetAutoscaleConfigurationsResponseData], error) {
	requestPath := "/v1/services/autoscale"

	var body []byte
	var err error
	if request != nil {
		body, err = json.Marshal(request)
		if err != nil {
			return nil, fmt.Errorf("error marshalling request body: %w", err)
		}
	}

	resp := &Response[GetAutoscaleConfigurationsResponseData]{}
	if err := doRequest[GetAutoscaleConfigurationsResponseData](ctx, c.request, http.MethodGet, requestPath, nil, bytes.NewReader(body), resp, ropts...); err != nil {
		return nil, fmt.Errorf("performing the request: %w", err)
	}
	return resp, nil
}

func (c *client) GetAutoscaleConfigurationsByServiceGroupUUID(ctx context.Context, uuid string, ropts ...RequestOption) (*Response[GetAutoscaleConfigurationsResponseData], error) {
	requestPath := "/v1/services/{uuid}/autoscale"
	requestPath = strings.ReplaceAll(requestPath, "{uuid}", url.PathEscape(uuid))

	resp := &Response[GetAutoscaleConfigurationsResponseData]{}
	if err := doRequest[GetAutoscaleConfigurationsResponseData](ctx, c.request, http.MethodGet, requestPath, nil, nil, resp, ropts...); err != nil {
		return nil, fmt.Errorf("performing the request: %w", err)
	}
	return resp, nil
}

func (c *client) CreateCertificate(ctx context.Context, request CreateCertificateRequest, ropts ...RequestOption) (*Response[CreateCertificateResponseData], error) {
	requestPath := "/v1/certificates"

	body, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("error marshalling request body: %w", err)
	}

	resp := &Response[CreateCertificateResponseData]{}
	if err := doRequest[CreateCertificateResponseData](ctx, c.request, http.MethodPost, requestPath, nil, bytes.NewReader(body), resp, ropts...); err != nil {
		return nil, fmt.Errorf("performing the request: %w", err)
	}
	return resp, nil
}

func (c *client) DeleteCertificateByUUID(ctx context.Context, uuid string, ropts ...RequestOption) (*Response[DeleteCertificatesResponseData], error) {
	requestPath := "/v1/certificates/{uuid}"
	requestPath = strings.ReplaceAll(requestPath, "{uuid}", url.PathEscape(uuid))

	resp := &Response[DeleteCertificatesResponseData]{}
	if err := doRequest[DeleteCertificatesResponseData](ctx, c.request, http.MethodDelete, requestPath, nil, nil, resp, ropts...); err != nil {
		return nil, fmt.Errorf("performing the request: %w", err)
	}
	return resp, nil
}

func (c *client) DeleteCertificates(ctx context.Context, request []NameOrUUID, ropts ...RequestOption) (*Response[DeleteCertificatesResponseData], error) {
	requestPath := "/v1/certificates"

	var body []byte
	var err error
	if request != nil {
		body, err = json.Marshal(request)
		if err != nil {
			return nil, fmt.Errorf("error marshalling request body: %w", err)
		}
	}

	resp := &Response[DeleteCertificatesResponseData]{}
	if err := doRequest[DeleteCertificatesResponseData](ctx, c.request, http.MethodDelete, requestPath, nil, bytes.NewReader(body), resp, ropts...); err != nil {
		return nil, fmt.Errorf("performing the request: %w", err)
	}
	return resp, nil
}

func (c *client) GetCertificateByUUID(ctx context.Context, uuid string, ropts ...RequestOption) (*Response[GetCertificatesResponseData], error) {
	requestPath := "/v1/certificates/{uuid}"
	requestPath = strings.ReplaceAll(requestPath, "{uuid}", url.PathEscape(uuid))

	resp := &Response[GetCertificatesResponseData]{}
	if err := doRequest[GetCertificatesResponseData](ctx, c.request, http.MethodGet, requestPath, nil, nil, resp, ropts...); err != nil {
		return nil, fmt.Errorf("performing the request: %w", err)
	}
	return resp, nil
}

func (c *client) GetCertificates(ctx context.Context, request []NameOrUUID, details bool, ropts ...RequestOption) (*Response[GetCertificatesResponseData], error) {
	requestPath := "/v1/certificates"

	query := make(url.Values)
	query.Add("details", fmt.Sprintf("%t", details))

	var body []byte
	var err error
	if request != nil {
		body, err = json.Marshal(request)
		if err != nil {
			return nil, fmt.Errorf("error marshalling request body: %w", err)
		}
	}

	resp := &Response[GetCertificatesResponseData]{}
	if err := doRequest[GetCertificatesResponseData](ctx, c.request, http.MethodGet, requestPath, query, bytes.NewReader(body), resp, ropts...); err != nil {
		return nil, fmt.Errorf("performing the request: %w", err)
	}
	return resp, nil
}

func (c *client) GetImageByDigest(ctx context.Context, digest string, ropts ...RequestOption) (*Response[GetImageResponseData], error) {
	requestPath := "/v1/images/digest/{digest}"
	requestPath = strings.ReplaceAll(requestPath, "{digest}", url.PathEscape(digest))

	resp := &Response[GetImageResponseData]{}
	if err := doRequest[GetImageResponseData](ctx, c.request, http.MethodGet, requestPath, nil, nil, resp, ropts...); err != nil {
		return nil, fmt.Errorf("performing the request: %w", err)
	}
	return resp, nil
}

func (c *client) GetImageByTag(ctx context.Context, tag string, ropts ...RequestOption) (*Response[GetImageResponseData], error) {
	requestPath := "/v1/images/tag/{tag}"
	requestPath = strings.ReplaceAll(requestPath, "{tag}", url.PathEscape(tag))

	resp := &Response[GetImageResponseData]{}
	if err := doRequest[GetImageResponseData](ctx, c.request, http.MethodGet, requestPath, nil, nil, resp, ropts...); err != nil {
		return nil, fmt.Errorf("performing the request: %w", err)
	}
	return resp, nil
}

func (c *client) CreateInstance(ctx context.Context, request CreateInstanceRequest, ropts ...RequestOption) (*Response[CreateInstanceResponseData], error) {
	requestPath := "/v1/instances"

	body, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("error marshalling request body: %w", err)
	}

	resp := &Response[CreateInstanceResponseData]{}
	if err := doRequest[CreateInstanceResponseData](ctx, c.request, http.MethodPost, requestPath, nil, bytes.NewReader(body), resp, ropts...); err != nil {
		return nil, fmt.Errorf("performing the request: %w", err)
	}
	return resp, nil
}

func (c *client) CreateTemplateInstances(ctx context.Context, request CreateTemplateInstancesRequest, ropts ...RequestOption) (*Response[CreateTemplateInstancesResponseData], error) {
	requestPath := "/v1/instances/templates"

	body, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("error marshalling request body: %w", err)
	}

	resp := &Response[CreateTemplateInstancesResponseData]{}
	if err := doRequest[CreateTemplateInstancesResponseData](ctx, c.request, http.MethodPost, requestPath, nil, bytes.NewReader(body), resp, ropts...); err != nil {
		return nil, fmt.Errorf("performing the request: %w", err)
	}
	return resp, nil
}

func (c *client) DeleteInstanceByUUID(ctx context.Context, uuid string, ropts ...RequestOption) (*Response[DeleteInstancesResponseData], error) {
	requestPath := "/v1/instances/{uuid}"
	requestPath = strings.ReplaceAll(requestPath, "{uuid}", url.PathEscape(uuid))

	resp := &Response[DeleteInstancesResponseData]{}
	if err := doRequest[DeleteInstancesResponseData](ctx, c.request, http.MethodDelete, requestPath, nil, nil, resp, ropts...); err != nil {
		return nil, fmt.Errorf("performing the request: %w", err)
	}
	return resp, nil
}

func (c *client) DeleteInstances(ctx context.Context, request []NameOrUUID, ropts ...RequestOption) (*Response[DeleteInstancesResponseData], error) {
	requestPath := "/v1/instances"

	var body []byte
	var err error
	if request != nil {
		body, err = json.Marshal(request)
		if err != nil {
			return nil, fmt.Errorf("error marshalling request body: %w", err)
		}
	}

	resp := &Response[DeleteInstancesResponseData]{}
	if err := doRequest[DeleteInstancesResponseData](ctx, c.request, http.MethodDelete, requestPath, nil, bytes.NewReader(body), resp, ropts...); err != nil {
		return nil, fmt.Errorf("performing the request: %w", err)
	}
	return resp, nil
}

func (c *client) DeleteTemplateInstanceByUUID(ctx context.Context, uuid string, ropts ...RequestOption) (*Response[DeleteTemplateInstancesResponseData], error) {
	requestPath := "/v1/instances/templates/{uuid}"
	requestPath = strings.ReplaceAll(requestPath, "{uuid}", url.PathEscape(uuid))

	resp := &Response[DeleteTemplateInstancesResponseData]{}
	if err := doRequest[DeleteTemplateInstancesResponseData](ctx, c.request, http.MethodDelete, requestPath, nil, nil, resp, ropts...); err != nil {
		return nil, fmt.Errorf("performing the request: %w", err)
	}
	return resp, nil
}

func (c *client) DeleteTemplateInstances(ctx context.Context, request []NameOrUUID, ropts ...RequestOption) (*Response[DeleteTemplateInstancesResponseData], error) {
	requestPath := "/v1/instances/templates"

	var body []byte
	var err error
	if request != nil {
		body, err = json.Marshal(request)
		if err != nil {
			return nil, fmt.Errorf("error marshalling request body: %w", err)
		}
	}

	resp := &Response[DeleteTemplateInstancesResponseData]{}
	if err := doRequest[DeleteTemplateInstancesResponseData](ctx, c.request, http.MethodDelete, requestPath, nil, bytes.NewReader(body), resp, ropts...); err != nil {
		return nil, fmt.Errorf("performing the request: %w", err)
	}
	return resp, nil
}

func (c *client) GetInstanceByUUID(ctx context.Context, uuid string, details bool, ropts ...RequestOption) (*Response[GetInstancesResponseData], error) {
	requestPath := "/v1/instances/{uuid}"
	requestPath = strings.ReplaceAll(requestPath, "{uuid}", url.PathEscape(uuid))

	query := make(url.Values)
	query.Add("details", fmt.Sprintf("%t", details))

	resp := &Response[GetInstancesResponseData]{}
	if err := doRequest[GetInstancesResponseData](ctx, c.request, http.MethodGet, requestPath, query, nil, resp, ropts...); err != nil {
		return nil, fmt.Errorf("performing the request: %w", err)
	}
	return resp, nil
}

func (c *client) GetInstanceLogs(ctx context.Context, request []GetInstancesLogsRequestItem, ropts ...RequestOption) (*Response[GetInstancesLogsResponseData], error) {
	requestPath := "/v1/instances/log"

	var body []byte
	var err error
	if request != nil {
		body, err = json.Marshal(request)
		if err != nil {
			return nil, fmt.Errorf("error marshalling request body: %w", err)
		}
	}

	resp := &Response[GetInstancesLogsResponseData]{}
	if err := doRequest[GetInstancesLogsResponseData](ctx, c.request, http.MethodGet, requestPath, nil, bytes.NewReader(body), resp, ropts...); err != nil {
		return nil, fmt.Errorf("performing the request: %w", err)
	}
	return resp, nil
}

func (c *client) GetInstanceLogsByUUID(ctx context.Context, uuid string, request GetInstanceLogsByUUIDRequestBody, ropts ...RequestOption) (*Response[GetInstancesLogsResponseData], error) {
	requestPath := "/v1/instances/{uuid}/log"
	requestPath = strings.ReplaceAll(requestPath, "{uuid}", url.PathEscape(uuid))

	body, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("error marshalling request body: %w", err)
	}

	resp := &Response[GetInstancesLogsResponseData]{}
	if err := doRequest[GetInstancesLogsResponseData](ctx, c.request, http.MethodGet, requestPath, nil, bytes.NewReader(body), resp, ropts...); err != nil {
		return nil, fmt.Errorf("performing the request: %w", err)
	}
	return resp, nil
}

func (c *client) GetInstanceMetrics(ctx context.Context, request []NameOrUUID, ropts ...RequestOption) (*Response[GetInstancesMetricsResponseData], error) {
	requestPath := "/v1/instances/metrics"

	var body []byte
	var err error
	if request != nil {
		body, err = json.Marshal(request)
		if err != nil {
			return nil, fmt.Errorf("error marshalling request body: %w", err)
		}
	}

	resp := &Response[GetInstancesMetricsResponseData]{}
	if err := doRequest[GetInstancesMetricsResponseData](ctx, c.request, http.MethodGet, requestPath, nil, bytes.NewReader(body), resp, ropts...); err != nil {
		return nil, fmt.Errorf("performing the request: %w", err)
	}
	return resp, nil
}

func (c *client) GetInstanceMetricsByUUID(ctx context.Context, uuid string, ropts ...RequestOption) (*Response[GetInstancesMetricsResponseData], error) {
	requestPath := "/v1/instances/{uuid}/metrics"
	requestPath = strings.ReplaceAll(requestPath, "{uuid}", url.PathEscape(uuid))

	resp := &Response[GetInstancesMetricsResponseData]{}
	if err := doRequest[GetInstancesMetricsResponseData](ctx, c.request, http.MethodGet, requestPath, nil, nil, resp, ropts...); err != nil {
		return nil, fmt.Errorf("performing the request: %w", err)
	}
	return resp, nil
}

func (c *client) GetInstances(ctx context.Context, request []NameOrUUID, details bool, ropts ...RequestOption) (*Response[GetInstancesResponseData], error) {
	requestPath := "/v1/instances"

	query := make(url.Values)
	query.Add("details", fmt.Sprintf("%t", details))

	var body []byte
	var err error
	if request != nil {
		body, err = json.Marshal(request)
		if err != nil {
			return nil, fmt.Errorf("error marshalling request body: %w", err)
		}
	}

	resp := &Response[GetInstancesResponseData]{}
	if err := doRequest[GetInstancesResponseData](ctx, c.request, http.MethodGet, requestPath, query, bytes.NewReader(body), resp, ropts...); err != nil {
		return nil, fmt.Errorf("performing the request: %w", err)
	}
	return resp, nil
}

func (c *client) GetTemplateInstanceByUUID(ctx context.Context, uuid string, details bool, ropts ...RequestOption) (*Response[GetTemplateInstancesResponseData], error) {
	requestPath := "/v1/instances/templates/{uuid}"
	requestPath = strings.ReplaceAll(requestPath, "{uuid}", url.PathEscape(uuid))

	query := make(url.Values)
	query.Add("details", fmt.Sprintf("%t", details))

	resp := &Response[GetTemplateInstancesResponseData]{}
	if err := doRequest[GetTemplateInstancesResponseData](ctx, c.request, http.MethodGet, requestPath, query, nil, resp, ropts...); err != nil {
		return nil, fmt.Errorf("performing the request: %w", err)
	}
	return resp, nil
}

func (c *client) GetTemplateInstances(ctx context.Context, details bool, fromUuid string, count int32, tags []string, ropts ...RequestOption) (*Response[GetTemplateInstancesResponseData], error) {
	requestPath := "/v1/instances/templates"

	query := make(url.Values)
	query.Add("details", fmt.Sprintf("%t", details))
	query.Add("from_uuid", fromUuid)
	query.Add("count", fmt.Sprintf("%d", count))
	for _, v := range tags {
		query.Add("tags", v)
	}

	resp := &Response[GetTemplateInstancesResponseData]{}
	if err := doRequest[GetTemplateInstancesResponseData](ctx, c.request, http.MethodGet, requestPath, query, nil, resp, ropts...); err != nil {
		return nil, fmt.Errorf("performing the request: %w", err)
	}
	return resp, nil
}

func (c *client) StartInstanceByUUID(ctx context.Context, uuid string, ropts ...RequestOption) (*Response[StartInstancesResponseData], error) {
	requestPath := "/v1/instances/{uuid}/start"
	requestPath = strings.ReplaceAll(requestPath, "{uuid}", url.PathEscape(uuid))

	resp := &Response[StartInstancesResponseData]{}
	if err := doRequest[StartInstancesResponseData](ctx, c.request, http.MethodPut, requestPath, nil, nil, resp, ropts...); err != nil {
		return nil, fmt.Errorf("performing the request: %w", err)
	}
	return resp, nil
}

func (c *client) StartInstances(ctx context.Context, request []NameOrUUID, ropts ...RequestOption) (*Response[StartInstancesResponseData], error) {
	requestPath := "/v1/instances/start"

	var body []byte
	var err error
	if request != nil {
		body, err = json.Marshal(request)
		if err != nil {
			return nil, fmt.Errorf("error marshalling request body: %w", err)
		}
	}

	resp := &Response[StartInstancesResponseData]{}
	if err := doRequest[StartInstancesResponseData](ctx, c.request, http.MethodPut, requestPath, nil, bytes.NewReader(body), resp, ropts...); err != nil {
		return nil, fmt.Errorf("performing the request: %w", err)
	}
	return resp, nil
}

func (c *client) StopInstanceByUUID(ctx context.Context, uuid string, force bool, drainTimeoutMs int32, ropts ...RequestOption) (*Response[StopInstancesResponseData], error) {
	requestPath := "/v1/instances/{uuid}/stop"
	requestPath = strings.ReplaceAll(requestPath, "{uuid}", url.PathEscape(uuid))

	query := make(url.Values)
	query.Add("force", fmt.Sprintf("%t", force))
	query.Add("drain_timeout_ms", fmt.Sprintf("%d", drainTimeoutMs))

	resp := &Response[StopInstancesResponseData]{}
	if err := doRequest[StopInstancesResponseData](ctx, c.request, http.MethodPut, requestPath, query, nil, resp, ropts...); err != nil {
		return nil, fmt.Errorf("performing the request: %w", err)
	}
	return resp, nil
}

func (c *client) StopInstances(ctx context.Context, request []StopInstancesRequestItem, ropts ...RequestOption) (*Response[StopInstancesResponseData], error) {
	requestPath := "/v1/instances/stop"

	var body []byte
	var err error
	if request != nil {
		body, err = json.Marshal(request)
		if err != nil {
			return nil, fmt.Errorf("error marshalling request body: %w", err)
		}
	}

	resp := &Response[StopInstancesResponseData]{}
	if err := doRequest[StopInstancesResponseData](ctx, c.request, http.MethodPut, requestPath, nil, bytes.NewReader(body), resp, ropts...); err != nil {
		return nil, fmt.Errorf("performing the request: %w", err)
	}
	return resp, nil
}

func (c *client) UpdateInstanceByUUID(ctx context.Context, uuid string, request UpdateInstanceByUUIDRequestBody, ropts ...RequestOption) (*Response[UpdateInstancesResponseData], error) {
	requestPath := "/v1/instances/{uuid}"
	requestPath = strings.ReplaceAll(requestPath, "{uuid}", url.PathEscape(uuid))

	body, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("error marshalling request body: %w", err)
	}

	resp := &Response[UpdateInstancesResponseData]{}
	if err := doRequest[UpdateInstancesResponseData](ctx, c.request, http.MethodPatch, requestPath, nil, bytes.NewReader(body), resp, ropts...); err != nil {
		return nil, fmt.Errorf("performing the request: %w", err)
	}
	return resp, nil
}

func (c *client) UpdateInstances(ctx context.Context, request []UpdateInstancesRequestItem, ropts ...RequestOption) (*Response[UpdateInstancesResponseData], error) {
	requestPath := "/v1/instances"

	var body []byte
	var err error
	if request != nil {
		body, err = json.Marshal(request)
		if err != nil {
			return nil, fmt.Errorf("error marshalling request body: %w", err)
		}
	}

	resp := &Response[UpdateInstancesResponseData]{}
	if err := doRequest[UpdateInstancesResponseData](ctx, c.request, http.MethodPatch, requestPath, nil, bytes.NewReader(body), resp, ropts...); err != nil {
		return nil, fmt.Errorf("performing the request: %w", err)
	}
	return resp, nil
}

func (c *client) UpdateTemplateInstanceByUUID(ctx context.Context, uuid string, request UpdateTemplateInstanceByUUIDRequestBody, ropts ...RequestOption) (*Response[UpdateTemplateInstancesResponseData], error) {
	requestPath := "/v1/instances/templates/{uuid}"
	requestPath = strings.ReplaceAll(requestPath, "{uuid}", url.PathEscape(uuid))

	body, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("error marshalling request body: %w", err)
	}

	resp := &Response[UpdateTemplateInstancesResponseData]{}
	if err := doRequest[UpdateTemplateInstancesResponseData](ctx, c.request, http.MethodPatch, requestPath, nil, bytes.NewReader(body), resp, ropts...); err != nil {
		return nil, fmt.Errorf("performing the request: %w", err)
	}
	return resp, nil
}

func (c *client) UpdateTemplateInstances(ctx context.Context, request []UpdateTemplateInstancesRequestItem, ropts ...RequestOption) (*Response[UpdateTemplateInstancesResponseData], error) {
	requestPath := "/v1/instances/templates"

	var body []byte
	var err error
	if request != nil {
		body, err = json.Marshal(request)
		if err != nil {
			return nil, fmt.Errorf("error marshalling request body: %w", err)
		}
	}

	resp := &Response[UpdateTemplateInstancesResponseData]{}
	if err := doRequest[UpdateTemplateInstancesResponseData](ctx, c.request, http.MethodPatch, requestPath, nil, bytes.NewReader(body), resp, ropts...); err != nil {
		return nil, fmt.Errorf("performing the request: %w", err)
	}
	return resp, nil
}

func (c *client) WaitInstanceByUUID(ctx context.Context, uuid string, request WaitInstanceByUUIDRequestBody, ropts ...RequestOption) (*Response[WaitInstancesResponseData], error) {
	requestPath := "/v1/instances/{uuid}/wait"
	requestPath = strings.ReplaceAll(requestPath, "{uuid}", url.PathEscape(uuid))

	body, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("error marshalling request body: %w", err)
	}

	resp := &Response[WaitInstancesResponseData]{}
	if err := doRequest[WaitInstancesResponseData](ctx, c.request, http.MethodGet, requestPath, nil, bytes.NewReader(body), resp, ropts...); err != nil {
		return nil, fmt.Errorf("performing the request: %w", err)
	}
	return resp, nil
}

func (c *client) WaitInstances(ctx context.Context, request []WaitInstancesRequestItem, ropts ...RequestOption) (*Response[WaitInstancesResponseData], error) {
	requestPath := "/v1/instances/wait"

	var body []byte
	var err error
	if request != nil {
		body, err = json.Marshal(request)
		if err != nil {
			return nil, fmt.Errorf("error marshalling request body: %w", err)
		}
	}

	resp := &Response[WaitInstancesResponseData]{}
	if err := doRequest[WaitInstancesResponseData](ctx, c.request, http.MethodGet, requestPath, nil, bytes.NewReader(body), resp, ropts...); err != nil {
		return nil, fmt.Errorf("performing the request: %w", err)
	}
	return resp, nil
}

func (c *client) Healthz(ctx context.Context, ropts ...RequestOption) (*Response[HealthzResponseData], error) {
	requestPath := "/v1/healthz"

	resp := &Response[HealthzResponseData]{}
	if err := doRequest[HealthzResponseData](ctx, c.request, http.MethodGet, requestPath, nil, nil, resp, ropts...); err != nil {
		return nil, fmt.Errorf("performing the request: %w", err)
	}
	return resp, nil
}

func (c *client) CreateServiceGroup(ctx context.Context, request CreateServiceGroupRequest, ropts ...RequestOption) (*Response[CreateServiceGroupResponseData], error) {
	requestPath := "/v1/services"

	body, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("error marshalling request body: %w", err)
	}

	resp := &Response[CreateServiceGroupResponseData]{}
	if err := doRequest[CreateServiceGroupResponseData](ctx, c.request, http.MethodPost, requestPath, nil, bytes.NewReader(body), resp, ropts...); err != nil {
		return nil, fmt.Errorf("performing the request: %w", err)
	}
	return resp, nil
}

func (c *client) DeleteServiceGroupByUUID(ctx context.Context, uuid string, ropts ...RequestOption) (*Response[DeleteServiceGroupsResponseData], error) {
	requestPath := "/v1/services/{uuid}"
	requestPath = strings.ReplaceAll(requestPath, "{uuid}", url.PathEscape(uuid))

	resp := &Response[DeleteServiceGroupsResponseData]{}
	if err := doRequest[DeleteServiceGroupsResponseData](ctx, c.request, http.MethodDelete, requestPath, nil, nil, resp, ropts...); err != nil {
		return nil, fmt.Errorf("performing the request: %w", err)
	}
	return resp, nil
}

func (c *client) DeleteServiceGroups(ctx context.Context, request []NameOrUUID, ropts ...RequestOption) (*Response[DeleteServiceGroupsResponseData], error) {
	requestPath := "/v1/services"

	var body []byte
	var err error
	if request != nil {
		body, err = json.Marshal(request)
		if err != nil {
			return nil, fmt.Errorf("error marshalling request body: %w", err)
		}
	}

	resp := &Response[DeleteServiceGroupsResponseData]{}
	if err := doRequest[DeleteServiceGroupsResponseData](ctx, c.request, http.MethodDelete, requestPath, nil, bytes.NewReader(body), resp, ropts...); err != nil {
		return nil, fmt.Errorf("performing the request: %w", err)
	}
	return resp, nil
}

func (c *client) GetServiceGroupByUUID(ctx context.Context, uuid string, details bool, ropts ...RequestOption) (*Response[GetServiceGroupsResponseData], error) {
	requestPath := "/v1/services/{uuid}"
	requestPath = strings.ReplaceAll(requestPath, "{uuid}", url.PathEscape(uuid))

	query := make(url.Values)
	query.Add("details", fmt.Sprintf("%t", details))

	resp := &Response[GetServiceGroupsResponseData]{}
	if err := doRequest[GetServiceGroupsResponseData](ctx, c.request, http.MethodGet, requestPath, query, nil, resp, ropts...); err != nil {
		return nil, fmt.Errorf("performing the request: %w", err)
	}
	return resp, nil
}

func (c *client) GetServiceGroups(ctx context.Context, request []NameOrUUID, details bool, ropts ...RequestOption) (*Response[GetServiceGroupsResponseData], error) {
	requestPath := "/v1/services"

	query := make(url.Values)
	query.Add("details", fmt.Sprintf("%t", details))

	var body []byte
	var err error
	if request != nil {
		body, err = json.Marshal(request)
		if err != nil {
			return nil, fmt.Errorf("error marshalling request body: %w", err)
		}
	}

	resp := &Response[GetServiceGroupsResponseData]{}
	if err := doRequest[GetServiceGroupsResponseData](ctx, c.request, http.MethodGet, requestPath, query, bytes.NewReader(body), resp, ropts...); err != nil {
		return nil, fmt.Errorf("performing the request: %w", err)
	}
	return resp, nil
}

func (c *client) UpdateServiceGroupByUUID(ctx context.Context, uuid string, request UpdateServiceGroupByUUIDRequestBody, ropts ...RequestOption) (*Response[UpdateServiceGroupsResponseData], error) {
	requestPath := "/v1/services/{uuid}"
	requestPath = strings.ReplaceAll(requestPath, "{uuid}", url.PathEscape(uuid))

	body, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("error marshalling request body: %w", err)
	}

	resp := &Response[UpdateServiceGroupsResponseData]{}
	if err := doRequest[UpdateServiceGroupsResponseData](ctx, c.request, http.MethodPatch, requestPath, nil, bytes.NewReader(body), resp, ropts...); err != nil {
		return nil, fmt.Errorf("performing the request: %w", err)
	}
	return resp, nil
}

func (c *client) UpdateServiceGroups(ctx context.Context, request []UpdateServiceGroupsRequestItem, ropts ...RequestOption) (*Response[UpdateServiceGroupsResponseData], error) {
	requestPath := "/v1/services"

	var body []byte
	var err error
	if request != nil {
		body, err = json.Marshal(request)
		if err != nil {
			return nil, fmt.Errorf("error marshalling request body: %w", err)
		}
	}

	resp := &Response[UpdateServiceGroupsResponseData]{}
	if err := doRequest[UpdateServiceGroupsResponseData](ctx, c.request, http.MethodPatch, requestPath, nil, bytes.NewReader(body), resp, ropts...); err != nil {
		return nil, fmt.Errorf("performing the request: %w", err)
	}
	return resp, nil
}

func (c *client) AddUsers(ctx context.Context, ropts ...RequestOption) (*Response[AddUsersResponseData], error) {
	requestPath := "/v1/users"

	resp := &Response[AddUsersResponseData]{}
	if err := doRequest[AddUsersResponseData](ctx, c.request, http.MethodPost, requestPath, nil, nil, resp, ropts...); err != nil {
		return nil, fmt.Errorf("performing the request: %w", err)
	}
	return resp, nil
}

func (c *client) GetUser(ctx context.Context, ropts ...RequestOption) (*Response[QuotasResponseData], error) {
	requestPath := "/v1/users/quotas"

	resp := &Response[QuotasResponseData]{}
	if err := doRequest[QuotasResponseData](ctx, c.request, http.MethodGet, requestPath, nil, nil, resp, ropts...); err != nil {
		return nil, fmt.Errorf("performing the request: %w", err)
	}
	return resp, nil
}

func (c *client) GetUserByUUID(ctx context.Context, uuid string, ropts ...RequestOption) (*Response[QuotasResponseData], error) {
	requestPath := "/v1/users/{uuid}/quotas"
	requestPath = strings.ReplaceAll(requestPath, "{uuid}", url.PathEscape(uuid))

	resp := &Response[QuotasResponseData]{}
	if err := doRequest[QuotasResponseData](ctx, c.request, http.MethodGet, requestPath, nil, nil, resp, ropts...); err != nil {
		return nil, fmt.Errorf("performing the request: %w", err)
	}
	return resp, nil
}

func (c *client) AttachVolumeByUUID(ctx context.Context, uuid string, request AttachVolumeByUUIDRequestBody, ropts ...RequestOption) (*Response[AttachVolumesResponseData], error) {
	requestPath := "/v1/volumes/{uuid}/attach"
	requestPath = strings.ReplaceAll(requestPath, "{uuid}", url.PathEscape(uuid))

	body, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("error marshalling request body: %w", err)
	}

	resp := &Response[AttachVolumesResponseData]{}
	if err := doRequest[AttachVolumesResponseData](ctx, c.request, http.MethodPut, requestPath, nil, bytes.NewReader(body), resp, ropts...); err != nil {
		return nil, fmt.Errorf("performing the request: %w", err)
	}
	return resp, nil
}

func (c *client) AttachVolumes(ctx context.Context, request AttachVolumesRequest, ropts ...RequestOption) (*Response[AttachVolumesResponseData], error) {
	requestPath := "/v1/volumes/attach"

	body, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("error marshalling request body: %w", err)
	}

	resp := &Response[AttachVolumesResponseData]{}
	if err := doRequest[AttachVolumesResponseData](ctx, c.request, http.MethodPut, requestPath, nil, bytes.NewReader(body), resp, ropts...); err != nil {
		return nil, fmt.Errorf("performing the request: %w", err)
	}
	return resp, nil
}

func (c *client) CreateVolume(ctx context.Context, request CreateVolumeRequest, ropts ...RequestOption) (*Response[CreateVolumeResponseData], error) {
	requestPath := "/v1/volumes"

	body, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("error marshalling request body: %w", err)
	}

	resp := &Response[CreateVolumeResponseData]{}
	if err := doRequest[CreateVolumeResponseData](ctx, c.request, http.MethodPost, requestPath, nil, bytes.NewReader(body), resp, ropts...); err != nil {
		return nil, fmt.Errorf("performing the request: %w", err)
	}
	return resp, nil
}

func (c *client) DeleteVolumeByUUID(ctx context.Context, uuid string, ropts ...RequestOption) (*Response[DeleteVolumesResponseData], error) {
	requestPath := "/v1/volumes/{uuid}"
	requestPath = strings.ReplaceAll(requestPath, "{uuid}", url.PathEscape(uuid))

	resp := &Response[DeleteVolumesResponseData]{}
	if err := doRequest[DeleteVolumesResponseData](ctx, c.request, http.MethodDelete, requestPath, nil, nil, resp, ropts...); err != nil {
		return nil, fmt.Errorf("performing the request: %w", err)
	}
	return resp, nil
}

func (c *client) DeleteVolumes(ctx context.Context, request []NameOrUUID, ropts ...RequestOption) (*Response[DeleteVolumesResponseData], error) {
	requestPath := "/v1/volumes"

	var body []byte
	var err error
	if request != nil {
		body, err = json.Marshal(request)
		if err != nil {
			return nil, fmt.Errorf("error marshalling request body: %w", err)
		}
	}

	resp := &Response[DeleteVolumesResponseData]{}
	if err := doRequest[DeleteVolumesResponseData](ctx, c.request, http.MethodDelete, requestPath, nil, bytes.NewReader(body), resp, ropts...); err != nil {
		return nil, fmt.Errorf("performing the request: %w", err)
	}
	return resp, nil
}

func (c *client) DetachVolumeByUUID(ctx context.Context, uuid string, request DetachVolumeByUUIDRequestBody, ropts ...RequestOption) (*Response[DetachVolumesResponseData], error) {
	requestPath := "/v1/volumes/{uuid}/detach"
	requestPath = strings.ReplaceAll(requestPath, "{uuid}", url.PathEscape(uuid))

	body, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("error marshalling request body: %w", err)
	}

	resp := &Response[DetachVolumesResponseData]{}
	if err := doRequest[DetachVolumesResponseData](ctx, c.request, http.MethodPut, requestPath, nil, bytes.NewReader(body), resp, ropts...); err != nil {
		return nil, fmt.Errorf("performing the request: %w", err)
	}
	return resp, nil
}

func (c *client) DetachVolumes(ctx context.Context, request DetachVolumesRequest, ropts ...RequestOption) (*Response[DetachVolumesResponseData], error) {
	requestPath := "/v1/volumes/detach"

	body, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("error marshalling request body: %w", err)
	}

	resp := &Response[DetachVolumesResponseData]{}
	if err := doRequest[DetachVolumesResponseData](ctx, c.request, http.MethodPut, requestPath, nil, bytes.NewReader(body), resp, ropts...); err != nil {
		return nil, fmt.Errorf("performing the request: %w", err)
	}
	return resp, nil
}

func (c *client) GetVolumeByUUID(ctx context.Context, uuid string, details bool, ropts ...RequestOption) (*Response[GetVolumesResponseData], error) {
	requestPath := "/v1/volumes/{uuid}"
	requestPath = strings.ReplaceAll(requestPath, "{uuid}", url.PathEscape(uuid))

	query := make(url.Values)
	query.Add("details", fmt.Sprintf("%t", details))

	resp := &Response[GetVolumesResponseData]{}
	if err := doRequest[GetVolumesResponseData](ctx, c.request, http.MethodGet, requestPath, query, nil, resp, ropts...); err != nil {
		return nil, fmt.Errorf("performing the request: %w", err)
	}
	return resp, nil
}

func (c *client) GetVolumes(ctx context.Context, request []NameOrUUID, details bool, ropts ...RequestOption) (*Response[GetVolumesResponseData], error) {
	requestPath := "/v1/volumes"

	query := make(url.Values)
	query.Add("details", fmt.Sprintf("%t", details))

	var body []byte
	var err error
	if request != nil {
		body, err = json.Marshal(request)
		if err != nil {
			return nil, fmt.Errorf("error marshalling request body: %w", err)
		}
	}

	resp := &Response[GetVolumesResponseData]{}
	if err := doRequest[GetVolumesResponseData](ctx, c.request, http.MethodGet, requestPath, query, bytes.NewReader(body), resp, ropts...); err != nil {
		return nil, fmt.Errorf("performing the request: %w", err)
	}
	return resp, nil
}

func (c *client) UpdateVolumeByUUID(ctx context.Context, uuid string, request UpdateVolumeByUUIDRequestBody, ropts ...RequestOption) (*Response[UpdateVolumesResponseData], error) {
	requestPath := "/v1/volumes/{uuid}"
	requestPath = strings.ReplaceAll(requestPath, "{uuid}", url.PathEscape(uuid))

	body, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("error marshalling request body: %w", err)
	}

	resp := &Response[UpdateVolumesResponseData]{}
	if err := doRequest[UpdateVolumesResponseData](ctx, c.request, http.MethodPatch, requestPath, nil, bytes.NewReader(body), resp, ropts...); err != nil {
		return nil, fmt.Errorf("performing the request: %w", err)
	}
	return resp, nil
}

func (c *client) UpdateVolumes(ctx context.Context, request []UpdateVolumesRequestItem, ropts ...RequestOption) (*Response[UpdateVolumesResponseData], error) {
	requestPath := "/v1/volumes"

	var body []byte
	var err error
	if request != nil {
		body, err = json.Marshal(request)
		if err != nil {
			return nil, fmt.Errorf("error marshalling request body: %w", err)
		}
	}

	resp := &Response[UpdateVolumesResponseData]{}
	if err := doRequest[UpdateVolumesResponseData](ctx, c.request, http.MethodPatch, requestPath, nil, bytes.NewReader(body), resp, ropts...); err != nil {
		return nil, fmt.Errorf("performing the request: %w", err)
	}
	return resp, nil
}
