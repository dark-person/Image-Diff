package main

import (
	"bufio"
	"log"
	"os"
	"strings"
)

// All data structure will store in here
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

func (data *GuiData) Valid() bool {
	return data.Processed1_filepath != "" && data.Processed2_filepath != ""
}

type ImagesQueue struct {
	// An array of GuiData objects
	image1_path []string
	image2_path []string
}

func (queue *ImagesQueue) Capacity() int {
	return len(queue.image1_path)
}

func (queue *ImagesQueue) Get(index int) (path1, path2 string) {
	if index < queue.Capacity() {
		return queue.image1_path[index], queue.image2_path[index]
	} else {
		return "", ""
	}

}

func (queue *ImagesQueue) Add(path1, path2 string) {
	queue.image1_path = append(queue.image1_path, path1)
	queue.image2_path = append(queue.image2_path, path2)
}

func (queue *ImagesQueue) Concat(queue2 *ImagesQueue, start_index int) {
	queue.image1_path = append(queue.image1_path, queue2.image1_path[start_index:]...)
	queue.image2_path = append(queue.image2_path, queue2.image2_path[start_index:]...)
}

func (queue *ImagesQueue) Remove(index int) {
	queue.image1_path = append(queue.image1_path[:index], queue.image2_path[:index]...)
	queue.image2_path = append(queue.image2_path[:index], queue.image2_path[:index]...)
}

func (queue *ImagesQueue) Empty() bool {
	return len(queue.image1_path) <= 0 || len(queue.image2_path) <= 0
}

func (queue *ImagesQueue) Clone() *ImagesQueue {
	var dest1 []string
	var dest2 []string

	dest1 = append(dest1, queue.image1_path...)
	dest2 = append(dest2, queue.image2_path...)

	return &ImagesQueue{
		image1_path: dest1,
		image2_path: dest2,
	}
}

func NewImagesQueue() *ImagesQueue {
	return &ImagesQueue{
		image1_path: make([]string, 0),
		image2_path: make([]string, 0),
	}
}

// Not Implmented
// Forrmat like filename1, filename2
func NewImagesQueueByFile(filepath string) *ImagesQueue {

	queue := NewImagesQueue()

	file, err := os.Open(filepath)
	if err != nil {
		log.Printf("Could not open file: %v\n", err)
	}
	defer file.Close()

	fileScanner := bufio.NewScanner(file)
	fileScanner.Split(bufio.ScanLines)

	for fileScanner.Scan() {
		if strings.Contains(fileScanner.Text(), "?") {
			path1, path2, check := strings.Cut(fileScanner.Text(), " ??? ")
			if check { // Prevent Invalid record
				queue.Add(path1, path2)
			}
		}
	}

	return queue
}
