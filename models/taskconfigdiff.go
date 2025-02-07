package models

import (
	"fmt"
	"strings"
)

// TaskConfigDiff all the changes in a task definition
type TaskConfigDiff struct {
	cpu                  *StringDiff
	memory               *StringDiff
	containerDefinitions map[string]*ContainerConfigDiff
}

// Empty check if there's no change on the task definition config
func (diff *TaskConfigDiff) Empty() bool {
	cpuChanged := !diff.cpu.Empty()
	memoryChanged := !diff.memory.Empty()
	if cpuChanged || memoryChanged {
		return false
	}
	for _, diff := range diff.containerDefinitions {
		if !diff.Empty() {
			return false
		}
	}
	return true
}

// ChangeCPU register a change in cpu
func (diff *TaskConfigDiff) ChangeCPU(was *string, isNow *string) {
	if diff.cpu == nil {
		diff.cpu = &StringDiff{}
	}
	diff.cpu.Change(was, isNow)
}

// ChangeMemory register a change in memory
func (diff *TaskConfigDiff) ChangeMemory(was *string, isNow *string) {
	if diff.memory == nil {
		diff.memory = &StringDiff{}
	}
	diff.memory.Change(was, isNow)
}

// ChangeContainer updates a container diff
func (diff *TaskConfigDiff) ChangeContainer(name string, containerDiff *ContainerConfigDiff) {
	if diff.containerDefinitions == nil {
		diff.containerDefinitions = make(map[string]*ContainerConfigDiff)
	}
	diff.containerDefinitions[name] = containerDiff
}

func (diff *TaskConfigDiff) String() string {
	var parts []string
	if !diff.cpu.Empty() {
		parts = append(parts, fmt.Sprintf("CPU %s", diff.cpu))
	}
	if !diff.memory.Empty() {
		parts = append(parts, fmt.Sprintf("memory %s", diff.memory))
	}
	for name, diff := range diff.containerDefinitions {
		if diff.Empty() {
			continue
		}
		parts = append(parts, fmt.Sprintf("the container definition \"%s\" changed in this way:\n%s\n", name, diff))
	}
	return strings.Join(parts, "\n")
}
