package main

import (
	"fmt"
	"os"
	"os/exec"
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
	YtDlpPath          string `yaml:"yt_dlp_path"`
	RateLimit          string `yaml:"rate_limit"`
	SleepInterval      int    `yaml:"sleep_interval"`
	SkipIfExists       bool   `yaml:"skip_if_exists"`
	WriteMetadata      bool   `yaml:"write_metadata"`
	WriteThumbnail     bool   `yaml:"write_thumbnail"`
	CookiesFromBrowser string `yaml:"cookies_from_browser"`
	CookiesFile        string `yaml:"cookies_file"`
	JsRuntime          string `yaml:"js_runtime"`
	RemoteComponents   string `yaml:"remote_components"` // e.g. "ejs:github", "ejs:npm"
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
	cfg.Advanced.WriteMetadata = true
	cfg.Advanced.WriteThumbnail = true
	cfg.Advanced.CookiesFromBrowser = "safari"
	cfg.Advanced.JsRuntime = detectJsRuntime()
	cfg.Advanced.RemoteComponents = "ejs:github"

	return cfg
}

// detectJsRuntime returns the first JS runtime found on PATH, or empty string.
func detectJsRuntime() string {
	for _, rt := range []string{"deno", "node"} {
		if _, err := exec.LookPath(rt); err == nil {
			return rt
		}
	}
	return ""
}

// ─── Loader ───────────────────────────────────────────────────────────────────

func loadConfig(path string) (Config, error) {
	cfg := defaultConfig()

	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		// First run: write the default config to disk so the user can edit it
		if writeErr := writeDefaultConfig(path, cfg); writeErr == nil {
			printInfo("Created default config at %s", cPath(path))
		}
		return cfg, nil
	}
	if err != nil {
		return cfg, err
	}

	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return cfg, fmt.Errorf("failed to parse %s: %w", path, err)
	}
	return cfg, nil
}

// writeDefaultConfig marshals cfg to YAML and writes it to path (creating
// parent directories as needed).
func writeDefaultConfig(path string, cfg Config) error {
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}
	data, err := yaml.Marshal(cfg)
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
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

// slugify converts a string to a safe directory-name segment:
//   - lowercase
//   - spaces → _
//   - strips characters that are not alphanumeric, _, or -
//   - collapses multiple consecutive underscores
//   - trims leading/trailing underscores
//   - caps length at maxLen runes
func slugify(s string, maxLen int) string {
	s = strings.ToLower(s)
	var b strings.Builder
	prevUnderscore := false
	for _, r := range s {
		var c rune
		switch {
		case r >= 'a' && r <= 'z', r >= '0' && r <= '9', r == '-':
			c = r
			prevUnderscore = false
		case r == ' ' || r == '_' || r == '.' || r == ':' || r == ',' || r == '(' || r == ')' || r == '[' || r == ']':
			c = '_'
		default:
			// skip everything else (emoji, accents handled below by keeping letters)
			if r > 127 {
				// keep unicode letters as-is so non-latin titles aren't wiped out
				c = r
				prevUnderscore = false
			} else {
				continue
			}
		}
		if c == '_' {
			if prevUnderscore {
				continue
			}
			prevUnderscore = true
		}
		b.WriteRune(c)
		if b.Len() >= maxLen {
			break
		}
	}
	return strings.Trim(b.String(), "_")
}

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

// authArgs returns the yt-dlp flags for cookies and JS runtime from config.
// Used by both youtube.go and instagram.go so they stay in sync.
func authArgs(adv AdvancedConfig) []string {
	var args []string
	if adv.CookiesFromBrowser != "" {
		args = append(args, "--cookies-from-browser", adv.CookiesFromBrowser)
	}
	if adv.CookiesFile != "" {
		args = append(args, "--cookies", expandHome(adv.CookiesFile))
	}
	if adv.JsRuntime != "" {
		args = append(args, "--js-runtimes", adv.JsRuntime)
	}
	if adv.RemoteComponents != "" {
		args = append(args, "--remote-components", adv.RemoteComponents)
	}
	return args
}

// printConfigStatus logs whether the config was loaded or had a parse error.
// Missing config is handled silently in loadConfig (auto-created).
func printConfigStatus(path string, err error) {
	if err == nil {
		return // no noise on success
	}
	printWarn("Could not load config (%s): %v — using defaults.", cPath(path), err)
}
