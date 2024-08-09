# Enhanced Directory Lister (GOLS)

![image](https://i.postimg.cc/htv8YKBp/golsshot.jpg)

This program lists the files and directories in the current working directory with colored and icon-based outputs for different file types.

## Features

- List files and directories plus symlinks.
- Supports showing hidden files or directories.
- Colored icons based on file types.
- List directories.
- Order (sort) files by time or size.
- Show a tree of current or any directory path.
- Size of files.
- Options to show directories or files only.
- Summary of files and directories.
- List file based on extention ```gols /path/to/dir/ .go``` to list all go files.
- Exlude files using there extention ```gols -x go,txt ...```.
- Use the extention to list files ```gols -e go``` to list golang files.


## Installation

### Dependencies

- **Go Compiler**: [Install Go](https://go.dev/dl/)
- **Nerd Fonts**: [Download Nerd Fonts](https://www.nerdfonts.com/font-downloads)

### Clone the repository
```bash
git clone https://github.com/Tigermen0710/gols
cd gols/
go build
sudo cp gols /usr/local/bin/
sudo cp gols.1 /usr/local/share/man/man1/ # To copy the man page.
gols
```
#### Or using the Makefile.

```bash
git clone https://github.com/Tigermen0710/gols
cd gols/
make
sudo make install
gols
```

#### Note: gols is in the [AUR](https://aur.archlinux.org/packages/gols) and a [template](https://github.com/elbachir-one/void-templates) for Void Linux (xbps-src).

##### Arch Linux:
```bash
yay -S gols
```

##### Void Linux:

Assuming you have void-packages.
```bash
git clone https://github.com/elbachir-one/void-templates
cp void-templates/gols/ void-packages/srcpkgs/
./xbps-src pkg gols
sudo xbps-install -R hostdir/binpkgs gols
```

## Usage
```bash
gols
```
Or
```bash
gols [FLAG] [DIRECTORY] [FILES]
```

### Flags

| flag | description                                                  | example                                                                                         |
|------|--------------------------------------------------------------|-------------------------------------------------------------------------------------------------|
| -?   | display help options or flags                                | ![image](https://i.postimg.cc/htsDBSD7/image.png)                                               |
| -a   | show hidden files or directories                             | ![image](https://i.postimg.cc/zGsDxgmV/a-flag.png)                                              |
| -A   | show only hidden directories and files                       | ![image](https://i.postimg.cc/SQYzhZCc/A.png)                                                   |
| -c   | show all files in one column                                 | ![image](https://github.com/user-attachments/assets/07ec7ab1-3740-487c-8602-03963b3c556d)       |
| -D   | list only directories                                        | ![image](https://i.postimg.cc/52M98M9g/D.png)                                                   |
| -e   | list files based on there extention                          | ![image](https://i.postimg.cc/fLxxT1NJ/e.png)                                                   |
| -f   | show a summary of file and directories                       | ![image](https://i.postimg.cc/gcL2ZFDf/ff.png)                                                  |
| -F   | list files only                                              | ![image](https://i.postimg.cc/Z5FbcDCS/F.png)                                                   |
| -i   | show directory icon on left                                  | ![image](https://i.postimg.cc/Z0tKKdX7/i.png)                                                   |
| -l   | long listing                                                 | ![image](https://github.com/user-attachments/assets/98a41e56-92b5-46ad-8780-e3c611476207)       |
| -m   | only show symbolik links                                     | ![image](https://i.postimg.cc/N2f5FZ1s/symlink.png)                                             |
| -o   | sort files by size                                           | ![image](https://github.com/user-attachments/assets/80e7ce61-b606-413e-9407-f71c812a54a3)       |
| -r   | tree like listing, and d number to do the depth (gols -rd 1) | ![image](https://i.postimg.cc/rsdQLxW4/tree.png) ![image](https://i.postimg.cc/PJ5NmZC4/rd.png) |
| -s   | show files size                                              | ![image](https://github.com/user-attachments/assets/433e18af-b869-4bfc-982a-6528341895a9)       |
| -t   | order all by time                                            | ![image](https://github.com/user-attachments/assets/7037b518-c08a-464c-847e-486966bfa7ff)       |
| -v   | version number                                               |                                                                                                 |
| -x   | exclude files from the listing using there extention         | ![image](https://i.postimg.cc/90Cy41m1/x.png)                                                   |
