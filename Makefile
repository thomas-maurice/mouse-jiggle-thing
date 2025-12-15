.PHONY: flash build usb-tools clean

# Environment variables with defaults
USB_MANUFACTURER ?= Razer USA, Ltd
USB_PRODUCT ?= Razer DeathAdder V2 Pro
USB_VID ?= 5426
USB_PID ?= 124
OUTPUT_FILE ?= firmware.uf2
TARGET ?= rp2350

# Target-specific settings
ifeq ($(TARGET),rp2040)
	TINYGO_TARGET := pico
	SOURCE_DIR := rp2040
else
	TINYGO_TARGET := pico2
	SOURCE_DIR := rp2350
endif

flash:
	cd $(SOURCE_DIR) && tinygo flash -target $(TINYGO_TARGET) \
		-ldflags="-X main.usbManufacturer='$(USB_MANUFACTURER)' -X main.usbProduct='$(USB_PRODUCT)' -X main.usbVID='$(USB_VID)' -X main.usbPID='$(USB_PID)'"

build:
	@cd $(SOURCE_DIR) && tinygo build -target $(TINYGO_TARGET) -o ../$(OUTPUT_FILE) \
		-ldflags="-X 'main.usbManufacturer=$(USB_MANUFACTURER)' -X 'main.usbProduct=$(USB_PRODUCT)' -X 'main.usbVID=$(USB_VID)' -X 'main.usbPID=$(USB_PID)'"

GO ?= go

usb-tools:
	mkdir -p bin
	cd scripts && $(GO) build -o ../bin/usb-tools .
	@echo "usb-tools built to bin/usb-tools"

clean:
	rm -rf bin/
	rm -rf output/
	rm -f firmware.uf2
	rm -f firmware.u2f
	rm -f scripts/usb-tools
	@echo "Cleaned build artifacts"