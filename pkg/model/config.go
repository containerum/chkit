package model

import (
	kubeClientModels "git.containerum.net/ch/kube-client/pkg/model"
)

type Config struct {
	ConfigPath string
	Tokens     kubeClientModels.Tokens
	Client     ClientConfig
}

type ClientConfig struct {
	APIaddr  string
	Username string
	Password string
}
