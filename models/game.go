package models

import (
	"fmt"
	"time"

	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/muesli/termenv"
)

const (
	width = 60.

	// charsPerWord is the average characters per word used by most typing tests
	// to calculate your WPM score.
	charsPerWord = 5.
)

// The game model is used for storing the bubbletea application's current state
// It contains mutliple progress bars, their percentages and is drawn on the
// screen
// All clients send their percentages to the central server.
type Game struct {
	// Percentages is a value from 0 to 1 that represents the current completion of the typing test
	Percentages  []float64 `json:"percentages"`
	Progress []*progress.Model `json:"-"`
	// this player's playerID
	PlayerID int `json:"player_id"`
	// Text is the randomly generated text for the user to type
	Text string `json:"-"`
	// Typed is the text that the user has typed so far
	Typed string `json:"-"`
	// Start and end are the start and end time of the typing test
	Start time.Time `json:"-"`
	// Mistakes is the number of characters that were mistyped by the user
	Mistakes int `json:"-"`
	// Score is the user's score calculated by correct characters typed
	Score float64 `json:"-"`
}

// Init inits the bubbletea model for use
func (m Game) Init() tea.Cmd {
	return nil
}

// updateProgress updates the percentage for this player
func (m Game) updateProgress() (tea.Model, tea.Cmd) {
	m.Percentages[m.PlayerID] = float64(len(m.Typed)) / float64(len(m.Text))
	if m.Percentages[m.PlayerID] >= 1.0 {
		return m, tea.Quit
	}
	return m, nil
}

// Update updates the bubbletea model by handling the progress bar update
// and adding typed characters to the state if they are valid typing characters
func (m Game) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		// Start counting time only after the first keystroke
		if m.Start.IsZero() {
			m.Start = time.Now()
		}

		// User wants to cancel the typing test
		if msg.Type == tea.KeyCtrlC {
			return m, tea.Quit
		}

		// Deleting characters
		if msg.Type == tea.KeyBackspace && len(m.Typed) > 0 {
			m.Typed = m.Typed[:len(m.Typed)-1]
			return m.updateProgress()
		}

		// Ensure we are adding characters only that we want the user to be able to type
		if msg.Type != tea.KeyRunes {
			return m, nil
		}

		char := msg.Runes[0]
		next := rune(m.Text[len(m.Typed)])

		// To properly account for line wrapping we need to always insert a new line
		// Where the next line starts to not break the user interface, even if the user types a random character
		if next == '\n' {
			m.Typed += "\n"

			// Since we need to perform a line break
			// if the user types a space we should simply ignore it.
			if char == ' ' {
				return m, nil
			}
		}

		m.Typed += msg.String()

		if char == next {
			m.Score += 1.
		}

		return m.updateProgress()
	//case tea.WindowSizeMsg:
	//	m.Progress.Width = msg.Width - 4
	//	if m.Progress.Width > width {
	//		m.Progress.Width = width
	//	}
	//	return m, nil

	default:
		return m, nil
	}
}

// View shows the current state of the typing test.
// It displays a progress bar for the progression of the typing test,
// the typed characters (with errors displayed in red) and remaining
// characters to be typed in a faint display
func (m Game) View() string {
	remaining := m.Text[len(m.Typed):]

	var typed string
	for i, c := range m.Typed {
		if c == rune(m.Text[i]) {
			typed += string(c)
		} else {
			typed += termenv.String(string(m.Typed[i])).Background(termenv.ANSIBrightRed).String()
		}
	}

	str := "\n "
	for i, p := range m.Percentages {
		str = fmt.Sprintf("%s\n\n%s", str, m.Progress[i].View(p))
	}

	s := fmt.Sprintf("%s\n\n%s%s", str, typed, termenv.String(remaining).Faint())

	var wpm float64
	// Start counting wpm after at least two characters are typed
	if len(m.Typed) > 1 {
		wpm = (m.Score / charsPerWord) / (time.Since(m.Start).Minutes())
	}
	s += fmt.Sprintf("\n\nWPM: %.2f\n", wpm)
	return s
}
