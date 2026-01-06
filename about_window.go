package main

import (
	"fmt"
	"image/color"

	"gioui.org/app"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

var (
	Version   = "v0.0.0"
	GitCommit = "unknown"
	BuildTime = "unknown"
)

// AboutWindow 关于窗口
type AboutWindow struct {
	theme       *material.Theme
	window      *app.Window
	closeBtn    widget.Clickable // 右上角关闭按钮
	closeBtnBot widget.Clickable // 底部关闭按钮
	scrollView  widget.List
}

func NewAboutWindow() *AboutWindow {
	theme := material.NewTheme()
	theme.Palette.Fg = white
	theme.Palette.Bg = darkGreen

	return &AboutWindow{
		theme:      theme,
		scrollView: widget.List{List: layout.List{Axis: layout.Vertical}},
	}
}

func (a *AboutWindow) Run(w *app.Window) error {
	a.window = w
	var ops op.Ops

	for {
		e := w.Event()
		switch e := e.(type) {
		case app.DestroyEvent:
			fmt.Println("sub window closed")
			return e.Err

		case app.FrameEvent:
			ops.Reset()
			gtx := app.NewContext(&ops, e)

			// 布局
			a.Layout(gtx)

			e.Frame(gtx.Ops)
		}
	}
}

func (a *AboutWindow) Layout(gtx layout.Context) layout.Dimensions {
	// 填充深绿色背景
	dialogBg := color.NRGBA{R: 45, G: 75, B: 65, A: 255}
	paint.Fill(gtx.Ops, dialogBg)

	return layout.Flex{
		Axis:    layout.Vertical,
		Spacing: layout.SpaceStart,
	}.Layout(gtx,
		// 标题栏
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return a.layoutTitleBar(gtx)
		}),
		// 内容区域
		layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
			return layout.Inset{
				Top:    unit.Dp(20),
				Bottom: unit.Dp(20),
				Left:   unit.Dp(20),
				Right:  unit.Dp(20),
			}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return a.layoutContent(gtx)
			})
		}),
	)
}

func (a *AboutWindow) layoutTitleBar(gtx layout.Context) layout.Dimensions {
	return layout.Inset{
		Top:    unit.Dp(15),
		Bottom: unit.Dp(15),
		Left:   unit.Dp(20),
		Right:  unit.Dp(20),
	}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return layout.Flex{
			Axis:    layout.Horizontal,
			Spacing: layout.SpaceBetween,
		}.Layout(gtx,
			// 标题
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				label := material.H6(a.theme, "About")
				label.Color = white
				label.Alignment = text.Start
				return label.Layout(gtx)
			}),
		)
	})
}

func (a *AboutWindow) layoutContent(gtx layout.Context) layout.Dimensions {
	return a.scrollView.Layout(gtx, 1, func(gtx layout.Context, index int) layout.Dimensions {
		return layout.Flex{
			Axis:    layout.Vertical,
			Spacing: layout.SpaceStart,
		}.Layout(gtx,
			// 作者
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return layout.Inset{Bottom: unit.Dp(15)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					return layout.Flex{
						Axis: layout.Horizontal,
					}.Layout(gtx,
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							label := material.Body1(a.theme, "Author: ")
							label.Color = lightGray
							label.TextSize = unit.Sp(14)
							return label.Layout(gtx)
						}),
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							label := material.Body1(a.theme, "panyingyun@gmail.com")
							label.Color = white
							label.TextSize = unit.Sp(14)
							return label.Layout(gtx)
						}),
					)
				})
			}),
			// 目的
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return layout.Inset{Bottom: unit.Dp(15)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					return layout.Flex{
						Axis: layout.Horizontal,
					}.Layout(gtx,
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							label := material.Body1(a.theme, "Purpose: ")
							label.Color = lightGray
							label.TextSize = unit.Sp(14)
							return label.Layout(gtx)
						}),
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							label := material.Body1(a.theme, "How to write Go GUI app with gio")
							label.Color = white
							label.TextSize = unit.Sp(14)
							return label.Layout(gtx)
						}),
					)
				})
			}),

			// 目的
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return layout.Inset{Bottom: unit.Dp(15)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					return layout.Flex{
						Axis: layout.Horizontal,
					}.Layout(gtx,
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							label := material.Body1(a.theme, "Version: ")
							label.Color = lightGray
							label.TextSize = unit.Sp(14)
							return label.Layout(gtx)
						}),
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							label := material.Body1(a.theme, Version)
							label.Color = white
							label.TextSize = unit.Sp(14)
							return label.Layout(gtx)
						}),
					)
				})
			}),

			// 目的
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return layout.Inset{Bottom: unit.Dp(15)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					return layout.Flex{
						Axis: layout.Horizontal,
					}.Layout(gtx,
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							label := material.Body1(a.theme, "GitCommit: ")
							label.Color = lightGray
							label.TextSize = unit.Sp(14)
							return label.Layout(gtx)
						}),
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							label := material.Body1(a.theme, GitCommit)
							label.Color = white
							label.TextSize = unit.Sp(14)
							return label.Layout(gtx)
						}),
					)
				})
			}),

			// 目的
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return layout.Inset{Bottom: unit.Dp(15)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					return layout.Flex{
						Axis: layout.Horizontal,
					}.Layout(gtx,
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							label := material.Body1(a.theme, "BuildTime: ")
							label.Color = lightGray
							label.TextSize = unit.Sp(14)
							return label.Layout(gtx)
						}),
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							label := material.Body1(a.theme, BuildTime)
							label.Color = white
							label.TextSize = unit.Sp(14)
							return label.Layout(gtx)
						}),
					)
				})
			}),
		)
	})
}
