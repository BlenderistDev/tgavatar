package avatar

import (
	"bytes"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"os"
)

type Generator struct {
}

func (g Generator) Generate(hour int) ([]byte, error) {
	img, err := g.getTemplate(hour)

	// Получаем размеры изображения
	bounds := img.Bounds()
	width := bounds.Max.X
	height := bounds.Max.Y

	// Создаем новое изображение с черными плоскостями
	newImg := image.NewRGBA(bounds)

	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			// Копируем пиксель из старого изображения в новое
			pixel := img.At(x, y)
			newImg.Set(x, y, pixel)

			if isEyeOpen(hour) {
				if isWhiteOfEye(x, y) {
					newImg.Set(x, y, color.White)
				}

				if isPupil(x, y) {
					newImg.Set(x, y, color.RGBA{
						R: 55,
						G: 38,
						B: 27,
					})
				}

				if isGlasses(x, y) {
					newImg.Set(x, y, color.RGBA{
						R: 55,
						G: 38,
						B: 27,
					})
				}
			}

			if isBread(x, y, hour) {
				if isBreadMouth(x, y) {
					newImg.Set(x, y, color.White)
				} else {
					newImg.Set(x, y, color.RGBA{
						R: 138,
						G: 93,
						B: 62,
					})
				}
			}

		}
	}

	// Кодирование картинки в байты
	buf := new(bytes.Buffer)
	err = jpeg.Encode(buf, newImg, nil)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (g Generator) getTemplate(hour int) (image.Image, error) {
	var path string

	// choose image with correct background
	switch {
	case hour <= 7:
		path = "img/night.png"
	case hour >= 8 && hour <= 11:
		path = "img/morning.png"
	case hour >= 12 && hour <= 16:
		path = "img/day.png"
	case hour >= 17:
		path = "img/evening.png"
	}

	// open image
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// decode image
	img, err := png.Decode(file)
	if err != nil {
		return nil, err
	}

	return img, nil
}

func isBread(x, y, hour int) bool {
	if hour == 9 {
		return false
	}

	if hour >= 10 {
		// left whisker
		if x >= 313 && x <= 338 && y >= 415 && y <= 510 {
			return true
		}

		// right whisker
		if x >= 664 && x <= 689 && y >= 415 && y <= 510 {
			return true
		}

		// bristle
		if x >= 313 && x <= 689 && y >= 466 && y <= 510 {
			return true
		}

		// mustaches
		if x >= 426 && x <= 576 && y >= 443 && y <= 465 {
			return true
		}

		// bread first layer
		if x >= 373 && x <= 629 && y >= 510 && y <= 542 {
			return true
		}

		breadStepSize := 24
		breadAdditionalLayers := hour - 10
		breadStartY := 543

		// growing layer
		if breadAdditionalLayers > 2 {
			if x >= 390 && x <= 610 && y >= breadStartY && y <= breadStartY+breadStepSize*(breadAdditionalLayers-2) {
				return true
			}
		}

		// first narrow layer
		if breadAdditionalLayers > 1 {
			if x >= 414 && x <= 586 && y >= breadStartY && y >= breadStartY+breadStepSize*(breadAdditionalLayers-2) && y <= breadStartY+breadStepSize*(breadAdditionalLayers-1) {
				return true
			}
		}

		// second narrow layer
		if breadAdditionalLayers > 0 {
			if x >= 460 && x <= 540 && y >= breadStartY && y >= breadStartY+breadStepSize*(breadAdditionalLayers-1) && y <= breadStartY+breadStepSize*breadAdditionalLayers {
				return true
			}
		}
	}

	return false
}

func isBreadMouth(x int, y int) bool {
	if x >= 452 && x <= 550 && y >= 464 && y <= 474 {
		return true
	}

	return false
}

func isEyeOpen(hour int) bool {
	if hour >= 10 || hour <= 1 {
		return true
	}

	return false
}

// isPupil is pixel in pupil
func isPupil(x, y int) bool {
	// left pupil
	if x >= 396 && x <= 433 && y >= 323 && y <= 356 {
		return true
	}

	// right pupil
	if x >= 569 && x <= 606 && y >= 323 && y <= 356 {
		return true
	}

	return false
}

func isWhiteOfEye(x, y int) bool {
	// left eye
	if x >= 365 && x <= 464 && y >= 323 && y <= 362 {
		return true
	}

	// right eye
	if x >= 538 && x <= 637 && y >= 323 && y <= 362 {
		return true
	}

	return false
}

func isGlasses(x, y int) bool {
	// left horizontal
	if x >= 343 && x <= 473 && ((y >= 312 && y <= 314) || (y >= 375 && y <= 377)) {
		return true
	}

	// left vertical
	if ((x >= 342 && x <= 344) || (x >= 472 && x <= 474)) && y >= 313 && y <= 376 {
		return true
	}

	// right horizontal
	if x >= 529 && x <= 659 && ((y >= 312 && y <= 314) || (y >= 375 && y <= 377)) {
		return true
	}

	// right vertical
	if ((x >= 528 && x <= 530) || (x >= 658 && x <= 660)) && y >= 313 && y <= 376 {
		return true
	}

	// center
	if x >= 475 && x <= 527 && y >= 343 && y <= 344 {
		return true
	}

	return false
}
