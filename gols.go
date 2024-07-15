package main

import (
    "fmt"
    "io/ioutil"
    "log"
    "os"
    "path/filepath"
    "strings"
)

const (
    reset   = "\033[0m"
    red     = "\033[31m"
    green   = "\033[32m"
    yellow  = "\033[33m"
    blue    = "\033[34m"
    magenta = "\033[35m"
    cyan    = "\033[36m"
    white   = "\033[97m"
)

var filesInLine = 0

func main() {
    directory, err := os.Getwd()
    if err != nil {
        log.Fatal(err)
    }

    sanitizedDir, err := sanitizePath(directory)
    if err != nil {
        log.Fatal(err)
    }

    files, err := ioutil.ReadDir(sanitizedDir)
    if err != nil {
        log.Fatal(err)
    }

    for _, file := range files {
        switch {
        case file.IsDir():
            fmt.Print(green + file.Name() + " " + reset)
        case strings.Contains(file.Name(), ".go"):
            fmt.Print(cyan + " " + reset + file.Name())
        case strings.Contains(file.Name(), ".sh"):
            fmt.Print(white + " " + reset + file.Name())
        case strings.Contains(file.Name(), ".cpp"):
            fmt.Print(blue + " " + reset + file.Name())
        case strings.Contains(file.Name(), ".css"):
            fmt.Print(blue + " " + reset + file.Name())
        case strings.Contains(file.Name(), ".c"):
            fmt.Print(blue + " " + reset + file.Name())
        case strings.Contains(file.Name(), ".png"), strings.Contains(file.Name(), ".jpg"), strings.Contains(file.Name(), ".webp"), strings.Contains(file.Name(), ".JPG"):
            fmt.Print(magenta + " " + reset + file.Name())
        case strings.Contains(file.Name(), ".xfc"):
            fmt.Print(white + " " + reset + file.Name())
        case strings.Contains(file.Name(), ".xml"), strings.Contains(file.Name(), ".htm"):
            fmt.Print(red + " " + reset + file.Name())
        case strings.Contains(file.Name(), ".txt"):
            fmt.Print(white + " " + reset + file.Name())
        case strings.Contains(file.Name(), ".mp3"), strings.Contains(file.Name(), ".ogg"):
            fmt.Print(cyan + " " + reset + file.Name())
        case strings.Contains(file.Name(), ".zip"), strings.Contains(file.Name(), ".tar"):
            fmt.Print(yellow + "󰿺 " + reset + file.Name())
        case strings.Contains(file.Name(), ".jar"), strings.Contains(file.Name(), ".java"):
            fmt.Print(white + " " + reset + file.Name())
        case strings.Contains(file.Name(), ".js"):
            fmt.Print(yellow + " " + reset + file.Name())
        case strings.Contains(file.Name(), ".py"):
            fmt.Print(yellow + " " + reset + file.Name())
        case strings.Contains(file.Name(), ".rs"):
            fmt.Print(white + " " + reset + file.Name())
        case strings.Contains(file.Name(), ".deb"):
            fmt.Print(red + " " + reset + file.Name())
        default:
            fmt.Print(white + " " + reset + file.Name())
        }
        filesInLine++
        if filesInLine > 2 || len(file.Name()) > 19 {
            fmt.Println()
            filesInLine = 0
        } else {
            fmt.Print(strings.Repeat(" ", 20-len(file.Name())))
        }
    }
    fmt.Println("")
}

func sanitizePath(path string) (string, error) {
    cleanPath := filepath.Clean(path)
    if !filepath.IsAbs(cleanPath) {
        absPath, err := filepath.Abs(cleanPath)
        if err != nil {
            return "", err
        }
        cleanPath = absPath
    }

    if strings.Contains(cleanPath, "..") {
        return "", fmt.Errorf("invalid path: %s", cleanPath)
    }
    return cleanPath, nil
}
