package action

import (
	"context"
	"errors"
	"testing"
	"time"

	mock "github.com/adroll/ecs-ship/action/mock"
	"github.com/adroll/ecs-ship/ecs"

	ecssdk "github.com/aws/aws-sdk-go-v2/service/ecs"
	"github.com/aws/aws-sdk-go-v2/service/ecs/types"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_ECSDeploy_NoArgs(t *testing.T) {
	err := ECSDeploy(context.Background(), "", "", nil, time.Nanosecond, nil, false, false)
	assert.Error(t, err)
}

func Test_ECSDeploy_TruePath(t *testing.T) {
	const clusterName = "clusterA"
	const serviceName = "serviceA"

	family := "test-family"
	serviceTaskArn := "test-family:0"
	taskArn := "test-family:1"
	timeout := time.Nanosecond
	service := &types.Service{
		TaskDefinition: &serviceTaskArn,
	}
	oldTask := &types.TaskDefinition{}
	taskCopy := &ecssdk.RegisterTaskDefinitionInput{}
	newTask := &ecssdk.RegisterTaskDefinitionInput{
		ContainerDefinitions: make([]types.ContainerDefinition, 0),
		Family:               &family,
	}
	newTaskDefinition := &types.TaskDefinition{
		TaskDefinitionArn: &taskArn,
	}
	newService := &types.Service{}
	diff := &ecs.TaskConfigDiff{}
	oldCpu := "1"
	newCpu := "10"
	diff.ChangeCPU(&oldCpu, &newCpu)

	//#region mock
	c := gomock.NewController(t)

	ctx := context.Background()
	client := mock.NewMockECSDeployClient(c)
	client.EXPECT().GetService(ctx, clusterName, serviceName).Return(service, nil).Times(1)
	client.EXPECT().LooksGood(ctx, service).Return(true, nil).Times(1)
	client.EXPECT().CopyTaskDefinition(ctx, service).Return(taskCopy, oldTask, nil).Times(1)
	client.EXPECT().RegisterTaskDefinition(ctx, newTask).Return(newTaskDefinition, nil).Times(1)
	client.EXPECT().UpdateTaskDefinition(ctx, service, newTaskDefinition).Return(newService, nil).Times(1)
	client.EXPECT().WaitUntilGood(ctx, newService, &timeout).Return(nil).Times(1)

	cfg := mock.NewMockECSDeployTaskConfig(c)
	cfg.EXPECT().ApplyTo(taskCopy).Return(newTask, diff).Times(1)

	//#endregion

	err := ECSDeploy(ctx, clusterName, serviceName, client, timeout, cfg, false, false)
	require.NoError(t, err)
}

func Test_ECSDeploy_NoWait(t *testing.T) {
	const clusterName = "clusterA"
	const serviceName = "serviceA"

	family := "test-family"
	serviceTaskArn := "test-family:0"
	taskArn := "test-family:1"
	timeout := time.Nanosecond
	service := &types.Service{
		TaskDefinition: &serviceTaskArn,
	}
	oldTask := &types.TaskDefinition{}
	taskCopy := &ecssdk.RegisterTaskDefinitionInput{}
	newTask := &ecssdk.RegisterTaskDefinitionInput{
		ContainerDefinitions: make([]types.ContainerDefinition, 0),
		Family:               &family,
	}
	newTaskDefinition := &types.TaskDefinition{
		TaskDefinitionArn: &taskArn,
	}
	newService := &types.Service{}
	diff := &ecs.TaskConfigDiff{}
	oldCpu := "1"
	newCpu := "10"
	diff.ChangeCPU(&oldCpu, &newCpu)

	//#region mock
	c := gomock.NewController(t)

	ctx := context.Background()
	client := mock.NewMockECSDeployClient(c)
	client.EXPECT().GetService(ctx, clusterName, serviceName).Return(service, nil).Times(1)
	client.EXPECT().LooksGood(ctx, service).Return(true, nil).Times(1)
	client.EXPECT().CopyTaskDefinition(ctx, service).Return(taskCopy, oldTask, nil).Times(1)
	client.EXPECT().RegisterTaskDefinition(ctx, newTask).Return(newTaskDefinition, nil).Times(1)
	client.EXPECT().UpdateTaskDefinition(ctx, service, newTaskDefinition).Return(newService, nil).Times(1)
	client.EXPECT().WaitUntilGood(ctx, newService, &timeout).Return(nil).Times(0)

	cfg := mock.NewMockECSDeployTaskConfig(c)
	cfg.EXPECT().ApplyTo(taskCopy).Return(newTask, diff).Times(1)

	//#endregion

	err := ECSDeploy(ctx, clusterName, serviceName, client, timeout, cfg, false, true)
	require.NoError(t, err)
}

func Test_ECSDeploy_DryRunPath(t *testing.T) {
	const clusterName = "clusterA"
	const serviceName = "serviceA"

	family := "test-family"
	serviceTaskArn := "test-family:0"
	timeout := time.Nanosecond
	service := &types.Service{
		TaskDefinition: &serviceTaskArn,
	}
	oldTask := &types.TaskDefinition{}
	taskCopy := &ecssdk.RegisterTaskDefinitionInput{}
	newTask := &ecssdk.RegisterTaskDefinitionInput{
		ContainerDefinitions: make([]types.ContainerDefinition, 0),
		Family:               &family,
	}
	diff := &ecs.TaskConfigDiff{}
	oldCpu := "1"
	newCpu := "10"
	diff.ChangeCPU(&oldCpu, &newCpu)

	//#region mock
	c := gomock.NewController(t)

	ctx := context.Background()
	client := mock.NewMockECSDeployClient(c)
	client.EXPECT().GetService(ctx, clusterName, serviceName).Return(service, nil).Times(1)
	client.EXPECT().LooksGood(ctx, service).Return(true, nil).Times(1)
	client.EXPECT().CopyTaskDefinition(ctx, service).Return(taskCopy, oldTask, nil).Times(1)

	cfg := mock.NewMockECSDeployTaskConfig(c)
	cfg.EXPECT().ApplyTo(taskCopy).Return(newTask, diff).Times(1)

	//#endregion

	err := ECSDeploy(ctx, clusterName, serviceName, client, timeout, cfg, true, false)
	require.NoError(t, err)
}

func Test_ECSDeploy_RollbackPath(t *testing.T) {
	const clusterName = "clusterA"
	const serviceName = "serviceA"

	timeout := time.Nanosecond
	timeoutErr := errors.New("timeout error")
	service := &types.Service{}
	oldTask := &types.TaskDefinition{}
	taskCopy := &ecssdk.RegisterTaskDefinitionInput{}
	newTask := &ecssdk.RegisterTaskDefinitionInput{}
	newTaskDefinition := &types.TaskDefinition{}
	newService := &types.Service{}
	rolledBackService := &types.Service{}
	diff := &ecs.TaskConfigDiff{}

	//#region mock
	c := gomock.NewController(t)

	ctx := context.Background()
	client := mock.NewMockECSDeployClient(c)
	client.EXPECT().GetService(ctx, clusterName, serviceName).Return(service, nil).Times(1)
	client.EXPECT().LooksGood(ctx, service).Return(true, nil).Times(1)
	client.EXPECT().CopyTaskDefinition(ctx, service).Return(taskCopy, oldTask, nil).Times(1)
	client.EXPECT().RegisterTaskDefinition(ctx, newTask).Return(newTaskDefinition, nil).Times(1)
	client.EXPECT().UpdateTaskDefinition(ctx, service, newTaskDefinition).Return(newService, nil).Times(1)
	client.EXPECT().WaitUntilGood(ctx, newService, timeout).Return(timeoutErr).Times(1)
	client.EXPECT().UpdateTaskDefinition(ctx, service, oldTask).Return(rolledBackService, nil).Times(1)
	client.EXPECT().WaitUntilGood(ctx, rolledBackService, timeout).Return(nil).Times(1)

	cfg := mock.NewMockECSDeployTaskConfig(c)
	cfg.EXPECT().ApplyTo(taskCopy).Return(newTask, diff).Times(1)

	//#endregion

	err := ECSDeploy(ctx, clusterName, serviceName, client, timeout, cfg, false, false)
	require.NoError(t, err)
}

func Test_ECSDeploy_GetServiceError(t *testing.T) {
	const clusterName = "clusterA"
	const serviceName = "serviceA"

	timeout := time.Nanosecond
	expectedError := errors.New("get service error")
	service := &types.Service{}

	//#region mock
	c := gomock.NewController(t)

	ctx := context.Background()
	client := mock.NewMockECSDeployClient(c)
	client.EXPECT().GetService(ctx, clusterName, serviceName).Return(service, expectedError).Times(1)

	cfg := mock.NewMockECSDeployTaskConfig(c)
	//#endregion

	err := ECSDeploy(ctx, clusterName, serviceName, client, timeout, cfg, false, false)
	require.Equal(t, expectedError, err)
}

func Test_ECSDeploy_LooksGoodError(t *testing.T) {
	const clusterName = "clusterA"
	const serviceName = "serviceA"

	timeout := time.Nanosecond
	expectedError := errors.New("looks good error")
	service := &types.Service{}

	//#region mock
	c := gomock.NewController(t)

	ctx := context.Background()
	client := mock.NewMockECSDeployClient(c)
	client.EXPECT().GetService(ctx, clusterName, serviceName).Return(service, nil).Times(1)
	client.EXPECT().LooksGood(ctx, service).Return(true, expectedError).Times(1)

	cfg := mock.NewMockECSDeployTaskConfig(c)
	//#endregion

	err := ECSDeploy(ctx, clusterName, serviceName, client, timeout, cfg, false, false)
	require.Equal(t, expectedError, err)
}

func Test_ECSDeploy_CopyTaskDefinitionError(t *testing.T) {
	const clusterName = "clusterA"
	const serviceName = "serviceA"

	timeout := time.Nanosecond
	expectedError := errors.New("copy task definition error")
	service := &types.Service{}

	//#region mock
	c := gomock.NewController(t)

	ctx := context.Background()
	client := mock.NewMockECSDeployClient(c)
	client.EXPECT().GetService(ctx, clusterName, serviceName).Return(service, nil).Times(1)
	client.EXPECT().LooksGood(ctx, service).Return(true, nil).Times(1)
	client.EXPECT().CopyTaskDefinition(ctx, service).Return(nil, nil, expectedError).Times(1)

	cfg := mock.NewMockECSDeployTaskConfig(c)
	//#endregion

	err := ECSDeploy(ctx, clusterName, serviceName, client, timeout, cfg, false, false)
	require.Equal(t, expectedError, err)
}
