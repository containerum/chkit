package clinamespace

import (
	"fmt"
	"strings"
	"time"

	"github.com/containerum/chkit/pkg/configuration"
	. "github.com/containerum/chkit/pkg/context"
	"github.com/containerum/chkit/pkg/model"
	"github.com/containerum/chkit/pkg/model/namespace"
	"github.com/containerum/chkit/pkg/util/animation"
	"github.com/containerum/chkit/pkg/util/trasher"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var aliases = []string{"ns", "namespaces"}
var getNamespaceDataConfig = struct {
	configuration.ExportConfig
}{}

var Get = &cobra.Command{
	Use:     "namespace",
	Aliases: aliases,
	Short:   `shows namespace data or namespace list`,
	Long:    `shows namespace data or namespace list. Aliases: ` + strings.Join(aliases, ", "),
	Example: "chkit get namespace_name... [-o yaml/json] [-f output_file]",
	Run: func(command *cobra.Command, args []string) {
		logrus.WithFields(logrus.Fields{
			"command": "get namespace",
		}).Debug("getting namespace data")

		anime := &animation.Animation{
			Framerate:      0.5,
			ClearLastFrame: true,
			Source:         trasher.NewSilly(),
		}
		go func() {
			time.Sleep(4 * time.Second)
			anime.Run()
		}()
		nsData, err := func() (model.Renderer, error) {
			defer anime.Stop()
			switch len(args) {
			case 1:
				namespaceLabel := args[0]
				logrus.Debugf("getting namespace %q", namespaceLabel)
				ns, err := Context.Client.GetNamespace(namespaceLabel)
				if err != nil {
					logrus.WithError(err).Errorf("unable to get namespace %q", namespaceLabel)
					return nil, err
				}
				return ns, nil
			default:
				var list namespace.NamespaceList
				logrus.Debugf("getting namespace list")
				list, err := Context.Client.GetNamespaceList()
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

func init() {
	Get.PersistentFlags().
		StringVarP((*string)(&getNamespaceDataConfig.Format), "output", "o", "", "output format (json/yaml)")
	Get.PersistentFlags().
		StringVarP(&getNamespaceDataConfig.Filename, "file", "f", "", "output file")
}
