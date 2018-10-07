package main

import (
	"github.com/gotk3/gotk3/gdk"
	"log"

	"github.com/gotk3/gotk3/cairo"
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"

	"github.com/kislenko-artem/gopaint/primitives"
	ln "github.com/kislenko-artem/gopaint/primitives/line"
	"github.com/kislenko-artem/gopaint/property/color"
)

var (
	baseColor   color.Color
	objects     []primitives.Primitive
	lineCounter = -1
)

func mainInit(mainWin *gtk.ApplicationWindow) {

	mainWin.Connect("destroy", func() {
		gtk.MainQuit()
	})

	mainWin.Connect("button-press-event", func(win *gtk.ApplicationWindow, ev *gdk.Event) {
		line := ln.New(baseColor)
		lineCounter++
		objects = append(objects, line)
		event := &gdk.EventButton{Event: ev}
		objects[lineCounter].SetStart(event.X(), event.Y())
	})

	mainWin.Connect("motion-notify-event", func(win *gtk.ApplicationWindow, ev *gdk.Event) {
		if lineCounter < 0 {
			return
		}
		if !objects[lineCounter].IsWait() {
			return
		}
		event := &gdk.EventButton{Event: ev}
		objects[lineCounter].SetStop(event.X(), event.Y())
		win.QueueDraw()
	})

	mainWin.Connect("button-release-event", func(win *gtk.ApplicationWindow, ev *gdk.Event) {
		if lineCounter < 0 {
			return
		}
		event := &gdk.EventButton{Event: ev}
		objects[lineCounter].SetStop(event.X(), event.Y())
		objects[lineCounter].Release()
		win.QueueDraw()
	})

	mainWin.SetDefaultSize(800, 600)

	mainWin.ShowAll()
}

func drawWindow(drawWindow *gtk.DrawingArea) {
	drawWindow.Connect("draw", func(da *gtk.DrawingArea, cr *cairo.Context) {
		if lineCounter < 0 {
			return
		}
		for i := range objects {
			objects[i].SetColor(cr)
			objects[i].Draw(cr)
		}
	})

}

func colorPicker(picker *gtk.ColorButton) {
	picker.Connect("color-set", func(obj *glib.Object) {
		rgb, err := obj.GetProperty("rgba")
		if err != nil {
			log.Fatal(err)
		}
		values := rgb.(*gdk.RGBA).Floats()
		baseColor.RGB.R = values[0]
		baseColor.RGB.G = values[1]
		baseColor.RGB.B = values[2]
	})
}

func main() {
	var (
		err error
		b   *gtk.Builder
		obj glib.IObject
	)

	gtk.Init(nil)

	if b, err = gtk.BuilderNew(); err != nil {
		log.Fatal("Ошибка:", err)
	}

	if err = b.AddFromFile("main.glade"); err != nil {
		log.Fatal("Ошибка:", err)
	}

	if obj, err = b.GetObject("window_main"); err != nil {
		log.Fatal("Ошибка:", err)
	}
	mainInit(obj.(*gtk.ApplicationWindow))

	if obj, err = b.GetObject("window_drawing"); err != nil {
		log.Fatal("Ошибка:", err)
	}
	drawWindow(obj.(*gtk.DrawingArea))

	if obj, err = b.GetObject("color_picker"); err != nil {
		log.Fatal("Ошибка:", err)
	}
	colorPicker(obj.(*gtk.ColorButton))

	gtk.Main()
}
