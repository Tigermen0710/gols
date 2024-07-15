# Enhanced Directory Lister (GOLS)

This Go program lists the files and directories in the current working directory with colored and icon-based outputs for different file types.

## Changes Made

### Code Refactoring
- **Function Separation**: Extracted the logic for getting file colors and icons into separate functions (`getFileColor` and `getFileIcon`), improving the readability of the `main` function.
- **Constants for Colors**: Defined the color codes as constants for better readability and maintainability.
- **String Repeat**: Used `strings.Repeat` for padding spaces instead of a manual loop.
- **Switch Statements**: Used switch statements instead of multiple `if-else` statements for better clarity and performance.

### Security Enhancements
- **Path Sanitization**: Added a `sanitizePath` function to ensure the directory path is clean and absolute, and to check for path traversal attempts.
- **Safe Path Handling**: Used the `path/filepath` package to handle paths safely.
- **Error Handling**: Improved error handling in the `sanitizePath` function.

## Problems Addressed
- **Unvalidated Path Input**: Previously, the program used unvalidated input for directory paths, which could lead to path injection vulnerabilities. This has been mitigated by sanitizing the path input.
- **Code Readability**: The original code had multiple nested if-else statements, making it difficult to read and maintain. Refactoring into separate functions has improved readability and maintainability.
- **Path Traversal Risk**: The original code did not check for path traversal attempts, which could have allowed unintended file or directory access. The added sanitization function addresses this issue.

## Installation

### Dependencies
- **GCC Go Compiler**: [Install GCC Go](https://go.dev/doc/install/gccgo)
- **Nerd Fonts**: [Download Nerd Fonts](https://www.nerdfonts.com/font-downloads)

### Bash
```bash
git clone https://github.com/Tigermen0710/gols
cd gols
gccgo gols.go -o ls
sudo mv gols.sh /bin/gols
```

To run it, use:
```bash
gols
```

### Improvements Made:
1. **Input Validation**: Ensured the provided directory exists and is valid.
2. **Safety**: Quoted variables to handle spaces and special characters in paths.
3. **Error Handling**: Added meaningful error messages if the directory change fails.
4. **Portability**: Maintained the use of `exec` to replace the shell with the specified command for better performance.

---

