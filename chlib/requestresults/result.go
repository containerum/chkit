package requestresults

import (
	"fmt"
	"os"

	"github.com/kfeofantov/chkit-v2/chlib"
	"github.com/olekukonko/tablewriter"
)

type responseProcessor func(resp []chlib.GenericJson) (ResultPrinter, error)

var resultKinds = make(map[string]responseProcessor)

type prettyPrintConfig struct {
	Columns []string
	Data    [][]string
	Align   int
}

type ResultPrinter interface {
	Print() error
}

func (p prettyPrintConfig) Print() error {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader(p.Columns)
	table.AppendBulk(p.Data)
	table.SetAlignment(p.Align)
	table.Render()
	return nil
}

func ProcessResponse(resp []chlib.GenericJson) (res ResultPrinter, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("Can`t find kind in response")
		}
	}()
	kind := resp[0]["data"].(map[string]interface{})["kind"].(string)
	prc, ok := resultKinds[kind]
	if !ok {
		return nil, fmt.Errorf("kind %s is not registered", kind)
	}
	return prc(resp)
}
