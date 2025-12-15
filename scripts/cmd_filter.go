package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var filterCmd = &cobra.Command{
	Use:   "filter",
	Short: "Filter USB devices by type",
	Long:  `Filters USB devices from the parsed YAML file by device type (e.g., mice).`,
	RunE: func(cmd *cobra.Command, args []string) error {
		inputFile, _ := cmd.Flags().GetString("input")
		outputFile, _ := cmd.Flags().GetString("output")
		deviceType, _ := cmd.Flags().GetString("type")

		return filterDevices(inputFile, outputFile, deviceType)
	},
}

func init() {
	filterCmd.Flags().StringP("input", "i", "usb_devices.yaml", "Input YAML file")
	filterCmd.Flags().StringP("output", "o", "usb_filtered.yaml", "Output YAML file")
	filterCmd.Flags().StringP("type", "t", "mice", "Device type to filter (mice, keyboard, etc.)")
}

func filterDevices(inputPath, outputPath, deviceType string) error {
	// Read the full YAML file
	data, err := os.ReadFile(inputPath)
	if err != nil {
		return fmt.Errorf("error reading file: %w", err)
	}

	var devices []USBDevice
	err = yaml.Unmarshal(data, &devices)
	if err != nil {
		return fmt.Errorf("error parsing YAML: %w", err)
	}

	// Get keywords based on device type
	keywords := getKeywordsForType(deviceType)

	// Filter devices
	var filtered []USBDevice
	for _, device := range devices {
		productLower := strings.ToLower(device.Product)

		// Check if product mentions any keyword
		match := false
		for _, keyword := range keywords {
			if strings.Contains(productLower, keyword) {
				match = true
				break
			}
		}

		if match {
			filtered = append(filtered, device)
		}
	}

	// Generate YAML output
	yamlData, err := yaml.Marshal(filtered)
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

	fmt.Printf("Successfully filtered %d %s from %d total USB devices\n", len(filtered), deviceType, len(devices))
	fmt.Printf("Output written to %s\n", outputPath)
	return nil
}

func getKeywordsForType(deviceType string) []string {
	switch strings.ToLower(deviceType) {
	case "mice", "mouse":
		return []string{
			"mouse", "mice", "trackball", "pointer",
			"mx master", "mx anywhere", "mx ergo", "mx vertical",
			"deathadder", "viper", "basilisk", "naga", "mamba", "orochi",
			"g502", "g305", "g703", "g pro", "g403", "g903", "g603", "g604",
			"rival", "sensei", "aerox", "prime",
			"model o", "model d",
			"intellimouse", "arc mouse", "surface mouse",
			"gaming mouse", "optical mouse", "wireless mouse", "laser mouse",
			"ergonomic mouse", "vertical mouse",
		}
	case "keyboard", "keyboards":
		return []string{
			"keyboard", "keypad", "numpad",
			"mechanical keyboard", "gaming keyboard",
		}
	case "headset", "headsets", "audio":
		return []string{
			"headset", "headphone", "headphones", "audio",
			"speaker", "speakers", "microphone", "mic",
		}
	case "webcam", "camera":
		return []string{
			"webcam", "camera", "video",
		}
	default:
		return []string{strings.ToLower(deviceType)}
	}
}
