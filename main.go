package main

import (
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"
)

func Test_Set(set_number int) (string, string) {
	switch set_number {
	case 1:
		return "test_set\\blue gk_1672650631.png", "test_set\\blue gk_1672650653.png"
	case 2:
		return "test_set\\an yasuri_1672722083_faceleft.png", "test_set\\an yasuri_1672922692_faceright.jpg"
	case 3:
		return "test_set\\akie (44265104)_1672682033.png", "test_set\\akie (44265104)_1673037743.jpg"
	}
	return "", ""

}

func main() {
	fmt.Println("Hello world!")

	queue := NewImagesQueueByFile("similar_data.txt")
	skipped_item := NewImagesQueue()

	for i := 1; i <= 3; i++ {
		path1, path2 := Test_Set(i)
		queue.Add(path1, path2)
	}

	if queue.Empty() {
		os.Exit(135)
	}

	// Diff checker Init
	diff_check := NewImageDiff()
	diff_check.SetDiffImageDir("Temp")
	err := diff_check.Init()
	if err != nil {
		fmt.Println("Image Magick is not installed.", err)
		os.Exit(1)
	}
	defer diff_check.Terminate()

	// Start Gui with first item in queue
	current_index := 0
	diff_check.SetImages(queue.Get(current_index))
	diff_check.Diff()

	window := NewGuiWindow()
	window.Update(NewGuiData(diff_check.GetInfo()))

	update := func() (data GuiData) {
		current_index++
		fmt.Println("Index: ", current_index)

		diff_check.SetImages(queue.Get(current_index))
		diff_check.Diff()

		// v1, v2, v3, v4, v5, v6, v7, v8 := diff_check.GetInfo()
		// fmt.Println("Update Callback 1: {", v1, v2, v3, v4, v5, v6, v7, v8, "}")
		// fmt.Printf("Update Callback2 : %v\n", NewGuiData(diff_check.GetInfo()))

		return NewGuiData(diff_check.GetInfo())
	}

	window.Next_Ontapped = func(perform_update bool) (data GuiData) {
		skipped_item.Add(queue.Get(current_index))
		if perform_update {
			data = update()
			return data
		} else {
			return data
		}
	}

	window.Same_Ontapped = func() (data GuiData) {
		_, path2 := queue.Get(current_index)
		//os.Remove(path1)
		os.Remove(path2)

		data = update()
		return data
	}

	window.Different_Ontapped = func() (data GuiData) {
		data = update()
		// fmt.Printf("Different Callback : %s\n", data)
		return data
	}

	window.IsLastItem = func() bool {
		return current_index >= queue.Capacity()-1
	}

	window.ShowAndRun()

	// Cleanup
	skipped_item.Concat(queue, current_index)

	fmt.Println("Skipped:", skipped_item)

	data_file, data_err := os.OpenFile("similar_data.txt", os.O_RDWR|os.O_CREATE, 0755)
	if data_err != nil {
		log.Errorf("[OutputReport] ERROR: Cannot Open Data File.\n")
		log.Errorf("[OutputReport] ------ Error   : %s\n", data_err)
	}
	defer data_file.Close()

	for i := 0; i < skipped_item.Capacity(); i++ {
		data_file.WriteString(fmt.Sprintf("%s ??? %s\n", skipped_item.image1_path[i], skipped_item.image2_path[i]))
	}
}
