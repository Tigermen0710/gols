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
    "unsafe"
)

var (
    longListing         bool
    humanReadable       bool
    fileSize            bool
    orderBySize         bool
    orderByTime         bool
    showOnlySymlinks    bool
    showHidden          bool
    recursiveListing    bool
    dirOnLeft           bool
    showSummary         bool
    showVersion         bool
    maxDepth            int = -1
    listDirsOnly        bool
    listFilesOnly       bool
    listHiddenOnly      bool
    oneColumn           bool
    fileExtensions      []string
    extFlag             string
    excludeExtensions   bool
    excludedExts        []string
    onlyPermissions     bool
    showOwner           bool
    getTime             bool
)

type winsize struct {
    Row    uint16
    Col    uint16
    Xpixel uint16
    Ypixel uint16
}

type fakeDirEntry struct {
    info os.FileInfo
}

func main() {
    args := os.Args[1:]
    nonFlagArgs, hasFlags, hasSpecificFlags := parseFlags(args)

    if showVersion {
        fmt.Println(version)
        return
    }

    var directory string

    if len(nonFlagArgs) > 0 {
        directory = nonFlagArgs[0]
    }

    if len(nonFlagArgs) > 1 {
        fileExtensions = strings.Split(nonFlagArgs[1], ",")
    } else if extFlag != "" {
        fileExtensions = strings.Split(extFlag, ",")
    }

    if directory == "" {
        directory = "."
    }

    var files []os.DirEntry
    var err error
    var info os.FileInfo

    info, err = os.Stat(directory)
    if err != nil {
        log.Fatalf("Error: %v", err)
    }

    if info.IsDir() {
        files, err = os.ReadDir(directory)
        if err != nil {
            log.Fatal(err)
        }
    } else {
        files = []os.DirEntry{&fakeDirEntry{info}}
        directory = filepath.Dir(directory)
    }

    if len(fileExtensions) > 0 {
        extSet := make(map[string]struct{})
        for _, ext := range fileExtensions {
            ext = strings.TrimPrefix(ext, ".")
            extSet[ext] = struct{}{}
        }
        var filteredFiles []os.DirEntry
        for _, file := range files {
            ext := strings.TrimPrefix(filepath.Ext(file.Name()), ".")
            if _, found := extSet[ext]; found {
                filteredFiles = append(filteredFiles, file)
            }
        }
        files = filteredFiles
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

    if listDirsOnly {
        files = filterDirectories(files)
    } else if listFilesOnly {
        files = filterNonDirectories(files)
    } else if listHiddenOnly {
        files = filterHiddenOnly(files)
    }

    if len(excludedExts) > 0 {
        files = filterExcludedExtensions(files, excludedExts)
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

    if onlyPermissions {
        printPermissionsWithIcons(files, directory)
    } else if showOwner {
        printOwner(files, directory)
    } else if getTime {
        printTime(files, directory)
    } else if recursiveListing {
        printTree(directory, "", true, 0, maxDepth)
    } else if longListing {
        printLongListing(files, directory, humanReadable)
    } else if fileSize {
        getFileSize(files, directory, humanReadable, dirOnLeft)
    } else {
        printFilesInColumns(files, directory, dirOnLeft, showSummary)
    }

    if (hasSpecificFlags && !longListing) || !hasFlags {
        fmt.Println()
    }
}

func (f *fakeDirEntry) Name() string               { return f.info.Name() }
func (f *fakeDirEntry) IsDir() bool                { return f.info.IsDir() }
func (f *fakeDirEntry) Type() os.FileMode          { return f.info.Mode().Type() }
func (f *fakeDirEntry) Info() (os.FileInfo, error) { return f.info, nil }

func filterByExtension(files []os.DirEntry, extension string) []os.DirEntry {
    var filtered []os.DirEntry
    for _, file := range files {
        if strings.TrimPrefix(filepath.Ext(file.Name()), ".") == extension {
            filtered = append(filtered, file)
        }
    }
    return filtered
}

func filterNonDirectories(files []os.DirEntry) []os.DirEntry {
    var nonDirs []os.DirEntry
    for _, file := range files {
        if !file.IsDir() {
            nonDirs = append(nonDirs, file)
        }
    }
    return nonDirs
}

func filterDirectories(entries []os.DirEntry) []os.DirEntry {
    var result []os.DirEntry
    for _, entry := range entries {
        if entry.IsDir() {
            result = append(result, entry)
        }
    }
    return result
}

func filterFiles(entries []os.DirEntry) []os.DirEntry {
    var result []os.DirEntry
    for _, entry := range entries {
        if !entry.IsDir() {
            result = append(result, entry)
        }
    }
    return result
}

func filterExcludedExtensions(files []os.DirEntry, excludedExts []string) []os.DirEntry {
    var filteredFiles []os.DirEntry
    for _, file := range files {
        ext := strings.TrimPrefix(filepath.Ext(file.Name()), ".")
        exclude := false
        for _, excludedExt := range excludedExts {
            if ext == excludedExt {
                exclude = true
                break
            }
        }
        if !exclude {
            filteredFiles = append(filteredFiles, file)
        }
    }
    return filteredFiles
}

func filterHiddenOnly(files []os.DirEntry) []os.DirEntry {
    var hiddenFiles []os.DirEntry
    for _, file := range files {
        if strings.HasPrefix(file.Name(), ".") {
            hiddenFiles = append(hiddenFiles, file)
        }
    }
    return hiddenFiles
}

func filterByExtensions(files []os.DirEntry, extensions []string) []os.DirEntry {
    var filtered []os.DirEntry
    extMap := make(map[string]struct{})
    for _, ext := range extensions {
        ext = strings.TrimPrefix(ext, ".")
        extMap[ext] = struct{}{}
    }
    for _, file := range files {
        ext := strings.TrimPrefix(filepath.Ext(file.Name()), ".")
        if _, found := extMap[ext]; found {
            filtered = append(filtered, file)
        }
    }
    return filtered
}

func getTerminalWidth() (int, error) {
    ws := &winsize{}
    _, _, err := syscall.Syscall(syscall.SYS_IOCTL, uintptr(syscall.Stdin), uintptr(syscall.TIOCGWINSZ), uintptr(unsafe.Pointer(ws)))
    if err != 0 {
        return 0, fmt.Errorf("failed to get terminal size")
    }
    return int(ws.Col), nil
}

func printPadding(name string, maxFileNameLength int) {
    padding := maxFileNameLength - len(name) + 1
    for i := 0; i < padding; i++ {
        fmt.Print(" ")
    }
}

func truncateName(name string, maxLength int) string {
    if len(name) > maxLength {
        return name[:maxLength-1] + "…"
    }
    return name
}

func printFile(file os.DirEntry, directory string, maxLength int, dirOnLeft bool) {
    info, err := file.Info()
    if err != nil {
        log.Fatal(err)
    }

    truncatedName := truncateName(file.Name(), maxLength)

    if file.IsDir() && dirOnLeft {
        icon := getDirectoryIcon(file.Name())
        fmt.Print(blue + icon + " " + truncatedName + reset)
    } else if file.IsDir() {
        icon := getDirectoryIcon(file.Name())
        fmt.Print(blue + truncatedName + " " + icon + reset)
    } else {
        fmt.Print(getFileIcon(file, info.Mode(), directory) + truncatedName)
    }
}

func truncateString(s string, maxLength int) string {
    if len(s) > maxLength {
        return s[:maxLength-3] + "..."
    }
    return s
}

func getMaxNameLength(files []os.DirEntry) int {
    maxLen := 0
    for _, file := range files {
        if len(file.Name()) > maxLen {
            maxLen = len(file.Name())
        }
    }
    return maxLen
}

func printFilesInColumns(files []os.DirEntry, directory string, dirOnLeft bool, showSummary bool) {
    if oneColumn {
        for _, file := range files {
            printFile(file, directory, getMaxNameLength(files), dirOnLeft)
            fmt.Println()
        }
    } else {
        terminalWidth, err := getTerminalWidth()
        if err != nil {
            fmt.Println("Error getting terminal width:", err)
            return
        }

        maxFileNameLength := getMaxNameLength(files)
        columnWidth := maxFileNameLength + 1

        maxFilesInLine := terminalWidth / columnWidth

        filesInLine := 0

        for _, file := range files {
            printFile(file, directory, maxFileNameLength, dirOnLeft)

            filesInLine++
            if filesInLine >= maxFilesInLine {
                fmt.Println()
                filesInLine = 0
            } else {
                printPadding(truncateName(file.Name(), maxFileNameLength), maxFileNameLength)
            }
        }
    }

    if showSummary {
        fmt.Println()
        fmt.Println()
        printSummary(files, directory)
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
            if len(arg) > 2 && arg[1] == '-' {
                switch arg {
                case "--version":
                    showVersion = true
                case "--help":
                    showHelp()
                    os.Exit(0)
                default:
                    fmt.Println("Unknown long flag:", arg)
                    showHelp()
                    os.Exit(1)
                }
            } else {

                for j := 1; j < len(arg); j++ {
                    switch arg[j] {
                    case 'l':
                        longListing = true
                    case 'c':
                        oneColumn = true
                    case 'h':
                        humanReadable = true
                        hasSpecificFlags = true
                    case 's':
                        fileSize = true
                    case 'o':
                        orderBySize = true
                        hasSpecificFlags = true
                    case 'p':
                        onlyPermissions = true
                    case 'O':
                        showOwner = true
                    case 't':
                        orderByTime = true
                        hasSpecificFlags = true
                    case 'T':
                        getTime = true
                    case 'm':
                        showOnlySymlinks = true
                        hasSpecificFlags = true
                    case 'a':
                        showHidden = true
                        hasSpecificFlags = true
                    case 'A':
                        listHiddenOnly = true
                        showHidden = true
                        hasSpecificFlags = true
                    case 'r':
                        recursiveListing = true
                    case 'i':
                        dirOnLeft = true
                        hasSpecificFlags = true
                        hasFlags = true
                    case 'f':
                        showSummary = true
                        hasSpecificFlags = true
                    case 'v':
                        showVersion = true
                    case 'D':
                        listDirsOnly = true
                        hasSpecificFlags = true
                    case 'F':
                        listFilesOnly = true
                        hasSpecificFlags = true
                    case 'x':
                        if j+1 < len(arg) && (arg[j+1] < '0' || arg[j+1] > '9') {
                            excludedExts = strings.Split(arg[j+1:], ",")
                            hasSpecificFlags = true
                            break
                        } else if i+1 < len(args) && args[i+1][0] != '-' {
                            excludedExts = strings.Split(args[i+1], ",")
                            hasSpecificFlags = true
                            i++
                            break
                        } else {
                            fmt.Println("Missing value for -x")
                            os.Exit(1)
                        }
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
                    case 'e':
                        if i+1 < len(args) {
                            extFlag = args[i+1]
                            i++
                            hasSpecificFlags = true
                        } else {
                            fmt.Println("Missing value for -e")
                            os.Exit(1)
                        }
                    default:
                        showHelp()
                        os.Exit(1)
                    }
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
    fmt.Println("	-? --help       Help")
    fmt.Println()
    fmt.Println("	-a              Show Hidden files")
    fmt.Println("	-A              Show only hidden files and directories")
    fmt.Println("	-e              Filter files based on extensions")
    fmt.Println("	-f              Show summary of directories and files")
    fmt.Println("	-F              List files only")
    fmt.Println("	-c              Don't use spacing, print all files in one column")
    fmt.Println("	-D              Only directories are showing")
    fmt.Println("	-h              Human-readable file sizes")
    fmt.Println("	-i              Show directory icon on left")
    fmt.Println("	-l              Long listing format")
    fmt.Println("	-m              Only symbolic links are showing")
    fmt.Println("	-o              Sort by size")
    fmt.Println("	-r d n          Tree like listing, set the depth of the directory tree (n is an integer)")
    fmt.Println("	-s              Print files size")
    fmt.Println("	-t              Order by time")
    fmt.Println("	-v --version    Show version")
    fmt.Println("	-x              Exclude specific extensions")
    fmt.Println()
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
        fmt.Println()
        printSummary(files, directory)
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

func printPermissionsWithIcons(files []os.DirEntry, directory string) {
    for _, file := range files {
        info, err := file.Info()
        if err != nil {
            log.Fatal(err)
        }

        permissions := formatPermissions(file, info.Mode(), directory)
        permissions = green + permissions + reset

        iconAndName := getFileIcon(file, info.Mode(), directory) + " " + file.Name()

        fmt.Printf("%s %s\n", permissions, iconAndName)
    }
}

func printOwner(files []os.DirEntry, directory string) {
    for _, file := range files {
        info, err := file.Info()
        if err != nil {
            log.Fatal(err)
        }

        owner, err := user.LookupId(fmt.Sprintf("%d", info.Sys().(*syscall.Stat_t).Uid))
        if err != nil {
            log.Fatal(err)
        }

        ownerStr := cyan + owner.Username + reset
        icon := getFileIcon(file, info.Mode(), directory)
        fileName := file.Name()

        fmt.Printf("%s %s %s\n", ownerStr, icon, fileName)
    }
}

func printTime(files []os.DirEntry, directory string) {
    for _, file := range files {
        info, err := file.Info()
        if err != nil {
            log.Fatal(err)
        }

        modTime := info.ModTime()
        timeStr := modTime.Format("15:04:05")
        dateStr := modTime.Format("2006-01-02")
        icon := getFileIcon(file, info.Mode(), directory)
        fileName := file.Name()

        fmt.Printf("%s %s %s %s\n", dateStr, timeStr, icon, fileName)
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
        fmt.Println()
        printSummary(filteredFiles, directory)
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

func getDirectoryIcon(directory string) string {
    for dirType, icon := range directoryIcons {
        if filepath.Base(directory) == dirType {
            return icon
        }
    }
    return directoryIcons["default"]
}

func getSpecialFileIcon(fileName string) (string, bool) {
    icon, found := specialFileIcons[fileName]
    return icon, found
}

func getFileIcon(file os.DirEntry, mode os.FileMode, directory string) string {
    if file.Type()&os.ModeSymlink != 0 {
        linkTarget, err := os.Readlink(filepath.Join(directory, file.Name()))
        if err == nil {
            symlinkTarget := filepath.Join(directory, linkTarget)
            targetInfo, err := os.Stat(symlinkTarget)
            if err == nil && targetInfo.IsDir() {
                return iconSymlinkDir
            } else {
                return iconSymlinkFile
            }
        }
    }

    if mode.IsDir() {
        icon := getDirectoryIcon(file.Name())
        return blue + icon + " " + reset
    }

    if icon, found := getSpecialFileIcon(file.Name()); found {
        return icon
    }

    ext := filepath.Ext(file.Name())
    icon, exists := fileIcons[ext]
    if exists {
        switch ext {
        case ".sh", ".ps1":
            if mode&os.ModePerm&0111 != 0 {
                return brightGreen + icon + reset
            } else {
                return white + icon + reset
            }
        case ".cpp", ".hpp", ".cxx", ".hxx", ".dart", ".gd", ".v":
            return blue + icon + reset
        case ".css", ".ml", ".rst", ".nix":
            return lightBlue + icon + reset
        case ".c", ".h", ".mp3", ".m4a", ".ogg", ".flac", ".php", ".lua", ".sql", ".m":
            return brightBlue + icon + reset
        case ".png", ".jpg", ".jpeg", ".JPG", ".webp", ".R", ".ts", ".bmp":
            return darkBlue + icon + reset
        case ".md", ".epub", ".obj", ".go":
            return cyan + icon + reset
        case ".xml":
            return lightCyan + icon + reset
        case ".exe", ".desktop", ".mk":
            return brightCyan + icon + reset
        case ".gif", ".xcf", ".el", ".lisp":
            return magenta + icon + reset
        case ".cs", ".mp4", ".mkv", ".webm", ".org", ".ejs":
            return darkMagenta + icon + reset
        case ".js", ".lock":
            return yellow + icon + reset
        case ".json", ".tiff", ".nim":
            return brightYellow + icon + reset
        case ".patch", ".diff", ".py":
            return darkYellow + icon + reset
        case ".yml", ".yaml", ".pdf", ".db":
            return brightRed + icon + reset
        case ".deb":
            return lightRed + icon + reset
        case ".rb", ".cmake", ".pl", ".scala", ".erl", ".build":
            return red + icon + reset
        case ".htm", ".html", ".java", ".jar", ".git", ".ps", ".eps", ".swift":
            return orange + icon + reset
        case ".toml", ".zig":
            return darkOrange + icon + reset
        case ".tmux.conf":
            return green + icon + reset
        case ".xbps", ".vim", ".jai":
            return darkGreen + icon + reset
        case ".iso", ".asm", ".f90", ".groovy", ".ini", ".cfg":
            return gray + icon + reset
        case ".conf", ".bat", ".rs":
            return darkGray + icon + reset
        case ".fish", ".o", ".m4":
            return lightGray + icon + reset
        case ".1", ".hs":
            return lightBrown + icon + reset
        case ".txt", ".app":
            return white + icon + reset
        case ".zip", ".tar", ".gz", ".bz2", ".xz", ".7z", ".svg", ".kt", ".ex", ".zst":
            return lightPurple + icon + reset
        default:
            return icon
        }
    }

    if mode&os.ModePerm&0111 != 0 {
        return green + " " + reset
    }

    return " " + reset
}

func getFileNameAndExtension(file os.DirEntry) (string, string) {
    ext := filepath.Ext(file.Name())
    name := strings.TrimSuffix(file.Name(), ext)
    return name, ext
}

func printSummary(files []os.DirEntry, directory string) {
    fileCount, dirCount, symlinkFileCount, symlinkDirCount := 0, 0, 0, 0

    for _, file := range files {
        if file.Type()&os.ModeSymlink != 0 {
            linkTarget, err := os.Readlink(filepath.Join(directory, file.Name()))
            if err == nil {
                targetInfo, err := os.Stat(filepath.Join(directory, linkTarget))
                if err == nil && targetInfo.IsDir() {
                    symlinkDirCount++
                } else {
                    symlinkFileCount++
                }
            }
        } else if file.IsDir() {
            dirCount++
        } else {
            fileCount++
        }
    }

    fmt.Printf(iconDirectory + " Directories: %s%d%s\n", blue, dirCount, reset)
    fmt.Printf(iconOther + " Files: %s%d%s\n", red, fileCount, reset)

    if symlinkDirCount > 0 {
        fmt.Printf(iconSymlinkDir + " Symlinked Directories: %s%d%s\n", magenta, symlinkDirCount, reset)
    }
    if symlinkFileCount > 0 {
        fmt.Printf(iconSymlinkFile + " Symlinked Files: %s%d%s\n", cyan, symlinkFileCount, reset)
    }

    total := dirCount + fileCount + symlinkDirCount + symlinkFileCount
    fmt.Printf(iconTotal + ":%s%d%s\n", brightGreen, total, reset)
}

func printTree(path, prefix string, isLast bool, currentDepth, maxDepth int) (totalFiles, totalDirs int) {
    if maxDepth != -1 && currentDepth > maxDepth {
        return 0, 0
    }

    files, err := os.ReadDir(path)
    if err != nil {
        if os.IsPermission(err) {
            fmt.Printf("%sError: Permission denied for %s%s\n", red, path, reset)
        } else {
            fmt.Printf("%sError reading directory %s: %v%s\n", red, path, err, reset)
        }
        return 0, 0
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

        maxFileNameLength := getMaxNameLength(filteredFiles)
        printFile(file, path, maxFileNameLength, true)
        fmt.Println()

        if file.Type()&os.ModeSymlink != 0 {
            linkTarget, err := os.Readlink(filepath.Join(path, file.Name()))
            if err == nil {
                fmt.Printf("%s%s ==> %s%s\n", prefix, cyan, linkTarget, reset)
            } else {
                fmt.Printf("%s%s %s%s\n", prefix, red, "==> error", reset)
            }
        }

        if file.IsDir() {
            newPrefix := prefix
            if isLastFile {
                newPrefix += "    "
            } else {
                newPrefix += "│   "
            }
            subPath := filepath.Join(path, file.Name())
            subFiles, subDirs := printTree(subPath, newPrefix, isLastFile, currentDepth+1, maxDepth)
            totalFiles += subFiles
            totalDirs += subDirs + 1
        } else {
            totalFiles++
        }
    }

    if currentDepth == 0 && showSummary {
        fmt.Println()
        printSummary(filteredFiles, path)
    }

    return totalFiles, totalDirs
}
