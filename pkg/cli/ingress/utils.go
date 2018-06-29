package clingress

import (
	"github.com/containerum/chkit/pkg/model/ingress"
	"github.com/spf13/cobra"
)

func buildIngress(cmd *cobra.Command, ingr ingress.Ingress) (ingress.Ingress, bool) {
	if !cmd.Flags().Parsed() {
		panic("[cli/clingress.buildIngress] attempt to build ingress.Ingress from non parsed flags")
	}
	ingr = ingr.Copy()
	var changed = false
	var rule, ruleChanged = buildRule(cmd, ingr.Rules.Head())
	if ruleChanged {
		ingr.Rules = ingress.RuleList{rule}
		changed = true
	}
	return ingr, changed
}

func buildRule(cmd *cobra.Command, rule ingress.Rule) (ingress.Rule, bool) {
	if !cmd.Flags().Parsed() {
		panic("[cli/clingress.buildRule] attempt to build ingress.Rule from non parsed flags")
	}
	rule = rule.Copy()
	var changed = false
	var path ingress.Path
	var flags = cmd.Flags()
	if len(rule.Paths) > 0 {
		path = rule.Paths[0]
	}
	path, pathChanged := buildPath(cmd, path)
	if pathChanged {
		rule.Paths = ingress.PathList{path}
		changed = true
	}
	if flags.Changed("host") {
		rule.Host, _ = flags.GetString("host")
	}
	if flags.Changed("tls-secret") {
		TLS, _ := flags.GetString("tls-secret")
		rule.TLSSecret = TLS
	}
	return rule, changed
}

func buildPath(cmd *cobra.Command, path ingress.Path) (ingress.Path, bool) {
	if !cmd.Flags().Parsed() {
		panic("[cli/clingress.buildPath] attempt to build ingress.Path from non parsed flags")
	}
	var changed = false
	flags := cmd.Flags()
	if flags.Changed("host") {
		path.Path, _ = flags.GetString("host")
		changed = true
	}
	if flags.Changed("service") {
		path.ServiceName, _ = flags.GetString("service")
		changed = true
	}
	if flags.Changed("port") {
		path.ServicePort, _ = flags.GetInt("port")
		changed = true
	}
	return path, changed
}
