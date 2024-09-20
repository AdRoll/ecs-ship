package ecs

import (
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/ecs"
	"github.com/aws/aws-sdk-go-v2/service/ecs/types"

	"github.com/stretchr/testify/assert"
)

func Test_CopyTaskDefinition(t *testing.T) {
	tagKey := "abc"
	tagVal := "123"
	cpu := "100"
	task := &ecs.DescribeTaskDefinitionOutput{
		TaskDefinition: &types.TaskDefinition{
			Cpu: &cpu,
		},
		Tags: []types.Tag{
			{Key: &tagKey, Value: &tagVal},
		},
	}

	res := copyTaskDefinition(task)

	assert.Equal(t, len(res.Tags), 1)
	assert.Equal(t, res.Tags[0].Key, &tagKey)
	assert.Equal(t, res.Tags[0].Value, &tagVal)
	assert.Equal(t, res.Cpu, &cpu)
}
