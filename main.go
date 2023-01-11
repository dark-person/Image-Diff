package main

import (
	"fmt"
	"os"
	"time"
)

func Test_Set(set_number int) (string, string) {
	switch set_number {
	case 3:
		return "test_set\\blue gk_1672650631.png", "test_set\\blue gk_1672650653.png"
	case 2:
		return "test_set\\an yasuri_1672722083_faceleft.png", "test_set\\an yasuri_1672922692_faceright.jpg"
	case 1:
		return "test_set\\akie (44265104)_1672682033.png", "test_set\\akie (44265104)_1673037743.jpg"
	}
	return "", ""

}

func main() {
	fmt.Println("Hello world!")

	diff_check := NewImageDiff()
	diff_check.SetDiffImageDir("Temp")
	err := diff_check.Init()
	if err != nil {
		fmt.Println("Image Magick is not installed.", err)
		os.Exit(1)
	}

	// for i := 1; i <= 3; i++ {
	// 	path1, path2 := Test_Set(i)
	// 	diff_check.SetImages(path1, path2)
	// 	diff_check.Diff()
	// 	diff_check.ClearData()
	// }

	path1, path2 := Test_Set(3)
	diff_check.SetImages(path1, path2)
	diff_check.Diff()

	window := NewGuiWindow()

	go func() {
		time.Sleep(time.Second * 1)
		fmt.Println("Trigger Update")
		window.Update(NewGuiData(diff_check.GetInfo()))
	}()

	window.ShowAndRun()

	// diff_check.ClearData()
}
