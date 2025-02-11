package models_test

import (
	"testing"

	"github.com/adroll/ecs-ship/models"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/stretchr/testify/assert"
)

func Test_StringDiff_Empty(t *testing.T) {
	stringDiff := &models.StringDiff{}
	assert.True(t, stringDiff.Empty())
}

func Test_StringDiff_Empty_WasNil(t *testing.T) {
	stringDiff := &models.StringDiff{}
	stringDiff.Change(nil, aws.String("value"))
	assert.False(t, stringDiff.Empty())
	assert.Equal(t, "was: <nil> and now is: \"value\"", stringDiff.String())
}

func Test_StringDiff_Empty_IsNowNil(t *testing.T) {
	stringDiff := &models.StringDiff{}
	stringDiff.Change(aws.String("value"), nil)
	assert.False(t, stringDiff.Empty())
	assert.Equal(t, "was: \"value\" and now is: <nil>", stringDiff.String())
}

func Test_StringDiff_Empty_WasAndIsNowNil(t *testing.T) {
	stringDiff := &models.StringDiff{}
	stringDiff.Change(nil, nil)
	assert.True(t, stringDiff.Empty())
}

func Test_StringDiff_Empty_WasAndIsNowNotNil(t *testing.T) {
	stringDiff := &models.StringDiff{}
	stringDiff.Change(aws.String("value"), aws.String("value"))
	assert.True(t, stringDiff.Empty())
}

func Test_StringDiff_Empty_WasAndIsNowNotNil_DifferentValues(t *testing.T) {
	stringDiff := &models.StringDiff{}
	stringDiff.Change(aws.String("value"), aws.String("another value"))
	assert.False(t, stringDiff.Empty())
	assert.Equal(t, "was: \"value\" and now is: \"another value\"", stringDiff.String())
}
