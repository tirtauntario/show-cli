package show

import (
	"sort"
	"strings"

	"github.com/alecthomas/chroma/v2"
	"github.com/alecthomas/chroma/v2/lexers"
)

func detectFileTypeFromExtension(path string) string {
	lexer := lexers.Match(path)
	if lexer == nil {
		return "unknown"
	}
	lexer = chroma.Coalesce(lexer)
	if lexer == nil {
		return "unknown"
	}
	if config := lexer.Config(); config != nil && config.Name != "" {
		return config.Name
	}
	return "unknown"
}

func SupportedFileTypes() []string {
	canonical := make(map[string]struct{})
	for _, name := range lexers.Names(false) {
		canonical[strings.ToLower(name)] = struct{}{}
	}

	var aliases []string
	for _, name := range lexers.Names(true) {
		lower := strings.ToLower(name)
		if _, ok := canonical[lower]; ok {
			continue
		}
		aliases = append(aliases, lower)
	}
	sort.Strings(aliases)
	return aliases
}
