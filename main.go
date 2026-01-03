package main

import (
	"fmt"
	"image"
	"image/color"
	"os"
	"strconv"
	"strings"

	"gioui.org/app"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

// 颜色定义
var (
	darkGreen   = color.NRGBA{R: 30, G: 60, B: 50, A: 255}    // 深绿色背景
	buttonGreen = color.NRGBA{R: 40, G: 80, B: 65, A: 255}    // 按钮深绿色
	brightGreen = color.NRGBA{R: 0, G: 200, B: 100, A: 255}   // 亮绿色（运算符和AC等）
	equalsGreen = color.NRGBA{R: 0, G: 220, B: 110, A: 255}   // 等号按钮绿色
	white       = color.NRGBA{R: 255, G: 255, B: 255, A: 255} // 白色文字
	lightGray   = color.NRGBA{R: 200, G: 200, B: 200, A: 255} // 浅灰色（之前的计算）
	red         = color.NRGBA{R: 255, G: 80, B: 80, A: 255}   // 红色（退格图标）
)

func main() {
	go func() {
		defer os.Exit(0)
		w := app.NewWindow(
			app.Decorated(false),
			app.Size(unit.Dp(400), unit.Dp(700)),
		)
		calc := NewCalculator()
		calc.Run(w)
	}()
	app.Main()
}

type Calculator struct {
	display *widget.Editor
	buttons [][]widget.Clickable
	theme   *material.Theme

	window *app.Window

	// 标题栏按钮
	menuBtn    widget.Clickable
	historyBtn widget.Clickable

	// 计算状态
	currentValue     float64
	previousValue    float64
	operation        string
	shouldReset      bool
	previousExprText string
}

func NewCalculator() *Calculator {
	display := &widget.Editor{
		SingleLine: true,
		ReadOnly:   true,
		Alignment:  text.End,
	}
	display.SetText("0")

	// 创建按钮网格 5行4列
	buttons := make([][]widget.Clickable, 5)
	for i := range buttons {
		buttons[i] = make([]widget.Clickable, 4)
	}

	theme := material.NewTheme()
	theme.Palette.Fg = white
	theme.Palette.Bg = darkGreen

	return &Calculator{
		display: display,
		buttons: buttons,
		theme:   theme,
	}
}

func (c *Calculator) Run(w *app.Window) error {
	c.window = w
	var ops op.Ops

	for {
		e := w.NextEvent()
		switch e := e.(type) {
		case app.DestroyEvent:
			return e.Err
		case app.FrameEvent:
			ops.Reset()
			gtx := app.NewContext(&ops, e)
			c.Layout(gtx)
			e.Frame(gtx.Ops)
		}
	}
}

func (c *Calculator) Layout(gtx layout.Context) layout.Dimensions {
	c.handleEvents(gtx)

	// 填充深绿色背景
	paint.Fill(gtx.Ops, darkGreen)

	return layout.Flex{
		Axis:    layout.Vertical,
		Spacing: layout.SpaceStart,
	}.Layout(gtx,
		// 标题栏
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return c.layoutTitleBar(gtx)
		}),
		// 显示区域
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return c.layoutDisplay(gtx)
		}),
		// 按钮网格
		layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
			return layout.Inset{
				Top:    unit.Dp(20),
				Bottom: unit.Dp(20),
				Left:   unit.Dp(20),
				Right:  unit.Dp(20),
			}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return c.layoutButtons(gtx)
			})
		}),
	)
}

