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
        fmt.Print(getFileColor(file) + getFileIcon(file) + file.Name() + reset)
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

func getFileColor(file os.FileInfo) string {
    switch {
    case file.IsDir():
        return green
    case strings.Contains(file.Name(), ".go"):
        return cyan
    case strings.Contains(file.Name(), ".sh"):
        return white
    case strings.Contains(file.Name(), ".cpp"), strings.Contains(file.Name(), ".c"):
        return blue
    case strings.Contains(file.Name(), ".css"), strings.Contains(file.Name(), ".xml"), strings.Contains(file.Name(), ".htm"):
        return red
    case strings.Contains(file.Name(), ".png"), strings.Contains(file.Name(), ".jpg"), strings.Contains(file.Name(), ".webp"), strings.Contains(file.Name(), ".JPG"):
        return magenta
    case strings.Contains(file.Name(), ".xfc"):
        return white
    case strings.Contains(file.Name(), ".txt"):
        return white
    case strings.Contains(file.Name(), ".mp3"), strings.Contains(file.Name(), ".ogg"):
        return cyan
    case strings.Contains(file.Name(), ".zip"), strings.Contains(file.Name(), ".tar"):
        return yellow
    case strings.Contains(file.Name(), ".jar"), strings.Contains(file.Name(), ".java"):
        return white
    case strings.Contains(file.Name(), ".js"):
        return yellow
    case strings.Contains(file.Name(), ".py"):
        return yellow
    case strings.Contains(file.Name(), ".rs"):
        return white
    case strings.Contains(file.Name(), ".deb"):
        return red
    default:
        return white
    }
}

func getFileIcon(file os.FileInfo) string {
    switch {
    case file.IsDir():
        return " "
    case strings.Contains(file.Name(), ".go"):
        return " "
    case strings.Contains(file.Name(), ".sh"):
        return " "
    case strings.Contains(file.Name(), ".cpp"):
        return " "
    case strings.Contains(file.Name(), ".css"):
        return " "
    case strings.Contains(file.Name(), ".c"):
        return " "
    case strings.Contains(file.Name(), ".png"), strings.Contains(file.Name(), ".jpg"), strings.Contains(file.Name(), ".webp"), strings.Contains(file.Name(), ".JPG"):
        return " "
    case strings.Contains(file.Name(), ".xfc"):
        return " "
    case strings.Contains(file.Name(), ".xml"), strings.Contains(file.Name(), ".htm"):
        return " "
    case strings.Contains(file.Name(), ".txt"):
        return " "
    case strings.Contains(file.Name(), ".mp3"), strings.Contains(file.Name(), ".ogg"):
        return " "
    case strings.Contains(file.Name(), ".zip"), strings.Contains(file.Name(), ".tar"):
        return "󰿺 "
    case strings.Contains(file.Name(), ".jar"), strings.Contains(file.Name(), ".java"):
        return " "
    case strings.Contains(file.Name(), ".js"):
        return " "
    case strings.Contains(file.Name(), ".py"):
        return " "
    case strings.Contains(file.Name(), ".rs"):
        return " "
    case strings.Contains(file.Name(), ".deb"):
        return " "
    default:
        return " "
    }
}
