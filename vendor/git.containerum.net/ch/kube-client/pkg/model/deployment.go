package model

type Deployment struct {
	Name            string            `json:"name" binding:"required"`
	Owner           *string           `json:"owner_id,omitempty"`
	Replicas        int               `json:"replicas" binding:"required"`
	Containers      []Container       `json:"containers" binding:"required"`
	ImagePullSecret map[string]string `json:"image_pull_secret,omitempty"`
	Status          *DeploymentStatus `json:"status,omitempty"`
	Hostname        *string           `json:"hostname,omitempty"`
}

type DeploymentStatus struct {
	Created             int64 `json:"created_at"`
	Updated             int64 `json:"updated_at"`
	Replicas            int   `json:"replicas"`
	ReadyReplicas       int   `json:"ready_replicas"`
	AvailableReplicas   int   `json:"available_replicas"`
	UnavailableReplicas int   `json:"unavailable_replicas"`
	UpdatedReplicas     int   `json:"updated_replicas"`
}
