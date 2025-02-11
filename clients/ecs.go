package clients

//go:generate mockgen -destination=mocks/ecs.go . ECSClient

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ecs"
	ecsTypes "github.com/aws/aws-sdk-go-v2/service/ecs/types"
	"github.com/joomcode/errorx"
)

const sleepTime = 5 * time.Second

type ECSClient interface {
	GetService(ctx context.Context, clusterName string, serviceName string) (*ecsTypes.Service, error)
	DoesServiceLookGood(ctx context.Context, service *ecsTypes.Service) (bool, error)
	WaitForServiceToLookGood(ctx context.Context, service *ecsTypes.Service, timeout time.Duration) error
	GetTaskDefinition(ctx context.Context, service *ecsTypes.Service) (*ecs.DescribeTaskDefinitionOutput, error)
	CopiedTaskDefinition(output *ecs.DescribeTaskDefinitionOutput) *ecs.RegisterTaskDefinitionInput
	RegisterTaskDefinition(ctx context.Context, input *ecs.RegisterTaskDefinitionInput) (*ecsTypes.TaskDefinition, error)
	UpdateTaskDefinition(ctx context.Context, service *ecsTypes.Service, task *ecsTypes.TaskDefinition) (*ecsTypes.Service, error)
}

type ecsClient struct {
	client *ecs.Client
}

func NewECSClient(client *ecs.Client) ECSClient {
	return &ecsClient{client: client}
}

func (c *ecsClient) GetService(ctx context.Context, clusterName string, serviceName string) (*ecsTypes.Service, error) {
	describeResult, err := c.client.DescribeServices(
		ctx,
		&ecs.DescribeServicesInput{
			Cluster:  aws.String(clusterName),
			Services: []string{serviceName},
		},
	)
	if err != nil {
		return nil, errorx.Decorate(err, "unable to describe services")
	}
	if len(describeResult.Services) == 0 {
		return nil, fmt.Errorf("service %s not found in cluster %s", serviceName, clusterName)
	}
	if len(describeResult.Services) > 1 {
		return nil, fmt.Errorf("many services %s found in cluster %s", serviceName, clusterName)
	}
	return &describeResult.Services[0], nil
}

func (c *ecsClient) DoesServiceLookGood(ctx context.Context, service *ecsTypes.Service) (bool, error) {
	if len(service.Deployments) != 1 {
		return false, nil
	}

	runningTasks, err := c.client.ListTasks(ctx, &ecs.ListTasksInput{
		Cluster:     service.ClusterArn,
		ServiceName: service.ServiceName,
	})
	if err != nil {
		return false, errorx.Decorate(err, "unable to list tasks")
	}
	if len(runningTasks.TaskArns) == 0 {
		return service.DesiredCount == 0, nil
	}
	runningTaskDetails, err := c.client.DescribeTasks(ctx, &ecs.DescribeTasksInput{
		Cluster: service.ClusterArn,
		Tasks:   runningTasks.TaskArns,
	})
	if err != nil {
		return false, errorx.Decorate(err, "unable to describe tasks")
	}
	matchCount := 0
	for _, task := range runningTaskDetails.Tasks {
		if *task.TaskDefinitionArn == *service.TaskDefinition && *task.LastStatus == string(ecsTypes.DesiredStatusRunning) {
			matchCount++
		}
	}
	return int32(matchCount) == service.DesiredCount, nil
}

func (c *ecsClient) WaitForServiceToLookGood(ctx context.Context, service *ecsTypes.Service, timeout time.Duration) error {
	refreshService, err := c.GetService(ctx, *service.ClusterArn, *service.ServiceName)
	if err != nil {
		return errorx.Decorate(err, "unable to get service")
	}

	alreadyLookingGood := false
	var errorMessages []string
	var checkedUntil *time.Time

	ticker := time.NewTicker(sleepTime)
	defer ticker.Stop()

	flushTicker := time.NewTicker(80 * sleepTime)
	defer flushTicker.Stop()

	timeoutChan := time.After(timeout)

	for {
		select {
		case <-ticker.C:
			if alreadyLookingGood, err = c.DoesServiceLookGood(ctx, refreshService); alreadyLookingGood {
				return nil
			}
			if err != nil {
				return errorx.Decorate(err, "unable to check if service looks good")
			}
			refreshService, err = c.GetService(ctx, *service.ClusterArn, *service.ServiceName)
			if err != nil {
				return errorx.Decorate(err, "unable to get service")
			}
			var newErrors []string
			if newErrors, checkedUntil, err = c.getRecentErrorMessages(refreshService, checkedUntil); err != nil {
				errorMessages = append(errorMessages, newErrors...)
			}
			fmt.Print(".")
		case <-flushTicker.C:
			fmt.Println()
		case <-timeoutChan:
			fmt.Println()
			if len(errorMessages) == 0 {
				return errors.New(strings.Join([]string{
					"We ran into a timeout while waiting for the service to reach steady state",
					"We found no errors, so consider just increasing the timeout you're using",
					"(or maybe the service doesn't have logging setup correctly)",
				}, "\n"))
			}
			return errors.New(strings.Join(append([]string{
				"We ran into a timeout while waiting for the service to reach steady state",
				"We found the following errors while trying to get your service up:",
			}, errorMessages...), "\n"))
		}
	}
}

