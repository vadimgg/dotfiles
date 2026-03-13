package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

func downloadYouTube(url, configPath, outputDirOverride, formatOverride string) {
	cfg, err := loadConfig(configPath)
	printConfigStatus(configPath, err)

	yt := cfg.YouTube

	// CLI overrides
	if outputDirOverride != "" {
		yt.OutputDir = outputDirOverride
	}
	if formatOverride != "" {
		yt.Format = formatOverride
	}

	// Ensure yt-dlp is available
	ytDlp := cfg.Advanced.YtDlpPath
	if _, err := exec.LookPath(ytDlp); err != nil {
		fmt.Printf("❌ yt-dlp not found at %q. Install it with:\n", ytDlp)
		fmt.Println("   pip install yt-dlp   or   brew install yt-dlp")
		os.Exit(1)
	}

	absOutputDir, err := prepareOutputDir(yt.OutputDir)
	if err != nil {
		fmt.Printf("❌ Could not prepare output directory: %v\n", err)
		os.Exit(1)
	}

	outputTemplate := filepath.Join(absOutputDir, yt.FileNamePrefix+"_%(title)s_%(id)s.%(ext)s")

	// Build format string — audio_only and resolution override the config format
	format := yt.Format
	if yt.AudioOnly {
		format = "bestaudio/best"
	} else if yt.Resolution != "" && yt.Resolution != "best" {
		format = fmt.Sprintf("bestvideo[height<=%s]+bestaudio/best[height<=%s]", yt.Resolution, yt.Resolution)
	}

	args := []string{
		"--format", format,
		"--output", outputTemplate,
		"--retries", fmt.Sprintf("%d", yt.Retries),
		"--progress",
	}

	// Audio-only: extract and convert to mp3
	if yt.AudioOnly {
		args = append(args,
			"--extract-audio",
			"--audio-format", "mp3",
			"--audio-quality", "0", // best quality
		)
		fmt.Println("🎵 Audio-only mode — extracting as mp3.")
	} else {
		args = append(args, "--merge-output-format", yt.MergeFormat)
	}

	// Subtitles
	if yt.Subtitles {
		args = append(args,
			"--write-subs",
			"--embed-subs",
			"--sub-langs", yt.SubtitleLangs,
		)
	}

	if cfg.Advanced.SkipIfExists {
		args = append(args, "--no-overwrites")
	}
	if cfg.Advanced.WriteMetadata {
		args = append(args, "--write-info-json")
	}
	if cfg.Advanced.WriteThumbnail {
		args = append(args, "--write-thumbnail")
	}
	if cfg.Advanced.RateLimit != "" {
		args = append(args, "--rate-limit", cfg.Advanced.RateLimit)
	}
	if cfg.Advanced.SleepInterval > 0 {
		args = append(args, "--sleep-interval", fmt.Sprintf("%d", cfg.Advanced.SleepInterval))
	}

	args = append(args, url)

	fmt.Printf("⬇️  Downloading: %s\n", url)
	fmt.Printf("📁 Saving to:   %s\n\n", absOutputDir)

	cmd := exec.Command(ytDlp, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		fmt.Printf("\n❌ Download failed: %v\n", err)
		fmt.Println("\nTips:")
		fmt.Println("  • Age-restricted videos may require cookies: add --cookies-from-browser chrome to yt_dlp_path args.")
		fmt.Println("  • Private videos are not downloadable.")
		fmt.Println("  • Keep yt-dlp updated: pip install -U yt-dlp")
		os.Exit(1)
	}

	fmt.Printf("\n✅ Done! File saved in: %s\n", absOutputDir)
}
