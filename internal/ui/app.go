package ui

import (
	"fmt"
	"os"
	"sort"
	"time"
	"image/color"

	"gioui.org/app"
	"gioui.org/font/gofont"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"

	"github.com/diskcern/diskcern/internal/models"
	"github.com/diskcern/diskcern/internal/providers"
	"github.com/diskcern/diskcern/internal/scanner"
)

type UIState struct {
	IsScanning   bool
	LogStream    []string
	InsightCards []InsightCard
	GameCards    []InsightCard
	CustomCards  []InsightCard

	WizTreePathEditor widget.Editor
	ImportWizTreeBtn  widget.Clickable
	BrowseWizTreeBtn  widget.Clickable
	ScanFolderBtn     widget.Clickable

	TotalDiskSpace  int64
	UsedDiskSpace   int64
	CurrentScanPath string
	
	RootNode     *TreeNode
	VisibleNodes []*TreeNode
	SelectedNode *TreeNode
	ShowDialog   bool
	CurrentTab   int

	Registry *providers.Registry
	Window   *app.Window

	InsightsList widget.List
	TreeList     widget.List
	CustomState  CustomState

	CloseDialogBtn    widget.Clickable
	IncBtn            widget.Clickable
	DecBtn            widget.Clickable
	RuleBtn           widget.Clickable
	SaveDescBtn       widget.Clickable
	DescriptionEditor widget.Editor

	TabDashBtn    widget.Clickable
	TabGamesBtn   widget.Clickable
	TabIgnoredBtn widget.Clickable
	TabTermBtn    widget.Clickable
}

type InsightCard struct {
	ProviderName string
	Risk         string
	Recoverable  int64
	Context      string
	CommandStr   string
	ActionBtn    widget.Clickable
	CardBtn      widget.Clickable
	Path         string
}

func RunApp() {
	go func() {
		w := new(app.Window)
		w.Option(app.Title("Diskcern Dashboard"), app.Size(1024, 768))
		if err := draw(w); err != nil {
			fmt.Fprintf(os.Stderr, "failed to run UI: %v\n", err)
			os.Exit(1)
		}
		os.Exit(0)
	}()
	app.Main()
}

