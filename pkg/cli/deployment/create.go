package clideployment

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"time"

	. "github.com/containerum/chkit/pkg/context"
	"github.com/containerum/chkit/pkg/model/deployment/deplactive"
	"github.com/containerum/chkit/pkg/util/activekit"
	"github.com/containerum/chkit/pkg/util/angel"
	"github.com/containerum/chkit/pkg/util/animation"
	"github.com/containerum/chkit/pkg/util/text"
	"github.com/containerum/chkit/pkg/util/trasher"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var createDeplConfig = struct {
	File  string
	Force bool
}{}

var Create = &cobra.Command{
	Use:     "deployment",
	Aliases: aliases,
	Run: func(cmd *cobra.Command, args []string) {
		depl := deplactive.DefaultDeployment()
		if cmd.Flag("file").Changed {
			var err error
			depl, err = deplactive.FromFile(createDeplConfig.File)
			if err != nil {
				logrus.WithError(err).Errorf("unable to load deployment data from file %s", createDeplConfig.File)
				fmt.Printf("Unable to load deployment data from file :(\n%v", err)
				os.Exit(1)
			}
		}
		depl, err := deplactive.ConstructDeployment(deplactive.Config{
			Deployment: &depl,
		})
		if err != nil {
			logrus.WithError(err).Errorf("unable to create deployment")
			fmt.Println(err)
			os.Exit(1)
		}
		for {
			_, err := (&activekit.Menu{
				Items: []*activekit.MenuItem{
					{
						Name: "Create deployment",
						Action: func() error {
							anime := &animation.Animation{
								Source: trasher.NewSilly(),
							}
							go func() {
								time.Sleep(4 * time.Second)
								anime.Run()
							}()
							err := func() error {
								defer anime.Stop()
								return Context.Client.CreateDeployment(Context.Namespace, depl)
							}()
							if err != nil {
								logrus.WithError(err).Errorf("unable to create deployment %q", depl.Name)
								fmt.Println(err)
							}
							return nil
						},
					},
					{
						Name: "Edit deployment",
						Action: func() error {
							var err error
							depl, err = deplactive.ConstructDeployment(deplactive.Config{
								Deployment: &depl,
							})
							if err != nil {
								logrus.WithError(err).Errorf("unable to create deployment")
								fmt.Println(err)
								os.Exit(1)
							}
							return nil
						},
					},
					{
						Name: "Print to terminal",
						Action: activekit.ActionWithErr(func() error {
							if data, err := depl.RenderYAML(); err != nil {
								return err
							} else {
								upBorders := strings.Repeat("_", text.Width(data))
								downBorders := strings.Repeat("_", text.Width(data))
								fmt.Printf("%s\n\n%s\n%s\n", upBorders, data, downBorders)
							}
							return nil
						}),
					},
					{
						Name: "Save to file",
						Action: func() error {
							filename, _ := activekit.AskLine("Print filename: ")
							data, err := depl.RenderJSON()
							if err != nil {
								return err
							}
							if err := ioutil.WriteFile(filename, []byte(data), os.ModePerm); err != nil {
								logrus.WithError(err).Errorf("unable to save deployment %q to file", depl.Name)
								fmt.Printf("Unable to save deployment to file :(\n%v", err)
								return nil
							}
							fmt.Printf("OK\n")
							return nil
						},
					},
					{
						Name: "Exit",
						Action: func() error {
							if yes, _ := activekit.Yes("Are you sure you want to exit?"); yes {
								os.Exit(0)
							}
							return nil
						},
					},
				},
			}).Run()
			if err != nil {
				logrus.WithError(err).Errorf("error while menu execution")
				angel.Angel(err)
				os.Exit(1)
			}
		}
	},
}

func textWidth(text string) int {
	width := 0
	for _, line := range strings.Split(text, "\n") {
		if len(line) > width {
			width = len(line)
		}
	}
	return width
}

func init() {
	Create.PersistentFlags().
		StringVar(&createDeplConfig.File, "file", "", "create deployment from file")
	Create.PersistentFlags().
		BoolVarP(&createDeplConfig.Force, "force", "f", false, "create from file without customisation")
}
