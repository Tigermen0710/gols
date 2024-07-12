package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
    //"time"
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
	maxFilesInLine    = 6
	maxFileNameLength = 19
)

var (
	longListing   bool
	humanReadable bool
	fileIcons     = map[string]string{
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
	binaryExtensions = map[string]bool{
		".exe": true, ".bin": true, ".o": true, ".so": true, ".dll": true, ".out": true,
	}
)

func main() {
    defer func() {
        fmt.Println(" ")
        fmt.Print("\033[0m") // Reset ANSI codes
        os.Stdout.Sync()
    }()

    parseFlags()

    directory, err := os.Getwd()
    if err != nil {
        log.Fatal(err)
    }

    files, err := os.ReadDir(directory)
    if err != nil {
        log.Fatal(err)
    }

    if longListing {
        printLongListing(files)
    } else {
        printFilesInColumns(files)
    }
}

func parseFlags() {
	for _, arg := range os.Args[1:] {
		if strings.HasPrefix(arg, "-") {
			for _, flag := range arg[1:] {
				switch flag {
				case 'l':
					longListing = true
				case 'h':
					if len(arg) == 2 { // If the flag is exactly "-h"
						showHelp()
						os.Exit(0)
					} else {
						humanReadable = true
					}
				default:
					log.Fatalf("unknown flag: %c", flag)
				}
			}
		}
	}
}

func showHelp() {
	fmt.Println("Usage: gols [options]")
	fmt.Println("Options:")
	fmt.Println("  -l    Long listing format")
	fmt.Println("  -h    Human-readable file sizes")
	fmt.Println("        Shows this help message if used alone")
}

func printFilesInColumns(files []os.DirEntry) {
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

func printLongListing(files []os.DirEntry) {
	for _, file := range files {
		info, err := file.Info()
		if err != nil {
			log.Fatal(err)
		}
		size := info.Size()
		sizeStr := fmt.Sprintf("%d", size)
		if humanReadable {
			sizeStr = humanizeSize(size)
		}
		modTime := info.ModTime().Format("Jan 02 15:04")
		fmt.Printf("%s %10s %s %s\n", modTime, sizeStr, getFileIcon(file), file.Name())
	}
}

func getFileIcon(file os.DirEntry) string {
	name := file.Name()
	ext := strings.ToLower(filepath.Ext(name))

	if ext == "" {
		if file.IsDir() {
			return green + "" + reset
		}
		return green + "" + reset
	}
	if icon, exists := fileIcons[ext]; exists {
		return icon
	}
	if isBinary(ext) {
		return green + "" + reset
	}
	return white + "" + reset
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

func humanizeSize(size int64) string {
	const unit = 1024
	if size < unit {
		return fmt.Sprintf("%dB", size)
	}
	div, exp := int64(unit), 0
	for n := size / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f%cB", float64(size)/float64(div), "KMGTPE"[exp])
}
