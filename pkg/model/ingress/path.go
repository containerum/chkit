package ingress

import (
	"fmt"

	kubeModels "github.com/containerum/kube-client/pkg/model"
)

type Path kubeModels.Path

func PathFromKube(kubePath kubeModels.Path) Path {
	return Path(kubePath)
}

func (path Path) ToKube() kubeModels.Path {
	return kubeModels.Path(path)
}

type PathList []Path

func PathListFromKube(kubeList []kubeModels.Path) PathList {
	var list PathList = make([]Path, 0, len(kubeList))
	for _, p := range kubeList {
		list = append(list, PathFromKube(p))
	}
	return list
}

func (list PathList) Copy() PathList {
	return append(make([]Path, 0, len(list)), list...)
}

func (list *PathList) Delete(i int) PathList {
	cp := list.Copy()
	return append(cp[:i], cp[i+1:]...)
}

func (list PathList) Append(paths ...Path) PathList {
	return append(list.Copy(), paths...)
}

func (list PathList) ToKube() []kubeModels.Path {
	kubeList := make([]kubeModels.Path, 0, len(list))
	for _, path := range list {
		kubeList = append(kubeList, path.ToKube())
	}
	return kubeList
}

type Service struct {
	Name string
	Port int
}

func (list PathList) Services() []Service {
	var services = make([]Service, 0, len(list))
	for _, path := range list {
		services = append(services, Service{
			Name: path.ServiceName,
			Port: path.ServicePort,
		})
	}
	return services
}

func (list PathList) ServicesNames() []string {
	var services = make([]string, 0, len(list))
	for _, path := range list {
		services = append(services, path.ServiceName)
	}
	return services
}

func (list PathList) ServicesPorts() []int {
	var ports = make([]int, 0, len(list))
	for _, path := range list {
		ports = append(ports, path.ServicePort)
	}
	return ports
}

func (list PathList) Paths() []string {
	var paths = make([]string, 0, len(list))
	for _, path := range list {
		paths = append(paths, path.Path)
	}
	return paths
}

func (list PathList) ServicesTableView() []string {
	var services = make([]string, 0, len(list))
	for _, path := range list {
		services = append(services,
			fmt.Sprintf("%q -> %s:%d",
				path.Path,
				path.ServiceName,
				path.ServicePort))
	}
	return services
}
