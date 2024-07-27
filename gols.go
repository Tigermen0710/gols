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
	longListing      bool
	humanReadable    bool
	fileSize         bool
    orderBySize      bool
    orderByTime      bool
    showOnlySymlinks bool
    showHidden       bool
    recursiveListing bool
    dirOnLeft	  	 bool

	// File icons based on extensions
	fileIcons = map[string]string{
		".go":   " ",
        ".mod":  " ",
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
		".JPG":  " ",
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
		".iso":  " ",
        ".exe":  " ",
        ".odt":  "󰷈 ",
        ".ods":  "󰱾 ",
        ".odp":  "󰈧 ",
        ".gif":  "󰵸 ",
        ".tiff": "󰋪 ",
        ".7z":   " ",
        ".bat":  " ",
        ".app":  " ",
        ".log":  " ",
        ".sql":  " ",
        ".db":   " ",
	}
)

func main() {
	args := os.Args[1:]
	nonFlagArgs, hasFlags, hasSpecificFlags := parseFlags(args)

	var directory string
	var fileExtension string

	if len(nonFlagArgs) > 0 {
		directory = nonFlagArgs[0]
	}
	if len(nonFlagArgs) > 1 {
		fileExtension = strings.TrimPrefix(filepath.Ext(nonFlagArgs[1]), ".")
	}

	if directory == "" {
		directory = "."
	}

	var files []os.DirEntry
	var err error

	if fileExtension != "" {
		files, err = listFilesWithExtension(directory, fileExtension)
		if err != nil {
			log.Fatalf("Error listing files with extension %s: %v", fileExtension, err)
		}
	} else {
		files, err = os.ReadDir(directory)
		if err != nil {
			log.Fatal(err)
		}
	}

	if len(files) == 0 {
		fmt.Println("No files found.")
		return
	}

	if !showHidden {
		files = filterHidden(files)
	}

	if showOnlySymlinks {
		files = filterSymlinks(files, directory)
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

	if recursiveListing {
		printTree(directory, "", true)
	} else if longListing {
		printLongListing(files, directory)
	} else if fileSize {
		getFileSize(files, directory)
	} else {
		printFilesInColumns(files, directory, dirOnLeft)
	}

	if (hasSpecificFlags && !longListing) || !hasFlags {
		fmt.Println()
	}
}

func listFilesWithExtension(dir string, ext string) ([]os.DirEntry, error) {
	var result []os.DirEntry

	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		if strings.HasSuffix(entry.Name(), "."+ext) {
			result = append(result, entry)
		}
	}
	return result, nil
}

func filterHidden(entries []os.DirEntry) []os.DirEntry {
	var result []os.DirEntry
	for _, entry := range entries {
		if !strings.HasPrefix(entry.Name(), ".") {
			result = append(result, entry)
		}
	}
	return result
}

func filterSymlinks(entries []os.DirEntry, dir string) []os.DirEntry {
	var result []os.DirEntry
	for _, entry := range entries {
		fullPath := filepath.Join(dir, entry.Name())
		if info, err := os.Lstat(fullPath); err == nil && info.Mode()&os.ModeSymlink != 0 {
			result = append(result, entry)
		}
	}
	return result
}

func parseFlags(args []string) ([]string, bool, bool) {
	var nonFlagArgs []string
	hasFlags := false
	hasSpecificFlags := false
	for i := 0; i < len(args); i++ {
		arg := args[i]
		if len(arg) > 1 && arg[0] == '-' {
			hasFlags = true
			for _, ch := range arg[1:] {
				switch ch {
				case 'l':
					longListing = true
				case 'h':
					humanReadable = true
				case 's':
					fileSize = true
				case 'o':
					orderBySize = true
					hasSpecificFlags = true
				case 't':
					orderByTime = true
					hasSpecificFlags = true
				case 'm':
					showOnlySymlinks = true
					hasSpecificFlags = true
				case 'a':
					showHidden = true
					hasSpecificFlags = true
				case 'r':
					recursiveListing = true
				case 'i':
					dirOnLeft = true
					hasSpecificFlags = true
				default:
					showHelp()
					os.Exit(1)
				}
			}
		} else {
			nonFlagArgs = append(nonFlagArgs, arg)
		}
	}
	return nonFlagArgs, hasFlags, hasSpecificFlags
}

