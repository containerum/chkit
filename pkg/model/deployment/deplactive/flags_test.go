package deplactive

import (
	"bytes"
	"fmt"
	"testing"
)

func TestFlags_BuildContainers(t *testing.T) {
	var containers, err = (&Flags{
		Image: []string{
			"gateway@nginx",
			"blog@wordpress",
		},
		Env: []string{
			"gateway@PASS:12345678",
			"gateway@USER:merlin",
			"blog@MY_SQL_SECRET:asdopasfoahufheaudjia892eqd",
		},
		Memory: []string{
			"gateway@500",
			"blog@2000",
		},
		CPU: []string{
			"gateway@1000",
			"blog@3000",
		},
		Configmap: []string{
			"gateway@nginx-config@/etc/nginx/",
			"blog@wordpress-thema",
		},
	}).BuildContainers()
	if err != nil {
		t.Fatal(err)
	}
	for _, container := range containers {
		var buf = &bytes.Buffer{}
		fmt.Fprintf(buf, "Name: %q\n"+
			"Image: %q\n"+
			"Limits: %#v\n"+
			"Ports: %v\n"+
			"Env: %v\n"+
			"Configmap: %v\n"+
			"Volumes: %v\n",
			container.Name,
			container.Image,
			container.Limits,
			container.Ports,
			container.Env,
			container.ConfigMaps,
			container.VolumeMounts)
		t.Log(buf)
	}
}
