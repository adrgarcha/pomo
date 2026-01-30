package main

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Styles
var (
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FF6B6B")).
			MarginBottom(1)

	subtitleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#95E1D3"))

	helpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#AAAAAA")).
			Italic(true)

	selectedStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF6B6B")).
			Bold(true).
			PaddingLeft(2)

	normalStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFFFF")).
			PaddingLeft(4)

	timeStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#4ECDC4"))

	completedStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#95E1D3")).
			Background(lipgloss.Color("#2D4A3E")).
			Padding(1, 2)
)

// Session types
type sessionType int

const (
	sessionWork sessionType = iota
	sessionBreak
)

// States
type state int

const (
	stateMenu state = iota
	stateTimer
	stateConfirm
	stateComplete
)

// Messages
type tickMsg time.Time
type completeMsg struct{}

// Model
type model struct {
	state          state
	sessionType    sessionType
	cursor         int
	choices        []string
	workDuration   time.Duration
	breakDuration  time.Duration
	remaining      time.Duration
	total          time.Duration
	progress       progress.Model
	spinner        spinner.Model
	quitting       bool
	sessionCount   int
	confirmBreak   bool
	confirmAnother bool
}

func initialModel() model {
	p := progress.New(
		progress.WithDefaultGradient(),
		progress.WithWidth(40),
		progress.WithoutPercentage(),
	)

	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("#FF6B6B"))

	return model{
		state:    stateMenu,
		cursor:   0,
		choices:  []string{"25/5 (25m work, 5m break)", "50/10 (50m work, 10m break)", "Exit"},
		progress: p,
		spinner:  s,
		quitting: false,
	}
}

func (m model) Init() tea.Cmd {
	return m.spinner.Tick
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch m.state {
		case stateMenu:
			return m.updateMenu(msg)
		case stateTimer:
			return m.updateTimer(msg)
		case stateConfirm:
			return m.updateConfirm(msg)
		case stateComplete:
			return m.updateComplete(msg)
		}

	case tickMsg:
		if m.state == stateTimer {
			m.remaining -= time.Second
			if m.remaining <= 0 {
				return m, m.completeTimer()
			}
			return m, m.tick()
		}

	case completeMsg:
		if m.sessionType == sessionWork {
			// Work session complete
			m.sessionCount++
			sendNotification("🍅 Pomodoro Complete", "Time to take a well-deserved break! 🧘", "Crystal")
			m.state = stateConfirm
			m.confirmBreak = true
		} else {
			// Break session complete
			sendNotification("⏰ Back to Work", "Break time is over. Ready to focus again? 💪", "Crystal")
			m.state = stateConfirm
			m.confirmAnother = true
		}
		return m, nil

	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd

	case tea.WindowSizeMsg:
		m.progress.Width = msg.Width - 4
		m.progress.Width = min(m.progress.Width, 80)
		return m, nil
	}

	return m, nil
}

func (m model) updateMenu(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c", "q":
		m.quitting = true
		return m, tea.Quit

	case "up", "k":
		if m.cursor > 0 {
			m.cursor--
		}

	case "down", "j":
		if m.cursor < len(m.choices)-1 {
			m.cursor++
		}

	case "enter", " ":
		switch m.cursor {
		case 0: // 25/5
			m.workDuration = 25 * time.Minute
			m.breakDuration = 5 * time.Minute
			return m.startWorkSession(), m.tick()
		case 1: // 50/10
			m.workDuration = 50 * time.Minute
			m.breakDuration = 10 * time.Minute
			return m.startWorkSession(), m.tick()
		case 2: // Exit
			m.quitting = true
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m model) updateTimer(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c", "q":
		m.quitting = true
		return m, tea.Quit
	}
	return m, nil
}

func (m model) updateConfirm(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c", "q":
		m.quitting = true
		return m, tea.Quit

	case "y", "Y":
		if m.confirmBreak {
			// Start break
			m.confirmBreak = false
			return m.startBreakSession(), m.tick()
		}
		if m.confirmAnother {
			// Start another pomodoro
			m.confirmAnother = false
			return m.startWorkSession(), m.tick()
		}

	case "n", "N":
		if m.confirmBreak {
			// Skip break, start another work session
			m.confirmBreak = false
			return m.startWorkSession(), m.tick()
		}
		if m.confirmAnother {
			// Done with pomodoros
			m.confirmAnother = false
			m.state = stateComplete
		}
	}
	return m, nil
}

