package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

// PostDownloadMenu shows an interactive menu after a successful download.
// Re-displays after every action except [c] (download another) and [q] (exit).
// Returns true if the user wants to download another video.
func PostDownloadMenu(filePath, outputDir string) bool {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Println()
		fmt.Println(strings.Repeat("─", 48))
		fmt.Println(cBold("  What would you like to do?"))
		fmt.Println(strings.Repeat("─", 48))
		fmt.Printf("  %s  Analyze scenes\n",   cBold(colorize(cyan, "[a]")))
		fmt.Printf("  %s  Open folder\n",      cBold(colorize(cyan, "[o]")))
		fmt.Printf("  %s  Play in IINA\n",     cBold(colorize(cyan, "[r]")))
		fmt.Printf("  %s  Download another\n", cBold(colorize(cyan, "[c]")))
		fmt.Printf("  %s  Exit\n",             cBold(colorize(cyan, "[q]")))
		fmt.Println(strings.Repeat("─", 48))
		fmt.Print(cBold("  → "))

		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println()
			return false
		}

		switch strings.TrimSpace(strings.ToLower(input)) {
		case "a":
			if filePath == "" {
				printWarn("Could not determine video file path — skipping scene analysis.")
			} else {
				RunSceneAnalysis(filePath, outputDir)
			}
			// loop: re-show menu

		case "o":
			openFolder(outputDir)
			// loop: re-show menu

		case "r":
			if filePath == "" {
				printWarn("Could not determine file path — opening folder instead.")
				openFolder(outputDir)
			} else {
				openIINA(filePath)
			}
			// loop: re-show menu

		case "c":
			fmt.Println()
			return true

		case "q", "":
			fmt.Println(cDim("  Bye! 👋"))
			return false

		default:
			fmt.Printf("  %s — press %s, %s, %s, %s, or %s\n",
				cWarn("Unknown option"),
				cBold(colorize(cyan, "a")),
				cBold(colorize(cyan, "o")),
				cBold(colorize(cyan, "r")),
				cBold(colorize(cyan, "c")),
				cBold(colorize(cyan, "q")),
			)
			// loop: re-show menu
		}
	}
}

// openFolder opens the output directory in the system file manager.
func openFolder(dir string) {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("open", dir)
	case "windows":
		cmd = exec.Command("explorer", dir)
	default:
		cmd = exec.Command("xdg-open", dir)
	}
	if err := cmd.Start(); err != nil {
		printError("Could not open folder: %v", err)
	} else {
		printSuccess("Opened folder: %s", cPath(dir))
	}
}

// openIINA opens the given file in IINA (macOS), falling back to system default.
func openIINA(filePath string) {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "darwin":
		if _, err := exec.LookPath("iina"); err == nil {
			cmd = exec.Command("iina", filePath)
		} else if _, err := exec.LookPath("/Applications/IINA.app/Contents/MacOS/iina"); err == nil {
			cmd = exec.Command("/Applications/IINA.app/Contents/MacOS/iina", filePath)
		} else {
			printWarn("IINA not found — opening with default player.")
			cmd = exec.Command("open", filePath)
		}
	case "windows":
		cmd = exec.Command("cmd", "/c", "start", "", filePath)
	default:
		cmd = exec.Command("xdg-open", filePath)
	}
	if err := cmd.Start(); err != nil {
		printError("Could not open file: %v", err)
	} else {
		printSuccess("Opening: %s", cPath(filepath.Base(filePath)))
	}
}
