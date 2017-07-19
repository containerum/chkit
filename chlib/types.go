package chlib

import (
	"net"
	"time"
)

type Metadata struct {
	Annotations       map[string]string `json:"annotations,omitempty"`
	CreationTimestamp *time.Time        `json:"creationTimestamp,omitempty"`
	GenerateName      string            `json:"generateName,omitempty"`
	Labels            map[string]string `json:"labels,omitempty"`
	Name              string            `json:"name,omitempty"`
	Namespace         string            `json:"namespace,omitempty"`
	OwnerReferences   []struct {
		APIVersion         string `json:"apiVersion,omitempty"`
		BlockOwnerDeletion bool   `json:"blockOwnerDeletion,omitempty"`
		Controller         bool   `json:"controller,omitempty"`
		Kind               string `json:"kind,omitempty"`
		Name               string `json:"name,omitempty"`
		UID                string `json:"uid,omitempty"`
	} `json:"ownerReferences,omitempty"`
	ResourceVersion string `json:"resourceVersion,omitempty"`
	SelfLink        string `json:"selfLink,omitempty"`
	UID             string `json:"uid,omitempty"`
}

type HwSpecs struct {
	LimitsCPU      string `json:"limits.cpu,omitempty"`
	LimitsMemory   string `json:"limits.memory,omitempty"`
	RequestsCPU    string `json:"requests.cpu,omitempty"`
	RequestsMemory string `json:"requests.memory,omitempty"`
}

type Specs struct {
	Hard HwSpecs `json:"hard,omitempty"`
}

type HwResources struct {
	CPU    string `json:"cpu,omitempty"`
	Memory string `json:"memory,omitempty"`
}

type Resources struct {
	Limits   *HwResources `json:"limits,omitempty"`
	Requests *HwResources `json:"requests,omitempty"`
}

type EnvVar struct {
	Name  string `json:"name,omitempty"`
	Value string `json:"value,omitempty"`
}

type Port struct {
	ContainerPort int    `json:"containerPort,omitempty"`
	Port          int    `json:"port,omitempty"`
	TargetPort    int    `json:"targetPort,omitempty"`
	Name          string `json:"name,omitempty"`
	Protocol      string `json:"protocol,omitempty"`
}

type Container struct {
	Command                  []string  `json:"command,omitempty"`
	Env                      []EnvVar  `json:"env,omitempty"`
	Image                    string    `json:"image,omitempty"`
	ImagePullPolicy          string    `json:"imagePullPolicy,omitempty"`
	Name                     string    `json:"name,omitempty"`
	Ports                    []Port    `json:"ports,omitempty"`
	Resources                Resources `json:"resources,omitempty"`
	TerminationMessagePath   string    `json:"terminationMessagePath,omitempty"`
	TerminationMessagePolicy string    `json:"terminationMessagePolicy,omitempty"`
}

type Condiniton struct {
	LastTransitionTime time.Time `json:"lastTransitionTime,omitempty"`
	LastUpdateTime     time.Time `json:"lastUpdateTime,omitempty"`
	Message            string    `json:"message,omitempty"`
	Reason             string    `json:"reason,omitempty"`
	Status             string    `json:"status,omitempty"`
	Type               string    `json:"type,omitempty"`
}

type ContainerStatus struct {
	ContainerID  string `json:"containerID,omitempty"`
	Image        string `json:"image,omitempty"`
	ImageID      string `json:"imageID,omitempty"`
	Name         string `json:"name,omitempty"`
	Ready        bool   `json:"ready,omitempty"`
	RestartCount int    `json:"restartCount,omitempty"`
	State        struct {
		Running struct {
			StartedAt time.Time `json:"startedAt,omitempty"`
		} `json:"running,omitempty"`
	} `json:"state,omitempty"`
}

type NodeSelector struct {
	Role string `json:"role,omitempty"`
}