func (c *ecsClient) getRecentErrorMessages(service *ecsTypes.Service, after *time.Time) ([]string, *time.Time, error) {
	var reportSince *time.Time
	if after == nil {
		var deployedAt *time.Time
		for _, deployment := range service.Deployments {
			if *deployment.Status == "PRIMARY" {
				deployedAt = deployment.CreatedAt
			}
		}
		if deployedAt == nil {
			return nil, nil, errors.New("we could not find a primary deployment for the given service")
		}
		reportSince = deployedAt
	} else {
		reportSince = after
	}

	var listErrs []string
	var reportedUntil *time.Time
	for _, event := range service.Events {
		if event.CreatedAt.After(*reportSince) && strings.Contains(*event.Message, "unable") {
			listErrs = append(listErrs, fmt.Sprintf("[ERROR] %s: %s", *event.CreatedAt, *event.Message))
		}
		if reportedUntil == nil || event.CreatedAt.After(*reportedUntil) {
			reportedUntil = event.CreatedAt
		}
	}

	return listErrs, reportedUntil, nil
}

func (c *ecsClient) GetTaskDefinition(ctx context.Context, service *ecsTypes.Service) (*ecs.DescribeTaskDefinitionOutput, error) {
	output, err := c.client.DescribeTaskDefinition(ctx, &ecs.DescribeTaskDefinitionInput{
		TaskDefinition: service.TaskDefinition,
		Include:        []ecsTypes.TaskDefinitionField{ecsTypes.TaskDefinitionFieldTags},
	})
	if err != nil {
		return nil, errorx.Decorate(err, "unable to describe task definition")
	}
	return output, nil
}

func (c *ecsClient) CopiedTaskDefinition(output *ecs.DescribeTaskDefinitionOutput) *ecs.RegisterTaskDefinitionInput {
	tags := output.Tags
	if len(tags) == 0 {
		tags = nil
	}
	return &ecs.RegisterTaskDefinitionInput{
		ContainerDefinitions:    output.TaskDefinition.ContainerDefinitions,
		Family:                  output.TaskDefinition.Family,
		Cpu:                     output.TaskDefinition.Cpu,
		EphemeralStorage:        output.TaskDefinition.EphemeralStorage,
		ExecutionRoleArn:        output.TaskDefinition.ExecutionRoleArn,
		InferenceAccelerators:   output.TaskDefinition.InferenceAccelerators,
		IpcMode:                 output.TaskDefinition.IpcMode,
		Memory:                  output.TaskDefinition.Memory,
		NetworkMode:             output.TaskDefinition.NetworkMode,
		PidMode:                 output.TaskDefinition.PidMode,
		PlacementConstraints:    output.TaskDefinition.PlacementConstraints,
		ProxyConfiguration:      output.TaskDefinition.ProxyConfiguration,
		RequiresCompatibilities: output.TaskDefinition.RequiresCompatibilities,
		RuntimePlatform:         output.TaskDefinition.RuntimePlatform,
		Tags:                    tags,
		TaskRoleArn:             output.TaskDefinition.TaskRoleArn,
		Volumes:                 output.TaskDefinition.Volumes,
	}
}

func (c *ecsClient) RegisterTaskDefinition(ctx context.Context, input *ecs.RegisterTaskDefinitionInput) (*ecsTypes.TaskDefinition, error) {
	output, err := c.client.RegisterTaskDefinition(ctx, input)
	if err != nil {
		return nil, errorx.Decorate(err, "unable to register task definition")
	}
	return output.TaskDefinition, nil
}

func (c *ecsClient) UpdateTaskDefinition(ctx context.Context, service *ecsTypes.Service, task *ecsTypes.TaskDefinition) (*ecsTypes.Service, error) {
	output, err := c.client.UpdateService(ctx, &ecs.UpdateServiceInput{
		Cluster:        service.ClusterArn,
		Service:        service.ServiceName,
		TaskDefinition: task.TaskDefinitionArn,
	})
	if err != nil {
		return nil, errorx.Decorate(err, "unable to update service")
	}
	return output.Service, nil
}
