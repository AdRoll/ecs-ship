package action

import (
	"errors"
	"log"
	"time"

	"github.com/adroll/ecs-ship/ecs"
	ecssdk "github.com/aws/aws-sdk-go/service/ecs"
	"github.com/fatih/color"
)

// ECSDeployClient defines a simple interface for our required ecs stuffs
type ECSDeployClient interface {
	GetService(clusterName string, serviceName string) (*ecssdk.Service, error)
	LooksGood(service *ecssdk.Service) (bool, error)
	CopyTaskDefinition(service *ecssdk.Service) (*ecssdk.RegisterTaskDefinitionInput, *ecssdk.TaskDefinition, error)
	WaitUntilGood(service *ecssdk.Service, timeout *time.Duration) error
	RegisterTaskDefinition(input *ecssdk.RegisterTaskDefinitionInput) (*ecssdk.TaskDefinition, error)
	UpdateTaskDefinition(service *ecssdk.Service, task *ecssdk.TaskDefinition) (*ecssdk.Service, error)
}

// ECSDeployTaskConfig defines a simple interface of how we want a config thing to do
type ECSDeployTaskConfig interface {
	ApplyTo(input *ecssdk.RegisterTaskDefinitionInput) (*ecssdk.RegisterTaskDefinitionInput, *ecs.TaskConfigDiff)
}

// ECSDeploy deploy an ecs service
func ECSDeploy(clusterName string, serviceName string, client ECSDeployClient, timeout time.Duration, config ECSDeployTaskConfig) error {
	if len(clusterName) == 0 {
		return errors.New("cluster was not provided")
	}
	if len(serviceName) == 0 {
		return errors.New("service was not provided")
	}
	if client == nil {
		return errors.New("client was not provided")
	}
	if config == nil {
		return errors.New("config was not provided")
	}

	service, err := client.GetService(clusterName, serviceName)
	if err != nil {
		return err
	}

	log.Printf("Updating service:\n  Cluster: %s\n  Service: %s\n", clusterName, serviceName)

	good, err := client.LooksGood(service)
	if err != nil {
		return err
	}
	if good {
		log.Println(color.GreenString("The service looks good to begin with"))
	} else {
		log.Println(color.YellowString("The service doesn't look good to begin with"))
	}

	copyTask, oldTaskDefinition, err := client.CopyTaskDefinition(service)
	if err != nil {
		return err
	}

	newTask, diff := config.ApplyTo(copyTask)
	if diff.Empty() {
		log.Println(color.GreenString("The service is up to date, we have nothing to do :D"))
		return nil
	}

	if err := newTask.Validate(); err != nil {
		return err
	}

	newTaskDefinition, err := client.RegisterTaskDefinition(newTask)
	if err != nil {
		return err
	}

	log.Printf("Changing task definition\n  Old: %s\n  New: %s\n", *service.TaskDefinition, *newTaskDefinition.TaskDefinitionArn)
	log.Println("These are the changes:")
	log.Println(diff)

	newService, err := client.UpdateTaskDefinition(service, newTaskDefinition)
	if err != nil {
		return err
	}

	log.Println("Waiting for the service to reflect the new changes...")
	if originalErr := client.WaitUntilGood(newService, &timeout); originalErr != nil {
		log.Println(color.RedString("There was an error updating the service :", originalErr.Error()))
		log.Println(color.YellowString("we are trying to roll back changes..."))
		rolledBackService, err := client.UpdateTaskDefinition(service, oldTaskDefinition)
		if err != nil {
			log.Println(color.RedString("You're unlucky we also failed to roll back the service with error:", err.Error()))
			return originalErr
		}
		log.Println("Waiting for rollback service to reflect the new changes...")
		if err = client.WaitUntilGood(rolledBackService, &timeout); good && err != nil {
			log.Println(color.RedString("stopped waiting with error:", err.Error()))
			return originalErr
		}
		log.Println(color.GreenString("Order restored, but still reporting on the original error."))
		return originalErr
	}
	log.Println(color.GreenString("Now everything looks good"))

	return nil
}
