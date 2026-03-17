package items

import (
	"fmt"
	"strings"

	"github.com/camilo-zuluaga/zotero-tui/ui/cmds"
	"github.com/camilo-zuluaga/zotero-tui/zotero"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	paneBox = lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("#6B7280"))

	activePaneBox = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("#FFFFFF"))

	detailTitleStyle = lipgloss.NewStyle().Bold(true)
	sectionStyle     = lipgloss.NewStyle().
				Bold(true).
				MarginTop(1).
				MarginBottom(1)
	detailLabelStyle = lipgloss.NewStyle().Faint(true)
	helpStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("#3C3C3C"))
	detailValueStyle = lipgloss.NewStyle()
)

type pane int

const (
	paneLeft pane = iota
	paneRight
)

type Model struct {
	height, width int
	list          list.Model
	detailVP      viewport.Model
	focus         pane
	helpText      string
}

func New() Model {
	items := []list.Item{}

	l := list.New(items, NewDelegate(), 0, 0)
	l.SetStatusBarItemName("Item", "Items")
	l.SetShowTitle(false)
	l.SetShowHelp(false)
	l.SetShowPagination(false)
	l.InfiniteScrolling = true

	l.FilterInput.PromptStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#FFFFFF"))
	l.FilterInput.Cursor.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("#FFFFFF"))

	l.Styles.StatusBar = l.Styles.StatusBar.
		Align(lipgloss.Left).
		Margin(0).
		MarginBottom(1).
		Padding(0)

	return Model{
		list:  l,
		focus: paneLeft,
	}
}

func (m *Model) SetZoteroItems(zItems []zotero.ZoteroItem) {
	items := make([]list.Item, 0, len(zItems))
	for _, z := range zItems {
		var title, desc, creators string

		title = z.Data.Title
		creators = z.Data.CreatorSummary
		desc = fmt.Sprintf("%s\n%s • [DOI] %s", creators, z.Data.Date, z.Data.DOI)
		if z.Data.ShortTitle != "" {
			title = z.Data.ShortTitle
		}

		c := item{
			title:      title,
			desc:       desc,
			ZoteroItem: z,
		}
		items = append(items, c)
	}
	m.list.SetItems(items)
	m.refreshDetails()
}

func (m *Model) AppendNote(parentKey, noteKey, newContent string) {
	items := m.list.Items()
	for i, listItem := range items {
		if itm, ok := listItem.(item); ok && itm.ZoteroItem.Key == parentKey {
			note := zotero.ZoteroNote{Key: noteKey, Note: newContent}
			itm.ZoteroItem.Data.Note = append(itm.ZoteroItem.Data.Note, note)
			m.list.SetItem(i, itm)
			break
		}
	}
	m.refreshDetails()
}

func (m *Model) UpdateNote(parentKey, noteKey, newContent string) {
	items := m.list.Items()
	for i, listItem := range items {
		if itm, ok := listItem.(item); ok && itm.ZoteroItem.Key == parentKey {
			for j, n := range itm.ZoteroItem.Data.Note {
				if n.Key == noteKey {
					itm.ZoteroItem.Data.Note[j].Note = newContent
					break
				}
			}
			m.list.SetItem(i, itm)
			break
		}
	}
	m.refreshDetails()
}

func (m *Model) AppendZoteroItems(zItems []zotero.ZoteroItem) {
	existing := m.list.Items()
	for _, z := range zItems {
		var title, desc, creators, date, doi string

		title = z.Data.Title
		creators = z.Data.CreatorSummary
		date = "No Date"
		if strings.TrimSpace(z.Data.Date) != "" {
			date = z.Data.Date
		}
		doi = "No DOI"
		if z.Data.DOI != "" {
			doi = z.Data.DOI
		}
		desc = fmt.Sprintf("%s\n%s • [DOI] %s", creators, date, doi)
		if z.Data.ShortTitle != "" {
			title = z.Data.ShortTitle
		}

		c := item{
			title:      title,
			desc:       desc,
			ZoteroItem: z,
		}
		existing = append(existing, c)
	}
	m.list.SetItems(existing)
	m.refreshDetails()
}

