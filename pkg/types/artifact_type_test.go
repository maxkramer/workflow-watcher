package types

import (
	"testing"

	"github.com/stretchr/testify/assert"
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
