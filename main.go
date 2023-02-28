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

	term = pflag.StringP("term", "t", "", "Terminal type: (default), tmux, screen.\n"+
		"Note: starting with tmux 3.3a, tmux users won't need to set this option.\n"+
		"Refer to the man page for more information.")
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
			WithLongDescription(ProjectName+" a utility that copies text to your clipboard from anywhere using ANSI OSC52 sequence.").
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
