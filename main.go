package main

import (
	"fmt"
	"log"
	"strings"
	"sync"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/ssh"
	"github.com/charmbracelet/wish"
	bm "github.com/charmbracelet/wish/bubbletea"
)

// Styles for the UI
var (
	titleStyle  = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#7d56f4")).Padding(0, 1)
	systemStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#FFD700")).Italic(true)
	userStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("#00ADD8")).Bold(true)
)

type Server struct {
	mu    sync.Mutex
	conns map[ssh.Session]chan string
}

func main() {
	s := &Server{
		conns: make(map[ssh.Session]chan string),
	}

	srv, err := wish.NewServer(
		wish.WithAddress("0.0.0.0:2222"),
		wish.WithHostKeyPath("id_ed25519"), // charmbracelet/wish generates this key on first time starting server using github.com/charmbracelet/keygen
		wish.WithMiddleware(
			// 1. Use the Bubble Tea middleware
			bm.Middleware(s.teaHandler),
		),
	)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Chat server starting on :2222")
	if err := srv.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}

// 2. The Tea Model: This represents the UI state for ONE user
type model struct {
	server   *Server
	sess     ssh.Session
	username string
	messages []string
	input    string
	sub      chan string // Local channel to receive broadcasts
}

// 3. ReceiveMsg: A custom message type for the Bubble Tea loop
type receiveMsg string

func (s *Server) teaHandler(sess ssh.Session) (tea.Model, []tea.ProgramOption) {
	user := sess.User()
	msgChan := make(chan string, 10)

	s.mu.Lock()
	s.conns[sess] = msgChan
	s.mu.Unlock()

	// Notify others
	s.broadcast(fmt.Sprintf("*** %s joined the chat ***", user), sess)

	m := model{
		server:   s,
		sess:     sess,
		username: user,
		messages: []string{fmt.Sprintf("Welcome %s!", user)},
		sub:      msgChan,
	}

	return m, []tea.ProgramOption{tea.WithAltScreen()}
}

func (m model) Init() tea.Cmd {
	// Start listening for messages immediately
	return m.waitForMsg()
}

// waitForMsg is a command that waits for a message on the channel
func (m model) waitForMsg() tea.Cmd {
	return func() tea.Msg {
		return receiveMsg(<-m.sub)
	}
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			m.cleanup()
			return m, tea.Quit

		case "enter":
			if m.input == "/quit" {
				m.cleanup()
				return m, tea.Quit
			}
			if m.input == "/who" {
				m.server.mu.Lock()
				var users []string
				for sess := range m.server.conns {
					users = append(users, sess.User())
				}
				m.server.mu.Unlock()
				m.messages = append(m.messages, systemStyle.Render("Users: "+strings.Join(users, ", ")))
				m.input = ""
				return m, nil
			}

			if strings.TrimSpace(m.input) != "" {
				m.server.broadcast(fmt.Sprintf("[%s]: %s", m.username, m.input), m.sess)
				m.messages = append(m.messages, fmt.Sprintf("You: %s", m.input))
				m.input = ""
			}
			return m, nil

		case "backspace":
			if len(m.input) > 0 {
				m.input = m.input[:len(m.input)-1]
			}
			return m, nil

		default:
			m.input += msg.String()
			return m, nil
		}

	case receiveMsg:
		// When a broadcast arrives, add it to our list and wait for the next one
		m.messages = append(m.messages, string(msg))
		return m, m.waitForMsg()
	}

	return m, nil
}

func (m model) View() string {
	var s strings.Builder
	s.WriteString(titleStyle.Render("SSH Chat Server") + "\n\n")

	// Display last 15 messages
	displayMsgs := m.messages
	if len(displayMsgs) > 15 {
		displayMsgs = displayMsgs[len(displayMsgs)-15:]
	}

	for _, msg := range displayMsgs {
		s.WriteString(msg + "\n")
	}

	s.WriteString(fmt.Sprintf("\n> %s█", m.input))
	return s.String()
}

func (m model) cleanup() {
	m.server.mu.Lock()
	delete(m.server.conns, m.sess)
	m.server.mu.Unlock()
	m.server.broadcast(fmt.Sprintf("*** %s left the chat ***", m.username), nil)
}

func (s *Server) broadcast(msg string, sender ssh.Session) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for sess, ch := range s.conns {
		if sender == nil || sess != sender {
			select {
			case ch <- msg:
			default:
			}
		}
	}
}
