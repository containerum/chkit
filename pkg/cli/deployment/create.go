package clideployment

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/containerum/chkit/pkg/context"
	"github.com/containerum/chkit/pkg/model/configmap"
	"github.com/containerum/chkit/pkg/model/deployment"
	"github.com/containerum/chkit/pkg/model/deployment/deplactive"
	"github.com/containerum/chkit/pkg/porta"
	"github.com/containerum/chkit/pkg/util/activekit"
	"github.com/containerum/chkit/pkg/util/angel"
	"github.com/containerum/chkit/pkg/util/coblog"
	"github.com/containerum/chkit/pkg/util/ferr"
	"github.com/octago/sflags/gen/gpflag"
	"github.com/spf13/cobra"
)

func Create(ctx *context.Context) *cobra.Command {
	var flags struct {
		porta.Importer
		porta.Exporter
		deplactive.Flags
	}
	command := &cobra.Command{
		Use:     "deployment",
		Aliases: aliases,
		Short:   "create deployment",
		//	Long:    help.MustGetString("create deployment"),
		Run: func(cmd *cobra.Command, args []string) {
			var logger = coblog.Logger(cmd)
			logger.Struct(flags)
			logger.Debugf("running create deployment command")
			var depl deployment.Deployment
			if flags.ImportActivated() {
				if err := flags.Import(&depl); err != nil {
					ferr.Printf("unable to import deployment:\n%v\n", err)
					ctx.Exit(1)
				}
				flags.Flags = deplactive.FlagsFromDeployment(flags.Flags, depl)
			} else {
				var err error
				depl, err = flags.Deployment()
				if err != nil {
					ferr.Println(err)
					ctx.Exit(1)
				}
			}
			deplactive.Fill(&depl)
			switch {
			case flags.Force && flags.ExporterActivated():
				if err := flags.Export(depl); err != nil {
					ferr.Println(err)
					ctx.Exit(1)
				}
				return
			case flags.Force && !flags.ExporterActivated():
				if err := deplactive.ValidateDeployment(depl); err != nil {
					ferr.Println(err)
					ctx.Exit(1)
				}
				if err := ctx.Client.CreateDeployment(ctx.GetNamespace().ID, depl); err != nil {
					ferr.Println(err)
					ctx.Exit(1)
				}
				fmt.Printf("Deployment %s created\n", depl.Name)
				return
			}
			/*
				var volumes = make(chan volume.VolumeList)
				go func() {
					logger := logger.Component("getting namespace list")
					logger.Debugf("START")
					defer logger.Debugf("END")
					defer close(volumes)
					var volumeList, err = ctx.Client.GetVolumeList(ctx.GetNamespace().ID)
					if err != nil {
						logger.WithError(err).Errorf("unable to get volume list from namespace %q", ctx.GetNamespace())
						ferr.Println(err)
						ctx.Exit(1)
					}
					volumes <- volumeList
				}()
			*/
			var configs = make(chan configmap.ConfigMapList)
			go func() {
				logger := logger.Component("getting configmap list")
				logger.Debugf("START")
				defer logger.Debugf("END")
				defer close(configs)
				configList, err := ctx.Client.GetConfigmapList(ctx.GetNamespace().ID)
				if err != nil {
					logger.WithError(err).Errorf("unable to get configmap list")
					ferr.Println(err)
					ctx.Exit(1)
				}
				configs <- configList
			}()

			var configmapList = <-configs

			fmt.Println(depl.RenderTable())
			depl = deplactive.Wizard{
				EditName:   true,
				Deployment: &depl,
				Configmaps: configmapList.Names(),
				//			Volumes:    (<-volumes).Names(),
			}.Run()
			fmt.Println(depl.RenderTable())

			var pushed = false
			if activekit.YesNo("Are you sure you want to  create deployment %s?", depl.Name) {
				if err := ctx.Client.CreateDeployment(ctx.GetNamespace().ID, depl); err != nil {
					ferr.Println(err)
					ctx.Exit(1)
				}
				fmt.Printf("Deployment %s created\n", depl.Name)
				pushed = true
			}

			(&activekit.Menu{
				Items: activekit.MenuItems{
					{
						Label: fmt.Sprintf("Edit deployment %s", depl.Name),
						Action: func() error {
							depl = deplactive.Wizard{
								Deployment: &depl,
								Configmaps: configmapList.Names(),
							}.Run()
							if activekit.YesNo("Push changes to server?") {
								if pushed {
									if err := ctx.Client.ReplaceDeployment(ctx.GetNamespace().ID, depl); err != nil {
										ferr.Println(err)
									}
								} else {
									if err := ctx.Client.CreateDeployment(ctx.GetNamespace().ID, depl); err != nil {
										ferr.Println(err)
									}
									pushed = true
								}
							}
							return nil
						},
					},
					{
						Label: fmt.Sprintf("Save deployment %s to file", depl.Name),
						Action: func() error {
							for {
								var fname = activekit.Promt("Type filename (if ext is yaml or yml then file encodes as YAML, JSON by default): ")
								fname = strings.TrimSpace(fname)
								var err error
								var data string
								switch strings.ToLower(filepath.Ext(fname)) {
								case ".yaml", ".yml":
									fmt.Println("Encoding deployment as YAML")
									data, err = depl.RenderYAML()
								default:
									fmt.Println("Encoding deployment as JSON")
									data, err = depl.RenderJSON()
								}
								if err != nil {
									ferr.Println(err)
								}
								if err := ioutil.WriteFile(fname, []byte(data), os.ModePerm); err != nil {
									ferr.Println(err)
								}
							}
							return nil
						},
					},
				},
			}).Run()

		},
	}
	if err := gpflag.ParseTo(&flags, command.Flags()); err != nil {
		angel.Angel(ctx, fmt.Errorf("it seems that the structure of the flags is set incorrectly: %v", err))
	}
	return command
}
