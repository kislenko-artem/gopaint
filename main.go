package main

import (
	"github.com/gotk3/gotk3/gdk"
	"log"

	"github.com/gotk3/gotk3/cairo"
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"

	"./primitives"
	ln "./primitives/line"
)

var (
	//line primitives.Line
	lines       []primitives.Line
	lineCounter = -1
)

func init() {
	//line = ln.New()
}

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
		lines = append(lines, line)
		event := &gdk.EventButton{ev}
		lines[lineCounter].SetStart(event.X(), event.Y())
	})

	mainWin.Connect("motion-notify-event", func(win *gtk.ApplicationWindow, ev *gdk.Event) {
		if !lines[lineCounter].IsWait() {
			return
		}
		event := &gdk.EventButton{ev}
		lines[lineCounter].SetEnd(event.X(), event.Y())
		win.QueueDraw()
	})

	mainWin.Connect("button-release-event", func(win *gtk.ApplicationWindow, ev *gdk.Event) {
		event := &gdk.EventButton{ev}
		lines[lineCounter].SetEnd(event.X(), event.Y())
		lines[lineCounter].Release()
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
		cr.SetSourceRGB(0, 0, 0)
		for i := range lines {
			cr.MoveTo(lines[i].GetStart())
			cr.LineTo(lines[i].GetEnd())
		}
		cr.Stroke()
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
