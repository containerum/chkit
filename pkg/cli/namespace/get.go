package clinamespace

import (
	"fmt"
	"os"
	"strings"

	"github.com/containerum/chkit/pkg/cli/prerun"
	"github.com/containerum/chkit/pkg/configuration"
	"github.com/containerum/chkit/pkg/context"
	"github.com/containerum/chkit/pkg/model"
	"github.com/containerum/chkit/pkg/model/namespace"
	"github.com/containerum/chkit/pkg/util/angel"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var aliases = []string{"ns", "namespaces"}

func Get(ctx *context.Context) *cobra.Command {
	var getNamespaceDataConfig = struct {
		configuration.ExportConfig
	}{}
	command := &cobra.Command{
		Use:     "namespace",
		Aliases: aliases,
		Short:   `shows namespace data or namespace list`,
		Long:    `shows namespace data or namespace list. Aliases: ` + strings.Join(aliases, ", "),
		Example: "chkit get $ID... [-o yaml/json] [-f output_file]",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			if err := prerun.PreRun(ctx); err != nil {
				angel.Angel(ctx, err)
				os.Exit(1)
			}
			if cmd.Flags().Changed("namespace") {
				ctx.Namespace, _ = cmd.Flags().GetString("namespace")
			}
		},
		Run: func(command *cobra.Command, args []string) {
			logrus.WithFields(logrus.Fields{
				"command": "get namespace",
			}).Debug("getting namespace data")
			nsData, err := func() (model.Renderer, error) {
				switch len(args) {
				case 1:
					namespaceLabel := args[0]
					logrus.Debugf("getting namespace %q", namespaceLabel)
					ns, err := ctx.Client.GetNamespace(namespaceLabel)
					if err != nil {
						logrus.WithError(err).Errorf("unable to get namespace %q", namespaceLabel)
						return nil, err
					}
					return ns, nil
				default:
					var list namespace.NamespaceList
					logrus.Debugf("getting namespace list")
					list, err := ctx.Client.GetNamespaceList()
					if err != nil {
						logrus.WithError(err).Errorf("unable to get namespace list")
						return nil, err
					}
					return list, nil
				}
			}()
			if err != nil {
				fmt.Printf("%v\n", err)
				return
			}
			err = configuration.ExportData(nsData, getNamespaceDataConfig.ExportConfig)
			if err != nil {
				logrus.Debugf("fatal error: %v", err)
				return
			}
			logrus.Debugf("OK")
		},
	}

	command.PersistentFlags().
		StringVarP((*string)(&getNamespaceDataConfig.Format), "output", "o", "", "output format (json/yaml)")
	command.PersistentFlags().
		StringVarP(&getNamespaceDataConfig.Filename, "file", "f", "", "output file")

	return command
}
