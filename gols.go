package main

import (
	"fmt"
	"log"
	"os"
	"os/user"
	"path/filepath"
    "sort"
	"strings"
	"syscall"
)

// ANSI escape codes for colors
const (
	reset         = "\033[0m"
	green         = "\033[32m"
	red           = "\033[31m"
	yellow        = "\033[33m"
	blue          = "\033[34m"
	magenta       = "\033[35m"
	white         = "\033[97m"
	cyan          = "\033[36m"
	orange        = "\033[38;5;208m"
	purple        = "\033[35m"
    gray          = "\033[37m"
    lightRed      = "\033[91m"
	lightgreen    = "\033[92m"
	lightyellow   = "\033[93m"
	lightblue     = "\033[94m"
	lightPurple   = "\033[95m"
    lightCyan     = "\033[38;5;87m"
	darkGreen     = "\033[38;5;22m"
	darkOrange    = "\033[38;5;208m"
	darkYellow    = "\033[38;5;172m"
	darkMagenta   = "\033[38;5;125m"
    darkGray      = "\033[90m"
    brightRed     = "\033[38;5;196m"
    brightGreen   = "\033[38;5;46m"
    brightYellow  = "\033[38;5;226m"
    brightBlue    = "\033[38;5;39m"
    brightMagenta = "\033[38;5;198m"
    brightCyan    = "\033[38;5;51m"
    brightWhite   = "\033[97m"
)

