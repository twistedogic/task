package search

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type mode uint

const (
	modeList  mode = iota // list delegate
	modeView              // vim delegate
	modeInput             // textinput delegate
)

var (
	docStyle  = lipgloss.NewStyle().Margin(1, 2)
	searchKey = key.NewBinding(
		key.WithKeys("w"),
		key.WithHelp("w", "search"),
	)
	newSearchKey = key.NewBinding(
		key.WithKeys("s"),
		key.WithHelp("s", "new search"),
	)
	additionKeys = func() []key.Binding { return []key.Binding{searchKey, newSearchKey} }
)

type Process struct {
	ctx       context.Context
	source    Source
	mode      mode
	list      list.Model
	textInput textinput.Model
	errCh     chan error
}

func New(ctx context.Context, src Source) *Process {
	resultList := list.New([]list.Item{}, list.NewDefaultDelegate(), 0, 0)
	resultList.Title = "Search results"
	resultList.AdditionalShortHelpKeys = additionKeys
	resultList.AdditionalFullHelpKeys = additionKeys
	resultList.StatusMessageLifetime = 5 * time.Second
	textInput := textinput.New()
	textInput.Placeholder = "r:github.com/prometheus f:ast.go Expr"
	return &Process{
		ctx:       ctx,
		source:    src,
		textInput: textInput,
		list:      resultList,
		errCh:     make(chan error, 1),
	}
}

func (p *Process) startVim(r Result) error {
	dir, err := ioutil.TempDir("", "search")
	if err != nil {
		return err
	}
	filename := filepath.Base(r.File)
	path := filepath.Join(dir, filename)
	p.mode = modeView
	defer func() {
		p.mode = modeList
		os.RemoveAll(dir)
	}()
	if err := p.source.Download(p.ctx, &r); err != nil {
		return err
	}
	if err := ioutil.WriteFile(path, r.Content, 0644); err != nil {
		return err
	}
	to := fmt.Sprintf("+%d", r.Matches[0].LineNumber+1)
	cmd := exec.Command("vim", to, path)
	cmd.Stdout, cmd.Stdin = os.Stdout, os.Stdin
	return cmd.Run()
}

func (p *Process) inputUpdate(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			p.search(p.textInput.Value())
			p.mode = modeList
			return nil
		case "esc":
			p.textInput.Reset()
			p.mode = modeList
			return nil
		}
	}
	var cmd tea.Cmd
	p.textInput, cmd = p.textInput.Update(msg)
	return cmd
}

func (p *Process) currentItem() Result {
	r, ok := p.list.SelectedItem().(Result)
	if !ok {
		return Result{}
	}
	return r
}

func (p *Process) setTitle(s string) { p.list.Title = s }

func (p *Process) listUpdate(msg tea.Msg) tea.Cmd {
	select {
	case err := <-p.errCh:
		return p.list.NewStatusMessage(err.Error())
	default:
	}
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if !p.list.SettingFilter() {
			switch {
			case key.Matches(msg, newSearchKey):
				p.mode = modeInput
				p.textInput.Reset()
				p.textInput.Focus()
				return nil
			case key.Matches(msg, searchKey):
				p.mode = modeInput
				q := p.currentItem().QueryString()
				p.textInput.SetValue(q)
				p.textInput.SetCursor(len(q))
				p.textInput.Focus()
				return nil
			}
		}
		switch msg.String() {
		case "enter":
			if !p.list.SettingFilter() {
				p.startVim(p.currentItem())
				return nil
			}
		}
	}
	var cmd tea.Cmd
	p.list, cmd = p.list.Update(msg)
	return cmd
}

func (p *Process) Init() tea.Cmd { return nil }

func (p *Process) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return p, tea.Quit
		}
	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		p.list.SetSize(msg.Width-h, msg.Height-v)
	}
	var cmd tea.Cmd
	switch p.mode {
	case modeList:
		cmd = p.listUpdate(msg)
	case modeInput:
		cmd = p.inputUpdate(msg)
	}
	return p, cmd
}

func (p *Process) View() string {
	switch p.mode {
	case modeList:
		return p.list.View()
	case modeInput:
		return fmt.Sprintf("query:\n\n%s", p.textInput.View())
	}
	return p.list.View()
}

func (p *Process) search(query ...string) {
	q, err := ParseQuery(query...)
	if err != nil {
		p.errCh <- err
		return
	}
	res := make(Results, 0)
	r, err := p.source.Search(p.ctx, q)
	if err != nil {
		p.errCh <- err
		return
	}
	res = append(res, r...)
	p.setTitle(fmt.Sprintf("Search results: %v", query))
	p.list.SetItems(res.Items())
}

func (p *Process) Run(query ...string) error {
	if len(query) != 0 {
		p.search(query...)
	}
	return tea.NewProgram(p).Start()
}
