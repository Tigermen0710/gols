package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

var (
	reset   = "\033[0m"
	red     = "\033[31m"
	green   = "\033[32m"
	yellow  = "\033[33m"
	blue    = "\033[34m"
	magenta = "\033[35m"
	cyan    = "\033[36m"
	white   = "\033[97m"
	orange  = "\033[38;5;208m"
	purple  = "\033[35m"
)

const (
	maxFilesInLine     = 3
	maxFileNameLength  = 19
)

var fileIcons = map[string]string{
	".go":   cyan + " " + reset,
	".sh":   white + " " + reset,
	".cpp":  blue + " " + reset,
	".hpp":  blue + " " + reset,
	".cxx":  blue + " " + reset,
	".hxx":  blue + " " + reset,
	".css":  blue + " " + reset,
	".c":    blue + " " + reset,
	".png":  magenta + " " + reset,
	".jpg":  magenta + " " + reset,
	".jpeg": magenta + " " + reset,
	".webp": magenta + " " + reset,
	".xcf":  white + " " + reset,
	".xml":  red + " " + reset,
	".htm":  red + " " + reset,
	".html": red + " " + reset,
	".txt":  white + " " + reset,
	".mp3":  cyan + " " + reset,
	".ogg":  cyan + " " + reset,
	".zip":  yellow + "󰿺 " + reset,
	".tar":  yellow + "󰿺 " + reset,
	".gz":   yellow + "󰿺 " + reset,
	".bz2":  yellow + "󰿺 " + reset,
	".xz":   yellow + "󰿺 " + reset,
	".jar":  white + " " + reset,
	".java": white + " " + reset,
	".js":   yellow + " " + reset,
	".py":   yellow + " " + reset,
	".rs":   orange + " " + reset,
	".deb":  red + " " + reset,
	".md":   blue + " " + reset,
	".rb":   red + " " + reset,
	".php":  purple + " " + reset,
	".pl":   orange + " " + reset,
	".svg":  magenta + " " + reset,
	".eps":  magenta + " " + reset,
	".ps":   magenta + " " + reset,
	".git":  orange + " " + reset,
}

var binaryExtensions = map[string]bool{
	".exe": true, ".bin": true, ".o": true, ".so": true, ".dll": true, ".out": true,
}

func main() {
	directory, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	files, err := os.ReadDir(directory)
	if err != nil {
		log.Fatal(err)
	}

	filesInLine := 0
	for _, file := range files {
		printFile(file)
		filesInLine++
		if filesInLine >= maxFilesInLine || len(file.Name()) > maxFileNameLength {
			fmt.Println()
			filesInLine = 0
		} else {
			printPadding(file.Name())
		}
	}
}

func printFile(file os.DirEntry) {
	name := file.Name()
	ext := strings.ToLower(filepath.Ext(name))
	icon, exists := fileIcons[ext]

	if ext == "" {
		if file.IsDir() {
			fmt.Print(green + name + " " + reset)
		} else {
			fmt.Print(green + " " + reset + name)
		}
	} else if exists {
		fmt.Print(icon + name)
	} else if isBinary(ext) {
		fmt.Print(green + " " + reset + name)
	} else {
		fmt.Print(white + " " + reset + name)
	}
}

func printPadding(fileName string) {
	padding := maxFileNameLength - len(fileName)
	for i := 0; i < padding; i++ {
		fmt.Print(" ")
	}
}

func isBinary(ext string) bool {
	return binaryExtensions[ext]
}
