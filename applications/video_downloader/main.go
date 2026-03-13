// Video Downloader
// Automatically detects YouTube or Instagram links and routes accordingly.
//
// Setup (one-time):
//   go mod init video_downloader
//   go mod tidy
//
// Run:
//   go run . <url>
//   go run . -config ~/.config/video-downloader/config.yml <url>
//   go run . -output ~/Videos <url>
//
// Build:
//   go build -o video_downloader .
//   ./video_downloader <url>

package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

func main() {
	configPath := flag.String("config", xdgConfigPath(), "Path to config.yml")
	outputDir  := flag.String("output", "", "Override output directory")
	format     := flag.String("format", "", "Override yt-dlp format string")
	flag.Parse()

	if flag.NArg() < 1 {
		fmt.Println("Usage:   video_downloader [flags] <url>")
		fmt.Println()
		fmt.Println("Supported platforms:")
		fmt.Println("  • YouTube   — youtube.com, youtu.be")
		fmt.Println("  • Instagram — instagram.com (posts, reels, stories)")
		fmt.Println()
		fmt.Println("Flags:")
		flag.PrintDefaults()
		fmt.Println()
		fmt.Println("Config files:")
		fmt.Printf("  %s\n", xdgConfigPath())
		os.Exit(1)
	}

	url := flag.Arg(0)

	switch detectPlatform(url) {
	case "youtube":
		fmt.Println("🎬 Detected: YouTube")
		downloadYouTube(url, *configPath, *outputDir, *format)

	case "instagram":
		fmt.Println("📸 Detected: Instagram")
		downloadInstagram(url, *configPath, *outputDir, *format)

	default:
		fmt.Printf("❌ Unsupported URL: %q\n", url)
		fmt.Println("   Supported platforms: YouTube, Instagram")
		os.Exit(1)
	}
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
