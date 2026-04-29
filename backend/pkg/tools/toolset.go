package tools

import "github.com/vxcontrol/langchaingo/llms"

// toolSetBuilder accumulates tool definitions and handlers for a customExecutor.
type toolSetBuilder struct {
	fte        *flowToolsExecutor
	taskID     *int64
	subtaskID  *int64
	defs       []llms.FunctionDefinition
	handlers   map[string]ExecutorHandler
	barriers   map[string]struct{}
	summarizer SummarizeHandler
}

// newToolSet creates a toolSetBuilder for the given flow executor and task context.
func (fte *flowToolsExecutor) newToolSet(taskID, subtaskID *int64, summarizer SummarizeHandler) *toolSetBuilder {
	return &toolSetBuilder{
		fte:        fte,
		taskID:     taskID,
		subtaskID:  subtaskID,
		handlers:   make(map[string]ExecutorHandler),
		barriers:   make(map[string]struct{}),
		summarizer: summarizer,
	}
}

// add registers a mandatory (always-present) tool.
func (tsb *toolSetBuilder) add(name string, handler ExecutorHandler) *toolSetBuilder {
	tsb.defs = append(tsb.defs, registryDefinitions[name])
	tsb.handlers[name] = handler
	return tsb
}

// addBarrier registers a tool that, when called, terminates the execution chain.
func (tsb *toolSetBuilder) addBarrier(name string, handler ExecutorHandler) *toolSetBuilder {
	tsb.add(name, handler)
	tsb.barriers[name] = struct{}{}
	return tsb
}

// addOptional appends each named tool only if t.IsAvailable() returns true.
// All names share the same handler (t.Handle).
func (tsb *toolSetBuilder) addOptional(t Tool, names ...string) *toolSetBuilder {
	if !t.IsAvailable() {
		return tsb
	}
	for _, name := range names {
		tsb.defs = append(tsb.defs, registryDefinitions[name])
		tsb.handlers[name] = t.Handle
	}
	return tsb
}

// addTerminalFrom appends TerminalToolName and FileToolName, both handled by term.Handle.
func (tsb *toolSetBuilder) addTerminalFrom(term Tool) *toolSetBuilder {
	tsb.defs = append(tsb.defs,
		registryDefinitions[TerminalToolName],
		registryDefinitions[FileToolName],
	)
	tsb.handlers[TerminalToolName] = term.Handle
	tsb.handlers[FileToolName] = term.Handle
	return tsb
}

// addHttpClient always appends the HTTP client tool (no availability guard).
func (tsb *toolSetBuilder) addHttpClient() *toolSetBuilder {
	httpclient := NewHttpClientTool(tsb.fte.flowID, tsb.taskID, tsb.subtaskID)
	tsb.defs = append(tsb.defs, registryDefinitions[HttpClientToolName])
	tsb.handlers[HttpClientToolName] = httpclient.Handle
	return tsb
}

// addBrowser optionally appends the browser tool.
func (tsb *toolSetBuilder) addBrowser() *toolSetBuilder {
	browser := NewBrowserTool(
		tsb.fte.flowID, tsb.taskID, tsb.subtaskID,
		tsb.fte.cfg.DataDir,
		tsb.fte.cfg.ScraperPrivateURL,
		tsb.fte.cfg.ScraperPublicURL,
		tsb.fte.scp,
	)
	return tsb.addOptional(browser, BrowserToolName)
}

// addGuide optionally appends both guide store and guide search tools (Store then Search order).
func (tsb *toolSetBuilder) addGuide() *toolSetBuilder {
	guide := NewGuideTool(
		tsb.fte.flowID, tsb.taskID, tsb.subtaskID,
		tsb.fte.replacer, tsb.fte.store, tsb.fte.vslp,
	)
	return tsb.addOptional(guide, StoreGuideToolName, SearchGuideToolName)
}

// addGuideSearch optionally appends the guide search tool only (no store).
func (tsb *toolSetBuilder) addGuideSearch() *toolSetBuilder {
	guide := NewGuideTool(
		tsb.fte.flowID, tsb.taskID, tsb.subtaskID,
		tsb.fte.replacer, tsb.fte.store, tsb.fte.vslp,
	)
	return tsb.addOptional(guide, SearchGuideToolName)
}

// addMemory optionally appends the in-memory similarity search tool.
func (tsb *toolSetBuilder) addMemory() *toolSetBuilder {
	memory := NewMemoryTool(tsb.fte.flowID, tsb.fte.store, tsb.fte.vslp)
	return tsb.addOptional(memory, SearchInMemoryToolName)
}

