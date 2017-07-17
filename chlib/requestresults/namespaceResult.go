package requestresults

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/kfeofantov/chkit-v2/chlib"
	"github.com/olekukonko/tablewriter"
)

type namespaceResult []struct {
	chlib.Namespace
}

func (n namespaceResult) Print() error {
	if len(n) < 2 {
		return fmt.Errorf("Invalid namespace response")
	}
	name := n[0].Data.Metadata.Name
	phase := n[0].Data.Status.Phase
	createdAt := n[0].Data.Metadata.CreationTimestamp
	hard := n[1].Data.Status.Hard
	used := n[1].Data.Status.Used
	fmt.Printf("%-20s %s\n", "Name:", name)
	fmt.Printf("%-20s %s\n", "Phase:", phase)
	fmt.Printf("%-20s %s\n", "Creation time:", createdAt.Format(time.RFC1123))
	fmt.Println("Hard:")
	fmt.Printf("\t%-20s %s\n", "CPU", hard.RequestsCPU)
	fmt.Printf("\t%-20s %s\n", "Memory", hard.RequestsMemory)
	fmt.Println("Used:")
	fmt.Printf("\t%-20s %s\n", "CPU", used.RequestsCPU)
	fmt.Printf("\t%-20s %s\n", "Memory", used.RequestsMemory)
	return nil
}

func (n namespaceResult) formatPrettyPrint() (ppc prettyPrintConfig) {
	ppc.Columns = []string{"NAME", "HARD CPU", "HARD MEMORY", "USED CPU", "USED MEMORY", "AGE"}
	for _, v := range n {
		row := []string{
			v.Data.Metadata.Namespace,
			v.Data.Status.Hard.LimitsCPU,
			v.Data.Status.Hard.LimitsMemory,
			v.Data.Status.Used.LimitsCPU,
			v.Data.Status.Used.LimitsMemory,
			fmt.Sprintf("%dd", int(time.Now().Sub(*v.Data.Metadata.CreationTimestamp).Hours()/24)),
		}
		ppc.Data = append(ppc.Data, row)
	}
	ppc.Align = tablewriter.ALIGN_LEFT
	return
}

func init() {
	resultKinds["ResourceQuota"] = func(resp []chlib.GenericJson) (ResultPrinter, error) {
		var res namespaceResult
		b, _ := json.Marshal(resp)
		if err := json.Unmarshal(b, &res); err != nil {
			return nil, fmt.Errorf("invalid namespace list response: %s", err)
		}
		return res.formatPrettyPrint(), nil
	}
	resultKinds["Namespace"] = func(resp []chlib.GenericJson) (ResultPrinter, error) {
		var res namespaceResult
		b, _ := json.Marshal(resp)
		if err := json.Unmarshal(b, &res); err != nil {
			return res, fmt.Errorf("invalid namespace response: %s", err)
		}
		return res, nil
	}
}
