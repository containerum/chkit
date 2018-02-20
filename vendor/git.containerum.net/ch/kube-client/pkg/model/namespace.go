package model

type Namespace struct {
	Name      string    `json:"name" binding:"required"`
	Owner     *string   `json:"owner,omitempty"`
	Resources Resources `json:"resources"`
}

type Resources struct {
	Hard Resource  `json:"hard"`
	Used *Resource `json:"used,omitempty"`
}

type Resource struct {
	CPU    string `json:"cpu"`
	Memory string `json:"memory"`
}
