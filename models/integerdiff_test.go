package models_test

import (
	"testing"

	"github.com/adroll/ecs-ship/models"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/stretchr/testify/assert"
)

func Test_IntegerDiff_Empty(t *testing.T) {
	integerDiff := &models.IntegerDiff{}
	assert.True(t, integerDiff.Empty())
}

func Test_IntegerDiff_Empty_WasNil(t *testing.T) {
	integerDiff := &models.IntegerDiff{}
	integerDiff.Change(nil, aws.Int32(10))
	assert.False(t, integerDiff.Empty())
	assert.Equal(t, "was: <nil> and now is: 10", integerDiff.String())
}

func Test_IntegerDiff_Empty_IsNowNil(t *testing.T) {
	integerDiff := &models.IntegerDiff{}
	integerDiff.Change(aws.Int32(10), nil)
	assert.False(t, integerDiff.Empty())
	assert.Equal(t, "was: 10 and now is: <nil>", integerDiff.String())
}

func Test_IntegerDiff_Empty_WasAndIsNowNil(t *testing.T) {
	integerDiff := &models.IntegerDiff{}
	integerDiff.Change(nil, nil)
	assert.True(t, integerDiff.Empty())
}

func Test_IntegerDiff_Empty_WasAndIsNowNotNil(t *testing.T) {
	integerDiff := &models.IntegerDiff{}
	integerDiff.Change(aws.Int32(10), aws.Int32(10))
	assert.True(t, integerDiff.Empty())
}

func Test_IntegerDiff_Empty_WasAndIsNowNotNil_DifferentValues(t *testing.T) {
	integerDiff := &models.IntegerDiff{}
	integerDiff.Change(aws.Int32(10), aws.Int32(20))
	assert.False(t, integerDiff.Empty())
	assert.Equal(t, "was: 10 and now is: 20", integerDiff.String())
}
