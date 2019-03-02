package test_utils

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"testing"
)

func MarshalToMap(t *testing.T, unmarshalled interface{}) map[string]interface{} {
	bytes, err := json.Marshal(unmarshalled)
	assert.Nil(t, err)

	res := make(map[string]interface{})
	encodingErr := json.Unmarshal(bytes, &res)

	assert.Nil(t, encodingErr)
	return res
}