func draw(w *app.Window) error {
	th := material.NewTheme()
	th.Shaper = text.NewShaper(text.WithCollection(gofont.Collection()))

	reg := providers.NewRegistry()
	reg.Register(&providers.NodeProvider{})
	reg.Register(&providers.SteamProvider{})
	reg.Register(&providers.DockerProvider{})
	reg.Register(&providers.EpicGamesProvider{})
	reg.Register(&providers.SDKProvider{})
	reg.Register(&providers.AppDataProvider{})
	reg.Register(&providers.DockerProvider{})

	state := &UIState{
		IsScanning:     false,
		Registry:       reg,
		Window:         w,
		UsedDiskSpace:  75,
		TotalDiskSpace: 100,
		CustomState:    LoadState(),
	}
	state.WizTreePathEditor.SetText("C:\\temp\\wiztree.csv")
	state.InsightsList.Axis = layout.Vertical
	state.TreeList.Axis = layout.Vertical
	state.DescriptionEditor.SingleLine = false

	var ops op.Ops
	for {
		switch e := w.Event().(type) {
		case app.DestroyEvent:
			return e.Err
		case app.FrameEvent:
			gtx := app.NewContext(&ops, e)

			if state.BrowseWizTreeBtn.Clicked(gtx) {
				go func() {
					path, err := OpenFileDialog("CSV Files (*.csv)|*.csv")
					if err == nil && path != "" {
						state.WizTreePathEditor.SetText(path)
						w.Invalidate()
					}
				}()
			}

			if state.ScanFolderBtn.Clicked(gtx) {
				go func() {
					path, err := OpenFolderDialog()
					if err == nil && path != "" {
						handleFolderScan(state, path)
					}
				}()
			}

			if state.TabDashBtn.Clicked(gtx) { state.CurrentTab = 0 }
			if state.TabGamesBtn.Clicked(gtx) { state.CurrentTab = 1 }
			if state.TabIgnoredBtn.Clicked(gtx) { state.CurrentTab = 2 }
			if state.TabTermBtn.Clicked(gtx) { state.CurrentTab = 3 }

			if state.ImportWizTreeBtn.Clicked(gtx) {
				path := state.WizTreePathEditor.Text()
				handleWizTreeImport(state, path)
			}
			
			if state.CloseDialogBtn.Clicked(gtx) {
				state.ShowDialog = false
			}
			
			if state.ShowDialog && state.SelectedNode != nil {
				if state.IncBtn.Clicked(gtx) {
					state.CustomState.Scores[state.SelectedNode.Path]++
					SaveState(state.CustomState)
				}
				if state.DecBtn.Clicked(gtx) {
					state.CustomState.Scores[state.SelectedNode.Path]--
					SaveState(state.CustomState)
				}
				if state.SaveDescBtn.Clicked(gtx) {
					state.CustomState.Descriptions[state.SelectedNode.Path] = state.DescriptionEditor.Text()
					SaveState(state.CustomState)
					state.LogStream = append(state.LogStream, "[STATE] Saved custom description for: " + state.SelectedNode.Path)
				}
				if state.RuleBtn.Clicked(gtx) {
					AddJSONRule(state.SelectedNode.Path)
					state.LogStream = append(state.LogStream, "[RULE] Bound JSON rule to: " + state.SelectedNode.Path)
					state.ShowDialog = false
				}
			}

			if !state.ShowDialog {
				for _, node := range state.VisibleNodes {
					if node.Clickable.Clicked(gtx) {
						now := time.Now()
						if now.Sub(node.LastClick) < 300*time.Millisecond {
							// Double Click
							state.SelectedNode = node
							state.DescriptionEditor.SetText(state.CustomState.Descriptions[node.Path])
							state.ShowDialog = true
						} else {
							// Single Click
							if node.IsDir {
								node.Expanded = !node.Expanded
								state.VisibleNodes = FlattenTree(state.RootNode, state.CustomState.Scores)
							}
						}
						node.LastClick = now
					}
				}
			}

			// Handle Insights clicks
			for i := range state.InsightCards {
				if state.InsightCards[i].ActionBtn.Clicked(gtx) {
					cmd := state.InsightCards[i].CommandStr
					if cmd != "" {
						state.LogStream = append(state.LogStream, "[EXEC] "+cmd)
						ch := make(chan string)
						go func() { ExecuteCommand(cmd, ch) }()
						go func() {
							for msg := range ch {
								state.LogStream = append(state.LogStream, msg)
								w.Invalidate()
							}
						}()
					}
				}
				if state.InsightCards[i].CardBtn.Clicked(gtx) {
					// Dummy node for dialog
					state.SelectedNode = &TreeNode{Path: state.InsightCards[i].Path}
					state.DescriptionEditor.SetText(state.CustomState.Descriptions[state.InsightCards[i].Path])
					state.ShowDialog = true
				}
			}

			// Main UI Layout
			m := op.Record(gtx.Ops)
			
			layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return DrawTopBar(gtx, th, state)
				}),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return DrawDiskOccupationBar(gtx, th, state.UsedDiskSpace, state.TotalDiskSpace)
				}),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					// Tabs
					return layout.UniformInset(unit.Dp(8)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
						return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
							layout.Rigid(func(gtx layout.Context) layout.Dimensions {
								btn := material.Button(th, &state.TabDashBtn, "Dashboard")
								if state.CurrentTab == 0 { btn.Background = color.NRGBA{R: 50, G: 100, B: 200, A: 255} }
								return btn.Layout(gtx)
							}),
							layout.Rigid(layout.Spacer{Width: unit.Dp(8)}.Layout),
							layout.Rigid(func(gtx layout.Context) layout.Dimensions {
								btn := material.Button(th, &state.TabGamesBtn, "Games")
								if state.CurrentTab == 1 { btn.Background = color.NRGBA{R: 50, G: 100, B: 200, A: 255} }
								return btn.Layout(gtx)
							}),
							layout.Rigid(layout.Spacer{Width: unit.Dp(8)}.Layout),
							layout.Rigid(func(gtx layout.Context) layout.Dimensions {
								btn := material.Button(th, &state.TabIgnoredBtn, "Custom / Ignored")
								if state.CurrentTab == 2 { btn.Background = color.NRGBA{R: 50, G: 100, B: 200, A: 255} }
								return btn.Layout(gtx)
							}),
							layout.Rigid(layout.Spacer{Width: unit.Dp(8)}.Layout),
							layout.Rigid(func(gtx layout.Context) layout.Dimensions {
								btn := material.Button(th, &state.TabTermBtn, "Execution Terminal")
								if state.CurrentTab == 3 { btn.Background = color.NRGBA{R: 50, G: 100, B: 200, A: 255} }
								return btn.Layout(gtx)
							}),
						)
					})
				}),
				layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
					if state.CurrentTab == 0 {
						return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
							layout.Flexed(0.5, func(gtx layout.Context) layout.Dimensions {
								return DrawInsightsPanel(gtx, th, state)
							}),
							layout.Flexed(0.5, func(gtx layout.Context) layout.Dimensions {
								return DrawTreeView(gtx, th, state)
							}),
						)
					} else if state.CurrentTab == 1 {
						return DrawGamesPanel(gtx, th, state)
					} else if state.CurrentTab == 2 {
						return DrawIgnoredPanel(gtx, th, state)
					}
					return DrawSandboxTerminal(gtx, th, state)
				}),
			)
			
			call := m.Stop()
			call.Add(gtx.Ops)
			
			if state.ShowDialog {
				DrawActionDialog(gtx, th, state)
			}

			e.Frame(gtx.Ops)
		}
	}
}

