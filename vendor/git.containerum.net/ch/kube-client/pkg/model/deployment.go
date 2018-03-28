package model

// DeploymentStatus -- represents deployment status
// CreatedAt and UpdatedAt -- Unix timestamps
type DeploymentStatus struct {
	CreatedAt           int64 `json:"created_at"`
	UpdatedAt           int64 `json:"updated_at"`
	Replicas            int   `json:"replicas"`
	ReadyReplicas       int   `json:"ready_replicas"`
	AvailableReplicas   int   `json:"available_replicas"`
	UnavailableReplicas int   `json:"unavailable_replicas"`
	UpdatedReplicas     int   `json:"updated_replicas"`
}

// UpdateReplicas -- contains new number of replicas
type UpdateReplicas struct {
	Replicas int `json:"replicas"`
}

// Deployment --
type Deployment struct {
	Status      *DeploymentStatus `json:"status,omitempty"`
	Containers  []Container       `json:"containers"`
	Labels      map[string]string `json:"labels,omitempty"`
	Name        string            `json:"name"`
	Replicas    int               `json:"replicas"`
	TotalCPU    string            `json:"total_cpu,omitempty"`
	TotalMemory string            `json:"total_memory,omitempty"`
}

// Container --
type Container struct {
	Image        string            `json:"image"`
	Name         string            `json:"name"`
	Limits       Resource          `json:"limits"`
	Env          []Env             `json:"env,omitempty"`
	Commands     []string          `json:"commands,omitempty"`
	Ports        []ContainerPort   `json:"ports,omitempty"`
	VolumeMounts []ContainerVolume `json:"volume_mounts,omitempty"`
	ConfigMaps   []ContainerVolume `json:"config_maps,omitempty"`
}

// Env -- represents key value pair of enviroment variable
type Env struct {
	Value string `json:"value"`
	Name  string `json:"name"`
}

// ContainerPort --
type ContainerPort struct {
	Name     string   `json:"name"`
	Port     int      `json:"port"`
	Protocol Protocol `json:"protocol"`
}

// ContainerVolume --
type ContainerVolume struct {
	Name      string  `json:"name"`
	Mode      *string `json:"mode,omitempty"`
	MountPath string  `json:"mount_path"`
	SubPath   *string `json:"sub_path,omitempty"`
}
