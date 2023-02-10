package main

import (
	"fmt"
	"os"
	"strings"

	log "github.com/sirupsen/logrus"
)

func main() {
	logger := prepare_logger("temp.log")
	defer logger_close(logger)

	fmt.Println("Hello world!")

	queue := NewImagesQueueByFile("similar_data.txt")
	skipped_item := NewImagesQueue()

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
	window.Update(NewGuiData(diff_check.GetInfo()), queue.Capacity())

	update := func() (data GuiData, remain int) {
		retry_count := diff_check.ClearTempFile()
		current_index++
		log.Debugf("Index: %d, Retry Pending: %d\n", current_index, retry_count)

		// Prevent Index out of range
		if current_index > queue.Capacity()-1 {
			return NewGuiData("", "", "", "", "", "", "", ""), -1
		}

		diff_check.SetImages(queue.Get(current_index))
		diff_check.Diff()

		// v1, v2, v3, v4, v5, v6, v7, v8 := diff_check.GetInfo()
		// fmt.Println("Update Callback 1: {", v1, v2, v3, v4, v5, v6, v7, v8, "}")
		// fmt.Printf("Update Callback2 : %v\n", NewGuiData(diff_check.GetInfo()))

		return NewGuiData(diff_check.GetInfo()), queue.Capacity() - current_index
	}

	window.Next_Ontapped = func(perform_update bool) (data GuiData, remain int) {
		skipped_item.Add(queue.Get(current_index))
		if perform_update {
			data, remain = update()
			return data, remain
		} else {
			return data, queue.Capacity() - current_index
		}
	}

	window.Same_Ontapped = func() (data GuiData, remain int) {
		_, path2 := queue.Get(current_index)
		err := os.Remove(path2)
		if err != nil {
			log.Warnf("Error when removing path2: %v\n", err)
			if !strings.Contains(err.Error(), "cannot find the file specified") {
				diff_check.AddRetryQueue(path2)
			}
		}

		data, remain = update()
		return data, remain
	}

	window.Different_Ontapped = func() (data GuiData, remain int) {
		data, remain = update()
		// fmt.Printf("Different Callback : %s\n", data)
		return data, remain
	}

	window.IsLastItem = func() bool {
		return current_index >= queue.Capacity()-1
	}

	window.Cleanup_Handler = func() {
		//diff_check.ClearTempFile()
	}

	window.ShowAndRun()

	// Cleanup
	diff_check.ClearTempFile()
	skipped_item.Concat(queue, current_index)

	os.Remove("similar_data.txt") // Remove file to prevent overwriting
	data_file, data_err := os.OpenFile("similar_data.txt", os.O_RDWR|os.O_CREATE, 0755)
	if data_err != nil {
		log.Errorf("[OutputReport] ERROR: Cannot Open Data File.\n")
		log.Errorf("[OutputReport] ------ Error   : %s\n", data_err)
	}
	defer data_file.Close()

	fmt.Println("Capacity:", skipped_item.Capacity(), "Original:", queue.Capacity())

	for i := 0; i < skipped_item.Capacity(); i++ {
		_, err := data_file.WriteString(fmt.Sprintf("%s ??? %s\n", skipped_item.image1_path[i], skipped_item.image2_path[i]))
		if err != nil {
			log.Errorf("Error when writing string to data file: %v\n", err)
		}
	}
}
