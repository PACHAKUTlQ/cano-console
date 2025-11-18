package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
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

	const udevRulesPath = "/usr/lib/udev"
	const udevSocketPath = "/run/udev/io.systemd.Udev"

	args := []string{
		"--unshare-all",
		"--share-net",
		"--die-with-parent",
		"--proc", "/proc",
		"--dev", "/dev",

		// "--ro-bind", "/sys", "/sys",
		"--ro-bind", "/sys/devices", "/sys/devices",
		"--ro-bind", "/sys/bus/usb", "/sys/bus/usb",

		"--bind", udevSocketPath, udevSocketPath,
		"--ro-bind", udevRulesPath, udevRulesPath,
		"--ro-bind-try", "/etc/udev", "/etc/udev",

		"--ro-bind", "/run/udev/data", "/run/udev/data",
		"--bind", "/run/dbus/system_bus_socket", "/run/dbus/system_bus_socket",

		"--dev-bind", "/dev/bus/usb", "/dev/bus/usb",

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

	const maxHidrawDevices = 16
	log.Printf("... pre-emptively binding /dev/hidraw0 to /dev/hidraw%d for hot-plug support", maxHidrawDevices-1)
	for i := range maxHidrawDevices {
		devPath := "/dev/hidraw" + strconv.Itoa(i)
		// fmt.Printf("... binding %s\n", devPath)
		args = append(args, "--dev-bind-try", devPath, devPath)
	}

	log.Println("--- Running initial device scan ---")
	_, _ = findHidrawDevicesByVIDPID(passkeyVID, passkeyPID)
	log.Println("--- Initial scan complete ---")

	args = append(args,
		"--setenv", "XDG_RUNTIME_DIR", runUserPath,
		"--setenv", "WAYLAND_DISPLAY", waylandDisplay,
		"--setenv", "DBUS_SESSION_BUS_ADDRESS", dbusAddr,
	)

	browserArgs := []string{
		browserBinary,
		"--user-data-dir=" + profileAndSocketDir,
		"--process-singleton-dir=" + profileAndSocketDir,
		"--new-window",
		"--no-first-run",
		"--disable-gpu",
		"--enable-features=WebUSB,UseOzonePlatform",
		"--app=" + url,
	}
	args = append(args, browserArgs...)

	cmd := exec.Command("bwrap", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}