func (c *Calculator) layoutTitleBar(gtx layout.Context) layout.Dimensions {
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
			// 关于按钮（左上角）
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return c.menuBtn.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					label := material.Body2(c.theme, "关于")
					label.Color = white
					label.Alignment = text.Start
					label.TextSize = unit.Sp(14)
					return label.Layout(gtx)
				})
			}),
			// Standard 文字（中间）
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				label := material.Body1(c.theme, "Calculator")
				label.Color = white
				label.Alignment = text.Middle
				label.TextSize = unit.Sp(16)
				return label.Layout(gtx)
			}),
			// 关闭按钮（右上角）
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				gtx.Constraints.Min = image.Pt(24, 24)
				gtx.Constraints.Max = image.Pt(24, 24)
				return c.historyBtn.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					c.drawCloseIcon(gtx)
					return layout.Dimensions{Size: gtx.Constraints.Max}
				})
			}),
		)
	})
}

func (c *Calculator) drawMenuIcon(gtx layout.Context) {
	// 绘制三条横线
	lineHeight := float32(2.0)
	lineWidth := float32(18.0)
	lineSpacing := float32(6.0)
	startY := float32(6.0)

	for i := 0; i < 3; i++ {
		y := startY + float32(i)*lineSpacing
		xOffset := float32(3.0)
		rect := image.Rectangle{
			Min: image.Pt(int(xOffset), int(y)),
			Max: image.Pt(int(xOffset+lineWidth), int(y+lineHeight)),
		}
		rr := clip.UniformRRect(rect, 1)
		paint.FillShape(gtx.Ops, white, rr.Op(gtx.Ops))
	}
}

func (c *Calculator) drawCloseIcon(gtx layout.Context) {
	// 绘制关闭图标（X）- 使用小矩形组成两条交叉线
	center := image.Pt(12, 12)
	size := 10 // X的边长的一半
	lineWidth := 2

	// 绘制两条交叉线形成X
	// 线1: 左上到右下
	for i := -size; i <= size; i++ {
		rect := image.Rectangle{
			Min: image.Pt(center.X+i-lineWidth/2, center.Y+i-lineWidth/2),
			Max: image.Pt(center.X+i+lineWidth/2, center.Y+i+lineWidth/2),
		}
		rr := clip.UniformRRect(rect, 1)
		paint.FillShape(gtx.Ops, white, rr.Op(gtx.Ops))
	}

	// 线2: 右上到左下
	for i := -size; i <= size; i++ {
		rect := image.Rectangle{
			Min: image.Pt(center.X+i-lineWidth/2, center.Y-i-lineWidth/2),
			Max: image.Pt(center.X+i+lineWidth/2, center.Y-i+lineWidth/2),
		}
		rr := clip.UniformRRect(rect, 1)
		paint.FillShape(gtx.Ops, white, rr.Op(gtx.Ops))
	}
}

func (c *Calculator) layoutDisplay(gtx layout.Context) layout.Dimensions {
	// 固定显示区域高度，足够显示三行内容
	fixedHeight := gtx.Dp(unit.Dp(180)) // 固定高度约180dp，可容纳三行

	return layout.Inset{
		Top:    unit.Dp(30),
		Bottom: unit.Dp(30),
		Left:   unit.Dp(30),
		Right:  unit.Dp(30),
	}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		// 设置固定高度约束
		gtx.Constraints.Min.Y = fixedHeight
		gtx.Constraints.Max.Y = fixedHeight

		return layout.Flex{
			Axis:      layout.Vertical,
			Spacing:   layout.SpaceEnd, // 内容靠底部，间距紧密
			Alignment: layout.End,      // 右对齐
		}.Layout(gtx,
			// 第一行：预留空间（顶部空白，推动内容到底部）
			layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
				return layout.Dimensions{
					Size: image.Pt(gtx.Constraints.Max.X, gtx.Constraints.Max.Y),
				}
			}),
			// 第二行：之前的计算表达式（小字、灰色）
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				if c.previousExprText != "" {
					return layout.Inset{
						Bottom: unit.Dp(5), // 表达式和结果之间的小间距
					}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
						label := material.Caption(c.theme, c.previousExprText)
						label.Color = lightGray
						label.Alignment = text.End
						label.TextSize = unit.Sp(14)
						return label.Layout(gtx)
					})
				}
				return layout.Dimensions{}
			}),
			// 第三行：当前结果（大字体、白色）- 紧挨着表达式
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				displayText := c.display.Text()
				label := material.H1(c.theme, displayText)
				label.Color = white
				label.Alignment = text.End
				label.TextSize = unit.Sp(48)
				return label.Layout(gtx)
			}),
		)
	})
}

