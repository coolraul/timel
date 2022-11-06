package main

import (
	"fmt"
	"image/color"
	"strconv"
	"time"

	"github.com/fogleman/gg"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font/gofont/goregular"
)

func drawDropShadow(gc *gg.Context, x1, y1, x2, y2 float64) {
	opacity := uint8(0x44)
	offset := 4.0

	gc.Push()
	gc.SetFillStyle(gg.NewSolidPattern(color.RGBA{0x00, 0x00, 0x00, opacity}))
	gc.MoveTo(x1+offset, y1+offset)
	gc.LineTo(x2+offset, y1+offset)
	gc.LineTo(x2+offset, y2+offset)
	gc.LineTo(x1+offset, y2+offset)
	gc.LineTo(x1+offset, y1+offset)
	gc.Fill()
	gc.Pop()
}

func drawBlock(d *Data, gc *gg.Context, x1, y1, x2, y2 float64, strokeColor, fillColor color.Color, label string, showLabel bool) {
	gc.Push()
	gc.SetFillStyle(gg.NewSolidPattern(fillColor))
	gc.SetStrokeStyle(gg.NewSolidPattern(strokeColor))
	gc.SetLineWidth(1)
	gc.MoveTo(x1, y1)
	gc.LineTo(x2, y1)
	gc.LineTo(x2, y2)
	gc.LineTo(x1, y2)
	gc.LineTo(x1, y1)
	gc.Stroke()
	gc.LineTo(x2, y1)
	gc.LineTo(x2, y2)
	gc.LineTo(x1, y2)
	gc.LineTo(x1, y1)
	gc.Fill()
	gc.Pop()

	if !showLabel {
		return
	}

	//label
	gc.Push()
	gc.SetFillStyle(gg.NewSolidPattern(d.FrameBorderColor))

	tx1, ty1, tx2, ty2 := bounds(gc, label)

	w := x2 - x1
	tw := tx2 - tx1
	th := ty2 - ty1

	adjustedTextY := y2 - th*0.5
	//shift or hide if label too long
	if tw > w || x1+tw > d.ChartW {
		label = ""
	}

	midX := x1 + (x2-x1)*0.5
	adjustedTextX := midX - tw*0.5

	//normalize label x-position
	if adjustedTextX < 0 {
		adjustedTextX = 0
	} else {
		for adjustedTextX+tw > d.ChartW {
			adjustedTextX--
		}
	}

	gc.DrawString(label, adjustedTextX, adjustedTextY)
	gc.Pop()
}
func bounds(gc *gg.Context, s string) (float64, float64, float64, float64) {
	//a, b, c, d := gc.GetStringBounds(s)
	//return a, b, c, d
	w, h := gc.MeasureString(s)
	return 0.0, 0.0, w, h
}
func drawArrowHeadE(d *Data, gc *gg.Context, x, y float64) {
	gc.Push()
	gc.SetStrokeStyle(gg.NewSolidPattern(color.RGBA{0x00, 0x00, 0x00, 0xff}))
	gc.SetFillStyle(gg.NewSolidPattern(color.RGBA{0x00, 0x00, 0x00, 0xff}))
	r := 3.0 * d.Scale
	gc.MoveTo(x, y)
	gc.LineTo(x-2*r, y-r)
	gc.LineTo(x-r, y)
	gc.LineTo(x-2*r, y+r)
	gc.LineTo(x, y)
	gc.Fill()
	gc.Pop()
}
func drawArrowHeadS(d *Data, gc *gg.Context, x, y float64) {
	gc.Push()
	gc.SetStrokeStyle(gg.NewSolidPattern(color.RGBA{0x00, 0x00, 0x00, 0xff}))
	gc.SetFillStyle(gg.NewSolidPattern(color.RGBA{0x00, 0x00, 0x00, 0xff}))
	r := 3.0 * d.Scale
	gc.MoveTo(x, y)
	gc.LineTo(x-r, y-2*r)
	gc.LineTo(x, y-r)
	gc.LineTo(x+r, y-2*r)
	gc.LineTo(x, y)
	gc.Fill()
	gc.Pop()
}

