package main

import (
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

type GuiData struct {
	Image1_filepath string
	Image2_filepath string

	Processed1_filepath string
	Processed2_filepath string

	Compare_jpg_filepath string
	Compare_gif_filepath string

	Image1_size string
	Image2_size string
}

func NewGuiData(original1, processed1, filesize1, original2, processed2, filesize2, compare_jpg, compare_gif string) GuiData {
	return GuiData{
		Image1_filepath:      original1,
		Image2_filepath:      original2,
		Processed1_filepath:  processed1,
		Processed2_filepath:  processed2,
		Compare_jpg_filepath: compare_jpg,
		Compare_gif_filepath: compare_gif,
		Image1_size:          filesize1,
		Image2_size:          filesize2,
	}
}

type GuiWindow struct {
	window                     fyne.Window
	tabs                       *container.AppTabs
	image1_label               *widget.Label
	image2_label               *widget.Label
	image1_canvas              *canvas.Image
	image2_canvas              *canvas.Image
	compare_jpg_canvas         *canvas.Image
	compare_gif_canvas         *widgetx.AnimatedGif
	compare_animated_container *fyne.Container
	same_button                *widget.Button
	different_button           *widget.Button
	next_button                *widget.Button
	button_container           *fyne.Container
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

	cc := container.NewBorder(nil, nil, nil, nil, animated_gif)

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

	return &GuiWindow{w, tabs,
		image1_fileinfo, image2_fileinfo, image1_canvas, image2_canvas, compare_jpg_canvas,
		animated_gif, cc,
		same_image_button, different_image_button, next_button, button_area}
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
	window.compare_animated_container.RemoveAll()
	new_gif_canvas, err := widgetx.NewAnimatedGif(storage.NewFileURI(data.Compare_gif_filepath))
	if err != nil {
		panic(err)
	}

	window.compare_animated_container.RemoveAll()
	window.compare_animated_container.Add(new_gif_canvas)
	new_gif_canvas.Start()
	window.compare_animated_container.Refresh()

	// Refresh the image area (tab)
	window.tabs.Refresh()

	// Button Handler setup
	window.same_button.OnTapped = func() {

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
