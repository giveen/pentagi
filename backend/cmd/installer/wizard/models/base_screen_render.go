package models

import (
	"strings"

	"pentagi/cmd/installer/wizard/locale"
	"pentagi/cmd/installer/wizard/styles"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	"github.com/charmbracelet/lipgloss"
)

// getViewportFormSize calculates viewport left panel dimensions
func (b *BaseScreen) getViewportFormSize() (int, int) {
	contentWidth, contentHeight := b.window.GetContentSize()
	if contentWidth <= 0 || contentHeight <= 0 {
		return 0, 0
	}

	if b.isVerticalLayout() {
		return contentWidth - PaddingWidth/2, contentHeight - PaddingHeight
	} else {
		leftWidth, rightWidth := MinMenuWidth, MinInfoWidth
		extraWidth := contentWidth - leftWidth - rightWidth - PaddingWidth
		if extraWidth > 0 {
			leftWidth = min(leftWidth+extraWidth/2, MaxMenuWidth)
		}
		return leftWidth, contentHeight - PaddingHeight
	}
}

// updateViewports updates the viewports with current content
func (b *BaseScreen) updateViewports() {
	contentWidth, contentHeight := b.window.GetContentSize()
	if contentWidth <= 0 || contentHeight <= 0 {
		return
	}

	b.updateFormContent()

	viewportWidth, viewportHeight := b.getViewportFormSize()
	formContentHeight := lipgloss.Height(b.formContent)

	b.viewportForm.Width = viewportWidth
	b.viewportForm.Height = min(viewportHeight, formContentHeight)
	b.viewportForm.SetContent(b.formContent) // force update of the viewport content

	helpContent := b.renderHelpContent()
	b.viewportHelp.Width = lipgloss.Width(helpContent)
	b.viewportHelp.Height = lipgloss.Height(helpContent)
	b.viewportHelp.SetContent(helpContent) // force update of the viewport content
}

// updateFormContent renders form content and calculates field heights
func (b *BaseScreen) updateFormContent() {
	var sections []string
	b.fieldHeights = []int{}
	inputWidth := b.GetInputWidth()

	if b.listHandler != nil {
		if listModel := b.listHandler.GetList(); listModel != nil {
			listStyle := b.styles.FormInput.Width(inputWidth)
			if b.focusedIndex == 0 {
				listStyle = listStyle.BorderForeground(styles.Primary)
			}

			listModel.SetWidth(inputWidth - 4)
			renderedList := listStyle.Render(listModel.View())

			// field title
			titleStyle := b.styles.FormLabel
			if b.getFieldIndex() == -1 {
				titleStyle = titleStyle.Foreground(styles.Primary)
			}
			title := titleStyle.Render(b.listHandler.GetListTitle())
			sections = append(sections, title)

			// field description
			description := b.styles.FormHelp.Render(b.listHandler.GetListDescription())
			sections = append(sections, description)

			sections = append(sections, renderedList)
			sections = append(sections, "")

			listHeight := lipgloss.Height(b.renderFormContent(sections[:3]))
			b.fieldHeights = append(b.fieldHeights, listHeight)
		}
	}

	for i, field := range b.fields {
		// check if this field is focused
		focused := b.getFieldIndex() == i

		// field title
		titleStyle := b.styles.FormLabel
		if focused {
			titleStyle = titleStyle.Foreground(styles.Primary)
		}
		title := titleStyle.Render(field.Title)
		sections = append(sections, title)

		// field description
		description := b.styles.FormHelp.Render(field.Description)
		sections = append(sections, description)

		// input field
		inputStyle := b.styles.FormInput.Width(inputWidth)
		if focused {
			inputStyle = inputStyle.BorderForeground(styles.Primary)
		}

		// configure input
		input := field.Input
		input.Width = inputWidth - 3
		input.SetValue(input.Value()) // force update of the input value

		// set up suggestions for tab completion
		if len(field.Suggestions) > 0 {
			input.ShowSuggestions = true
			input.SetSuggestions(field.Suggestions)
		}

		// apply masking if needed and not showing values
		if field.Masked && !b.showValues {
			input.EchoMode = textinput.EchoPassword
		} else {
			input.EchoMode = textinput.EchoNormal
		}

		// ensure focus state is correct
		if focused {
			input.Focus()
		} else {
			input.Blur()
		}

		renderedInput := inputStyle.Render(input.View())
		sections = append(sections, renderedInput)
		sections = append(sections, "")

		// update the field with configured input
		b.fields[i].Input = input

		// calculate field height
		renderedField := b.renderFormContent([]string{title, description, renderedInput})
		b.fieldHeights = append(b.fieldHeights, lipgloss.Height(renderedField))
	}

	// update list styles
	if b.listHandler != nil {
		if listModel := b.listHandler.GetList(); listModel != nil {
			listModel.Styles.PaginationStyle = b.styles.FormPagination.Width(inputWidth)
		}
		if listDelegate := b.listHandler.GetListDelegate(); listDelegate != nil {
			listDelegate.SetWidth(inputWidth)
		}
	}

	statusMessage := ""
	if b.hasChanges {
		statusMessage = b.styles.Warning.Render(locale.UIUnsavedChanges)
	} else {
		statusMessage = b.styles.Success.Render(locale.UIConfigSaved)
	}

	sections = append(sections, statusMessage)
	bottomSections := []string{statusMessage}

	if summary := b.handler.GetFormSummary(); summary != "" {
		sections = append(sections, "", summary)
		bottomSections = append(bottomSections, "", summary)
	}
	b.bottomHeight = lipgloss.Height(b.renderFormContent(bottomSections))

	// update form content
	b.formContent = b.renderFormContent(sections)
}

