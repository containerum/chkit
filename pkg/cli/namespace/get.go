package clinamespace

import (
	"fmt"
	"strings"

	"github.com/containerum/chkit/pkg/configuration"
	"github.com/containerum/chkit/pkg/context"
	"github.com/containerum/chkit/pkg/model"
	"github.com/containerum/chkit/pkg/model/namespace"
	"github.com/go-siris/siris/core/errors"
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
		Run: func(command *cobra.Command, args []string) {
			var logger = ctx.Log.Command("get namespace")
			logger.Debugf("START")
			defer logger.Debugf("END")
			nsData, err := func() (model.Renderer, error) {
				switch len(args) {
				case 1:
					namespaceLabel := args[0]
					logrus.Debugf("getting namespace %q", namespaceLabel)
					nsList, err := ctx.Client.GetNamespaceList()
					if err != nil {
						logrus.WithError(err).Errorf("unable to get namespace list")
						return nil, err
					}
					ns, ok := nsList.GetByUserFriendlyID(namespaceLabel)
					if !ok {
						logrus.WithError(errors.New("not found")).
							Errorf("unable to get namespace %q", namespaceLabel)
						return nil, fmt.Errorf("namespace %q not found", namespaceLabel)
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
			logger.Debugf("exporting data")
			err = configuration.ExportData(nsData, getNamespaceDataConfig.ExportConfig)
			if err != nil {
				logger.WithError(err).Errorf("fatal error: %v", err)
				return
			}
		},
	}

	command.PersistentFlags().
		StringVarP((*string)(&getNamespaceDataConfig.Format), "output", "o", "", "output format (json/yaml)")
	command.PersistentFlags().
		StringVarP(&getNamespaceDataConfig.Filename, "file", "f", "", "output file")

	return command
}
