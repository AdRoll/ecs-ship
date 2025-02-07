package services

import (
	"context"
	"log"
	"time"

	"github.com/adroll/ecs-ship/clients"
	"github.com/adroll/ecs-ship/models"
	"github.com/fatih/color"
	"github.com/joomcode/errorx"
)

// DeployInput represents the input for the DeployerService
type DeployInput struct {
	// Cluster is the name of the ECS cluster
	Cluster string
	// Service is the name of the ECS service
	Service string
	// NewConfig is the new configuration for the service
	NewConfig models.TaskConfig
	// DryRun will only show what would change in the remote service
	DryRun bool
	// Timeout is the time to wait for the service to be correctly updated
	Timeout time.Duration
	// NoWait will disable waiting for updates to be completed
	NoWait bool
}

// DeployerService is the interface for the deployer service
type DeployerService interface {
	// Deploy will deploy the new configuration to the service
	Deploy(ctx context.Context, input *DeployInput) error
}

type deployerService struct {
	client clients.ECSClient
}

// NewDeployerService creates a new DeployerService
func NewDeployerService(client clients.ECSClient) DeployerService {
	return &deployerService{client: client}
}

func (s *deployerService) Deploy(ctx context.Context, input *DeployInput) error {
	log.Printf("updating service:\n  cluster: %s\n  service: %s\n", input.Cluster, input.Service)
	service, err := s.client.GetService(ctx, input.Cluster, input.Service)
	if err != nil {
		return errorx.Decorate(err, "unable to get service")
	}

	looksGood, err := s.client.DoesServiceLookGood(ctx, service)
	if err != nil {
		return errorx.Decorate(err, "unable to check if service looks good")
	}
	if looksGood {
		log.Println(color.GreenString("the service looks good to begin with"))
	}

	output, err := s.client.GetTaskDefinition(ctx, service)
	if err != nil {
		return errorx.Decorate(err, "unable to get task definition")
	}

	oldTaskDefinitionInput := s.client.CopiedTaskDefinition(output)

	newTaskDefinitionInput, diff := input.NewConfig.ApplyTo(oldTaskDefinitionInput)

	if diff.Empty() {
		log.Println(color.GreenString("the service is up to date, we have nothing to do :d"))
		return nil
	}

	log.Println("these are the changes:")
	log.Println(diff)

	if input.DryRun {
		log.Println("not proceeding with the updates because this is a dry run :)")
		return nil
	}

	newTaskDefinition, err := s.client.RegisterTaskDefinition(ctx, newTaskDefinitionInput)
	if err != nil {
		return errorx.Decorate(err, "unable to register new task definition")
	}

	newService, err := s.client.UpdateTaskDefinition(ctx, service, newTaskDefinition)
	if err != nil {
		return errorx.Decorate(err, "unable to update service")
	}

	if input.NoWait {
		log.Println("not waiting for your service to reflect the cahnges, check the console instead.")
		return nil
	}

	log.Println("waiting for the service to reflect the changes")

	err = s.client.WaitForServiceToLookGood(ctx, newService, input.Timeout)

	if err != nil {
		return errorx.Decorate(err, "service did not reflect changes")
	}

	log.Println(color.GreenString("service has been updated successfully!"))

	return nil
}
