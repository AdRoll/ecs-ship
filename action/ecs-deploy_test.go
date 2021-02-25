package action

import (
	"errors"
	"testing"
	"time"

	mock "github.com/adroll/ecs-ship/action/mock"
	"github.com/adroll/ecs-ship/ecs"

	ecssdk "github.com/aws/aws-sdk-go/service/ecs"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_ECSDeploy_NoArgs(t *testing.T) {
	err := ECSDeploy("", "", nil, time.Nanosecond, nil)
	assert.Error(t, err)
}

func Test_ECSDeploy_TruePath(t *testing.T) {
	const clusterName = "clusterA"
	const serviceName = "serviceA"

	timeout := time.Nanosecond
	service := &ecssdk.Service{}
	oldTask := &ecssdk.TaskDefinition{}
	taskCopy := &ecssdk.RegisterTaskDefinitionInput{}
	newTask := &ecssdk.RegisterTaskDefinitionInput{}
	newTaskDefinition := &ecssdk.TaskDefinition{}
	newService := &ecssdk.Service{}
	diff := &ecs.TaskConfigDiff{}

	//#region mock
	c := gomock.NewController(t)

	client := mock.NewMockECSDeployClient(c)
	client.EXPECT().GetService(clusterName, serviceName).Return(service, nil).Times(1)
	client.EXPECT().LooksGood(service).Return(true, nil).Times(1)
	client.EXPECT().CopyTaskDefinition(service).Return(taskCopy, oldTask, nil).Times(1)
	client.EXPECT().RegisterTaskDefinition(newTask).Return(newTaskDefinition, nil).Times(1)
	client.EXPECT().UpdateTaskDefinition(service, newTaskDefinition).Return(newService, nil).Times(1)
	client.EXPECT().WaitUntilGood(newService, timeout).Return(nil).Times(1)

	cfg := mock.NewMockECSDeployTaskConfig(c)
	cfg.EXPECT().ApplyTo(taskCopy).Return(newTask, diff).Times(1)

	//#endregion

	err := ECSDeploy(clusterName, serviceName, client, timeout, cfg)
	require.NoError(t, err)
}

func Test_ECSDeploy_RollbackPath(t *testing.T) {
	const clusterName = "clusterA"
	const serviceName = "serviceA"

	timeout := time.Nanosecond
	timeoutErr := errors.New("timeout error")
	service := &ecssdk.Service{}
	oldTask := &ecssdk.TaskDefinition{}
	taskCopy := &ecssdk.RegisterTaskDefinitionInput{}
	newTask := &ecssdk.RegisterTaskDefinitionInput{}
	newTaskDefinition := &ecssdk.TaskDefinition{}
	newService := &ecssdk.Service{}
	rolledBackService := &ecssdk.Service{}
	diff := &ecs.TaskConfigDiff{}

	//#region mock
	c := gomock.NewController(t)

	client := mock.NewMockECSDeployClient(c)
	client.EXPECT().GetService(clusterName, serviceName).Return(service, nil).Times(1)
	client.EXPECT().LooksGood(service).Return(true, nil).Times(1)
	client.EXPECT().CopyTaskDefinition(service).Return(taskCopy, oldTask, nil).Times(1)
	client.EXPECT().RegisterTaskDefinition(newTask).Return(newTaskDefinition, nil).Times(1)
	client.EXPECT().UpdateTaskDefinition(service, newTaskDefinition).Return(newService, nil).Times(1)
	client.EXPECT().WaitUntilGood(newService, timeout).Return(timeoutErr).Times(1)
	client.EXPECT().UpdateTaskDefinition(service, oldTask).Return(rolledBackService, nil).Times(1)
	client.EXPECT().WaitUntilGood(rolledBackService, timeout).Return(nil).Times(1)

	cfg := mock.NewMockECSDeployTaskConfig(c)
	cfg.EXPECT().ApplyTo(taskCopy).Return(newTask, diff).Times(1)

	//#endregion

	err := ECSDeploy(clusterName, serviceName, client, timeout, cfg)
	require.NoError(t, err)
}

func Test_ECSDeploy_GetServiceError(t *testing.T) {
	const clusterName = "clusterA"
	const serviceName = "serviceA"

	timeout := time.Nanosecond
	expectedError := errors.New("get service error")
	service := &ecssdk.Service{}

	//#region mock
	c := gomock.NewController(t)

	client := mock.NewMockECSDeployClient(c)
	client.EXPECT().GetService(clusterName, serviceName).Return(service, expectedError).Times(1)

	cfg := mock.NewMockECSDeployTaskConfig(c)
	//#endregion

	err := ECSDeploy(clusterName, serviceName, client, timeout, cfg)
	require.Equal(t, expectedError, err)
}

func Test_ECSDeploy_LooksGoodError(t *testing.T) {
	const clusterName = "clusterA"
	const serviceName = "serviceA"

	timeout := time.Nanosecond
	expectedError := errors.New("looks good error")
	service := &ecssdk.Service{}

	//#region mock
	c := gomock.NewController(t)

	client := mock.NewMockECSDeployClient(c)
	client.EXPECT().GetService(clusterName, serviceName).Return(service, nil).Times(1)
	client.EXPECT().LooksGood(service).Return(true, expectedError).Times(1)

	cfg := mock.NewMockECSDeployTaskConfig(c)
	//#endregion

	err := ECSDeploy(clusterName, serviceName, client, timeout, cfg)
	require.Equal(t, expectedError, err)
}

func Test_ECSDeploy_CopyTaskDefinitionError(t *testing.T) {
	const clusterName = "clusterA"
	const serviceName = "serviceA"

	timeout := time.Nanosecond
	expectedError := errors.New("copy task definition error")
	service := &ecssdk.Service{}

	//#region mock
	c := gomock.NewController(t)

	client := mock.NewMockECSDeployClient(c)
	client.EXPECT().GetService(clusterName, serviceName).Return(service, nil).Times(1)
	client.EXPECT().LooksGood(service).Return(true, nil).Times(1)
	client.EXPECT().CopyTaskDefinition(service).Return(nil, nil, expectedError).Times(1)

	cfg := mock.NewMockECSDeployTaskConfig(c)
	//#endregion

	err := ECSDeploy(clusterName, serviceName, client, timeout, cfg)
	require.Equal(t, expectedError, err)
}