func (m *Model) ClearItems() {
	m.list.SetItems([]list.Item{})
}

func (m *Model) HelpText(msgType HelpMsgType) tea.Cmd {
	var helpText string
	switch msgType {
	case ModeNormal:
		helpText = "[tab] Switch pane • [enter] Load details • [n] Create/Edit note • [r] Read PDF • [b] Bibliography • [esc] Back • [q] quit"

	case ModeClipboard:
		m.helpText = helpStyle.Render("Copied to clipboard!")
		return cmds.ResetHelpCmd()
	}
	m.helpText = helpStyle.Render(helpText)
	return nil
}

func (m *Model) UpdateChildrenItems(parentKey string, attachments []zotero.ZoteroAttachment,
	notes []zotero.ZoteroNote) {
	items := m.list.Items()
	for i, listItem := range items {
		if itm, ok := listItem.(item); ok && itm.ZoteroItem.Key == parentKey {
			itm.ZoteroItem.Data.Attachment = attachments
			itm.ZoteroItem.Data.Note = notes
			m.list.SetItem(i, itm)
			break
		}
	}
	m.refreshDetails()
}

func (m Model) SelectedZoteroItem() *zotero.ZoteroItem {
	it, ok := m.list.SelectedItem().(item)
	if !ok {
		return nil
	}
	return &it.ZoteroItem
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "tab":
			if m.focus == paneLeft {
				m.focus = paneRight
			} else {
				m.focus = paneLeft
			}
			return m, nil
		}

		if m.focus == paneRight {
			switch msg.String() {
			case "up", "k":
				m.detailVP.ScrollUp(1)
			case "down", "j":
				m.detailVP.ScrollDown(1)
			case "pgup":
				m.detailVP.HalfPageUp()
			case "pgdown":
				m.detailVP.HalfPageDown()
			}
			return m, nil
		}

		var cmd tea.Cmd
		m.list, cmd = m.list.Update(msg)
		m.refreshDetails()
		return m, cmd
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m *Model) refreshDetails() {
	contentWidth := max(m.detailVP.Width, 10)
	m.detailVP.SetContent(m.buildDetailsContent(contentWidth))
	m.detailVP.GotoTop()
}

func (m Model) buildDetailsContent(contentWidth int) string {
	zi := m.SelectedZoteroItem()
	if zi == nil {
		return detailLabelStyle.Render("No attachment selected")
	}

	d := zi.Data
	title := d.Title

	titleLine := detailTitleStyle.
		Width(contentWidth).
		Render(title)

	lines := []string{
		titleLine,
		"",
	}

	if d.ItemType != "" {
		lines = append(lines, fmt.Sprintf(
			"%s %s",
			detailLabelStyle.Render("ItemType:"),
			d.ItemType,
		))
	}

	lines = append(lines, "")

	if len(d.Creators) != 0 {
		for _, c := range d.Creators {
			lines = append(lines, fmt.Sprintf(
				"%s %s",
				detailLabelStyle.Render("Author:"),
				fmt.Sprintf("%s, %s", c.LastName, c.FirstName),
			))
		}
	}

	lines = append(lines, "")

	if d.Date != "" {
		lines = append(lines, fmt.Sprintf(
			"%s %s",
			detailLabelStyle.Render("Date:"),
			d.Date,
		))
	}

	if d.DOI != "" {
		lines = append(lines, fmt.Sprintf(
			"%s %s",
			detailLabelStyle.Render("DOI:"),
			d.DOI,
		))
	}

	if d.URL != "" {
		lines = append(lines, fmt.Sprintf(
			"%s %s",
			detailLabelStyle.Render("URL:"),
			d.URL,
		))
	}

	lines = append(lines, sectionStyle.Render("󰁦 Attachments"))

	if len(d.Attachment) != 0 {
		for _, a := range d.Attachment {
			lines = append(lines, fmt.Sprintf(
				"%s %s",
				detailLabelStyle.Render("Attachment:"),
				a.Title,
			))
		}
	} else {
		lines = append(lines, detailLabelStyle.Faint(true).Render("Empty Attachments"))
	}

	lines = append(lines, sectionStyle.Render("󰎛 Notes"))

	if len(d.Note) != 0 {
		formatNote := func(note string) string {
			noteHTMLstripped := StripHTML(note)
			if len(noteHTMLstripped) > 200 {
				return noteHTMLstripped[:200] + "..."
			}
			return noteHTMLstripped
		}

		for _, n := range d.Note {
			lines = append(lines, fmt.Sprintf(
				"%s %s",
				detailLabelStyle.Render("Note:"),
				detailValueStyle.Width(contentWidth).Render(
					formatNote(n.Note)+"\n"),
			))
		}
	} else {
		lines = append(lines, detailLabelStyle.Faint(true).Render("Empty notes"))
	}

	lines = append(lines, "")

	if d.DateModified != "" {
		lines = append(lines, fmt.Sprintf(
			"%s %s",
			detailLabelStyle.Render("Modified:"),
			d.DateModified,
		))
	}

	return lipgloss.JoinVertical(lipgloss.Left, lines...)
}

