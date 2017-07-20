package requestresults

import (
	"fmt"
	"os"
	"time"

	"chkit-v2/chlib"
	"github.com/olekukonko/tablewriter"
	jww "github.com/spf13/jwalterweatherman"
	"sort"
	"strconv"
	"strings"
)

type responseProcessor func(resp []chlib.GenericJson) (ResultPrinter, error)

var resultKinds = make(map[string]responseProcessor)

var fieldToSort string

type prettyPrintConfig struct {
	Columns []string
	Data    [][]string
	Align   int
}

type ResultPrinter interface {
	Print() error
}

func (p prettyPrintConfig) getFieldIndex(fieldName string) int {
	if fieldName == "" {
		return -1
	}
	for k, v := range p.Columns {
		if v == fieldName {
			return k
		}
	}
	return -1
}

// compareAges returns -1 if age1 < age2, 0 if age1 == age2, +1 if age1 > age2
func compareAges(age1, age2 string) int {
	ageTypes := map[byte]int{'s': 0, 'm': 1, 'h': 2, 'd': 3, 'M': 4, 'Y': 5}
	age1Pf := age1[len(age1)-1]
	age2Pf := age2[len(age2)-1]
	if ageTypes[age1Pf] > ageTypes[age2Pf] {
		return 1
	} else if ageTypes[age1Pf] < ageTypes[age2Pf] {
		return -1
	} else {
		age1D, _ := strconv.Atoi(age1[:len(age1)-1])
		age2D, _ := strconv.Atoi(age2[:len(age2)-1])
		diff := age1D - age2D
		if diff < 0 {
			return -1
		} else if diff > 0 {
			return 1
		}
		return 0
	}
}

func (p prettyPrintConfig) sortByField(fieldName string) {
	fieldIndex := p.getFieldIndex(fieldName)
	if fieldIndex == -1 {
		return
	}
	if fieldName == "AGE" {
		sort.Slice(p.Data, func(i, j int) bool {
			return compareAges(p.Data[i][fieldIndex], p.Data[j][fieldIndex]) <= 0
		})
	} else {
		sort.Slice(p.Data, func(i, j int) bool {
			return strings.Compare(p.Data[i][fieldIndex], p.Data[j][fieldIndex]) <= 0
		})
	}
}

func (p prettyPrintConfig) Print() error {
	table := tablewriter.NewWriter(os.Stdout)
	p.sortByField(fieldToSort)
	table.SetHeader(p.Columns)
	table.AppendBulk(p.Data)
	table.SetAlignment(p.Align)
	table.Render()
	return nil
}

func ageFormat(d time.Duration) string {
	durations := []struct {
		num float64
		pf  rune
	}{
		{d.Hours() / 24 / 365, 'Y'},
		{d.Hours() / 24 / 30, 'M'},
		{d.Hours() / 24, 'd'},
		{d.Hours(), 'h'},
		{d.Minutes(), 'm'},
		{d.Seconds(), 's'},
	}
	for _, v := range durations {
		if v.num > 1 {
			return fmt.Sprintf("%d%c", int(v.num), v.pf)
		}
	}
	return ""
}

func ProcessResponse(resp []chlib.GenericJson, sortField string, np *jww.Notepad) (res ResultPrinter, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("Can`t find kind in response")
		}
	}()
	np.SetPrefix("ResponseProcessor")
	fieldToSort = sortField
	kind := resp[0]["data"].(map[string]interface{})["kind"].(string)
	np.DEBUG.Println("Response kind", kind)
	prc, ok := resultKinds[kind]
	if !ok {
		return nil, fmt.Errorf("kind %s is not registered", kind)
	}
	return prc(resp)
}
