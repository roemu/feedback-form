package main

// An example Bubble Tea server. This will put an ssh session into alt screen
// and continually print up to date terminal information.

import (
	"context"
	"database/sql"
	_ "embed"
	"errors"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/log"
	"github.com/charmbracelet/ssh"
	"github.com/charmbracelet/wish"
	"github.com/charmbracelet/wish/activeterm"
	"github.com/charmbracelet/wish/bubbletea"
	"github.com/charmbracelet/wish/logging"
	"gopkg.in/yaml.v3"
)

//go:embed questions.yaml
var questions []byte

var database *sql.DB

const (
	host = "0.0.0.0"
	port = "23234"
)

func main() {
	database = CreateDatabase()
	defer database.Close()
	
	s, err := wish.NewServer(
		wish.WithAddress(net.JoinHostPort(host, port)),
		wish.WithHostKeyPath(".ssh/id_ed25519"),
		wish.WithMiddleware(
			bubbletea.Middleware(teaHandler),
			activeterm.Middleware(), // Bubble Tea apps usually require a PTY.
			logging.Middleware(),
		),
	)
	if err != nil {
		log.Error("Could not start server", "error", err)
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	log.Info("Starting SSH server", "host", host, "port", port)
	go func() {
		if err = s.ListenAndServe(); err != nil && !errors.Is(err, ssh.ErrServerClosed) {
			log.Error("Could not start server", "error", err)
			done <- nil
		}
	}()

	<-done
	log.Info("Stopping SSH server")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer func() { cancel() }()
	if err := s.Shutdown(ctx); err != nil && !errors.Is(err, ssh.ErrServerClosed) {
		log.Error("Could not stop server", "error", err)
	}
}

func teaHandler(s ssh.Session) (tea.Model, []tea.ProgramOption) {
	pty, _, _ := s.Pty()

	var questionConfig FeedbackConfig
	err := yaml.Unmarshal(questions, &questionConfig)
	if err != nil {
		log.Fatal("Unable to Unmarshal questions yaml: ", "err", err)
	}
	if questionConfig.WelcomeText == "" {
		log.Fatal("WelcomeText must not be empyt. Set welcome: <welcomeText> in question.yaml")
	}
	if questionConfig.GoodbyeText == "" {
		log.Fatal("GoodbyeText must not be empyt. Set goodbye: <goodbyeText> in question.yaml")
	}

	questionConfig.Questions = Map(questionConfig.Questions, func(q Question) Question {
		q.Answer = textarea.New()
		q.Answer.SetWidth(50)
		q.Answer.SetHeight(10)
		q.Answer.ShowLineNumbers = false;
		return q
	})

	f := Feedback{
		pty.Window.Width,
		pty.Window.Height,
		s.User(),
		questionConfig.Questions,
		-1,
		questionConfig,
	}
	return f, []tea.ProgramOption{tea.WithAltScreen()}
}
