package models

import (
	"pentagi/cmd/installer/wizard/controller"
	"pentagi/cmd/installer/wizard/locale"
	"pentagi/cmd/installer/wizard/styles"
	"pentagi/cmd/installer/wizard/window"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
)

var (
	verticalLayoutPaddings   = []int{0, 4, 0, 2}
	horizontalLayoutPaddings = []int{0, 2, 0, 2}
)

// BaseScreenModel defines methods that concrete screens must implement
type BaseScreenModel interface {
	// GetFormTitle returns the title for the form (layout header)
	GetFormTitle() string

	// GetFormDescription returns the description for the form (right panel)
	GetFormDescription() string

	// GetFormName returns the name for the form (right panel)
	GetFormName() string

	// GetFormOverview returns form overview for list screens (right panel)
	GetFormOverview() string

	// GetCurrentConfiguration returns text with current configuration for the list screens
	GetCurrentConfiguration() string

	// IsConfigured returns true if the form is configured
	IsConfigured() bool

	// GetFormHotKeys returns the hotkeys for the form (layout footer)
	GetFormHotKeys() []string

	tea.Model // for common interface logic
}

// BaseScreenHandler defines methods that concrete screens must implement
type BaseScreenHandler interface {
	// BuildForm builds the specific form fields for this screen
	BuildForm() tea.Cmd

	// GetFormSummary returns optional summary for the form bottom
	GetFormSummary() string

	// GetHelpContent returns the right panel help content
	GetHelpContent() string

	// HandleSave handles saving the form data
	HandleSave() error

	// HandleReset handles resetting the form to default values
	HandleReset()

	// OnFieldChanged is called when a form field value changes
	OnFieldChanged(fieldIndex int, oldValue, newValue string)

	// GetFormFields returns the current form fields
	GetFormFields() []FormField

	// SetFormFields sets the form fields
	SetFormFields(fields []FormField)
}

// BaseListHandler defines methods for screens that use lists (optional)
type BaseListHandler interface {
	// GetList returns the list model if this screen uses a list
	GetList() *list.Model

	// GetListDelegate returns the list delegate if this screen uses a list
	GetListDelegate() *BaseListDelegate

	// OnListSelectionChanged is called when list selection changes
	OnListSelectionChanged(oldSelection, newSelection string)

	// GetListTitle returns the title of the list
	GetListTitle() string

	// GetListDescription returns the description of the list
	GetListDescription() string
}

// FormField represents a single form field
type FormField struct {
	Key         string
	Title       string
	Description string
	Placeholder string
	Required    bool
	Masked      bool
	Input       textinput.Model
	Value       string
	Suggestions []string
}

// BaseScreen provides common functionality for installer form screens
type BaseScreen struct {
	// Dependencies
	controller controller.Controller
	styles     styles.Styles
	window     window.Window

	// State
	initialized  bool
	hasChanges   bool
	focusedIndex int
	showValues   bool

	// Form data
	fields       []FormField
	fieldHeights []int
	bottomHeight int

	// UI components
	viewportForm viewport.Model
	viewportHelp viewport.Model
	formContent  string

	// Handlers - must be set by concrete implementations
	handler     BaseScreenHandler
	listHandler BaseListHandler // optional, can be nil

	// Common utilities
	listHelper BaseListHelper
}

// NewBaseScreen creates a new base screen instance
func NewBaseScreen(
	c controller.Controller, s styles.Styles, w window.Window,
	h BaseScreenHandler, lh BaseListHandler, // can be nil
) *BaseScreen {
	return &BaseScreen{
		controller:   c,
		styles:       s,
		window:       w,
		showValues:   false,
		viewportForm: viewport.New(w.GetContentSize()),
		viewportHelp: viewport.New(w.GetContentSize()),
		handler:      h,
		listHandler:  lh,
		fieldHeights: []int{},
		listHelper:   BaseListHelper{},
	}
}

// Init initializes the base screen
func (b *BaseScreen) Init() tea.Cmd {
	cmd := b.handler.BuildForm()
	b.fields = b.handler.GetFormFields()
	b.updateViewports()
	return cmd
}

