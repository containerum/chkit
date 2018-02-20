package model

type Deployment struct {
	Name            string             `json:"name" binding:"required"`
	Replicas        int                `json:"replicas" binding:"required"`
	Containers      []Container        `json:"containers" binding:"required,dive"`
	ImagePullSecret *map[string]string `json:"image_pull_secret,omitempty"`
	Status          *DeploymentStatus  `json:"status,omitempty"`
	Hostname        *string            `json:"hostname,omitempty"`
}

type DeploymentStatus struct {
	CreatedAt           int64 `json:"created_at"`
	UpdatedAt           int64 `json:"updated_at"`
	Replicas            int   `json:"replicas"`
	ReadyReplicas       int   `json:"ready_replicas"`
	AvailableReplicas   int   `json:"available_replicas"`
	UnavailableReplicas int   `json:"unavailable_replicas"`
	UpdatedReplicas     int   `json:"updated_replicas"`
}

type UpdateReplicas struct {
	Replicas int `json:"replicas" binding:"required"`
}

type ResourceDeployment struct {
	Containers []Container       `json:"containers"`
	Owner      *string           `json:"owner,omitempty"`
	Labels     map[string]string `json:"labels"`
	Name       string            `json:"name"`
	Replicas   int               `json:"replicas"`
}

type ResourceContainer struct {
	Image     string `json:"image"`
	Name      string `json:"name"`
	Resources struct {
		Requests Resource `json:"requests"`
	} `json:"resources"`
	Env          []ResourceEnv         `json:"env"`
	Commands     []string              `json:"commands"`
	Ports        []ResourcePort        `json:"ports"`
	VolumeMounts []ResourceVolumeMount `json:"volumeMounts"`
}

type ResourceEnv struct {
	Value string `json:"value"`
	Name  string `json:"name"`
}

type ResourcePort struct {
	ContainerPort int `json:"containerPort"`
}

type ResourceVolumeMount struct {
	Name      string `json:"name"`
	MountPath string `json:"mountPath"`
	SubPath   string `json:"subPath"`
}
