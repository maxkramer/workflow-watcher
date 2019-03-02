package types

type ArtifactType string

const (
	Input  ArtifactType = "input"
	Output ArtifactType = "output"
	Logs   ArtifactType = "logs"
)
