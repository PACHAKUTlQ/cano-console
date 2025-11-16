package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
)

const (
	baseURL    = "https://console.canokeys.org/"
	appDir     = "canokey-local-console"
	gstaticURL = "https://www.gstatic.com/flutter-canvaskit"
	userAgent  = "CanoKey-Local-Launcher/2.0"
)

var (
	engineRevisionRegex = regexp.MustCompile(`"engineRevision":\s*"([a-f0-9]{40})"`)
	loaderCallRegex     = regexp.MustCompile(`_flutter\.loader\.load\({[\s\S]*?}\);`)
)

type FontAsset struct {
	Asset  string `json:"asset"`
	Weight int    `json:"weight,omitempty"`
	Style  string `json:"style,omitempty"`
}

type FontFamily struct {
	Family string      `json:"family"`
	Fonts  []FontAsset `json:"fonts"`
}

func UpdateAssets(force bool) error {
	if _, err := os.Stat(appDir); os.IsNotExist(err) || force {
		log.Println("🚀 Assets not found or update forced. Starting download...")
		if force {
			_ = os.RemoveAll(appDir)
		}
		if err := fetchAllAssets(); err != nil {
			return err
		}
		log.Println("✅ All assets fetched and localized successfully.")
	} else {
		log.Println("✅ Local assets found. Skipping download. Use -update to refresh.")
	}
	return nil
}

func fetchAllAssets() error {
	coreFiles := []string{
		"index.html", "flutter_bootstrap.js", "flutter.js", "main.dart.js", "manifest.json",
		"favicon.png", "icons/Icon-192.png", "icons/Icon-512.png", "img/splash.png",
		"pkg/rust_lib_canokey_console.js",
		"pkg/rust_lib_canokey_console_bg.wasm",
		"version.json",
	}
	for _, file := range coreFiles {
		if err := downloadFile(file, file); err != nil {
			return fmt.Errorf("failed to download core file %s: %w", file, err)
		}
	}

	if err := fetchCanvasKitFiles(); err != nil {
		return err
	}

	// Download the self-hosted icon fonts (MaterialIcons, etc.)
	if err := fetchSelfHostedFonts(); err != nil {
		return err
	}

	// Download Google Fonts and patch the manifest to use our local copies.
	if err := fetchAndPatchFontManifest(); err != nil {
		return err
	}

	fetchOtherManifestAssets()

	// Patch JS to use local CanvasKit.
	if err := applyJSPatch(); err != nil {
		return err
	}

	return nil
}

// Downloads the fonts that are self-hosted on console.canokeys.org.
func fetchSelfHostedFonts() error {
	log.Println("🖋️  Downloading self-hosted UI icon fonts...")
	requiredFonts := []string{
		"fonts/MaterialIcons-Regular.otf",
		"packages/cupertino_icons/assets/CupertinoIcons.ttf",
		"packages/font_awesome_flutter/lib/fonts/fa-brands-400.ttf",
		"packages/font_awesome_flutter/lib/fonts/fa-regular-400.ttf",
		"packages/font_awesome_flutter/lib/fonts/fa-solid-900.ttf",
		"packages/lucide_icons/assets/lucide.ttf",
	}

	var wg sync.WaitGroup
	errs := make(chan error, len(requiredFonts))
	for _, fontPath := range requiredFonts {
		wg.Add(1)
		go func(p string) {
			defer wg.Done()
			if err := downloadFile("assets/"+p, "assets/"+p); err != nil {
				errs <- fmt.Errorf("failed to download critical font %s: %w", p, err)
			}
		}(fontPath)
	}
	wg.Wait()
	close(errs)

	if len(errs) > 0 {
		var errStrings []string
		for err := range errs {
			errStrings = append(errStrings, err.Error())
		}
		return errors.New(strings.Join(errStrings, "\n"))
	}
	return nil
}

