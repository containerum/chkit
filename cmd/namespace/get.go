package namespace

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/containerum/chkit/pkg/chkitErrors"

	"github.com/containerum/chkit/cmd/util"
	"github.com/containerum/chkit/pkg/model"
	"gopkg.in/urfave/cli.v2"
)

type renderer interface {
	model.TableRenderer
	model.YAMLrenderer
	model.JSONrenderer
}

// GetNamespace -- commmand 'get' entity data
var GetNamespace = &cli.Command{
	Name:        "ns",
	Description: `show namespace or namespace list`,
	Usage:       `Shows namespace data or namespace list`,
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name: "json",
		},
	},
	Action: func(ctx *cli.Context) error {
		log := util.GetLog(ctx)
		client := util.GetClient(ctx)
		log.Debugf("get ns from %q", client.APIaddr)
		var namespaceData renderer
		if ctx.NArg() > 0 {
			name := ctx.Args().First()
			ns, err := client.GetNamespace(name)
			if err != nil {
				return err
			}
			namespaceData = &ns
		} else {
			list, err := client.GetNamespaceList()
			if err != nil {
				return err
			}
			namespaceData = list
		}
		ok, flag := isSet(ctx.Args().Slice(), "json", "yaml")
		switch {
		case ok && flag.Name == "json":
			log.Debugf("rendering namespace to JSON")
			data, err := namespaceData.RenderJSON()
			if err != nil {
				return err
			}
			if flag.Value == "" {
				fmt.Println(data)
			} else {
				if err := ioutil.WriteFile(flag.Value, []byte(data), os.ModePerm); err != nil {
					return chkitErrors.NewExitCoder(err)
				}
			}
		default:
			log.Debugf("rendering namespace to table")
			fmt.Println(namespaceData.RenderTable())
		}
		return nil
	},
}

func isSet(args []string, flags ...string) (bool, *cli.StringFlag) {
	for _, flag := range flags {
		for i, arg := range args {
			if len(flag) == 1 {
				arg = strings.TrimPrefix(arg, "-")
			} else if len(flag) > 1 {
				arg = strings.TrimPrefix(arg, "--")
			}
			if arg == flag {
				value := ""
				if i+1 < len(args) {
					value = args[i+1]
				}
				return true, &cli.StringFlag{
					Name:  flag,
					Value: value,
				}
			}
		}
	}
	return false, nil
}
