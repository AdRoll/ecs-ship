package models

import (
	"fmt"
	"strings"
)

// ContainerConfigDiff all the changes in a task definition
type ContainerConfigDiff struct {
	cpu               *IntegerDiff
	environment       map[string]*StringDiff
	image             *StringDiff
	memory            *IntegerDiff
	memoryReservation *IntegerDiff
}

// Empty check if there's no change on the container config
func (diff *ContainerConfigDiff) Empty() bool {
	cpuChanged := !diff.cpu.Empty()
	imageChanged := !diff.image.Empty()
	memoryChanged := !diff.memory.Empty()
	memoryReservationChanged := !diff.memoryReservation.Empty()
	if cpuChanged || imageChanged || memoryChanged || memoryReservationChanged {
		return false
	}
	for _, diff := range diff.environment {
		if !diff.Empty() {
			return false
		}
	}
	return true
}

func (diff *ContainerConfigDiff) String() string {
	var parts []string
	if !diff.cpu.Empty() {
		parts = append(parts, fmt.Sprintf("CPU %s", diff.cpu))
	}
	if !diff.image.Empty() {
		parts = append(parts, fmt.Sprintf("image %s", diff.image))
	}
	if !diff.memory.Empty() {
		parts = append(parts, fmt.Sprintf("memory %s", diff.memory))
	}
	if !diff.memoryReservation.Empty() {
		parts = append(parts, fmt.Sprintf("memoryReservation %s", diff.memoryReservation))
	}
	for name, diff := range diff.environment {
		if !diff.Empty() {
			parts = append(parts, fmt.Sprintf("environment variable \"%s\" %s", name, diff))
		}
	}
	return strings.Join(parts, "\n")
}

// ChangeCPU register a change in cpu
func (diff *ContainerConfigDiff) ChangeCPU(was *int32, isNow *int32) {
	if diff.cpu == nil {
		diff.cpu = &IntegerDiff{}
	}
	diff.cpu.Change(was, isNow)
}

// ChangeImage register a change in the image
func (diff *ContainerConfigDiff) ChangeImage(was *string, isNow *string) {
	if diff.image == nil {
		diff.image = &StringDiff{}
	}
	diff.image.Change(was, isNow)
}

// ChangeMemory register a change in memory
func (diff *ContainerConfigDiff) ChangeMemory(was *int32, isNow *int32) {
	if diff.memory == nil {
		diff.memory = &IntegerDiff{}
	}
	diff.memory.Change(was, isNow)
}

// ChangeMemoryReservation register a change in memory reservation
func (diff *ContainerConfigDiff) ChangeMemoryReservation(was *int32, isNow *int32) {
	if diff.memoryReservation == nil {
		diff.memoryReservation = &IntegerDiff{}
	}
	diff.memoryReservation.Change(was, isNow)
}

// ChangeEnvironment register a change in environment variables
func (diff *ContainerConfigDiff) ChangeEnvironment(variable string, was *string, isNow *string) {
	if diff.environment == nil {
		diff.environment = make(map[string]*StringDiff)
	}
	if _, ok := diff.environment[variable]; !ok {
		diff.environment[variable] = &StringDiff{}
	}
	diff.environment[variable].Change(was, isNow)
}
