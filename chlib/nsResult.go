package chlib

import (
	"encoding/json"
	"fmt"
	"time"
)

type nsResult struct {
	Data struct {
		Metadata struct {
			CreatedAt time.Time `json:"creationTimestamp"`
			Namespace string    `json:"namespace"`
		} `json:"metadata"`
		Status struct {
			Hard struct {
				LimitsCpu string `json:"limits.cpu"`
				LimitsMem string `json:"limits.memory"`
			} `json:"hard"`
			Used struct {
				LimitsCpu string `json:"limits.cpu"`
				LimitsMem string `json:"limits.memory"`
			} `json:"used"`
		} `json:"status"`
	} `json:"data"`
}

func extractNsResults(data []GenericJson) (res []nsResult, err error) {
	b, _ := json.Marshal(data)
	if err := json.Unmarshal(b, &res); err != nil {
		return res, fmt.Errorf("invalid namespace response: %s", err)
	}
	return
}

func formatNamespacePrettyPrint(data []nsResult) (ppc PrettyPrintConfig) {
	ppc.Columns = []string{"NAME", "HARD CPU", "HARD MEMORY", "USED CPU", "USED MEMORY", "AGE"}
	for _, v := range data {
		row := []string{
			v.Data.Metadata.Namespace,
			v.Data.Status.Hard.LimitsCpu,
			v.Data.Status.Hard.LimitsMem,
			v.Data.Status.Used.LimitsCpu,
			v.Data.Status.Used.LimitsMem,
			fmt.Sprintf("%dd", int(time.Now().Sub(v.Data.Metadata.CreatedAt).Hours()/24)),
		}
		ppc.Data = append(ppc.Data, row)
	}
	return
}
