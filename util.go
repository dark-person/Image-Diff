package main

import (
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

// Using Command line to exec magick

func OutputImageDiffJpg(first_image_path, second_image_path string) (output_image_path string, err error) {
	output_image_path = "Temp/compare" + time.Now().Format("20060102150405") + ".jpg"

	// Create image file and directory first
	os.MkdirAll(filepath.Dir(output_image_path), 0755)

	// Check version
	out, err := exec.Command("magick", "-version").Output()
	if err != nil {
		fmt.Println("No magick output found. Possibly not installed.")
		panic(err)
	}
	if strings.Contains(string(out), "not found") {
		return "", fmt.Errorf("magick not installed")
	}

	// Run command
	cmd := exec.Command("magick", "compare", "-metric", "AE", "-fuzz", "55%", first_image_path, second_image_path, "-compose",
		"minussrc", "-colorspace", "Gray", output_image_path)

	fmt.Println("Start Generating image..")
	out, err = cmd.Output()
	if err != nil {
		fmt.Println("Image Generate Failed")
		panic(err)
	}
	fmt.Println(string(out))

	fmt.Println("Command run successfully.")

	return "", nil
}

func OutputImageDiffJpg2(first_image_path, second_image_path string) (output_image_path string, err error) {
	output_image_path = "Temp/compare" + time.Now().Format("20060102150405") + ".jpg"

	// Create image file and directory first
	os.MkdirAll(filepath.Dir(output_image_path), 0755)

	// Check version
	out, err := exec.Command("magick", "-version").Output()
	if err != nil {
		fmt.Println("No magick output found. Possibly not installed.")
		panic(err)
	}
	if strings.Contains(string(out), "not found") {
		return "", fmt.Errorf("magick not installed")
	}

	// Run command
	cmd := exec.Command("magick", "compare", "-metric", "AE", "-fuzz", "55%", first_image_path, second_image_path, "-compose",
		"minussrc", "-colorspace", "Gray", output_image_path)

	fmt.Println("Start Generating image..")
	out, err = cmd.Output()
	if err != nil {
		fmt.Println("Image Generate Failed")
		panic(err)
	}
	fmt.Println(string(out))

	fmt.Println("Command run successfully.")

	return "", nil
}

func OutputImageDiffGif(first_image_path, second_image_path string) (output_image_path string, err error) {
	output_image_path = "Temp/compare_" + time.Now().Format("20060102150405") + ".gif"

	// Create image file and directory first
	os.MkdirAll(filepath.Dir(output_image_path), 0755)

	// Check version
	out, err := exec.Command("magick", "-version").Output()
	if err != nil {
		fmt.Println("No magick output found. Possibly not installed.")
		panic(err)
	}
	if strings.Contains(string(out), "not found") {
		return "", fmt.Errorf("magick not installed")
	}

	// Resize both image to same dimension
	first_file, err_open1 := os.Open(first_image_path)
	second_file, err_open2 := os.Open(second_image_path)

	if err_open1 != nil {
		log.Println("Error opening file", err_open1)
		return "", err_open1
	} else if err_open2 != nil {
		log.Println("Error opening file", err_open2)
		return "", err_open2
	}

	image1, _, err_decode1 := image.Decode(first_file)
	image2, _, err_decode2 := image.Decode(second_file)
	if err_decode1 != nil {
		log.Println("Error decode file", err_decode1)
		return "", err_open1
	} else if err_decode2 != nil {
		log.Println("Error decode file", err_decode2)
		return "", err_open2
	}

	var new_bound image.Rectangle
	if image1.Bounds().Dx() > image2.Bounds().Dx() {
		fmt.Println("Image 1 is larger.")
		new_bound = image1.Bounds()
	} else {
		fmt.Println("Image 2 is larger.")
		new_bound = image2.Bounds()
	}

	if new_bound == image1.Bounds() {
		fmt.Println("Resizing image 2..")
		cmd := exec.Command("magick", second_image_path, "-resize", strconv.Itoa(new_bound.Bounds().Dx())+"x"+strconv.Itoa(new_bound.Bounds().Dy())+"!",
			"Temp/temp_"+filepath.Base(second_image_path))
		out, _ := cmd.Output()
		fmt.Println(string(out))

		second_image_path = "Temp/temp_" + filepath.Base(second_image_path)
	} else {
		fmt.Println("Resizing image 1..")
		cmd := exec.Command("magick", first_image_path, "-resize", strconv.Itoa(new_bound.Bounds().Dx())+"x"+strconv.Itoa(new_bound.Bounds().Dy())+"!",
			"Temp/temp_"+filepath.Base(first_image_path))
		out, _ := cmd.Output()
		fmt.Println(string(out))

		first_image_path = "Temp/temp_" + filepath.Base(first_image_path)
	}

	// Run command
	cmd := exec.Command("magick", "-delay", "50", first_image_path, second_image_path, "-loop", "0", output_image_path)

	fmt.Println("Start Generating gif..")
	out, err = cmd.Output()
	if err != nil {
		fmt.Println("Image Generate gif")
		panic(err)
	}
	fmt.Println(string(out))

	fmt.Println("Command run successfully.")

	return "", nil
}
