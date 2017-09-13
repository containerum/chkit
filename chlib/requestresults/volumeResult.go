package requestresults

import "fmt"

type SingleVolumeResult struct {
	Replica int    `json:"replica"`
	Status  string `json:"status"`
	Label   string `json:"label"`
	Df      struct {
		Available  string `json:"available"`
		Total      string `json:"total"`
		UsePercent string `json:"use_percent"`
		Used       string `json:"used"`
	} `json:"df"`
}

type VolumeListResult []struct {
	Label  string `json:"label"`
	Size   int    `json:"size"`
	Status string `json:"status"`
}

var VolumeColumns = []string{"LABEL", "SIZE (GiB)", "STATUS"}

func (l VolumeListResult) Print() error {
	var ppc prettyPrintConfig
	ppc.Columns = VolumeColumns
	for _, vol := range l {
		row := []string{
			vol.Label,
			fmt.Sprintf("%d", vol.Size),
			vol.Status,
		}
		ppc.Data = append(ppc.Data, row)
	}
	return ppc.Print()
}

func (s SingleVolumeResult) Print() error {
	fmt.Printf("%-30s %s\n", "Label:", s.Label)
	fmt.Printf("%-30s %d\n", "Replicas:", s.Replica)
	fmt.Printf("%-30s %s\n", "Status:", s.Status)
	fmt.Printf("%-30s %s\n", "Capacity:", s.Df.Total)
	fmt.Printf("%-30s %s (%s)\n", "Used:", s.Df.Used, s.Df.UsePercent)
	fmt.Printf("%-30s %s\n", "Available:", s.Df.Available)
	return nil
}
