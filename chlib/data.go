package chlib

import (
	"fmt"
	"os/user"
	"path"
	"regexp"
)

var homeDir string

const (
	KindDeployments = "deployments"
	KindNamespaces  = "namespaces"
	KindPods        = "pods"
	KindService     = "services"
)

const (
	KeyImage    = "image"
	KeyReplicas = "replicas"
)

const (
	DefaultReplicas      = 1
	DefaultCPURequest    = "100m"
	DefaultMemoryRequest = "128Mi"
)

func init() {
	currentUser, err := user.Current()
	if err != nil {
		panic(fmt.Errorf("get current user: %s", err))
	}
	homeDir = currentUser.HomeDir
	ConfigDir = path.Join(homeDir, ".containerum")
	ConfigFile = path.Join(ConfigDir, "config.db")
	SrcFolder = path.Join(ConfigDir, "src")
	TemplatesFolder = path.Join(SrcFolder, "json_templates")
	RunFile = path.Join(TemplatesFolder, "run.json")
	ExposeFile = path.Join(TemplatesFolder, "expose.json")
}

var (
	ConfigDir       string
	ConfigFile      string
	SrcFolder       string
	TemplatesFolder string
	RunFile         string
	ExposeFile      string
)

const DefaultProto = "TCP"

const (
	nameRegex       = `[a-z0-9]([-a-z0-9]*[a-z0-9])?`
)

var (
	LabelRegex      = regexp.MustCompile(`^` + nameRegex + `$`)
	ImageRegex      = regexp.MustCompile(`(?:.+/)?([^:]+)(?::.+)?`)
	CpuRegex        = regexp.MustCompile(`^\d+(.\d+)?m?$`)
	MemRegex        = regexp.MustCompile(`^\d+(.\d+)?(Mi|Gi)$`)
	ObjectNameRegex = LabelRegex
	PortRegex       = regexp.MustCompile(`^(\D+):(\d+)(:(\d+))?(:(TCP|UDP))?$`)
	PortNameRegex   = LabelRegex
	VolumesRegex    = regexp.MustCompile(`([^"]\S*|".+?")\s*`)
	VolumeRegex     = regexp.MustCompile(`^"?(` + nameRegex + `)(\/([^/][^\x00]+))?"?="?(\/[^\x00]+)"?$`) // format: "volumeLabel/subPath"="/mountPath"
)