package activekit

import (
	"strings"

	"fmt"

	"github.com/containerum/chkit/pkg/util/text"
)

func Attention(txt string) {
	border := strings.Repeat("!", text.Width(txt))
	fmt.Printf("%s\n%s\n%s\n", border, txt, border)
}
