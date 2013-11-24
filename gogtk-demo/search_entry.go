// Entry/Search Entry
//
// GtkEntry allows to display icons and progress information.
// This demo shows how to use these features in a search entry.
package search_entry

import "gobject/gtk-3.0"
import "gobject/gdk-3.0"
import "gobject/gobject-2.0"
import "time"

var dialog *gtk.Dialog
var notebook *gtk.Notebook
var entry *gtk.Entry

const (
	find_button_page = iota
	cancel_button_page
)

const (
	start_search_cmd = iota
	cancel_search_cmd
	stop_daemon_cmd
)

func start_search_daemon(msgqueue chan int) {
	go func() {
		var progress_ticker *time.Ticker
		var finish_ticker *time.Ticker
		var progress_tick <-chan time.Time
		var finish_tick <-chan time.Time

		start_search := func() {
			if progress_ticker != nil {
				println("search in progress, start command ignored")
				return
			}
			progress_ticker = time.NewTicker(100 * time.Millisecond)
			finish_ticker = time.NewTicker(10 * time.Second)
			progress_tick = progress_ticker.C
			finish_tick = finish_ticker.C
			gdk.ThreadsEnter()
			notebook.SetCurrentPage(cancel_button_page)
			gdk.ThreadsLeave()
		}

		stop_search := func() {
			if progress_ticker == nil {
				return
			}

			progress_ticker.Stop()
			finish_ticker.Stop()
			progress_ticker = nil
			progress_tick = nil
			finish_ticker = nil
			finish_tick = nil
		}

		reset_gui := func() {
			gdk.ThreadsEnter()
			entry.SetProgressFraction(0)
			notebook.SetCurrentPage(find_button_page)
			gdk.ThreadsLeave()
		}

		done: for {
			select {
			case cmd := <-msgqueue:
				switch cmd {
				case start_search_cmd:
					start_search()
				case cancel_search_cmd:
					stop_search()
					reset_gui()
				case stop_daemon_cmd:
					stop_search()
					break done
				}
			case <-progress_tick:
				gdk.ThreadsEnter()
				entry.ProgressPulse()
				gdk.ThreadsLeave()
			case <-finish_tick:
				stop_search()
				reset_gui()
			}
		}

		msgqueue <- stop_daemon_cmd
	}()
}

func search_by_name() {
	entry.SetIconFromStock(gtk.EntryIconPositionPrimary, gtk.StockFind)
	entry.SetIconTooltipText(gtk.EntryIconPositionPrimary, "Search by name\n" +
		"Click here to change the search type")
	entry.SetPlaceholderText("name")
}

func search_by_description() {
	entry.SetIconFromStock(gtk.EntryIconPositionPrimary, gtk.StockEdit)
	entry.SetIconTooltipText(gtk.EntryIconPositionPrimary, "Search by description\n" +
		"Click here to change the search type")
	entry.SetPlaceholderText("description")
}

func search_by_file() {
	entry.SetIconFromStock(gtk.EntryIconPositionPrimary, gtk.StockOpen)
	entry.SetIconTooltipText(gtk.EntryIconPositionPrimary, "Search by file name\n" +
		"Click here to change the search type")
	entry.SetPlaceholderText("file name")
}

func create_search_menu() *gtk.Menu {
	menu := gtk.NewMenu()

	item := gtk.NewImageMenuItemWithMnemonic("Search by _name")
	image := gtk.NewImageFromStock(gtk.StockFind, int(gtk.IconSizeMenu))
	item.SetImage(image)
	item.SetAlwaysShowImage(true)
	item.Connect("activate", search_by_name)
	menu.Append(item)

	item = gtk.NewImageMenuItemWithMnemonic("Search by _description")
	image = gtk.NewImageFromStock(gtk.StockEdit, int(gtk.IconSizeMenu))
	item.SetImage(image)
	item.SetAlwaysShowImage(true)
	item.Connect("activate", search_by_description)
	menu.Append(item)

	item = gtk.NewImageMenuItemWithMnemonic("Search by _file name")
	image = gtk.NewImageFromStock(gtk.StockOpen, int(gtk.IconSizeMenu))
	item.SetImage(image)
	item.SetAlwaysShowImage(true)
	item.Connect("activate", search_by_file)
	menu.Append(item)

	menu.ShowAll()
	return menu
}

