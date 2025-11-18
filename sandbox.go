package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
)

const browserBinary = "chromium"

const passkeyVID = "20a0"
const passkeyPID = "42d4"

func LaunchSandboxedBrowser(url string) error {
	for _, cmd := range []string{"bwrap", browserBinary} {
		if _, err := exec.LookPath(cmd); err != nil {
			return fmt.Errorf("required command '%s' not found in PATH", cmd)
		}
	}

	log.Println("🔎 Scanning for passkey devices...")
	_, _ = findHidrawDevicesByVIDPID(passkeyVID, passkeyPID)

	profileAndSocketDir, err := os.MkdirTemp("", "canokey-chromium-data.*")
	if err != nil {
		return fmt.Errorf("failed to create temp profile dir: %w", err)
	}
	defer os.RemoveAll(profileAndSocketDir)

	uid := os.Getuid()
	runUserPath := fmt.Sprintf("/run/user/%d", uid)
	waylandDisplay := os.Getenv("WAYLAND_DISPLAY")
	dbusAddr := os.Getenv("DBUS_SESSION_BUS_ADDRESS")

	log.Println("🔒 Preparing strict sandbox and launching browser...")

	args := []string{
		"--unshare-all",
		"--share-net",
		"--die-with-parent",
		"--proc", "/proc",
		"--dev-bind", "/dev", "/dev",

		"--ro-bind", "/sys/devices", "/sys/devices",
		"--ro-bind", "/sys/bus/usb", "/sys/bus/usb",

		"--bind", "/run/dbus/system_bus_socket", "/run/dbus/system_bus_socket",

		"--ro-bind", "/usr/bin", "/usr/bin",
		"--ro-bind", "/usr/lib", "/usr/lib",
		"--ro-bind", "/usr/lib64", "/usr/lib64",
		"--ro-bind", "/usr/share", "/usr/share",
		"--ro-bind", "/etc/machine-id", "/etc/machine-id",

		"--symlink", "usr/bin", "/bin",
		"--symlink", "usr/lib", "/lib",
		"--symlink", "usr/lib64", "/lib64",
		"--symlink", "usr/bin", "/sbin",

		"--bind", runUserPath, runUserPath,
	}

	args = append(args,
		"--setenv", "XDG_RUNTIME_DIR", runUserPath,
		"--setenv", "WAYLAND_DISPLAY", waylandDisplay,
		"--setenv", "DBUS_SESSION_BUS_ADDRESS", dbusAddr,

		browserBinary,
		"--user-data-dir="+profileAndSocketDir,
		"--process-singleton-dir="+profileAndSocketDir,
		"--new-window",
		"--no-first-run",
		"--disable-gpu",
		"--enable-features=WebUSB,UseOzonePlatform",
		"--app="+url,
	)

	cmd := exec.Command("bwrap", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}
