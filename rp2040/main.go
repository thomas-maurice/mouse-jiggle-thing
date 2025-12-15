// Copyright (C) 2025 Thomas Maurice <thomas@maurice.fr>
// This work is free. You can redistribute it and/or modify it under the
// terms of the Do What The Fuck You Want To Public License, Version 2,
// as published by Sam Hocevar. See the LICENSE file for more details.

package main

import (
	"machine"
	"machine/usb"
	"machine/usb/hid/mouse"
	"math/rand"
	"strconv"
	"time"
)

const (
	// Probs should not be too high if you value your eyesight
	// (don't ask me how I found out)
	MAX_BRIGHTNESS = 20
	MIN_BRIGHTNESS = 0

	// 5000 ms of jitter
	MAX_MS_DELAY_JITTER = 5000
	// Base wait time between 2 movements
	BASE_WAIT_TIME_MS = 500

	// Jitter for the increment
	BRIGHTNESS_INCREMENT_JITTER = 4
)

var usbVID, usbPID string
var usbManufacturer, usbProduct string

type Brightness struct {
	Increment int
	Value     int
}

func init() {
	if usbVID != "" {
		vid, _ := strconv.ParseUint(usbVID, 0, 16)
		usb.VendorID = uint16(vid)
	}

	if usbPID != "" {
		pid, _ := strconv.ParseUint(usbPID, 0, 16)
		usb.ProductID = uint16(pid)
	}

	if usbManufacturer != "" {
		usb.Manufacturer = usbManufacturer
	}

	if usbProduct != "" {
		usb.Product = usbProduct
	}
}

func getJiggleDirection() int {
	if rand.Int()%2 == 0 {
		return 1
	}
	return -1
}

func breathe() {
	machine.LED.Configure(
		machine.PinConfig{
			Mode: machine.PinOutput,
		},
	)

	for {
		machine.LED.High()
		time.Sleep(time.Millisecond * 10)
		machine.LED.Low()
		time.Sleep(time.Millisecond * 10)
		machine.LED.High()
		time.Sleep(time.Millisecond * 10)
		machine.LED.Low()
		time.Sleep(time.Millisecond * 10)
		time.Sleep(time.Millisecond * 1000)
	}
}

func main() {
	go breathe()

	mouse := mouse.Port()

	for {
		randIncX, err := machine.GetRNG()
		if err != nil {
			panic(err)
		}
		randIncY, err := machine.GetRNG()
		if err != nil {
			panic(err)
		}

		incX := int(randIncX%10) * getJiggleDirection()
		incY := int(randIncY%10) * getJiggleDirection()
		mouse.Move(incX, incY)

		randomWait, err := machine.GetRNG()
		if err != nil {
			panic(err)
		}

		msDelay := time.Duration(randomWait%MAX_MS_DELAY_JITTER+BASE_WAIT_TIME_MS) * time.Millisecond
		time.Sleep(msDelay)
	}
}