func showHelp() {
	fmt.Println("Usage: gols [FLAG] [DIRECTORY]")
	fmt.Println()
	fmt.Println("FLAGS:")
	fmt.Println()
    fmt.Println("  -?        Options")
	fmt.Println()
	fmt.Println("  -l        Long listing format")
	fmt.Println("  -h        Human-readable file sizes")
	fmt.Println("  -s        Print files size")
    fmt.Println("  -o        Sort by size")
    fmt.Println("  -t        Order by time")
    fmt.Println("  -m        Only symbolic links are showing")
    fmt.Println("  -a        Show Hidden files")
    fmt.Println("  -r        Tree like listing")
    fmt.Println("  -i        Show directory icon on left")
	fmt.Println()
}

func printFilesInColumns(files []os.DirEntry, directory string, dirOnLeft bool) {
	maxFilesInLine := 4
	maxFileNameLength := 19

	filesInLine := 0
	for _, file := range files {
		if file.IsDir() && dirOnLeft {
			fmt.Print(blue + "  " + file.Name() + reset)
		} else {
			printFile(file, directory)
		}
		filesInLine++
		if filesInLine >= maxFilesInLine || len(file.Name()) > maxFileNameLength {
			fmt.Println()
			filesInLine = 0
		} else {
			printPadding(file.Name(), maxFileNameLength)
		}
	}
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
        	if dirOnLeft {
        		fmt.Println(blue + " " + file.Name() + reset)
        	} else {
     			fmt.Println(blue + file.Name() + " " + reset)
        	}
        } else {
            fmt.Println(getFileIcon(file, info.Mode(), directory)+ " " + file.Name())
        }
    }
}

func printLongListing(files []os.DirEntry, directory string) {
    maxLen := map[string]int{
        "permissions": 0,
        "size":        0,
        "owner":       0,
        "group":       0,
        "month":       0,
        "day":         0,
        "time":        0,
    }

    var filteredFiles []os.DirEntry
    if showOnlySymlinks {
        for _, file := range files {
            if file.Type()&os.ModeSymlink != 0 {
                filteredFiles = append(filteredFiles, file)
            }
        }
    } else {
        filteredFiles = files
    }

    for _, file := range filteredFiles {
        info, err := file.Info()
        if err != nil {
            log.Fatal(err)
        }

        permissions := formatPermissions(file, info.Mode(), directory)
        size := info.Size()
        sizeStr := fmt.Sprintf("%d", size)
        if humanReadable {
            sizeStr = humanizeSize(size)
        }

        owner, err := user.LookupId(fmt.Sprintf("%d", info.Sys().(*syscall.Stat_t).Uid))
        if err != nil {
            log.Fatal(err)
        }
        group, err := user.LookupGroupId(fmt.Sprintf("%d", info.Sys().(*syscall.Stat_t).Gid))
        if err != nil {
            log.Fatal(err)
        }

        modTime := info.ModTime()
        month := modTime.Format("Jan")
        day := fmt.Sprintf("%2d", modTime.Day())
        timeStr := modTime.Format("15:04:05 2006")

        if len(permissions) > maxLen["permissions"] {
            maxLen["permissions"] = len(permissions)
        }
        if len(sizeStr) > maxLen["size"] {
            maxLen["size"] = len(sizeStr)
        }
        if len(owner.Username) > maxLen["owner"] {
            maxLen["owner"] = len(owner.Username)
        }
        if len(group.Name) > maxLen["group"] {
            maxLen["group"] = len(group.Name)
        }
        if len(month) > maxLen["month"] {
            maxLen["month"] = len(month)
        }
        if len(day) > maxLen["day"] {
            maxLen["day"] = len(day)
        }
        if len(timeStr) > maxLen["time"] {
            maxLen["time"] = len(timeStr)
        }
    }

    for _, file := range filteredFiles {
        info, err := file.Info()
        if err != nil {
            log.Fatal(err)
        }

        permissions := formatPermissions(file, info.Mode(), directory)
        size := info.Size()
        sizeStr := fmt.Sprintf("%d", size)
        if humanReadable {
            sizeStr = humanizeSize(size)
        }

        owner, err := user.LookupId(fmt.Sprintf("%d", info.Sys().(*syscall.Stat_t).Uid))
        if err != nil {
            log.Fatal(err)
        }
        group, err := user.LookupGroupId(fmt.Sprintf("%d", info.Sys().(*syscall.Stat_t).Gid))
        if err != nil {
            log.Fatal(err)
        }

        modTime := info.ModTime()
        month := modTime.Format("Jan")
        day := fmt.Sprintf("%2d", modTime.Day())
        timeStr := modTime.Format("15:04:05 2006")

        // Print long listing format with icons and details
        line := fmt.Sprintf("%-*s %*s %-*s %-*s %-*s %-*s %-*s", maxLen["permissions"], permissions, maxLen["size"], sizeStr, maxLen["owner"], owner.Username, maxLen["group"], group.Name, maxLen["month"], month, maxLen["day"], day, maxLen["time"], timeStr)
        line += fmt.Sprintf(" %s %s%s", getFileIcon(file, info.Mode(), directory), file.Name(), reset)

        // Check if the file is a symbolic link
        if file.Type()&os.ModeSymlink != 0 {
            linkTarget, err := os.Readlink(filepath.Join(directory, file.Name()))
            if err == nil {
                line += fmt.Sprintf(" %s==> %s%s", cyan, linkTarget, reset)
            }
        }

        fmt.Println(line)
    }
}

