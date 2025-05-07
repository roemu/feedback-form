package main

import (
	"strconv"

	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/gamut"
)

var blends = gamut.Blends(lipgloss.Color("#F25D94"), lipgloss.Color("#EDFF82"), 25)

type Feedback struct {
	TermHeight     int
	TermWidth      int
	Host           string
	Questions      []Question
	QuestionIndex  int
	QuestionConfig QuestionConfig
	DebugMode      bool
}

type Question struct {
	ID     int    `json:"id"` // to identify in output file
	Title  string `json:"title"`
	Answer textarea.Model
}

func (f Feedback) Init() tea.Cmd {
	return nil
}

func (f Feedback) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		f.TermHeight = msg.Height
		f.TermWidth = msg.Width
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return f, tea.Quit
		case "ctrl+d":
			f.DebugMode = !f.DebugMode
			return f, nil
		case "up", "shift+tab":
			f.QuestionIndex = Clamp(0, f.QuestionIndex-1, len(f.Questions))
		case "down", "tab":
			f.QuestionIndex = Clamp(0, f.QuestionIndex+1, len(f.Questions))
		case "enter":
			if f.QuestionIndex >= len(f.Questions) {
				
				InsertFeedback(database, f)
				return f, tea.Quit
			}
			fallthrough
		default:
			if -1 < f.QuestionIndex && f.QuestionIndex < len(f.Questions) {
				updated, cmd := f.Questions[f.QuestionIndex].Answer.Update(msg)
				f.Questions[f.QuestionIndex].Answer = updated
				return f, cmd
			}
		}
	}
	if -1 < f.QuestionIndex && f.QuestionIndex < len(f.Questions) {
		cmd := f.Questions[f.QuestionIndex].Answer.Focus()
		return f, cmd
	}
	return f, nil
}

func (f Feedback) View() string {
	if f.QuestionIndex < 0 {
		return lipgloss.Place(
			f.TermWidth,
			f.TermHeight,
			lipgloss.Center,
			lipgloss.Center,
			Rainbow(lipgloss.NewStyle(), f.QuestionConfig.WelcomeText, blends),
		)
	}
	if f.QuestionIndex >= len(f.Questions) {
		return lipgloss.Place(f.TermWidth,
			f.TermHeight,
			lipgloss.Center,
			lipgloss.Center,
			lipgloss.JoinVertical(
				lipgloss.Center,
				Button(),
				lipgloss.NewStyle().MarginBottom(2).Italic(true).Render("<press enter>"),
				f.QuestionConfig.GoodbyeText,
			),
		)
	}
	information := lipgloss.JoinVertical(
		lipgloss.Center,
		lipgloss.NewStyle().MarginBottom(2).Render(f.Questions[f.QuestionIndex].Title),
		lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder(), true).
			BorderForeground(lipgloss.Color("#FFFFFF")).
			Render(
				f.Questions[f.QuestionIndex].Answer.View()),
	)
	if f.DebugMode {
		information = lipgloss.JoinVertical(
			lipgloss.Center,
			f.Host,
			strconv.Itoa(f.TermWidth),
			strconv.Itoa(f.TermHeight),
			strconv.Itoa(f.QuestionIndex),
		)
	}
	return lipgloss.Place(f.TermWidth, f.TermHeight, lipgloss.Center, lipgloss.Center, information)
}

type QuestionConfig struct {
	WelcomeText string     `yaml:"welcome"`
	GoodbyeText string     `yaml:"goodbye"`
	Questions   []Question `yaml:"questions"`
}

func Button() string {
	border := lipgloss.NewStyle().
		Padding(1, 4).
		Border(lipgloss.RoundedBorder(), true).
		BorderForeground(lipgloss.Color("#874BFD"))

	return border.Render(Rainbow(lipgloss.NewStyle(), "Send feedback", blends))
}
