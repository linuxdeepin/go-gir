// Assistant
//
// Demonstrates a sample multi-step assistant. Assistants are used to divide
// an operation into several simpler sequential steps, and to guide the user
// through these steps.
package assistant

import "gobject/gdk-3.0"
import "gobject/gtk-3.0"
import "fmt"
import "time"

var assistant *gtk.Assistant
var progress_bar *gtk.ProgressBar

func on_assistant_apply() {
	// Start a timer to simulate changes taking a few seconds to apply.
	go func() {
		for {
			gdk.ThreadsEnter()

			// Work, work, work...
			fraction := progress_bar.GetFraction()
			fraction += 0.05

			if fraction < 1.0 {
				progress_bar.SetFraction(fraction)
			} else {
				// Close automatically once changes are fully applied.
				assistant.Destroy()
				gdk.ThreadsLeave()
				return
			}
			gdk.ThreadsLeave()

			time.Sleep(100 * time.Millisecond)
		}
	}()
}

func on_assistant_prepare() {
	n_pages := assistant.GetNPages()
	page_n := assistant.GetCurrentPage()

	title := fmt.Sprintf("Sample assistant (%d of %d)", page_n+1, n_pages)
	assistant.SetTitle(title)

	// The fourth page (counting from zero) is the progress page.  The
	// user clicked Apply to get here so we tell the assistant to commit,
	// which means the changes up to this point are permanent and cannot
	// be cancelled or revisited.
	if page_n == 3 {
		assistant.Commit()
	}
}

func create_page1(assistant *gtk.Assistant) {
	box := gtk.NewBox(gtk.OrientationHorizontal, 12)
	box.SetBorderWidth(12)

	label := gtk.NewLabel("You must fill out this entry to continue:")
	box.PackStart(label, false, false, 0)

	entry := gtk.NewEntry()
	entry.SetActivatesDefault(true)
	box.PackStart(entry, true, true, 0)
	entry.Connect("changed", func() {
		page_n := assistant.GetCurrentPage()
		cur_page := assistant.GetNthPage(page_n)
		if entry.GetTextLength() > 0 {
			assistant.SetPageComplete(cur_page, true)
		} else {
			assistant.SetPageComplete(cur_page, false)
		}
	})

	box.ShowAll()
	assistant.AppendPage(box)
	assistant.SetPageTitle(box, "Page 1")
	assistant.SetPageType(box, gtk.AssistantPageTypeIntro)
}

func create_page2(assistant *gtk.Assistant) {
	box := gtk.NewBox(gtk.OrientationVertical, 12)
	box.SetBorderWidth(12)

	checkbutton := gtk.NewCheckButtonWithLabel("This is optional data, you may continue " +
		"even if you do not check this")
	box.PackStart(checkbutton, false, false, 0)

	box.ShowAll()
	assistant.AppendPage(box)
	assistant.SetPageComplete(box, true)
	assistant.SetPageTitle(box, "Page 2")
}

func create_page3(assistant *gtk.Assistant) {
	label := gtk.NewLabel("This is confirmation page, press 'Apply' to apply changes")

	label.Show()
	assistant.AppendPage(label)
	assistant.SetPageType(label, gtk.AssistantPageTypeConfirm)
	assistant.SetPageComplete(label, true)
	assistant.SetPageTitle(label, "Confirmation")
}

func create_page4(assistant *gtk.Assistant) {
	progress_bar = gtk.NewProgressBar()
	progress_bar.SetHAlign(gtk.AlignCenter)
	progress_bar.SetVAlign(gtk.AlignCenter)

	progress_bar.Show()
	assistant.AppendPage(progress_bar)
	assistant.SetPageType(progress_bar, gtk.AssistantPageTypeProgress)
	assistant.SetPageTitle(progress_bar, "Applying changes")

	// This prevents the assistant window from being
	// closed while we're "busy" applying changes.
	assistant.SetPageComplete(progress_bar, false)
}

func Do(mainwin *gtk.Window) *gtk.Window {
	if assistant == nil {
		assistant = gtk.NewAssistant()
		assistant.SetDefaultSize(-1, 300)
		assistant.SetScreen(mainwin.GetScreen())
		assistant.Connect("destroy", func() {
			assistant = nil
			progress_bar = nil
		})

		create_page1(assistant)
		create_page2(assistant)
		create_page3(assistant)
		create_page4(assistant)

		close_cancel := func() {
			assistant.Destroy()
		}
		assistant.Connect("cancel", close_cancel)
		assistant.Connect("close", close_cancel)
		assistant.Connect("apply", on_assistant_apply)
		assistant.Connect("prepare", on_assistant_prepare)
	}

	if !assistant.GetVisible() {
		assistant.ShowAll()
	} else {
		assistant.Destroy()
	}
	return gtk.ToWindow(assistant)
}
