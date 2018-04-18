package servactive

import (
	"encoding/json"
	"io/ioutil"

	"github.com/containerum/chkit/pkg/model/service"
)

func FromFile(filename string) (service.Service, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return service.Service{}, err
	}
	kubeServ := (&service.Service{}).ToKube()
	if err = json.Unmarshal(data, &kubeServ); err != nil {
		return service.Service{}, err
	}
	serv := service.ServiceFromKube(kubeServ)
	return serv, ValidateService(serv)
}
