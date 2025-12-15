package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var parseCmd = &cobra.Command{
	Use:   "parse",
	Short: "Parse usb.ids file and generate YAML",
	Long:  `Parses the usb.ids file and generates a YAML file with all USB devices, converting hex IDs to decimal.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Ensure usb.ids file exists
		if err := ensureUSBIDs(); err != nil {
			return err
		}

		outputFile, _ := cmd.Flags().GetString("output")
		return parseUSBIDs(outputFile)
	},
}

func init() {
	parseCmd.Flags().StringP("output", "o", "usb_devices.yaml", "Output YAML file")
}

func parseUSBIDs(outputPath string) error {
	file, err := os.Open(usbIDsFile)
	if err != nil {
		return fmt.Errorf("error opening file: %w", err)
	}
	defer file.Close()

	var devices []USBDevice
	var currentVendor string
	var currentVID int64

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		// Skip comments and empty lines
		if strings.HasPrefix(line, "#") || len(strings.TrimSpace(line)) == 0 {
			continue
		}

		// Check if it's a vendor line (starts at column 0 with hex ID)
		if len(line) > 0 && line[0] != '\t' && line[0] != ' ' {
			// Parse vendor line: "XXXX  Vendor Name"
			parts := strings.SplitN(line, "  ", 2)
			if len(parts) == 2 {
				vendorID := strings.TrimSpace(parts[0])
				vendorName := strings.TrimSpace(parts[1])

				// Convert hex vendor ID to decimal
				vid, err := strconv.ParseInt(vendorID, 16, 64)
				if err != nil {
					continue // Skip if not a valid hex number
				}

				currentVendor = vendorName
				currentVID = vid
			}
		} else if strings.HasPrefix(line, "\t") && !strings.HasPrefix(line, "\t\t") {
			// This is a product line (single tab indentation)
			// Remove the leading tab
			productLine := strings.TrimPrefix(line, "\t")

			// Parse product line: "XXXX  Product Name"
			parts := strings.SplitN(productLine, "  ", 2)
			if len(parts) == 2 && currentVendor != "" {
				productID := strings.TrimSpace(parts[0])
				productName := strings.TrimSpace(parts[1])

				// Convert hex product ID to decimal
				pid, err := strconv.ParseInt(productID, 16, 64)
				if err != nil {
					continue // Skip if not a valid hex number
				}

				// Remove all quotes to avoid shell escaping issues
				vendorSafe := strings.ReplaceAll(strings.ReplaceAll(currentVendor, `"`, ``), `'`, ``)
				productSafe := strings.ReplaceAll(strings.ReplaceAll(productName, `"`, ``), `'`, ``)

				device := USBDevice{
					Vendor:  vendorSafe,
					Product: productSafe,
					VID:     currentVID,
					PID:     pid,
				}
				devices = append(devices, device)
			}
		}
		// Skip interface lines (double tab indentation)
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading file: %w", err)
	}

	// Generate YAML output
	yamlData, err := yaml.Marshal(devices)
	if err != nil {
		return fmt.Errorf("error marshaling to YAML: %w", err)
	}

	// Write to output file
	out, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("error creating output file: %w", err)
	}
	defer out.Close()

	_, err = out.Write(yamlData)
	if err != nil {
		return fmt.Errorf("error writing to output file: %w", err)
	}

	fmt.Printf("Successfully parsed %d USB devices to %s\n", len(devices), outputPath)
	return nil
}
