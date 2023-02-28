package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/aymanbagabas/go-osc52/v2"
	"github.com/muesli/mango"
	mpflag "github.com/muesli/mango-pflag"
	"github.com/muesli/roff"
	"github.com/spf13/pflag"
)

var (
	ProjectName = "shcopy"
	Version     = "unknown"
	CommitSHA   = "build from source"

	term    = pflag.StringP("term", "t", "", "Terminal type: (default), tmux, screen.")
	clear   = pflag.BoolP("clear", "c", false, "Clear the clipboard and exit.")
	primary = pflag.BoolP("primary", "p", false, "Use the primary clipboard instead system clipboard.")
	version = pflag.BoolP("version", "v", false, "Print version and exit.")
	help    = pflag.BoolP("help", "h", false, "Print help and exit.")
	debug   = pflag.BoolP("debug", "d", false, "Print debug information.")
	man     = pflag.Bool("man", false, "Generate man pages.")
)

func usage(isError bool) {
	out := os.Stdout
	if isError {
		out = os.Stderr
	}
	fmt.Fprintf(out, `Usage:
  %[1]s [options] [text]
  %[1]s [options] < [file]

  Copy text to the system clipboard from any supported terminal using
  ANSI OSC 52 sequence.

Options:
`, ProjectName)
	pflag.PrintDefaults()
}

func main() {
	pflag.Usage = func() {
		usage(true)
	}
	pflag.Parse()
	pflag.Lookup("debug").Hidden = true
	pflag.Lookup("man").Hidden = true

	if *version {
		fmt.Printf("%s version %s (%s)", ProjectName, Version, CommitSHA)
		return
	}

	if *help {
		usage(false)
		return
	}

	if *man {
		manPage := mango.NewManPage(1, ProjectName, "Copy text to the system clipboard from any supported terminal using ANSI OSC 52 sequence.").
			WithLongDescription(ProjectName+" a utility that copies text to your clipboard from anywhere using ANSI OSC 52 sequence.").
			WithSection("Terminal", "shcopy should work in any terminal that supports OSC 52. There are some exceptions below.").
			WithSection("Kitty", `
Kitty, version 0.22.0 and below, had a bug where it appends to the clipboard instead of replacing it. To workaround this bug, clear the clipboard before copying any text.

shcopy -c; shcopy "Hello World"
`,
			).
			WithSection("Screen", `
To use shcopy within a screen session, make sure that the outer terminal supports OSC 52. If your '$TERM' environment variable is not set to 'screen-*', use '--term screen' to force shcopy to work with screen.
`,
			).
			WithSection("Tmux", `
To use shcopy within a tmux session, make sure that the outer terminal supports OSC 52, and use one of the following options:

1. Configure tmux to allow programs to access the clipboard (recommended). The tmux 'set-clipboard' option was added in tmux 1.5 with a default of 'on'; the default was changed to 'external' when 'external' was added in tmux 2.6. Setting 'set-clipboard' to 'on' allows external programs in tmux to access the clipboard. To enable this option, add 'set -s set-clipboard on' to your tmux config.

2. Use '--term tmux' option to force shcopy to work with tmux. This option requires the 'allow-passthrough' option to be enabled in tmux. Starting with tmux 3.3a, the 'allow-passthrough' option is no longer enabled by default. This option allows tmux to pass an ANSI escape sequence to the outer terminal by wrapping it in another special tmux escape sequence. This means the '--term tmux' option won't work unless you're running an older version of tmux or you have enabled 'allow-passthrough' in tmux. Add the following to your tmux config to enable passthrough 'set -g allow-passthrough on'.

Refer to https://github.com/tmux/tmux/wiki/Clipboard for more info.
`,
			).
			WithSection("Bugs", "Report bugs to https://github.com/aymanbagabas/shcopy/issues").
			WithSection("Copyright", "(C) 2023 Ayman Bagabas.\n"+
				"Released under MIT license.")

		pflag.VisitAll(mpflag.PFlagVisitor(manPage))
		fmt.Println(manPage.Build(roff.NewDocument()))
		return
	}

	var str string
	args := pflag.Args()
	// read from stdin if no arguments are provided and we are not clearing the
	// clipboard or reading the clipboard contents.
	if len(args) == 0 && !*clear {
		reader := bufio.NewReader(os.Stdin)
		var b strings.Builder

		for {
			r, _, err := reader.ReadRune()
			if err != nil && err == io.EOF {
				break
			}
			_, err = b.WriteRune(r)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Failed to write rune: %v", err)
				os.Exit(1)
			}
		}

		// input
		str = b.String()
	} else {
		str = strings.Join(args, " ")
	}

	// the sequence string to be sent to the terminal
	seq := osc52.New(str)

	if *primary {
		seq = seq.Primary()
	}

	// Detect `screen` terminal type
	if term := os.Getenv("TERM"); term != "" {
		if strings.HasPrefix(term, "screen") {
			seq = seq.Screen()
		}
	}

	if *term != "" {
		switch strings.ToLower(*term) {
		case "screen":
			seq = seq.Screen()
		case "tmux":
			seq = seq.Tmux()
		default:
			seq = seq.Mode(osc52.DefaultMode)
		}
	}

	if *clear {
		seq = seq.Clear()
	}

	if *debug {
		log.Printf("Sequence: %q", seq)
	}

	// send the sequence to the terminal
	_, _ = seq.WriteTo(os.Stderr)
}
