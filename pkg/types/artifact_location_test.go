package types

import (
	"github.com/project-interstellar/workflow-watcher/internal/test_utils"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestShouldRenderCorrectJson_ArtifactLocation(t *testing.T) {
	location := ArtifactLocation{Bucket: "bucket", Key: "key"}
	res := test_utils.MarshalToMap(t, location)

	assert.Equal(t, location.Bucket, res["bucket"])
	assert.Equal(t, location.Key, res["key"])
}
