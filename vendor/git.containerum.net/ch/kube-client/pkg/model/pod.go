package model

// Pod --
type Pod struct {
	Name            string             `json:"name"`
	Containers      []Container        `json:"containers"`
	ImagePullSecret *map[string]string `json:"image_pull_secret,omitempty"`
	Status          *PodStatus         `json:"status,omitempty"`
	Hostname        *string            `json:"hostname,omitempty"`
}

// PodStatus --
type PodStatus struct {
	Phase        string `json:"phase"`
	RestartCount int    `json:"restart_count"`
	StartAt      int64  `json:"start_at"`
}

// UpdateImage --
type UpdateImage struct {
	Container string `json:"container_name"`
	Image     string `json:"image"`
}