func (c *Calculator) layoutButtons(gtx layout.Context) layout.Dimensions {
	buttonLabels := [][]string{
		{"AC", "±", "%", "⌫"},
		{"7", "8", "9", "÷"},
		{"4", "5", "6", "×"},
		{"1", "2", "3", "+"},
		{".", "0", "=", ""}, // 最后一行的=会跨两列
	}

	// 计算按钮大小：4列布局，留出间距
	availableWidth := gtx.Constraints.Max.X
	availableHeight := gtx.Constraints.Max.Y
	buttonGap := gtx.Dp(unit.Dp(10))
	buttonSize := (availableWidth - buttonGap*3) / 4
	maxHeight := (availableHeight - buttonGap*4) / 5
	if buttonSize > maxHeight {
		buttonSize = maxHeight
	}
	if buttonSize < 40 {
		buttonSize = 40
	}

	return layout.Flex{
		Axis:    layout.Vertical,
		Spacing: layout.SpaceBetween,
	}.Layout(gtx,
		// 第一行：AC ± % ⌫
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return c.layoutButtonRow(gtx, buttonLabels[0], 0, buttonSize, buttonGap)
		}),
		// 第二行：7 8 9 ÷
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return c.layoutButtonRow(gtx, buttonLabels[1], 1, buttonSize, buttonGap)
		}),
		// 第三行：4 5 6 ×
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return c.layoutButtonRow(gtx, buttonLabels[2], 2, buttonSize, buttonGap)
		}),
		// 第四行：1 2 3 +
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return c.layoutButtonRow(gtx, buttonLabels[3], 3, buttonSize, buttonGap)
		}),
		// 第五行：. 0 =（=跨两列）
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{
				Axis:    layout.Horizontal,
				Spacing: layout.SpaceBetween,
			}.Layout(gtx,
				// . 按钮
				c.button(gtx, &c.buttons[4][0], ".", buttonSize),
				// 0 按钮
				c.button(gtx, &c.buttons[4][1], "0", buttonSize),
				// = 按钮（跨两列，占据第3和第4列的位置）
				c.buttonWide(gtx, &c.buttons[4][2], "=", buttonSize*2+buttonGap, buttonSize),
			)
		}),
	)
}

func (c *Calculator) layoutButtonRow(gtx layout.Context, labels []string, row int, buttonSize, buttonGap int) layout.Dimensions {
	return layout.Flex{
		Axis:    layout.Horizontal,
		Spacing: layout.SpaceBetween,
	}.Layout(gtx,
		c.button(gtx, &c.buttons[row][0], labels[0], buttonSize),
		c.button(gtx, &c.buttons[row][1], labels[1], buttonSize),
		c.button(gtx, &c.buttons[row][2], labels[2], buttonSize),
		c.button(gtx, &c.buttons[row][3], labels[3], buttonSize),
	)
}

func (c *Calculator) buttonWide(gtx layout.Context, btn *widget.Clickable, label string, width, height int) layout.FlexChild {
	return layout.Rigid(func(gtx layout.Context) layout.Dimensions {
		gtx.Constraints = layout.Exact(image.Pt(width, height))

		// 确定按钮样式
		bgColor, textColor := c.getButtonColors(label)

		// 绘制圆角矩形背景
		r := op.Record(gtx.Ops)
		rect := image.Rectangle{Max: gtx.Constraints.Max}
		radius := gtx.Dp(unit.Dp(40)) // 更圆润的圆角
		rr := clip.UniformRRect(rect, radius)
		paint.FillShape(gtx.Ops, bgColor, rr.Op(gtx.Ops))
		call := r.Stop()

		// 按钮文字 - 使用 Center 布局确保完全居中
		labelWidget := material.Body1(c.theme, label)
		labelWidget.Color = textColor
		if label == "⌫" {
			// 退格按钮使用特殊颜色
			labelWidget.Color = red
		}
		labelWidget.Alignment = text.Middle
		labelWidget.TextSize = unit.Sp(28) // 更大的字体

		// 点击区域和文字居中布局
		clickable := btn.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			call.Add(gtx.Ops)
			// 使用 Center 布局确保文字在按钮中完全居中
			return layout.Center.Layout(gtx, labelWidget.Layout)
		})

		return clickable
	})
}

