package show

import "github.com/alecthomas/chroma/v2/styles"

func SupportedThemes() []string {
	return styles.Names()
}

func IsSupportedTheme(name string) bool {
	for _, theme := range styles.Names() {
		if theme == name {
			return true
		}
	}
	return false
}