func formatPermissions(file os.DirEntry, mode os.FileMode, directory string) string {
    perms := make([]byte, 10)
    for i := range perms {
        perms[i] = '-'
    }

    if file.Type()&os.ModeSymlink != 0 {
        linkTarget, err := os.Readlink(filepath.Join(directory, file.Name()))
        if err == nil {
            symlinkTarget := filepath.Join(directory, linkTarget)
            targetInfo, err := os.Stat(symlinkTarget)
            if err == nil && targetInfo.IsDir() {
                perms[0] = 'l'
                perms[1] = 'd'
            } else {
                perms[0] = 'l'
            }
        }
    } else if mode.IsDir() {
        perms[0] = 'd'
    }

    for i, s := range []struct {
        bit os.FileMode
        char byte
    }{
        {0400, 'r'}, {0200, 'w'}, {0100, 'x'},
        {0040, 'r'}, {0020, 'w'}, {0010, 'x'},
        {0004, 'r'}, {0002, 'w'}, {0001, 'x'},
    } {
        if mode&s.bit != 0 {
            perms[i+1] = s.char
        }
    }

    return string(perms)
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
		if perm&os.ModeSetuid != 0 {
			b.WriteString("s")
		} else {
			b.WriteString("x")
		}
	} else {
		if perm&os.ModeSetuid != 0 {
			b.WriteString("S")
		} else {
			b.WriteString("-")
		}
	}

	return b.String()
}

func printFile(file os.DirEntry, directory string) {
	info, err := file.Info()
	if err != nil {
		log.Fatal(err)
	}
	if file.IsDir() {
		fmt.Print(blue + file.Name() + " " + reset)
	} else {
		fmt.Print(getFileIcon(file, info.Mode(), directory) + file.Name())
	}
	fmt.Print(" ")
}

