package main
import (
    "fmt"
    "strings"
    "log"
    "io/ioutil"
    "os"
)
var reset = "\033[0m"
var red = "\033[31m"
var green = "\033[32m"
var yellow = "\033[33m"
var blue = "\033[34m"
var magenta = "\033[35m"
var cyan = "\033[36m"
var white = "\033[97m"
var filesInLine = 0
func main()  {
    directory, err := os.Getwd()
    if err != nil {
        log.Fatal(err)
    }
    files, err := ioutil.ReadDir(directory)

    if err != nil {
        log.Fatal(err)
    }
    for _, file := range files {
        if file.IsDir() {
            fmt.Print(green + file.Name() + " " + reset)
        } else if strings.Contains(file.Name(), ".go") {
            fmt.Print(cyan + " " + reset + file.Name())
        } else if strings.Contains(file.Name(), ".sh") {
            fmt.Print(white + " " + reset + file.Name())
        } else if strings.Contains(file.Name(), ".cpp") {
            fmt.Print(blue + " " + reset + file.Name())
        } else if strings.Contains(file.Name(), ".css") {
            fmt.Print(blue + " " + reset + file.Name())
        } else if strings.Contains(file.Name(), ".c") {
            fmt.Print(blue + " " + reset + file.Name())
        } else if strings.Contains(file.Name(), ".png") || strings.Contains(file.Name(), ".jpg") || strings.Contains(file.Name(), ".webp") || strings.Contains(file.Name(), ".JPG") {
            fmt.Print(magenta + " " + reset + file.Name())
        } else if strings.Contains(file.Name(), ".xfc") {
            fmt.Print(white + " " + reset + file.Name())
        } else if strings.Contains(file.Name(), ".xml")|| strings.Contains(file.Name(), ".htm") {
            fmt.Print(red + " " + reset + file.Name())
        } else if strings.Contains(file.Name(), ".txt") {
            fmt.Print(white + " " + reset + file.Name())
        } else if strings.Contains(file.Name(), ".mp3") ||  strings.Contains(file.Name(), ".ogg") {
            fmt.Print(cyan + " " + reset + file.Name())
        } else if strings.Contains(file.Name(), ".zip")||strings.Contains(file.Name(), ".tar") {
            fmt.Print(yellow + "󰿺 " + reset + file.Name())
        } else if strings.Contains(file.Name(), ".jar")||strings.Contains(file.Name(), ".java") {
            fmt.Print(white + " " + reset + file.Name())
        } else if strings.Contains(file.Name(), ".js") {
            fmt.Print(yellow + " " + reset + file.Name())
        } else if strings.Contains(file.Name(), ".py") {
            fmt.Print(yellow + " " + reset + file.Name())
        } else if strings.Contains(file.Name(), ".rs") {
            fmt.Print(white + " " + reset + file.Name())
        }  else if strings.Contains(file.Name(), ".deb") {
            fmt.Print(red + " " + reset + file.Name())
        } else {
            fmt.Print(white + " " + reset + file.Name())
        }
        filesInLine++
        if filesInLine > 2 || len(file.Name()) > 19 {
            fmt.Println(" ")
            filesInLine = 0
        } else {
            for i := 0; i < 20 - len(file.Name()); i++ {
                fmt.Print(" ")
            }
        }
    }
    fmt.Println("")
}
