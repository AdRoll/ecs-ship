package ecs

import (
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ecs"
)

// Client thin wrapper around ecs
type Client struct {
	service *ecs.ECS
}

// NewClient build a new Client out of a session
func NewClient(session *session.Session) *Client {
	return &Client{
		service: ecs.New(session),
	}
}

// BuildDefaultClient provides new Client with default session config
func BuildDefaultClient() *Client {
	sess := session.Must(session.NewSession())
	return NewClient(sess)
}

// GetService grabs the first service matching the provided cluster and service names
func (client *Client) GetService(clusterName string, serviceName string) (*ecs.Service, error) {
	describeResult, err := client.service.DescribeServices(&ecs.DescribeServicesInput{
		Cluster:  aws.String(clusterName),
		Services: []*string{aws.String(serviceName)},
	})
	if err != nil {
		return nil, err
	}
	if len(describeResult.Services) == 0 {
		return nil, fmt.Errorf("service %s not found in cluster %s", serviceName, clusterName)
	}
	if len(describeResult.Services) > 1 {
		return nil, fmt.Errorf("many services %s found in cluster %s", serviceName, clusterName)
	}
	return describeResult.Services[0], nil
}

// LooksGood checks if a service looks good
func (client *Client) LooksGood(service *ecs.Service) (bool, error) {
	if len(service.Deployments) != 1 {
		return false, nil
	}

	runningTasks, err := client.service.ListTasks(&ecs.ListTasksInput{
		Cluster:     service.ClusterArn,
		ServiceName: service.ServiceName,
	})
	if err != nil {
		return false, err
	}
	if len(runningTasks.TaskArns) == 0 {
		return *service.DesiredCount == 0, nil
	}

	runningTaskDetails, err := client.service.DescribeTasks(&ecs.DescribeTasksInput{
		Cluster: service.ClusterArn,
		Tasks:   runningTasks.TaskArns,
	})
	if err != nil {
		return false, err
	}

	matchCount := 0
	for _, task := range runningTaskDetails.Tasks {
		if *task.TaskDefinitionArn == *service.TaskDefinition && *task.LastStatus == "RUNNING" {
			matchCount++
		}
	}
	return int64(matchCount) == *service.DesiredCount, nil
}

const (
	sleepTime      = 5 * time.Second
	defaultTimeout = 5 * time.Minute
)

// WaitUntilGood wait for a service to look good
func (client *Client) WaitUntilGood(service *ecs.Service, timeout *time.Duration) error {
	refreshService, err := client.GetService(*service.ClusterArn, *service.ServiceName)
	if err != nil {
		return err
	}

	var deadline time.Duration
	if timeout == nil {
		deadline = *timeout
	} else {
		deadline = defaultTimeout
	}

	alreadyLookingGood := false
	var errorMessages []string
	var checkedUntil *time.Time

	defer func() {
		log.Println("")
	}()

	for {
		select {
		case <-time.After(sleepTime):
			if alreadyLookingGood, err = client.LooksGood(refreshService); alreadyLookingGood {
				return nil
			}
			if err != nil {
				return err
			}
			refreshService, err = client.GetService(*service.ClusterArn, *service.ServiceName)
			if err != nil {
				return err
			}
			var newErrors []string
			if newErrors, checkedUntil, err = client.GetRecentErrorMessages(refreshService, checkedUntil); err != nil {
				errorMessages = append(errorMessages, newErrors...)
			}
			fmt.Print(".")
		case <-time.After(deadline):
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

// GetTaskDefinition grabs the first task definition for a service
// NOTE: Tags are in the top level of teh task definition output we need those for later
func (client *Client) GetTaskDefinition(service *ecs.Service) (*ecs.TaskDefinition, error) {
	output, err := client.service.DescribeTaskDefinition(&ecs.DescribeTaskDefinitionInput{
		TaskDefinition: service.TaskDefinition,
		Include:        []*string{aws.String("TAGS")},
	})
	if err != nil {
		return nil, err
	}
	return output.TaskDefinition, nil
}

// GetRecentErrorMessages grabs the first task definition for a service
func (client *Client) GetRecentErrorMessages(service *ecs.Service, after *time.Time) ([]string, *time.Time, error) {
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

// CopyTaskDefinition grabs the first task definition for a service & make copy of it
// NOTE: Tags are in the top level of the task definition output we need those for later
func (client *Client) CopyTaskDefinition(service *ecs.Service) (*ecs.RegisterTaskDefinitionInput, *ecs.TaskDefinition, error) {
	task, err := client.service.DescribeTaskDefinition(&ecs.DescribeTaskDefinitionInput{
		TaskDefinition: service.TaskDefinition,
		Include:        []*string{aws.String("TAGS")},
	})
	if err != nil {
		return nil, nil, err
	}
	return copyTaskDefinition(task), task.TaskDefinition, nil
}

// RegisterTaskDefinition Registers a task definition and returns it
func (client *Client) RegisterTaskDefinition(input *ecs.RegisterTaskDefinitionInput) (*ecs.TaskDefinition, error) {
	output, err := client.service.RegisterTaskDefinition(input)
	if err != nil {
		return nil, err
	}
	return output.TaskDefinition, err
}

// UpdateTaskDefinition Registers a task definition and returns it
func (client *Client) UpdateTaskDefinition(service *ecs.Service, task *ecs.TaskDefinition) (*ecs.Service, error) {
	output, err := client.service.UpdateService(&ecs.UpdateServiceInput{
		Cluster:        service.ClusterArn,
		Service:        service.ServiceName,
		TaskDefinition: task.TaskDefinitionArn,
	})
	if err != nil {
		return nil, err
	}
	return output.Service, nil
}

func copyTaskDefinition(task *ecs.DescribeTaskDefinitionOutput) *ecs.RegisterTaskDefinitionInput {
	tags := task.Tags
	if len(tags) == 0 {
		tags = nil
	}
	return &ecs.RegisterTaskDefinitionInput{
		ContainerDefinitions:    task.TaskDefinition.ContainerDefinitions,
		Cpu:                     task.TaskDefinition.Cpu,
		ExecutionRoleArn:        task.TaskDefinition.ExecutionRoleArn,
		Family:                  task.TaskDefinition.Family,
		InferenceAccelerators:   task.TaskDefinition.InferenceAccelerators,
		IpcMode:                 task.TaskDefinition.IpcMode,
		Memory:                  task.TaskDefinition.Memory,
		NetworkMode:             task.TaskDefinition.NetworkMode,
		PidMode:                 task.TaskDefinition.PidMode,
		PlacementConstraints:    task.TaskDefinition.PlacementConstraints,
		ProxyConfiguration:      task.TaskDefinition.ProxyConfiguration,
		RequiresCompatibilities: task.TaskDefinition.RequiresCompatibilities,
		Tags:                    tags,
		TaskRoleArn:             task.TaskDefinition.TaskRoleArn,
		Volumes:                 task.TaskDefinition.Volumes,
	}
}
