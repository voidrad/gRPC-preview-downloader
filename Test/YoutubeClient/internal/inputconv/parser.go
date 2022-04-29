package inputconv

import (
	"log"
	"strings"
)

func ParseCommand(in []string) (out []string, async bool) {
	async = false
	if in[len(in)-1] == "--async" {
		async = true
		in = in[0 : len(in)-1]
	}
	for i := range in {
		linkParsed := strings.Split(in[i], "?v=")
		if (linkParsed[0] != "www.youtube.com/watch") || (len(in[i]) == len(linkParsed[0])) {
			log.Println("invalid link", in[i])
			continue
		}
		out = append(out, linkParsed[1])
	}
	return out, async
}
