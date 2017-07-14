package requestresults

import "time"

type metadata struct {
	Annotations       map[string]string `json:"annotations"`
	CreationTimestamp *time.Time        `json:"creationTimestamp"`
	GenerateName      string            `json:"generateName"`
	Labels            map[string]string `json:"labels"`
	Name              string            `json:"name"`
	Namespace         string            `json:"namespace"`
	OwnerReferences   []struct {
		APIVersion         string `json:"apiVersion"`
		BlockOwnerDeletion bool   `json:"blockOwnerDeletion"`
		Controller         bool   `json:"controller"`
		Kind               string `json:"kind"`
		Name               string `json:"name"`
		UID                string `json:"uid"`
	} `json:"ownerReferences"`
	ResourceVersion string `json:"resourceVersion"`
	SelfLink        string `json:"selfLink"`
	UID             string `json:"uid"`
}

type hwSpecs struct {
	LimitsCPU      string `json:"limits.cpu"`
	LimitsMemory   string `json:"limits.memory"`
	RequestsCPU    string `json:"requests.cpu"`
	RequestsMemory string `json:"requests.memory"`
}

type specs struct {
	Hard hwSpecs `json:"hard"`
}

type resources struct {
	Limits struct {
		CPU    string `json:"cpu"`
		Memory string `json:"memory"`
	} `json:"limits"`
	Requests struct {
		CPU    string `json:"cpu"`
		Memory string `json:"memory"`
	} `json:"requests"`
}

type envVar struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type port struct {
	ContainerPort int    `json:"containerPort"`
	Port          int    `json:"port"`
	TargetPort    int    `json:"targetPort"`
	Name          string `json:"name"`
	Protocol      string `json:"protocol"`
}

type container struct {
	Command                  []string  `json:"command"`
	Env                      []envVar  `json:"env"`
	Image                    string    `json:"image"`
	ImagePullPolicy          string    `json:"imagePullPolicy"`
	Name                     string    `json:"name"`
	Ports                    []port    `json:"ports"`
	Resources                resources `json:"resources"`
	TerminationMessagePath   string    `json:"terminationMessagePath"`
	TerminationMessagePolicy string    `json:"terminationMessagePolicy"`
}

type condiniton struct {
	LastTransitionTime time.Time `json:"lastTransitionTime"`
	LastUpdateTime     time.Time `json:"lastUpdateTime"`
	Message            string    `json:"message"`
	Reason             string    `json:"reason"`
	Status             string    `json:"status"`
	Type               string    `json:"type"`
}

type containerStatus struct {
	ContainerID  string `json:"containerID"`
	Image        string `json:"image"`
	ImageID      string `json:"imageID"`
	Name         string `json:"name"`
	Ready        bool   `json:"ready"`
	RestartCount int    `json:"restartCount"`
	State        struct {
		Running struct {
			StartedAt time.Time `json:"startedAt"`
		} `json:"running"`
	} `json:"state"`
}

type nodeSelector struct {
	Role string `json:"role"`
}

type spec struct {
	AutomountServiceAccountToken  bool         `json:"automountServiceAccountToken"`
	Containers                    []container  `json:"containers"`
	DNSPolicy                     string       `json:"dnsPolicy"`
	NodeName                      string       `json:"nodeName"`
	NodeSelector                  nodeSelector `json:"nodeSelector"`
	RestartPolicy                 string       `json:"restartPolicy"`
	SchedulerName                 string       `json:"schedulerName"`
	ServiceAccount                string       `json:"serviceAccount"`
	ServiceAccountName            string       `json:"serviceAccountName"`
	TerminationGracePeriodSeconds int          `json:"terminationGracePeriodSeconds"`
	Tolerations                   []struct {
		Effect            string `json:"effect"`
		Key               string `json:"key"`
		Operator          string `json:"operator"`
		TolerationSeconds int    `json:"tolerationSeconds"`
	} `json:"tolerations"`
}