func (b *BaseScreen) renderFormContent(sections []string) string {
	content := strings.Join(sections, "\n")

	contentWidth, contentHeight := b.window.GetContentSize()
	viewportHeight := contentHeight - PaddingHeight // for final rendering
	approximateMaxHeight := lipgloss.Height(content)
	vp := viewport.New(contentWidth, max(viewportHeight, approximateMaxHeight*3))

	if b.isVerticalLayout() {
		xAxisPadding := verticalLayoutPaddings[1] + verticalLayoutPaddings[3]
		verticalStyle := lipgloss.NewStyle().Width(contentWidth - xAxisPadding)
		content = verticalStyle.Render(content)
	} else {
		leftWidth, rightWidth := MinMenuWidth, MinInfoWidth
		extraWidth := contentWidth - leftWidth - rightWidth - PaddingWidth
		if extraWidth > 0 {
			leftWidth = min(leftWidth+extraWidth/2, MaxMenuWidth)
		}
		xAxisPadding := horizontalLayoutPaddings[1] + horizontalLayoutPaddings[3]
		content = lipgloss.NewStyle().Width(leftWidth - xAxisPadding).Render(content)
	}

	vp.SetContent(content)
	vp.Height = vp.VisibleLineCount()

	return vp.View()
}

func (b *BaseScreen) renderHelpContent() string {
	helpContent := b.handler.GetHelpContent()
	if b.isVerticalLayout() {
		return b.renderFormContent([]string{helpContent})
	}

	contentWidth, _ := b.window.GetContentSize()
	leftWidth, rightWidth := MinMenuWidth, MinInfoWidth
	extraWidth := contentWidth - leftWidth - rightWidth - PaddingWidth
	if extraWidth > 0 {
		leftWidth = min(leftWidth+extraWidth/2, MaxMenuWidth)
		rightWidth = contentWidth - leftWidth - PaddingWidth/2
	}

	return lipgloss.NewStyle().Width(rightWidth - 2).Render(helpContent)
}

// renderForm renders the left panel with the form
func (b *BaseScreen) renderForm() string {
	if !b.initialized {
		return locale.UILoading
	}
	return b.viewportForm.View()
}

// renderHelp renders the right panel with help content
func (b *BaseScreen) renderHelp() string {
	if !b.initialized {
		return ""
	}
	return b.viewportHelp.View()
}

// isVerticalLayout determines if vertical layout should be used
func (b *BaseScreen) isVerticalLayout() bool {
	contentWidth := b.window.GetContentWidth()
	return contentWidth < (MinMenuWidth + MinInfoWidth + PaddingWidth)
}

// renderVerticalLayout renders content in vertical layout
func (b *BaseScreen) renderVerticalLayout(leftPanel, rightPanel string, width, height int) string {
	verticalStyle := lipgloss.NewStyle().Width(width).Padding(verticalLayoutPaddings...)

	leftStyled := verticalStyle.Render(leftPanel)
	rightStyled := verticalStyle.Render(rightPanel)
	if lipgloss.Height(leftStyled)+lipgloss.Height(rightStyled)+3 < height {
		return lipgloss.JoinVertical(lipgloss.Left,
			verticalStyle.Render(leftPanel),
			verticalStyle.Height(2).Render("\n"),
			verticalStyle.Render(rightPanel),
		)
	}

	return verticalStyle.Render(leftPanel)
}

// renderHorizontalLayout renders content in horizontal layout
func (b *BaseScreen) renderHorizontalLayout(leftPanel, rightPanel string, width, height int) string {
	leftWidth, rightWidth := MinMenuWidth, MinInfoWidth
	extraWidth := width - leftWidth - rightWidth - PaddingWidth
	if extraWidth > 0 {
		leftWidth = min(leftWidth+extraWidth/2, MaxMenuWidth)
		rightWidth = width - leftWidth - PaddingWidth/2
	}

	leftStyled := lipgloss.NewStyle().Width(leftWidth).Padding(horizontalLayoutPaddings...).Render(leftPanel)
	rightStyled := lipgloss.NewStyle().Width(rightWidth).PaddingLeft(2).Render(rightPanel)

	vp := viewport.New(width, height-PaddingHeight)
	vp.SetContent(lipgloss.JoinHorizontal(lipgloss.Top, leftStyled, rightStyled))

	return vp.View()
}
