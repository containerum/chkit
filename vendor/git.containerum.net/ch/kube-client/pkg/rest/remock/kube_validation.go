package remock

import (
	"fmt"

	"git.containerum.net/ch/kube-client/pkg/model"
	"git.containerum.net/ch/kube-client/pkg/rest/remock/ermockerr"
	api_resource "k8s.io/apimachinery/pkg/api/resource"
	api_validation "k8s.io/apimachinery/pkg/util/validation"
)

const (
	minDeployCPU    = "10m"
	minDeployMemory = "10Mi"
	maxDeployCPU    = "4"
	maxDeployMemory = "4Gi"

	maxDeployReplicas = 10

	minport = 11000
	maxport = 65535
)

const (
	fieldShouldExist   = "Field %v should be provided"
	invalidReplicas    = "Invalid replicas number: %v. It must be between 1 and %v"
	invalidPort        = "Invalid port: %v. It must be between %v and %v"
	invalidProtocol    = "Invalid protocol: %v. It must be TCP or UDP"
	noOwner            = "Owner should be provided"
	invalidOwner       = "Owner should be UUID"
	noContainer        = "Container %v is not found in deployment"
	invalidName        = "Invalid name: %v. It must consist of lower case alphanumeric characters, '-' or '.', and must start and end with an alphanumeric character"
	invalidKey         = "Invalid key: %v. It must consist of alphanumeric characters, '-', '_' or '.'"
	invalidIP          = "Invalid IP: %v. It must be a valid IP address, (e.g. 10.9.8.7)"
	invalidCPUQuota    = "Invalid CPU quota: %v. It must be between %v and %v"
	invalidMemoryQuota = "Invalid memory quota: %v. It must be between %v and %v"
)

func ValidateDeployment(deployment model.Deployment) error {
	errs := []error{}
	if len(api_validation.IsDNS1123Subdomain(deployment.Name)) > 0 {
		errs = append(errs, fmt.Errorf(invalidName, deployment.Name))
	}
	if len(api_validation.IsInRange(deployment.Replicas, 1, maxDeployReplicas)) > 0 {
		errs = append(errs, fmt.Errorf(invalidReplicas, deployment.Replicas, maxDeployReplicas))
	}
	if deployment.Containers == nil || len(deployment.Containers) == 0 {
		errs = append(errs, fmt.Errorf(fieldShouldExist, "Containers"))
	}
	if len(errs) > 0 {
		return ermockerr.ErrInvalidDeployment().
			AddDetailsErr(errs...)
	}
	return nil
}

func ValidateContainer(container model.Container) error {
	errs := []error{}
	cpu, err := api_resource.ParseQuantity(container.Limits.CPU)
	if err != nil {
		return ermockerr.ErrInvalidContainer().
			AddDetailsErr(err)
	}
	mem, err := api_resource.ParseQuantity(container.Limits.Memory)
	if err != nil {
		return ermockerr.ErrInvalidContainer().
			AddDetailsErr(err)
	}
	mincpu, _ := api_resource.ParseQuantity(minDeployCPU)
	maxcpu, _ := api_resource.ParseQuantity(maxDeployCPU)
	minmem, _ := api_resource.ParseQuantity(minDeployMemory)
	maxmem, _ := api_resource.ParseQuantity(maxDeployMemory)

	if len(api_validation.IsDNS1123Subdomain(container.Name)) > 0 {
		errs = append(errs, fmt.Errorf(invalidName, container.Name))
	}

	if cpu.Cmp(mincpu) == -1 || cpu.Cmp(maxcpu) == 1 {
		errs = append(errs, fmt.Errorf(invalidCPUQuota, cpu.String(), minDeployCPU, maxDeployCPU))
	}

	if mem.Cmp(minmem) == -1 || mem.Cmp(maxmem) == 1 {
		errs = append(errs, fmt.Errorf(invalidMemoryQuota, mem.String(), minDeployMemory, maxDeployMemory))
	}

	for _, v := range container.Ports {
		if len(api_validation.IsValidPortName(v.Name)) > 0 {
			errs = append(errs, fmt.Errorf(invalidName, v.Name))
		}
		if v.Protocol != model.UDP && v.Protocol != model.TCP {
			errs = append(errs, fmt.Errorf(invalidProtocol, v.Protocol))
		}
		if len(api_validation.IsValidPortNum(v.Port)) > 0 {
			errs = append(errs, fmt.Errorf(invalidPort, v.Port, minport, maxport))
		}
	}

	for _, v := range container.Env {
		if len(api_validation.IsEnvVarName(v.Value)) > 0 {
			errs = append(errs, fmt.Errorf(fieldShouldExist, "Env: Value"))
		}
		if v.Name == "" {
			errs = append(errs, fmt.Errorf(fieldShouldExist, "Env: Name"))
		}
	}

	for _, v := range container.VolumeMounts {
		if len(api_validation.IsDNS1123Subdomain(v.Name)) > 0 {
			errs = append(errs, fmt.Errorf(invalidName, v.Name))
		}
		if v.MountPath == "" {
			errs = append(errs, fmt.Errorf(fieldShouldExist, "Volume: Mount path"))
		}
	}

	for _, v := range container.ConfigMaps {
		if len(api_validation.IsDNS1123Subdomain(v.Name)) > 0 {
			errs = append(errs, fmt.Errorf(invalidName, v.Name))
		}
		if v.MountPath == "" {
			errs = append(errs, fmt.Errorf(fieldShouldExist, "Config: Map mount path"))
		}
	}

	if len(errs) > 0 {
		return ermockerr.ErrInvalidContainer().
			AddDetailsErr(errs...)
	}
	return nil
}
