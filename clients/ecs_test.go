package clients_test

import (
	"testing"

	"github.com/adroll/ecs-ship/clients"
	"github.com/aws/aws-sdk-go-v2/service/ecs"
	ecsTypes "github.com/aws/aws-sdk-go-v2/service/ecs/types"
	"github.com/stretchr/testify/assert"
)

func Test_CopyTaskDefinition(t *testing.T) {
	client := clients.NewECSClient(nil)
	tagKey := "abc"
	tagVal := "123"
	cpu := "100"
	output := &ecs.DescribeTaskDefinitionOutput{
		TaskDefinition: &ecsTypes.TaskDefinition{
			Cpu: &cpu,
		},
		Tags: []ecsTypes.Tag{
			{Key: &tagKey, Value: &tagVal},
		},
	}

	input := client.CopiedTaskDefinition(output)

	assert.Equal(t, len(input.Tags), 1)
	assert.Equal(t, input.Tags[0].Key, &tagKey)
	assert.Equal(t, input.Tags[0].Value, &tagVal)
	assert.Equal(t, input.Cpu, &cpu)
}
