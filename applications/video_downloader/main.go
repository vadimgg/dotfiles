// Video Downloader & Scene Analyzer
// Automatically detects YouTube or Instagram links and routes accordingly.
// Can also analyze scenes in a local video file directly.
//
// Setup (one-time):
//   go mod init video_downloader
//   go mod tidy
//
// Run:
//   go run . <url>                      — download YouTube or Instagram
//   go run . --analyze <video.mp4>      — analyze scenes in a local file
//   go run . --auto-analyze <url>       — download then auto-run scene analysis
//
// Build:
//   go build -o video_downloader .
//   ./video_downloader <url>

package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	configPath   := flag.String("config",       xdgConfigPath(), "Path to config.yml")
	outputDir    := flag.String("output",       "",              "Override output directory")
	format       := flag.String("format",       "",              "Override yt-dlp format string")
	analyzeOnly  := flag.Bool("analyze",        false,           "Analyze scenes in a local file instead of downloading")
	autoAnalyze  := flag.Bool("auto-analyze",   false,           "Download then immediately run scene analysis")
	flag.Parse()

	// ── Mode: analyze a local file directly ──────────────────────────────────
	if *analyzeOnly {
		if flag.NArg() < 1 {
			printError("Please provide a video file path.")
			fmt.Println("  Usage: video_downloader --analyze <video.mp4>")
			os.Exit(1)
		}
		videoPath := flag.Arg(0)
		videoDir  := filepath.Dir(videoPath)
		if videoDir == "" || videoDir == "." {
			videoDir, _ = os.Getwd()
		}
		RunSceneAnalysis(videoPath, videoDir)
		return
	}

	// ── Mode: download (and optionally analyze) ───────────────────────────────
	if flag.NArg() < 1 {
		fmt.Println(cBold("Usage:") + "   video_downloader [flags] " + cURL("<url or file>"))
		fmt.Println()
		fmt.Println(cBold("Modes:"))
		fmt.Printf("  %-36s %s\n", cURL("<url>"),                      "Download YouTube or Instagram video")
		fmt.Printf("  %-36s %s\n", "--analyze "+cURL("<video.mp4>"),   "Analyze scenes in a local file")
		fmt.Printf("  %-36s %s\n", "--auto-analyze "+cURL("<url>"),    "Download then auto-run scene analysis")
		fmt.Println()
		fmt.Println(cBold("Supported platforms:"))
		fmt.Println("  " + cSuccess("YouTube")   + "   — youtube.com, youtu.be")
		fmt.Println("  " + cSuccess("Instagram") + " — instagram.com (posts, reels, stories)")
		fmt.Println()
		fmt.Println(cBold("Flags:"))
		flag.PrintDefaults()
		fmt.Println()
		fmt.Println(cBold("Config:"))
		fmt.Printf("  %s\n", cPath(xdgConfigPath()))
		os.Exit(1)
	}

	url    := flag.Arg(0)
	reader := bufio.NewReader(os.Stdin)

	for {
		var filePath, outDir string

		switch detectPlatform(url) {
		case "youtube":
			fmt.Println(cSuccess("🎬 Detected: ") + cBold("YouTube"))
			filePath, outDir = downloadYouTube(url, *configPath, *outputDir, *format)

		case "instagram":
			fmt.Println(cSuccess("📸 Detected: ") + cBold("Instagram"))
			filePath, outDir = downloadInstagram(url, *configPath, *outputDir, *format)

		default:
			printError("Unsupported URL: %s", cURL(url))
			fmt.Println(cDim("   Supported platforms: YouTube, Instagram"))
			os.Exit(1)
		}

		// --auto-analyze: skip the menu, run scene analysis immediately
		if *autoAnalyze && filePath != "" {
			RunSceneAnalysis(filePath, outDir)
			return
		}

		// Interactive post-download menu
		if another := PostDownloadMenu(filePath, outDir); another {
			url = promptURL(reader)
			if url == "" {
				fmt.Println(cDim("  Bye! 👋"))
				return
			}
		} else {
			return
		}
	}
}

// promptURL asks the user to paste a new URL and returns it (trimmed).
func promptURL(reader *bufio.Reader) string {
	fmt.Printf("\n%s ", cBold("  Paste URL → "))
	line, err := reader.ReadString('\n')
	if err != nil {
		return ""
	}
	return strings.TrimSpace(line)
}

// detectPlatform returns "youtube", "instagram", or "unknown".
func detectPlatform(url string) string {
	switch {
	case strings.Contains(url, "youtube.com") || strings.Contains(url, "youtu.be"):
		return "youtube"
	case strings.Contains(url, "instagram.com"):
		return "instagram"
	default:
		return "unknown"
	}
}
