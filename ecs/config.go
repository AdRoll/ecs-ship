package ecs

import (
	"github.com/aws/aws-sdk-go/service/ecs"
)

// ContainerConfig represents changes we can make to containers
type ContainerConfig struct {
	CPU               *int64            `json:"cpu" yaml:"cpu"`
	Environment       map[string]string `json:"environment" yaml:"environment"`
	Image             *string           `json:"image" yaml:"image"`
	Memory            *int64            `json:"memory" yaml:"memory"`
	MemoryReservation *int64            `json:"memoryReservation" yaml:"memoryReservation"`
}

// ApplyTo apply a config to a container definition
func (config *ContainerConfig) ApplyTo(input *ecs.ContainerDefinition) (*ecs.ContainerDefinition, *ContainerConfigDiff) {
	diff := &ContainerConfigDiff{}
	newDef := &ecs.ContainerDefinition{
		Command:                input.Command,
		Cpu:                    input.Cpu,
		DependsOn:              input.DependsOn,
		DisableNetworking:      input.DisableNetworking,
		DnsSearchDomains:       input.DnsSearchDomains,
		DnsServers:             input.DnsServers,
		DockerLabels:           input.DockerLabels,
		DockerSecurityOptions:  input.DockerSecurityOptions,
		Environment:            input.Environment,
		EnvironmentFiles:       input.EnvironmentFiles,
		Essential:              input.Essential,
		ExtraHosts:             input.ExtraHosts,
		HealthCheck:            input.HealthCheck,
		Hostname:               input.Hostname,
		Image:                  input.Image,
		Interactive:            input.Interactive,
		Links:                  input.Links,
		LinuxParameters:        input.LinuxParameters,
		LogConfiguration:       input.LogConfiguration,
		Memory:                 input.Memory,
		MemoryReservation:      input.MemoryReservation,
		MountPoints:            input.MountPoints,
		Name:                   input.Name,
		PortMappings:           input.PortMappings,
		Privileged:             input.Privileged,
		PseudoTerminal:         input.PseudoTerminal,
		ReadonlyRootFilesystem: input.ReadonlyRootFilesystem,
		RepositoryCredentials:  input.RepositoryCredentials,
		ResourceRequirements:   input.ResourceRequirements,
		Secrets:                input.Secrets,
		StartTimeout:           input.StartTimeout,
		StopTimeout:            input.StopTimeout,
		SystemControls:         input.SystemControls,
		Ulimits:                input.Ulimits,
		User:                   input.User,
		VolumesFrom:            input.VolumesFrom,
		WorkingDirectory:       input.WorkingDirectory,
	}
	updateInt(newDef.Cpu, config.CPU, func(val int64) { newDef.SetCpu(val) }, diff.ChangeCPU)
	updateString(newDef.Image, config.Image, func(val string) { newDef.SetImage(val) }, diff.ChangeImage)
	updateInt(newDef.Memory, config.Memory, func(val int64) { newDef.SetMemory(val) }, diff.ChangeMemory)
	updateInt(newDef.MemoryReservation, config.MemoryReservation, func(val int64) { newDef.SetMemoryReservation(val) }, diff.ChangeMemoryReservation)

	newEnvironment := make([]*ecs.KeyValuePair, 0)
	used := make(map[string]struct{})
	usedFlag := struct{}{}
	// Update existing environment variables
	for _, pair := range newDef.Environment {
		if value, prs := config.Environment[*pair.Name]; prs {
			valueCopy := value[:]
			newEnvironment = append(newEnvironment, &ecs.KeyValuePair{Name: pair.Name, Value: &valueCopy})
			diff.ChangeEnvironment(*pair.Name, pair.Value, &valueCopy)
			used[*pair.Name] = usedFlag
		} else {
			newEnvironment = append(newEnvironment, pair)
		}
	}

	// Create new environment variables
	for name, value := range config.Environment {
		if _, prs := used[name]; prs {
			continue
		}
		nameCopy := name[:]
		valueCopy := value[:]
		newEnvironment = append(newEnvironment, &ecs.KeyValuePair{Name: &nameCopy, Value: &valueCopy})
	 diff.ChangeEnvironment(name, nil, &valueCopy)
	}
	newDef.SetEnvironment(newEnvironment)

	return newDef, diff
}

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
		Cpu:                     input.Cpu,
		ExecutionRoleArn:        input.ExecutionRoleArn,
		Family:                  input.Family,
		InferenceAccelerators:   input.InferenceAccelerators,
		IpcMode:                 input.IpcMode,
		Memory:                  input.Memory,
		NetworkMode:             input.NetworkMode,
		PidMode:                 input.PidMode,
		PlacementConstraints:    input.PlacementConstraints,
		ProxyConfiguration:      input.ProxyConfiguration,
		RequiresCompatibilities: input.RequiresCompatibilities,
		Tags:                    input.Tags,
		TaskRoleArn:             input.TaskRoleArn,
		Volumes:                 input.Volumes,
	}
	updateString(newInput.Cpu, config.CPU, func(val string) { newInput.SetCpu(val) }, diff.ChangeCPU)
	updateString(newInput.Memory, config.Memory, func(val string) { newInput.SetMemory(val) }, diff.ChangeCPU)

	// Update container definitions
	newDefs := make([]*ecs.ContainerDefinition, 0, len(newInput.ContainerDefinitions))
	for _, definition := range newInput.ContainerDefinitions {
		if config, prs := config.ContainerDefinitions[*definition.Name]; prs {
			newDef, newDiff := config.ApplyTo(definition)
			newDefs = append(newDefs, newDef)
			diff.ChangeContainer(*definition.Name, newDiff)
		} else {
			newDefs = append(newDefs, definition)
		}
	}
	newInput.SetContainerDefinitions(newDefs)

	return newInput, diff
}

func updateString(old *string, new *string, apply func(string), record func(*string, *string)) {
	if old == nil && new == nil || new == nil {
		return
	}
	apply(*new)
	record(old, new)
}

func updateInt(old *int64, new *int64, apply func(int64), record func(*int64, *int64)) {
	if old == nil && new == nil || new == nil {
		return
	}
	apply(*new)
	record(old, new)
}