func getFileIcon(file os.DirEntry, mode os.FileMode, directory string) string {
	if file.Type()&os.ModeSymlink != 0 {
		linkTarget, err := os.Readlink(filepath.Join(directory, file.Name()))
		if err == nil {
			symlinkTarget := filepath.Join(directory, linkTarget)
			targetInfo, err := os.Stat(symlinkTarget)
			if err == nil && targetInfo.IsDir() {
				return brightMagenta + " " + reset // Symbolic link to directory icon
			} else {
				return brightCyan + " " + reset // Symbolic link to file icon
			}
		}
	}

	if mode.IsDir() {
		return blue + " " + reset
	}

	ext := filepath.Ext(file.Name())
	icon, exists := fileIcons[ext]
	if exists {
		switch ext {
		case ".go":
			return cyan + icon + reset
        case ".sh":
			if mode&os.ModePerm&0111 != 0 {
				return brightGreen + icon + reset
			} else {
				return white + icon + reset
			}
		case ".cpp", ".hpp", ".cxx", ".hxx":
			return blue + icon + reset
		case ".css":
			return lightblue + icon + reset
		case ".c", ".h":
			return blue + icon + reset
		case ".cs":
			return darkMagenta + icon + reset
		case ".png", ".jpg", ".jpeg", ".JPG", ".webp":
			return brightMagenta + icon + reset
		case ".gif":
			return magenta + icon + reset
		case ".xcf":
			return purple + icon + reset
		case ".xml":
			return lightCyan + icon + reset
		case ".htm", ".html":
			return orange + icon + reset
		case ".txt", ".app":
			return white + icon + reset
		case ".mp3", ".m4a", ".ogg", ".flac":
			return brightBlue + icon + reset
		case ".mp4", ".mkv", ".webm":
			return brightMagenta + icon + reset
		case ".zip", ".tar", ".gz", ".bz2", ".xz", ".7z":
			return lightPurple + icon + reset
		case ".jar", ".java":
			return orange + icon + reset
		case ".js":
			return yellow + icon + reset
		case ".json", ".tiff":
			return brightYellow + icon + reset
		case ".py":
			return darkYellow + icon + reset
		case ".rs":
			return darkGray + icon + reset
		case ".yml", ".yaml":
			return brightRed + icon + reset
		case ".toml":
			return darkOrange + icon + reset
		case ".deb":
			return lightRed + icon + reset
		case ".md":
			return cyan + icon + reset
		case ".rb":
			return red + icon + reset
		case ".php":
			return brightBlue + icon + reset
		case ".pl":
			return red + icon + reset
		case ".svg":
			return lightPurple + icon + reset
		case ".eps", ".ps":
			return orange + icon + reset
		case ".git":
			return orange + icon + reset
		case ".zig":
			return darkOrange + icon + reset
		case ".xbps":
			return darkGreen + icon + reset
		case ".el":
			return purple + icon + reset
		case ".vim":
			return darkGreen + icon + reset
		case ".lua", ".sql":
			return brightBlue + icon + reset
		case ".pdf", ".db":
			return brightRed + icon + reset
		case ".epub":
			return cyan + icon + reset
		case ".conf", ".bat":
			return darkGray + icon + reset
		case ".iso":
			return gray + icon + reset
		case ".exe":
			return brightCyan + icon + reset
		default:
			return icon
		}
	}

	if mode&os.ModePerm&0111 != 0 {
		return green + " " + reset
	}

	return " " + reset
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

func printTree(path, prefix string, isLast bool) {
	files, err := os.ReadDir(path)
	if err != nil {
		log.Fatal(err)
	}

	// Filter out hidden files if showHidden is false
	var filteredFiles []os.DirEntry
	for _, file := range files {
		if showHidden || !strings.HasPrefix(file.Name(), ".") {
			filteredFiles = append(filteredFiles, file)
		}
	}

	for i, file := range filteredFiles {
		isLastFile := i == len(filteredFiles)-1
		if isLastFile {
			fmt.Printf("%s└── ", prefix)
		} else {
			fmt.Printf("%s├── ", prefix)
		}

		printFile(file, path)
		fmt.Println()

		if file.IsDir() {
			newPrefix := prefix
			if isLastFile {
				newPrefix += "    "
			} else {
				newPrefix += "│   "
			}
			printTree(filepath.Join(path, file.Name()), newPrefix, isLastFile)
		}
	}
}
