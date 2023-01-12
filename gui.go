package main

import (
	"fmt"
	"log"
	"path/filepath"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/widget"

	_ "image/gif"

	widgetx "fyne.io/x/fyne/widget"
)

var Default_Image_Size = fyne.NewSize(600, 600)

type GuiWindow struct {
	window                     fyne.Window
	tabs                       *container.AppTabs
	image1_label               *widget.Label
	image2_label               *widget.Label
	image1_canvas              *canvas.Image
	image2_canvas              *canvas.Image
	compare_jpg_canvas         *canvas.Image
	compare_gif_canvas         *widgetx.AnimatedGif
	button_restart_animate     *widget.Button
	compare_animated_container *fyne.Container
	same_button                *widget.Button
	different_button           *widget.Button
	next_button                *widget.Button
	button_container           *fyne.Container

	loading_dialog *dialog.ProgressInfiniteDialog
	exit_dialog    dialog.Dialog

	Next_Ontapped      func(perform_update bool) (data GuiData)
	Same_Ontapped      func() GuiData
	Different_Ontapped func() GuiData
	IsLastItem         func() bool

	Exception_Handler func(message string)
	Try_Next_button   *widget.Button // Button when error appear and let user to select next image
}

func NewGuiWindow() *GuiWindow {
	a := app.New()
	w := a.NewWindow("Hello World") // Due to set fixed size has bug, not change this

	// Title
	title := widget.NewLabel("Check Image Different")
	title.Alignment = fyne.TextAlignCenter
	title.TextStyle.Bold = true

	// Original Image 1
	image1_canvas := canvas.NewImageFromFile("resource/loading.jpg")
	image1_canvas.FillMode = canvas.ImageFillContain
	image1_canvas.SetMinSize(Default_Image_Size)

	image1_fileinfo := widget.NewLabel("")
	image1_fileinfo.Alignment = fyne.TextAlignCenter

	// Original Image 2
	image2_canvas := canvas.NewImageFromFile("resource/loading.jpg")
	image2_canvas.FillMode = canvas.ImageFillContain
	image2_canvas.SetMinSize(Default_Image_Size)

	image2_fileinfo := widget.NewLabel(filepath.Base(""))
	image2_fileinfo.Alignment = fyne.TextAlignCenter

	// Compare Image by Gif
	animated_gif, err := widgetx.NewAnimatedGif(storage.NewFileURI("resource/loading.gif"))
	if err != nil {
		panic(err)
	}
	animated_gif.Start()

	// Comapre Image by Jpg
	compare_jpg_canvas := canvas.NewImageFromFile("resource/loading.jpg")
	compare_jpg_canvas.FillMode = canvas.ImageFillContain
	compare_jpg_canvas.SetMinSize(Default_Image_Size)

	// Tabs
	originalTab := container.NewTabItem(
		"Original Image Compare",
		container.NewCenter(container.NewHBox(
			container.NewVBox(image1_fileinfo, image1_canvas),
			container.NewVBox(image2_fileinfo, image2_canvas),
		)),
	)

	comparejpgTab := container.NewTabItem(
		"Difference",
		compare_jpg_canvas,
	)

	button_restart_animate := widget.NewButton("Restart Again", func() {
		fmt.Println("Button Not set function.")
	})

	cc := container.NewBorder(container.NewCenter(button_restart_animate), nil, nil, nil, animated_gif)

	comparegifTab := container.NewTabItem(
		"Animated",
		cc,
	)

	tabs := container.NewAppTabs(
		originalTab, comparejpgTab, comparegifTab,
	)
	tabs.Select(comparejpgTab)

	// Buttons Field
	next_button := widget.NewButton("Next Set", func() {

	})
	next_button.Importance = widget.WarningImportance

	same_image_button := widget.NewButton("They are Same, Delete image with smaller size", func() {

	})
	same_image_button.Importance = widget.DangerImportance

	different_image_button := widget.NewButton("They are Different, Keep Both", func() {

	})

	button_area := container.NewVBox(
		container.NewCenter(
			container.NewVBox(
				same_image_button,
				different_image_button,
			)),
		container.NewBorder(nil, nil, nil, next_button),
	)

	// Layout
	content := container.NewVBox(
		title,
		tabs,
		widget.NewLabel(""), //Act as seperate space
		button_area,
	)

	// Key Binding
	w.Canvas().SetOnTypedKey(func(ke *fyne.KeyEvent) {
		if ke.Physical.ScanCode == 331 { // Lef arrow key pressed
			current_index := tabs.SelectedIndex()
			if current_index == 0 {
				current_index = len(tabs.Items) - 1
			} else {
				current_index--
			}
			//fmt.Println(current_index)
			tabs.SelectIndex(current_index)
			return
		}

		if ke.Physical.ScanCode == 333 { // Right Arrow Key Pressed
			current_index := tabs.SelectedIndex()
			if current_index == len(tabs.Items)-1 {
				current_index = 0
			} else {
				current_index++
			}
			tabs.SelectIndex(current_index)
			return
		}
	})

	// Loading Dialog
	loading_dialog := dialog.NewProgressInfinite("Loading", "Please wait until next image is loaded", w)
	loading_dialog.Hide()

	exit_dialog := dialog.NewInformation("Exiting", "That is Last Image. Program will terminate.", w)
	exit_dialog.Hide()

	exit_dialog.SetOnClosed(func() {
		w.Close()
	})

	try_next_button := widget.NewButton("Try Next", func() {

	})
	try_next_button.Importance = widget.WarningImportance

	// Set up error when exception
	exception_handler := func(message string) {
		fmt.Println("Critial Error Occur.")
		w.SetContent(container.NewCenter(
			container.NewVBox(
				container.NewCenter(widget.NewLabel("Error Occur.\n"+message)),
				container.NewCenter(container.NewHBox(
					widget.NewButton("Copy Error Message", func() {
						w.Clipboard().SetContent(message)
					}),
					widget.NewButton("Close", func() {
						w.Close()
					}),
					try_next_button,
				)),
			),
		))
		w.Resize(fyne.NewSize(200, 200))
		w.CenterOnScreen()
	}

	w.SetContent(content)
	w.CenterOnScreen()

	empty_data := func() GuiData { // Dummy function
		return NewGuiData("", "", "", "", "", "", "", "")
	}

	empty_data2 := func(bool) GuiData { // Dummy function
		return NewGuiData("", "", "", "", "", "", "", "")
	}

	empty_bool := func() bool { // Dummy function
		return false
	}

	return &GuiWindow{w, tabs,
		image1_fileinfo, image2_fileinfo, image1_canvas, image2_canvas, compare_jpg_canvas, animated_gif, button_restart_animate, cc,
		same_image_button, different_image_button, next_button, button_area,
		loading_dialog, exit_dialog,
		empty_data2, empty_data, empty_data, empty_bool,
		exception_handler, try_next_button}
}

