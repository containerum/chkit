package requestresults

import (
	"fmt"
	"os"
	"time"

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

func ageFormat(d time.Duration) string {
	durations := []struct {
		num float64
		pf  string
	}{
		{d.Hours() / 24 / 365, "Y"},
		{d.Hours() / 24 / 30, "M"},
		{d.Hours() / 24, "d"},
		{d.Hours(), "d"},
		{d.Minutes(), "m"},
		{d.Seconds(), "s"},
	}
	for _, v := range durations {
		if v.num > 1 {
			return fmt.Sprintf("%d%s", int(v.num), v.pf)
		}
	}
	return ""
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
