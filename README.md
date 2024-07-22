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

#### Note: gols is in the [AUR](https://aur.archlinux.org/packages/gols).

```bash
yay -Sy gols
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

| flag | description                                     | example                                                                                   |
|------|-------------------------------------------------|-------------------------------------------------------------------------------------------|
| -    | display help options or flags                   | ![image](https://i.postimg.cc/Ff3fByr4/flags.png)                                         |
| -a   | show hidden files or directories                | ![image](https://i.postimg.cc/zGsDxgmV/a-flag.png)                                        |
| -l   | long listing (-lh for human-readable formt)     | ![image](https://github.com/user-attachments/assets/98a41e56-92b5-46ad-8780-e3c611476207) |
| -m   | only show symbolik links                        | ![image](https://i.postimg.cc/N2f5FZ1s/symlink.png)                                       |
| -o   | sort files by size                              | ![image](https://github.com/user-attachments/assets/80e7ce61-b606-413e-9407-f71c812a54a3) |
| -s   | show files size (-hs for human-readable format) | ![image](https://github.com/user-attachments/assets/433e18af-b869-4bfc-982a-6528341895a9) |
| -t   | order all by time                               | ![image](https://github.com/user-attachments/assets/7037b518-c08a-464c-847e-486966bfa7ff) |
