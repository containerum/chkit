package model

type Secret struct {
	Name      string            `json:"name" binding:"required"`
	CreatedAt *int64            `json:"created_at,omitempty"`
	Data      map[string]string `json:"data" binding:"required"`
}
