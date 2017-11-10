package cmd

import (
	"os"
	"path"
	"strings"

	"github.com/containerum/chkit/chlib"
	"github.com/containerum/chkit/helpers"
	"github.com/spf13/cobra"
)

var solutionCmd = &cobra.Command{
	Use:     "solution",
	Short:   "Solutions management",
	Aliases: []string{"sln"},
}

var solutionCmdName string
var solutionRunCmd = &cobra.Command{
	Use:   "run <solution name>",
	Short: "Run a solution",
	PreRun: func(cmd *cobra.Command, args []string) {
		if len(args) < 2 {
			np.FEEDBACK.Println("Solution name must be specified")
			cmd.Usage()
			os.Exit(1)
		}
		solutionCmdName = args[0]
	},
	Run: func(cmd *cobra.Command, args []string) {
		np.SetPrefix("Solution")
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
		err = client.RunSolution(solutionPath, envArg, ns)
		exitOnErr(err)
		np.FEEDBACK.Println("OK")
	},
}

var solutionListCmd = &cobra.Command{
	Use:   "list",
	Short: "Show solutions made by Containerum",
	Run: func(cmd *cobra.Command, args []string) {
		np.SetPrefix("Solution")
		exitOnErr(helpers.ShowSolutionList())
	},
}

func init() {
	solutionRunCmd.PersistentFlags().StringSliceP("env", "e", []string{}, "Environment variables. Format: key1=value1 ... keyN=valueN")
	solutionRunCmd.PersistentFlags().StringP("namespace", "n", "", "Namespace")
	solutionRunCmd.PersistentFlags().StringP("branch", "b", "master", "Branch in remote repo")
	solutionCmd.AddCommand(solutionRunCmd)
	solutionCmd.AddCommand(solutionListCmd)
	RootCmd.AddCommand(solutionCmd)
}
