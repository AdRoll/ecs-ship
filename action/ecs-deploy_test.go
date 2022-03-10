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
	err := ECSDeploy("", "", nil, time.Nanosecond, nil, false, false)
	assert.Error(t, err)
}

func Test_ECSDeploy_TruePath(t *testing.T) {
	const clusterName = "clusterA"
	const serviceName = "serviceA"

	family := "test-family"
	serviceTaskArn := "test-family:0"
	taskArn := "test-family:1"
	timeout := time.Nanosecond
	service := &ecssdk.Service{
		TaskDefinition: &serviceTaskArn,
	}
	oldTask := &ecssdk.TaskDefinition{}
	taskCopy := &ecssdk.RegisterTaskDefinitionInput{}
	newTask := &ecssdk.RegisterTaskDefinitionInput{
		ContainerDefinitions: make([]*ecssdk.ContainerDefinition, 0),
		Family:               &family,
	}
	newTaskDefinition := &ecssdk.TaskDefinition{
		TaskDefinitionArn: &taskArn,
	}
	newService := &ecssdk.Service{}
	diff := &ecs.TaskConfigDiff{}
	oldCpu := "1"
	newCpu := "10"
	diff.ChangeCPU(&oldCpu, &newCpu)

	//#region mock
	c := gomock.NewController(t)

	client := mock.NewMockECSDeployClient(c)
	client.EXPECT().GetService(clusterName, serviceName).Return(service, nil).Times(1)
	client.EXPECT().LooksGood(service).Return(true, nil).Times(1)
	client.EXPECT().CopyTaskDefinition(service).Return(taskCopy, oldTask, nil).Times(1)
	client.EXPECT().RegisterTaskDefinition(newTask).Return(newTaskDefinition, nil).Times(1)
	client.EXPECT().UpdateTaskDefinition(service, newTaskDefinition).Return(newService, nil).Times(1)
	client.EXPECT().WaitUntilGood(newService, &timeout).Return(nil).Times(1)

	cfg := mock.NewMockECSDeployTaskConfig(c)
	cfg.EXPECT().ApplyTo(taskCopy).Return(newTask, diff).Times(1)

	//#endregion

	err := ECSDeploy(clusterName, serviceName, client, timeout, cfg, false, false)
	require.NoError(t, err)
}

func Test_ECSDeploy_NoWait(t *testing.T) {
	const clusterName = "clusterA"
	const serviceName = "serviceA"

	family := "test-family"
	serviceTaskArn := "test-family:0"
	taskArn := "test-family:1"
	timeout := time.Nanosecond
	service := &ecssdk.Service{
		TaskDefinition: &serviceTaskArn,
	}
	oldTask := &ecssdk.TaskDefinition{}
	taskCopy := &ecssdk.RegisterTaskDefinitionInput{}
	newTask := &ecssdk.RegisterTaskDefinitionInput{
		ContainerDefinitions: make([]*ecssdk.ContainerDefinition, 0),
		Family:               &family,
	}
	newTaskDefinition := &ecssdk.TaskDefinition{
		TaskDefinitionArn: &taskArn,
	}
	newService := &ecssdk.Service{}
	diff := &ecs.TaskConfigDiff{}
	oldCpu := "1"
	newCpu := "10"
	diff.ChangeCPU(&oldCpu, &newCpu)

	//#region mock
	c := gomock.NewController(t)

	client := mock.NewMockECSDeployClient(c)
	client.EXPECT().GetService(clusterName, serviceName).Return(service, nil).Times(1)
	client.EXPECT().LooksGood(service).Return(true, nil).Times(1)
	client.EXPECT().CopyTaskDefinition(service).Return(taskCopy, oldTask, nil).Times(1)
	client.EXPECT().RegisterTaskDefinition(newTask).Return(newTaskDefinition, nil).Times(1)
	client.EXPECT().UpdateTaskDefinition(service, newTaskDefinition).Return(newService, nil).Times(1)
	client.EXPECT().WaitUntilGood(newService, &timeout).Return(nil).Times(0)

	cfg := mock.NewMockECSDeployTaskConfig(c)
	cfg.EXPECT().ApplyTo(taskCopy).Return(newTask, diff).Times(1)

	//#endregion

	err := ECSDeploy(clusterName, serviceName, client, timeout, cfg, false, true)
	require.NoError(t, err)
}

func Test_ECSDeploy_DryRunPath(t *testing.T) {
	const clusterName = "clusterA"
	const serviceName = "serviceA"

	family := "test-family"
	serviceTaskArn := "test-family:0"
	timeout := time.Nanosecond
	service := &ecssdk.Service{
		TaskDefinition: &serviceTaskArn,
	}
	oldTask := &ecssdk.TaskDefinition{}
	taskCopy := &ecssdk.RegisterTaskDefinitionInput{}
	newTask := &ecssdk.RegisterTaskDefinitionInput{
		ContainerDefinitions: make([]*ecssdk.ContainerDefinition, 0),
		Family:               &family,
	}
	diff := &ecs.TaskConfigDiff{}
	oldCpu := "1"
	newCpu := "10"
	diff.ChangeCPU(&oldCpu, &newCpu)

	//#region mock
	c := gomock.NewController(t)

	client := mock.NewMockECSDeployClient(c)
	client.EXPECT().GetService(clusterName, serviceName).Return(service, nil).Times(1)
	client.EXPECT().LooksGood(service).Return(true, nil).Times(1)
	client.EXPECT().CopyTaskDefinition(service).Return(taskCopy, oldTask, nil).Times(1)

	cfg := mock.NewMockECSDeployTaskConfig(c)
	cfg.EXPECT().ApplyTo(taskCopy).Return(newTask, diff).Times(1)

	//#endregion

	err := ECSDeploy(clusterName, serviceName, client, timeout, cfg, true, false)
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

	err := ECSDeploy(clusterName, serviceName, client, timeout, cfg, false, false)
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

	err := ECSDeploy(clusterName, serviceName, client, timeout, cfg, false, false)
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

	err := ECSDeploy(clusterName, serviceName, client, timeout, cfg, false, false)
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

	err := ECSDeploy(clusterName, serviceName, client, timeout, cfg, false, false)
	require.Equal(t, expectedError, err)
}
