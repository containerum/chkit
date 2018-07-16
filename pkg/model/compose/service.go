package compose

type Service struct {
	Name      string `yaml:"-"`
	Deploy    Deploy
	Build     Build
	Ports     []string
	Expose    []string
	Volumes   []string
	DependsOn Dependecies
	EnvFile   string
	Secrets   []string
	Extends   Parent
	Command   Command
}

func (service Service) Copy() Service {
	var cp = service
	{
		var portsCp = make([]string, 0, len(service.Ports))
		cp.Ports = append(portsCp, service.Ports...)
	}
	{
		var volumesCp = make([]string, len(service.Volumes))
		cp.Volumes = append(volumesCp, service.Volumes...)
	}
	cp.Build = service.Build.Copy()
	return cp
}

type Services map[string]Service

func (services Services) Slice() []Service {
	var slice = make([]Service, 0, len(services))
	for _, service := range services {
		slice = append(slice, service)
	}
	return slice
}

func (services Services) New() Services {
	return make(Services, len(services))
}

func (services Services) Copy() Services {
	var cp = services.New()
	for name, service := range services {
		cp[name] = service.Copy()
	}
	return cp
}
