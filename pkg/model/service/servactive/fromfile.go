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
	var serv service.Service
	err = json.Unmarshal(data, &serv)
	return serv, err
}
