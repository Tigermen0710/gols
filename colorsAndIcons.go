package main

const (
	reset         = "\033[0m"
	black         = "\033[30m"
	red           = "\033[31m"
	green         = "\033[32m"
	yellow        = "\033[33m"
	blue          = "\033[34m"
	magenta       = "\033[35m"
	cyan          = "\033[36m"
	white         = "\033[37m"
	gray          = "\033[90m"
	orange        = "\033[38;5;208m"
	lightRed      = "\033[91m"
	lightGreen    = "\033[92m"
	lightYellow   = "\033[93m"
	lightBlue     = "\033[94m"
	lightMagenta  = "\033[95m"
	lightCyan     = "\033[96m"
	lightWhite    = "\033[97m"
	lightGray     = "\033[37m"
	lightOrange   = "\033[38;5;214m"
	lightPink     = "\033[38;5;218m"
	lightPurple   = "\033[38;5;183m"
	lightBrown    = "\033[38;5;180m"
	lightCyanBlue = "\033[38;5;117m"
	brightOrange  = "\033[38;5;214m"
	brightPink    = "\033[38;5;213m"
	brightCyan    = "\033[38;5;51m"
	brightPurple  = "\033[38;5;135m"
	brightYellow  = "\033[38;5;226m"
	brightGreen   = "\033[38;5;46m"
	brightBlue    = "\033[38;5;33m"
	brightRed     = "\033[38;5;196m"
	brightMagenta = "\033[38;5;198m"
	darkGray      = "\033[38;5;236m"
	darkOrange    = "\033[38;5;208m"
	darkGreen     = "\033[38;5;22m"
	darkCyan      = "\033[38;5;23m"
	darkMagenta   = "\033[38;5;90m"
	darkYellow    = "\033[38;5;172m"
	darkRed       = "\033[38;5;124m"
	darkBlue      = "\033[38;5;18m"
)

var (

	fileIcons = map[string]string{
		".go":    " ",
		".mod":   " ",
		".sh":    " ",
		".cpp":   " ",
		".hpp":   " ",
		".cxx":   " ",
		".hxx":   " ",
		".css":   " ",
		".c":     " ",
		".h":     " ",
		".cs":    "󰌛 ",
		".png":   " ",
		".jpg":   "󰈥 ",
		".JPG":   "󰈥 ",
		".jpeg":  " ",
		".webp":  " ",
		".xcf":   " ",
		".xml":   "󰗀 ",
		".htm":   " ",
		".html":  " ",
		".txt":   " ",
		".mp3":   " ",
		".m4a":   "󱀞 ",
		".ogg":   " ",
		".flac":  " ",
		".mp4":   " ",
		".mkv":   " ",
		".webm":  "󰃽 ",
		".zip":   "󰿺 ",
		".tar":   "󰛫 ",
		".gz":    " ",
		".bz2":   "󰿺 ",
		".xz":    "󰿺 ",
		".jar":   " ",
		".java":  " ",
		".js":    " ",
		".json":  " ",
		".py":    " ",
		".rs":    " ",
		".yml":   " ",
		".yaml":  " ",
		".toml":  " ",
		".deb":   " ",
		".md":    " ",
		".rb":    " ",
		".php":   " ",
		".pl":    " ",
		".svg":   "󰜡 ",
		".eps":   " ",
		".ps":    " ",
		".git":   " ",
		".zig":   " ",
		".xbps":  " ",
		".el":    " ",
		".vim":   " ",
		".lua":   " ",
		".pdf":   " ",
		".epub":  "󰂺 ",
		".conf":  " ",
		".iso":   " ",
		".exe":   " ",
		".odt":   "󰷈 ",
		".ods":   "󰱾 ",
		".odp":   "󰈧 ",
		".gif":   "󰵸 ",
		".tiff":  "󰋪 ",
		".7z":    " ",
		".bat":   " ",
		".app":   " ",
		".log":   " ",
		".sql":   " ",
		".db":    " ",
		".org":   " ",
		".ini":   "󱁻 ",
		".zst":   " ",
		".tex":   " ",
		".bash":  " ",
		".jai":   "󱢢 ",
		".r":     " ",
		".swift": "󰛥 ",
		".hs":    "󰲒 ",
		".v":     " ",
		".patch": " ",
		".diff":  " ",
		".lock":  "󰈡 ",
		".ts":    " ",
	}

	directoryIcons = map[string]string{
		"default":   "",
		"Music":     "󱍙",
		"Downloads": "󰉍",
		"Videos":    "",
		"Documents": "",
		"Pictures":  "",
		"dotfiles":  "󱗜",
		"Public":    "",
        "Movies":    "󰎁",
        "src":       "󱧼",
		"bin":       "",
		"docs":      "",
		"lib":       "",
		".github":   "",
		".git":      "",
		".config":   "",
		".ssh":      "󰣀",
		".gnupg":    "󰢬",
		".icons":    "",
		".fonts":    "",
		".cache":    "󰃨",
		".emacs.d":  "",
        ".themes":   "󰔎",
        ".npm":      "",
		".vim":      "",
	}
)

var (
	iconOther              = "\033[38;5;39m \033[0m"
	iconDirectory          = "\033[34;1m \033[0m"
)
