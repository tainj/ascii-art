package services

import (
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"os"
	"path"
	"math"
	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"
)

var asciiChars = []rune{ // Палитра
    '█',
    '▓', 
    '▒', 
    '░', 
    '#', 
    '&', 
    '@', 
    '%', 
    '$', 
    '+', 
    '=', 
    'o', 
    '*', 
    ':', 
    '-', 
    '.', 
    ' ', // Пробел (самый светлый)
}
var charLen = len(asciiChars)

func ScaleGrayImage(filePath string, widthSymbols int) (image.Image, error) {
	// Открываем файл
	imgIn, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer imgIn.Close()

	// Получаем файл
	img, _, err := image.Decode(imgIn)
	if err != nil {
		return nil, err
	}

	// Считаем высоту
	heightSymbols := int(math.Round(float64(img.Bounds().Dy()) / float64(img.Bounds().Dx()) * float64(widthSymbols)) / 2.0)

	// Новое изображение
	dst := image.NewGray(image.Rect(0, 0, widthSymbols, heightSymbols))

	// Преобразуем в серый цвет
	bounds := img.Bounds()
	for y := range heightSymbols {
        // Определяем границы пикселей, которые мы рассматриваем

        startY := int(math.Floor(float64(y) * float64(bounds.Dy() / heightSymbols)))
        endY := int(math.Ceil(float64(y + 1) * float64((bounds.Dy() / heightSymbols))))
		for x := range widthSymbols {
            // Определяем границы пикселей, которые мы рассматриваем

            startX := int(math.Floor(float64(x) * float64(bounds.Dx() / widthSymbols)))
            endX := int(math.Ceil(float64(x + 1) * float64((bounds.Dx() / widthSymbols))))
			var flow uint32

			for i := startY; i < endY; i++ {
				for j := startX; j < endX; j++ {
					r, g, b, _ := img.At(j, i).RGBA()
					flow += (r*299 + g*587 + b*114) / 1000
				}
			}
            avgGray := uint8((flow / uint32((endX - startX) * (endY - startY))) >> 8)
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

func ConvertImgToASCII(img image.Image) []string {
    lines := make([]string, 0)
    bounds := img.Bounds()
    for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
        s := ""
        for x := bounds.Min.X; x < bounds.Max.X; x++ {
            r, _, _, _ := img.At(x, y).RGBA() // Достаточно одного компонента
            brightness := uint8(r >> 8)          // 0-255
            index := int(brightness) * (len(asciiChars) - 1) / 255
            if index >= len(asciiChars) { // Защита от выхода за границы
                index = len(asciiChars) - 1
            }
            s += string(asciiChars[index])
        }
        lines = append(lines, s)
    }
    return lines
}

func CreateImgFromASCII(lines []string) (string, error) {
    // Параметры шрифта (реальные размеры)
    const fontWidth = 7
    const fontHeight = 13

    // Создание холста
    imgWidth := len(lines[0]) * fontWidth
    imgHeight := len(lines) * fontHeight
    canvas := image.NewRGBA(image.Rect(0, 0, imgWidth, imgHeight))
    draw.Draw(canvas, canvas.Bounds(), &image.Uniform{color.White}, image.Point{}, draw.Src)

    // Отрисовка
    d := &font.Drawer{
        Dst:  canvas,
        Src:  image.NewUniform(color.Black),
        Face: basicfont.Face7x13,
    }

    for i, line := range lines {
        d.Dot = fixed.Point26_6{
            X: 0,
            Y: fixed.Int26_6((i+1)*fontHeight*64 - 3*64), // Вертикальное позиционирование
        }
        d.DrawString(line)
    }

    // Сохранение результата
    outPath := path.Join("./temp", "ascii_art3.jpg")
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
