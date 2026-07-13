package ui

import (
	"fmt"
	"image"
	"image/color"
	"strings"

	"gioui.org/font"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget/material"
)

func DrawTopBar(gtx layout.Context, th *material.Theme, state *UIState) layout.Dimensions {
	return layout.UniformInset(unit.Dp(8)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}.Layout(gtx,
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				btn := material.Button(th, &state.ScanFolderBtn, "Live Scan Folder...")
				if state.IsScanning {
					btn.Text = "Scanning..."
				}
				return btn.Layout(gtx)
			}),
			layout.Rigid(layout.Spacer{Width: unit.Dp(16)}.Layout),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return material.Body1(th, "WizTree CSV: ").Layout(gtx)
			}),
			layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
				ed := material.Editor(th, &state.WizTreePathEditor, "C:\\path\\to\\wiztree.csv")
				return ed.Layout(gtx)
			}),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				btn := material.Button(th, &state.BrowseWizTreeBtn, "Browse...")
				return btn.Layout(gtx)
			}),
			layout.Rigid(layout.Spacer{Width: unit.Dp(8)}.Layout),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				btn := material.Button(th, &state.ImportWizTreeBtn, "Import WizTree")
				if state.IsScanning {
					btn.Text = "..."
				}
				return btn.Layout(gtx)
			}),
			layout.Rigid(layout.Spacer{Width: unit.Dp(16)}.Layout),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				if state.IsScanning {
					txt := state.CurrentScanPath
					if len(txt) > 50 {
						txt = "..." + txt[len(txt)-47:]
					}
					return material.Body2(th, "Scanning: "+txt).Layout(gtx)
				}
				return layout.Dimensions{}
			}),
		)
	})
}

func DrawDiskOccupationBar(gtx layout.Context, th *material.Theme, used int64, total int64) layout.Dimensions {
	return layout.UniformInset(unit.Dp(16)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				title := material.H5(th, "Disk Occupation")
				return title.Layout(gtx)
			}),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				// Draw background bar
				height := gtx.Dp(unit.Dp(24))
				width := gtx.Constraints.Max.X

				bgRect := image.Rect(0, 0, width, height)
				paint.FillShape(gtx.Ops, color.NRGBA{R: 220, G: 220, B: 220, A: 255}, clip.Rect(bgRect).Op())

				// Draw used bar
				pct := float32(used) / float32(total)
				if pct > 1 {
					pct = 1
				}
				usedWidth := int(float32(width) * pct)
				usedRect := image.Rect(0, 0, usedWidth, height)
				paint.FillShape(gtx.Ops, color.NRGBA{R: 100, G: 150, B: 255, A: 255}, clip.Rect(usedRect).Op())

				// Return dimensions of the bar
				return layout.Dimensions{Size: image.Point{X: width, Y: height}}
			}),
		)
	})
}

func DrawInsightsPanel(gtx layout.Context, th *material.Theme, state *UIState) layout.Dimensions {
	return DrawInsightCardList(gtx, th, state, state.InsightCards, "Actionable Insights")
}

func DrawGamesPanel(gtx layout.Context, th *material.Theme, state *UIState) layout.Dimensions {
	return DrawInsightCardList(gtx, th, state, state.GameCards, "Detected Games")
}

func DrawIgnoredPanel(gtx layout.Context, th *material.Theme, state *UIState) layout.Dimensions {
	return DrawInsightCardList(gtx, th, state, state.CustomCards, "Custom / Ignored")
}

