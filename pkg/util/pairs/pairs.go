package pairs

import (
	"fmt"
	"regexp"
	"strings"
)

type Pair struct {
	Key   string
	Value string
}

func Parse(txt, kvDelim string) ([]Pair, error) {
	var pairs []Pair
	re, err := regexp.Compile(fmt.Sprintf(`(\w+)\s*%q\s*(\w+|"(.+?[^\\])")`, kvDelim))
	if err != nil {
		return nil, err
	}
	for _, match := range re.FindAllString(txt, -1) {
		kv := strings.SplitN(match, kvDelim, 1)
		if len(kv) == 2 {
			pairs = append(pairs, Pair{
				Key:   kv[0],
				Value: kv[1],
			})
		}
	}
	return pairs, nil
}

func ParseMap(txt, kvDelim string) (map[string]string, error) {
	pairs, err := Parse(txt, kvDelim)
	if err != nil {
		return nil, err
	}
	pairsMap := make(map[string]string, len(pairs))
	for _, pair := range pairs {
		pairsMap[pair.Key] = pair.Value
	}
	return pairsMap, nil
}
