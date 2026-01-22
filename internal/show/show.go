package show

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

type Deps struct {
	FileReader FileReader
}

type FileReader interface {
	ReadFile(path string) ([]byte, error)
}

type OSFileReader struct{}

func (OSFileReader) ReadFile(path string) ([]byte, error) {
	return os.ReadFile(path)
}

type ShowOptions struct {
	Path     string
	FileType string
	Theme    string
	Debug    bool
}

type ShowResult struct {
	Content []byte
}

func RunShow(ctx context.Context, deps Deps, opts ShowOptions) (ShowResult, error) {
	if opts.Path == "" {
		return ShowResult{}, errors.New("path is required")
	}
	if deps.FileReader == nil {
		return ShowResult{}, errors.New("file reader is required")
	}
	if opts.Theme != "" && !IsSupportedTheme(opts.Theme) {
		return ShowResult{}, fmt.Errorf("unknown theme: %s", opts.Theme)
	}

	data, err := deps.FileReader.ReadFile(opts.Path)
	if err != nil {
		return ShowResult{}, fmt.Errorf("read file: %w", err)
	}

	content := string(data)
	highlighted, err := highlightContent(opts.Path, content, opts.FileType, opts.Theme)
	if err != nil {
		return ShowResult{}, fmt.Errorf("highlight content: %w", err)
	}

	content = addLineNumbers(highlighted)
	if opts.Debug {
		content = wrapWithDebugFileType(opts.Path, data, content)
	}
	_ = ctx
	return ShowResult{Content: []byte(content)}, nil
}

func wrapWithDebugFileType(path string, data []byte, content string) string {
	fileType := detectFileType(path, data)
	line := fmt.Sprintf("DEBUG file type: %s", fileType)

	var b strings.Builder
	b.WriteString(line)
	b.WriteByte('\n')
	b.WriteString(content)
	if content != "" && !strings.HasSuffix(content, "\n") {
		b.WriteByte('\n')
	}
	b.WriteByte('\n')
	b.WriteString(line)
	b.WriteByte('\n')
	b.WriteByte('\n')
	return b.String()
}

func detectFileType(path string, data []byte) string {
	return detectFileTypeFromExtension(path)
}

func addLineNumbers(input string) string {
	if input == "" {
		return ""
	}

	width := lineNumberWidth(input)
	var b strings.Builder
	reader := bufio.NewReader(strings.NewReader(input))
	lineNum := 1
	lineStarted := false
	sep := lineSeparator()
	reset := "\x1b[0m"
	white := "\x1b[37m"
	useColor := !noColor()
	for {
		line, err := reader.ReadString('\n')
		if line != "" {
			if errors.Is(err, io.EOF) && !strings.HasSuffix(line, "\n") {
				b.WriteString(line)
				break
			}
			if !lineStarted {
				if useColor {
					fmt.Fprintf(&b, "%s%s%*d %s%s ", reset, white, width, lineNum, sep, reset)
				} else {
					fmt.Fprintf(&b, "%*d %s ", width, lineNum, sep)
				}
				lineStarted = true
			}
			b.WriteString(line)
			if strings.HasSuffix(line, "\n") {
				lineNum++
				lineStarted = false
			}
		}
		if errors.Is(err, bufio.ErrBufferFull) {
			continue
		}
		if errors.Is(err, io.EOF) || err != nil {
			break
		}
	}
	return b.String()
}

func lineSeparator() string {
	if isUTF8Locale() {
		return "│"
	}
	return "|"
}

func isUTF8Locale() bool {
	return isUTF8Env("LC_ALL") || isUTF8Env("LC_CTYPE") || isUTF8Env("LANG")
}

func isUTF8Env(key string) bool {
	value := os.Getenv(key)
	value = strings.ToLower(value)
	return strings.Contains(value, "utf-8") || strings.Contains(value, "utf8")
}

func noColor() bool {
	value, ok := os.LookupEnv("NO_COLOR")
	return ok && value != ""
}

func lineNumberWidth(input string) int {
	if input == "" {
		return 1
	}
	lines := strings.Count(input, "\n")
	if !strings.HasSuffix(input, "\n") {
		lines++
	}
	if lines == 0 {
		lines = 1
	}
	return len(strconv.Itoa(lines))
}
