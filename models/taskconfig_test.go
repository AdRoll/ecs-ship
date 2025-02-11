package models_test

import (
	"testing"

	"github.com/adroll/ecs-ship/models"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ecs"
	"github.com/aws/aws-sdk-go-v2/service/ecs/types"
	"github.com/stretchr/testify/assert"
)

func Test_TaskConfig_ApplyTo_Empty(t *testing.T) {
	taskConfig := &models.TaskConfig{}
	input := &ecs.RegisterTaskDefinitionInput{}
	_, diff := taskConfig.ApplyTo(input)
	assert.True(t, diff.Empty())
}

func Test_TaskConfig_ApplyTo_CPU(t *testing.T) {
	taskConfig := &models.TaskConfig{
		CPU: aws.String("256"),
	}
	input := &ecs.RegisterTaskDefinitionInput{}
	newInput, diff := taskConfig.ApplyTo(input)
	assert.False(t, diff.Empty())
	assert.Equal(t, "CPU was: <nil> and now is: \"256\"", diff.String())
	assert.NotNil(t, newInput.Cpu)
	assert.Equal(t, "256", *newInput.Cpu)
}

func Test_TaskConfig_ApplyTo_Memory(t *testing.T) {
	taskConfig := &models.TaskConfig{
		Memory: aws.String("512"),
	}
	input := &ecs.RegisterTaskDefinitionInput{}
	newInput, diff := taskConfig.ApplyTo(input)
	assert.False(t, diff.Empty())
	assert.Equal(t, "memory was: <nil> and now is: \"512\"", diff.String())
	assert.NotNil(t, newInput.Memory)
	assert.Equal(t, "512", *newInput.Memory)
}

func Test_TaskConfig_ApplyTo_ContainerDefinitions_NonExistent(t *testing.T) {
	taskConfig := &models.TaskConfig{
		ContainerDefinitions: map[string]models.ContainerConfig{
			"container": {
				CPU: aws.Int32(256),
			},
		},
	}
	input := &ecs.RegisterTaskDefinitionInput{}
	_, diff := taskConfig.ApplyTo(input)
	assert.True(t, diff.Empty())
}

func Test_TaskConfig_ApplyTo_ContainerDefinitions_Existent(t *testing.T) {
	taskConfig := &models.TaskConfig{
		ContainerDefinitions: map[string]models.ContainerConfig{
			"container": {
				CPU: aws.Int32(256),
			},
		},
	}
	input := &ecs.RegisterTaskDefinitionInput{
		ContainerDefinitions: []types.ContainerDefinition{
			{
				Name: aws.String("container"),
			},
		},
	}
	newInput, diff := taskConfig.ApplyTo(input)
	assert.False(t, diff.Empty())
	assert.Equal(t, "the container definition \"container\" changed in this way:\nCPU was: 0 and now is: 256\n", diff.String())
	assert.Equal(t, 1, len(newInput.ContainerDefinitions))
	assert.Equal(t, int32(256), newInput.ContainerDefinitions[0].Cpu)
}
