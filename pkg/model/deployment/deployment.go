package deployment

type Deployment struct {
	Name     string
	Replicas int
	Status   *Status
	Volumes  []Volume
}
