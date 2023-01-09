package main

import (
	"fmt"
	"image"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
)

type ImageDiff struct {
	first_image_path  string // Relative path to first image
	second_image_path string // Relative path to second image

	processed_image1 string //Relative path to first image (Processed)
	processed_image2 string //Relative path to second image (Processed)

	diff_jpg_path string // Difference with gray scale
	diff_gif_path string // Difference with animated GIF of two images

	initialized    bool   // Security lock to check object validity
	diff_directory string // Directory that store temporary image
}

var Err_ImageDiff_Not_Initialized = fmt.Errorf("imageDiff not initalised yet")

const time_format_string = "20060102150405"

// Create a new ImageDiff object
func NewImageDiff() *ImageDiff {
	return &ImageDiff{
		first_image_path:  "",
		second_image_path: "",

		processed_image1: "",
		processed_image2: "",

		diff_jpg_path: "",
		diff_gif_path: "",

		initialized:    false,
		diff_directory: "",
	}
}

// Initalize a new ImageDiff, must be called after set up
// Check the ImageMagick is installed or not
func (diff *ImageDiff) Init() error {
	// Check version
	out, err := exec.Command("magick", "-version").Output()
	if err != nil {
		log.Errorf("No magick output found. Possibly not installed.\n")
		return fmt.Errorf("magick command run failed: %v", err)
	}
	if strings.Contains(string(out), "not found") {
		return fmt.Errorf("magick not not found")
	}

	diff.initialized = true
	return nil
}

// Set the directory that store those diff images result
func (diff *ImageDiff) SetDiffImageDir(directory_relative_path string) error {
	diff.diff_directory = directory_relative_path
	err := os.MkdirAll(directory_relative_path, 0755)
	return err
}

// Set the two image (relative path) to the ImageDiff, also processed image
func (diff *ImageDiff) SetImages(image1_path string, image2_path string) error {
	if !diff.initialized {
		return Err_ImageDiff_Not_Initialized
	}

	diff.first_image_path = image1_path
	diff.second_image_path = image2_path

	// Resize both image to same dimension
	first_file, err_open1 := os.Open(diff.first_image_path)
	second_file, err_open2 := os.Open(diff.second_image_path)

	if err_open1 != nil {
		log.Errorf("Error opening file 1\n", err_open1)
		return err_open1
	} else if err_open2 != nil {
		log.Errorf("Error opening file 2\n", err_open2)
		return err_open2
	}

	image1, _, err_decode1 := image.Decode(first_file)
	image2, _, err_decode2 := image.Decode(second_file)
	if err_decode1 != nil {
		log.Errorf("Error decode file 1\n", err_decode1)
		return err_decode1
	} else if err_decode2 != nil {
		log.Errorf("Error decode file 2\n", err_decode2)
		return err_decode2
	}

	var new_bound image.Rectangle
	if image1.Bounds().Dx() > image2.Bounds().Dx() {
		log.Tracef("Image 1 is larger.\n")
		new_bound = image1.Bounds()
	} else {
		log.Tracef("Image 2 is larger.\n")
		new_bound = image2.Bounds()
	}

	if new_bound == image1.Bounds() {
		diff.processed_image1 = diff.first_image_path
		diff.processed_image2 = diff.diff_directory + "/temp_" + filepath.Base(diff.second_image_path)

		log.Debugf("Resizing image 2..\n")
		cmd := exec.Command("magick", diff.second_image_path, "-resize", strconv.Itoa(new_bound.Bounds().Dx())+"x"+strconv.Itoa(new_bound.Bounds().Dy())+"!",
			diff.processed_image2)
		out, _ := cmd.Output()
		log.Tracef(string(out) + "\n")

	} else {
		diff.processed_image1 = diff.diff_directory + "/temp_" + filepath.Base(diff.first_image_path)
		diff.processed_image2 = diff.second_image_path

		log.Debugf("Resizing image 1..\n")
		cmd := exec.Command("magick", diff.first_image_path, "-resize", strconv.Itoa(new_bound.Bounds().Dx())+"x"+strconv.Itoa(new_bound.Bounds().Dy())+"!",
			diff.processed_image1)
		out, _ := cmd.Output()
		log.Tracef(string(out) + "\n")

	}

	return nil

}

// Start the diff process, will output a jpg show the difference of two image, and a gif maked by these two images
func (diff *ImageDiff) Diff() error {
	if !diff.initialized {
		return Err_ImageDiff_Not_Initialized
	}

	// Gif creation
	gif_output := filepath.Join(diff.diff_directory, "gif_compare_"+time.Now().Format(time_format_string)+".gif")
	err := output_diff_gif(diff.processed_image1, diff.processed_image2, gif_output)
	if err != nil {
		return err
	}
	diff.diff_gif_path = gif_output

	// Jpg creation
	jpg_output := filepath.Join(diff.diff_directory, "jpg_compare_"+time.Now().Format(time_format_string)+".jpg")
	err = output_diff_jpg(diff.processed_image1, diff.processed_image2, jpg_output)
	if err != nil {
		return err
	}
	diff.diff_jpg_path = jpg_output

	// Remove the temporary files
	if strings.Contains(diff.processed_image1, "temp_") {
		os.Remove(diff.processed_image1)
	}

	if strings.Contains(diff.processed_image2, "temp_") {
		os.Remove(diff.processed_image2)
	}

	return nil
}

// Clear all image related data, allow the next comparison of other images
func (diff *ImageDiff) ClearData() {

	diff.first_image_path = ""
	diff.second_image_path = ""
	diff.processed_image1 = ""
	diff.processed_image2 = ""
	diff.diff_gif_path = ""
	diff.diff_jpg_path = ""
}

// Get the related path of diff image
func (diff *ImageDiff) GetDiffJpg() (string, error) {
	if !diff.initialized {
		return "", Err_ImageDiff_Not_Initialized
	}

	return diff.diff_jpg_path, nil
}

// Get the related path of diff gif image
func (diff *ImageDiff) GetDiffGif() (string, error) {
	if !diff.initialized {
		return "", Err_ImageDiff_Not_Initialized
	}

	return diff.diff_gif_path, nil

}

// Get all needed information for data storing/processing
func (diff *ImageDiff) GetInfo() (first_image_path, second_image_path, diff_jpg, diff_gif string) {
	return diff.first_image_path, diff.second_image_path, diff.diff_jpg_path, diff.diff_gif_path
}

// Helper function to create new diff gif
func output_diff_gif(first_image_path, second_image_path string, output_image_path string) error {
	// Run command
	cmd := exec.Command("magick", "-delay", "50", first_image_path, second_image_path, "-loop", "0", output_image_path)

	log.Debugf("Start Generating gif..\n")
	out, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("command of generate gif failed: %v", err)
	}
	log.Tracef(string(out) + "\n")

	log.Debugf("Command run successfully.\n")

	return nil
}

// Helper function to create new diff jpg
func output_diff_jpg(first_image_path, second_image_path string, output_image_path string) error {
	// Run command: magick composite "Temp\temp_blue gk_1672650653.png" "test_set\blue gk_1672650631.png" -compose difference difference.jpg
	cmd := exec.Command("magick", "composite", first_image_path, second_image_path, "-compose",
		"difference", output_image_path)

	log.Debugf("Start Generating image..")
	out, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("fail to generate diff jpg: %v", err)
	}
	log.Tracef(string(out) + "\n")

	log.Debugf("Command run successfully.")

	return nil
}
