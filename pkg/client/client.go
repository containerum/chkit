package client

import (
	kubeClient "git.containerum.net/ch/kube-client/pkg/cmd"
	kubeModels "git.containerum.net/ch/kube-client/pkg/model"
	"github.com/containerum/chkit/pkg/model"
)

type ChkitClient struct {
	kube   kubeClient.Client
	tokens kubeModels.Tokens
	config model.Config
}
