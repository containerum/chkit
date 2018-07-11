package clinamespace

import (
	"fmt"
	"strings"

	"github.com/containerum/chkit/pkg/context"
	"github.com/containerum/chkit/pkg/model/namespace"
	"github.com/containerum/chkit/pkg/porta"
	"github.com/containerum/chkit/pkg/util/ferr"
	"github.com/ninedraft/boxofstuff/str"
	"github.com/octago/sflags/gen/gpflag"
	"github.com/spf13/cobra"
)

var aliases = []string{"ns", "namespaces"}

func Get(ctx *context.Context) *cobra.Command {
	var flags struct {
		Names bool   `desc:"print namespace names"`
		Find  string `desc:"find namespace which name contains substring"`
		porta.Exporter
	}

	command := &cobra.Command{
		Use:     "namespace",
		Aliases: aliases,
		Short:   `show namespace data or namespace list`,
		Long:    "show namespace data or namespace list.",
		Example: "chkit get $ID... [-o yaml/json] [-f output_file]",
		Run: func(command *cobra.Command, args []string) {
			var logger = ctx.Log.Command("get namespace")
			logger.Debugf("START")
			defer logger.Debugf("END")
			var namespaces, err = ctx.Client.GetNamespaceList()
			if err != nil {
				ferr.Println(err)
				ctx.Exit(1)
			}
			if len(args) > 0 {
				namespaces = namespaces.Filter(func(namespace namespace.Namespace) bool {
					return str.Vector(args).Contains(namespace.Label) ||
						str.Vector(args).Contains(namespace.OwnerAndLabel())
				})
			} else if flags.Find != "" {
				namespaces = namespaces.Filter(func(namespace namespace.Namespace) bool {
					return strings.Contains(namespace.OwnerAndLabel(), flags.Find)
				})
			}
			if flags.Names {
				var names = str.Vector(namespaces.OwnersAndLabels())
				fmt.Println(names.Join("\n"))
				return
			}
			if err := flags.Export(namespaces); err != nil {
				ferr.Printf("unable to export data:\n%v\n", err)
				ctx.Exit(1)
			}
		},
	}
	if err := gpflag.ParseTo(&flags, command.PersistentFlags()); err != nil {
		panic(err)
	}
	return command
}
