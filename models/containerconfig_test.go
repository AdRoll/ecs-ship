package models_test

import (
	"testing"

	"github.com/adroll/ecs-ship/models"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ecs/types"
	"github.com/stretchr/testify/assert"
)

func Test_ContainerConfig_ApplyTo_Empty(t *testing.T) {
	containerConfig := &models.ContainerConfig{}
	containerDefinition := &types.ContainerDefinition{}
	_, diff := containerConfig.ApplyTo(containerDefinition)
	assert.True(t, diff.Empty())
}

func Test_ContainerConfig_ApplyTo_CPU(t *testing.T) {
	var newCpu int32 = 100
	containerConfig := &models.ContainerConfig{CPU: &newCpu}
	containerDefinition := &types.ContainerDefinition{}
	newDefinition, diff := containerConfig.ApplyTo(containerDefinition)
	assert.Equal(t, newDefinition.Cpu, newCpu)
	assert.False(t, diff.Empty())
	assert.Equal(t, "CPU was: 0 and now is: 100", diff.String())
}

func Test_ContainerConfig_ApplyTo_Environment(t *testing.T) {
	newEnv := map[string]string{"key": "value"}
	containerConfig := &models.ContainerConfig{Environment: newEnv}
	containerDefinition := &types.ContainerDefinition{}
	newDefinition, diff := containerConfig.ApplyTo(containerDefinition)
	assert.Equal(t, len(newDefinition.Environment), 1)
	assert.False(t, diff.Empty())
	assert.Equal(t, "environment variable \"key\" was: <nil> and now is: \"value\"", diff.String())
}

func Test_ContainerConfig_ApplyTo_ExistingEnvironment(t *testing.T) {
	newEnv := map[string]string{"key": "value"}
	containerConfig := &models.ContainerConfig{Environment: newEnv}
	containerDefinition := &types.ContainerDefinition{Environment: []types.KeyValuePair{{Name: aws.String("key"), Value: aws.String("oldValue")}}}
	newDefinition, diff := containerConfig.ApplyTo(containerDefinition)
	assert.Equal(t, len(newDefinition.Environment), 1)
	assert.False(t, diff.Empty())
	assert.Equal(t, "environment variable \"key\" was: \"oldValue\" and now is: \"value\"", diff.String())
}

func Test_ContainerConfig_ApplyTo_Image(t *testing.T) {
	newImage := "newImage"
	containerConfig := &models.ContainerConfig{Image: &newImage}
	containerDefinition := &types.ContainerDefinition{}
	newDefinition, diff := containerConfig.ApplyTo(containerDefinition)
	assert.Equal(t, *newDefinition.Image, newImage)
	assert.False(t, diff.Empty())
	assert.Equal(t, "image was: <nil> and now is: \"newImage\"", diff.String())
}

func Test_ContainerConfig_ApplyTo_Memory(t *testing.T) {
	var newMemory int32 = 100
	containerConfig := &models.ContainerConfig{Memory: &newMemory}
	containerDefinition := &types.ContainerDefinition{}
	newDefinition, diff := containerConfig.ApplyTo(containerDefinition)
	assert.NotNil(t, newDefinition.Memory)
	assert.Equal(t, *newDefinition.Memory, newMemory)
	assert.False(t, diff.Empty())
	assert.Equal(t, "memory was: <nil> and now is: 100", diff.String())
}

func Test_ContainerConfig_ApplyTo_MemoryReservation(t *testing.T) {
	var newMemoryReservation int32 = 100
	containerConfig := &models.ContainerConfig{MemoryReservation: &newMemoryReservation}
	containerDefinition := &types.ContainerDefinition{}
	newDefinition, diff := containerConfig.ApplyTo(containerDefinition)
	assert.NotNil(t, newDefinition.MemoryReservation)
	assert.Equal(t, *newDefinition.MemoryReservation, newMemoryReservation)
	assert.False(t, diff.Empty())
	assert.Equal(t, "memoryReservation was: <nil> and now is: 100", diff.String())
}
