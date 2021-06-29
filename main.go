package main

import (
    "fmt"
    "image"
    "os"
    "io/fs"
    "strings"
    // "flag"
    
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
func FibonacciSeed() func() int {
    term1 := 0
    term2 := 1
    return func() int {
        term2 = term1 + term2
        term1 = term2 - term1
        return term1
    }
}

func Blank(width, height int) [][][]uint8 {
    mat := make([][][]uint8, height)
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
    return iot.RenderPNGTensor(mat)
}
// func CatMapOnce(im image.Image) string {
//     width := im.Bounds().Max.X
//     height := im.Bounds().Max.Y
//     m := PxMat(im)
//     mat := Blank(width, height)
//     for y := 0; y < height; y++ {
//         for x := 0; x < width; x++ {
//             nx := (2*x + y) % width
//             ny := (x + y) % height
//             cell := m[y][x]
//             A := cell[3]
//             (mat)[ny][nx] = []uint8{
//                 uint8(A * cell[0]), 
//                 uint8(A * cell[1]), 
//                 uint8(A * cell[2]), 
//                 uint8(A),            
//             }
//         }
//     }
//     return iot.RenderPNG(mat)
// }
func CatMapOnceTensor(in [][][]uint8) [][][]uint8 {
    width := len(in[0])
    height := len(in)
    out := Blank(width, height)
    for y := 0; y < height; y++ {
        for x := 0; x < width; x++ {
            go transferTensor(&in, &out, x, y, width, height)
        }
    }
    return out
}
func CatMapOnce(in []uint8, width, height, stride int) []uint8 {
    out := make([]uint8, len(in))
    for y := 0; y < height; y++ {
        for x := 0; x < width; x++ {
            go transfer(&in, &out, x, y, width, height, stride)
        }
    }
    return out
}
func FullArnoldTensor(pth string) (string, int) {
    original := PxMat(iot.LoadImageFace(pth)) // [][][]uint8
    current := CatMapOnceTensor(original) // [][][]uint8
    
    home := pt.Pwd()
    proj := pt.ProjDir(pth)
    fmt.Println(proj)
    os.MkdirAll(proj, fs.ModeDir)
    pt.Cd(proj)
    defer pt.Cd(home)
    
    var ctr int
    for ; !it.SameImageTensor(current, original); ctr++ {
        current = CatMapOnceTensor(current)
        iot.RenderPNGTensor(current)
        fmt.Printf("\r%v", ctr)
    }
    fmt.Println()
    return proj, ctr
}
func FullArnold(pth string) (string, int) {
    original := iot.LoadImage(pth)
    width := original.Bounds().Max.X
    height := original.Bounds().Max.Y
    current := CatMapOnce(original.Pix, width, height, original.Stride)
    
    home := pt.Pwd()
    proj := pt.ProjDir(pth)
    proj = strings.Replace(proj, "/", string(os.PathSeparator), -1)
    // proj := strings.Replace(pt.ProjDir(pth), "\\", "/", -1)
    fmt.Println(proj)
    os.MkdirAll(proj, fs.ModeDir)
    pt.Cd(proj)
    defer pt.Cd(home)
    iot.RenderPNG(current, width, height)
    
    ctr := 1
    for ; !it.SameImage(current, original.Pix); ctr++ {
        current = CatMapOnce(current, width, height, original.Stride)
        iot.RenderPNG(current, width, height)
        fmt.Printf("\r%v", ctr)
    }
    fmt.Println()
    return proj, ctr
}
func Period(width, height int) int {
    if width == height {
        fin := false
        ctr := 0
        fib := FibonacciSeed()
        for ;!fin; ctr++ {
            n := fib()
            // if (1 == ((2*n - 1)%width)) && (0 == ((2*n)%height)) {
            if (1 == ((2*n + 1)%width)) && (1 == ((2*n + 2)%width)) {
                fin = true
            }
        }
        return ctr
    } else {
        msg := fmt.Sprintf("Period is undefined for non-square toroids\n\t(w=%v, h=%v)", width, height)
        panic(msg)
    }
}
func transferTensor(source, target *[][][]uint8, x, y, width, height int) {
    nx := (2*x + y) % width
    ny := (x + y) % height
    
    cell := (*source)[y][x]
    A := cell[3]
    
    
    (*target)[ny][nx] = []uint8{
        // uint8(A * cell[0]), // post-alpha compensation
        uint8(cell[0]), // post-alpha compensation
        // uint8(A * cell[1]), // post-alpha compensation
        uint8(cell[1]), // post-alpha compensation
        // uint8(A * cell[2]), // post-alpha compensation
        uint8(cell[2]), // post-alpha compensation
        // uint8(A), // post-alpha compensation
        uint8(A), // post-alpha compensation
    
    }
}
func transfer(source, target *[]uint8, x, y, width, height, stride int) {
    i := y*stride + 4*x
    nx := (2*x + y) % width
    ny := (x + y) % height
    ni := ny*stride + 4*nx
    
    // R := (*source)[i]
    // G := (*source)[i+1]
    // B := (*source)[i+2]
    // A := (*source)[i+3]
    
    // (*target)[ni] = uint8(A * R) // post-alpha compensation
    (*target)[ni] = uint8((*source)[i]) // post-alpha compensation
    // (*target)[ni+1] = uint8(A * G) // post-alpha compensation
    (*target)[ni+1] = uint8((*source)[i+1]) // post-alpha compensation
    // (*target)[ni+2] = uint8(A * B) // post-alpha compensation
    (*target)[ni+2] = uint8((*source)[i+2]) // post-alpha compensation
    // (*target)[ni+3] = uint8(A) // post-alpha compensation
    (*target)[ni+3] = uint8((*source)[i+3]) // post-alpha compensation
}

// func main() {
//     root := "c:/users/kenneth/pictures"
    
//     folders := pt.Folders(root)
//     fmt.Println(len(folders))
    
//     folder := folders[it.Randex(len(folders))]
//     fmt.Println(folder)
//     files := pt.Files(folder)
//     file := files[it.Randex(len(files))]
//     fmt.Println(file)
// }
// func main() {
//     folder := "E:/Projects/audio/krule_errors/slices"
    
//     files := pt.Files(folder)
//     file := files[it.Randex(len(files))]
//     fmt.Println(file)
//     pic := iot.LoadImage(file)
//     width := pic.Bounds().Max.X
//     height := pic.Bounds().Max.Y
//     fmt.Println(width, height)
//     iot.Pop(iot.RenderPNG(CatMapOnce(pic)))
//     // fmt.Println(Period(width, height))
// }
// func main() {
//     folder := "E:/Projects/audio/krule_errors/slices"
//     files := pt.Files(folder)
//     file := files[it.Randex(len(files))]
//     fmt.Println(file)
//     dir, period := FullArnold(file)
//     iot.Pop(dir)
//     fmt.Println(period)
// }
func handle(pth string) {
    fmt.Println(pth)
    dir, period := FullArnold(pth)
    iot.Pop(dir)
    fmt.Println(period)
}
func main() {
    if len(os.Args) != 2 {
        fmt.Println("Usage:\n \t./catlab image_file")
        os.Exit(0)
    }
    // file := os.Args[1]
    // fmt.Printf("%T", iot.LoadImage(file))
    // fmt.Printf("%T", iot.LoadImage(file).Pix)
    // return
    if len(os.Args[1:]) > 0 {
        for _, pth := range os.Args[1:] {
            switch {
                case pt.IsFile(pth):
                    handle(pth)
                case pt.IsDir(pth):
                    // files := pt.Files(pth)
                    // pth = files[it.Randex(len(files))]
                    // handle(pth)
                    handle(pt.SuitableImage(pth))
            }
        }
        return
    } 
    // files := pt.Files(pt.Pwd())
    // pth = files[it.Randex(len(files))]
    // handle(pth)
    handle(pt.SuitableImage(pt.Pwd()))
}