type Spec struct {
	AutomountServiceAccountToken  bool        `json:"automountServiceAccountToken,omitempty"`
	Containers                    []Container `json:"containers,omitempty"`
	DNSPolicy                     string      `json:"dnsPolicy,omitempty"`
	NodeName                      string      `json:"nodeName,omitempty"`
	RestartPolicy                 string      `json:"restartPolicy,omitempty"`
	SchedulerName                 string      `json:"schedulerName,omitempty"`
	ServiceAccount                string      `json:"serviceAccount,omitempty"`
	ServiceAccountName            string      `json:"serviceAccountName,omitempty"`
	TerminationGracePeriodSeconds int         `json:"terminationGracePeriodSeconds,omitempty"`
	Tolerations                   []struct {
		Effect            string `json:"effect,omitempty"`
		Key               string `json:"key,omitempty"`
		Operator          string `json:"operator,omitempty"`
		TolerationSeconds int    `json:"tolerationSeconds,omitempty"`
	} `json:"tolerations,omitempty"`
}

type Service struct {
	Kind     string   `json:"kind"`
	Metadata Metadata `json:"metadata,omitempty"`
	Spec     struct {
		ClusterIP           net.IP            `json:"clusterIP,omitempty"`
		DeprecatedPublicIPs []net.IP          `json:"deprecatedPublicIPs,omitempty"`
		ExternalHosts       []string          `json:"externalHosts,omitempty"`
		DomainHosts         []string          `json:"domainHosts,omitempty"`
		Ports               []Port            `json:"ports,omitempty"`
		Selector            map[string]string `json:"selector,omitempty"`
		SessionAffinity     string            `json:"sessionAffinity,omitempty"`
		Type                string            `json:"type,omitempty"`
	} `json:"spec,omitempty"`
	Status map[string]interface{} `json:"status,omitempty"`
}

type Deploy struct {
	Kind     string   `json:"kind"`
	Metadata Metadata `json:"metadata,omitempty"`
	Spec     struct {
		ProgressDeadlineSeconds int `json:"progressDeadlineSeconds,omitempty"`
		Replicas                int `json:"replicas,omitempty"`
		RevisionHistoryLimit    int `json:"revisionHistoryLimit,omitempty"`
		Selector                *struct {
			MatchLabels map[string]string `json:"matchLabels,omitempty"`
		} `json:"selector,omitempty"`
		Strategy map[string]interface{} `json:"strategy,omitempty"`
		Template struct {
			Metadata Metadata `json:"metadata,omitempty"`
			Spec     Spec     `json:"spec,omitempty"`
		} `json:"template,omitempty"`
	} `json:"spec,omitempty"`
	Status *struct {
		AvailableReplicas   int          `json:"availableReplicas,omitempty"`
		Conditions          []Condiniton `json:"conditions,omitempty"`
		ObservedGeneration  int          `json:"observedGeneration,omitempty"`
		ReadyReplicas       int          `json:"readyReplicas,omitempty"`
		Replicas            int          `json:"replicas,omitempty"`
		UpdatedReplicas     int          `json:"updatedReplicas,omitempty"`
		UnavaliableReplicas int          `json:"unavaliableReplicas,omitempty"`
	} `json:"status,omitempty"`
}

type Namespace struct {
	Kind     string   `json:"kind"`
	Metadata Metadata `json:"metadata,omitempty"`
	DataType string   `json:"DataType,omitempty"`
	Data     struct {
		APIVersion string   `json:"apiVersion,omitempty"`
		Kind       string   `json:"kind,omitempty"`
		Metadata   Metadata `json:"metadata,omitempty"`
		Spec       Specs    `json:"spec,omitempty"`
		Status     struct {
			Hard  HwSpecs `json:"hard,omitempty"`
			Used  HwSpecs `json:"used,omitempty"`
			Phase string  `json:"phase,omitempty"`
		} `json:"status,omitempty"`
	} `json:"data,omitempty"`
}

type Pod struct {
	Kind     string   `json:"kind"`
	Metadata Metadata `json:"metadata,omitempty"`
	Spec     Spec     `json:"spec,omitempty"`
	Status   struct {
		Conditions        []Condiniton      `json:"conditions,omitempty"`
		ContainerStatuses []ContainerStatus `json:"containerStatuses,omitempty"`
		HostIP            net.IP            `json:"hostIP,omitempty"`
		Phase             string            `json:"phase,omitempty"`
		PodIP             net.IP            `json:"podIP,omitempty"`
		QosClass          string            `json:"qosClass,omitempty"`
		StartTime         time.Time         `json:"startTime,omitempty"`
	} `json:"status,omitempty"`
}
