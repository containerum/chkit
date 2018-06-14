package model

import (
	"fmt"

	"github.com/blang/semver"
	"github.com/docker/distribution/reference"
)

// DeploymentStatus -- kubernetes status of deployment
//
// swagger:model
type DeploymentStatus struct {
	Replicas            int `json:"replicas"`
	ReadyReplicas       int `json:"ready_replicas"`
	AvailableReplicas   int `json:"available_replicas"`
	UnavailableReplicas int `json:"unavailable_replicas"`
	UpdatedReplicas     int `json:"updated_replicas"`
}

// DeploymentVersion -- model for deployment version update
//
// swagger:model
type DeploymentVersion struct {
	Version string `json:"version"`
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
	//creation date in RFC3339 format
	CreatedAt string `json:"created_at,omitempty"`
	//delete date in RFC3339 format
	DeletedAt string            `json:"deleted_at,omitempty"`
	Status    *DeploymentStatus `json:"status,omitempty"`
	// required: true
	Containers []Container `json:"containers"`
	// required: true
	Name string `json:"name"`
	// required: true
	Replicas int `json:"replicas"`
	//total CPU usage by all containers in this deployment
	TotalCPU uint `json:"total_cpu,omitempty"`
	//Solution ID (only if deployment is part of solution)
	SolutionID string `json:"solution_id,omitempty"`
	//total RAM usage by all containers in this deployment
	TotalMemory uint           `json:"total_memory,omitempty"`
	Owner       string         `json:"owner,omitempty"`
	Active      bool           `json:"active"`
	Version     semver.Version `json:"version"`
}

func (deployment Deployment) ImagesNames() []string {
	var images = make([]string, 0, len(deployment.Containers))
	for _, container := range deployment.Containers {
		images = append(images, container.Image)
	}
	return images
}

func (deployment Deployment) Images() []Image {
	var images = make([]Image, 0, len(deployment.Containers))
	for _, container := range deployment.Containers {
		var img, err = ImageFromString(container.Image)
		if err == nil {
			images = append(images, img)
		}
	}
	return images
}

func (deployment Deployment) ContainersNames() []string {
	var names = make([]string, 0, len(deployment.Containers))
	for _, container := range deployment.Containers {
		names = append(names, container.Name)
	}
	return names
}

func (deployment Deployment) ContainersAndImages() []string {
	var items = make([]string, 0, len(deployment.Containers))
	for _, container := range deployment.Containers {
		items = append(items, fmt.Sprintf("%s [%s]", container.Name, container.Image))
	}
	return items
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

func (container Container) Version() string {
	var ref, err = reference.Parse(container.Image)
	if err != nil {
		return ""
	}
	if tagged, ok := ref.(reference.Tagged); ok && tagged != nil {
		return tagged.Tag()
	}
	return ""
}

func (container *Container) AddEnv(env Env) {
	for i, cont := range container.Env {
		if cont.Name == env.Name {
			container.Env[i].Value = env.Value
			return
		}
	}
	container.Env = append(container.Env, env)
}

func (container *Container) GetEnv(name string) (Env, bool) {
	for _, env := range container.Env {
		if env.Name == name {
			return env, true
		}
	}
	return Env{}, false
}

func (container *Container) GetEnvMap() map[string]string {
	var envs = make(map[string]string, len(container.Env))
	for _, env := range container.Env {
		envs[env.Name] = env.Value
	}
	return envs
}

func (container *Container) PutEnvMap(envs map[string]string) {
	for k, v := range envs {
		container.AddEnv(Env{
			Name:  k,
			Value: v,
		})
	}
}

type Image struct {
	Name string
	Tag  string
}

func ImageFromString(str string) (Image, error) {
	var img, err = reference.ParseNamed(str)
	if err != nil {
		return Image{}, err
	}
	if tagged, ok := img.(reference.NamedTagged); tagged != nil && ok {
		return Image{
			Name: tagged.Name(),
			Tag:  tagged.Tag(),
		}, nil
	}
	return Image{
		Name: img.Name(),
	}, nil
}

func (image Image) String() string {
	return image.Name + ":" + image.Tag
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
