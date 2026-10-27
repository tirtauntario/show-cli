package cli

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/urfave/cli/v2"
	"golang.org/x/term"

	"github.com/tirtauntario/show-cli/internal/show"
)

type BuildInfo struct {
	Version string
	Commit  string
	Date    string
}

type CLI struct {
	deps   show.Deps
	info   BuildInfo
	out    io.Writer
	errOut io.Writer
	// isTerminal reports whether w is a terminal. Tests replace it.
	isTerminal func(w io.Writer) bool
}

func New(deps show.Deps, info BuildInfo, out io.Writer, errOut io.Writer) *CLI {
	return &CLI{deps: deps, info: info, out: out, errOut: errOut, isTerminal: isTerminal}
}

// isTerminal reports whether w is a file descriptor attached to a terminal.
func isTerminal(w io.Writer) bool {
	f, ok := w.(interface{ Fd() uintptr })
	return ok && term.IsTerminal(int(f.Fd()))
}

// useColor resolves --color. auto colors only when stdout is a terminal and
// NO_COLOR is unset or empty (https://no-color.org).
func (c *CLI) useColor(mode string) (bool, error) {
	switch strings.ToLower(mode) {
	case "always":
		return true, nil
	case "never":
		return false, nil
	case "auto", "":
		if os.Getenv("NO_COLOR") != "" {
			return false, nil
		}
		return c.isTerminal(c.out), nil
	default:
		return false, fmt.Errorf("invalid --color value %q (want auto, always or never)", mode)
	}
}

func (c *CLI) Run(args []string) error {
	app := &cli.App{
		Name:      "show",
		Usage:     "display text file contents with syntax highlighting and line numbers",
		Writer:    c.out,
		ErrWriter: c.errOut,
		ArgsUsage: "<path>", HideHelp: true,
		HideVersion:     true,
		HideHelpCommand: true,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "t",
				Aliases: []string{"filetype"},
				Usage:   "force syntax highlighting file type (e.g. go, python, yaml)",
			},
			&cli.BoolFlag{
				Name:    "d",
				Aliases: []string{"debug"},
				Usage:   "print debug metadata",
			},
			&cli.BoolFlag{
				Name:    "h",
				Aliases: []string{"help"},
				Usage:   "show help",
			},
			&cli.BoolFlag{
				Name:    "v",
				Aliases: []string{"version"},
				Usage:   "print version",
			},
			&cli.StringFlag{
				Name:  "color",
				Value: "auto",
				Usage: "when to use color: auto, always or never",
			},
			&cli.IntFlag{
				Name:  "tab-width",
				Value: 4,
				Usage: "expand tabs to this many columns (0 keeps tabs)",
			},
			&cli.StringFlag{
				Name:  "theme",
				Usage: "set syntax highlighting theme (default: onedark, see --list-themes)",
			},
			&cli.BoolFlag{
				Name:  "list-file-types",
				Usage: "print supported file type aliases",
			},
			&cli.BoolFlag{
				Name:  "list-themes",
				Usage: "print supported syntax highlighting themes",
			},
			&cli.StringFlag{
				Name:  "install-completion",
				Usage: "print shell completion script (bash|zsh|fish)",
			},
		},
		Action: func(ctx *cli.Context) error {
			if ctx.Bool("help") || ctx.Bool("h") {
				return cli.ShowAppHelp(ctx)
			}
			if ctx.Bool("version") || ctx.Bool("v") {
				return c.printVersion()
			}
			if ctx.IsSet("install-completion") {
				return c.runCompletion(ctx.String("install-completion"))
			}
			if ctx.Bool("list-file-types") {
				return c.runSupportedTypes()
			}
			if ctx.Bool("list-themes") {
				return c.runListThemes()
			}
			return c.runShow(ctx)
		},
	}

	argv := append([]string{"show"}, normalizeArgs(args)...)
	return app.Run(argv)
}

func (c *CLI) runCompletion(shell string) error {
	if shell == "" {
		return errors.New("usage: show --install-completion <bash|zsh|fish>")
	}

	script, err := completionScript(shell)
	if err != nil {
		return err
	}
	_, err = io.WriteString(c.out, script)
	return err
}

func (c *CLI) runShow(ctx *cli.Context) error {
	if ctx.NArg() != 1 {
		return errors.New("usage: show <path>\n-h for help")
	}

	opts := show.ShowOptions{Path: ctx.Args().Get(0)}
	opts.FileType = ctx.String("filetype")
	if opts.FileType == "" {
		opts.FileType = ctx.String("t")
	}
	opts.Theme = ctx.String("theme")
	opts.Debug = ctx.Bool("debug") || ctx.Bool("d")
	opts.TabWidth = ctx.Int("tab-width")
	color, err := c.useColor(ctx.String("color"))
	if err != nil {
		return err
	}
	opts.Color = color
	result, err := show.RunShow(context.Background(), c.deps, opts)
	if err != nil {
		return err
	}

	_, err = c.out.Write(result.Content)
	return err
}

