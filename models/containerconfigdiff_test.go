package models_test

import (
	"testing"

	"github.com/adroll/ecs-ship/models"
	"github.com/stretchr/testify/assert"
)

func Test_ContainerConfigDiff_Empty(t *testing.T) {
	containerConfigDiff := &models.ContainerConfigDiff{}
	assert.True(t, containerConfigDiff.Empty())
}

func Test_ContainerConfigDiff_CPU(t *testing.T) {
	var oldCpu int32 = 0
	var newCpu int32 = 100
	containerConfigDiff := &models.ContainerConfigDiff{}
	containerConfigDiff.ChangeCPU(&oldCpu, &newCpu)
	assert.False(t, containerConfigDiff.Empty())
	assert.Equal(t, "CPU was: 0 and now is: 100", containerConfigDiff.String())
}

func Test_ContainerConfigDiff_Image(t *testing.T) {
	oldImage := ""
	newImage := "newImage"
	containerConfigDiff := &models.ContainerConfigDiff{}
	containerConfigDiff.ChangeImage(&oldImage, &newImage)
	assert.False(t, containerConfigDiff.Empty())
	assert.Equal(t, "image was: \"\" and now is: \"newImage\"", containerConfigDiff.String())
}

func Test_ContainerConfigDiff_Memory(t *testing.T) {
	var oldMemory int32 = 0
	var newMemory int32 = 100
	containerConfigDiff := &models.ContainerConfigDiff{}
	containerConfigDiff.ChangeMemory(&oldMemory, &newMemory)
	assert.False(t, containerConfigDiff.Empty())
	assert.Equal(t, "memory was: 0 and now is: 100", containerConfigDiff.String())
}

func Test_ContainerConfigDiff_MemoryReservation(t *testing.T) {
	var oldMemoryReservation int32 = 0
	var newMemoryReservation int32 = 100
	containerConfigDiff := &models.ContainerConfigDiff{}
	containerConfigDiff.ChangeMemoryReservation(&oldMemoryReservation, &newMemoryReservation)
	assert.False(t, containerConfigDiff.Empty())
	assert.Equal(t, "memoryReservation was: 0 and now is: 100", containerConfigDiff.String())
}

func Test_ContainerConfigDiff_Environment(t *testing.T) {
	oldEnv := ""
	newEnv := "value"
	containerConfigDiff := &models.ContainerConfigDiff{}
	containerConfigDiff.ChangeEnvironment("variable", &oldEnv, &newEnv)
	assert.False(t, containerConfigDiff.Empty())
	assert.Equal(t, "environment variable \"variable\" was: \"\" and now is: \"value\"", containerConfigDiff.String())
}
