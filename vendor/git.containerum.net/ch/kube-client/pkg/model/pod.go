package model

type Pod struct {
	Name            string             `json:"name" binding:"required"`
	Containers      []Container        `json:"containers" binding:"dive"`
	ImagePullSecret *map[string]string `json:"image_pull_secret,omitempty"`
	Status          *PodStatus         `json:"status,omitempty"`
	Hostname        *string            `json:"hostname,omitempty"`
}

type PodStatus struct {
	Phase        string `json:"phase"`
	RestartCount int    `json:"restart_count"`
	StartAt      int64  `json:"start_at"`
}

type UpdateImage struct {
	Container string `json:"container_name" binding:"required"`
	Image     string `json:"image" binding:"required"`
}
