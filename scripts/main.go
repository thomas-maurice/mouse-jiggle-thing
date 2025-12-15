package main

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/spf13/cobra"
)

const usbIDsURL = "http://www.linux-usb.org/usb.ids"
const usbIDsFile = "usb.ids"

type USBDevice struct {
	Vendor  string `yaml:"vendor"`
	Product string `yaml:"product"`
	VID     int64  `yaml:"vid"`
	PID     int64  `yaml:"pid"`
}

var rootCmd = &cobra.Command{
	Use:   "usb-tools",
	Short: "USB device database parser and filter",
	Long:  `A CLI tool to parse USB device IDs and filter specific device types like mice.`,
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

// downloadUSBIDs downloads the usb.ids file if it doesn't exist
func downloadUSBIDs() error {
	if _, err := os.Stat(usbIDsFile); err == nil {
		fmt.Printf("File %s already exists, skipping download\n", usbIDsFile)
		return nil
	}

	fmt.Printf("Downloading %s from %s...\n", usbIDsFile, usbIDsURL)

	resp, err := http.Get(usbIDsURL)
	if err != nil {
		return fmt.Errorf("failed to download: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", resp.Status)
	}

	out, err := os.Create(usbIDsFile)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	fmt.Printf("Successfully downloaded %s\n", usbIDsFile)
	return nil
}

// ensureUSBIDs checks if usb.ids exists and downloads if needed
func ensureUSBIDs() error {
	if _, err := os.Stat(usbIDsFile); os.IsNotExist(err) {
		return downloadUSBIDs()
	}
	return nil
}

func init() {
	rootCmd.AddCommand(parseCmd)
	rootCmd.AddCommand(filterCmd)
	rootCmd.AddCommand(downloadCmd)
}
