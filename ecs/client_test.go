package ecs

import (
	"testing"

	"github.com/aws/aws-sdk-go/service/ecs"

	"github.com/stretchr/testify/assert"
)

func Test_CopyTaskDefinition(t *testing.T) {
	tagKey := "abc"
	tagVal := "123"
	cpu := "100"
	task := &ecs.DescribeTaskDefinitionOutput{
		TaskDefinition: &ecs.TaskDefinition{
			Cpu: &cpu,
		},
		Tags: []*ecs.Tag{
			{Key: &tagKey, Value: &tagVal},
		},
	}

	res := copyTaskDefinition(task)

	assert.Equal(t, len(res.Tags), 1)
	assert.Equal(t, res.Tags[0].Key, &tagKey)
	assert.Equal(t, res.Tags[0].Value, &tagVal)
	assert.Equal(t, res.Cpu, &cpu)
}
