package types

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTypeReturnsCorrectValue_ArtifactType_Input(t *testing.T) {
	assert.Equal(t, string(Input), "input")
}

func TestTypeReturnsCorrectValue_ArtifactType_Output(t *testing.T) {
	assert.Equal(t, string(Output), "output")
}

func TestTypeReturnsCorrectValue_ArtifactType_Logs(t *testing.T) {
	assert.Equal(t, string(Logs), "logs")
}