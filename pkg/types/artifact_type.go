package types

type ArtifactType int

const (
	Input  ArtifactType = 0
	Output ArtifactType = 1
	Logs   ArtifactType = 2
)
