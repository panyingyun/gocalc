# go-calc

A cross-platform calculator application implemented using the Gio GUI library, demonstrating the process from design to functional debugging and showcasing the effectiveness of AI in programming.

## Features

- Basic arithmetic operations (addition, subtraction, multiplication, division)
- Clear functions (AC: clear all, CE: clear current input)
- Backspace function (⌫)
- Sign toggle (±)
- Decimal point support
- Modern Material Design interface

![Calculator](docs/ui.png)

## Requirements

- Go 1.21 or higher

## Installation and Running

1. After cloning or downloading the project, navigate to the project directory

2. Install dependencies:
```bash
go mod tidy
```

3. Run the application:
```bash
go run main.go
```

## Usage

- Click number buttons to input numbers
- Click operator buttons (+, -, ×, ÷) to select operations
- Click equals (=) to perform calculations
- Click AC to clear all data and operations
- Click CE to clear current input
- Click ⌫ to delete the last digit
- Click ± to toggle sign

## Environment Setup

### Ubuntu 

```bash
apt install gcc pkg-config libwayland-dev libx11-dev libx11-xcb-dev libxkbcommon-x11-dev libgles2-mesa-dev libegl1-mesa-dev libffi-dev libxcursor-dev libvulkan-dev
```

### Windows

```bash
nothing need to install
```

### Mac

```bash
Xcode is required for Apple platforms.
```

## Building Executables

### Windows
```bash
go build -ldflags="-H windowsgui" -o calculator.exe main.go
```

### Linux/macOS
```bash
go build -o calculator main.go
```

### Android or other
```
TODO
```
