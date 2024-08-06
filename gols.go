package main

import (
	"fmt"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"syscall"
)

const (
	version = "gols: 1.3.3"
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
	dirOnLeft        bool
	oneColumn        bool
	showSummary      bool
	showVersion      bool
	maxDepth         int = -1
)

func main() {
	args := os.Args[1:]
	nonFlagArgs, hasFlags, hasSpecificFlags := parseFlags(args)

	if showVersion {
		fmt.Println(version)
		return
	}

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
		printTree(directory, "", true, 0, maxDepth)
	} else if longListing {
		printLongListing(files, directory, humanReadable)
	} else if fileSize {
		getFileSize(files, directory, humanReadable, dirOnLeft)
	} else {
		printFilesInColumns(files, directory, dirOnLeft, showSummary, oneColumn)
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
			for j := 1; j < len(arg); j++ {
				switch arg[j] {
				case 'l':
					longListing = true
				case 'h':
					humanReadable = true
                    hasSpecificFlags = true
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
				case 'c':
					oneColumn = true
				case 'f':
					showSummary = true
				case 'v':
					showVersion = true
				case 'd':
					if j+1 < len(arg) && arg[j+1] >= '0' && arg[j+1] <= '9' {
						depthValue := arg[j+1:]
						maxDepthValue, err := strconv.Atoi(depthValue)
						if err != nil {
							fmt.Println("Invalid value for -d")
							os.Exit(1)
						}
						maxDepth = maxDepthValue
						hasSpecificFlags = true
						break
					} else if i+1 < len(args) {
						depthValue := args[i+1]
						maxDepthValue, err := strconv.Atoi(depthValue)
						if err != nil {
							fmt.Println("Invalid value for -d")
							os.Exit(1)
						}
						maxDepth = maxDepthValue
						hasSpecificFlags = true
						i++
						break
					} else {
						fmt.Println("Missing value for -d")
						os.Exit(1)
					}
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
	fmt.Println()
	fmt.Println("Usage: gols [FLAG] [DIRECTORY] [FILES]")
	fmt.Println()
	fmt.Println("FLAGS:")
	fmt.Println()
	fmt.Println("	-?        Help")
	fmt.Println()
	fmt.Println("	-a        Show Hidden files")
	fmt.Println("	-c        Don't use spacing, print all files in one column")
	fmt.Println("	-f        Show summary of directories and files")
	fmt.Println("	-h        Human-readable file sizes")
	fmt.Println("	-i        Show directory icon on left")
	fmt.Println("	-l        Long listing format")
	fmt.Println("	-m        Only symbolic links are showing")
	fmt.Println("	-o        Sort by size")
	fmt.Println("	-r d n    Tree like listing, set the depth of the directory tree (n is an integer)")
	fmt.Println("	-s        Print files size")
	fmt.Println("	-t        Order by time")
	fmt.Println("	-v        Version")
	fmt.Println()
}

func printFilesInColumns(files []os.DirEntry, directory string, dirOnLeft bool, showSummary bool, oneColumn bool) {
    const (
        maxFilesInLine = 4
        maxFileNameLength = 19
    )

    filesInLine := 0
    dirCount := 0
    fileCount := 0

    for _, file := range files {
        if file.IsDir() {
            dirCount++
            if dirOnLeft {
                fmt.Print(blue + "  " + file.Name() + reset)
            } else {
                printFile(file, directory)
            }
        } else {
            fileCount++
            printFile(file, directory)
        }

        if !oneColumn {
            filesInLine++
            if filesInLine >= maxFilesInLine || len(file.Name()) > maxFileNameLength {
                fmt.Println()
                filesInLine = 0
            } else {
                printPadding(file.Name(), maxFileNameLength)
            }
        } else {
            fmt.Println()
        }
    }

    if showSummary {
        fmt.Println()
        dirCount, fileCount := countFilesAndDirs(files)
        fmt.Printf(iconDirectory + " Directories: %s%d%s\n", blue, dirCount, reset)
        fmt.Printf(iconOther + " Files: %s%d%s\n", red, fileCount, reset)
    }
}

func getFileSize(files []os.DirEntry, directory string, humanReadable, dirOnLeft bool) {
	const sizeFieldWidth = 10
	const spaceBetweenSizeAndIcon = 2

	for _, file := range files {
		info, err := file.Info()
		if err != nil {
			log.Fatal(err)
		}

		size := info.Size()
		sizeStr := formatSize(size, humanReadable)

		sizeStr = fmt.Sprintf("%*s", sizeFieldWidth, sizeStr)

		fmt.Print(sizeStr)
		for i := 0; i < spaceBetweenSizeAndIcon; i++ {
			fmt.Print(" ")
		}

		if file.IsDir() {
			if dirOnLeft {
				fmt.Println(iconDirectory + " " + blue + file.Name() + reset)
			} else {
				fmt.Println(blue + file.Name() + " " + iconDirectory + " " + reset)
			}
		} else {
			fmt.Println(getFileIcon(file, info.Mode(), directory) + " " + file.Name())
		}
	}

	if showSummary {
		fileCount, dirCount := countFilesAndDirs(files)
		fmt.Printf(iconDirectory + " Directories: %s%d%s\n", blue, dirCount, reset)
		fmt.Printf(iconOther + " Files: %s%d%s\n", red, fileCount, reset)
	}
}

func padRight(str string, length int) string {
	for len(str) < length {
		str += " "
	}
	return str
}

func formatSize(size int64, humanReadable bool) string {
	const (
		_  = iota
		KB = 1 << (10 * iota)
		MB
		GB
		TB
	)

	if humanReadable {
		switch {
		case size >= TB:
			return fmt.Sprintf("%.2f TB", float64(size)/float64(TB))
		case size >= GB:
			return fmt.Sprintf("%.2f GB", float64(size)/float64(GB))
		case size >= MB:
			return fmt.Sprintf("%.2f MB", float64(size)/float64(MB))
		case size >= KB:
			return fmt.Sprintf("%.2f KB", float64(size)/float64(KB))
		default:
			return fmt.Sprintf("%d B", size)
		}
	} else {
		switch {
		case size >= TB:
			return fmt.Sprintf("%d TB", size)
		case size >= GB:
			return fmt.Sprintf("%d GB", size)
		case size >= MB:
			return fmt.Sprintf("%d MB", size)
		case size >= KB:
			return fmt.Sprintf("%d KB", size)
		default:
			return fmt.Sprintf("%d B", size)
		}
	}
}

func printLongListing(files []os.DirEntry, directory string, humanReadable bool) {
	maxLen := map[string]int{
		"permissions": 0,
		"size":        0,
		"owner":       0,
		"group":       0,
		"month":       0,
		"day":         0,
		"time":        0,
		"linkTarget":  0,
	}

	var filteredFiles []os.DirEntry
	for _, file := range files {
		info, err := file.Info()
		if err != nil {
			log.Fatal(err)
		}

		permissions := formatPermissions(file, info.Mode(), directory)
		size := info.Size()
		sizeStr := formatSize(size, humanReadable)
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

		maxLen["permissions"] = max(maxLen["permissions"], len(permissions))
		maxLen["size"] = max(maxLen["size"], len(sizeStr))
		maxLen["owner"] = max(maxLen["owner"], len(owner.Username))
		maxLen["group"] = max(maxLen["group"], len(group.Name))
		maxLen["month"] = max(maxLen["month"], len(month))
		maxLen["day"] = max(maxLen["day"], len(day))
		maxLen["time"] = max(maxLen["time"], len(timeStr))

		if file.Type()&os.ModeSymlink != 0 {
			linkTarget, err := os.Readlink(filepath.Join(directory, file.Name()))
			if err == nil {
				maxLen["linkTarget"] = max(maxLen["linkTarget"], len(linkTarget)+5)
			}
		}

		filteredFiles = append(filteredFiles, file)
	}

	for _, file := range filteredFiles {
		info, err := file.Info()
		if err != nil {
			log.Fatal(err)
		}

		permissions := formatPermissions(file, info.Mode(), directory)
		size := info.Size()
		sizeStr := formatSize(size, humanReadable)
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

		permissions = green + permissions + reset
		sizeStr = fmt.Sprintf("%*s", maxLen["size"], sizeStr)
		ownerStr := cyan + owner.Username + reset
		groupStr := brightBlue + group.Name + reset
		monthStr := magenta + month + reset
		dayStr := magenta + day + reset
		timeStr = magenta + timeStr + reset

		line := fmt.Sprintf(
			"%-*s  %s  %-*s  %-*s %-*s %-*s %-*s %s %s",
			maxLen["permissions"], permissions,
			sizeStr,
			maxLen["owner"], ownerStr,
			maxLen["group"], groupStr,
			maxLen["month"], monthStr,
			maxLen["day"], dayStr,
			maxLen["time"], timeStr,
			getFileIcon(file, info.Mode(), directory), file.Name(),
		)

		if file.Type()&os.ModeSymlink != 0 {
			linkTarget, err := os.Readlink(filepath.Join(directory, file.Name()))
			if err == nil {
				line += fmt.Sprintf(" %s==> %s%s", cyan, linkTarget, reset)
			}
		}

		fmt.Println(line)
	}

    if showSummary {
		fileCount, dirCount := countFilesAndDirs(files)
		fmt.Printf(iconDirectory + " Directories: %s%d%s\n", blue, dirCount, reset)
		fmt.Printf(iconOther + " Files: %s%d%s\n", red, fileCount, reset)
	}
}

func countFilesAndDirs(files []os.DirEntry) (int, int) {
	fileCount := 0
	dirCount := 0
	for _, file := range files {
		if file.IsDir() {
			dirCount++
		} else {
			fileCount++
		}
	}
	return fileCount, dirCount
}

func colorize(char byte) string {
	switch char {
	case 'r':
		return reset + magenta + string(char) + reset
	case 'w':
		return lightGreen + string(char) + reset
	case 'x':
		return reset + red + string(char) + reset
	case 'd':
		return blue + string(char) + reset
	case 'l':
		return brightCyan + string(char) + reset
	case '-':
		return reset + lightWhite + string(char) + reset
	default:
		return string(char)
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
		bit  os.FileMode
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

	coloredPerms := ""
	for _, perm := range perms {
		coloredPerms += colorize(perm)
	}

	return coloredPerms
}

func rwx(perm os.FileMode) string {
	var b strings.Builder

	if perm&0400 != 0 {
		b.WriteString(magenta + "r" + reset)
	} else {
		b.WriteString(orange + "-" + reset)
	}
	if perm&0200 != 0 {
		b.WriteString(green + "w" + reset)
	} else {
		b.WriteString(orange + "-" + reset)
	}
	if perm&0100 != 0 {
		if perm&os.ModeSetuid != 0 {
			b.WriteString(red + "s" + reset)
		} else {
			b.WriteString(red + "x" + reset)
		}
	} else {
		if perm&os.ModeSetuid != 0 {
			b.WriteString(red + "S" + reset)
		} else {
			b.WriteString(orange + "-" + reset)
		}
	}

	return b.String()
}

func printEntry(file os.DirEntry, mode os.FileMode, directory string) {
	perms := formatPermissions(file, mode, directory)
	name := file.Name()
	fmt.Printf("%s %s\n", perms, name)
}

func printFile(file os.DirEntry, directory string) {
	info, err := file.Info()
	if err != nil {
		log.Fatal(err)
	}
	if file.IsDir() && dirOnLeft {
		icon := getDirectoryIcon(file.Name())
		fmt.Print(blue + icon + " " + file.Name() + reset)
	} else if file.IsDir() {
		icon := getDirectoryIcon(file.Name())
		fmt.Print(blue + file.Name() + " " + icon + reset)
	} else {
		fmt.Print(getFileIcon(file, info.Mode(), directory) + file.Name())
	}
	fmt.Print(" ")
}

func getDirectoryIcon(directory string) string {
	for dirType, icon := range directoryIcons {
		if filepath.Base(directory) == dirType {
			return icon
		}
	}
	return directoryIcons["default"]
}

func getFileIcon(file os.DirEntry, mode os.FileMode, directory string) string {
	if file.Type()&os.ModeSymlink != 0 {
		linkTarget, err := os.Readlink(filepath.Join(directory, file.Name()))
		if err == nil {
			symlinkTarget := filepath.Join(directory, linkTarget)
			targetInfo, err := os.Stat(symlinkTarget)
			if err == nil && targetInfo.IsDir() {
				return brightMagenta + " " + reset
			} else {
				return brightCyan + " " + reset
			}
		}
	}

	if mode.IsDir() {
		icon := getDirectoryIcon(file.Name())
		return blue + icon + " " + reset
	}

	switch file.Name() {
	case "Makefile":
		return darkBlue + " " + reset
	case "Dockerfile":
		return lightBlue + " " + reset
	case "LICENSE":
		return gray + " " + reset
	case "config":
		return lightGray + " " + reset
	case "PKGBUILD":
		return brightBlue + "󰣇 " + reset
	case ".gitconfig", ".gitignore":
		return darkOrange + " " + reset
	case ".xinitrc":
		return lightGray + " " + reset
	case ".bashrc", ".zshrc":
		return lightGray + "󱆃 " + reset
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
			return lightBlue + icon + reset
		case ".c", ".h", ".mp3", ".m4a", ".ogg", ".flac", ".php", ".lua", ".sql":
			return brightBlue + icon + reset
		case ".cs", ".mp4", ".mkv", ".webm", ".org":
			return darkMagenta + icon + reset
		case ".png", ".jpg", ".jpeg", ".JPG", ".webp":
			return darkBlue + icon + reset
		case ".gif", ".xcf", ".el":
			return magenta + icon + reset
		case ".xml":
			return lightCyan + icon + reset
		case ".htm", ".html", ".java", ".jar", ".git", ".ps", ".eps":
			return orange + icon + reset
		case ".txt", ".app":
			return white + icon + reset
		case ".zip", ".tar", ".gz", ".bz2", ".xz", ".7z", ".svg":
			return lightPurple + icon + reset
		case ".js":
			return yellow + icon + reset
		case ".json", ".tiff":
			return brightYellow + icon + reset
		case ".py":
			return darkYellow + icon + reset
		case ".yml", ".yaml", ".pdf", ".db":
			return brightRed + icon + reset
		case ".toml", ".zig":
			return darkOrange + icon + reset
		case ".deb":
			return lightRed + icon + reset
		case ".md", ".epub":
			return cyan + icon + reset
		case ".rb", ".cmake", ".pl":
			return red + icon + reset
		case ".xbps", ".vim":
			return darkGreen + icon + reset
		case ".conf", ".bat", ".rs":
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

func printPadding(name string, maxFileNameLength int) {
	padding := maxFileNameLength - len(name)
	fmt.Print(strings.Repeat(" ", padding))
}

func getFileNameAndExtension(file os.DirEntry) (string, string) {
	ext := filepath.Ext(file.Name())
	name := strings.TrimSuffix(file.Name(), ext)
	return name, ext
}

func printTree(path, prefix string, isLast bool, currentDepth, maxDepth int) {
	if maxDepth != -1 && currentDepth > maxDepth {
		return
	}

	files, err := os.ReadDir(path)
	if err != nil {
		log.Fatal(err)
	}

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
			printTree(filepath.Join(path, file.Name()), newPrefix, isLastFile, currentDepth+1, maxDepth)
		}
	}

	if showSummary && currentDepth == 0 {
		fileCount, dirCount := countFilesAndDirs(files)
		fmt.Printf(iconDirectory + " Directories: %s%d%s\n", blue, dirCount, reset)
		fmt.Printf(iconOther + " Files: %s%d%s\n", red, fileCount, reset)
	}
}
