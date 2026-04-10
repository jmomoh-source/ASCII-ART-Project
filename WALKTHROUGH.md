# Deep-Dive Developer Walkthrough: Complete ASCII-Art Project

Welcome to the **Complete ASCII-Art Walkthrough**. This document serves as an exhaustive, technical deep-dive into every single mechanic, algorithm, and file powering the unified ASCII-Art project.

Whether you are trying to understand the mathematical alignment formulas, how ANSI string manipulation bypasses length checks, or how variables flow step-by-step from input to terminal, this guide covers **every single detail** for beginners and advanced developers alike.

---

## 1. Project Amalgamation Context
This project merges five previously standalone capabilities into a single powerful application:
1. **Base Framework**: Translates typed letters into 8x8 character art.
2. **File Systems (`-fs`)**: Interacting dynamically with various typography styles (`standard.txt`, `shadow.txt`).
3. **Color Processing (`-color`)**: Sub-scanning words to paint exact letters with ANSI codes.
4. **Terminal Justification (`-justify`)**: Automatically resizing spacing outputs matching your terminal `OS` bounds.
5. **Data Output (`-output`)**: Bypassing the terminal and writing formatted art entirely to `.txt` files.

---

## 2. Directory Structure Breakdown

The codebase lives completely inside the `ascii-art/` folder. It uses a clean separation of concerns:

```text
ascii-art/
├── main.go               // The Primary CLI router dealing with inputs and OS system exits.
├── template/             // The raw graphical datasets.
│   ├── standard.txt
│   ├── shadow.txt
│   └── thinkertoy.txt
└── ascii/                // The internal domain library where logic lives.
    ├── types.go          // Definitions, constants, and data structure rules.
    ├── width.go          // OS-level checks catching terminal dimensions.
    ├── color.go          // Matches strings and translates them to ANSI colors.
    ├── align.go          // Mathematics for Left, Right, and Center alignment formatting.
    ├── font.go           // Disk scanner that cleans up standard .txt typography.
    ├── generate.go       // The core Engine orchestrator bridging all systems above.
    ├── ascii.go          // Simple entry function mapped specifically backwards for legacy Testing.
    └── ascii_test.go     // Automated Test cases ensuring nothing breaks over time.
```

---

## 3. Step-By-Step Execution Pathway

What happens the millisecond you press "Enter" on a command like:
`go run . --color=red --align=justify "hello world" shadow`

### Phase A: Decoding User Intent (`main.go`)
Before we draw anything, we have to isolate your arguments safely avoiding index panics.
1. The `main()` function fetches `os.Args[1:]` (the arguments ignoring the program name itself).
2. It spins up an array loop targeting the inputs. We utilize robust `strings.HasPrefix(arg, "--[tag]")` checks to find your options. 
3. **Sorting**: 
   - `--color=red` gets parsed. `red` is grabbed and placed into a slice tracking pending colors.
   - `--align=justify` tells our padding system what mode to fall into.
   - Once all `--` flags are isolated, whatever remains falls predictably: The first remaining string (`"hello world"`) becomes the **Target Text**, and the string directly behind it (`"shadow"`) becomes the **Banner Type**.
4. **Why this matters**: Because of this sophisticated logic, a user can write `go run . "hello" --color=green` OR `go run . --color=green "hello"`, and the program behaves perfectly regardless of order.

### Phase B: Loading The Aesthetics (`ascii/font.go`)
Once `main.go` identifies the banner (`shadow`), it calls `LoadBanner()`.
1. **Safety Append**: It checks if the user wrote `shadow.txt`. If they didn't, it automatically appends `.txt`.
2. **File Reading**: `os.ReadFile` runs pointing to the nested `/template` folder.
3. **Endings Standardization**: This is a critical bug-stopper! Different operating systems break lines differently. Windows uses `\r\n` and Unix uses `\n`. We run a `ReplaceAll` to force all `\r\n` instances safely into standard `\n`.
4. Finally, the raw text is chopped up into a massive Go Slice using `strings.Split()`.