var (
	longListing   bool
	humanReadable bool
	fileSize      bool
    orderBySize   bool
    orderByTime   bool

	// File icons based on extensions
	fileIcons = map[string]string{
		".go":   " ",
		".sh":   " ",
		".cpp":  " ",
		".hpp":  " ",
		".cxx":  " ",
		".hxx":  " ",
		".css":  " ",
		".c":    " ",
		".h":    " ",
		".cs":   "󰌛 ",
		".png":  " ",
		".jpg":  " ",
		".jpeg": " ",
		".webp": " ",
		".xcf":  " ",
		".xml":  "󰗀 ",
		".htm":  " ",
		".html": " ",
		".txt":  " ",
		".mp3":  " ",
		".m4a":  " ",
		".ogg":  " ",
		".flac": " ",
		".mp4":  " ",
		".mkv":  " ",
		".webm": " ",
		".zip":  "󰿺 ",
		".tar":  "󰛫 ",
		".gz":   "󰛫 ",
		".bz2":  "󰿺 ",
		".xz":   "󰿺 ",
		".jar":  " ",
		".java": " ",
		".js":   " ",
		".json": " ",
		".py":   " ",
		".rs":   " ",
		".yml":  " ",
		".yaml": " ",
		".toml": " ",
		".deb":  " ",
		".md":   " ",
		".rb":   " ",
		".php":  " ",
		".pl":   " ",
		".svg":  "󰜡 ",
		".eps":  " ",
		".ps":   " ",
		".git":  " ",
		".zig":  " ",
		".xbps": " ",
		".el":   " ",
		".vim":  " ",
		".lua":  " ",
		".pdf":  " ",
		".epub": "󰂺 ",
		".conf": " ",
		".iso":  "󰗮 ",
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
    if orderBySize {
		sort.Slice(files, func(i, j int) bool {
			info1, _ := files[i].Info()
			info2, _ := files[j].Info()
			return info1.Size() < info2.Size()
		})
	}
    if orderByTime {
    sort.Slice(files, func(i, j int) bool {
        info1, _ := files[i].Info()
        info2, _ := files[j].Info()
        return info1.ModTime().Before(info2.ModTime())
    })
    }
	if longListing {
		printLongListing(files, directory)
	} else if fileSize {
		getFileSize(files, directory)
	} else {
		printFilesInColumns(files, directory)
	}
}

func parseFlags() {
	for i := 1; i < len(os.Args); i++ {
		arg := os.Args[i]
		switch arg {
		case "-l":
			longListing = true
		case "-":
			showHelp()
			os.Exit(0)
		case "-lh", "-hl":
			longListing = true
			humanReadable = true
		case "-s":
			fileSize = true
		case "-hs", "-sh":
			fileSize = true
			humanReadable = true
        case "-o":
            orderBySize = true
            longListing = true
            humanReadable = true
        case "-t":
            orderByTime = true
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
	fmt.Println("  -l        Long listing format")
	fmt.Println("  -lh       Human-readable file sizes")
	fmt.Println("  -hl       Human-readable file sizes")
	fmt.Println("  -s        Print files size")
	fmt.Println("  -hs       Print files size human-readable")
	fmt.Println("  -sh       Print files size human-readable")
    fmt.Println("  -o        Sort by size")
    fmt.Println("  -t        Order by time")
    fmt.Println("  -         Show options")
}

func printFilesInColumns(files []os.DirEntry, directory string) {
	maxFilesInLine := 3
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

func getFileSize(files []os.DirEntry, directory string) {
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
		var spaces = 10 - len(sizeStr)
		fmt.Print(sizeStr)
		for i := 0; i < spaces; i++ {
			fmt.Print(" ")
		}
		if file.IsDir() {
			fmt.Println(blue + file.Name() + "  " + reset)
		} else {
			fmt.Println(getFileIcon(file, info.Mode()) + file.Name())
		}
	}
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

		fmt.Printf(" %s %s\n", getFileIcon(file, info.Mode()), file.Name())
	}
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
		b.WriteString("r")
	} else {
		b.WriteString("-")
	}
	if perm&0200 != 0 {
		b.WriteString("w")
	} else {
		b.WriteString("-")
	}
	if perm&0100 != 0 {
		b.WriteString("x")
	} else {
		b.WriteString("-")
	}

	return b.String()
}

func printFile(file os.DirEntry, directory string) {
	info, err := file.Info()
	if err != nil {
		log.Fatal(err)
	}
	if file.IsDir() {
		fmt.Print(blue + file.Name() + "  " + reset)
	} else {
		fmt.Print(getFileIcon(file, info.Mode()) + file.Name())
	}
	fmt.Print(" ")
}

func getFileIcon(file os.DirEntry, mode os.FileMode) string {
	if mode.IsDir() {
		return blue + " " + reset // Directory icon
	}

	ext := filepath.Ext(file.Name())
	icon, exists := fileIcons[ext]
	if exists {
		return getIconColor(icon, mode)
	}
	// Default icon for files without known extensions
	if mode&0111 != 0 {
		return green + " " + reset // Executable file icon
	}
	return " " + reset // Regular file icon
}

func getIconColor(icon string, mode os.FileMode) string {
	// For executable files, return green icon
	if mode&0111 != 0 {
		return green + icon + reset
	}
	// For other icons, return default color based on extension
	switch icon {
	case " ":
		return cyan + icon + reset // .go files
	case " ":
		return white + icon + reset // .sh files
	case " ":
		return blue + icon + reset // .cpp, .hpp, .cxx, .hxx files
	case " ":
		return lightblue + icon + reset // .css files
	case " ":
		return blue + icon + reset // .c .h files
    case "󰌛 ":
        return darkMagenta + icon + reset // .cs files
	case " ":
		return lightRed + icon + reset // .png, .jpg, .jpeg, .webp files
	case " ":
		return purple + icon + reset // .xcf files
	case " ":
		return orange + icon + reset // .htm files
    case "󰗀 ":
        return lightCyan + icon + reset // .xml
    case " ":
        return orange + icon + reset // .html
    case " ":
        return yellow + icon + reset // .flac
	case " ":
		return white + icon + reset // .txt files
	case " ":
		return brightBlue + icon + reset // .mp3, .ogg files
	case " ":
		return brightMagenta + icon + reset // .mp4 .mp4 .webm files
	case "󰿺 ":
		return brightYellow + icon + reset // .zip, .bz2, .xz files
    case "󰛫 ":
		return yellow + icon + reset // .tar .gz files
	case " ":
		return orange + icon + reset // .jar, .java files
	case " ":
		return yellow + icon + reset // .js files
    case " ":
		return brightYellow + icon + reset // .json files
	case " ":
		return darkYellow + icon + reset // .py files
	case " ":
		return darkGray + icon + reset // .rs files
    case " ":
        return brightRed + icon + reset // .yml .yaml files
    case " ":
        return darkOrange + icon + reset // .toml files
	case " ":
		return red + icon + reset // .deb files
	case " ":
		return cyan + icon + reset // .md files
	case " ":
		return red + icon + reset // .rb files
	case " ":
		return brightBlue + icon + reset // .php files
	case " ":
		return red + icon + reset // .pl files
	case " ":
		return orange + icon + reset // .eps, .ps files
    case "󰜡 ":
        return orange + icon + reset // .svg files 
	case " ":
		return orange + icon + reset // .git files
	case " ":
		return darkOrange + icon + reset // .zig files
	case " ":
		return darkGreen + icon + reset // .xbps files
    case "i":
        return purple + icon + reset // .el files
    case " ":
        return green + icon + reset // .vim files
    case " ":
        return blue + icon + reset // .lua files
    case " ":
        return red + icon + reset // .pdf files
    case "󰂺 ":
        return blue + icon + reset // .epub files
    case " ":
        return gray + icon + reset // .conf files
    case "󰗮 ":
        return gray + icon + reset // .iso files
	default:
		return icon + reset // Default to icon without color for unknown extensions
	}
}

func humanizeSize(size int64) string {
	const (
		_  = iota
		KB = 1 << (10 * iota)
		MB
		GB
		TB
	)
	switch {
	case size >= TB:
		return fmt.Sprintf("%.2fTB", float64(size)/float64(TB))
	case size >= GB:
		return fmt.Sprintf("%.2fGB", float64(size)/float64(GB))
	case size >= MB:
		return fmt.Sprintf("%.2fMB", float64(size)/float64(MB))
	case size >= KB:
		return fmt.Sprintf("%.2fKB", float64(size)/float64(KB))
	default:
		return fmt.Sprintf("%dB", size)
	}
}

func printPadding(name string, maxFileNameLength int) {
	padding := maxFileNameLength - len(name)
	fmt.Print(strings.Repeat(" ", padding))
}

func getFileNameAndExtension(file os.DirEntry) (string, string) {
	ext := filepath.Ext(file.Name())
	name := strings.TrimSuffix(file.Name(), ext)
	return name, ext
}
