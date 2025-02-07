package models

import (
	"fmt"
	"strings"
)

// TaskConfigDiff all the changes in a task definition
type TaskConfigDiff struct {
	CPU                  *StringDiff                     `json:"cpu"`
	Memory               *StringDiff                     `json:"memory"`
	ContainerDefinitions map[string]*ContainerConfigDiff `json:"containerDefinitions"`
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
