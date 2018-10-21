package main

import (
	"github.com/gotk3/gotk3/gdk"
	"log"

	"github.com/gotk3/gotk3/cairo"
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"

	"github.com/kislenko-artem/gopaint/primitives"
	ln "github.com/kislenko-artem/gopaint/primitives/line"
	"github.com/kislenko-artem/gopaint/primitives/pencil"
	"github.com/kislenko-artem/gopaint/property/color"
)

var (
	baseColor      color.Color
	objects        []primitives.Primitive
	primitivesList []glib.IObject
	objCounter     = -1
)

func mainInit(mainWin *gtk.ApplicationWindow) {

	mainWin.Connect("destroy", func() {
		gtk.MainQuit()
	})

	mainWin.Connect("button-press-event", func(win *gtk.ApplicationWindow, ev *gdk.Event) {
		log.Println(objCounter, len(objects))
		if objCounter < 0 {
			return
		}
		event := &gdk.EventButton{Event: ev}
		objects[objCounter].SetStart(event.X(), event.Y())
	})

	mainWin.Connect("motion-notify-event", func(win *gtk.ApplicationWindow, ev *gdk.Event) {
		if objCounter < 0 {
			return
		}
		if !objects[objCounter].IsWait() {
			return
		}
		event := &gdk.EventButton{Event: ev}
		objects[objCounter].SetStop(event.X(), event.Y())
		win.QueueDraw()
	})

	mainWin.Connect("button-release-event", func(win *gtk.ApplicationWindow, ev *gdk.Event) {
		if objCounter < 0 {
			return
		}
		event := &gdk.EventButton{Event: ev}
		objects[objCounter].SetStop(event.X(), event.Y())
		objects[objCounter].Release()
		win.QueueDraw()
	})

	mainWin.SetDefaultSize(800, 600)

	mainWin.ShowAll()
}

func drawWindow(drawWindow *gtk.DrawingArea) {
	drawWindow.Connect("draw", func(da *gtk.DrawingArea, cr *cairo.Context) {
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

func listenCheckingObjects(glibs []glib.IObject, IDs []string) {

	glibs[0].(*gtk.Button).Connect("button-press-event", func(btn *gtk.Button, ev *gdk.Event) {
		objCounter++
		log.Println("pencil_btn")
		pencil := pencil.New(baseColor)
		pencil.RGB = baseColor.RGB
		objects = append(objects, pencil)
	})
	glibs[1].(*gtk.Button).Connect("button-press-event", func(btn *gtk.Button, ev *gdk.Event) {
		objCounter++
		log.Println("line_btn")
		line := ln.New(baseColor)
		line.RGB = baseColor.RGB
		objects = append(objects, line)
	})

}

func main() {
	var (
		err error
		b   *gtk.Builder
		obj glib.IObject
	)
	primitivesList = make([]glib.IObject, 2)

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

	if primitivesList[0], err = b.GetObject("pencil_btn"); err != nil {
		log.Fatal("Ошибка:", err)
	}

	if primitivesList[1], err = b.GetObject("line_btn"); err != nil {
		log.Fatal("Ошибка:", err)
	}

	listenCheckingObjects(primitivesList, []string{"pencil_btn", "line_btn"})
	gtk.Main()
}
