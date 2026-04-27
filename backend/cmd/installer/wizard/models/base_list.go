package models

import (
	"fmt"
	"io"

	"pentagi/cmd/installer/wizard/styles"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// BaseListOption represents a generic list option that can be used in any list
type BaseListOption struct {
	Value   string // the actual value
	Display string // the display text (can be different from value)
}

func (d BaseListOption) FilterValue() string { return d.Value }

// BaseListDelegate handles rendering of generic list options
type BaseListDelegate struct {
	style      lipgloss.Style
	width      int
	selectedFg lipgloss.Color
	normalFg   lipgloss.Color
}

// NewBaseListDelegate creates a new generic list delegate
func NewBaseListDelegate(style lipgloss.Style, width int) *BaseListDelegate {
	return &BaseListDelegate{
		style:      style,
		width:      width,
		selectedFg: styles.Primary,
		normalFg:   lipgloss.Color(""),
	}
}

// SetColors allows customizing the colors
func (d *BaseListDelegate) SetColors(selectedFg, normalFg lipgloss.Color) {
	d.selectedFg = selectedFg
	d.normalFg = normalFg
}

func (d *BaseListDelegate) SetWidth(width int)                     { d.width = width }
func (d BaseListDelegate) Height() int                             { return 1 }
func (d BaseListDelegate) Spacing() int                            { return 0 }
func (d BaseListDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d BaseListDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	option, ok := listItem.(BaseListOption)
	if !ok {
		return
	}

	str := option.Display
	if index == m.Index() {
		str = d.style.Width(d.width).Foreground(d.selectedFg).Render(str)
	} else {
		str = d.style.Width(d.width).Foreground(d.normalFg).Render(str)
	}

	fmt.Fprint(w, str)
}

// BaseListHelper provides utility functions for working with lists
type BaseListHelper struct{}

// CreateList creates a new list with the given options and delegate
func (h BaseListHelper) CreateList(options []BaseListOption, delegate list.ItemDelegate, width, height int) list.Model {
	items := make([]list.Item, len(options))
	for i, option := range options {
		items[i] = option
	}

	listModel := list.New(items, delegate, width, height)
	listModel.SetShowStatusBar(false)
	listModel.SetFilteringEnabled(false)
	listModel.SetShowHelp(false)
	listModel.SetShowTitle(false)

	return listModel
}

// SelectByValue selects the list item that matches the given value
func (h BaseListHelper) SelectByValue(listModel *list.Model, value string) {
	items := listModel.Items()
	for i, item := range items {
		if option, ok := item.(BaseListOption); ok && option.Value == value {
			listModel.Select(i)
			break
		}
	}
}

// GetSelectedValue returns the value of the currently selected item
func (h BaseListHelper) GetSelectedValue(listModel *list.Model) string {
	selectedItem := listModel.SelectedItem()
	if selectedItem == nil {
		return ""
	}

	if option, ok := selectedItem.(BaseListOption); ok {
		return option.Value
	}

	return ""
}

// GetSelectedDisplay returns the display text of the currently selected item
func (h BaseListHelper) GetSelectedDisplay(listModel *list.Model) string {
	selectedItem := listModel.SelectedItem()
	if selectedItem == nil {
		return ""
	}

	if option, ok := selectedItem.(BaseListOption); ok {
		return option.Display
	}

	return ""
}