func drawDateStamp(d *Data, gc *gg.Context, x, y float64, label string) {
	//vertical line
	gc.Push()
	gc.SetStrokeStyle(gg.NewSolidPattern(color.RGBA{0x00, 0x00, 0x00, 0xff}))
	gc.SetFillStyle(gg.NewSolidPattern(color.RGBA{0x00, 0x00, 0x00, 0xff}))
	r := 6 * d.Scale
	gc.MoveTo(x, y-2*r)
	gc.LineTo(x, y+2*r)
	gc.Stroke()
	gc.Pop()

	//label
	gc.Push()
	gc.SetFillStyle(gg.NewSolidPattern(color.RGBA{0x00, 0x00, 0x00, 0xff}))

	x1, y1, x2, y2 := bounds(gc, label)

	w := x2 - x1
	h := y2 - y1

	adjustedTextY := y + 2.5*h
	adjustedTextX := x - w*0.5

	//normalize label x-position
	if adjustedTextX < 0 {
		adjustedTextX = 0
	} else {
		for adjustedTextX+w > d.ChartW {
			adjustedTextX--
		}
	}

	//strip out year as it's shown at the top
	gc.DrawString(label, adjustedTextX, adjustedTextY)
	gc.Pop()
}

// drawCalendarGuides draws lines on each period change as indicated
// by the result of fn(time)
func drawCalendarGuides(d *Data, gc *gg.Context, y1, y2 float64, fn func(time.Time) string) {
	var period string

	gc.Push()
	gc.SetStrokeStyle(gg.NewSolidPattern(d.GridColor))

	for i := 0; i <= d.Days; i++ {
		// t := time.Date(d.First.Year(), d.First.Month(), d.First.Day()+i, 0, 0, 0, 0, time.UTC)
		t := d.First.AddDate(0, 0, i)

		currentPeriod := fn(t)
		println("currentPeriod:", currentPeriod, "t:", t.Format("Mon Jan 2 15:04:05 -0700 MST 2006"))

		if currentPeriod != period {
			x := float64(i) * d.DayW

			fmt.Println("period change, drawing line at x:", x)
			gc.Push()
			gc.SetDash(2.0, 2.0)
			gc.MoveTo(x, y1)
			gc.LineTo(x, float64(y2))
			gc.Stroke()
			gc.Pop()
			period = currentPeriod
		}
	}

	gc.Pop()
}

func drawCalendarRow(d *Data, gc *gg.Context, y float64, strokeColor, fillColor color.Color, fn func(time.Time) string) {
	var period string
	var from int
	for i := 0; i <= d.Days; i++ {
		t := time.Date(d.First.Year(), d.First.Month(), d.First.Day()+i, 0, 0, 0, 0, time.UTC)
		currentPeriod := fn(t)

		last := i == d.Days

		if i == 0 {
			prevT := time.Date(d.First.Year(), d.First.Month(), d.First.Day()-1, 0, 0, 0, 0, time.UTC)
			period = fn(prevT)
		}

		// new period or end of timeline
		if currentPeriod != period || last {
			x1 := float64(from) * d.DayW
			x2 := float64(i) * d.DayW
			if last {
				x2 += d.DayW
			}
			y2 := y + d.RowH

			drawBlock(d, gc, x1, y, x2, y2, strokeColor, fillColor, period, true)

			// now update from for next section
			period = currentPeriod
			from = i
		}
	}
}

func drawStripe(d *Data, gc *gg.Context, index int, y1, y2 float64) {
	color := d.StripeColorDark
	if index%2 != 0 {
		color = d.StripeColorLight
	}

	y1 -= d.RowH / 2
	y2 += d.RowH / 2

	gc.Push()
	gc.SetStrokeStyle(gg.NewSolidPattern(color))
	gc.SetFillStyle(gg.NewSolidPattern(color))
	gc.MoveTo(0, y1)
	gc.LineTo(d.ChartW, y1)
	gc.LineTo(d.ChartW, y2)
	gc.LineTo(0, y2)
	gc.LineTo(0, y1)
	gc.Stroke()
	gc.LineTo(d.ChartW, y1)
	gc.LineTo(d.ChartW, y2)
	gc.LineTo(0, y2)
	gc.LineTo(0, y1)
	gc.Fill()
	gc.Pop()
}

func setFontFace(gc *gg.Context) {
	font, _ := truetype.Parse(goregular.TTF)
	face := truetype.NewFace(font, &truetype.Options{
		Size: 14,
	})
	gc.SetFontFace(face)
}

