package okr

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"time"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/twistedogic/store"
	"github.com/twistedogic/store/bolt"
)

const dbName = "okr.db"

var (
	docStyle = lipgloss.NewStyle().Margin(1, 2)
	addKey   = key.NewBinding(
		key.WithKeys("a"),
		key.WithHelp("a", "add item"),
	)
	deleteKey = key.NewBinding(
		key.WithKeys("d"),
		key.WithHelp("d", "delete item"),
	)
	editKey = key.NewBinding(
		key.WithKeys("e"),
		key.WithHelp("edit", "edit item"),
	)
	additionKeys = func() []key.Binding {
		return []key.Binding{
			addKey,
			deleteKey,
			editKey,
		}
	}
)

type Process struct {
	list  list.Model
	input textinput.Model
	store store.StoreWithContext
}

func New(ctx context.Context) (*Process, error) {
	l := list.New([]list.Item{}, list.NewDefaultDelegate(), 0, 0)
	l.Title = "Objectives"
	l.AdditionalShortHelpKeys = additionKeys
	l.AdditionalFullHelpKeys = additionKeys
	l.StatusMessageLifetime = 5 * time.Second
	input := textinput.New()

	dir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	path := filepath.Join(dir, ".task", dbName)
	s, err := bolt.New(path)
	if err != nil {
		return nil, err
	}
	p := &Process{
		list:  l,
		input: input,
		store: store.NewStoreWithContext(ctx, s),
	}
	obj, err := p.list()
	if err != nil {
		return nil, err
	}
	l.SetItems([]list.Item(obj))
	return p, nil
}

func (p *Process) list() ([]*Objective, error) {
	items, err := p.store.PrefixScan(nil)
	if err != nil {
		return nil, err
	}
	obj := make([]*Objective, len(items))
	for i, o := range obj {
		if err := json.Unmarshal(item[i].Data, &o); err != nil {
			return nil, err
		}
		obj[i] = o
	}
	return obj, nil
}

func (p *Process) save(obj *Objective) error {
	b, err := json.Marshal(obj)
	if err != nil {
		return err
	}
	i := store.NewStringKeyItem(obj.Name, b)
	return p.store.Set(i)
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
	return p.list.View()
}

func (p *Process) Run() error {
	return tea.NewProgram(p).Start()
}
