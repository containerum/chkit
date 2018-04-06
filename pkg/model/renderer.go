package model

type Renderer interface {
	TableRenderer
	YAMLrenderer
	JSONrenderer
}