### Phase C: Scanning the Viewport (`ascii/width.go`)
If the user requested formatting (Center, Right, Justify), we need to know the exact physical constraints of your terminal monitor.
1. `GetTerminalWidth()` invokes an OS-Level execution mapping to Unix `stty size`.
2. Standard out returns something like `24 80` (Rows, Columns).
3. The method uses `strings.Fields()` to extract the columns, converting `"80"` into a tangible standard integer via `strconv.Atoi`.
4. **Fallback Handling**: If `stty` fails (e.g., if you pipe data through SSH, or you're on base Windows), the code safely suppresses the crash and defaults gracefully to an 80-character width column space limit.

### Phase D: Color Processing algorithms (`ascii/color.go`)
Colors in terminals are created using invisible **ANSI Escape Codes** like so: `\033[31m` (Turns text red), and `\033[0m` (Disables the color back to white).
When you provide a substring parameter, our program must identify where in the text your target sits.
1. Inside `GetCharColors()`, we generate an empty generic array whose exact length identically maps standard byte-indexes of your current target word.
2. In a loop, it slides over your target word comparing segments `word[i:i+len(substring)] == substring`. 
3. If it returns true, it fills the exact character bytes of that index slice with the target color. 

### Phase E: Generating the Artwork (`ascii/generate.go`)
The true heartbeat of the application lies in `GenerateAsciiArt()`.

#### Step 1: Mapping Newlines
We pass in your target text, but we need to split it if you supplied `\n` line breaks.
- **Edgecase**: Shell parsers like bash pass the literal string `\n` (two keys: a backslash & 'n'), not an actual newline break. 
- The loop calls `strings.Split(input, "\\n")` ensuring it actually slices your strings when it encounters those keys. If a slice maps as explicitly empty, it prints an exact newline to your terminal cleanly handling `\n\n` double spacing.

#### Step 2: Extracting Art Coordinates
For a targeted word like `hello`, it executes an iteration across every character one by one (`h`, then `e`, `l`...):
- Every graphical ascii letter takes up uniformly **8 Vertical Rows**. 
- It processes the mathematical formula `index = int(char-' ')*9 + 1` to find exactly where the `h` sequence resides within the massive `LoadBanner` text slice mapping algorithm.

It creates 8 empty line buckets using `strings.Builder`. 
For row 0: it grabs row 0 of `h`, appends it to bucket 0. Then attaches row 0 of `e` directly beside it, building horizontally line-by-line avoiding messy newline line breaks until the entire letter loop finalizes!

#### Step 3: ANSI Color Injection
During the bucket building, before appending an `h`, it checks the map built earlier by `color.go`. If `h` was targeted in your argument, it wraps the bucket append identically like:
`ANSI RED CODE` + `Graphic h` + `ANSI RESET CODE`.

---

## 4. The Mathematics Of Alignments (`ascii/align.go`)
Alignment relies on taking full terminal width and padding strings with spacing `("   ")`. This presents a massive issue: **ANSI Code Length Interference**. Let's break this detail down entirely:

### Problem: Invisible Bytes
If the word "test" is 20 physical block characters long, measuring `len("test")` normally returns 20. But, if it was colored Red, the ANSI codes inject 9 invisible bytes into the string. Go evaluates `len()` as 29 bytes. **This completely destroys geometric centering math loops making words drift wildly off-center!**

### Solution: `VisibleLen()`
In `align.go`, we created an algorithm called `VisibleLen`. It parses strings character by character. If it detects `\033`, it silences the byte counter recursively until it hits the character `m` (ANSI string closer). The alignment tool only calculates spacing using bytes that are purely visible to human eyes!

### Base Layout Equations (`AlignLine`)
1. **Right Alignment**:
   `PADDING_SPACING = (TerminalWidth - WordWidth)`
   The program generates that difference in empty spaces, adding them directly to the `LEFT` of the graphic pushing it perfectly aligned right.
2. **Center Alignment**:
   `PADDING_SPACING = (TerminalWidth - WordWidth) / 2`
   We distribute space equally.

### Advanced Equations (`JustifyLine`)
Justifying separates targeted multiple words across margins stretching edge-to-edge seamlessly.
1. The engine (`AlignJustify` code block) splits your text into physical word batches.
2. It accumulates precisely how wide **each** word graphic currently calculates (`totalWordWidth`).
3. We discover our usable white-space: `TotalSpaces = (TerminalWidth - totalWordWidth)`.
4. We evaluate integer distributions across boundaries calculating gaps. (e.g. 5 words have exactly 4 spacer gaps).
   `Spaces_Per_Gap = TotalSpaces / Gaps`
5. **Handling Mathematics Remainders:** Since division decimals don't translate to block-spacing uniformly, the remainder logic `%` triggers passing 1 supplementary space cleanly starting sequentially leftwards `ExtraSpaces = TotalSpaces % Gaps` ensuring seamless boundaries uniformly edge-to-edge.

---

## 5. Storage Processing & Output (`main.go`)

### Writing Data Instead of Printing
Finally, the massive variable holding all 8 rows filled via spatial padding, ansi color injections, and multi-newline tracking evaluates if you supplied the `--output` command.
- If no destination was identified: `fmt.Print` triggers, unleashing the structure directly to your current console output successfully.
- If `--output=out.txt` triggered logic checks earlier: Our code skips the terminal render, formatting via system standard libraries `os.WriteFile(outputFile, []byte(result), 0644)` piping the block precisely to disk enabling you to archive exactly what generated seamlessly.

---

## 6. Testing Philosophy (`ascii_test.go`)
How do we know 100% of these parameters calculate effectively simultaneously without overriding each other causing chaotic regressions? 

**The Test-Runner Component**:
In the development pipeline, we preserved `ascii.go` as a facade method pointing to our newly crafted `GenerateAsciiArt` engine safely loaded to explicitly default parameters. 
- Over 7 robust test configurations target the package triggering via `go test ./ascii/... -v`. 
- This automatically evaluates edge-cases structurally, tracking input patterns like `"1 \n\n world "`, assessing graphic block formations line-by-line confirming our 8-row layout builders aren't shifting bytes inappropriately. 

---

### Endnote
By separating functions into single-job tools—`main.go` sorting definitions, `width.go` talking to operating systems, and `generate.go` executing math formulas over strings—we effectively designed an unbreakable pattern architecture scaling ASCII arts optimally.
