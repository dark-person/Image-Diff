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

// func OutputDifferenceImage(first_image_path, second_image_path string) (output_image_path string, err error) {

// 	output_image_path = "Temp/compare.jpg"

// 	first_file, err_open1 := os.Open(first_image_path)
// 	second_file, err_open2 := os.Open(second_image_path)

// 	if err_open1 != nil {
// 		log.Println("Error opening file", err_open1)
// 		return "", err_open1
// 	} else if err_open2 != nil {
// 		log.Println("Error opening file", err_open2)
// 		return "", err_open2
// 	}

// 	// Create new image for visualization
// 	os.MkdirAll(filepath.Dir(output_image_path), 0755)
// 	compare_file, err_create := os.Create(output_image_path)
// 	if err_create != nil {
// 		return "", err_create
// 	}
// 	writer := bufio.NewWriter(compare_file)

// 	image1, _, err_decode1 := image.Decode(first_file)
// 	image2, _, err_decode2 := image.Decode(second_file)
// 	if err_decode1 != nil {
// 		log.Println("Error decode file", err_decode1)
// 		return "", err_open1
// 	} else if err_decode2 != nil {
// 		log.Println("Error decode file", err_decode2)
// 		return "", err_open2
// 	}

// 	var image_Bound_x int
// 	var image_Bound_y int
// 	if image1.Bounds().Dx() > image2.Bounds().Dx() {
// 		image_Bound_x = image1.Bounds().Dx()
// 	} else {
// 		image_Bound_x = image2.Bounds().Dx()
// 	}

// 	if image1.Bounds().Dy() > image2.Bounds().Dy() {
// 		image_Bound_y = image1.Bounds().Dy()
// 	} else {
// 		image_Bound_y = image2.Bounds().Dy()
// 	}

// 	resized_image1 := image.NewRGBA(image.Rect(0, 0, image_Bound_x, image_Bound_y))
// 	resized_image2 := image.NewRGBA(image.Rect(0, 0, image_Bound_x, image_Bound_y))

// 	draw.ApproxBiLinear.Scale(resized_image1, resized_image1.Rect, image1, image1.Bounds(), draw.Over, nil)
// 	draw.ApproxBiLinear.Scale(resized_image2, resized_image2.Rect, image2, image2.Bounds(), draw.Over, nil)

// 	compared_image := image.NewRGBA(image.Rect(0, 0, image_Bound_x, image_Bound_y))

// 	for x := 0; x < resized_image1.Bounds().Dx(); x++ {
// 		for y := 0; y < resized_image1.Bounds().Dy(); y++ {
// 			if resized_image1.At(x, y) == resized_image2.At(x, y) {
// 				compared_image.Set(x, y, color.Black)
// 			} else {
// 				compared_image.Set(x, y, resized_image2.At(x, y))
// 			}
// 		}
// 	}

// 	jpeg.Encode(writer, compared_image, &jpeg.Options{Quality: 100})
// 	return output_image_path, nil
// // }

// Bash command:
// magick compare -metric AE -fuzz 55% test_set/akie\ \(44265104\)_1672682033.png test_set/akie\ \(44265104\)_1673037743.jpg -compose minussrc -colorspace Gray compare.jpg

/* func OutputDifferenceImage2(first_image_path, second_image_path string) (output_image_path string, err error) {
	output_image_path = "Temp/compare" + time.Now().Format("20060102150405") + ".jpg"

	// Create image file and directory first
	os.MkdirAll(filepath.Dir(output_image_path), 0755)
	os.Create(output_image_path)

	imagick.Initialize()
	defer imagick.Terminate()

	imagick.Initialize()
	defer imagick.Terminate()

	magicwand1 := imagick.NewMagickWand()
	magicwand2 := imagick.NewMagickWand()

	if err := magicwand1.ReadImage(first_image_path); err != nil {
		panic(err)
	}

	if err := magicwand2.ReadImage(second_image_path); err != nil {
		panic(err)
	}
	// _, similarity, result_wand := magicwand1.SimilarityImage(magicwand2, imagick.METRIC_FUZZ_ERROR, 100.0)
	// fmt.Println("similarity:", similarity)

	// result_wand.SetCompression(imagick.COMPRESSION_BZIP)
	// fmt.Println("Compress Completed.")
	//magicwand1.SetColorspace(imagick.COLORSPACE_GRAY)

	result_wand, distortion := magicwand1.CompareImages(magicwand2, imagick.METRIC_ABSOLUTE_ERROR)

	err = result_wand.SetImageFuzz(55.0)
	if err != nil {
		panic(err)
	}

	err = result_wand.SetImageCompose(imagick.COMPOSITE_OP_MINUS_SRC)
	if err != nil {
		panic(err)
	}

	err = result_wand.SetImageColorspace(imagick.COLORSPACE_GRAY)
	if err != nil {
		panic(err)
	}

	fmt.Println("Distort:", distortion)

	if err := result_wand.WriteImage(output_image_path); err != nil {
		panic(err)
	}

	return output_image_path, nil
}
*/
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
