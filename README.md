# shcopy

<p>
    <a href="https://github.com/aymanbagabas/shcopy/actions"><img src="https://github.com/aymanbagabas/shcopy/workflows/build/badge.svg" alt="Build Status"></a>
    <a href="https://github.com/aymanbagabas/go-osc52/releases"><img src="https://img.shields.io/github/release/aymanbagabas/go-osc52.svg" alt="Latest Release"></a>
</p>

**Sh**ell **Copy** is a simple utility that copies text to the clipboard from anywhere using [ANSI OSC52](https://invisible-island.net/xterm/ctlseqs/ctlseqs.html#h3-Operating-System-Commands) sequence. It works with local terminals (/dev/tty*) and remote terminals (SSH, Telnet).

## Installation

### Go Install

```sh
go install github.com/aymanbagabas/shcopy@latest
```

### Homebrew

```sh
brew install aymanbagabas/tap/shcopy
```

## Supported Terminals

This is a non-exhaustive list of the status of popular terminal emulators regarding OSC52 [^1]:

| Terminal | OSC52 support |
|----------|:-------------:|
| [Alacritty](https://github.com/alacritty/alacritty) | **yes** |
| [foot](https://codeberg.org/dnkl/foot) | **yes** |
| [GNOME Terminal](https://github.com/GNOME/gnome-terminal) (and other VTE-based terminals) | [not yet](https://bugzilla.gnome.org/show_bug.cgi?id=795774) |
| [hterm (Chromebook)](https://chromium.googlesource.com/apps/libapps/+/master/README.md) | [**yes**](https://chromium.googlesource.com/apps/libapps/+/master/nassh/doc/FAQ.md#Is-OSC-52-aka-clipboard-operations_supported) |
| [iTerm2](https://iterm2.com/) | **yes** |
| [kitty](https://github.com/kovidgoyal/kitty) | **yes** |
| [Konsole](https://konsole.kde.org/) | [not yet](https://bugs.kde.org/show_bug.cgi?id=372116) |
| [QTerminal](https://github.com/lxqt/qterminal#readme) | [not yet](https://github.com/lxqt/qterminal/issues/839)
| [screen](https://www.gnu.org/software/screen/) | **yes** |
| [st](https://st.suckless.org/) | **yes** (but needs to be enabled, see [here](https://git.suckless.org/st/commit/a2a704492b9f4d2408d180f7aeeacf4c789a1d67.html)) |
| [Terminal.app](https://en.wikipedia.org/wiki/Terminal_(macOS)) | no, but see [workaround](https://github.com/roy2220/osc52pty) |
| [tmux](https://github.com/tmux/tmux) | **yes** |
| [Windows Terminal](https://github.com/microsoft/terminal) | **yes** |
| [rxvt](http://rxvt.sourceforge.net/) | **yes** (to be confirmed) |
| [urxvt](http://software.schmorp.de/pkg/rxvt-unicode.html) | **yes** (with a script, see [here](https://github.com/ojroques/vim-oscyank/issues/4)) |
| [xterm.js](https://xtermjs.org/) (Hyper terminal) | [not yet](https://github.com/xtermjs/xterm.js/issues/3260) |
| [wezterm](https://github.com/wez/wezterm) | [**yes**](https://wezfurlong.org/wezterm/escape-sequences.html#operating-system-command-sequences) |

[^1]: Originally copied from [vim-oscyank](https://github.com/ojroques/vim-oscyank)

### Tmux

Starting with tmux 3.3, the `allow-passthrough` option is disabled by default. This means that `shcopy` will not work in tmux by default. To enable it, add the following to your tmux config:

```tmux
set -g allow-passthrough on
```

or use the following if you have `set-clipboard on` in your `~/.tmux.conf`:

```sh
# set the terminal to default
echo "Hello, World!" | shcopy -t default
```

## Credits

This project is built on top of [go-osc52](https://github.com/aymanbagabas/go-osc52), based on [vim-oscyank](https://github.com/ojroques/vim-oscyank).