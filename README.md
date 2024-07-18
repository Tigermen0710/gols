![image](https://github.com/user-attachments/assets/90efeb71-b0dd-451c-8c4e-09eec752db76)

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
|  -   | display help options or flags                   | ![image](https://i.postimg.cc/R6yM3q4S/gols-help.png) |
|  -s  | show files size (-hs for human-readable format) | ![image](https://github.com/user-attachments/assets/433e18af-b869-4bfc-982a-6528341895a9) |
|  -l  | long listing (-lh for human-readable formt)     | ![image](https://github.com/user-attachments/assets/98a41e56-92b5-46ad-8780-e3c611476207) |
|  -o  | sort files by size                              | ![image](https://github.com/user-attachments/assets/80e7ce61-b606-413e-9407-f71c812a54a3) |
|  -t  | order all by time                               | ![image](https://github.com/user-attachments/assets/7037b518-c08a-464c-847e-486966bfa7ff) |
|  -m  | only show symbolik links                        | ![image](https://i.postimg.cc/hzfDPVFZ/gols-symb.png) |

