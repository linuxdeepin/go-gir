// Menus
//
// There are several widgets involved in displaying menus. The
// GtkMenuBar widget is a menu bar, which normally appears horizontally
// at the top of an application, but can also be layed out vertically.
// The GtkMenu widget is the actual menu that pops up. Both GtkMenuBar
// and GtkMenu are subclasses of GtkMenuShell; a GtkMenuShell contains
// menu items (GtkMenuItem). Each menu item contains text and/or images
// and can be selected by the user.
//
// There are several kinds of menu item, including plain GtkMenuItem,
// GtkCheckMenuItem which can be checked/unchecked, GtkRadioMenuItem
// which is a check menu item that's in a mutually exclusive group,
// GtkSeparatorMenuItem which is a separator bar, GtkTearoffMenuItem
// which allows a GtkMenu to be torn off, and GtkImageMenuItem which
// can place a GtkImage or other widget next to the menu text.
//
// A GtkMenuItem can have a submenu, which is simply a GtkMenu to pop
// up when the menu item is selected. Typically, all menu items in a menu bar
// have submenus.
//
// GtkUIManager provides a higher-level interface for creating menu bars
// and menus; while you can construct menus manually, most people don't
// do that. There's a separate demo for GtkUIManager.
package menus

import "gobject/gtk-3.0"
import "fmt"

func create_menu(depth int, tearoff bool) *gtk.Menu {
	if depth < 1 {
		return nil
	}

	menu := gtk.NewMenu()
	if tearoff {
		menuitem := gtk.NewTearoffMenuItem()
		menu.Append(menuitem)
		menuitem.Show()
	}

//	var group *gtk.RadioMenuItem

	for i, j := 0, 1; i < 5; i, j = i+1, j+1 {
		label := fmt.Sprintf("item %2d - %d", depth, j)
		// TODO: groups
		var menuitem *gtk.RadioMenuItem
		menuitem = gtk.NewRadioMenuItemWithLabel(nil, label)
//		if i == 0 {
//			group = menuitem
//		}

		menu.Append(menuitem)
		menuitem.Show()

		if i == 3 {
			menuitem.SetSensitive(false)
		}

		menuitem.SetSubmenu(create_menu(depth-1, true))
	}

	return menu
}

func change_orientation(button *gtk.Widget, menubar *gtk.MenuBar) {
	parent := gtk.ToOrientable(menubar.GetParent())
	orientation := parent.GetOrientation()
	if orientation == gtk.OrientationVertical {
		menubar.SetProperty("pack-direction", gtk.PackDirectionTtb)
		parent.SetOrientation(gtk.OrientationHorizontal)
	} else {
		menubar.SetProperty("pack-direction", gtk.PackDirectionLtr)
		parent.SetOrientation(gtk.OrientationVertical)
	}
}

var window *gtk.Window

func Do(mainwin *gtk.Window) *gtk.Window {
	if window == nil {
		window = gtk.NewWindow(gtk.WindowTypeToplevel)
		window.SetScreen(mainwin.GetScreen())
		window.SetTitle("Menus")
		window.Connect("destroy", func() { window = nil })

		accel_group := gtk.NewAccelGroup()
		window.AddAccelGroup(accel_group)

		window.SetBorderWidth(0)

		box := gtk.NewBox(gtk.OrientationHorizontal, 0)
		window.Add(box)
		box.Show()

		box1 := gtk.NewBox(gtk.OrientationVertical, 0)
		box.Add(box1)
		box1.Show()

		menubar := gtk.NewMenuBar()
		box1.PackStart(menubar, false, true, 0)
		menubar.Show()

		menu := create_menu(2, true)

		menuitem := gtk.NewMenuItemWithLabel("test\nline2")
		menuitem.SetSubmenu(menu)
		menubar.Append(menuitem)
		menuitem.Show()

		menuitem = gtk.NewMenuItemWithLabel("foo")
		menuitem.SetSubmenu(create_menu(3, true))
		menubar.Append(menuitem)
		menuitem.Show()

		menuitem = gtk.NewMenuItemWithLabel("bar")
		menuitem.SetSubmenu(create_menu(4, true))
		menubar.Append(menuitem)
		menuitem.Show()

		box2 := gtk.NewBox(gtk.OrientationVertical, 10)
		box2.SetBorderWidth(10)
		box1.PackStart(box2, false, true, 0)
		box2.Show()

		button := gtk.NewButtonWithLabel("Flip")
		button.Connect("clicked", func(button *gtk.Widget) {
			change_orientation(button, menubar)
		})
		box2.PackStart(button, true, true, 0)
		button.Show()

		button = gtk.NewButtonWithLabel("Close")
		button.Connect("clicked", func() { window.Destroy() })
		box2.PackStart(button, true, true, 0)
		button.SetCanDefault(true)
		button.GrabDefault()
		button.Show()
	}

	if !window.GetVisible() {
		window.ShowAll()
	} else {
		window.Destroy()
	}
	return window
}