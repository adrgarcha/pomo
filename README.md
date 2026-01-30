# Pomo - Beautiful Cross-Platform Pomodoro Timer

An interactive Pomodoro timer for macOS, Linux, and Windows with a beautiful terminal UI, native notifications, and sound alerts.

## ✨ Features

- 🎨 **Beautiful TUI** - Gorgeous terminal interface powered by Bubble Tea
- 🍅 Simple Pomodoro technique implementation
- 🔔 Native notifications on all platforms
- 🎵 Sound alerts (where supported)
- ⚡ Fast and lightweight
- 🎯 Interactive mode with preset time splits (25/5 or 50/10)
- 📊 Real-time progress bars and session tracking
- 🔁 Chain multiple pomodoros together
- ⌨️ Keyboard navigation (vim-style shortcuts supported)
- 🖥️ Cross-platform: macOS, Linux, Windows
- 💻 Minimal dependencies

## 🎥 Preview

```
🍅 Pomodoro Timer

Choose your session:

› 25/5 (25m work, 5m break)
  50/10 (50m work, 10m break)
  Exit

↑/↓: navigate • enter: select • q: quit

✓ Completed sessions: 3
```

When timer is running:
```
🍅 Work Session

25:00

████████████████████████████████░░░░░░░░

Session: 25 minutes
Completed: 3 sessions

q: quit
```

## Installation

### From Source

```bash
# Clone the repository
git clone https://github.com/adrgarcha/pomo.git
cd pomo

# Build
go build -o pomo

# Install (optional)
# macOS/Linux:
sudo mv pomo /usr/local/bin/

# Windows:
# Move pomo.exe to a directory in your PATH
```

### Using Go Install

```bash
go install github.com/adrgarcha/pomo@latest
```

### Using Homebrew (macOS/Linux)

```bash
brew tap adrgarcha/tap
brew install pomo
```

### Pre-built Binaries

Download from the [releases page](https://github.com/adrgarcha/pomo/releases).

## Platform-Specific Setup

### macOS
No additional setup required! Uses built-in `osascript` for notifications.

### Linux
For desktop notifications, install `libnotify`:

```bash
# Ubuntu/Debian
sudo apt install libnotify-bin

# Fedora
sudo dnf install libnotify

# Arch
sudo pacman -S libnotify
```

### Windows
No additional setup required! Uses PowerShell for toast notifications.

## Usage

### Interactive TUI Mode (Default)

Simply run:

```bash
pomo
```

Navigate with:
- **↑/↓** or **k/j** - Move selection
- **Enter** or **Space** - Select option
- **y/n** - Confirm choices
- **q** or **Ctrl+C** - Quit

The beautiful TUI will guide you through:
1. Selecting your pomodoro split (25/5 or 50/10)
2. Displaying a real-time progress bar during work/break sessions
3. Confirming whether to take breaks
4. Tracking your completed sessions

### Quick Commands (Legacy Mode)

For simple timers without the TUI:

60-minute work session:
```bash
pomo work
```

10-minute break:
```bash
pomo rest
```

### Environment Variables

Skip the prompt by setting a default split:

```bash
# macOS/Linux
export POMO_SPLIT="25/5"
pomo

# Windows (PowerShell)
$env:POMO_SPLIT="25/5"
pomo

# Windows (CMD)
set POMO_SPLIT=25/5
pomo
```

Add to your shell config file (`~/.zshrc`, `~/.bashrc`, or PowerShell profile) to make it permanent.

## How It Works

1. **Work Session**: Timer counts down your work duration
2. **Notification**: Get a notification with sound when work is complete
3. **Break Prompt**: Choose to take a break or continue
4. **Break Session**: Timer counts down your break duration
5. **Repeat**: Option to start another pomodoro

## Building for Different Platforms

```bash
# Build for current platform
go build -o pomo

# Cross-compile for macOS
GOOS=darwin GOARCH=amd64 go build -o pomo-darwin-amd64
GOOS=darwin GOARCH=arm64 go build -o pomo-darwin-arm64

# Cross-compile for Linux
GOOS=linux GOARCH=amd64 go build -o pomo-linux-amd64
GOOS=linux GOARCH=arm64 go build -o pomo-linux-arm64

# Cross-compile for Windows
GOOS=windows GOARCH=amd64 go build -o pomo-windows-amd64.exe
```

## Requirements

**For Running:**
- Go 1.21+ (for building from source)

**Platform-specific:**
- **macOS**: Built-in (uses `osascript`)
- **Linux**: `libnotify-bin` (recommended for notifications)
- **Windows**: Windows 10+ (uses PowerShell)

**Go Dependencies:**
- [Bubble Tea](https://github.com/charmbracelet/bubbletea) - TUI framework
- [Bubbles](https://github.com/charmbracelet/bubbles) - TUI components
- [Lipgloss](https://github.com/charmbracelet/lipgloss) - Terminal styling

All Go dependencies are automatically downloaded when building.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

MIT License - see LICENSE file for details
