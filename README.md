# Enhanced Directory Lister (GOLS)

This Go program lists the files and directories in the current working directory with colored and icon-based outputs for different file types.

## Installation

### Dependencies

- **Go Compiler**: [Install Go](https://go.dev/dl/)
- **Nerd Fonts**: [Download Nerd Fonts](https://www.nerdfonts.com/font-downloads)

### Clone the repository
```bash
git clone https://github.com/Tigermen0710/gols
cd gols
go build gols.go
sudo cp gols /usr/local/bin/
```
### Usage
```bash
gols
```
Or
```bash
gols [FLAG] [DIRECTORY]
```

### Flags

| flag |          description                            |       example        |
|------|-------------------------------------------------|----------------------|
|  -   | display help options or flags                   |                      |
|  -s  | show files size (-hs for human-readable format) | ![image](https://github.com/user-attachments/assets/433e18af-b869-4bfc-982a-6528341895a9) |
|  -l  | long listing (-lh for human-readable formt)     | ![image](https://github.com/user-attachments/assets/98a41e56-92b5-46ad-8780-e3c611476207) |
|  -o  | sort files by size                              |                      |
|  -t  | order all by time                             |                      |
