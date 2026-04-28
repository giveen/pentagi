package locale

// Core locale constants shared across screens.
// Includes UI labels, common statuses, and navigation hints.

// Section: Common Status And UI
const (
	// Common status and UI strings
	UIStatistics       = "Statistics"
	UIStatus           = "Status: "
	UIMode             = "Mode: "
	UINoConfigSelected = "No configuration selected"
	UILoading          = "Loading..."
	UINotImplemented   = "Not implemented yet"
	UIUnsavedChanges   = "Unsaved changes"
	UIConfigSaved      = "Configuration saved"

	// Status labels
	StatusEnabled       = "Enabled"
	StatusDisabled      = "Disabled"
	StatusConfigured    = "Configured"
	StatusNotConfigured = "Not configured"
	StatusEmbedded      = "Embedded"
	StatusExternal      = "External"

	// Success/Warning messages
	MessageSearchEnginesNone       = "⚠ No search engines configured"
	MessageSearchEnginesConfigured = "✓ %d search engines configured"
	MessageDockerConfigured        = "✓ Docker environment configured"
	MessageDockerNotConfigured     = "⚠ Docker environment not configured"
)

// Section: Legend
const (
	LegendConfigured    = "✓ Configured"
	LegendNotConfigured = "✗ Not configured"
)

// Section: Common Navigation Actions
const (
	NavBack       = "Esc: Back"
	NavExit       = "Ctrl+Q: Exit"
	NavUpDown     = "↑/↓: Scroll/Select"
	NavLeftRight  = "←/→: Move"
	NavPgUpPgDown = "PgUp/PgDn: Page"
	NavHomeEnd    = "Home/End: Start/End"
	NavEnter      = "Enter: Continue"
	NavYn         = "Y/N: Accept/Reject"
	NavCtrlC      = "Ctrl+C: Cancel"
	NavCtrlS      = "Ctrl+S: Save"
	NavCtrlR      = "Ctrl+R: Reset"
	NavCtrlH      = "Ctrl+H: Show/Hide"
	NavTab        = "Tab: Complete"
	NavSeparator  = " • "
)