func DrawInsightCardList(gtx layout.Context, th *material.Theme, state *UIState, cards []InsightCard, titleStr string) layout.Dimensions {
	return layout.UniformInset(unit.Dp(8)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				title := material.H6(th, titleStr)
				return layout.UniformInset(unit.Dp(8)).Layout(gtx, title.Layout)
			}),
			layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
				return material.List(th, &state.InsightsList).Layout(gtx, len(cards), func(gtx layout.Context, index int) layout.Dimensions {
					card := &cards[index]
					return layout.UniformInset(unit.Dp(8)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
						bgColor := color.NRGBA{R: 240, G: 240, B: 240, A: 255}
						if card.Risk == "Safe" { bgColor = color.NRGBA{R: 200, G: 255, B: 200, A: 255} }
						if card.Risk == "Warning" { bgColor = color.NRGBA{R: 255, G: 255, B: 200, A: 255} }
						if card.Risk == "High" { bgColor = color.NRGBA{R: 255, G: 200, B: 200, A: 255} }

						m := op.Record(gtx.Ops)
						dims := layout.UniformInset(unit.Dp(12)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
							return material.Clickable(gtx, &card.CardBtn, func(gtx layout.Context) layout.Dimensions {
								mb := float64(card.Recoverable) / 1024 / 1024
								txt := fmt.Sprintf("[%s] %s - %.2f MB\n%s", card.Risk, card.ProviderName, mb, card.Context)
								label := material.Body1(th, txt)
								label.Alignment = text.Start
								btn := material.Button(th, &card.ActionBtn, "Reclaim Space")
								
								return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
									layout.Rigid(label.Layout),
									layout.Rigid(layout.Spacer{Height: unit.Dp(8)}.Layout),
									layout.Rigid(btn.Layout),
								)
							})
						})
						call := m.Stop()
						rect := image.Rect(0, 0, dims.Size.X, dims.Size.Y)
						paint.FillShape(gtx.Ops, bgColor, clip.Rect(rect).Op())
						call.Add(gtx.Ops)
						return dims
					})
				})
			}),
		)
	})
}

func DrawTreeView(gtx layout.Context, th *material.Theme, state *UIState) layout.Dimensions {
	return layout.UniformInset(unit.Dp(8)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				title := material.H6(th, "Largest Uncategorized Folders")
				desc := material.Body2(th, "Folders not claimed by any Provider.")
				return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
					layout.Rigid(title.Layout),
					layout.Rigid(desc.Layout),
					layout.Rigid(layout.Spacer{Height: unit.Dp(8)}.Layout),
				)
			}),
			layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
				return material.List(th, &state.TreeList).Layout(gtx, len(state.VisibleNodes), func(gtx layout.Context, index int) layout.Dimensions {
					node := state.VisibleNodes[index]
					mb := float64(node.Size) / 1024 / 1024
					if mb < 1 {
						return layout.Dimensions{}
					}
					
					depth := strings.Count(node.Path, "\\")
					if node.Path == "" { depth = 0 }
					padding := unit.Dp(float32(depth * 16))
					
					return layout.Inset{Left: padding, Top: unit.Dp(4), Bottom: unit.Dp(4)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
						return material.Clickable(gtx, &node.Clickable, func(gtx layout.Context) layout.Dimensions {
							return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}.Layout(gtx,
								layout.Rigid(func(gtx layout.Context) layout.Dimensions {
									prefix := "  "
									if node.IsDir {
										if node.Expanded {
											prefix = "v "
										} else {
											prefix = "> "
										}
									}
									lbl := material.Body1(th, prefix)
									lbl.Font.Weight = font.Bold
									return lbl.Layout(gtx)
								}),
								layout.Rigid(func(gtx layout.Context) layout.Dimensions {
									sizeLabel := material.Body1(th, fmt.Sprintf("%.1f MB", mb))
									sizeLabel.Font.Weight = font.Bold
									return sizeLabel.Layout(gtx)
								}),
								layout.Rigid(layout.Spacer{Width: unit.Dp(16)}.Layout),
								layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
									pathLabel := material.Body2(th, node.Name)
									if node.Name == "" { pathLabel = material.Body2(th, node.Path) }
									return pathLabel.Layout(gtx)
								}),
								layout.Rigid(func(gtx layout.Context) layout.Dimensions {
									score := state.CustomState.Scores[node.Path]
									if score != 0 {
										lbl := material.Body1(th, fmt.Sprintf("  [Score: %d]", score))
										lbl.Color = color.NRGBA{R: 150, G: 150, B: 150, A: 255}
										return lbl.Layout(gtx)
									}
									return layout.Dimensions{}
								}),
								layout.Rigid(func(gtx layout.Context) layout.Dimensions {
									desc := state.CustomState.Descriptions[node.Path]
									if desc != "" {
										lbl := material.Body1(th, "  ("+desc+")")
										lbl.Color = color.NRGBA{R: 50, G: 150, B: 50, A: 255}
										return lbl.Layout(gtx)
									}
									return layout.Dimensions{}
								}),
							)
						})
					})
				})
			}),
		)
	})
}