func normalizeArgs(args []string) []string {
	if len(args) == 0 {
		return args
	}

	flags := make([]string, 0, len(args))
	positionals := make([]string, 0, len(args))
	for i := 0; i < len(args); i++ {
		arg := args[i]
		if arg == "--" {
			positionals = append(positionals, args[i+1:]...)
			break
		}
		if strings.HasPrefix(arg, "-") {
			flags = append(flags, arg)
			if flagNeedsValue(arg) && !strings.Contains(arg, "=") && i+1 < len(args) {
				flags = append(flags, args[i+1])
				i++
			}
			continue
		}
		positionals = append(positionals, arg)
	}
	return append(flags, positionals...)
}

func flagNeedsValue(arg string) bool {
	switch arg {
	case "-t", "--filetype", "--install-completion", "--theme", "--tab-width", "--color":
		return true
	default:
		return false
	}
}

func (c *CLI) printVersion() error {
	var b strings.Builder
	fmt.Fprintf(&b, "version: %s\n", c.info.Version)
	if c.info.Commit != "" {
		fmt.Fprintf(&b, "commit: %s\n", c.info.Commit)
	}
	if c.info.Date != "" {
		fmt.Fprintf(&b, "date: %s\n", c.info.Date)
	}
	_, err := io.WriteString(c.out, b.String())
	return err
}

func (c *CLI) runSupportedTypes() error {
	for _, name := range show.SupportedFileTypes() {
		if _, err := fmt.Fprintln(c.out, name); err != nil {
			return err
		}
	}
	return nil
}

func (c *CLI) runListThemes() error {
	for _, name := range show.SupportedThemes() {
		if _, err := fmt.Fprintln(c.out, name); err != nil {
			return err
		}
	}
	return nil
}

func completionScript(shell string) (string, error) {
	switch strings.ToLower(shell) {
	case "bash":
		return bashCompletion(), nil
	case "zsh":
		return zshCompletion(), nil
	case "fish":
		return fishCompletion(), nil
	default:
		return "", fmt.Errorf("unknown shell: %s", shell)
	}
}

func bashCompletion() string {
	return `# bash completion for show
_show() {
  local cur opts
  cur="${COMP_WORDS[COMP_CWORD]}"
  opts="-h --help -v --version -d --debug -t --filetype --color --tab-width --theme --list-file-types --list-themes --install-completion"
  if [[ "$cur" == -* ]]; then
    COMPREPLY=( $(compgen -W "${opts}" -- "${cur}") )
    return 0
  fi
  COMPREPLY=( $(compgen -f -- "${cur}") )
}
complete -F _show show
`
}

func zshCompletion() string {
	return `#compdef show

_arguments \
  '-h[show help]' \
  '--help[show help]' \
  '-v[print version]' \
  '--version[print version]' \
  '-d[print debug metadata]' \
  '--debug[print debug metadata]' \
  '-t[force syntax highlighting file type]:type:' \
  '--filetype[force syntax highlighting file type]:type:' \
  '--color[when to use color]:when:(auto always never)' \
  '--tab-width[expand tabs to this many columns]:width:' \
  '--theme[set syntax highlighting theme]:theme:' \
  '--list-file-types[print supported file type aliases]' \
  '--list-themes[print supported syntax highlighting themes]' \
  '--install-completion[print shell completion script (bash|zsh|fish)]:shell:(bash zsh fish)' \
  '1: :_files'
`
}

func fishCompletion() string {
	return `complete -c show -f
complete -c show -s h -d "show help"
complete -c show -l help -d "show help"
complete -c show -s v -d "print version"
complete -c show -l version -d "print version"
complete -c show -s d -d "print debug metadata"
complete -c show -l debug -d "print debug metadata"
complete -c show -s t -d "force syntax highlighting file type"
complete -c show -l filetype -d "force syntax highlighting file type"
complete -c show -l color -d "when to use color" -xa "auto always never"
complete -c show -l tab-width -x -d "expand tabs to this many columns"
complete -c show -l theme -d "set syntax highlighting theme"
complete -c show -l list-file-types -d "print supported file type aliases"
complete -c show -l list-themes -d "print supported syntax highlighting themes"
complete -c show -l install-completion -d "print shell completion script" -xa "bash zsh fish"
`
}