func (m model) updateComplete(_ tea.KeyMsg) (tea.Model, tea.Cmd) {
	m.quitting = true
	return m, tea.Quit
}

func (m model) startWorkSession() tea.Model {
	m.state = stateTimer
	m.sessionType = sessionWork
	m.remaining = m.workDuration
	m.total = m.workDuration
	return m
}

func (m model) startBreakSession() tea.Model {
	m.state = stateTimer
	m.sessionType = sessionBreak
	m.remaining = m.breakDuration
	m.total = m.breakDuration
	return m
}

func (m model) tick() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func (m model) completeTimer() tea.Cmd {
	return func() tea.Msg {
		return completeMsg{}
	}
}

func (m model) View() string {
	if m.quitting {
		return ""
	}

	switch m.state {
	case stateMenu:
		return m.viewMenu()
	case stateTimer:
		return m.viewTimer()
	case stateConfirm:
		return m.viewConfirm()
	case stateComplete:
		return m.viewComplete()
	}
	return ""
}

func (m model) viewMenu() string {
	var b strings.Builder

	b.WriteString(titleStyle.Render("🍅 Pomodoro Timer") + "\n\n")
	b.WriteString(subtitleStyle.Render("Choose your session:") + "\n\n")

	for i, choice := range m.choices {
		cursor := " "
		if m.cursor == i {
			cursor = "›"
			b.WriteString(selectedStyle.Render(cursor+" "+choice) + "\n")
		} else {
			b.WriteString(normalStyle.Render(cursor+" "+choice) + "\n")
		}
	}

	b.WriteString("\n" + helpStyle.Render("↑/↓: navigate • enter: select • q: quit"))

	if m.sessionCount > 0 {
		b.WriteString("\n\n" + completedStyle.Render(fmt.Sprintf("✓ Completed sessions: %d", m.sessionCount)))
	}

	return b.String()
}

func (m model) viewTimer() string {
	var b strings.Builder

	var emoji, sessionName string
	if m.sessionType == sessionWork {
		emoji = "🍅"
		sessionName = "Work Session"
	} else {
		emoji = "☕"
		sessionName = "Break Time"
	}

	b.WriteString(titleStyle.Render(emoji+" "+sessionName) + "\n\n")

	// Time remaining
	mins := int(m.remaining.Minutes())
	secs := int(m.remaining.Seconds()) % 60
	timeStr := fmt.Sprintf("%02d:%02d", mins, secs)
	b.WriteString(timeStyle.Render(timeStr) + "\n\n")

	// Progress bar
	percent := 1.0 - (float64(m.remaining) / float64(m.total))
	b.WriteString(m.progress.ViewAs(percent) + "\n\n")

	// Stats
	totalMins := int(m.total.Minutes())
	b.WriteString(subtitleStyle.Render(fmt.Sprintf("Session: %d minutes", totalMins)) + "\n")

	if m.sessionCount > 0 {
		b.WriteString(subtitleStyle.Render(fmt.Sprintf("Completed: %d sessions", m.sessionCount)) + "\n")
	}

	b.WriteString("\n" + helpStyle.Render("q: quit"))

	return b.String()
}

func (m model) viewConfirm() string {
	var b strings.Builder

	if m.confirmBreak {
		b.WriteString(titleStyle.Render("🎉 Work session complete!") + "\n\n")
		b.WriteString(subtitleStyle.Render("Ready for a break?") + "\n\n")
		b.WriteString(normalStyle.Render("  y - Take a break") + "\n")
		b.WriteString(normalStyle.Render("  n - Continue working") + "\n")
	} else if m.confirmAnother {
		b.WriteString(titleStyle.Render("✨ Break complete!") + "\n\n")
		b.WriteString(subtitleStyle.Render("Start another pomodoro?") + "\n\n")
		b.WriteString(normalStyle.Render("  y - Start another session") + "\n")
		b.WriteString(normalStyle.Render("  n - All done") + "\n")
	}

	b.WriteString("\n" + helpStyle.Render("y/n: choose • q: quit"))

	return b.String()
}

