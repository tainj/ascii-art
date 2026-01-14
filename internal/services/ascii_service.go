package services

import (
	"errors"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"os"
	"path"
	"math"

	"github.com/disintegration/imaging"
	"github.com/nfnt/resize"
	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"
)

var asciiChars = []rune("$@B%8&WM#*oahkbdpqwmZO0QLCJUYXzcvunxrjft/\\|()1{}[]?-_+~<>i!lI;:,\"^`'. ")
var charLen = len(asciiChars)

func ScaleGrayImageJpg(filePath string, widthSymbols int) (image.Image, error) {
	// Открываем файл
	imgIn, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer imgIn.Close()

	// Получаем файл
	imgJpg, err := jpeg.Decode(imgIn)
	if err != nil {
		return nil, err
	}

	// Считаем высоту
	heightSymbols := int(imgJpg.Bounds().Dy() / imgJpg.Bounds().Dx() * widthSymbols)
	heightSymbolInPixel := int(math.Round(float64(imgJpg.Bounds().Dy()) / float64(heightSymbols)))

	widthSymbolInPixel := int(math.Round(float64(float64(imgJpg.Bounds().Dx()) / float64(widthSymbols))))

	// Новое изображение
	dst := image.NewGray(image.Rect(0, 0, widthSymbols, heightSymbols))

	// Преобразуем в серый цвет
	bounds := imgJpg.Bounds()
	for y := 0; y < heightSymbols; y++ {
		for x := 0; x < widthSymbols; x++ {
			var flow uint32
			for i := 0; i < heightSymbolInPixel; i++ {
				for j := 0; j < widthSymbolInPixel; j++ {
					r, g, b, _ := imgJpg.At(bounds.Min.X + x * widthSymbolInPixel + j, bounds.Min.Y + y * heightSymbolInPixel + i).RGBA()
					flow += (r*299 + g*587 + b*114) / 1000
				}
			}
			avgGray := uint8((flow / uint32(heightSymbolInPixel * widthSymbolInPixel)) >> 8) // 0-255
            dst.SetGray(x, y, color.Gray{avgGray})
		}
	}

	// Сохранение результата
    outPath := path.Join("./temp", "ascii_art2.jpg")
    if err := os.MkdirAll("./temp", 0755); err != nil {
        return nil, err
    }
    
    f, err := os.Create(outPath)
    if err != nil {
        return nil, err
    }
    defer f.Close()
    
    if err := jpeg.Encode(f, dst, nil); err != nil {
        return nil, err
    }
	return dst, nil
}

func ConvertToASCII(filePath string) (string, error) {
    imgIn, _ := os.Open(filePath)
    imgJpg, _ := jpeg.Decode(imgIn)
    imgIn.Close()

    imgJpg = resize.Resize(120, 120, imgJpg, resize.Bicubic)
	imgJpg = imaging.Grayscale(imgJpg)

    // Преобразование в ASCII
    data := make([]string, 0)
    bounds := imgJpg.Bounds()
    for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
        s := ""
        for x := bounds.Min.X; x < bounds.Max.X; x++ {
            r, _, _, _ := imgJpg.At(x, y).RGBA() // Достаточно одного компонента
            brightness := uint8(r >> 8)          // 0-255
            index := int(brightness) * (len(asciiChars) - 1) / 255
            if index >= len(asciiChars) { // Защита от выхода за границы
                index = len(asciiChars) - 1
            }
            s += string(asciiChars[index])
        }
        data = append(data, s)
    }

    // Проверка на пустые данные
    if len(data) == 0 || len(data[0]) == 0 {
        return "", errors.New("ascii conversion produced empty data")
    }

    // Параметры шрифта (реальные размеры)
    const fontWidth = 7
    const fontHeight = 13

    // Создание холста
    imgWidth := len(data[0]) * fontWidth
    imgHeight := len(data) * fontHeight
    canvas := image.NewRGBA(image.Rect(0, 0, imgWidth, imgHeight))
    draw.Draw(canvas, canvas.Bounds(), &image.Uniform{color.White}, image.Point{}, draw.Src)

    // Отрисовка
    d := &font.Drawer{
        Dst:  canvas,
        Src:  image.NewUniform(color.Black),
        Face: basicfont.Face7x13,
    }

    for i, line := range data {
        d.Dot = fixed.Point26_6{
            X: 0,
            Y: fixed.Int26_6((i+1)*fontHeight*64 - 3*64), // Вертикальное позиционирование
        }
        d.DrawString(line)
    }

    // Сохранение результата
    outPath := path.Join("./temp", "ascii_art.jpg")
    if err := os.MkdirAll("./temp", 0755); err != nil {
        return "", err
    }
    
    f, err := os.Create(outPath)
    if err != nil {
        return "", err
    }
    defer f.Close()
    
    if err := jpeg.Encode(f, canvas, nil); err != nil {
        return "", err
    }

    return outPath, nil
}