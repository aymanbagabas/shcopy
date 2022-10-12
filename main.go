package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/aymanbagabas/go-osc52"
	"github.com/spf13/pflag"
)

var (
	ProjectName = "shcopy"
	Version     = "unknown"
	CommitSHA   = "build from source"

	term    = pflag.StringP("term", "t", os.Getenv("TERM"), "Terminal type: screen, tmux, etc.")
	clear   = pflag.BoolP("clear", "c", false, "Clear the clipboard.")
	primary = pflag.BoolP("primary", "p", false, "Use the primary clipboard instead system clipboard.")
	version = pflag.BoolP("version", "v", false, "Print version and exit.")
	help    = pflag.BoolP("help", "h", false, "Print help and exit.")
	debug   = pflag.BoolP("debug", "d", false, "Print debug information.")
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

	if *version {
		fmt.Printf("%s version %s (%s)", ProjectName, Version, CommitSHA)
		os.Exit(0)
	}

	if *help {
		usage(false)
		os.Exit(0)
	}

	var str string
	args := pflag.Args()
	// read from stdin if no arguments are provided and we are not clearing the
	// clipboard
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
				fmt.Printf("Failed to write rune: %v", err)
				os.Exit(1)
			}
		}

		// input
		str = b.String()
	} else {
		str = strings.Join(args, " ")
	}

	if *debug {
		log.Printf("Input: %q", str)
	}
	term := strings.ToLower(*term)
	clip := osc52.ClipboardC
	if *primary {
		clip = osc52.ClipboardP
	}
	if *debug {
		log.Printf("Clipboard: %v", clip)
	}
	if strings.Contains(term, "kitty") {
		// Flush the keyboard before copying, this is required for
		// Kitty < 0.22.0.
		os.Stdout.WriteString(osc52.Clear(term, clip))
	}
	// the sequence string to be sent to the terminal
	seq := osc52.Sequence(str, term, clip)
	if *debug {
		log.Printf("Sequence: %q", seq)
	}

	// send the sequence to the terminal
	os.Stdout.WriteString(seq)
}
