# USB Tools CLI

A command-line tool to parse USB device IDs from the official USB ID database and filter devices by type (e.g., mice, keyboards, etc.).

## Prerequisites

- Standard Go (not TinyGo) must be installed
- Dependencies (`cobra`, `yaml.v3`) will be automatically downloaded via `go get`

## Building

```bash
cd scripts
go build -o usb-tools .
```

Or use the full Go path:
```bash
/Users/thomas/.nix-profile/bin/go build -o usb-tools .
```

## Commands

### `usb-tools download`

Downloads the latest `usb.ids` file from http://www.linux-usb.org/usb.ids

```bash
# Download if not present
./usb-tools download

# Force re-download even if file exists
./usb-tools download --force
```

### `usb-tools parse`

Parses the `usb.ids` file and generates a YAML file with all USB devices. The file will be automatically downloaded if not present.

```bash
# Parse and create usb_devices.yaml (default)
./usb-tools parse

# Parse and specify custom output file
./usb-tools parse --output custom_output.yaml
```

**Features:**
- Automatically downloads `usb.ids` if not present
- Converts hex vendor/product IDs to decimal
- Outputs structured YAML with 21,813+ USB devices

### `usb-tools filter`

Filters USB devices by type from the parsed YAML file.

```bash
# Filter mice (default)
./usb-tools filter

# Filter with custom options
./usb-tools filter --type mice --input usb_devices.yaml --output usb_mice.yaml

# Filter keyboards
./usb-tools filter --type keyboard --output usb_keyboards.yaml

# Filter headsets
./usb-tools filter --type headset --output usb_headsets.yaml
```

**Supported device types:**
- `mice` / `mouse` - Mouse devices (803 devices found)
- `keyboard` / `keyboards` - Keyboard devices
- `headset` / `headsets` / `audio` - Audio devices
- `webcam` / `camera` - Camera devices
- Any custom keyword

## Output Format

All YAML files follow this structure:

```yaml
- vendor: Vendor Name
  product: Product Name
  vid: 1234      # Vendor ID in decimal (converted from hex)
  pid: 5678      # Product ID in decimal (converted from hex)
```

## Example

From the usb.ids file:
```
0001  Fry's Electronics
	7778  Counterfeit flash drive [Kingston]
```

Converts to YAML:
```yaml
- vendor: Fry's Electronics
  product: Counterfeit flash drive [Kingston]
  vid: 1
  pid: 30584
```

Where `0001` (hex) = `1` (decimal) and `7778` (hex) = `30584` (decimal).

## Common Workflow

```bash
# Build the tool
go build -o usb-tools .

# Parse USB IDs (auto-downloads if needed)
./usb-tools parse

# Filter to get only mice
./usb-tools filter --type mice --output usb_mice.yaml

# Update the database
./usb-tools download --force
./usb-tools parse
```

## Files Generated

- `usb.ids` - Raw USB ID database (730KB, downloaded from linux-usb.org)
- `usb_devices.yaml` - All parsed USB devices (21,813 entries)
- `usb_mice.yaml` - Filtered mouse devices (803 entries)
- `usb-tools` - The compiled CLI binary

## Top Mouse Vendors

The filter captures mice from major brands including:
- **Logitech** (125 mice) - G502, G Pro, MX Master, MX Anywhere, etc.
- **Razer** (72 mice) - DeathAdder, Viper, Basilisk, Naga, Mamba, etc.
- **Elecom** (74 mice)
- **Kensington** (51 mice)
- **Microsoft** (36 mice) - IntelliMouse, Arc Mouse, Surface Mouse, etc.
- **ROCCAT** (23 mice)
- Plus many others (HP, Dell, ASUS, Corsair, SteelSeries, etc.)