// Downloads Google Fonts and injects them into the application's FontManifest.
func fetchAndPatchFontManifest() error {
	log.Println("🔩 Patching FontManifest with local Google Fonts...")
	manifestPath := filepath.Join(appDir, "assets", "FontManifest.json")
	if err := downloadFile("assets/FontManifest.json", "assets/FontManifest.json"); err != nil {
		return fmt.Errorf("could not download original FontManifest.json: %w", err)
	}
	content, err := os.ReadFile(manifestPath)
	if err != nil {
		return fmt.Errorf("failed to read downloaded FontManifest.json: %w", err)
	}

	var manifest []FontFamily
	if err := json.Unmarshal(content, &manifest); err != nil {
		return fmt.Errorf("failed to parse original FontManifest.json: %w", err)
	}

	googleFonts := map[string]string{
		"Roboto-Regular.woff2": "https://fonts.gstatic.com/s/roboto/v32/KFOmCnqEu92Fr1Me4GZLCzYlKw.woff2",
		// "Roboto-Bold.woff2":    "https://fonts.gstatic.com/s/roboto/v32/.woff2",
		// "Roboto-Italic.woff2":  "https://fonts.gstatic.com/s/roboto/v32/.woff2",
	}

	googleFontDir := filepath.Join(appDir, "assets", "google_fonts")
	for filename, url := range googleFonts {
		if err := downloadFileFromURL(url, filepath.Join(googleFontDir, filename)); err != nil {
			return fmt.Errorf("failed to download google font %s: %w", filename, err)
		}
	}

	robotoFamily := FontFamily{
		Family: "Roboto",
		Fonts: []FontAsset{
			{Asset: "google_fonts/Roboto-Regular.woff2", Weight: 400},
			// {Asset: "google_fonts/Roboto-Bold.woff2", Weight: 700},
			// {Asset: "google_fonts/Roboto-Italic.woff2", Weight: 400, Style: "italic"},
		},
	}

	found := false
	for i, family := range manifest {
		if family.Family == "Roboto" {
			manifest[i] = robotoFamily
			found = true
			break
		}
	}
	if !found {
		manifest = append(manifest, robotoFamily)
	}

	patchedContent, _ := json.MarshalIndent(manifest, "", "  ")
	if err := os.WriteFile(manifestPath, patchedContent, 0644); err != nil {
		return fmt.Errorf("failed to write patched font manifest: %w", err)
	}
	log.Println("✅ FontManifest.json patched successfully with local Roboto font.")
	return nil
}

func fetchCanvasKitFiles() error {
	bootstrapPath := filepath.Join(appDir, "flutter_bootstrap.js")
	bootstrapContent, err := os.ReadFile(bootstrapPath)
	if err != nil {
		return err
	}
	match := engineRevisionRegex.FindStringSubmatch(string(bootstrapContent))
	if len(match) < 2 {
		return errors.New("could not find engineRevision in flutter_bootstrap.js")
	}
	engineRevision := match[1]
	canvasKitDir := filepath.Join(appDir, "canvaskit", "chromium")
	canvasKitAssets := []string{"canvaskit.js", "canvaskit.wasm"}
	for _, asset := range canvasKitAssets {
		url := fmt.Sprintf("%s/%s/chromium/%s", gstaticURL, engineRevision, asset)
		dest := filepath.Join(canvasKitDir, asset)
		if err := downloadFileFromURL(url, dest); err != nil {
			return fmt.Errorf("failed to download CanvasKit asset %s: %w", asset, err)
		}
	}
	return nil
}

func fetchOtherManifestAssets() {
	if err := downloadFile("assets/AssetManifest.json", "assets/AssetManifest.json"); err != nil {
		return
	}
	manifestContent, err := os.ReadFile(filepath.Join(appDir, "assets", "AssetManifest.json"))
	if err != nil {
		return
	}
	var assetMap map[string][]string
	if err := json.Unmarshal(manifestContent, &assetMap); err != nil {
		return
	}
	for _, assetPaths := range assetMap {
		for _, p := range assetPaths {
			// Skip fonts, as we handle them separately and more robustly.
			if strings.Contains(p, "font") {
				continue
			}
			go func(assetPath string) {
				_ = downloadFile(assetPath, assetPath)
			}(p)
		}
	}
}

func applyJSPatch() error {
	bootstrapPath := filepath.Join(appDir, "flutter_bootstrap.js")
	content, err := os.ReadFile(bootstrapPath)
	if err != nil {
		return fmt.Errorf("failed to read flutter_bootstrap.js for patching: %w", err)
	}
	const replacement = `_flutter.loader.load({serviceWorkerSettings: null, config: {canvasKitBaseUrl: "/canvaskit/"}});`
	if !loaderCallRegex.Match(content) {
		return errors.New("could not find _flutter.loader.load call in flutter_bootstrap.js")
	}
	newContent := loaderCallRegex.ReplaceAll(content, []byte(replacement))
	if err := os.WriteFile(bootstrapPath, newContent, 0644); err != nil {
		return fmt.Errorf("failed to write patched flutter_bootstrap.js: %w", err)
	}
	log.Println("🔧 Patched JS to use local CanvasKit and disable service worker.")
	return nil
}

func downloadFile(urlPath, localPath string) error {
	return downloadFileFromURL(baseURL+urlPath, filepath.Join(appDir, localPath))
}

func downloadFileFromURL(url string, destPath string) error {
	log.Printf("  -> Downloading %s", filepath.Base(destPath))
	if err := os.MkdirAll(filepath.Dir(destPath), 0755); err != nil {
		return err
	}
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("User-Agent", userAgent)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status giving %d for %s", resp.StatusCode, url)
	}
	out, err := os.Create(destPath)
	if err != nil {
		return err
	}
	defer out.Close()
	_, err = io.Copy(out, resp.Body)
	return err
}
