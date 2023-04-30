package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/widget"
)

var listData = []string{
	"Item 1",
	"Item 2",
	"Item 3",
}

func main() {
	a := app.New()
	w := a.NewWindow("List App")

	list := widget.NewList(
		func() int {
			return len(listData)
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("template")
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			o.(*widget.Label).SetText(listData[i])
		},
	)

	w.SetContent(list)
	w.SetFixedSize(true)
	w.Resize(fyne.NewSize(200, 200))
	w.ShowAndRun()
}
