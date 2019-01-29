// initialize main window

package cmd

import (
	"github.com/gotk3/gotk3/cairo"
	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
	"log"

	"github.com/kislenko-artem/gopaint/primitives"
	ln "github.com/kislenko-artem/gopaint/primitives/line"
	"github.com/kislenko-artem/gopaint/primitives/pencil"
	"github.com/kislenko-artem/gopaint/property/color"
)

func New() *Window {
	w := Window{}
	return w.Create()
}

type Window struct {
	baseColor      color.Color
	objects        []primitives.Primitive
	primitivesList []glib.IObject
	objCounter     int
	lastType       string
}

func (w *Window) Create() *Window {
	w.objCounter = -1
	return w
}

// create instance for primitive and start using it
func (w *Window) addObject(objType string) {
	w.lastType = objType
	var newObject primitives.Primitive
	switch objType {
	case "line":
		newObject = ln.New(w.baseColor)
	case "pencil":
		newObject = pencil.New(w.baseColor)
	default:
		return
	}
	w.objCounter++
	w.objects = append(w.objects, newObject)
}

// events and functions to drawing
func (w *Window) drawWindow(b *gtk.Builder) {
	var (
		err error
		obj glib.IObject
		ok  bool
		drawWindow *gtk.DrawingArea
	)

	if obj, err = b.GetObject("window_drawing"); err != nil {
		log.Fatal("Ошибка:", err)
	}

	if drawWindow, ok = obj.(*gtk.DrawingArea); ok == false {
		return
	}

	drawWindow.Connect("draw", func(da *gtk.DrawingArea, cr *cairo.Context) {
		for i := range w.objects {
			w.objects[i].SetColor(cr)
			w.objects[i].Draw(cr)
		}
	})

}

// events and functions to pick colors
func (w *Window) colorPicker(b *gtk.Builder) {
	var (
		err error
		obj glib.IObject
		ok  bool
		picker *gtk.ColorButton
	)

	if obj, err = b.GetObject("color_picker"); err != nil {
		log.Fatal("Ошибка:", err)
	}

	if picker, ok = obj.(*gtk.ColorButton); ok == false {
		return
	}

	picker.Connect("color-set", func(obj *glib.Object) {
		rgb, err := obj.GetProperty("rgba")
		if err != nil {
			log.Fatal(err)
		}
		values := rgb.(*gdk.RGBA).Floats()
		w.baseColor.RGB.R = values[0]
		w.baseColor.RGB.G = values[1]
		w.baseColor.RGB.B = values[2]
		w.addObject(w.lastType)
	})
}


// create events for drawing
func (w *Window) drawingInit(b *gtk.Builder) {
	var (
		err error
		obj glib.IObject
		ok  bool
		mainWin *gtk.ApplicationWindow
	)
	if obj, err = b.GetObject("window_main"); err != nil {
		log.Fatal("Ошибка:", err)
	}
	if mainWin, ok = obj.(*gtk.ApplicationWindow); ok == false {
		return
	}

	mainWin.Connect("destroy", func() {
		gtk.MainQuit()
	})

	mainWin.Connect("button-press-event", func(win *gtk.ApplicationWindow, ev *gdk.Event) {
		if w.objCounter < 0 {
			return
		}
		event := &gdk.EventButton{Event: ev}
		w.objects[w.objCounter].SetStart(event.X(), event.Y())
	})

	mainWin.Connect("motion-notify-event", func(win *gtk.ApplicationWindow, ev *gdk.Event) {
		if w.objCounter < 0 {
			return
		}
		if !w.objects[w.objCounter].IsWait() {
			return
		}
		event := &gdk.EventButton{Event: ev}
		w.objects[w.objCounter].SetStop(event.X(), event.Y())
		win.QueueDraw()
	})

	mainWin.Connect("button-release-event", func(win *gtk.ApplicationWindow, ev *gdk.Event) {
		if w.objCounter < 0 {
			return
		}
		event := &gdk.EventButton{Event: ev}
		w.objects[w.objCounter].SetStop(event.X(), event.Y())
		w.objects[w.objCounter].Release()
		w.addObject(w.lastType)
		win.QueueDraw()
	})

	mainWin.SetDefaultSize(800, 600)

	mainWin.ShowAll()
}

// create events for buttons
func (w *Window) primitivesBtnInit(b *gtk.Builder) {
	var (
		err error
	)
	w.primitivesList = make([]glib.IObject, 2)
	if w.primitivesList[0], err = b.GetObject("pencil_btn"); err != nil {
		log.Fatal("Ошибка:", err)
	}

	if w.primitivesList[1], err = b.GetObject("line_btn"); err != nil {
		log.Fatal("Ошибка:", err)
	}

	w.primitivesList[0].(*gtk.Button).Connect("button-press-event", func(btn *gtk.Button, ev *gdk.Event) {
		w.addObject("pencil")
	})
	w.primitivesList[1].(*gtk.Button).Connect("button-press-event", func(btn *gtk.Button, ev *gdk.Event) {
		w.addObject("line")
	})
}

// GtkInit initialize all GTK object and default values for starting application
func (w *Window) GtkInit() {
	var (
		err error
		b   *gtk.Builder
	)
	gtk.Init(nil)
	if b, err = gtk.BuilderNew(); err != nil {
		log.Fatal("Ошибка:", err)
	}
	if err = b.AddFromFile("main.glade"); err != nil {
		log.Fatal("Ошибка:", err)
	}

	w.drawingInit(b)
	w.drawWindow(b)
	w.colorPicker(b)
	w.primitivesBtnInit(b)

	gtk.Main()
}
