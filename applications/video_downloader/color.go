package main

import (
	"fmt"
	"os"
)

// ANSI escape codes
const (
	reset   = "\033[0m"
	bold    = "\033[1m"
	dim     = "\033[2m"

	red     = "\033[31m"
	green   = "\033[32m"
	yellow  = "\033[33m"
	blue    = "\033[34m"
	magenta = "\033[35m"
	cyan    = "\033[36m"
	white   = "\033[37m"
)

// colorsSupported returns false when stdout is not a real terminal
// (e.g. redirected to a file), so we never pollute piped output with
// raw escape codes.
func colorsSupported() bool {
	fi, err := os.Stdout.Stat()
	if err != nil {
		return false
	}
	return (fi.Mode() & os.ModeCharDevice) != 0
}

func colorize(color, s string) string {
	if !colorsSupported() {
		return s
	}
	return color + s + reset
}

// Semantic helpers used throughout the app.

func cSuccess(s string) string { return colorize(bold+green, s) }
func cError(s string) string   { return colorize(bold+red, s) }
func cWarn(s string) string    { return colorize(bold+yellow, s) }
func cInfo(s string) string    { return colorize(cyan, s) }
func cDim(s string) string     { return colorize(dim, s) }
func cBold(s string) string    { return colorize(bold, s) }
func cURL(s string) string     { return colorize(magenta, s) }
func cPath(s string) string    { return colorize(blue, s) }

// Convenience print wrappers

func printSuccess(format string, a ...any) {
	fmt.Println(cSuccess("✅ " + fmt.Sprintf(format, a...)))
}

func printError(format string, a ...any) {
	fmt.Println(cError("❌ " + fmt.Sprintf(format, a...)))
}

func printWarn(format string, a ...any) {
	fmt.Println(cWarn("⚠️  " + fmt.Sprintf(format, a...)))
}

func printInfo(format string, a ...any) {
	fmt.Println(cInfo("ℹ️  " + fmt.Sprintf(format, a...)))
}

func printTip(s string) {
	fmt.Println(cDim("  • " + s))
}
