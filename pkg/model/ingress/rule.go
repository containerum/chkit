package ingress

import (
	kubeModels "git.containerum.net/ch/kube-client/pkg/model"
)

type Rule struct {
	Host      string
	TLSSecret *string
	Paths     PathList
}

func RuleFromKube(kubeRule kubeModels.Rule) Rule {
	return Rule{
		Host:      kubeRule.Host,
		TLSSecret: kubeRule.TLSSecret,
		Paths:     PathListFromKube(kubeRule.Path),
	}
}

func (rule Rule) ToKube() kubeModels.Rule {
	return kubeModels.Rule{
		Host:      rule.Host,
		TLSSecret: rule.TLSSecret,
		Path:      rule.Paths.ToKube(),
	}
}

func (rule Rule) Copy() Rule {
	return Rule{
		Host:      rule.Host,
		TLSSecret: rule.TLSSecret,
		Paths:     rule.Paths.Copy(),
	}
}

type RuleList []Rule

func RuleListFromKube(kubeList []kubeModels.Rule) RuleList {
	var list RuleList = make([]Rule, 0, len(kubeList))
	for _, rule := range kubeList {
		list = append(list, RuleFromKube(rule))
	}
	return list
}

func (list RuleList) ToKube() []kubeModels.Rule {
	kubeList := make([]kubeModels.Rule, 0, len(list))
	for _, rule := range list {
		kubeList = append(kubeList, rule.ToKube())
	}
	return kubeList
}

func (list RuleList) Copy() RuleList {
	cp := append(RuleList{}, list...)
	for i, rule := range cp {
		cp[i] = rule.Copy()
	}
	return cp
}

func (list RuleList) Delete(i int) RuleList {
	cp := list.Copy()
	return append(cp[:i], cp[i+1:]...)
}

func (list RuleList) Append(rules ...Rule) RuleList {
	return append(list.Copy(), rules...)
}

func (list RuleList) Hosts() []string {
	hosts := make([]string, 0, len(list))
	for _, rule := range list {
		hosts = append(hosts, rule.Host)
	}
	return hosts
}

func (list RuleList) Paths() PathList {
	var paths = make(PathList, 0, len(list))
	for _, rule := range list {
		paths = append(paths, rule.Paths.Copy()...)
	}
	return paths
}

func (list RuleList) Services() []Service {
	var services = make([]Service, 0, len(list))
	for _, rule := range list {
		services = append(services, rule.Paths.Services()...)
	}
	return services
}

func (list RuleList) ServicesNames() []string {
	var services = make([]string, 0, len(list))
	for _, rule := range list {
		services = append(services, rule.Paths.ServicesNames()...)
	}
	return services
}

func (list RuleList) ServicesTableView() []string {
	var services = make([]string, 0, len(list))
	for _, rule := range list {
		services = append(services, rule.Paths.ServicesTableView()...)
	}
	return services
}
