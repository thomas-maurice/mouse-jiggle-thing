// Copyright (C) 2025 Thomas Maurice <thomas@maurice.fr>
// This work is free. You can redistribute it and/or modify it under the
// terms of the Do What The Fuck You Want To Public License, Version 2,
// as published by Sam Hocevar. See the LICENSE file for more details.

package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var downloadCmd = &cobra.Command{
	Use:   "download",
	Short: "Download the usb.ids file",
	Long:  `Downloads the latest usb.ids file from http://www.linux-usb.org/usb.ids`,
	RunE: func(cmd *cobra.Command, args []string) error {
		force, _ := cmd.Flags().GetBool("force")

		if force {
			// Remove existing file to force re-download
			fmt.Println("Force download enabled, removing existing file...")
			if err := removeUSBIDs(); err != nil {
				fmt.Printf("Warning: could not remove existing file: %v\n", err)
			}
		}

		return downloadUSBIDs()
	},
}

func init() {
	downloadCmd.Flags().BoolP("force", "f", false, "Force download even if file exists")
}

func removeUSBIDs() error {
	return os.Remove(usbIDsFile)
}
