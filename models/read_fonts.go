package models

import (
	"fmt"
	image2 "image"
	"image/color"
	"image/draw"
	"image/png"
	"log"
	"os"

	"golang.org/x/text/encoding/charmap"
)

const mlength = 70

var smap = createSymbMap()

const fontsFname = "./files/fonts/fonta.fnt"

type ImageData []byte

func createSymbMap() map[byte]ImageData {
	file, err := os.Open(fontsFname)
	if err != nil {
		fmt.Println(err)
	}
	buffer := make([]byte, 10240)
	n, err := file.Read(buffer)
	symbMap := make(map[byte]ImageData)
	data := buffer[:n]
	i := 6
	for i < len(data) {
		num := data[i]
		i += 1
		count := data[i]
		i += 1
		symbMap[num] = data[i : i+int(count*2)]
		i += int(count) * 2
	}
	return symbMap
}

func (i ImageData) GetLen() int {
	return len([]byte(i))
}

func (id ImageData) GetImage() image2.Image {
	data := []byte(id)
	ln := len(data)
	image := image2.NewRGBA(image2.Rect(0, 0, ln+2, 32))
	for i := ln; i < ln+2; i++ {
		for k := 0; k < 8; k++ {
			pos := i*8 + k
			image.Set((pos/16)*2, 2*(pos%16), color.Black)
			image.Set((pos/16)*2, 2*(pos%16)+1, color.Black)
			image.Set((pos/16)*2+1, 2*(pos%16), color.Black)
			image.Set((pos/16)*2+1, 2*(pos%16)+1, color.Black)
		}
	}

	for i := 0; i < ln; i++ {
		dat := data[i]
		var mask byte = 1
		for k := 0; k < 8; k++ {
			pos := i*8 + k
			if mask&dat > 0 {
				image.Set((pos/16)*2, 2*(pos%16), color.White)
				image.Set((pos/16)*2, 2*(pos%16)+1, color.White)
				image.Set((pos/16)*2+1, 2*(pos%16), color.White)
				image.Set((pos/16)*2+1, 2*(pos%16)+1, color.White)
			} else {
				image.Set((pos/16)*2, 2*(pos%16), color.Black)
				image.Set((pos/16)*2, 2*(pos%16)+1, color.Black)
				image.Set((pos/16)*2+1, 2*(pos%16), color.Black)
				image.Set((pos/16)*2+1, 2*(pos%16)+1, color.Black)
			}
			mask = mask << 1
		}

	}
	return image

}

/*getImageFromStr Получение image из одной строки
 */
func getImageFromStr(str string) image2.Image {
	data := HandleSymb(ToCp1251(str))
	images := make([]image2.Image, len(data))
	ln := 0 //суммарная длина
	for i, s := range data {
		bt := s
		bf := smap[bt]
		img := bf.GetImage()
		images[i] = img
		ln += bf.GetLen() + 2
	}
	sImage := image2.NewRGBA(image2.Rect(0, 0, ln, 32))
	pnt := 0
	for _, img := range images {
		draw.Draw(sImage, image2.Rect(pnt, 0, pnt+img.Bounds().Dx(), 32), img, image2.Point{0, 0}, draw.Src)
		pnt += img.Bounds().Dx()
	}
	return sImage

}

/*SaveImageToPng - сохранение в формате png */
func SaveImageToPng(img image2.Image, path string) {
	f, err := os.Create(path)
	if err != nil {
		log.Fatal(err)
	}
	if err := png.Encode(f, img); err != nil {
		f.Close()
		log.Fatal(err)
	}

	if err := f.Close(); err != nil {
		log.Fatal(err)
	}
}

/*GetPaintImage - получение картинки из списка строк */
func GetPaintImage(data []string) image2.Image {
	image := image2.NewRGBA(image2.Rect(0, 0, mlength*8, 32*len(data)))
	for i, s := range data {
		if len(s) < mlength {
			n := len(s)
			for n < mlength {
				n++
				s = s + " "
			}
		}
		im := getImageFromStr(s)
		draw.Draw(image, image2.Rect(0, 32*i, im.Bounds().Dx(), 32*(i+1)), im, image2.Point{0, 0}, draw.Src)
	}
	return image
}

/*HandleSymb -
<-60
>-62
*/
func HandleSymb(data []byte) []byte {
	bg := -1
	count := 0
	for i := 0; i < len(data); i++ {
		if data[i] == 60 {
			bg = i
		}
		if bg > -1 {
			if data[i] > 47 && data[i] < 58 {
				count++
			}
		}
		if bg > -1 && data[i] == 62 && count == 3 {
			data[bg] = byte(int(data[bg+1]-48)*100 + int(data[bg+2]-48)*10 + int(data[bg+3]-48))
			copy(data[bg+1:], data[bg+5:])
			return HandleSymb(data[:len(data)-4])
		}
		if count > 3 {
			bg = -1
			count = 0
		}
	}
	return data
}

/*ToCp1251 - преобразование в формат 1251 */
func ToCp1251(str string) []byte {
	dec := charmap.Windows1251.NewEncoder()
	dst := make([]byte, len(str)*2)
	n, _, err := dec.Transform(dst, []byte(str), false)
	if err != nil {
		fmt.Println("Dec error:", err)
	}
	return dst[:n]
}
