package model

type JSONrenderer interface {
	RenderJSON() (string, error)
}
