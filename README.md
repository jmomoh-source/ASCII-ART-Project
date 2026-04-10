# ASCII-ART-Project

Welcome to the **Unified ASCII Art** application! This tool converts standard string input into graphical ASCII formatting using predefined textual templates, combined with robust formatting capabilities.

## Features Encompassed
- **Base Application**: Map standard typography to an 8-character-tall graphical string.
- **File System (`fs`)**: Parse external banner text arrays and convert them perfectly. Support for `shadow`, `standard`, and `thinkertoy` out of the box.
- **Color Codes (`color`)**: Advanced ANSI mapping allowing users to paint specific characters or entire outputs with targeted color profiles, including native `hsl(h, s, l)` translation.
- **Terminal Justification (`justify`)**: Alignment scaling that reads your exact OS Terminal Width matching standard `left`, `right`, `center`, and `justify` bounds.
- **Data Export (`output`)**: Skips the console entirely framing the built ascii bytes into a standardized `.txt` localized disk output.

## Usage Guide

The application acts through a singular terminal entry parsing flags modularly:

```bash
# Basic string interpretation maps to standard font
go run . "hello world"

# Providing a custom font template 
go run . "hello world" shadow

# Aligning exactly perfectly right inside your Terminal Width
go run . --align=right "hello world" standard

# Justifying string distribution (Spaces letters automatically mapping bounds perfectly)
go run . --align=justify "a b c d" standard

# Coloring specific targeted substrings
go run . --color=red "world" "hello world"

# Exporting directly to text instead of logging to console
go run . --output=out.txt "hello"
```

## Directory Structure
- `main.go` - The core application entry that parses all arguments.
- `main_test.go` - Comprehensive tests running End-to-End checks on CLI behavior.
- `ascii/` - The unified package managing alignment, logic, bounding, coloring, and error handling.
- `template/` - Banner storage area holding the standard text resources.

## Documentation
For an exhaustive, technical deep-dive into every single mechanic ranging from the scaling width bounds to handling non-visible ANSI byte bugs: read the [WALKTHROUGH.md](./WALKTHROUGH.md).

> **Note**: This repository was consolidated from 5 separate legacy branches (`ascii-art`, `ascii-art-color`, `ascii-art-fs`, `ascii-art-justify`, `ascii-art-output`) into a clean cohesive structure.