#!/bin/bash
set -e

# Colors
YELLOW='\033[1;33m'
GREEN='\033[0;32m'
BLUE='\033[0;34m'
RED='\033[0;31m'
NC='\033[0m'

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"
MICE_YAML="${SCRIPT_DIR}/usb_mice.yaml"
USB_TOOLS="${PROJECT_ROOT}/bin/usb-tools"

# Find the correct go binary (not tinygo's wrapper)
find_go() {
    # Check common locations
    for go_path in /usr/bin/go /usr/local/go/bin/go; do
        if [ -f "$go_path" ]; then
            echo "$go_path"
            return
        fi
    done
    # Fallback to PATH
    echo "go"
}

# Check if usb-tools exists
if [ ! -f "$USB_TOOLS" ]; then
    echo -e "${YELLOW}usb-tools not found. Building...${NC}"
    cd "$PROJECT_ROOT"
    GO=$(find_go) make usb-tools
fi

# Check if usb_mice.yaml exists
if [ ! -f "$MICE_YAML" ]; then
    echo -e "${YELLOW}usb_mice.yaml not found. Generating...${NC}"
    "$USB_TOOLS" parse
    "$USB_TOOLS" filter --type mice --output "$MICE_YAML"
fi

# Determine target (rp2040, rp2350, or both)
TARGET="${TARGET:-both}"

# Create output directories
if [ "$TARGET" = "both" ] || [ "$TARGET" = "rp2040" ]; then
    mkdir -p "$PROJECT_ROOT/output/rp2040"
fi
if [ "$TARGET" = "both" ] || [ "$TARGET" = "rp2350" ]; then
    mkdir -p "$PROJECT_ROOT/output/rp2350"
fi

echo -e "${BLUE}=== Firmware Generation ===${NC}"
echo -e "Target: $TARGET"
echo -e "Parsing mice from $MICE_YAML\n"

cd "$PROJECT_ROOT"

# Function to build a single firmware
build_firmware() {
    vendor="$1"
    product="$2"
    vid="$3"
    pid="$4"
    target="$5"

    # Create safe filename with target prefix
    safe_name=$(echo "${vendor}_${product}" | sed 's/[^a-zA-Z0-9_-]/_/g' | sed 's/__*/_/g' | cut -c1-100)
    # Convert decimal VID/PID to hexadecimal (without 0x prefix)
    vid_hex=$(printf "%04x" "$vid")
    pid_hex=$(printf "%04x" "$pid")
    output_file="output/${target}/${target}-${safe_name}_${vid_hex}.${pid_hex}.uf2"
    
    # Skip if exists
    if [ -f "$output_file" ]; then
        echo -e "${YELLOW}SKIP${NC}: $vendor - $product"
        return 0
    fi
    
    echo -e "${BLUE}BUILD${NC}: $vendor - $product (VID=$vid PID=$pid)"
    
    # Build firmware
    USB_MANUFACTURER="$vendor" \
    USB_PRODUCT="$product" \
    USB_VID="$vid" \
    USB_PID="$pid" \
    OUTPUT_FILE="$output_file" \
    TARGET="$target" \
    make build 2>&1 | grep -v "^tinygo" || true
    
    if [ -f "$output_file" ]; then
        echo -e "${GREEN}✓ SUCCESS${NC}: $output_file"
    else
        echo -e "${RED}✗ FAILED${NC}: $vendor - $product"
        return 1
    fi
}

export -f build_firmware
export PROJECT_ROOT YELLOW GREEN BLUE RED NC

# Get total number of devices
TOTAL_DEVICES=$(yq eval 'length' "$MICE_YAML")

# Determine which targets to build
TARGETS=()
if [ "$TARGET" = "both" ]; then
    TARGETS=("rp2040" "rp2350")
elif [ "$TARGET" = "rp2040" ] || [ "$TARGET" = "rp2350" ]; then
    TARGETS=("$TARGET")
fi

# Build for each target
for current_target in "${TARGETS[@]}"; do
    echo -e "\n${BLUE}=== Building for $current_target ===${NC}\n"

    # Apply sharding if SHARD_INDEX and SHARD_TOTAL are set
    if [ -n "$SHARD_INDEX" ] && [ -n "$SHARD_TOTAL" ]; then
        echo "Building shard $SHARD_INDEX of $SHARD_TOTAL (total devices: $TOTAL_DEVICES)"

        # Calculate which devices this shard should build
        # Use modulo to distribute evenly
        if command -v parallel &> /dev/null; then
            echo "Using GNU parallel for faster builds..."
            yq eval '.[] | [.vendor, .product, .vid, .pid] | @tsv' "$MICE_YAML" | \
                awk -v shard="$SHARD_INDEX" -v total="$SHARD_TOTAL" 'NR % total == shard' | \
                parallel --colsep '\t' --jobs 8 --line-buffer build_firmware {1} {2} {3} {4} "$current_target"
        else
            echo "GNU parallel not found, building sequentially..."
            yq eval '.[] | [.vendor, .product, .vid, .pid] | @tsv' "$MICE_YAML" | \
                awk -v shard="$SHARD_INDEX" -v total="$SHARD_TOTAL" 'NR % total == shard' | \
                while IFS=$'\t' read -r vendor product vid pid; do
                    build_firmware "$vendor" "$product" "$vid" "$pid" "$current_target"
                done
        fi
    else
        # No sharding, build all
        if command -v parallel &> /dev/null; then
            echo "Using GNU parallel for faster builds..."
            yq eval '.[] | [.vendor, .product, .vid, .pid] | @tsv' "$MICE_YAML" | \
                parallel --colsep '\t' --jobs 8 --line-buffer build_firmware {1} {2} {3} {4} "$current_target"
        else
            echo "GNU parallel not found, building sequentially..."
            yq eval '.[] | [.vendor, .product, .vid, .pid] | @tsv' "$MICE_YAML" | while IFS=$'\t' read -r vendor product vid pid; do
                build_firmware "$vendor" "$product" "$vid" "$pid" "$current_target"
            done
        fi
    fi
done

echo -e "\n${GREEN}Done!${NC}"
RP2040_COUNT=$(ls output/rp2040/*.uf2 2>/dev/null | wc -l)
RP2350_COUNT=$(ls output/rp2350/*.uf2 2>/dev/null | wc -l)
TOTAL_COUNT=$((RP2040_COUNT + RP2350_COUNT))
echo "RP2040 firmwares: $RP2040_COUNT"
echo "RP2350 firmwares: $RP2350_COUNT"
echo "Total firmwares: $TOTAL_COUNT"
