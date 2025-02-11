package models_test

import (
	"testing"

	"github.com/adroll/ecs-ship/models"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/stretchr/testify/assert"
)

func Test_TaskConfigDiff_Empty(t *testing.T) {
	taskConfigDiff := &models.TaskConfigDiff{}
	assert.True(t, taskConfigDiff.Empty())
}

func Test_TaskConfigDiff_CPU(t *testing.T) {
	taskConfigDiff := &models.TaskConfigDiff{}
	taskConfigDiff.ChangeCPU(nil, aws.String("256"))
	assert.False(t, taskConfigDiff.Empty())
	assert.Equal(t, "CPU was: <nil> and now is: \"256\"", taskConfigDiff.String())
}

func Test_TaskConfigDiff_Memory(t *testing.T) {
	taskConfigDiff := &models.TaskConfigDiff{}
	taskConfigDiff.ChangeMemory(nil, aws.String("512"))
	assert.False(t, taskConfigDiff.Empty())
	assert.Equal(t, "memory was: <nil> and now is: \"512\"", taskConfigDiff.String())
}

func Test_TaskConfigDiff_ContainerDefinitions_NonChanges(t *testing.T) {
	taskConfigDiff := &models.TaskConfigDiff{}
	containerConfigDiff := &models.ContainerConfigDiff{}
	taskConfigDiff.ChangeContainer("container", containerConfigDiff)
	assert.True(t, taskConfigDiff.Empty())
	assert.Equal(t, "", taskConfigDiff.String())
}

func Test_TaskConfigDiff_ContainerDefinitions_WithChanges(t *testing.T) {
	taskConfigDiff := &models.TaskConfigDiff{}
	containerConfigDiff := &models.ContainerConfigDiff{}
	containerConfigDiff.ChangeEnvironment("variable", aws.String("oldValue"), aws.String("newValue"))
	taskConfigDiff.ChangeContainer("container", containerConfigDiff)
	assert.False(t, taskConfigDiff.Empty())
	assert.Equal(t, "the container definition \"container\" changed in this way:\nenvironment variable \"variable\" was: \"oldValue\" and now is: \"newValue\"\n", taskConfigDiff.String())
}
