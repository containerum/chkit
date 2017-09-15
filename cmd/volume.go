package cmd

import (
	"encoding/json"
	"os"

	"github.com/containerum/chkit.v2/chlib"
	"github.com/containerum/chkit.v2/chlib/requestresults"
	"github.com/containerum/chkit.v2/helpers"
	"github.com/spf13/cobra"
)

var volumeName string

var volumeCmd = &cobra.Command{
	Use:   "volume [NAME]",
	Short: "Show user volumes",
	PreRun: func(cmd *cobra.Command, args []string) {
		if len(args) > 1 {
			np.FEEDBACK.Println("Single parameter (volume name) must be specified")
			cmd.Usage()
			os.Exit(1)
		}
		if len(args) == 1 {
			volumeName = args[0]
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		client, err := chlib.NewClient(db, helpers.CurrentClientVersion, helpers.UuidV4(), np)
		if err != nil {
			np.ERROR.Println(err)
			return
		}
		resp, err := client.GetVolume(volumeName)
		if err != nil {
			np.ERROR.Println(err)
			return
		}
		respjson, _ := json.Marshal(resp)
		var p requestresults.ResultPrinter
		if volumeName == "" {
			var volList requestresults.VolumeListResult
			err = json.Unmarshal(respjson, &volList)
			p = volList
		} else {
			var singleVol requestresults.SingleVolumeResult
			err = json.Unmarshal(respjson, &singleVol)
			p = singleVol
		}
		if err != nil {
			np.ERROR.Println(err)
			return
		}
		if err := p.Print(); err != nil {
			np.ERROR.Println(err)
		}
	},
}

func init() {
	RootCmd.AddCommand(volumeCmd)
}
