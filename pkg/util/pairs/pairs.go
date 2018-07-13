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

var (
	wordRe = regexp.MustCompile(`^\w+$`)
)

func Parse(txt, kvDelim string) ([]Pair, error) {
	var pairs []Pair
	re, err := regexp.Compile(fmt.Sprintf(`(\w+)\s*%s\s*(\w+|"(.+?[^\\])")`, kvDelim))
	if err != nil {
		return nil, err
	}
	for _, match := range re.FindAllString(txt, -1) {
		kv := strings.SplitN(match, kvDelim, 2)
		if len(kv) == 2 {
			v := kv[1]
			if !wordRe.MatchString(v) {
				v = strings.TrimPrefix(v, "\"")
				v = strings.TrimSuffix(v, "\"")
			}
			pairs = append(pairs, Pair{
				Key:   kv[0],
				Value: v,
			})
		} else {
			panic(fmt.Sprintf("[chkit/pkg/util/pairs.Parse] unexpected number of tokens %d in %+v", len(kv), kv))
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
