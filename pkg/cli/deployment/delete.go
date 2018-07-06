package clideployment

import (
	"fmt"
	"runtime"
	"sync/atomic"

	"github.com/containerum/chkit/pkg/context"
	"github.com/containerum/chkit/pkg/model/deployment"
	"github.com/containerum/chkit/pkg/util/activekit"
	"github.com/containerum/chkit/pkg/util/ferr"
	"github.com/containerum/chkit/pkg/util/limiter"
	"github.com/ninedraft/boxofstuff/str"
	"github.com/octago/sflags/gen/gpflag"
	"github.com/spf13/cobra"
)

func Delete(ctx *context.Context) *cobra.Command {
	var flags = struct {
		Force       bool `dec:"suppress confirmation"`
		Concurrency uint `desc:"how much concurrent requeste can be performed at once" flag:"concurrency c"`
	}{
		Concurrency: uint(runtime.NumCPU()),
	}
	command := &cobra.Command{
		Use:     "deployment",
		Aliases: aliases,
		Short:   "delete deployment in specific namespace",
		// Long:    help.MustGetString("delete deployment"),
		Run: func(cmd *cobra.Command, args []string) {
			var logger = ctx.Log.Command("delete deployment")
			if flags.Concurrency == 0 {
				flags.Concurrency = 1
			}
			var deploymentsToDelete str.Vector = args
			if len(args) == 0 {
				if flags.Force {
					ferr.Printf("At least one deployment name must be provided in --force mode!\n")
					ctx.Exit(1)
				}
				logger.Debugf("getting deployment list from namespace %q", ctx.GetNamespace())
				var deploymentList, err = ctx.Client.GetDeploymentList(ctx.GetNamespace().ID)
				if err != nil {
					logger.WithError(err).Errorf("unable to get deployment list from namespace %q", ctx.GetNamespace())
					ferr.Println(err)
					ctx.Exit(1)
				}
				logger.Debugf("running deployment selection menu")
				for exit := false; !exit; {
					(&activekit.Menu{
						Title: "Select deployment:\nSelected: " + deploymentsToDelete.Join(" "),
						Items: activekit.StringSelector(deploymentList.Filter(func(depl deployment.Deployment) bool {
							return !deploymentsToDelete.Contains(depl.Name)
						}).Names(), func(s string) error {
							deploymentsToDelete = append(deploymentsToDelete, s)
							logger.Debugf("deployment %s selected", s)
							return nil
						}).Append(&activekit.MenuItem{
							Label: "Confirm",
							Action: func() error {
								exit = true
								return nil
							},
						}),
					}).Run()
				}
			}
			if flags.Force || activekit.YesNo("Are you really want to delete to delete %s?", deploymentsToDelete.Join(", ")) {
				logger.Debugf("deleting %d deployments: %v", deploymentsToDelete.Len(), deploymentsToDelete)
				var deleted uint64
				var limit = limiter.New(flags.Concurrency)
				for _, deplName := range deploymentsToDelete {
					go func(done func(), depl string) {
						defer done()
						logger.Debugf("deleting deployment %s in namespace %v", depl, ctx.GetNamespace())
						if err := ctx.Client.DeleteDeployment(ctx.GetNamespace().ID, depl); err != nil {
							logger.WithError(err).Errorf("unable to delete deployment %s in namespace %v", depl, ctx.GetNamespace())
							ferr.Println(err)
							return
						}
						logger.Debugf("deployment %s from namespace %s is deleted", depl, ctx.GetNamespace())
						fmt.Printf("Deployment %s is deleted\n", depl)
						atomic.AddUint64(&deleted, 1)
					}(limit.Start(), deplName)
				}
				logger.Debugf("wait until all tasks are completed")
				limit.Wait()
				fmt.Printf("%d deployments are deleted\n", deleted)
			}
		},
	}
	if err := gpflag.ParseTo(&flags, command.PersistentFlags()); err != nil {
		panic(err)
	}
	return command
}