func (c *Calculator) button(gtx layout.Context, btn *widget.Clickable, label string, buttonSize int) layout.FlexChild {
	return layout.Rigid(func(gtx layout.Context) layout.Dimensions {
		gtx.Constraints = layout.Exact(image.Pt(buttonSize, buttonSize))

		// 确定按钮样式
		bgColor, textColor := c.getButtonColors(label)

		// 绘制圆角矩形背景
		r := op.Record(gtx.Ops)
		rect := image.Rectangle{Max: gtx.Constraints.Max}
		radius := gtx.Dp(unit.Dp(40)) // 更圆润的圆角
		rr := clip.UniformRRect(rect, radius)
		paint.FillShape(gtx.Ops, bgColor, rr.Op(gtx.Ops))
		call := r.Stop()

		// 按钮文字 - 使用 Center 布局确保完全居中
		labelWidget := material.Body1(c.theme, label)
		labelWidget.Color = textColor
		if label == "⌫" {
			// 退格按钮使用特殊颜色
			labelWidget.Color = red
		}
		labelWidget.Alignment = text.Middle
		labelWidget.TextSize = unit.Sp(28) // 更大的字体

		// 点击区域和文字居中布局
		clickable := btn.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			call.Add(gtx.Ops)
			// 使用 Center 布局确保文字在按钮中完全居中
			return layout.Center.Layout(gtx, labelWidget.Layout)
		})

		return clickable
	})
}

func (c *Calculator) getButtonColors(label string) (bgColor, textColor color.NRGBA) {
	switch label {
	case "=":
		return equalsGreen, white
	case "+", "-", "×", "÷":
		return buttonGreen, brightGreen
	case "AC", "%":
		return buttonGreen, brightGreen
	case "⌫":
		return buttonGreen, red
	case "±":
		return buttonGreen, brightGreen
	default: // 数字和点
		return buttonGreen, white
	}
}

func (c *Calculator) handleEvents(gtx layout.Context) {
	buttonLabels := [][]string{
		{"AC", "±", "%", "⌫"},
		{"7", "8", "9", "÷"},
		{"4", "5", "6", "×"},
		{"1", "2", "3", "+"},
		{".", "0", "=", ""},
	}

	for i := range c.buttons {
		for j := range c.buttons[i] {
			if c.buttons[i][j].Clicked(gtx) {
				label := buttonLabels[i][j]
				if label != "" {
					c.handleButtonClick(label)
				}
			}
		}
	}

	// 处理标题栏按钮
	if c.menuBtn.Clicked(gtx) {
		// 打开关于窗口
		go func() {
			aboutWindow := app.NewWindow(
				app.Decorated(false),
				app.Title("关于"),
				app.Size(unit.Dp(380), unit.Dp(500)),
			)
			about := NewAboutWindow()
			about.Run(aboutWindow)
		}()
	}
	if c.historyBtn.Clicked(gtx) {
		// 关闭窗口
		if c.window != nil {
			fmt.Println("history button clicked")
			os.Exit(0)
		}
	}
}

