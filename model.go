package main

import (
	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/gamut"
)

var blends = gamut.Blends(lipgloss.Color("#F25D94"), lipgloss.Color("#EDFF82"), 90)
var blendsShort = gamut.Blends(lipgloss.Color("#F25D94"), lipgloss.Color("#EDFF82"), 35)

type Feedback struct {
	TermHeight     int
	TermWidth      int
	Renderer       *lipgloss.Renderer
	Host           string
	Questions      []Question
	QuestionIndex  int
	FeedbackConfig FeedbackConfig
}

type Question struct {
	Title  string `json:"title"`
	Answer textarea.Model
}

type FeedbackConfig struct {
	WelcomeText string     `yaml:"welcome"`
	GoodbyeText string     `yaml:"goodbye"`
	Questions   []Question `yaml:"questions"`
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
		case "left", "shift+tab":
			f.QuestionIndex = Clamp(0, f.QuestionIndex-1, len(f.Questions))
		case "right", "tab":
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
		return Welcome(f)
	}
	if f.QuestionIndex >= len(f.Questions) {
		return Goodbye(f)
	}
	return f.Renderer.Place(
		f.TermWidth,
		f.TermHeight,
		lipgloss.Center,
		lipgloss.Center,
		f.Questions[f.QuestionIndex].View(f.Renderer))
}

func (q Question) View(renderer *lipgloss.Renderer) string {
	return lipgloss.JoinVertical(
		lipgloss.Center,
		renderer.NewStyle().MarginBottom(2).Render(Rainbow(renderer, q.Title, blends)),
		renderer.NewStyle().
			Border(lipgloss.RoundedBorder(), true).
			BorderForeground(lipgloss.Color("#AAAAAA")).
			MarginBottom(2).
			Render(q.Answer.View()))
}

func Welcome(f Feedback) string {
	return lipgloss.Place(
		f.TermWidth,
		f.TermHeight,
		lipgloss.Center,
		lipgloss.Center,
		lipgloss.JoinVertical(
			lipgloss.Center,
			Rainbow(f.Renderer, f.FeedbackConfig.WelcomeText, blendsShort),
			f.Renderer.NewStyle().Foreground(lipgloss.Color("#AAAAAA")).MarginTop(2).Italic(true).Render("press <tab>/<shift+tab>")))
}

func Goodbye(f Feedback) string {
	return lipgloss.Place(
		f.TermWidth,
		f.TermHeight,
		lipgloss.Center,
		lipgloss.Center,
		lipgloss.JoinVertical(
			lipgloss.Center,
			Button(f.Renderer),
			f.Renderer.NewStyle().Foreground(lipgloss.Color("#AAAAAA")).Margin(2, 2).Italic(true).Render("press <enter> to send or <shift+tab> to go back"),
			Rainbow(f.Renderer, f.FeedbackConfig.GoodbyeText, blendsShort)))
}

func Button(renderer *lipgloss.Renderer) string {
	border := renderer.NewStyle().
		Padding(1, 4).
		Border(lipgloss.RoundedBorder(), true).
		BorderForeground(lipgloss.Color("#AAAAAA"))

	return border.Render(Rainbow(renderer, "Send feedback", blendsShort))
}
