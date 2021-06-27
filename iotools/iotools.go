package iotools // catlab/iotools
// package main

import (
    "fmt"
    "image"
    "image/png"
    "os"
    "os/exec"
    
    pt"catlab/pathtools"
)

func LoadImage(filePath string) (image.Image) {
    f, err := os.Open(filePath)
    defer f.Close()
    if err != nil {
        panic(fmt.Sprintf("Couldn't load image file:\n\t\"%s\"\n%s", filePath, err))
    }
    img, _, err := image.Decode(f)
    if err != nil {
        panic(fmt.Sprintf("Couldn't decode image file:\n\t\"%s\"\n%s", filePath, err))
    }
    return img
}
func SavePNG(m *image.NRGBA) string {
    name := pt.NameSpacer("image.png")
    w, err := os.Create(name)
    defer w.Close()
    encoder := &png.Encoder{CompressionLevel: png.BestCompression}
    err = encoder.Encode(w, m)
    if err != nil {
        panic(err)
    }
    return name
}
func set_at(x, y int, im *image.NRGBA, vals []uint8) {
    i := y*im.Stride + x*4
    (*im).Pix[i] = (vals)[0] // red
    (*im).Pix[i+1] = (vals)[1] // green
    (*im).Pix[i+2] = (vals)[2] // blue
    (*im).Pix[i+3] = (vals)[3] // alpha
}
func RenderPNG(mat [][][]uint8) string {
    height := len(mat)
    width := len(mat[0])
    im := image.NewNRGBA(image.Rect(0, 0, width, height))
    for y := 0; y < height; y++ {
        for x := 0; x < width; x++ {
            go set_at(x, y, im, mat[y][x])
        }
    }
    return SavePNG(im)
}
func Pop(name string) {
    err := exec.Command("explorer", name).Run()
    if err != nil {
        fmt.Println(err)
    }
}








func main() {
    folder := "c:/users/kenneth/pictures"
    folders := pt.Folders(folder)
    for _, f := range folders {
        fmt.Println(f)
    }
    files := pt.Files(folder)
    for _, f := range files {
        fmt.Println(f)
    }
    // Show(Files(folder))
}


