package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

// downloadInstagram downloads an Instagram video and returns the file path and output dir.
func downloadInstagram(url, configPath, outputDirOverride, formatOverride string) (filePath, outputDir string) {
	cfg, err := loadConfig(configPath)
	printConfigStatus(configPath, err)

	ig := cfg.Instagram

	if outputDirOverride != "" {
		ig.OutputDir = outputDirOverride
	}
	if formatOverride != "" {
		ig.Format = formatOverride
	}

	if v := os.Getenv("IG_USERNAME"); v != "" {
		ig.Auth.Username = v
	}
	if v := os.Getenv("IG_PASSWORD"); v != "" {
		ig.Auth.Password = v
	}

	ytDlp := cfg.Advanced.YtDlpPath
	if _, err := exec.LookPath(ytDlp); err != nil {
		printError("yt-dlp not found at %q. Install it with:", ytDlp)
		printTip("pip install yt-dlp   or   brew install yt-dlp")
		os.Exit(1)
	}

	// ── Fetch metadata ───────────────────────────────────────────────────────
	fmt.Printf("%s %s\n", cBold("🔍 Fetching metadata:"), cURL(url))

	metaExtraArgs := authArgs(cfg.Advanced)
	if ig.Auth.Username != "" && ig.Auth.Password != "" {
		metaExtraArgs = append(metaExtraArgs, "--username", ig.Auth.Username, "--password", ig.Auth.Password)
	}

	meta, err := fetchMeta(ytDlp, url, metaExtraArgs)
	if err != nil {
		printWarn("Could not fetch metadata (%v) — using fallback folder name.", err)
		meta = videoMeta{Title: "untitled", UploadDate: time.Now().Format("20060102")}
	} else {
		account := meta.Channel
		if account == "" {
			account = meta.Uploader
		}
		fmt.Printf("  %s %s\n",   cDim("Account:"), cBold(account))
		fmt.Printf("  %s %s\n\n", cDim("Title:  "), cBold(meta.Title))
	}

	folderName := buildFolderName(meta)

	// ── Prepare directories ──────────────────────────────────────────────────
	baseDir, err := prepareOutputDir(ig.OutputDir)
	if err != nil {
		printError("Could not prepare output directory: %v", err)
		os.Exit(1)
	}

	videoDir := filepath.Join(baseDir, folderName)
	if err := os.MkdirAll(videoDir, 0755); err != nil {
		printError("Could not create video folder: %v", err)
		os.Exit(1)
	}
	outputDir = videoDir

	outputTemplate := filepath.Join(videoDir, "%(id)s.%(ext)s")

	args := []string{
		"--format", ig.Format,
		"--merge-output-format", ig.MergeFormat,
		"--output", outputTemplate,
		"--retries", fmt.Sprintf("%d", ig.Retries),
		"--progress",
	}

	if cfg.Advanced.CookiesFromBrowser != "" {
		args = append(args, "--cookies-from-browser", cfg.Advanced.CookiesFromBrowser)
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
	if ig.Auth.Username != "" && ig.Auth.Password != "" {
		printInfo("Using Instagram credentials.")
		args = append(args, "--username", ig.Auth.Username, "--password", ig.Auth.Password)
	}

	args = append(args, url)

	fmt.Printf("%s %s\n",   cBold("⬇️  Downloading:"), cURL(url))
	fmt.Printf("%s %s\n\n", cBold("📁 Saving to:  "), cPath(videoDir))

	cmd := exec.Command(ytDlp, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		printError("Download failed: %v", err)
		fmt.Println(cBold("\nTips:"))
		printTip("Private posts require credentials in config.yml or IG_USERNAME/IG_PASSWORD env vars.")
		printTip("Stories expire after 24 hours.")
		printTip("Keep yt-dlp updated: " + cDim("pip install -U yt-dlp"))
		os.Exit(1)
	}

	filePath = resolveDownloadedFile(videoDir)
	printSuccess("Done! Saved to: %s", cPath(videoDir))
	return filePath, outputDir
}
