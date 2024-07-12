![gols Logo](https://github.com/Tigermen0710/gols/assets/139358448/0fb308e6-df80-4cb5-840b-1221ff0b4478)

# gols - Enhanced `ls` with Icons, Written in Go

`gols` is a file listing tool for Unix-like systems that uses Nerd Fonts 
to display icons, providing a more visually appealing and informative output 
compared to the traditional `ls` command.

## Installation

### Dependencies

Before installing `gols`, ensure you have the following dependencies installed:

- **Go**: [Install Go](https://go.dev/doc/install)
- **Nerd Fonts**: [Download Nerd Fonts](https://www.nerdfonts.com/font-downloads)

### Installation Steps

To install `gols`, follow these steps:

### Option 1
```bash
git clone https://github.com/Tigermen0710/gols
cd gols
go build gols.go
sudo mv gols /bin/gols
```
To use gols, simply run:
```bash
gols
```
### Option 2
```bash
go install github.com/Tigermen0710/gols@latest
```
NOTE: Your Go path must be set. Here's how to set it up: [Setup go path](https://go.dev/wiki/SettingGOPATH)
