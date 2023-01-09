package main

import "fmt"

func Test_Set(set_number int) (string, string) {
	switch set_number {
	case 3:
		return "test_set\\blue gk_1672650631.png", "test_set\\blue gk_1672650653.png"
	case 2:
		return "test_set\\an yasuri_1672722083.png", "test_set\\an yasuri_1672922692.jpg"
	case 1:
		return "test_set\\akie (44265104)_1672682033.png", "test_set\\akie (44265104)_1673037743.jpg"
	}
	return "", ""

}

func main() {
	fmt.Println("Hello world!")

	diff_check := NewImageDiff()
	diff_check.SetDiffImageDir("Temp")
	diff_check.Init()

	for i := 1; i <= 3; i++ {
		path1, path2 := Test_Set(i)
		diff_check.SetImages(path1, path2)
		diff_check.Diff()
		diff_check.ClearData()
	}

	// _, err := OutputImageDiffJpg(path1, path2)
	// if err != nil {
	// 	fmt.Println(err)
	// }

	// _, err = OutputImageDiffGif(path1, path2)
	// if err != nil {
	// 	fmt.Println(err)
	// }
}
