package main

import (
	"fmt"
	"image"
    iot"catlab/iotools"
    it"catlab/itertools"
    pt"catlab/pathtools"
)
func dummy(args...interface{}) {
	if len(args) < 0 {
		iot.LoadImage("dummy")
		pt.Files("dummy")
		it.Show([]string{"dummy"})
	}
}
func Blank(width, height int) [][][]uint8 {
	mat := make([][][]uint8, height)
	// mat := [][][]uint8{
	// 	{
	// 		{
				
	// 		},
	// 	},
	// }
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			mat[y] = append(mat[y], []uint8{})
		}
	}
	return mat
}
func PxMat(im image.Image) [][][]uint8 {
	width := im.Bounds().Max.X
	height := im.Bounds().Max.Y
	mat := Blank(width, height)
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			// i := ny*new.Stride + nx*4
			// px := im.At(x, y).RGBA()
			R, G, B, A := im.At(x, y).RGBA()
			r := uint8(R * A)
			g := uint8(G * A)
			b := uint8(B * A)
			a := uint8(A)
			(mat)[y][x] = []uint8{r, g, b, a}
		}
	}
	return mat
}
func CatMap(im image.Image) string {
	width := im.Bounds().Max.X
	height := im.Bounds().Max.Y
	// new := image.NewNRGBA(image.Rect(0, 0, dx, dy))
	mat := Blank(width, height)
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			nx := (2*x + y) % width
			ny := (x + y) % height
			// i := ny*new.Stride + nx*4
			// px := im.At(x, y).RGBA()
			R, G, B, A := im.At(x, y).RGBA()
			// r := px.
			// mat[ny][nx] = []uint8{px.R, px.G, px.B, px.A}
			r := uint8(R * A)
			g := uint8(G * A)
			b := uint8(B * A)
			a := uint8(A)
			(mat)[ny][nx] = []uint8{r, g, b, a}
		}
	}
	return iot.RenderPNG(mat)
}
func transfer(source, target *[][][]uint8, x, y, nx, ny int) {
	
}

func main() {
    // root := "c:/users/kenneth/pictures"
    // fmt.Println(root)
    // iot.Pop(root)
    
    // folders := pt.Folders(root)
    // fmt.Println(len(folders))
    // folder := "c:/users/kenneth/pictures/filtershoppes"
    folder := "E:/Projects/audio/krule_errors/slices"
    fmt.Println(folder)
    files := pt.Files(folder)
    file := files[it.Randex(len(files))]
    fmt.Println(file)
    pic := iot.LoadImage(file)
    fmt.Println(pic.Bounds().Max.X, pic.Bounds().Max.Y)
    // fmt.Println(pic)
    // iot.Pop(file)
    iot.Pop(CatMap(pic))
}
