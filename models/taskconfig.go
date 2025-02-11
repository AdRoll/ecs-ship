package models

import (
	"github.com/aws/aws-sdk-go-v2/service/ecs"
	"github.com/aws/aws-sdk-go-v2/service/ecs/types"
)

// TaskConfig represents changes we can make to task definitions
type TaskConfig struct {
	CPU                  *string                    `json:"cpu" yaml:"cpu"`
	Memory               *string                    `json:"memory" yaml:"memory"`
	ContainerDefinitions map[string]ContainerConfig `json:"containerDefinitions" yaml:"containerDefinitions"`
}

// ApplyTo apply a config to register task definition input
func (config *TaskConfig) ApplyTo(input *ecs.RegisterTaskDefinitionInput) (*ecs.RegisterTaskDefinitionInput, *TaskConfigDiff) {
	diff := &TaskConfigDiff{}
	newInput := &ecs.RegisterTaskDefinitionInput{
		ContainerDefinitions:    input.ContainerDefinitions,
		Family:                  input.Family,
		Cpu:                     input.Cpu,
		EphemeralStorage:        input.EphemeralStorage,
		ExecutionRoleArn:        input.ExecutionRoleArn,
		InferenceAccelerators:   input.InferenceAccelerators,
		IpcMode:                 input.IpcMode,
		Memory:                  input.Memory,
		NetworkMode:             input.NetworkMode,
		PidMode:                 input.PidMode,
		PlacementConstraints:    input.PlacementConstraints,
		ProxyConfiguration:      input.ProxyConfiguration,
		RequiresCompatibilities: input.RequiresCompatibilities,
		RuntimePlatform:         input.RuntimePlatform,
		Tags:                    input.Tags,
		TaskRoleArn:             input.TaskRoleArn,
		Volumes:                 input.Volumes,
	}
	updateString(newInput.Cpu, config.CPU, func(val string) { newInput.Cpu = &val }, diff.ChangeCPU)
	updateString(newInput.Memory, config.Memory, func(val string) { newInput.Memory = &val }, diff.ChangeMemory)

	// Update container definitions
	newDefs := make([]types.ContainerDefinition, 0, len(newInput.ContainerDefinitions))
	for _, definition := range newInput.ContainerDefinitions {
		if config, ok := config.ContainerDefinitions[*definition.Name]; ok {
			newDef, newDiff := config.ApplyTo(&definition)
			newDefs = append(newDefs, newDef)
			diff.ChangeContainer(*definition.Name, newDiff)
		} else {
			newDefs = append(newDefs, definition)
		}
	}
	newInput.ContainerDefinitions = newDefs

	return newInput, diff
}
