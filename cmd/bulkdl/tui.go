package main

import (
	"fmt"
	"time"

	"github.com/anasrar/bulkdl/downloader"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var TUIDownloadFilterStartIndex = 0

type (
	TUIMsgUpdateFrame struct{}
	TUIMsgFinish      struct{}
)

type TUIModel struct {
	Frame  int
	Finish bool
}

func TUICmdUpdateFrame() tea.Cmd {
	return tea.Tick(time.Second/10, func(time.Time) tea.Msg {
		return TUIMsgUpdateFrame{}
	})
}

func (m TUIModel) Init() tea.Cmd {
	return TUICmdUpdateFrame()
}

func (m TUIModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case TUIMsgUpdateFrame:
		m.Frame++
		return m, TUICmdUpdateFrame()

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		}
	case TUIMsgFinish:
		if m.Finish {
			return m, tea.Quit
		}
		m.Finish = true
		return m, func() tea.Msg { return TUIMsgFinish{} }
	}
	return m, nil
}

func (m TUIModel) View() string {
	if m.Finish {
		return ""
	}

	var (
		result       = TUIStyleHelp.Render("ctrl+c/q: quit") + "\n\n" + fmt.Sprintf("Downloads — %s (%s/%s)", TUIStyleFgYellow.Render(fmt.Sprintf("%5.1f%%", (*D.ItemsFinish/D.TOTAL_ITEMS_FLOAT)*100)), TUIStyleFgGreen.Render(fmt.Sprintf("%*.0f", D.FINISH_STRING_WIDTH, *D.ItemsFinish)), TUIStyleFgCyan.Render(fmt.Sprintf("%d", D.TOTAL_ITEMS_INT))) + "\n" + fmt.Sprintf("%s Error(s)", TUIStyleFgRed.Render(fmt.Sprintf("%0.f", *D.ItemsError))) + "\n\n"
		spinner      = TUIStyleFgBlue.Bold(true).Render(TUIRenderSpinner3(m.Frame))
		downloads    = []downloader.DownloaderItem{}
		runningTotal = 0
	)

	for index := TUIDownloadFilterStartIndex; index < D.TOTAL_ITEMS_INT; index++ {
		item := D.Items[index]
		if item.Status == downloader.DownloaderItemStatusRunning {
			// NOTE: Sometime skip 1-3 queue when it start because goroutine not start in a sequence
			if runningTotal == 0 && TUIDownloadFilterStartIndex < index {
				TUIDownloadFilterStartIndex = index
			}

			downloads = append(downloads, item)
			runningTotal++
			if runningTotal == D.Max {
				break
			}
		}
	}

	for _, item := range downloads {
		result += fmt.Sprintf("%s %s %s\n", spinner, TUIStyleFgLight.Render(fmt.Sprintf("%6.2f%%", item.Progress)), TUIStyleFgBrightLight.Render(item.Filename))
	}

	return TUIStyleContainer.Render(result)
}

// # Render

var (
	TUI_SPINNER3        = []string{"⠉⠁ ", "⠈⠉⠁", " ⠈⠉", "  ⠙", "  ⠸", "  ⠴", " ⠠⠤", "⠠⠤ ", "⠦  ", "⠇  ", "⠋  "}
	TUI_SPINNER3_LENGTH = len(TUI_SPINNER3)
)

func TUIRenderSpinner3(n int) string {
	return TUI_SPINNER3[n%TUI_SPINNER3_LENGTH]
}

// # Style

var (
	TUIStyleHelp      = lipgloss.NewStyle().Foreground(lipgloss.Color("8"))
	TUIStyleContainer = lipgloss.NewStyle().Margin(1, 2)
)

var (
	TUIStyleFgRed         = lipgloss.NewStyle().Foreground(lipgloss.Color("1"))
	TUIStyleFgGreen       = lipgloss.NewStyle().Foreground(lipgloss.Color("2"))
	TUIStyleFgYellow      = lipgloss.NewStyle().Foreground(lipgloss.Color("3"))
	TUIStyleFgBlue        = lipgloss.NewStyle().Foreground(lipgloss.Color("4"))
	TUIStyleFgCyan        = lipgloss.NewStyle().Foreground(lipgloss.Color("6"))
	TUIStyleFgLight       = lipgloss.NewStyle().Foreground(lipgloss.Color("7"))
	TUIStyleFgBrightLight = lipgloss.NewStyle().Foreground(lipgloss.Color("15"))
)
