package model

// DeploymentStatus -- kubernetes status of deployment
//
// swagger:model
type DeploymentStatus struct {
	//creation date in RFC3339 format
	CreatedAt string `json:"created_at"`
	//update date in RFC3339 format
	UpdatedAt           string `json:"updated_at"`
	Replicas            int    `json:"replicas"`
	ReadyReplicas       int    `json:"ready_replicas"`
	AvailableReplicas   int    `json:"available_replicas"`
	UnavailableReplicas int    `json:"unavailable_replicas"`
	UpdatedReplicas     int    `json:"updated_replicas"`
}

// UpdateReplicas -- contains new number of replicas
//
// swagger:model
type UpdateReplicas struct {
	// required: true
	Replicas int `json:"replicas"`
}

// DeploymentsList -- model for deployments list
//
// swagger:model
type DeploymentsList struct {
	Deployments []Deployment `json:"deployments"`
}

// Deployment -- model for deployments
//
// swagger:model
type Deployment struct {
	Status *DeploymentStatus `json:"status,omitempty"`
	// required: true
	Containers []Container `json:"containers"`
	// required: true
	Name string `json:"name"`
	// required: true
	Replicas int `json:"replicas"`
	//total CPU usage by all containers in this deployment
	TotalCPU uint `json:"total_cpu,omitempty"`
	//total RAM usage by all containers in this deployment
	TotalMemory uint   `json:"total_memory,omitempty"`
	Owner       string `json:"owner,omitempty"`
}

// Container -- model for container in deployment
//
// swagger:model
type Container struct {
	// required: true
	Image string `json:"image"`
	// required: true
	Name string `json:"name"`
	// required: true
	Limits       Resource          `json:"limits"`
	Env          []Env             `json:"env,omitempty"`
	Commands     []string          `json:"commands,omitempty"`
	Ports        []ContainerPort   `json:"ports,omitempty"`
	VolumeMounts []ContainerVolume `json:"volume_mounts,omitempty"`
	ConfigMaps   []ContainerVolume `json:"config_maps,omitempty"`
}

// Env -- key-value pair of environment variables
//
// swagger:model
type Env struct {
	// required: true
	Value string `json:"value"`
	// required: true
	Name string `json:"name"`
}

// ContainerPort -- model for port in container
//
// swagger:model
type ContainerPort struct {
	// required: true
	Name string `json:"name"`
	// required: true
	Port int `json:"port"`
	// required: true
	Protocol Protocol `json:"protocol"`
}

// ContainerVolume -- volume (or config map) mounted in container
//
// swagger:model
type ContainerVolume struct {
	// required: true
	Name string  `json:"name"`
	Mode *string `json:"mode,omitempty"`
	// required: true
	MountPath                 string  `json:"mount_path"`
	SubPath                   *string `json:"sub_path,omitempty"`
	PersistentVolumeClaimName *string `json:"pvc_name,omitempty"`
}

// Mask removes information not interesting for users
func (deploy *Deployment) Mask() {
	deploy.Owner = ""
}
