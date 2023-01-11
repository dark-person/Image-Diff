package main

import (
	"fmt"
	"path/filepath"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
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

	Next_Set func() (GuiData, bool) // function for setting next button
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

	w.SetContent(content)
	w.CenterOnScreen()

	empty := func() (GuiData, bool) { // Dummy function
		return NewGuiData("", "", "", "", "", "", "", ""), true
	}

	return &GuiWindow{w, tabs,
		image1_fileinfo, image2_fileinfo, image1_canvas, image2_canvas, compare_jpg_canvas, animated_gif, button_restart_animate, cc,
		same_image_button, different_image_button, next_button, button_area, empty}
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
	if err != nil {
		panic(err)
	}

	window.compare_animated_container.Add(new_gif_canvas)
	new_gif_canvas.Start()
	window.compare_gif_canvas = new_gif_canvas

	window.button_restart_animate.OnTapped = func() {
		fmt.Println("Button restart pressed.")
		new_gif_canvas.Start()
	}

	//window.button_restart_animate.Hide()

	window.compare_animated_container.Refresh()

	// Refresh the image area (tab)
	window.tabs.Refresh()

	// Button Handler setup
	window.same_button.OnTapped = func() {

	}

	window.next_button.OnTapped = func() {
		//fmt.Println(window.Next_Set())
		data, last := window.Next_Set()
		if last {
			window.next_button.Disable()
		}
		window.Update(data)
	}
}

func (window *GuiWindow) ShowAndRun() {
	window.window.ShowAndRun()
}

func WindowConstruct(data GuiData) fyne.Window {

	image1 := data.Image1_filepath
	image2 := data.Image2_filepath
	compare_jpg := data.Compare_jpg_filepath
	compare_gif := data.Compare_gif_filepath

	a := app.New()
	w := a.NewWindow("Hello World") // Due to set fixed size has bug, not change this

	// Title
	title := widget.NewLabel("Check Image Different")
	title.Alignment = fyne.TextAlignCenter
	title.TextStyle.Bold = true

	// Original Image 1
	image1_canvas := canvas.NewImageFromFile(image1)
	image1_canvas.FillMode = canvas.ImageFillContain
	image1_canvas.SetMinSize(Default_Image_Size)

	image1_filename := widget.NewLabel(filepath.Base(image1))
	image1_filename.Alignment = fyne.TextAlignCenter

	// Original Image 2
	image2_canvas := canvas.NewImageFromFile(image2)
	image2_canvas.FillMode = canvas.ImageFillContain
	image2_canvas.SetMinSize(Default_Image_Size)

	image2_filename := widget.NewLabel(filepath.Base(image2))
	image2_filename.Alignment = fyne.TextAlignCenter

	// Compare Image by Gif
	animated_gif, err := widgetx.NewAnimatedGif(storage.NewFileURI(compare_gif))
	animated_gif.Start()
	if err != nil {
		panic(err)
	}

	// Comapre Image by Jpg
	compare_jpg_canvas := canvas.NewImageFromFile(compare_jpg)
	compare_jpg_canvas.FillMode = canvas.ImageFillContain
	compare_jpg_canvas.SetMinSize(Default_Image_Size)

	// Tabs
	originalTab := container.NewTabItem(
		"Original Image Compare",
		container.NewCenter(container.NewHBox(
			container.NewVBox(image1_filename, image1_canvas),
			container.NewVBox(image2_filename, image2_canvas),
		)),
	)

	comparejpgTab := container.NewTabItem(
		"Difference",
		compare_jpg_canvas,
	)

	comparegifTab := container.NewTabItem(
		"Animated",
		animated_gif,
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

	// Layout
	content := container.NewVBox(
		title,
		tabs,
		container.NewCenter(
			container.NewVBox(
				same_image_button,
				different_image_button,
			)),
		container.NewBorder(nil, nil, nil, next_button),
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

	w.SetContent(content)
	w.ShowAndRun()

	return w
}