func entry_populate_popup(entry *gtk.Entry, menu *gtk.Menu) {
	has_text := entry.GetTextLength() > 0

	item := gtk.ToMenuItem(gtk.NewSeparatorMenuItem())
	item.Show()
	menu.Append(item)

	item = gtk.NewMenuItemWithMnemonic("C_lear")
	item.Show()
	item.Connect("activate", func() { entry.SetText("") })
	menu.Append(item)
	item.SetSensitive(has_text)

	search_menu := create_search_menu()
	item = gtk.NewMenuItemWithLabel("Search by")
	item.Show()
	item.SetSubmenu(search_menu)
	menu.Append(item)
}

func Do(mainwin *gtk.Window) *gtk.Window {
	if dialog == nil {
		msgqueue := make(chan int)
		start_search_daemon(msgqueue)

		dialog = gtk.NewDialogWithButtons("Search Entry", mainwin, 0,
			gtk.StockClose, gtk.ResponseTypeNone)
		dialog.SetResizable(false)

		dialog.Connect("response", func() { dialog.Destroy() })
		dialog.Connect("destroy", func() {
			msgqueue <- stop_daemon_cmd
			<-msgqueue
			dialog = nil
			notebook = nil
			entry = nil
		})

		content_area := gtk.ToBox(dialog.GetContentArea())

		vbox := gtk.NewBox(gtk.OrientationVertical, 5)
		content_area.PackStart(vbox, true, true, 0)
		vbox.SetBorderWidth(5)

		label := gtk.NewLabel(gobject.NilString)
		label.SetMarkup("Search entry demo")
		vbox.PackStart(label, false, false, 0)

		hbox := gtk.NewBox(gtk.OrientationHorizontal, 10)
		vbox.PackStart(hbox, true, true, 0)
		hbox.SetBorderWidth(0)

		// Create our entry
		entry = gtk.NewEntry()
		hbox.PackStart(entry, false, false, 0)

		// Create the find and cancel buttons
		notebook = gtk.NewNotebook()
		notebook.SetShowTabs(false)
		notebook.SetShowBorder(false)
		hbox.PackStart(notebook, false, false, 0)

		find_button := gtk.NewButtonWithLabel("Find")
		find_button.Connect("clicked", func() { msgqueue <- start_search_cmd })
		notebook.AppendPage(find_button, nil)
		find_button.Show()

		cancel_button := gtk.NewButtonWithLabel("Cancel")
		cancel_button.Connect("clicked", func() { msgqueue <- cancel_search_cmd })
		notebook.AppendPage(cancel_button, nil)
		cancel_button.Show()

		// Set up the search icon
		search_by_name()

		// Set up the clear icon
		entry.SetIconFromStock(gtk.EntryIconPositionSecondary, gtk.StockClear)
		entry.SetIconSensitive(gtk.EntryIconPositionSecondary, false)
		find_button.SetSensitive(false)

		// Create the menu
		menu := create_search_menu()
		menu.AttachToWidget(entry)

		entry.Connect("icon-press", func(e *gtk.Entry, pos gtk.EntryIconPosition, ev *gdk.EventButton) {
			if pos == gtk.EntryIconPositionPrimary {
				menu.Popup(nil, nil, nil, int(ev.Button), int(ev.Time))
			} else {
				e.SetText("")
			}
		})
		entry.Connect("notify::text", func() {
			has_text := entry.GetTextLength() > 0
			entry.SetIconSensitive(gtk.EntryIconPositionSecondary, has_text)
			find_button.SetSensitive(has_text)
		})
		entry.Connect("activate", func() { msgqueue <- start_search_cmd })
		entry.Connect("populate-popup", entry_populate_popup)

		button := dialog.GetWidgetForResponse(int(gtk.ResponseTypeNone))
		button.GrabFocus()
	}

	if !dialog.GetVisible() {
		dialog.ShowAll()
	} else {
		dialog.Destroy()
	}
	return gtk.ToWindow(dialog)
}