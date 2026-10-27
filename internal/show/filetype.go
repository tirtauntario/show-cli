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

// SupportedFileTypes returns every value accepted by --filetype: all lexer
// aliases plus lexer names that contain no spaces, lowercased and sorted.
func SupportedFileTypes() []string {
	seen := make(map[string]struct{})
	add := func(name string) {
		name = strings.ToLower(name)
		if name == "" || strings.ContainsAny(name, " \t") {
			return
		}
		seen[name] = struct{}{}
	}
	for _, lexer := range lexers.GlobalLexerRegistry.Lexers {
		config := lexer.Config()
		if config == nil {
			continue
		}
		add(config.Name)
		for _, alias := range config.Aliases {
			add(alias)
		}
	}

	types := make([]string, 0, len(seen))
	for name := range seen {
		types = append(types, name)
	}
	sort.Strings(types)
	return types
}
