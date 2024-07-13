![gols Logo](https://github.com/Tigermen0710/gols/assets/139358448/0fb308e6-df80-4cb5-840b-1221ff0b4478)
[gols](gols.png)

# gols - Enhanced `ls` with Icons, Written in Go

`gols` is a file listing tool for Unix-like systems that uses Nerd Fonts 
to display icons, providing a more visually appealing and informative output 
compared to the traditional `ls` command.

## Features

- Lists files in columns with icons based on file extensions.
- Supports long listing format (`-l`) similar to `ls -l`.
- Displays file sizes in human-readable format (`-h`) similar to `ls -lh`.
- Provides a help option (`-h`) to display usage instructions.


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

## Examples

```bash 
gols # normal list
```
```bash 
gols -l # longlist 
```
```bash 
gols -lh # longlist human-readable
```
```bash
gols -h # to show help options
```