func handleWizTreeImport(state *UIState, path string) {
	if state.IsScanning {
		return
	}
	state.IsScanning = true
	state.InsightCards = nil
	state.LogStream = append(state.LogStream, "Importing WizTree CSV from: "+path)

	go func() {
		defer func() {
			state.IsScanning = false
			state.Window.Invalidate()
		}()

		records, err := scanner.ImportWizTreeCSV(path)
		if err != nil {
			state.LogStream = append(state.LogStream, "[ERROR] Failed to import: "+err.Error())
			return
		}

		var cards []InsightCard
		var rawTopFolders []models.FileRecord

		for _, r := range records {
			if !r.IsDir {
				continue
			}
			p, dir := state.Registry.Detect(r.Path, true)

			if p != nil && dir == providers.StopTraversal {
				analysis, _ := p.Analyze(r.Path)
				actions := p.GetCleanupActions(r.Path)
				var cmdStr string
				if len(actions) > 0 {
					cmdStr = actions[0].Command
				}
				cards = append(cards, InsightCard{
					ProviderName: p.Name(),
					Risk:         analysis.Risk,
					Recoverable:  r.Size,
					Context:      analysis.Context + " (" + r.Path + ")",
					CommandStr:   cmdStr,
					Path:         r.Path,
				})
			} else if p == nil {
				rawTopFolders = append(rawTopFolders, r)
			}
		}

		// Add Global Providers
		for _, p := range state.Registry.GlobalProviders() {
			analysis, _ := p.Analyze("")
			actions := p.GetCleanupActions("")
			var cmdStr string
			if len(actions) > 0 {
				cmdStr = actions[0].Command
			}
			cards = append(cards, InsightCard{
				ProviderName: p.Name(),
				Risk:         analysis.Risk,
				Recoverable:  0, // Would need CLI to get actual size
				Context:      analysis.Context,
				CommandStr:   cmdStr,
			})
		}

		// Sort Cards: Safest first, then largest
		sort.Slice(cards, func(i, j int) bool {
			riskWeight := func(risk string) int {
				if risk == "Safe" {
					return 0
				}
				if risk == "Warning" {
					return 1
				}
				return 2
			}
			wI := riskWeight(cards[i].Risk)
			wJ := riskWeight(cards[j].Risk)
			if wI != wJ {
				return wI < wJ
			}
			return cards[i].Recoverable > cards[j].Recoverable
		})

		var normalCards, gameCards, customCards []InsightCard
		for _, c := range cards {
			if state.CustomState.Scores[c.Path] < 0 || state.CustomState.Descriptions[c.Path] != "" {
				customCards = append(customCards, c)
			} else if c.ProviderName == "Steam Game" || c.ProviderName == "Epic Games" {
				gameCards = append(gameCards, c)
			} else {
				normalCards = append(normalCards, c)
			}
		}
		
		state.InsightCards = normalCards
		state.GameCards = gameCards
		state.CustomCards = customCards

		state.RootNode = BuildTree(rawTopFolders)
		state.VisibleNodes = FlattenTree(state.RootNode, state.CustomState.Scores)

		state.LogStream = append(state.LogStream, fmt.Sprintf("Imported %d records. Found %d insights.", len(records), len(cards)))
	}()
}

