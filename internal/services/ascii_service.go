package services

import (
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"os"
	"path"

	"log"

	"github.com/disintegration/imaging"
	"github.com/nfnt/resize"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
	"golang.org/x/image/font/basicfont"
)

var asciiChars = []rune("$@B%8&WM#*oahkbdpqwmZO0QLCJUYXzcvunxrjft/\\|()1{}[]?-_+~<>i!lI;:,\"^`'. ")
var charLen = len(asciiChars)

func ConvertToASCII(filePath string) (string, error) {
	imgIn, _ := os.Open(filePath)
    imgJpg, _ := jpeg.Decode(imgIn)
    imgIn.Close()

    imgJpg = resize.Resize(120, 120, imgJpg, resize.Bicubic)
	imgJpg = imaging.Grayscale(imgJpg)

    imgOut, err := os.Create(path.Join("./temp", "test2.jpg"))
	if err != nil {
		log.Println(err)
		return "", err
	}

	data := make([]string, 0)
	bounds := imgJpg.Bounds()
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		s := ""
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			pixel := imgJpg.At(x, y)
			if r, g, b, _ := pixel.RGBA(); r == g && g == b {
				brightness := uint8(r >> 8) // 0-255
				index := int(brightness) * (charLen - 1) / 255
				s += string(asciiChars[index])
			}
		}
		data = append(data, s)
	}

	// Параметры шрифта
	charWidth := 6   // Ширина символа в basicfont
	charHeight := 10 // Высота строки

	// Расчёт размеров изображения
	imgWidth := len(data[0]) * charWidth
	imgHeight := len(data) * charHeight

	// Создаём холст (белый фон)
	canvas := image.NewRGBA(image.Rect(0, 0, imgWidth, imgHeight))
	bgColor := color.RGBA{255, 255, 255, 255}
	draw.Draw(canvas, canvas.Bounds(), &image.Uniform{bgColor}, image.Point{}, draw.Src)

	// Рисуем каждую строку
	point := fixed.Point26_6{
		X: fixed.Int26_6(0 * 64), // Начало по X
		Y: fixed.Int26_6(0 * 64), // Начало по Y
	}

	d := &font.Drawer{
		Dst:  canvas,
		Src:  image.NewUniform(color.Black), // Чёрный текст
		Face: basicfont.Face7x13,            // Моноширинный шрифт
		Dot:  point,
	}

	for i, line := range data {
		d.Dot.Y = fixed.Int26_6((i + 1) * charHeight * 64)
		d.DrawString(line)
	}



    jpeg.Encode(imgOut, canvas, nil)
    imgOut.Close()
	return "", nil
}