func (m model) viewComplete() string {
	var b strings.Builder

	b.WriteString(completedStyle.Render("🎉 Great work today!") + "\n\n")
	b.WriteString(titleStyle.Render(fmt.Sprintf("You completed %d pomodoro session(s)!", m.sessionCount)) + "\n\n")
	b.WriteString(subtitleStyle.Render("Press any key to exit..."))

	return b.String()
}

// Notification functions
func sendNotification(title, message, sound string) error {
	switch runtime.GOOS {
	case "darwin":
		return sendNotificationMacOS(title, message, sound)
	case "linux":
		return sendNotificationLinux(title, message)
	case "windows":
		return sendNotificationWindows(title, message)
	default:
		return nil
	}
}

func sendNotificationMacOS(title, message, sound string) error {
	script := fmt.Sprintf(`display notification "%s" with title "%s" sound name "%s"`, message, title, sound)
	cmd := exec.Command("osascript", "-e", script)
	return cmd.Run()
}

func sendNotificationLinux(title, message string) error {
	cmd := exec.Command("notify-send", title, message, "-u", "normal")
	err := cmd.Run()
	if err != nil {
		fmt.Printf("\n🔔 %s: %s\n", title, message)
	}
	return nil
}

func sendNotificationWindows(title, message string) error {
	script := fmt.Sprintf(`
		[Windows.UI.Notifications.ToastNotificationManager, Windows.UI.Notifications, ContentType = WindowsRuntime] | Out-Null
		[Windows.Data.Xml.Dom.XmlDocument, Windows.Data.Xml.Dom.XmlDocument, ContentType = WindowsRuntime] | Out-Null

		$template = @"
		<toast>
			<visual>
				<binding template="ToastText02">
					<text id="1">%s</text>
					<text id="2">%s</text>
				</binding>
			</visual>
		</toast>
"@

		$xml = New-Object Windows.Data.Xml.Dom.XmlDocument
		$xml.LoadXml($template)
		$toast = New-Object Windows.UI.Notifications.ToastNotification $xml
		[Windows.UI.Notifications.ToastNotificationManager]::CreateToastNotifier("Pomo").Show($toast)
	`, title, message)

	cmd := exec.Command("powershell", "-Command", script)
	err := cmd.Run()
	if err != nil {
		fmt.Printf("\n🔔 %s: %s\n", title, message)
	}
	return nil
}

// Simple command-line modes for backward compatibility
func workMode() {
	fmt.Println("🍅 Starting 60-minute work session...")
	time.Sleep(60 * time.Minute)
	sendNotification("🍅 Pomodoro Complete", "Time to take a well-deserved break! 🧘", "Crystal")
	fmt.Println("\n✓ Work session complete!")
}

func restMode() {
	fmt.Println("☕ Starting 10-minute break...")
	time.Sleep(10 * time.Minute)
	sendNotification("⏰ Back to Work", "Break time is over. Ready to focus again? 💪", "Crystal")
	fmt.Println("\n✓ Break complete!")
}

func printUsage() {
	fmt.Println(titleStyle.Render("🍅 Pomo - Pomodoro Timer"))
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  pomo          - Interactive TUI mode (default)")
	fmt.Println("  pomo work     - Quick 60-minute work timer")
	fmt.Println("  pomo rest     - Quick 10-minute break timer")
	fmt.Println("  pomo --help   - Show this help message")
	fmt.Println()
	fmt.Println("Environment Variables:")
	fmt.Println("  POMO_SPLIT    - Set to '25/5' or '50/10' to skip the prompt")
}

func main() {
	// Handle command-line arguments
	if len(os.Args) >= 2 {
		switch os.Args[1] {
		case "work":
			workMode()
			return
		case "rest":
			restMode()
			return
		case "--help", "-h", "help":
			printUsage()
			return
		default:
			fmt.Printf("Unknown command: %s\n\n", os.Args[1])
			printUsage()
			os.Exit(1)
		}
	}

	// Run interactive TUI mode
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}