func drawScene(d *Data) *gg.Context {
	width, rowH := d.W, d.RowH

	//TODO: calculate expected height - requires method
	heigth := d.RowH * (float64(len(d.Tasks)*2) + 4.9)
	fmt.Println("row height:", rowH, "heigth:", heigth, "width:", width)

	gc := gg.NewContext(int(width), int(heigth))

	setFontFace(gc)

	// calendar functions
	fnYear := func(t time.Time) string {
		return strconv.Itoa(t.Year())
	}
	fnMonth := func(t time.Time) string {
		return t.Format("Jan")
	}
	fnWeek := func(t time.Time) string {
		_, week := t.ISOWeek()
		return strconv.Itoa(week)
	}

	// label block
	var maxLabelWidth float64
	for _, task := range d.Tasks {
		label := task.Label

		x1, _, x2, _ := bounds(gc, label)

		labelWidth := x2 - x1
		if labelWidth > maxLabelWidth {
			maxLabelWidth = labelWidth
		}
	}

	d.LabelW = maxLabelWidth + 4    //leave room at start for milestones
	d.ChartW = width - d.LabelW - 4 //leave room at end for milestones

	//var dayW float64
	d.DayW = d.ChartW / float64(d.Days)
	var x, y float64

	// increment y as needed
	y = 0

	gc.Push()
	gc.Translate(d.LabelW, 0)
	// year
	drawCalendarRow(d, gc, y, d.FrameBorderColor, d.FrameFillColor, fnYear)
	// month
	y += rowH
	drawCalendarRow(d, gc, y, d.FrameBorderColor, d.FrameFillColor, fnMonth)

	// weeks
	if d.Days < d.MySettings.HideWeeksFrom {
		y += rowH
		drawCalendarRow(d, gc, y, d.FrameBorderColor, d.FrameFillColor, fnWeek)
	}

	// days
	if d.Days < d.MySettings.HideDaysFrom {
		y += rowH
		var weekend bool
		for i := 0; i <= d.Days; i++ {
			//determine if weekend
			t := time.Date(d.First.Year(), d.First.Month(), d.First.Day()+i, 0, 0, 0, 0, time.UTC)
			weekend = t.Weekday() == 0 || t.Weekday() == 6
			shade := d.StripeColorLight
			if weekend {
				shade = d.StripeColorDark
			}

			x = float64(i) * d.DayW

			drawBlock(d, gc, x, y, x+d.DayW, y+rowH, d.GridColor, shade, "", true)
		}
	}

	// stripes
	stripeY := y + 1.5*rowH + 1.0
	for index, _ := range d.Tasks {
		drawStripe(d, gc, index, stripeY, stripeY+d.RowH)
		stripeY += rowH * 2
	}

	// draw guides from calendar block onwards
	// guides
	var fnGuide func(t time.Time) string
	if d.Days < d.MySettings.HideDaysFrom {
		fnGuide = fnWeek
	} else if d.Days < d.MySettings.HideWeeksFrom {
		fnGuide = fnMonth
	} else {
		fnGuide = fnYear
	}

	y += rowH
	drawCalendarGuides(d, gc, y, y+float64(len(d.Tasks))*2.0*rowH+2*rowH, fnGuide)

	gc.Pop()

	y += d.RowH / 2

	for _, task := range d.Tasks {
		start := dayIndex(task.StartTime, d.First, d.Last)
		end := dayIndex(task.EndTime, d.First, d.Last)

		if start == -1 || end == -1 {
			fmt.Println("can't place task on timeline")
			continue
		}

		//one-day tasks: draw full day
		if start == end {
			end++
		}

		x1 := float64(start) * d.DayW
		x2 := float64(end) * d.DayW
		y1 := y
		y2 := y + rowH

		gc.Push()
		gc.Translate(d.LabelW, 0)
		drawDropShadow(gc, x1, y1, x2, y2)
		drawBlock(d, gc, x1, y1, x2, y2, task.BorderColor, task.FillColor, task.Label, false)
		gc.Pop()

		//write out label
		label := task.Label
		gc.Push()
		gc.SetFillStyle(gg.NewSolidPattern(color.RGBA{0x00, 0x00, 0x00, 0xff}))

		_, ty1, _, ty2 := bounds(gc, label)

		th := ty2 - ty1
		adjustedTextY := y2 - th*0.5

		gc.DrawString(label, 0, adjustedTextY)
		gc.Pop()

		y += d.RowH * 2
	}

	y += rowH
	//TODO: crop if less tall than precalculated image

	return gc //cropped
}
