# Enhanced Directory Lister (GOLS)

This Go program lists the files and directories in the current working directory with colored and icon-based outputs for different file types.

## Installation

### Dependencies
- **GCC Go Compiler**: [Install GCC Go](https://go.dev/doc/install/gccgo)
- **Nerd Fonts**: [Download Nerd Fonts](https://www.nerdfonts.com/font-downloads)

### Bash
```bash
git clone https://github.com/Tigermen0710/gols
cd gols
gccgo gols.go -o gols
chmod +x gols
sudo mv gols /bin/
```

To run it, use:
```bash
gols
```
## flags:

-s  - show files size (-hs for human-readable format)

![image](https://github.com/user-attachments/assets/433e18af-b869-4bfc-982a-6528341895a9)

-l  - long listing (-lh for human-readable format)

![image](https://github.com/user-attachments/assets/98a41e56-92b5-46ad-8780-e3c611476207)
