package main

import (
	"fmt"
	"image"
	"image/color"

	"gioui.org/app"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
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
			fmt.Println("子窗口已关闭")
			return e.Err

		case app.FrameEvent:
			ops.Reset()
			gtx := app.NewContext(&ops, e)

			// 先布局，让按钮渲染
			a.Layout(gtx)

			// fmt.Println("FrameEvent Clicked ...", a.closeBtn.Clicked(gtx), a.closeBtnBot.Clicked(gtx))
			// fmt.Println("FrameEvent Hovered ...", a.closeBtn.Hovered(), a.closeBtnBot.Hovered())
			// fmt.Println("FrameEvent Pressed ...", a.closeBtn.Pressed(), a.closeBtnBot.Pressed())
			// fmt.Println("FrameEvent History ...", a.closeBtn.History(), a.closeBtnBot.History())

			// 然后检查按钮点击事件
			if a.closeBtn.Pressed() || a.closeBtnBot.Pressed() {
				w.Perform(system.ActionClose)
				fmt.Println("close button clicked")
				return nil // 正常关闭窗口
			}

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
		// 底部关闭按钮
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Inset{
				Top:    unit.Dp(10),
				Bottom: unit.Dp(20),
				Left:   unit.Dp(20),
				Right:  unit.Dp(20),
			}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return a.layoutCloseButton(gtx)
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
				label := material.H6(a.theme, "关于")
				label.Color = white
				label.Alignment = text.Start
				return label.Layout(gtx)
			}),
			// 关闭按钮
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				gtx.Constraints.Min = image.Pt(24, 24)
				gtx.Constraints.Max = image.Pt(24, 24)
				return a.closeBtn.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					a.drawCloseIcon(gtx)
					return layout.Dimensions{Size: gtx.Constraints.Max}
				})
			}),
		)
	})
}

func (a *AboutWindow) drawCloseIcon(gtx layout.Context) {
	center := image.Pt(12, 12)
	size := 10
	lineWidth := 2

	// 绘制X图标
	for i := -size; i <= size; i++ {
		rect := image.Rectangle{
			Min: image.Pt(center.X+i-lineWidth/2, center.Y+i-lineWidth/2),
			Max: image.Pt(center.X+i+lineWidth/2, center.Y+i+lineWidth/2),
		}
		rr := clip.UniformRRect(rect, 1)
		paint.FillShape(gtx.Ops, white, rr.Op(gtx.Ops))
	}

	for i := -size; i <= size; i++ {
		rect := image.Rectangle{
			Min: image.Pt(center.X+i-lineWidth/2, center.Y-i-lineWidth/2),
			Max: image.Pt(center.X+i+lineWidth/2, center.Y-i+lineWidth/2),
		}
		rr := clip.UniformRRect(rect, 1)
		paint.FillShape(gtx.Ops, white, rr.Op(gtx.Ops))
	}
}

func (a *AboutWindow) layoutCloseButton(gtx layout.Context) layout.Dimensions {
	// 底部关闭按钮
	return layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		buttonWidth := gtx.Dp(unit.Dp(120))
		buttonHeight := gtx.Dp(unit.Dp(40))
		gtx.Constraints.Min = image.Pt(buttonWidth, buttonHeight)
		gtx.Constraints.Max = image.Pt(buttonWidth, buttonHeight)

		// 先记录背景绘制操作
		r := op.Record(gtx.Ops)
		rr := clip.UniformRRect(image.Rectangle{Max: gtx.Constraints.Max}, gtx.Dp(unit.Dp(5)))
		paint.FillShape(gtx.Ops, buttonGreen, rr.Op(gtx.Ops))
		call := r.Stop()

		// 使用 closeBtnBot 处理点击，内部绘制背景和文字
		return a.closeBtnBot.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			// 绘制按钮背景
			call.Add(gtx.Ops)
			// 绘制按钮文字
			label := material.Body1(a.theme, "关闭")
			label.Color = white
			label.Alignment = text.Middle
			label.TextSize = unit.Sp(16)
			return layout.Center.Layout(gtx, label.Layout)
		})
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
							label := material.Body1(a.theme, "作者：")
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
							label := material.Body1(a.theme, "目的：")
							label.Color = lightGray
							label.TextSize = unit.Sp(14)
							return label.Layout(gtx)
						}),
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							label := material.Body1(a.theme, "学习gioui编写Go GUI应用")
							label.Color = white
							label.TextSize = unit.Sp(14)
							return label.Layout(gtx)
						}),
					)
				})
			}),
			// 许可证标题
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return layout.Inset{Bottom: unit.Dp(10)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					label := material.Body1(a.theme, "许可证：")
					label.Color = lightGray
					label.TextSize = unit.Sp(14)
					return label.Layout(gtx)
				})
			}),
			// GPLv3许可证文本
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				licenseText := "GNU GENERAL PUBLIC LICENSE\nVersion 3, 29 June 2007\n\n" +
					"Copyright (C) 2007 Free Software Foundation, Inc.\n\n" +
					"This program is free software: you can redistribute it and/or modify it under the terms of the GNU General Public License as published by the Free Software Foundation, either version 3 of the License, or (at your option) any later version.\n\n" +
					"This program is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU General Public License for more details.\n\n" +
					"You should have received a copy of the GNU General Public License along with this program. If not, see <https://www.gnu.org/licenses/>."

				label := material.Body2(a.theme, licenseText)
				label.Color = white
				label.Alignment = text.Start
				label.TextSize = unit.Sp(11)
				return label.Layout(gtx)
			}),
		)
	})
}
