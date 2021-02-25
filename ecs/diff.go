package ecs

import (
	"fmt"
	"strings"
)

// IntegerDiff represents a difference on an integer value
type IntegerDiff struct {
	was   *int64
	isNow *int64
}

// StringDiff represents a difference on an integer value
type StringDiff struct {
	was   *string
	isNow *string
}

// ContainerConfigDiff all the changes in a task definition
type ContainerConfigDiff struct {
	CPU               *IntegerDiff           `json:"cpu"`
	Environment       map[string]*StringDiff `json:"environment"`
	Image             *StringDiff            `json:"image"`
	Memory            *IntegerDiff           `json:"memory"`
	MemoryReservation *IntegerDiff           `json:"memoryReservation"`
}

// TaskConfigDiff all the changes in a task definition
type TaskConfigDiff struct {
	CPU                  *StringDiff                     `json:"cpu"`
	Memory               *StringDiff                     `json:"memory"`
	ContainerDefinitions map[string]*ContainerConfigDiff `json:"containerDefinitions"`
}

// Empty check if there's no change on the value
func (diff *IntegerDiff) Empty() bool {
	if diff == nil || diff.was == diff.isNow {
		return true
	}
	if diff.was == nil || diff.isNow == nil {
		return false
	}
	return *diff.was == *diff.isNow
}

func (diff *IntegerDiff) change(was *int64, isNow *int64) {
	diff.was = was
	diff.isNow = isNow
}

func (diff *IntegerDiff) String() string {
	if diff.Empty() {
		return ""
	}
	if diff.was == nil && diff.isNow != nil {
		return fmt.Sprintf("was: <nil> and is now: %d", *diff.isNow)
	}
	if diff.was != nil && diff.isNow == nil {
		return fmt.Sprintf("was: %d and now is: <nil>", *diff.was)
	}
	return fmt.Sprintf("was: %d and now is: %d", *diff.was, *diff.isNow)
}

// Empty check if there's no change on the value
func (diff *StringDiff) Empty() bool {
	if diff == nil || diff.was == diff.isNow {
		return true
	}
	if diff.was == nil || diff.isNow == nil {
		return false
	}
	return *diff.was == *diff.isNow
}

func (diff *StringDiff) change(was *string, isNow *string) {
	diff.was = was
	diff.isNow = isNow
}

func (diff *StringDiff) String() string {
	if diff.Empty() {
		return ""
	}
	if diff.was == nil && diff.isNow != nil {
		return fmt.Sprintf("was: <nil> and is now: \"%s\"", *diff.isNow)
	}
	if diff.was != nil && diff.isNow == nil {
		return fmt.Sprintf("was: \"%s\" and is now: <nil>", *diff.was)
	}
	return fmt.Sprintf("was: \"%s\"  and is now: \"%s\"", *diff.was, *diff.isNow)
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
func (diff *ContainerConfigDiff) ChangeCPU(was *int64, isNow *int64) {
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
func (diff *ContainerConfigDiff) ChangeMemory(was *int64, isNow *int64) {
	if diff.Memory == nil {
		diff.Memory = &IntegerDiff{}
	}
	diff.Memory.change(was, isNow)
}

// ChangeMemoryReservation register a change in memory reservation
func (diff *ContainerConfigDiff) ChangeMemoryReservation(was *int64, isNow *int64) {
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

// Empty check if there's no change on the task definition config
func (diff *TaskConfigDiff) Empty() bool {
	cpuChanged := !diff.CPU.Empty()
	memoryChanged := !diff.Memory.Empty()
	if cpuChanged || memoryChanged {
		return false
	}
	for _, diff := range diff.ContainerDefinitions {
		if !diff.Empty() {
			return false
		}
	}
	return true
}

// ChangeCPU register a change in cpu
func (diff *TaskConfigDiff) ChangeCPU(was *string, isNow *string) {
	if diff.CPU == nil {
		diff.CPU = &StringDiff{}
	}
	diff.CPU.change(was, isNow)
}

// ChangeMemory register a change in memory
func (diff *TaskConfigDiff) ChangeMemory(was *string, isNow *string) {
	if diff.Memory == nil {
		diff.Memory = &StringDiff{}
	}
	diff.Memory.change(was, isNow)
}

// ChangeContainer updates a container diff
func (diff *TaskConfigDiff) ChangeContainer(name string, containerDiff *ContainerConfigDiff) {
	if diff.ContainerDefinitions == nil {
		diff.ContainerDefinitions = make(map[string]*ContainerConfigDiff)
	}
	diff.ContainerDefinitions[name] = containerDiff
}

func (diff *TaskConfigDiff) String() string {
	var parts []string
	if !diff.CPU.Empty() {
		parts = append(parts, fmt.Sprintf("cpu %s", diff.CPU))
	}
	if !diff.Memory.Empty() {
		parts = append(parts, fmt.Sprintf("memory %s", diff.Memory))
	}
	for name, diff := range diff.ContainerDefinitions {
		parts = append(parts, fmt.Sprintf("the container definition %s changed in this way:\n%s\n", name, diff))
	}
	return strings.Join(parts, "\n")
}
