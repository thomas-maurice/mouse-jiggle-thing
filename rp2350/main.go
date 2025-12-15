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

	pio "github.com/tinygo-org/pio/rp2-pio"
	"github.com/tinygo-org/pio/rp2-pio/piolib"
)

const (
	// Probs should not be too high if you value your eyesight
	// (don't ask me how I found out)
	MAX_BRIGHTNESS = 10
	MIN_BRIGHTNESS = 0

	// 5000 ms of jitter
	MAX_MS_DELAY_JITTER = 5000
	// Base wait time between 2 movements
	BASE_WAIT_TIME_MS = 500

	// Jitter for the increment
	BRIGHTNESS_INCREMENT_JITTER = 3
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

func clamp(val int) uint8 {
	if val >= MAX_BRIGHTNESS {
		return MAX_BRIGHTNESS
	}

	if val <= MIN_BRIGHTNESS {
		return MIN_BRIGHTNESS
	}

	// maybe not required
	return uint8(val & 0xff)
}

func breathe() {
	bri := [3]Brightness{
		{
			Value:     0,
			Increment: 1,
		},
		{
			Value:     MAX_BRIGHTNESS / 3 * 1,
			Increment: 1,
		},
		{
			Value:     MAX_BRIGHTNESS / 3 * 2,
			Increment: 1,
		},
	}

	sm, _ := pio.PIO0.ClaimStateMachine()
	ws, err := piolib.NewWS2812B(sm, machine.GP22)
	if err != nil {
		panic(err.Error())
	}

	for {
		ws.PutRGB(
			clamp(bri[0].Value),
			clamp(bri[1].Value),
			clamp(bri[2].Value),
		)
		time.Sleep(time.Millisecond * 150)

		for idx, val := range bri {
			brightness := val.Value
			if brightness >= MAX_BRIGHTNESS {
				// randomises a bit the pixel increment so we have
				// new colours
				randomIncrement, err := machine.GetRNG()
				if err != nil {
					panic(err)
				}
				bri[idx].Increment = -1 * int(randomIncrement%BRIGHTNESS_INCREMENT_JITTER)
			} else if brightness <= MIN_BRIGHTNESS {
				// randomises a bit the pixel increment so we have
				// new colours
				randomIncrement, err := machine.GetRNG()
				if err != nil {
					panic(err)
				}
				bri[idx].Increment = int(randomIncrement % BRIGHTNESS_INCREMENT_JITTER)
			}
			bri[idx].Value = bri[idx].Value + bri[idx].Increment
		}
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
