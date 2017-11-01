package cmd

import (
	"os"
	"path"
	"strings"

	"github.com/containerum/chkit/chlib"
	"github.com/containerum/chkit/helpers"
	"github.com/spf13/cobra"
)

var solutionCmdName string
var solutionCmdValidArgs = []string{"run"}
var solutionCmd = &cobra.Command{
	Use:       "solution run <solution name>",
	Short:     "Run a solution",
	Aliases:   []string{"sln"},
	ValidArgs: solutionCmdValidArgs,
	PreRun: func(cmd *cobra.Command, args []string) {
		if len(args) < 2 {
			np.FEEDBACK.Println("Solution name must be specified")
			cmd.Usage()
			os.Exit(1)
		}
		switch args[0] {
		case "run":
		default:
			np.FEEDBACK.Println("Valid subcommands are:", strings.Join(solutionCmdValidArgs, ", "))
			cmd.Usage()
			os.Exit(1)
		}
		solutionCmdName = args[1]
	},
	Run: func(cmd *cobra.Command, args []string) {
		nameParts := strings.Split(solutionCmdName, "/")
		solutionDirName := strings.TrimSuffix(nameParts[len(nameParts)-1], ".git")
		solutionPath := path.Join(chlib.SolutionsDir, solutionDirName)
		np.FEEDBACK.Println("Download...")
		branch, _ := cmd.Flags().GetString("branch")
		err := helpers.DownloadSolution(solutionCmdName, solutionPath, branch)
		exitOnErr(err)
		np.FEEDBACK.Println("Parse and run")
		envArg, _ := cmd.Flags().GetStringSlice("env")
		ns, _ := cmd.Flags().GetString("namespace")
		prefix, _ := cmd.Flags().GetString("prefix")
		err = client.RunSolution(solutionPath, envArg, ns, prefix)
		exitOnErr(err)
		np.FEEDBACK.Println("OK")
	},
}

func init() {
	solutionCmd.PersistentFlags().StringSliceP("env", "e", []string{}, "Environment variables. Format: key1=value1 ... keyN=valueN")
	solutionCmd.PersistentFlags().StringP("namespace", "n", "", "Namespace")
	solutionCmd.PersistentFlags().StringP("prefix", "p", "", "Services and deployments name prefix")
	solutionCmd.PersistentFlags().StringP("branch", "b", "master", "Branch in remote repo")
	RootCmd.AddCommand(solutionCmd)
}
