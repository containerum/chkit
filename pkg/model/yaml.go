package model

type YAMLrenderer interface {
	RenderYAML() (string, error)
}
