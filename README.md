# Image-Diff

An Application for finding image difference. 

Graphic User Interface is included (by `fyne.io`).

![ezgif com-optimize](https://user-images.githubusercontent.com/66205094/218022360-b63245c0-35a0-40f8-8be9-375207343616.gif)

## Graphic User Interface

This app will compare the image with three different aspect:
1. Original Image Compare: Filename and Filesize
2. Difference: Overlay two images, visualize the difference
3. Animated: Use GIF to visualize the difference of two image

### Keyboard Hotkeys
Use Arrow Key (Left/Right) can switch tabs.

## Run

### Pre-requisite
**Please ensure ImageMagick is installed.** [Download Link](https://imagemagick.org/script/download.php)

### Data
This app will import the data from an external file `similar_data.txt`. The format will like:
```
D:\Private\SetA\first_image.png ??? D:\SetA\second_image.png
D:\Private\SetB\first_image.png ??? D:\Private\SetB\second_image.png
```
The string ` ??? ` (*Please notice there is space before and after ???*) is used as seperator, 

while line break is represent another image set.

### Command

```
go run .
```
Compile and Run:
```
go build
.\image-diff.exe
```

## Advanced
The directory tree of necessary file
```
Image-Diff
 ┣ resource
 ┃ ┣ filename_placeholder.gif
 ┃ ┣ loading.gif
 ┃ ┗ loading.jpg
 ┣ Temp
 ┣ .gitignore
 ┣ data.go
 ┣ go.mod
 ┣ go.sum
 ┣ gui.go
 ┣ ImageDiff.go
 ┣ log_custom.go
 ┣ main.go
 ┗ util.go
```
The program is seperate to three different parts:
* ImageDiff : The Basic Part of generate image and handle image queue
* gui: Graphic User Interface Part
* Logger: the custom format of logger

If you need to customizer GUI, or implement the logic to other program, you can ignore these files:
* All file in `resouce` folder
* `gui.go`
* `log_custom.go` if you already/want to create logger yourself

