package show

import (
	"bytes"
	"fmt"

	"github.com/alecthomas/chroma/v2"
	"github.com/alecthomas/chroma/v2/formatters"
	"github.com/alecthomas/chroma/v2/lexers"
	"github.com/alecthomas/chroma/v2/styles"
)

func highlightContent(path string, content string, fileType string, theme string) (string, error) {
	var lexer chroma.Lexer
	if fileType != "" {
		lexer = lexers.Get(fileType)
		if lexer == nil {
			return "", fmt.Errorf("unknown file type: %s", fileType)
		}
	} else {
		lexer = lexers.Match(path)
		if lexer == nil {
			lexer = lexers.Analyse(content)
		}
	}
	if lexer == nil {
		lexer = lexers.Fallback
	}
	lexer = chromaCoalesce(lexer)

	formatter := formatters.Get("terminal16m")
	if formatter == nil {
		formatter = formatters.Get("terminal256")
	}
	if formatter == nil {
		formatter = formatters.Get("terminal")
	}
	if formatter == nil {
		return "", ErrNoFormatter
	}

	if theme == "" {
		theme = "onedark"
	}
	style := styles.Get(theme)
	if style == nil {
		style = styles.Fallback
	}

	iterator, err := lexer.Tokenise(nil, content)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err := formatter.Format(&buf, style, iterator); err != nil {
		return "", err
	}
	return buf.String(), nil
}

var ErrNoFormatter = errNoFormatter{}

type errNoFormatter struct{}

func (errNoFormatter) Error() string {
	return "no terminal formatter available"
}

func chromaCoalesce(lexer chroma.Lexer) chroma.Lexer {
	if lexer == nil {
		return nil
	}
	return chroma.Coalesce(lexer)
}