// addAnswers optionally appends both the search-answer and store-answer tools.
func (tsb *toolSetBuilder) addAnswers() *toolSetBuilder {
	s := NewSearchTool(
		tsb.fte.flowID, tsb.taskID, tsb.subtaskID,
		tsb.fte.replacer, tsb.fte.store, tsb.fte.vslp,
	)
	return tsb.addOptional(s, SearchAnswerToolName, StoreAnswerToolName)
}

// addAnswerSearch optionally appends the search-answer tool only (no store).
func (tsb *toolSetBuilder) addAnswerSearch() *toolSetBuilder {
	s := NewSearchTool(
		tsb.fte.flowID, tsb.taskID, tsb.subtaskID,
		tsb.fte.replacer, tsb.fte.store, tsb.fte.vslp,
	)
	return tsb.addOptional(s, SearchAnswerToolName)
}

// addCode optionally appends both the code search and store tools.
func (tsb *toolSetBuilder) addCode() *toolSetBuilder {
	code := NewCodeTool(
		tsb.fte.flowID, tsb.taskID, tsb.subtaskID,
		tsb.fte.replacer, tsb.fte.store, tsb.fte.vslp,
	)
	return tsb.addOptional(code, SearchCodeToolName, StoreCodeToolName)
}

// addCodeSearch optionally appends the code search tool only (no store).
func (tsb *toolSetBuilder) addCodeSearch() *toolSetBuilder {
	code := NewCodeTool(
		tsb.fte.flowID, tsb.taskID, tsb.subtaskID,
		tsb.fte.replacer, tsb.fte.store, tsb.fte.vslp,
	)
	return tsb.addOptional(code, SearchCodeToolName)
}

// addGraphiti optionally appends the graphiti knowledge-graph search tool.
func (tsb *toolSetBuilder) addGraphiti() *toolSetBuilder {
	gs := NewGraphitiSearchTool(
		tsb.fte.flowID, tsb.taskID, tsb.subtaskID,
		tsb.fte.graphitiClient,
	)
	return tsb.addOptional(gs, GraphitiSearchToolName)
}

// addSploitus optionally appends the Sploitus exploit-search tool.
func (tsb *toolSetBuilder) addSploitus() *toolSetBuilder {
	sploitus := NewSploitusTool(
		tsb.fte.cfg, tsb.fte.flowID, tsb.taskID, tsb.subtaskID,
		tsb.fte.slp,
	)
	return tsb.addOptional(sploitus, SploitusToolName)
}

// addWebSearch optionally appends all configured web-search engines in order:
// Searxng, Google, DuckDuckGo, Tavily, Traversaal, Perplexity.
func (tsb *toolSetBuilder) addWebSearch() *toolSetBuilder {
	searxng := NewSearxngTool(tsb.fte.cfg, tsb.fte.flowID, tsb.taskID, tsb.subtaskID, tsb.fte.slp, tsb.summarizer)
	tsb.addOptional(searxng, SearxngToolName)

	google := NewGoogleTool(tsb.fte.cfg, tsb.fte.flowID, tsb.taskID, tsb.subtaskID, tsb.fte.slp)
	tsb.addOptional(google, GoogleToolName)

	duckduckgo := NewDuckDuckGoTool(tsb.fte.cfg, tsb.fte.flowID, tsb.taskID, tsb.subtaskID, tsb.fte.slp)
	tsb.addOptional(duckduckgo, DuckDuckGoToolName)

	tavily := NewTavilyTool(tsb.fte.cfg, tsb.fte.flowID, tsb.taskID, tsb.subtaskID, tsb.fte.slp, tsb.summarizer)
	tsb.addOptional(tavily, TavilyToolName)

	traversaal := NewTraversaalTool(tsb.fte.cfg, tsb.fte.flowID, tsb.taskID, tsb.subtaskID, tsb.fte.slp)
	tsb.addOptional(traversaal, TraversaalToolName)

	perplexity := NewPerplexityTool(tsb.fte.cfg, tsb.fte.flowID, tsb.taskID, tsb.subtaskID, tsb.fte.slp, tsb.summarizer)
	tsb.addOptional(perplexity, PerplexityToolName)

	return tsb
}

// build constructs and returns a customExecutor from the accumulated tool set.
func (tsb *toolSetBuilder) build() *customExecutor {
	return &customExecutor{
		flowID:      tsb.fte.flowID,
		taskID:      tsb.taskID,
		subtaskID:   tsb.subtaskID,
		mlp:         tsb.fte.mlp,
		vslp:        tsb.fte.vslp,
		db:          tsb.fte.db,
		store:       tsb.fte.store,
		definitions: tsb.defs,
		handlers:    tsb.handlers,
		barriers:    tsb.barriers,
		summarizer:  tsb.summarizer,
	}
}
