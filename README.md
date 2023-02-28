# shcopy

<p>
    <a href="https://github.com/aymanbagabas/shcopy/actions"><img src="https://github.com/aymanbagabas/shcopy/workflows/build/badge.svg" alt="Build Status"></a>
    <a href="https://github.com/aymanbagabas/shcopy/releases"><img src="https://img.shields.io/github/release/aymanbagabas/shcopy.svg" alt="Latest Release"></a>
</p>

**Sh**ell **Copy** is a simple utility that copies text to the clipboard from anywhere using [ANSI OSC52](https://invisible-island.net/xterm/ctlseqs/ctlseqs.html#h3-Operating-System-Commands) sequence. It works with local terminals (/dev/tty*) and remote terminals (SSH, Telnet).

Think of this as a tool like `xclip` or `pbcopy` but also works over SSH.

## Example

```sh
# Copy text to clipboard
shcopy "Hello World"

# Copy text to primary clipboard (X11)
shcopy -p "Hello World"

# Copy command output to clipboard
echo -n "Hello World" | shcopy

# Copy file content to clipboard
shcopy < file.txt

# Copy from stdin until EOF
# Ctrl+D to finish
shcopy

# Need help?
shcopy --help
```

## Installation

### Go Install

```sh
go install github.com/aymanbagabas/shcopy@latest
```

### Homebrew

```sh
brew install aymanbagabas/tap/shcopy
```

### Debian/Ubuntu

```sh
echo 'deb [trusted=yes] https://repo.aymanbagabas.com/apt/ /' | sudo tee /etc/apt/sources.list.d/aymanbagabas.list
sudo apt update && sudo apt install shcopy
```

### Fedora

```sh
echo '[aymanbagabas]
name=Ayman Bagabas
baseurl=https://repo.aymanbagabas.com/yum/
enabled=1
gpgcheck=0' | sudo tee /etc/yum.repos.d/aymanbagabas.repo
sudo yum install shcopy
```

### Arch Linux

```sh
yay -S shcopy-bin
```

### Scoop (Windows)

```sh
scoop bucket add aymanbagabas https://github.com/aymanbagabas/scoop-bucket.git
scoop install aymanbagabas/shcopy
```

You can also download the latest binaries and packages from the [releases page](https://github.com/aymanbagabas/shcopy/releases).

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

To use shcopy within a tmux session, make sure that the outer terminal supports
OSC 52, and use one of the following options:

1. Configure tmux to allow programs to access the clipboard (recommended). The
   tmux `set-clipboard` option was added in tmux 1.5 with a default of `on`;
   the default was changed to `external` when `external` was added in tmux 2.6.
   Setting `set-clipboard` to `on` allows external programs in tmux to access
   the clipboard. To enable this option, add `set -s set-clipboard on` to your
   tmux config.

2. Use `--term tmux` option to force shcopy to work with tmux. This option
   requires the `allow-passthrough` option to be enabled in tmux. Starting with
   tmux 3.3a, the `allow-passthrough` option is no longer enabled by default.
   This option allows tmux to pass an ANSI escape sequence to the outer
   terminal by wrapping it in another special tmux escape sequence. This means
   the `--term tmux` option won't work unless you're running an older version
   of tmux or you have enabled `allow-passthrough` in tmux. Add the following
   to your tmux config to enable passthrough `set -g allow-passthrough on`.

Refer to https://github.com/tmux/tmux/wiki/Clipboard for more info.

## Credits

This project is built on top of [go-osc52](https://github.com/aymanbagabas/go-osc52), based on [vim-oscyank](https://github.com/ojroques/vim-oscyank).