func handleFolderScan(state *UIState, path string) {
	if state.IsScanning {
		return
	}
	state.IsScanning = true
	state.InsightCards = nil
	state.LogStream = append(state.LogStream, "Scanning directory: "+path)

	go func() {
		defer func() {
			state.IsScanning = false
			state.Window.Invalidate()
		}()

		engine := scanner.NewScanner(nil, state.Registry)
		records, err := engine.Scan(path, func(p string) {
			state.CurrentScanPath = p
			state.Window.Invalidate()
		})
		if err != nil {
			state.LogStream = append(state.LogStream, "[ERROR] Scan failed: "+err.Error())
			return
		}

		var cards []InsightCard
		var rawTopFolders []models.FileRecord

		for _, r := range records {
			if r.ProviderID != "" {
				p, _ := state.Registry.Detect(r.Path, true)
				if p != nil {
					analysis, _ := p.Analyze(r.Path)
					actions := p.GetCleanupActions(r.Path)
					var cmdStr string
					if len(actions) > 0 {
						cmdStr = actions[0].Command
					}
					cards = append(cards, InsightCard{
						ProviderName: p.Name(),
						Risk:         analysis.Risk,
						Recoverable:  r.Size,
						Context:      analysis.Context + " (" + r.Path + ")",
						CommandStr:   cmdStr,
						Path:         r.Path,
					})
				}
			} else if r.IsDir && r.MatchedRule == "" {
				rawTopFolders = append(rawTopFolders, r)
			}
		}

		for _, p := range state.Registry.GlobalProviders() {
			analysis, _ := p.Analyze("")
			actions := p.GetCleanupActions("")
			var cmdStr string
			if len(actions) > 0 {
				cmdStr = actions[0].Command
			}
			cards = append(cards, InsightCard{
				ProviderName: p.Name(),
				Risk:         analysis.Risk,
				Recoverable:  0,
				Context:      analysis.Context,
				CommandStr:   cmdStr,
			})
		}

		sort.Slice(cards, func(i, j int) bool {
			riskWeight := func(risk string) int {
				if risk == "Safe" {
					return 0
				}
				if risk == "Warning" {
					return 1
				}
				return 2
			}
			wI := riskWeight(cards[i].Risk)
			wJ := riskWeight(cards[j].Risk)
			if wI != wJ {
				return wI < wJ
			}
			return cards[i].Recoverable > cards[j].Recoverable
		})

		var normalCards, gameCards, customCards []InsightCard
		for _, c := range cards {
			if state.CustomState.Scores[c.Path] < 0 || state.CustomState.Descriptions[c.Path] != "" {
				customCards = append(customCards, c)
			} else if c.ProviderName == "Steam Game" || c.ProviderName == "Epic Games" {
				gameCards = append(gameCards, c)
			} else {
				normalCards = append(normalCards, c)
			}
		}
		
		state.InsightCards = normalCards
		state.GameCards = gameCards
		state.CustomCards = customCards

		state.RootNode = BuildTree(rawTopFolders)
		state.VisibleNodes = FlattenTree(state.RootNode, state.CustomState.Scores)

		state.LogStream = append(state.LogStream, fmt.Sprintf("Scanned %d records. Found %d insights.", len(records), len(cards)))
	}()
}
