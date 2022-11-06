package main

import (
	"fmt"
	"image/color"
)

func enrichData(d *Data) {
	defaultLang := "en-us"
	defaultZoom := 150 //req'd for screens with high pixel density
	defaultHideDaysFrom := 90
	defaultHideWeeksFrom := 180

	if d.MySettings == nil {
		//sensible default settings: lang, end, zoom, hideDaysFrom, hideWeeksFrom
		d.MySettings = &Settings{defaultLang, "", defaultZoom, defaultHideDaysFrom, defaultHideWeeksFrom}
	}

	//partial settings are allowed, so supply defaults for each individually
	if d.MySettings.End != "" {
		//custom end date or placeholder found
		d.Last = calcLast(d.MySettings.End)
	}

	if d.MySettings.Lang == "" {
		d.MySettings.Lang = defaultLang
	}

	if d.MySettings.Zoom == 0 {
		d.MySettings.Zoom = defaultZoom
	}

	if d.MySettings.HideDaysFrom == 0 {
		d.MySettings.HideDaysFrom = defaultHideDaysFrom
	}

	if d.MySettings.HideWeeksFrom == 0 {
		d.MySettings.HideWeeksFrom = defaultHideWeeksFrom
	}

	//convert datestamps first
	for _, task := range d.Tasks {
		processTask(d, task)
	}

	setDefaults(d)
	applyTheme(d)

	// this assumes no leap years
	// d.Days = (d.Last.Year()*365 + d.Last.YearDay()) - (d.First.Year()*365 + d.First.YearDay())
	d.Days = int(d.Last.Sub(d.First).Hours()) / 24
	fmt.Println("d.Days new: ", d.Days)

	// safe layout defaults
	// show days if < 90 days; show months if < 180 days
	if d.MySettings.HideDaysFrom == 0 || d.MySettings.HideWeeksFrom == 0 {
		d.MySettings.HideDaysFrom = 90
		d.MySettings.HideWeeksFrom = 180
	}

	//zoom property defaults to 100%
	if d.MySettings.Zoom == 0 {
		d.MySettings.Zoom = 100
	}

	d.Scale = float64(d.MySettings.Zoom) / 100
	d.W, d.H, d.RowH, d.FontSize = 1024.0*d.Scale, 768.0*d.Scale, 20.0*d.Scale, 10.0*d.Scale
}

// TODO: correct for go error returns,
// make receiver:
func (d *Data) Validate() error {
	if d.Scale <= 0.0 || d.W <= 0.0 || d.H <= 0.0 || d.RowH <= 0.0 || d.FontSize <= 0.0 {
		return fmt.Errorf("parameter unitialized or invalid: Scale, W, H, RowH, FontSize = %.2f, %.2f, %.2f, %.2f, %.2f", d.Scale, d.W, d.H, d.RowH, d.FontSize)
	}

	if d.Days <= 0 {
		return fmt.Errorf("invalid number of days: %d", d.Days)
	}

	length := len(d.Tasks)
	if length == 0 {
		return fmt.Errorf("no tasks specified")
	}

	for index, task := range d.Tasks {
		if task.StartTime.Unix() > task.EndTime.Unix() {
			return fmt.Errorf("task #%d ends before it begins: %v", index+1, task)
		}

		//blank labels are allowed
		if task.StartTime.IsZero() || task.EndTime.IsZero() {
			return fmt.Errorf("task #%d is incomplete: %v", index+1, task)
		}

	}
	return nil
}

func setDefaults(d *Data) {
	d.FrameBorderColor = color.RGBA{0x00, 0x00, 0x00, 0xff}
	d.FrameFillColor = color.RGBA{0xff, 0xff, 0xff, 0xff}
	d.StripeColorDark = color.RGBA{0xdd, 0xdd, 0xdd, 0xff}
	d.StripeColorLight = color.RGBA{0xee, 0xee, 0xee, 0xff}
	d.GridColor = color.RGBA{0x99, 0x99, 0x99, 0xff}
}

func processTask(d *Data, t *Task) {
	t.StartTime = parseDateStamp(t.Start)

	//end time may be placeholder; if so, use currently known last date
	if t.End == "-" {
		t.EndTime = d.Last
	} else {
		t.EndTime = parseDateStamp(t.End)
	}

	if d.First.IsZero() || t.StartTime.Unix() < d.First.Unix() {
		d.First = t.StartTime
	}

	if d.Last.IsZero() || t.EndTime.Unix() > d.Last.Unix() {
		d.Last = t.EndTime
	}

	t.BorderColor = color.RGBA{0x55, 0x55, 0x55, 0xff}
	t.FillColor = color.RGBA{0xff, 0xff, 0xff, 0xff}
}
