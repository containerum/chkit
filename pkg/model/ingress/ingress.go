package ingress

import kubeModels "github.com/containerum/kube-client/pkg/model"

type Ingress struct {
	Name  string   `json:"name"`
	Rules RuleList `json:"rules"`
}

func IngressFromKube(kubeIngress kubeModels.Ingress) Ingress {
	return Ingress{
		Name:  kubeIngress.Name,
		Rules: RuleListFromKube(kubeIngress.Rules),
	}
}

func (ingr Ingress) ToKube() kubeModels.Ingress {
	return kubeModels.Ingress{
		Name:  ingr.Name,
		Rules: ingr.Rules.ToKube(),
	}
}

func (ingr Ingress) Copy() Ingress {
	return Ingress{
		Name:  ingr.Name,
		Rules: ingr.Rules.Copy(),
	}
}
