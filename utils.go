package main

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// findHidrawDevicesByVIDPID scans /sys/class/hidraw to find all device nodes
// that match the specified vendor and product ID.
func findHidrawDevicesByVIDPID(targetVID, targetPID string) ([]string, error) {
	var devices []string
	hidrawPath := "/sys/class/hidraw"

	// Parse the target hex strings into integers for reliable comparison.
	targetVendorID, err := strconv.ParseInt(targetVID, 16, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid target vendor ID: %s", targetVID)
	}
	targetProductID, err := strconv.ParseInt(targetPID, 16, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid target product ID: %s", targetPID)
	}

	entries, err := os.ReadDir(hidrawPath)
	if err != nil {
		if os.IsNotExist(err) {
			log.Printf("WARN: %s does not exist, skipping device scan.", hidrawPath)
			return nil, nil
		}
		return nil, fmt.Errorf("failed to read %s: %w", hidrawPath, err)
	}

	for _, entry := range entries {
		// The path to the uevent file containing device identifiers.
		ueventPath := filepath.Join(hidrawPath, entry.Name(), "device", "uevent")

		content, err := os.ReadFile(ueventPath)
		if err != nil {
			continue // Device might have been unplugged, or path doesn't exist.
		}

		scanner := bufio.NewScanner(bytes.NewReader(content))
		for scanner.Scan() {
			line := scanner.Text()
			if !strings.HasPrefix(line, "HID_ID=") {
				continue
			}

			// Line format is "HID_ID=BUS:VENDOR:PRODUCT"
			// Example: "HID_ID=0003:0000046D:0000C077"
			parts := strings.Split(line, ":")
			if len(parts) != 3 {
				continue
			}

			vendorHex := parts[1]
			productHex := parts[2]

			vendorID, err := strconv.ParseInt(vendorHex, 16, 64)
			if err != nil {
				continue
			}
			productID, err := strconv.ParseInt(productHex, 16, 64)
			if err != nil {
				continue
			}

			if vendorID == targetVendorID && productID == targetProductID {
				deviceNode := filepath.Join("/dev", entry.Name())
				log.Printf("✅ Found matching passkey device: %s", deviceNode)
				devices = append(devices, deviceNode)
				break // Found match for this hidraw device, move to the next one.
			}
		}
	}

	if len(devices) == 0 {
		log.Println("⚠️ No matching passkey devices found. Is the key plugged in?")
	}

	return devices, nil
}