func (m Model) View() string {
	totalInner := m.width - 4
	leftWidth := int(float64(totalInner) * 0.6)
	rightWidth := totalInner - leftWidth

	if leftWidth < 1 {
		leftWidth = 1
	}
	if rightWidth < 1 {
		rightWidth = 1
	}

	innerHeight := m.height - 3

	/* ── LEFT PANE ── */
	leftContent := lipgloss.NewStyle().
		PaddingLeft(1).
		Width(leftWidth - 2).
		Height(innerHeight - 2).
		MaxHeight(innerHeight - 2).
		Render(m.list.View())

	leftBox := paneBox.
		Width(leftWidth).
		Height(innerHeight).
		Render(leftContent)

	/* ── RIGHT PANE ── */
	scrollInfo := detailLabelStyle.Render(
		fmt.Sprintf(" %d%%", int(m.detailVP.ScrollPercent()*100)),
	)

	rightInner := lipgloss.NewStyle().
		PaddingLeft(1).
		PaddingTop(1).
		Width(rightWidth - 2).
		Height(innerHeight - 3).
		MaxHeight(innerHeight - 3).
		Render(m.detailVP.View())

	rightInner = lipgloss.JoinVertical(lipgloss.Left, rightInner, scrollInfo)

	rightBoxStyle := paneBox
	if m.focus == paneRight {
		rightBoxStyle = activePaneBox
	}
	rightBox := rightBoxStyle.
		Width(rightWidth).
		Height(innerHeight).
		Render(rightInner)

	panes := lipgloss.JoinHorizontal(lipgloss.Top, leftBox, rightBox)

	return lipgloss.JoinVertical(lipgloss.Left, panes, m.helpText)
}

func (m *Model) SetSize(width, height int) {
	m.width = width
	m.height = height

	totalInner := width - 4
	leftWidth := int(float64(totalInner) * 0.6)
	rightWidth := totalInner - leftWidth

	listWidth := leftWidth - 2
	listHeight := height - 11

	if listWidth < 1 {
		listWidth = 1
	}
	if listHeight < 1 {
		listHeight = 1
	}
	m.list.SetSize(listWidth, listHeight)

	// Viewport
	vpWidth := rightWidth - 4
	vpHeight := height - 5
	if vpWidth < 1 {
		vpWidth = 1
	}
	if vpHeight < 1 {
		vpHeight = 1
	}
	m.detailVP.Width = vpWidth
	m.detailVP.Height = vpHeight

	m.refreshDetails()
}
