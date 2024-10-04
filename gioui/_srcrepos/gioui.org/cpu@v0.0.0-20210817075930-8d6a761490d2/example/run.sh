#!/bin/sh

# SPDX-License-Identifier: Unlicense OR MIT

set -e

SWIFTSHADER=$HOME/.cache/swiftshader/build.64bit/Linux/vk_swiftshader_icd.json
SWIFTSHADER_TRIPLE=x86_64-linux-unknown VK_ICD_FILENAMES=$SWIFTSHADER go run ../cmd/compile -layout "0:buffer,1:image" -arch amd64 example.comp
# Build and run driver
CGO_ENABLED=1 go run ../cmd/example
