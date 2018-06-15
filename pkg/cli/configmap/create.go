package cliconfigmap

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"

	"github.com/containerum/chkit/pkg/context"
	"github.com/containerum/chkit/pkg/model/configmap"
	"github.com/containerum/chkit/pkg/model/configmap/activeconfigmap"
	"github.com/containerum/chkit/pkg/util/activekit"
	"github.com/containerum/chkit/pkg/util/coblog"
	"github.com/containerum/chkit/pkg/util/namegen"
	"github.com/containerum/kube-client/pkg/model"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	flag "github.com/spf13/pflag"
	"gopkg.in/yaml.v2"
)

var aliases = []string{"cm", "confmap", "conf-map", "comap"}

func Create(ctx *context.Context) *cobra.Command {
	comand := &cobra.Command{
		Use:     "configmap",
		Aliases: aliases,
		Run: func(cmd *cobra.Command, args []string) {
			var logger = coblog.Logger(cmd)
			var flags = cmd.Flags()
			var config, err = buildConfigMapFromFlags(flags, logger)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			force, _ := flags.GetBool("force")
			if !force {
				config = activeconfigmap.Config{
					EditName:  true,
					ConfigMap: &config,
				}.Wizard()
				fmt.Println(config.RenderTable())
			}
			if force || activekit.YesNo("Are you sure you want to create configmap %s?", config.Name) {
				if err := config.Validate(); err != nil {
					fmt.Println(err)
					os.Exit(1)
				}
				if err := ctx.Client.CreateConfigMap(ctx.Namespace.ID, config); err != nil {
					logger.WithError(err).Errorf("unable to create configmap %q", config.Name)
					fmt.Println(err)
					os.Exit(1)
				}
				fmt.Println("OK")
			} else if !force {
				config = activeconfigmap.Config{
					EditName:  false,
					ConfigMap: &config,
				}.Wizard()
				fmt.Println(config.RenderTable())
			}
		},
	}
	var persistentFlags = comand.PersistentFlags()
	persistentFlags.String("name", namegen.Aster()+"-"+namegen.Physicist(), "configmap name")
	persistentFlags.StringSlice("item-string", nil, "configmap item, KEY:VALUE string pair")
	persistentFlags.StringSlice("item-file", nil, "configmap file, KEY:FILE_PATH or FILE_PATH")
	persistentFlags.String("file", "", "file with configmap data")
	persistentFlags.BoolP("force", "f", false, "suppress confirmation")
	return comand
}

func buildConfigMapFromFlags(flags *flag.FlagSet, logger logrus.FieldLogger) (configmap.ConfigMap, error) {
	var config = configmap.ConfigMap{
		Data: make(model.ConfigMapData, 16),
	}
	if flags.Changed("file") {
		var err error
		var fName, _ = flags.GetString("file")
		data, err := ioutil.ReadFile(fName)
		if err != nil {
			logger.WithError(err).Error("unable to load configmap data from file")
			fmt.Println(err)
			os.Exit(1)
		}
		switch path.Ext(fName) {
		case "json":
			err = json.Unmarshal(data, &config)
		case "yaml":
			err = yaml.Unmarshal(data, &config)
		}
		return config, err
	} else {
		config.Name, _ = flags.GetString("name")
		if flags.Changed("item-string") {
			rawItems, _ := flags.GetStringSlice("item-string")
			items, err := getStringItems(rawItems)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			config = config.AddItems(items...)
		}
		if flags.Changed("item-file") {
			rawItems, _ := flags.GetStringSlice("item-file")
			items, err := getFileItems(rawItems)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			config = config.AddItems(items...)
		}
		return config, nil
	}
}

func getFileItems(rawItems []string) ([]configmap.Item, error) {
	var items = make([]configmap.Item, 0, len(rawItems))
	for _, rawItem := range rawItems {
		var filepath string
		var key string
		if tokens := strings.SplitN(rawItem, ":", 2); len(tokens) == 2 {
			key = strings.TrimSpace(tokens[0])
			filepath = tokens[1]
		} else if len(tokens) == 1 {
			key = path.Base(tokens[0])
			filepath = tokens[0]
		} else {
			logrus.Panicf("[chkit/pkg/cli/configmap.getFileItems] ivalid token number in raw file item", len(tokens))
		}
		value, err := ioutil.ReadFile(filepath)
		if err != nil {
			return nil, err
		}
		items = append(items, configmap.Item{
			Key:   key,
			Value: base64.StdEncoding.EncodeToString(value),
		})
	}
	return items, nil
}

func getStringItems(rawItems []string) ([]configmap.Item, error) {
	var items = make([]configmap.Item, 0, len(rawItems))
	for _, rawItem := range rawItems {
		var key string
		var value string
		if tokens := strings.SplitN(rawItem, ":", 2); len(tokens) == 2 {
			key = strings.TrimSpace(tokens[0])
			value = strings.TrimSpace(tokens[1])
		} else {
			logrus.Panicf("[chkit/pkg/cli/configmap.getStringItems] ivalid token number in raw string item", len(tokens))
		}
		items = append(items, configmap.Item{
			Key:   key,
			Value: value,
		})
	}
	return items, nil
}
