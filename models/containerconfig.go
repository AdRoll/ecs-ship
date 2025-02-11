package models

import "github.com/aws/aws-sdk-go-v2/service/ecs/types"

// ContainerConfig represents changes we can make to containers
type ContainerConfig struct {
	CPU               *int32            `json:"cpu" yaml:"cpu"`
	Environment       map[string]string `json:"environment" yaml:"environment"`
	Image             *string           `json:"image" yaml:"image"`
	Memory            *int32            `json:"memory" yaml:"memory"`
	MemoryReservation *int32            `json:"memoryReservation" yaml:"memoryReservation"`
}

// ApplyTo apply a config to a container definition
func (config *ContainerConfig) ApplyTo(input *types.ContainerDefinition) (types.ContainerDefinition, *ContainerConfigDiff) {
	diff := &ContainerConfigDiff{}
	newDef := types.ContainerDefinition{
		Command:                input.Command,
		Cpu:                    input.Cpu,
		CredentialSpecs:        input.CredentialSpecs,
		DependsOn:              input.DependsOn,
		DisableNetworking:      input.DisableNetworking,
		DnsSearchDomains:       input.DnsSearchDomains,
		DnsServers:             input.DnsServers,
		DockerLabels:           input.DockerLabels,
		DockerSecurityOptions:  input.DockerSecurityOptions,
		EntryPoint:             input.EntryPoint,
		Environment:            input.Environment,
		EnvironmentFiles:       input.EnvironmentFiles,
		Essential:              input.Essential,
		ExtraHosts:             input.ExtraHosts,
		FirelensConfiguration:  input.FirelensConfiguration,
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
	updateInt(&newDef.Cpu, config.CPU, func(val int32) { newDef.Cpu = val }, diff.ChangeCPU)
	updateString(newDef.Image, config.Image, func(val string) { newDef.Image = &val }, diff.ChangeImage)
	// FIXME: We should have UpdateIntPtr instead
	updateInt(newDef.Memory, config.Memory, func(val int32) { newDef.Memory = &val }, diff.ChangeMemory)
	updateInt(newDef.MemoryReservation, config.MemoryReservation, func(val int32) { newDef.MemoryReservation = &val }, diff.ChangeMemoryReservation)

	newEnvironment := make([]types.KeyValuePair, 0)
	used := make(map[string]struct{})
	usedFlag := struct{}{}
	// Update existing environment variables
	for _, pair := range newDef.Environment {
		if value, prs := config.Environment[*pair.Name]; prs {
			valueCopy := value[:]
			newEnvironment = append(newEnvironment, types.KeyValuePair{Name: pair.Name, Value: &valueCopy})
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
		newEnvironment = append(newEnvironment, types.KeyValuePair{Name: &nameCopy, Value: &valueCopy})
		diff.ChangeEnvironment(name, nil, &valueCopy)
	}
	newDef.Environment = newEnvironment

	return newDef, diff
}
