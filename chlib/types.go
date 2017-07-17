package chlib

import (
	"net"
	"time"
)

type Metadata struct {
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

type HwSpecs struct {
	LimitsCPU      string `json:"limits.cpu"`
	LimitsMemory   string `json:"limits.memory"`
	RequestsCPU    string `json:"requests.cpu"`
	RequestsMemory string `json:"requests.memory"`
}

type Specs struct {
	Hard HwSpecs `json:"hard"`
}

type Resources struct {
	Limits struct {
		CPU    string `json:"cpu"`
		Memory string `json:"memory"`
	} `json:"limits"`
	Requests struct {
		CPU    string `json:"cpu"`
		Memory string `json:"memory"`
	} `json:"requests"`
}

type EnvVar struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type Port struct {
	ContainerPort int    `json:"containerPort"`
	Port          int    `json:"port"`
	TargetPort    int    `json:"targetPort"`
	Name          string `json:"name"`
	Protocol      string `json:"protocol"`
}

type Container struct {
	Command                  []string  `json:"command"`
	Env                      []EnvVar  `json:"env"`
	Image                    string    `json:"image"`
	ImagePullPolicy          string    `json:"imagePullPolicy"`
	Name                     string    `json:"name"`
	Ports                    []Port    `json:"ports"`
	Resources                Resources `json:"resources"`
	TerminationMessagePath   string    `json:"terminationMessagePath"`
	TerminationMessagePolicy string    `json:"terminationMessagePolicy"`
}

type Condiniton struct {
	LastTransitionTime time.Time `json:"lastTransitionTime"`
	LastUpdateTime     time.Time `json:"lastUpdateTime"`
	Message            string    `json:"message"`
	Reason             string    `json:"reason"`
	Status             string    `json:"status"`
	Type               string    `json:"type"`
}

type ContainerStatus struct {
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

type NodeSelector struct {
	Role string `json:"role"`
}

type Spec struct {
	AutomountServiceAccountToken  bool         `json:"automountServiceAccountToken"`
	Containers                    []Container  `json:"containers"`
	DNSPolicy                     string       `json:"dnsPolicy"`
	NodeName                      string       `json:"nodeName"`
	NodeSelector                  NodeSelector `json:"nodeSelector"`
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

type Service struct {
	Metadata Metadata `json:"metadata"`
	Spec     struct {
		ClusterIP           net.IP            `json:"clusterIP"`
		DeprecatedPublicIPs []net.IP          `json:"deprecatedPublicIPs"`
		ExternalHosts       []string          `json:"externalHosts"`
		DomainHosts         []string          `json:"externalHosts"`
		Ports               []Port            `json:"ports"`
		Selector            map[string]string `json:"selector"`
		SessionAffinity     string            `json:"sessionAffinity"`
		Type                string            `json:"type"`
	} `json:"spec"`
	Status struct {
		LoadBalancer map[string]interface{} `json:"loadBalancer"`
	} `json:"status"`
}

type Deploy struct {
	Metadata Metadata `json:"metadata"`
	Spec     struct {
		ProgressDeadlineSeconds int `json:"progressDeadlineSeconds"`
		Replicas                int `json:"replicas"`
		RevisionHistoryLimit    int `json:"revisionHistoryLimit"`
		Selector                struct {
			MatchLabels map[string]string `json:"matchLabels"`
		} `json:"selector"`
		Strategy map[string]interface{} `json:"strategy"`
		Template struct {
			Metadata Metadata `json:"metadata"`
			Spec     Spec     `json:"spec"`
		} `json:"template"`
	} `json:"spec"`
	Status struct {
		AvailableReplicas   int          `json:"availableReplicas"`
		Conditions          []Condiniton `json:"conditions"`
		ObservedGeneration  int          `json:"observedGeneration"`
		ReadyReplicas       int          `json:"readyReplicas"`
		Replicas            int          `json:"replicas"`
		UpdatedReplicas     int          `json:"updatedReplicas"`
		UnavaliableReplicas int          `json:"unavaliableReplicas"`
	} `json:"status"`
}

type Namespace struct {
	Metadata Metadata `json:"metadata"`
	DataType string   `json:"DataType"`
	Data     struct {
		APIVersion string   `json:"apiVersion"`
		Kind       string   `json:"kind"`
		Metadata   Metadata `json:"metadata"`
		Spec       Specs    `json:"spec"`
		Status     struct {
			Hard  HwSpecs `json:"hard"`
			Used  HwSpecs `json:"used"`
			Phase string  `json:"phase"`
		} `json:"status"`
	} `json:"data"`
}

type Pod struct {
	Metadata Metadata `json:"metadata"`
	Spec     Spec     `json:"spec"`
	Status   struct {
		Conditions        []Condiniton      `json:"conditions"`
		ContainerStatuses []ContainerStatus `json:"containerStatuses"`
		HostIP            net.IP            `json:"hostIP"`
		Phase             string            `json:"phase"`
		PodIP             net.IP            `json:"podIP"`
		QosClass          string            `json:"qosClass"`
		StartTime         time.Time         `json:"startTime"`
	} `json:"status"`
}
