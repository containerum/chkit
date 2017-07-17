package chlib

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
	configDir       = ".containerum"
	configFile      = "config.db"
	srcFolder       = "src"
	templatesFolder = "json_templates"
	runFile         = "run.json"
	exposeFile      = "expose.json"
)

const DefaultProto = "TCP"
