package services_test

import (
	"context"
	"testing"

	mock_clients "github.com/adroll/ecs-ship/clients/mocks"
	"github.com/adroll/ecs-ship/models"
	"github.com/adroll/ecs-ship/services"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ecs"
	"github.com/aws/aws-sdk-go-v2/service/ecs/types"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func Test_Deployer_Deploy_UnableToGetService(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockClient := mock_clients.NewMockECSClient(ctrl)
	deployer := services.NewDeployerService(mockClient)
	ctx := context.Background()
	input := &services.DeployInput{
		Cluster:   "cluster",
		Service:   "service",
		NewConfig: models.TaskConfig{},
		DryRun:    false,
		Timeout:   0,
		NoWait:    false,
	}

	mockClient.EXPECT().GetService(ctx, input.Cluster, input.Service).Return(nil, assert.AnError)

	err := deployer.Deploy(context.Background(), input)
	assert.Error(t, err)
	assert.Equal(t, "unable to get service, cause: assert.AnError general error for testing", err.Error())
}

func Test_Deployer_Deploy_UnableToCheckIfServiceLooksGood(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockClient := mock_clients.NewMockECSClient(ctrl)
	deployer := services.NewDeployerService(mockClient)
	ctx := context.Background()
	input := &services.DeployInput{
		Cluster:   "cluster",
		Service:   "service",
		NewConfig: models.TaskConfig{},
		DryRun:    false,
		Timeout:   0,
		NoWait:    false,
	}

	service := &types.Service{
		ServiceName: aws.String("service"),
		ClusterArn:  aws.String("arn::cluster"),
	}
	mockClient.EXPECT().GetService(ctx, input.Cluster, input.Service).Return(service, nil)
	mockClient.EXPECT().DoesServiceLookGood(ctx, service).Return(false, assert.AnError)

	err := deployer.Deploy(context.Background(), input)
	assert.Error(t, err)
	assert.Equal(t, "unable to check if service looks good, cause: assert.AnError general error for testing", err.Error())
}

func Test_Deployer_Deploy_UnableToGetTaskDefinition(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockClient := mock_clients.NewMockECSClient(ctrl)
	deployer := services.NewDeployerService(mockClient)
	ctx := context.Background()
	input := &services.DeployInput{
		Cluster:   "cluster",
		Service:   "service",
		NewConfig: models.TaskConfig{},
		DryRun:    false,
		Timeout:   0,
		NoWait:    false,
	}

	service := &types.Service{
		ServiceName: aws.String("service"),
		ClusterArn:  aws.String("arn::cluster"),
	}
	mockClient.EXPECT().GetService(ctx, input.Cluster, input.Service).Return(service, nil)
	mockClient.EXPECT().DoesServiceLookGood(ctx, service).Return(true, nil)
	mockClient.EXPECT().GetTaskDefinition(ctx, service).Return(nil, assert.AnError)

	err := deployer.Deploy(context.Background(), input)
	assert.Error(t, err)
	assert.Equal(t, "unable to get task definition, cause: assert.AnError general error for testing", err.Error())
}

func Test_Deployer_EmptyChanges(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockClient := mock_clients.NewMockECSClient(ctrl)
	deployer := services.NewDeployerService(mockClient)
	ctx := context.Background()
	input := &services.DeployInput{
		Cluster:   "cluster",
		Service:   "service",
		NewConfig: models.TaskConfig{},
		DryRun:    false,
		Timeout:   0,
		NoWait:    false,
	}

	service := &types.Service{
		ServiceName: aws.String("service"),
		ClusterArn:  aws.String("arn::cluster"),
	}
	taskDefinitonOutput := &ecs.DescribeTaskDefinitionOutput{}
	mockClient.EXPECT().GetService(ctx, input.Cluster, input.Service).Return(service, nil)
	mockClient.EXPECT().DoesServiceLookGood(ctx, service).Return(true, nil)
	mockClient.EXPECT().GetTaskDefinition(ctx, service).Return(taskDefinitonOutput, nil)
	mockClient.EXPECT().CopiedTaskDefinition(taskDefinitonOutput).Return(&ecs.RegisterTaskDefinitionInput{})

	err := deployer.Deploy(context.Background(), input)
	assert.Nil(t, err)
}

func Test_Deployer_UnableToRegisterTaskDefinition(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockClient := mock_clients.NewMockECSClient(ctrl)
	deployer := services.NewDeployerService(mockClient)
	ctx := context.Background()
	input := &services.DeployInput{
		Cluster: "cluster",
		Service: "service",
		NewConfig: models.TaskConfig{
			CPU: aws.String("256"),
		},
		DryRun:  false,
		Timeout: 0,
		NoWait:  false,
	}

	service := &types.Service{
		ServiceName: aws.String("service"),
		ClusterArn:  aws.String("arn::cluster"),
	}
	taskDefinitonOutput := &ecs.DescribeTaskDefinitionOutput{}
	mockClient.EXPECT().GetService(ctx, input.Cluster, input.Service).Return(service, nil)
	mockClient.EXPECT().DoesServiceLookGood(ctx, service).Return(true, nil)
	mockClient.EXPECT().GetTaskDefinition(ctx, service).Return(taskDefinitonOutput, nil)
	mockClient.EXPECT().CopiedTaskDefinition(taskDefinitonOutput).Return(&ecs.RegisterTaskDefinitionInput{})
	mockClient.EXPECT().RegisterTaskDefinition(ctx, gomock.Any()).Return(nil, assert.AnError)

	err := deployer.Deploy(context.Background(), input)
	assert.Error(t, err)
	assert.Equal(t, "unable to register new task definition, cause: assert.AnError general error for testing", err.Error())
}

func Test_Deployer_UnableToUpdateTaskDefinition(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockClient := mock_clients.NewMockECSClient(ctrl)
	deployer := services.NewDeployerService(mockClient)
	ctx := context.Background()
	input := &services.DeployInput{
		Cluster: "cluster",
		Service: "service",
		NewConfig: models.TaskConfig{
			CPU: aws.String("256"),
		},
		DryRun:  false,
		Timeout: 0,
		NoWait:  false,
	}

	service := &types.Service{
		ServiceName: aws.String("service"),
		ClusterArn:  aws.String("arn::cluster"),
	}
	taskDefinitonOutput := &ecs.DescribeTaskDefinitionOutput{}
	newTaskDefiniton := &types.TaskDefinition{}
	mockClient.EXPECT().GetService(ctx, input.Cluster, input.Service).Return(service, nil)
	mockClient.EXPECT().DoesServiceLookGood(ctx, service).Return(true, nil)
	mockClient.EXPECT().GetTaskDefinition(ctx, service).Return(taskDefinitonOutput, nil)
	mockClient.EXPECT().CopiedTaskDefinition(taskDefinitonOutput).Return(&ecs.RegisterTaskDefinitionInput{})
	mockClient.EXPECT().RegisterTaskDefinition(ctx, gomock.Any()).Return(newTaskDefiniton, nil)
	mockClient.EXPECT().UpdateTaskDefinition(ctx, service, newTaskDefiniton).Return(nil, assert.AnError)

	err := deployer.Deploy(context.Background(), input)
	assert.Error(t, err)
	assert.Equal(t, "unable to update service, cause: assert.AnError general error for testing", err.Error())
}

func Test_Deployer_SuccessButNoWait(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockClient := mock_clients.NewMockECSClient(ctrl)
	deployer := services.NewDeployerService(mockClient)
	ctx := context.Background()
	input := &services.DeployInput{
		Cluster: "cluster",
		Service: "service",
		NewConfig: models.TaskConfig{
			CPU: aws.String("256"),
		},
		DryRun:  false,
		Timeout: 0,
		NoWait:  true,
	}

	service := &types.Service{
		ServiceName: aws.String("service"),
		ClusterArn:  aws.String("arn::cluster"),
	}
	taskDefinitonOutput := &ecs.DescribeTaskDefinitionOutput{}
	newTaskDefiniton := &types.TaskDefinition{}
	newService := &types.Service{
		ServiceName: aws.String("service"),
		ClusterArn:  aws.String("arn::cluster"),
	}
	mockClient.EXPECT().GetService(ctx, input.Cluster, input.Service).Return(service, nil)
	mockClient.EXPECT().DoesServiceLookGood(ctx, service).Return(true, nil)
	mockClient.EXPECT().GetTaskDefinition(ctx, service).Return(taskDefinitonOutput, nil)
	mockClient.EXPECT().CopiedTaskDefinition(taskDefinitonOutput).Return(&ecs.RegisterTaskDefinitionInput{})
	mockClient.EXPECT().RegisterTaskDefinition(ctx, gomock.Any()).Return(newTaskDefiniton, nil)
	mockClient.EXPECT().UpdateTaskDefinition(ctx, service, newTaskDefiniton).Return(newService, nil)

	err := deployer.Deploy(context.Background(), input)
	assert.Nil(t, err)
}

func Test_Deployer_ServiceDoesntReflectTheChanges(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockClient := mock_clients.NewMockECSClient(ctrl)
	deployer := services.NewDeployerService(mockClient)
	ctx := context.Background()
	input := &services.DeployInput{
		Cluster: "cluster",
		Service: "service",
		NewConfig: models.TaskConfig{
			CPU: aws.String("256"),
		},
		DryRun:  false,
		Timeout: 0,
		NoWait:  false,
	}

	service := &types.Service{
		ServiceName: aws.String("service"),
		ClusterArn:  aws.String("arn::cluster"),
	}
	taskDefinitonOutput := &ecs.DescribeTaskDefinitionOutput{}
	newTaskDefiniton := &types.TaskDefinition{}
	newService := &types.Service{
		ServiceName: aws.String("service"),
		ClusterArn:  aws.String("arn::cluster"),
	}
	mockClient.EXPECT().GetService(ctx, input.Cluster, input.Service).Return(service, nil)
	mockClient.EXPECT().DoesServiceLookGood(ctx, service).Return(true, nil)
	mockClient.EXPECT().GetTaskDefinition(ctx, service).Return(taskDefinitonOutput, nil)
	mockClient.EXPECT().CopiedTaskDefinition(taskDefinitonOutput).Return(&ecs.RegisterTaskDefinitionInput{})
	mockClient.EXPECT().RegisterTaskDefinition(ctx, gomock.Any()).Return(newTaskDefiniton, nil)
	mockClient.EXPECT().UpdateTaskDefinition(ctx, service, newTaskDefiniton).Return(newService, nil)
	mockClient.EXPECT().WaitForServiceToLookGood(ctx, newService, input.Timeout).Return(assert.AnError)

	err := deployer.Deploy(context.Background(), input)
	assert.Error(t, err)
	assert.Equal(t, "service did not reflect changes, cause: assert.AnError general error for testing", err.Error())
}

func Test_Deployer_ServiceReflectsChanges(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockClient := mock_clients.NewMockECSClient(ctrl)
	deployer := services.NewDeployerService(mockClient)
	ctx := context.Background()
	input := &services.DeployInput{
		Cluster: "cluster",
		Service: "service",
		NewConfig: models.TaskConfig{
			CPU: aws.String("256"),
		},
		DryRun:  false,
		Timeout: 0,
		NoWait:  false,
	}

	service := &types.Service{
		ServiceName: aws.String("service"),
		ClusterArn:  aws.String("arn::cluster"),
	}
	taskDefinitonOutput := &ecs.DescribeTaskDefinitionOutput{}
	newTaskDefiniton := &types.TaskDefinition{}
	newService := &types.Service{
		ServiceName: aws.String("service"),
		ClusterArn:  aws.String("arn::cluster"),
	}
	mockClient.EXPECT().GetService(ctx, input.Cluster, input.Service).Return(service, nil)
	mockClient.EXPECT().DoesServiceLookGood(ctx, service).Return(true, nil)
	mockClient.EXPECT().GetTaskDefinition(ctx, service).Return(taskDefinitonOutput, nil)
	mockClient.EXPECT().CopiedTaskDefinition(taskDefinitonOutput).Return(&ecs.RegisterTaskDefinitionInput{})
	mockClient.EXPECT().RegisterTaskDefinition(ctx, gomock.Any()).Return(newTaskDefiniton, nil)
	mockClient.EXPECT().UpdateTaskDefinition(ctx, service, newTaskDefiniton).Return(newService, nil)
	mockClient.EXPECT().WaitForServiceToLookGood(ctx, newService, input.Timeout).Return(nil)

	err := deployer.Deploy(context.Background(), input)
	assert.Nil(t, err)
}
