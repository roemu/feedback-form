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
	cmds := []tea.Cmd{}
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		f.TermHeight = msg.Height
		f.TermWidth = msg.Width
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			if f.QuestionIndex >= len(f.Questions) {
				// TODO: Save to file
				return f, tea.Quit
			}
		case "ctrl+c":
			return f, tea.Quit
		case "ctrl+d":
			f.DebugMode = !f.DebugMode
			return f, nil
		case "up", "shift+tab":
			// TODO: factor out to below switch and upate with reassigning the model back into main model
			f.QuestionIndex = Clamp(0, f.QuestionIndex-1, len(f.Questions))
			if f.QuestionIndex > len(f.Questions) {
				cmd := f.Questions[f.QuestionIndex].Answer.Focus()
				cmds = append(cmds, cmd)
			}
			return f, nil
		case "down", "tab":
			f.QuestionIndex = Clamp(0, f.QuestionIndex+1, len(f.Questions))
			if f.QuestionIndex > len(f.Questions) {
				cmd := f.Questions[f.QuestionIndex].Answer.Focus()
				cmds = append(cmds, cmd)
			}
			return f, nil
		}
	}
	return f, tea.Batch(cmds...)
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
		f.Questions[f.QuestionIndex].Answer.View(),
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
