package model

// PodsList -- model for pods list
//
// swagger:model
type PodsList struct {
	Pods []Pod `json:"pods"`
}

// Pod -- model for pod
//
// swagger:model
type Pod struct {
	//creation date in RFC3339 format
	CreatedAt       *string            `json:"created_at,omitempty"`
	Name            string             `json:"name"`
	Owner           string             `json:"owner"`
	Containers      []Container        `json:"containers"`
	ImagePullSecret *map[string]string `json:"image_pull_secret,omitempty"`
	Status          *PodStatus         `json:"status,omitempty"`
	Deploy          *string            `json:"deploy,omitempty"`
	//total CPU usage by all containers in this pod
	TotalCPU uint `json:"total_cpu,omitempty"`
	//total RAM usage by all containers in this pod
	TotalMemory uint `json:"total_memory,omitempty"`
}

// PodStatus -- kubernetes status of pod
//
// swagger:model
type PodStatus struct {
	Phase        string `json:"phase"`
	RestartCount int    `json:"restart_count"`
	//pod start date in RFC3339 format
	StartAt string `json:"start_at"`
}

// UpdateImage -- model for update container image request
//
// swagger:model
type UpdateImage struct {
	// required: true
	Container string `json:"container_name"`
	// required: true
	Image string `json:"image"`
}

// Mask removes information not interesting for users
func (pod *Pod) Mask() {
	pod.Owner = ""
}
