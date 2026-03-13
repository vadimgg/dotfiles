package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

func downloadInstagram(url, configPath, outputDirOverride, formatOverride string) {
	cfg, err := loadConfig(configPath)
	printConfigStatus(configPath, err)

	ig := cfg.Instagram

	// CLI overrides
	if outputDirOverride != "" {
		ig.OutputDir = outputDirOverride
	}
	if formatOverride != "" {
		ig.Format = formatOverride
	}

	// Env-var credentials override config file
	if v := os.Getenv("IG_USERNAME"); v != "" {
		ig.Auth.Username = v
	}
	if v := os.Getenv("IG_PASSWORD"); v != "" {
		ig.Auth.Password = v
	}

	// Ensure yt-dlp is available
	ytDlp := cfg.Advanced.YtDlpPath
	if _, err := exec.LookPath(ytDlp); err != nil {
		fmt.Printf("❌ yt-dlp not found at %q. Install it with:\n", ytDlp)
		fmt.Println("   pip install yt-dlp   or   brew install yt-dlp")
		os.Exit(1)
	}

	absOutputDir, err := prepareOutputDir(ig.OutputDir)
	if err != nil {
		fmt.Printf("❌ Could not prepare output directory: %v\n", err)
		os.Exit(1)
	}

	outputTemplate := filepath.Join(absOutputDir, ig.FileNamePrefix+"_%(id)s.%(ext)s")

	args := []string{
		"--format", ig.Format,
		"--merge-output-format", ig.MergeFormat,
		"--output", outputTemplate,
		"--retries", fmt.Sprintf("%d", ig.Retries),
		"--progress",
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
		fmt.Println("🔐 Using Instagram credentials.")
		args = append(args, "--username", ig.Auth.Username, "--password", ig.Auth.Password)
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
		fmt.Println("  • Private posts require credentials in config.yml or IG_USERNAME/IG_PASSWORD env vars.")
		fmt.Println("  • Stories expire after 24 hours.")
		fmt.Println("  • Keep yt-dlp updated: pip install -U yt-dlp")
		os.Exit(1)
	}

	fmt.Printf("\n✅ Done! File saved in: %s\n", absOutputDir)
}
