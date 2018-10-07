package main

import (
	"github.com/gotk3/gotk3/gdk"
	"log"

	"github.com/gotk3/gotk3/cairo"
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"

	"github.com/kislenko-artem/gopaint/primitives"
	ln "github.com/kislenko-artem/gopaint/primitives/line"
)

var (
	//line primitives.Line
	objects     []primitives.Primitive
	lineCounter = -1
)

func mainInit(mainWin *gtk.ApplicationWindow) {
	// Преобразуем из объекта именно окно типа gtk.Window
	// и соединяем с сигналом "destroy" чтобы можно было закрыть
	// приложение при закрытии окна
	mainWin.Connect("destroy", func() {
		gtk.MainQuit()
	})

	mainWin.Connect("button-press-event", func(win *gtk.ApplicationWindow, ev *gdk.Event) {
		line := ln.New()
		lineCounter++
		objects = append(objects, line)
		event := &gdk.EventButton{Event: ev}
		objects[lineCounter].SetStart(event.X(), event.Y())
	})

	mainWin.Connect("motion-notify-event", func(win *gtk.ApplicationWindow, ev *gdk.Event) {
		if !objects[lineCounter].IsWait() {
			return
		}
		event := &gdk.EventButton{Event: ev}
		objects[lineCounter].SetStop(event.X(), event.Y())
		win.QueueDraw()
	})

	mainWin.Connect("button-release-event", func(win *gtk.ApplicationWindow, ev *gdk.Event) {
		event := &gdk.EventButton{Event: ev}
		objects[lineCounter].SetStop(event.X(), event.Y())
		objects[lineCounter].Release()
		win.QueueDraw()
	})
	// Set the default window size.
	mainWin.SetDefaultSize(800, 600)

	// Отображаем все виджеты в окне
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

func main() {
	var (
		err error
		b   *gtk.Builder
		obj glib.IObject
	)
	// Инициализируем GTK.
	gtk.Init(nil)

	// Создаём билдер
	if b, err = gtk.BuilderNew(); err != nil {
		log.Fatal("Ошибка:", err)
	}

	// Загружаем в билдер окно из файла Glade
	if err = b.AddFromFile("main.glade"); err != nil {
		log.Fatal("Ошибка:", err)
	}

	// Получаем объект главного окна по ID
	if obj, err = b.GetObject("window_main"); err != nil {
		log.Fatal("Ошибка:", err)
	}

	mainInit(obj.(*gtk.ApplicationWindow))

	// Получаем объект главного окна по ID
	if obj, err = b.GetObject("window_drawing"); err != nil {
		log.Fatal("Ошибка:", err)
	}
	drawWindow(obj.(*gtk.DrawingArea))
	// Выполняем главный цикл GTK (для отрисовки). Он остановится когда
	// выполнится gtk.MainQuit()
	gtk.Main()
}
