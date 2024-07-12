package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

var reset = "\033[0m"
var red = "\033[31m"
var green = "\033[32m"
var yellow = "\033[33m"
var blue = "\033[34m"
var magenta = "\033[35m"
var cyan = "\033[36m"
var white = "\033[97m"
var orange = "\033[38;5;208m"
var purple = "\033[35m"

const maxFilesInLine = 3
const maxFileNameLength = 19

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

	if ext == "" {
		if file.IsDir() {
			fmt.Print(green + name + " " + reset)
		} else {
			fmt.Print(green + " " + reset + name)
		}
	} else {
		switch ext {
		case ".go":
			fmt.Print(cyan + " " + reset + name)
		case ".sh":
			fmt.Print(white + " " + reset + name)
		case ".cpp", ".hpp", ".cxx", ".hxx":
			fmt.Print(blue + " " + reset + name)
		case ".css":
			fmt.Print(blue + " " + reset + name)
		case ".c":
			fmt.Print(blue + " " + reset + name)
		case ".png", ".jpg", ".jpeg", ".webp":
			fmt.Print(magenta + " " + reset + name)
		case ".xcf":
			fmt.Print(white + " " + reset + name)
		case ".xml", ".htm", ".html":
			fmt.Print(red + " " + reset + name)
		case ".txt":
			fmt.Print(white + " " + reset + name)
		case ".mp3", ".ogg":
			fmt.Print(cyan + " " + reset + name)
		case ".zip", ".tar", ".gz", ".bz2", ".xz":
			fmt.Print(yellow + "󰿺 " + reset + name)
		case ".jar", ".java":
			fmt.Print(white + " " + reset + name)
		case ".js":
			fmt.Print(yellow + " " + reset + name)
		case ".py":
			fmt.Print(yellow + " " + reset + name)
		case ".rs":
			fmt.Print(orange + " " + reset + name)
		case ".deb":
			fmt.Print(red + " " + reset + name)
		case ".md":
			fmt.Print(blue + " " + reset + name)
		case ".rb":
			fmt.Print(red + " " + reset + name)
		case ".php":
			fmt.Print(purple + " " + reset + name)
		case ".pl":
			fmt.Print(orange + " " + reset + name)
		case ".svg", ".eps", ".ps":
			fmt.Print(magenta + " " + reset + name)
		case ".git":
			fmt.Print(orange + " " + reset + name)
		default:
			if isBinary(file) {
				fmt.Print(green + " " + reset + name)
			} else {
				fmt.Print(white + " " + reset + name)
			}
		}
	}
}

func printPadding(fileName string) {
	padding := maxFileNameLength - len(fileName)
	for i := 0; i < padding; i++ {
		fmt.Print(" ")
	}
}

func isBinary(file os.DirEntry) bool {
	binaryExtensions := []string{".exe", ".bin", ".o", ".so", ".dll", ".out"}
	ext := strings.ToLower(filepath.Ext(file.Name()))
	for _, binExt := range binaryExtensions {
		if ext == binExt {
			return true
		}
	}
	return false
}