// Update handles common update logic and returns commands only
// Concrete implementations should call this and return themselves as the model
func (b *BaseScreen) Update(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		b.updateViewports()

	case tea.KeyMsg:
		return b.handleKeyPress(msg)
	}

	return nil
}

// GetFormHotKeys returns the hotkeys for the form
func (b *BaseScreen) GetFormHotKeys() []string {
	haveMaskedFields := false
	for _, field := range b.fields {
		if field.Masked {
			haveMaskedFields = true
			break
		}
	}

	haveFieldsWithSuggestions := false
	for _, field := range b.fields {
		if len(field.Suggestions) > 0 {
			haveFieldsWithSuggestions = true
			break
		}
	}

	hasList := b.listHandler != nil && b.listHandler.GetList() != nil

	var hotkeys []string
	if len(b.fields) > 0 || hasList {
		hotkeys = append(hotkeys, "down|up")
		if hasList {
			hotkeys = append(hotkeys, "left|right")
		}
		hotkeys = append(hotkeys, "ctrl+s")
		hotkeys = append(hotkeys, "ctrl+r")

	}
	if haveMaskedFields {
		hotkeys = append(hotkeys, "ctrl+h")
	}
	if haveFieldsWithSuggestions {
		hotkeys = append(hotkeys, "tab")
	}
	hotkeys = append(hotkeys, "enter")

	return hotkeys
}

// handleKeyPress handles common keyboard interactions
func (b *BaseScreen) handleKeyPress(msg tea.KeyMsg) tea.Cmd {
	switch msg.String() {
	case "down":
		b.focusNext()
		b.updateViewports()
		b.ensureFocusVisible()

	case "up":
		b.focusPrev()
		b.updateViewports()
		b.ensureFocusVisible()

	case "ctrl+s":
		return b.saveConfiguration()

	case "ctrl+r":
		b.resetForm()
		b.updateViewports()

	case "ctrl+h":
		b.toggleShowValues()
		b.updateViewports()

	case "tab":
		b.handleTabCompletion()

	case "enter":
		return b.saveAndReturn()
	}

	return nil
}

// View renders the screen
func (b *BaseScreen) View() string {
	contentWidth, contentHeight := b.window.GetContentSize()
	if contentWidth <= 0 || contentHeight <= 0 {
		return locale.UILoading
	}

	if !b.initialized {
		b.handler.BuildForm()
		b.fields = b.handler.GetFormFields()
		b.updateViewports()
		b.initialized = true
	}

	leftPanel := b.renderForm()
	rightPanel := b.renderHelp()

	if b.isVerticalLayout() {
		return b.renderVerticalLayout(leftPanel, rightPanel, contentWidth, contentHeight)
	}

	return b.renderHorizontalLayout(leftPanel, rightPanel, contentWidth, contentHeight)
}

// Common form methods

// GetInputWidth calculates the appropriate input width
func (b *BaseScreen) GetInputWidth() int {
	viewportWidth, _ := b.getViewportFormSize()
	inputWidth := viewportWidth - 6
	if b.isVerticalLayout() {
		inputWidth = viewportWidth - 4
	}
	return inputWidth
}

// GetController returns the state controller
func (b *BaseScreen) GetController() controller.Controller {
	return b.controller
}

// GetStyles returns the styles
func (b *BaseScreen) GetStyles() styles.Styles {
	return b.styles
}

// GetWindow returns the window
func (b *BaseScreen) GetWindow() window.Window {
	return b.window
}

// SetHasChanges sets the hasChanges flag
func (b *BaseScreen) SetHasChanges(hasChanges bool) {
	b.hasChanges = hasChanges
}

// GetHasChanges returns the hasChanges flag
func (b *BaseScreen) GetHasChanges() bool {
	return b.hasChanges
}

// GetShowValues returns the showValues flag
func (b *BaseScreen) GetShowValues() bool {
	return b.showValues
}

// GetFocusedIndex returns the currently focused field index
func (b *BaseScreen) GetFocusedIndex() int {
	return b.focusedIndex
}

// SetFocusedIndex sets the focused field index
func (b *BaseScreen) SetFocusedIndex(index int) {
	b.focusedIndex = index
}

// GetListHelper returns the list helper utility
func (b *BaseScreen) GetListHelper() *BaseListHelper {
	return &b.listHelper
}
