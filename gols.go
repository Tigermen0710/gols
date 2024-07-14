package main

import (
	"fmt"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"strings"
	"syscall"
)

// ANSI escape codes for colors
const (
	reset       = "\033[0m"
	red         = "\033[31m"
	green       = "\033[32m"
	yellow      = "\033[33m"
	blue        = "\033[34m"
	magenta     = "\033[35m"
	white       = "\033[97m"
	cyan        = "\033[36m"
	orange      = "\033[38;5;208m"
	purple      = "\033[35m"
	lightRed    = "\033[91m"
	lightPurple = "\033[95m"
	darkGreen   = "\033[38;5;22m"
	darkOrange  = "\033[38;5;208m"
	darkYellow  = "\033[38;5;172m"
	darkMagenta = "\033[38;5;125m"
)

var (
	longListing   bool
	humanReadable bool

	// File icons based on extensions
	fileIcons = map[string]string{
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
		".mp4":  cyan + " " + reset,
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
		".zig":  darkOrange + " " + reset,
		".xbps": darkGreen + " " + reset,
	}

	// Binary file extensions
	binaryExtensions = map[string]bool{
		".exe": true, ".bin": true, ".o": true, ".so": true, ".dll": true, ".out": true,
	}
)

func main() {
	parseFlags()

	directory := "."
	if len(os.Args) > 1 && !strings.HasPrefix(os.Args[len(os.Args)-1], "-") {
		directory = os.Args[len(os.Args)-1]
	}

	files, err := os.ReadDir(directory)
	if err != nil {
		log.Fatal(err)
	}

	if len(files) == 0 {
		fmt.Println("No files found.")
		return
	}

	if longListing {
		printLongListing(files, directory)
	} else {
		printFilesInColumns(files, directory)
	}
}

func parseFlags() {
	for _, arg := range os.Args[1:] {
		switch arg {
		case "-l":
			longListing = true
		case "-h":
			showHelp()
			os.Exit(0)
		case "-lh":
			longListing = true
			humanReadable = true
		default:
			if !strings.HasPrefix(arg, "-") {
				continue
			}
			showHelp()
			os.Exit(1)
		}
	}
}

func showHelp() {
	fmt.Println("Usage: gols [options] [directory]")
	fmt.Println("Options:")
	fmt.Println("  -l    Long listing format")
	fmt.Println("  -lh   Human-readable file sizes")
	fmt.Println("  -h    Show options")
}

func printFilesInColumns(files []os.DirEntry, directory string) {
	maxFilesInLine := 6
	maxFileNameLength := 19

	filesInLine := 0
	for _, file := range files {
		printFile(file, directory)
		filesInLine++
		if filesInLine >= maxFilesInLine || len(file.Name()) > maxFileNameLength {
			fmt.Println()
			filesInLine = 0
		} else {
			printPadding(file.Name(), maxFileNameLength)
		}
	}
	fmt.Println() // Ensure a newline at the end
}

func printLongListing(files []os.DirEntry, directory string) {
	for _, file := range files {
		info, err := file.Info()
		if err != nil {
			log.Fatal(err)
		}

		permissions := formatPermissions(info.Mode())
		size := info.Size()
		sizeStr := fmt.Sprintf("%d", size)
		if humanReadable {
			sizeStr = humanizeSize(size)
		}

		// Get owner and group names
		owner, err := user.LookupId(fmt.Sprintf("%d", info.Sys().(*syscall.Stat_t).Uid))
		if err != nil {
			log.Fatal(err)
		}
		group, err := user.LookupGroupId(fmt.Sprintf("%d", info.Sys().(*syscall.Stat_t).Gid))
		if err != nil {
			log.Fatal(err)
		}

		// Print long listing format with icons
		fmt.Printf("%s %10s %s %s", permissions, sizeStr, owner.Username, group.Name)
		fmt.Printf(" %s", info.ModTime().Format("Jan 02 15:04"))

		fmt.Printf(" %s %s\n", getFileIcon(file.Name()), file.Name())
	}
	fmt.Println() // Ensure a newline at the end
}

func formatPermissions(mode os.FileMode) string {
	var b strings.Builder

	if mode.IsDir() {
		b.WriteString("d")
	} else {
		b.WriteString("-")
	}

	b.WriteString(rwx(mode.Perm() >> 6)) // Owner permissions
	b.WriteString(rwx(mode.Perm() >> 3)) // Group permissions
	b.WriteString(rwx(mode.Perm()))      // Other permissions

	return b.String()
}

func rwx(perm os.FileMode) string {
	var b strings.Builder

	if perm&0400 != 0 {
		b.WriteString(green + "r")
	} else {
		b.WriteString("-")
	}
	if perm&0200 != 0 {
		b.WriteString(yellow + "w")
	} else {
		b.WriteString("-")
	}
	if perm&0100 != 0 {
		b.WriteString(red + "x")
	} else {
		b.WriteString("-")
	}

	b.WriteString(reset) // Reset colors

	return b.String()
}

func printFile(file os.DirEntry, directory string) {
	name := file.Name()
	ext := strings.ToLower(filepath.Ext(name))
	icon, exists := fileIcons[ext]

	if ext == "" {
		if file.IsDir() {
			fmt.Print(green + " " + reset + name)
		} else {
			fmt.Print(white + " " + reset + name)
		}
	} else if exists {
		fmt.Print(icon + name)
	} else if isBinary(ext) {
		fmt.Print(green + " " + reset + name)
	} else {
		fmt.Print(white + " " + reset + name)
	}

	// Print modification time only if long listing is specified and info.ModTime() is not zero
	if longListing && !file.IsDir() {
		info, err := os.Stat(filepath.Join(directory, file.Name()))
		if err == nil && !info.ModTime().IsZero() {
			fmt.Printf("  %s", info.ModTime().Format("Jan 02 15:04"))
		}
	}
}

func printPadding(fileName string, maxFileNameLength int) {
	padding := maxFileNameLength - len(fileName) + 2
	for i := 0; i < padding; i++ {
		fmt.Print(" ")
	}
}

func getFileIcon(fileName string) string {
	ext := strings.ToLower(filepath.Ext(fileName))

	if ext == "" {
		return green + " " + reset
	}

	icon, exists := fileIcons[ext]
	if exists {
		return icon
	}

	return white + " " + reset
}

func isBinary(ext string) bool {
	return binaryExtensions[ext]
}

func humanizeSize(size int64) string {
	const unit = 1024
	if size < unit {
		return fmt.Sprintf("%d B", size)
	}
	div, exp := int64(unit), 0
	for n := size / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(size)/float64(div), "KMGTPE"[exp])
}
