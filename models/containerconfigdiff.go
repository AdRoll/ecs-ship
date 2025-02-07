package models

import (
	"fmt"
	"strings"
)

// ContainerConfigDiff all the changes in a task definition
type ContainerConfigDiff struct {
	CPU               *IntegerDiff           `json:"cpu"`
	Environment       map[string]*StringDiff `json:"environment"`
	Image             *StringDiff            `json:"image"`
	Memory            *IntegerDiff           `json:"memory"`
	MemoryReservation *IntegerDiff           `json:"memoryReservation"`
}

// Empty check if there's no change on the container config
func (diff *ContainerConfigDiff) Empty() bool {
	cpuChanged := !diff.CPU.Empty()
	imageChanged := !diff.Image.Empty()
	memoryChanged := !diff.Memory.Empty()
	memoryReservationChanged := !diff.MemoryReservation.Empty()
	if cpuChanged || imageChanged || memoryChanged || memoryReservationChanged {
		return false
	}
	for _, diff := range diff.Environment {
		if !diff.Empty() {
			return false
		}
	}
	return true
}

func (diff *ContainerConfigDiff) String() string {
	var parts []string
	if !diff.CPU.Empty() {
		parts = append(parts, fmt.Sprintf("cpu %s", diff.CPU))
	}
	if !diff.Image.Empty() {
		parts = append(parts, fmt.Sprintf("image %s", diff.Image))
	}
	if !diff.Memory.Empty() {
		parts = append(parts, fmt.Sprintf("memory %s", diff.Memory))
	}
	if !diff.MemoryReservation.Empty() {
		parts = append(parts, fmt.Sprintf("memoryReservation %s", diff.MemoryReservation))
	}
	for name, diff := range diff.Environment {
		if !diff.Empty() {
			parts = append(parts, fmt.Sprintf("environment variable %s %s", name, diff))
		}
	}
	return strings.Join(parts, "\n")
}

// ChangeCPU register a change in cpu
func (diff *ContainerConfigDiff) ChangeCPU(was *int32, isNow *int32) {
	if diff.CPU == nil {
		diff.CPU = &IntegerDiff{}
	}
	diff.CPU.change(was, isNow)
}

// ChangeImage register a change in the image
func (diff *ContainerConfigDiff) ChangeImage(was *string, isNow *string) {
	if diff.Image == nil {
		diff.Image = &StringDiff{}
	}
	diff.Image.change(was, isNow)
}

// ChangeMemory register a change in memory
func (diff *ContainerConfigDiff) ChangeMemory(was *int32, isNow *int32) {
	if diff.Memory == nil {
		diff.Memory = &IntegerDiff{}
	}
	diff.Memory.change(was, isNow)
}

// ChangeMemoryReservation register a change in memory reservation
func (diff *ContainerConfigDiff) ChangeMemoryReservation(was *int32, isNow *int32) {
	if diff.MemoryReservation == nil {
		diff.MemoryReservation = &IntegerDiff{}
	}
	diff.MemoryReservation.change(was, isNow)
}

// ChangeEnvironment register a change in environment variables
func (diff *ContainerConfigDiff) ChangeEnvironment(variable string, was *string, isNow *string) {
	if diff.Environment == nil {
		diff.Environment = make(map[string]*StringDiff)
	}
	if _, prs := diff.Environment[variable]; !prs {
		diff.Environment[variable] = &StringDiff{}
	}
	diff.Environment[variable].change(was, isNow)
}