// Update the gui by GuiData, also set up handler and refresh container
func (window *GuiWindow) Update(data GuiData) {
	window.image1_canvas.File = data.Image1_filepath

	window.image2_canvas.File = data.Image2_filepath

	window.image1_label.Text = filepath.Base(data.Image1_filepath) + "\n" + data.Image1_size
	window.image1_label.Refresh()
	window.image2_label.Text = filepath.Base(data.Image2_filepath) + "\n" + data.Image2_size
	window.image2_label.Refresh()

	window.compare_jpg_canvas.File = data.Compare_jpg_filepath

	// Gif image reset
	window.compare_animated_container.Remove(window.compare_gif_canvas) // Remove the gif only
	new_gif_canvas, err := widgetx.NewAnimatedGif(storage.NewFileURI(data.Compare_gif_filepath))

	window.compare_animated_container.Add(new_gif_canvas)
	if err == nil {
		new_gif_canvas.Start() //Only Start animation if not error
	} else {
		log.Printf("Error: %v, target: %s\n", err, data)
		original_content := window.window.Content()
		window.Try_Next_button.OnTapped = func() {
			window.loading_dialog.Show()
			window.window.SetContent(original_content) // Recover the window structure
			data := window.Different_Ontapped()        // Use Different to not remain record to skipped item
			window.Update(data)
			window.loading_dialog.Hide()
		}
		window.Exception_Handler(
			fmt.Sprintf("%v.\n\nImage Error: %s\n\nImage Original:\n%s, \n%s\n", err, data.Compare_gif_filepath, data.Image1_filepath, data.Image2_filepath))

	}
	window.compare_gif_canvas = new_gif_canvas

	window.button_restart_animate.OnTapped = func() {
		fmt.Println("Button restart pressed.")
		new_gif_canvas.Start()
	}

	window.compare_animated_container.Refresh()

	// Button Handler setup
	window.same_button.OnTapped = func() {
		confirmCallback := func(confirm bool) {
			if confirm {
				window.loading_dialog.Show()
				data := window.Same_Ontapped()
				window.Update(data)
				window.loading_dialog.Hide()
			}
		}

		cnf := dialog.NewConfirm("Confirmation", "Are you sure to delete file: "+data.Image2_filepath+" ?", confirmCallback, window.window) // Image2 must be the image with smaller size
		cnf.SetDismissText("No")
		cnf.SetConfirmText("Yes")
		cnf.Show()
	}

	window.different_button.OnTapped = func() {
		window.loading_dialog.Show()
		data := window.Different_Ontapped()
		window.Update(data)
		window.loading_dialog.Hide()
	}

	window.next_button.OnTapped = func() {
		window.loading_dialog.Show()
		data := window.Next_Ontapped(true)

		window.Update(data)
		window.loading_dialog.Hide()
	}

	// Prevent Last Item Error
	if window.IsLastItem() {
		window.next_button.OnTapped = func() {
			window.Next_Ontapped(false)
			window.exit_dialog.Show()
		}
		window.next_button.SetText("Close Program")

		window.different_button.OnTapped = func() { //Overwirte orignial function to prevent error
			window.exit_dialog.Show()
		}

		window.same_button.OnTapped = func() {
			confirmCallback := func(confirm bool) {
				if confirm {
					window.exit_dialog.Show()
				}
			}

			cnf := dialog.NewConfirm("Confirmation", "Are you sure to delete file: "+data.Image2_filepath+" ?", confirmCallback, window.window) // Image2 must be the image with smaller size
			cnf.SetDismissText("No")
			cnf.SetConfirmText("Yes")
			cnf.Show()
		}
	}

	// Refresh the image area (tab) and button area
	window.tabs.Refresh()
	window.button_container.Refresh()
}

func (window *GuiWindow) ShowAndRun() {
	window.window.ShowAndRun()
}