func DrawSandboxTerminal(gtx layout.Context, th *material.Theme, state *UIState) layout.Dimensions {
	return layout.UniformInset(unit.Dp(8)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		txt := "Sandbox Terminal:\n"
		start := len(state.LogStream) - 30
		if start < 0 { start = 0 }
		for _, log := range state.LogStream[start:] {
			txt += log + "\n"
		}
		label := material.Body2(th, txt)
		
		m := op.Record(gtx.Ops)
		dims := label.Layout(gtx)
		call := m.Stop()
		
		rect := image.Rect(0, 0, gtx.Constraints.Max.X, dims.Size.Y+16)
		paint.FillShape(gtx.Ops, color.NRGBA{R: 255, G: 255, B: 230, A: 255}, clip.Rect(rect).Op())
		call.Add(gtx.Ops)
		return layout.Dimensions{Size: rect.Max}
	})
}

func DrawActionDialog(gtx layout.Context, th *material.Theme, state *UIState) layout.Dimensions {
	paint.Fill(gtx.Ops, color.NRGBA{R: 0, G: 0, B: 0, A: 150})
	
	return layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		m := op.Record(gtx.Ops)
		dims := layout.UniformInset(unit.Dp(24)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{Axis: layout.Vertical, Alignment: layout.Middle}.Layout(gtx,
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return material.H6(th, "Folder Actions").Layout(gtx)
				}),
				layout.Rigid(layout.Spacer{Height: unit.Dp(16)}.Layout),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return material.Body1(th, state.SelectedNode.Path).Layout(gtx)
				}),
				layout.Rigid(layout.Spacer{Height: unit.Dp(24)}.Layout),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}.Layout(gtx,
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							return material.Button(th, &state.DecBtn, "  -  ").Layout(gtx)
						}),
						layout.Rigid(layout.Spacer{Width: unit.Dp(16)}.Layout),
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							score := state.CustomState.Scores[state.SelectedNode.Path]
							return material.Body1(th, fmt.Sprintf("Priority Score: %d", score)).Layout(gtx)
						}),
						layout.Rigid(layout.Spacer{Width: unit.Dp(16)}.Layout),
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							return material.Button(th, &state.IncBtn, "  +  ").Layout(gtx)
						}),
					)
				}),
				layout.Rigid(layout.Spacer{Height: unit.Dp(16)}.Layout),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					ed := material.Editor(th, &state.DescriptionEditor, "Custom Description...")
					return ed.Layout(gtx)
				}),
				layout.Rigid(layout.Spacer{Height: unit.Dp(8)}.Layout),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return material.Button(th, &state.SaveDescBtn, "Save Description").Layout(gtx)
				}),
				layout.Rigid(layout.Spacer{Height: unit.Dp(16)}.Layout),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return material.Button(th, &state.RuleBtn, "Create JSON Rule").Layout(gtx)
				}),
				layout.Rigid(layout.Spacer{Height: unit.Dp(16)}.Layout),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					btn := material.Button(th, &state.CloseDialogBtn, "Close")
					btn.Background = color.NRGBA{R: 200, G: 50, B: 50, A: 255}
					return btn.Layout(gtx)
				}),
			)
		})
		call := m.Stop()
		rect := image.Rect(0, 0, dims.Size.X, dims.Size.Y)
		paint.FillShape(gtx.Ops, color.NRGBA{R: 255, G: 255, B: 255, A: 255}, clip.Rect(rect).Op())
		call.Add(gtx.Ops)
		return dims
	})
}
