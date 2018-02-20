package model

type Pod struct {
	Name            string            `json:"name" binding:"required"`
	Owner           *string           `json:"owner_id,omitempty"`
	Containers      []Container       `json:"containers"`
	ImagePullSecret map[string]string `json:"image_pull_secret,omitempty"`
	Status          *PodStatus        `json:"status,omitempty"`
	Hostname        *string           `json:"hostname,omitempty"`
}

type PodStatus struct {
	Phase string `json:"phase"`
}

type Container struct {
	Name   string    `json:"name" binding:"required"`
	Env    *[]Env    `json:"env,omitempty"`
	Image  string    `json:"image" binding:"required"`
	Volume *[]Volume `json:"volume,omitempty"`
}

type Env struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type Volume struct {
	Name      string  `json:"name"`
	MountPath string  `json:"mount_path"`
	SubPath   *string `json:"sub_path,omitempty"`
}