func (c *Calculator) handleButtonClick(label string) {
	currentText := c.display.Text()

	switch label {
	case "AC": // 全部清除
		c.reset()
		c.updateDisplay("0")
		c.previousExprText = ""
	case "⌫": // 退格
		if len(currentText) > 1 {
			c.updateDisplay(currentText[:len(currentText)-1])
		} else {
			c.updateDisplay("0")
		}
	case "±": // 正负号
		if currentText != "0" {
			if strings.HasPrefix(currentText, "-") {
				c.updateDisplay(currentText[1:])
			} else {
				c.updateDisplay("-" + currentText)
			}
		}
	case "%": // 百分比
		val := c.parseDisplay()
		c.updateDisplay(c.formatNumber(val / 100))
		c.shouldReset = true
	case "=":
		c.calculate()
	case "+", "-", "×", "÷":
		if c.operation != "" {
			c.calculate()
		}
		c.previousValue = c.parseDisplay()
		c.operation = label
		c.shouldReset = true
		// 更新之前的表达式显示
		c.previousExprText = c.formatNumber(c.previousValue) + " " + label
	case ".":
		if c.shouldReset {
			c.updateDisplay("0.")
			c.shouldReset = false
		} else if !strings.Contains(currentText, ".") {
			c.updateDisplay(currentText + ".")
		}
	case "0", "1", "2", "3", "4", "5", "6", "7", "8", "9":
		if c.shouldReset || currentText == "0" {
			c.updateDisplay(label)
			c.shouldReset = false
		} else {
			c.updateDisplay(currentText + label)
		}
		// 格式化显示（添加千位分隔符）
		val := c.parseDisplay()
		c.updateDisplay(c.formatNumber(val))
	}
}

func (c *Calculator) calculate() {
	current := c.parseDisplay()
	var result float64

	switch c.operation {
	case "+":
		result = c.previousValue + current
	case "-":
		result = c.previousValue - current
	case "×":
		result = c.previousValue * current
	case "÷":
		if current != 0 {
			result = c.previousValue / current
		} else {
			c.updateDisplay("错误")
			c.previousValue = 0
			c.operation = ""
			c.shouldReset = true
			c.previousExprText = ""
			return
		}
	default:
		return
	}

	// 更新之前的表达式显示
	c.previousExprText = c.formatNumber(c.previousValue) + " " + c.operation + " " + c.formatNumber(current)
	c.updateDisplay(c.formatNumber(result))
	c.previousValue = result
	c.operation = ""
	c.shouldReset = true
}

func (c *Calculator) parseDisplay() float64 {
	text := strings.ReplaceAll(strings.TrimSpace(c.display.Text()), ",", "")
	if text == "错误" {
		return 0
	}
	val, err := strconv.ParseFloat(text, 64)
	if err != nil {
		return 0
	}
	return val
}

func (c *Calculator) formatNumber(num float64) string {
	// 如果是整数，显示为整数；否则显示为小数
	var str string
	if num == float64(int64(num)) {
		str = fmt.Sprintf("%.0f", num)
	} else {
		str = fmt.Sprintf("%g", num)
	}

	// 添加千位分隔符
	parts := strings.Split(str, ".")
	intPart := parts[0]
	negative := false
	if strings.HasPrefix(intPart, "-") {
		intPart = intPart[1:]
		negative = true
	}
	result := c.addCommas(intPart)
	if len(parts) > 1 {
		result += "." + parts[1]
	}
	if negative {
		return "-" + result
	}
	return result
}

func (c *Calculator) addCommas(s string) string {
	if len(s) <= 3 {
		return s
	}
	var result strings.Builder
	for i, r := range s {
		if i > 0 && (len(s)-i)%3 == 0 {
			result.WriteString(",")
		}
		result.WriteRune(r)
	}
	return result.String()
}

func (c *Calculator) updateDisplay(text string) {
	c.display.SetText(text)
	if c.window != nil {
		c.window.Invalidate()
	}
}

func (c *Calculator) reset() {
	c.currentValue = 0
	c.previousValue = 0
	c.operation = ""
	c.shouldReset = false
	c.previousExprText = ""
}
