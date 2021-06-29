package pathtools // catlab/pathtools
// package main

import (
    "fmt"
    "log"
    "os"
    "io/fs"
    "path"
    "strings"
    
    it"catlab/itertools"
    et"catlab/errortools"
)
func Cd(pth string) {
    err := os.Chdir(pth)
    et.Assert(err)
}
func Pwd() string {
    home, err := os.Getwd()
    et.Assert(err)
    return home
}
// func Normalize(path) {
    
// }
func IsDir(pth string) bool {
    f, err := os.Lstat(pth)
    et.Check(err)
    mode := f.Mode()
    return mode.IsDir()
}
func IsFile(pth string) bool {
    f, err := os.Lstat(pth)
    et.Check(err)
    mode := f.Mode()
    return mode.IsRegular()
}
func IsSymLink(pth string) bool {
    f, err := os.Lstat(pth)
    et.Check(err)
    mode := f.Mode()
    return mode&fs.ModeSymlink != 0
}
func IsNamedPipe(pth string) bool {
    f, err := os.Lstat(pth)
    et.Check(err)
    mode := f.Mode()
    return mode&fs.ModeNamedPipe != 0
}
func ExtensionSupported(pth string) bool {
    switch path.Ext(pth) {
        case ".jpeg":
            return true
        case ".jpg":
            return true
        case ".png":
            return true
        default:
            return false
    }
}
func SuitableImage(dir string) string {
    files := Files(dir)
    var ctr int
    var pth string
    finding:
        for ; !ExtensionSupported(pth); ctr++ {
            if ctr >= len(files) {
                log.Fatalf("Directory \"%s\" doesn't contain any suitable images", dir)
                break finding
            }
            pth = files[it.Randex(len(files))]
        }
    return pth
}
// func ExampleFileMode(path string) {
//     fi, err := os.Lstat(path)
//     if err != nil {
//         log.Fatal(err)
//     }

//     fmt.Printf("permissions: %#o\n", fi.Mode().Perm()) // 0400, 0777, etc.
//     switch mode := fi.Mode(); {
//     case mode.IsRegular():
//         fmt.Println("regular file")
//         fmt.Println(mode.IsRegular())
//     case mode.IsDir():
//         fmt.Println("directory")
//         fmt.Println(mode.IsDir())
//     case mode&fs.ModeSymlink != 0:
//         fmt.Println("symbolic link")
//     case mode&fs.ModeNamedPipe != 0:
//         fmt.Println("named pipe")
//     }
// }
func ProjDir(pth string) string {
    pwd, err := os.Getwd()
    et.Assert(err)
    pth = strings.Replace(pth, "\\", "/", -1)
    _, name := path.Split(pth)
    // fmt.Println(name)
    parts := strings.Split(name, ".")
    ext := parts[len(parts)-1]
    name = strings.Join(parts[:len(parts)-1], ".") + "-" + ext
    new := path.Join(pwd, "cat_maps", name)
    return NameSpacer(new)
}
func Files(root string) []string {
    paths := make([]string, 0)
    for _, name := range Listdir(root) {
        pth := path.Join(root, name)
        stat, err := os.Lstat(pth)
        if err != nil {
            log.Fatal(err)
        }
        switch mode := stat.Mode(); {
            case mode.IsRegular():
                paths = append(paths, pth)
            case mode.IsDir():
                go it.Merge(&paths, Files(pth))
        }
    }
    return paths
}
func Folders(root string) []string {
    paths := make([]string, 0)
    for _, name := range Listdir(root) {
        pth := path.Join(root, name)
        // fmt.Println(pth)
        stat, err := os.Lstat(pth)
        if err != nil {
            log.Fatal(err)
        }
        switch mode := stat.Mode(); {
            case mode.IsRegular():
                continue
            case mode.IsDir():
                paths = append(paths, pth)
                go it.Merge(&paths, Folders(pth))
        }
    }
    return paths
}
func Listdir(pth string) []string {
    files, err := os.ReadDir(pth)
    et.Check(err)
    rack := make([]string, 0)
    for _, file := range files {
        rack = append(rack, file.Name())
    }
    return rack
}
type namespacer struct {
    Format string
    Index int
}
func (ns *namespacer) space(pth string) string {
    _, err := os.Stat(pth)
    if err != nil {
        return pth
    } else {
        n := ns.new(pth)
        _, err = os.Stat(n)
        if err != nil {
            return n
        } else {
            ns.Index += 1
            return ns.space(pth)
        }
    }
}
func (ns namespacer) new(pth string) string {
    ext := path.Ext(pth)
    name := pth[:len(pth)-len(ext)]
    return (name + fmt.Sprintf(ns.Format, ns.Index) + ext)
}
func NameSpacer(pth string) string {
    ns := namespacer{"_%v", 2}
    return ns.space(pth)
}








func main() {
    folder := "c:/users/kenneth/pictures"
    folders := Folders(folder)
    it.Show(folders)
    // for _, f := range folders {
    //     fmt.Println(f)
    // }
    files := Files(folder)
    it.Show(files)
    // for _, f := range files {
    //     fmt.Println(f)
    // }
    // Show(Files(folder))
}


