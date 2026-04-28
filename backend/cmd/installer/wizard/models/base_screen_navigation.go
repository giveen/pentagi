package models

import (
	"slices"

	"pentagi/cmd/installer/wizard/logger"

	tea "github.com/charmbracelet/bubbletea"
)

// ensureFocusVisible scrolls the viewport to ensure focused field is visible
func (b *BaseScreen) ensureFocusVisible() {
	if b.focusedIndex >= len(b.fieldHeights) {
		return
	}

	// calculate y position of focused field
	focusY := 0
	if b.focusedIndex == len(b.fieldHeights)-1 {
		focusY = b.bottomHeight
	}
	for i := range b.focusedIndex {
		focusY += b.fieldHeights[i] + 1 // empty line between fields
	}

	// get viewport dimensions
	visibleRows := b.viewportForm.Height
	offset := b.viewportForm.YOffset

	// if focused field is above visible area, scroll up
	if focusY < offset {
		b.viewportForm.YOffset = focusY
	}

	// if focused field is below visible area, scroll down
	if focusY+b.fieldHeights[b.focusedIndex] >= offset+visibleRows {
		b.viewportForm.YOffset = focusY + b.fieldHeights[b.focusedIndex] - visibleRows + 1
	}
}

// focusNext moves focus to the next field
func (b *BaseScreen) focusNext() {
	totalElements := b.getTotalElements()
	if totalElements == 0 {
		return
	}

	// blur current field
	b.blurCurrentField()

	// move to next element (with wrapping)
	b.focusedIndex = (b.focusedIndex + 1) % totalElements

	// focus new field
	b.focusCurrentField()
	b.updateFormContent()
}

// focusPrev moves focus to the previous field
func (b *BaseScreen) focusPrev() {
	totalElements := b.getTotalElements()
	if totalElements == 0 {
		return
	}

	// blur current field
	b.blurCurrentField()

	// move to previous element (with wrapping)
	b.focusedIndex = (b.focusedIndex - 1 + totalElements) % totalElements

	// focus new field
	b.focusCurrentField()
	b.updateFormContent()
}

// getTotalElements returns the total number of navigable elements
func (b *BaseScreen) getTotalElements() int {
	total := len(b.fields)

	// add 1 for list if present
	if b.listHandler != nil && b.listHandler.GetList() != nil {
		total++
	}

	return total
}

// blurCurrentField removes focus from the currently focused field
func (b *BaseScreen) blurCurrentField() {
	fieldIndex := b.getFieldIndex()
	if fieldIndex >= 0 && fieldIndex < len(b.fields) {
		b.fields[fieldIndex].Input.Blur()
	}
}

// focusCurrentField sets focus on the currently focused field
func (b *BaseScreen) focusCurrentField() {
	fieldIndex := b.getFieldIndex()
	if fieldIndex >= 0 && fieldIndex < len(b.fields) {
		b.fields[fieldIndex].Input.Focus()
	}
}

// getFieldIndex returns the field index for the current focusedIndex (-1 if focused on list)
func (b *BaseScreen) getFieldIndex() int {
	if b.listHandler != nil && b.listHandler.GetList() != nil {
		// list is at index 0, fields start at index 1
		return b.focusedIndex - 1
	}
	// no list, fields start at index 0
	return b.focusedIndex
}

// toggleShowValues toggles visibility of masked values
func (b *BaseScreen) toggleShowValues() {
	b.showValues = !b.showValues
	b.updateFormContent()
}

// handleTabCompletion handles tab completion for focused field
func (b *BaseScreen) handleTabCompletion() {
	fieldIndex := b.getFieldIndex()

	// check if we're focused on a valid field
	if fieldIndex >= 0 && fieldIndex < len(b.fields) {
		field := &b.fields[fieldIndex]

		// only handle tab completion if field has suggestions
		if len(field.Suggestions) > 0 {
			// use textinput's built-in suggestion functionality
			if suggestion := field.Input.CurrentSuggestion(); suggestion != "" {
				oldValue := field.Input.Value()
				field.Input.SetValue(suggestion)
				field.Input.CursorEnd()
				field.Value = suggestion
				b.hasChanges = true

				// notify handler about the change
				b.handler.OnFieldChanged(fieldIndex, oldValue, suggestion)

				// update the fields array
				b.fields[fieldIndex] = *field
				b.handler.SetFormFields(b.fields)
				b.updateViewports()
			}
		}
	}
}

// resetForm resets the form to default values
func (b *BaseScreen) resetForm() {
	b.handler.HandleReset()
	b.fields = b.handler.GetFormFields()
	b.hasChanges = false
	b.updateFormContent()
}

// saveConfiguration saves the current configuration
func (b *BaseScreen) saveConfiguration() tea.Cmd {
	if err := b.handler.HandleSave(); err != nil {
		logger.Errorf("[BaseScreen] SAVE: error: %v", err)
		return nil
	}

	b.hasChanges = false
	b.updateViewports()

	return nil
}

// saveAndReturn saves and returns to previous screen
func (b *BaseScreen) saveAndReturn() tea.Cmd {
	// save first
	cmd := b.saveConfiguration()
	if cmd != nil {
		return cmd
	}

	// return to previous screen
	return func() tea.Msg {
		return NavigationMsg{GoBack: true}
	}
}

// HandleFieldInput handles input for a specific field
func (b *BaseScreen) HandleFieldInput(msg tea.KeyMsg) tea.Cmd {
	fieldIndex := b.getFieldIndex()

	// all hotkeys are handled by handleKeyPress(msg tea.KeyMsg) method
	// inherit screen must call HandleUpdate(msg tea.Msg) for all uncaught messages
	if slices.Contains(b.GetFormHotKeys(), msg.String()) {
		return nil
	}

	// check if we're focused on a valid field
	if fieldIndex >= 0 && fieldIndex < len(b.fields) {
		var cmd tea.Cmd
		oldValue := b.fields[fieldIndex].Input.Value()
		b.fields[fieldIndex].Input, cmd = b.fields[fieldIndex].Input.Update(msg)
		newValue := b.fields[fieldIndex].Input.Value()

		if oldValue != newValue {
			b.fields[fieldIndex].Value = newValue
			b.hasChanges = true
			b.handler.OnFieldChanged(fieldIndex, oldValue, newValue)
		}

		b.updateViewports()
		return cmd
	}

	return nil
}

// HandleListInput handles input for the list component
func (b *BaseScreen) HandleListInput(msg tea.KeyMsg) tea.Cmd {
	// check if we have a list and we're focused on it (skip if not)
	if b.listHandler == nil {
		return nil
	}

	// check if focused on list (index 0 when list is present)
	isFocusedOnList := b.listHandler.GetList() != nil && b.focusedIndex == 0
	if !isFocusedOnList {
		return nil
	}

	// filter list input keys to slide the list
	switch msg.String() {
	case "left", "right":
		break
	default:
		return nil
	}

	listModel := b.listHandler.GetList()
	if listModel == nil {
		return nil
	}

	// get old selection
	oldSelection := ""
	if selectedItem := listModel.SelectedItem(); selectedItem != nil {
		oldSelection = selectedItem.FilterValue()
	}

	// update list
	var cmd tea.Cmd
	*listModel, cmd = listModel.Update(msg)

	// get new selection
	newSelection := ""
	if selectedItem := listModel.SelectedItem(); selectedItem != nil {
		newSelection = selectedItem.FilterValue()
	}

	// notify handler if selection changed
	if oldSelection != newSelection {
		b.listHandler.OnListSelectionChanged(oldSelection, newSelection)
		b.hasChanges = true
		b.updateViewports()
	}

	return cmd
}
