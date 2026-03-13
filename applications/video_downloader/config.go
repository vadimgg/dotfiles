package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

// ─── Config structs ───────────────────────────────────────────────────────────

type Config struct {
	Instagram InstagramConfig `yaml:"instagram"`
	YouTube   YouTubeConfig   `yaml:"youtube"`
	Advanced  AdvancedConfig  `yaml:"advanced"`
}

type InstagramConfig struct {
	OutputDir      string `yaml:"output_dir"`
	Format         string `yaml:"format"`
	MergeFormat    string `yaml:"merge_format"`
	Retries        int    `yaml:"retries"`
	FileNamePrefix string `yaml:"file_name_prefix"`
	Auth           struct {
		Username string `yaml:"username"`
		Password string `yaml:"password"`
	} `yaml:"auth"`
}

type YouTubeConfig struct {
	OutputDir      string `yaml:"output_dir"`
	Format         string `yaml:"format"`
	MergeFormat    string `yaml:"merge_format"`
	Retries        int    `yaml:"retries"`
	FileNamePrefix string `yaml:"file_name_prefix"`
	Resolution     string `yaml:"resolution"`       // e.g. "1080", "720", "best"
	AudioOnly      bool   `yaml:"audio_only"`       // download as mp3
	Subtitles      bool   `yaml:"subtitles"`        // embed subtitles if available
	SubtitleLangs  string `yaml:"subtitle_langs"`   // e.g. "en,es"
}

type AdvancedConfig struct {
	YtDlpPath      string `yaml:"yt_dlp_path"`
	RateLimit      string `yaml:"rate_limit"`
	SleepInterval  int    `yaml:"sleep_interval"`
	SkipIfExists   bool   `yaml:"skip_if_exists"`
	WriteMetadata  bool   `yaml:"write_metadata"`
	WriteThumbnail bool   `yaml:"write_thumbnail"`
}

// ─── Defaults ─────────────────────────────────────────────────────────────────

func defaultConfig() Config {
	var cfg Config

	cfg.Instagram.OutputDir = "~/Downloads/Instagram"
	cfg.Instagram.Format = "bestvideo+bestaudio/best"
	cfg.Instagram.MergeFormat = "mp4"
	cfg.Instagram.Retries = 3
	cfg.Instagram.FileNamePrefix = "instagram"

	cfg.YouTube.OutputDir = "~/Downloads/YouTube"
	cfg.YouTube.Format = "bestvideo+bestaudio/best"
	cfg.YouTube.MergeFormat = "mp4"
	cfg.YouTube.Retries = 3
	cfg.YouTube.FileNamePrefix = "youtube"
	cfg.YouTube.Resolution = "best"
	cfg.YouTube.AudioOnly = false
	cfg.YouTube.Subtitles = false
	cfg.YouTube.SubtitleLangs = "en"

	cfg.Advanced.YtDlpPath = "yt-dlp"
	cfg.Advanced.SkipIfExists = true
	cfg.Advanced.WriteMetadata = false
	cfg.Advanced.WriteThumbnail = false

	return cfg
}

// ─── Loader ───────────────────────────────────────────────────────────────────

func loadConfig(path string) (Config, error) {
	cfg := defaultConfig()

	data, err := os.ReadFile(path)
	if err != nil {
		return cfg, err
	}

	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return cfg, fmt.Errorf("failed to parse %s: %w", path, err)
	}
	return cfg, nil
}

// ─── XDG path ─────────────────────────────────────────────────────────────────

// xdgConfigPath returns the platform-appropriate config file path:
//   Linux/macOS : $XDG_CONFIG_HOME/video-downloader/config.yml
//                 (defaults to ~/.config/video-downloader/config.yml)
//   Windows     : %APPDATA%\video-downloader\config.yml
func xdgConfigPath() string {
	if base := os.Getenv("XDG_CONFIG_HOME"); base != "" {
		return filepath.Join(base, "video-downloader", "config.yml")
	}
	if base := os.Getenv("APPDATA"); base != "" {
		return filepath.Join(base, "video-downloader", "config.yml")
	}
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".config", "video-downloader", "config.yml")
}

// ─── Shared helpers ───────────────────────────────────────────────────────────

// expandHome expands a leading ~/ to the user's home directory.
func expandHome(path string) string {
	if strings.HasPrefix(path, "~/") {
		home, _ := os.UserHomeDir()
		return filepath.Join(home, path[2:])
	}
	return path
}

// prepareOutputDir expands ~ and creates the directory if it doesn't exist.
func prepareOutputDir(dir string) (string, error) {
	dir = expandHome(dir)
	abs, err := filepath.Abs(dir)
	if err != nil {
		return "", err
	}
	if err := os.MkdirAll(abs, 0755); err != nil {
		return "", err
	}
	return abs, nil
}

// printConfigStatus logs whether the config was loaded, missing, or broken.
func printConfigStatus(path string, err error) {
	if err == nil {
		fmt.Printf("✅ Config loaded: %s\n", path)
	} else if os.IsNotExist(err) {
		fmt.Printf("ℹ️  No config found at %q — using defaults.\n", path)
	} else {
		fmt.Printf("⚠️  Could not load config: %v\n", err)
	}
}
