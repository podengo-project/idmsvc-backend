#!/bin/bash

##
# This script merge the right format property to the openapi specification
# so that, we can keep using the api-designer tool to edit the files,
# and add later the format property using this script.
##
INPUT_PATH="$1"
PATCH_PATH="$2"
OUTPUT_PATH="$3"

./bin/yq eval-all 'select(fileIndex == 0) * select(fileIndex == 1)' "${INPUT_PATH}" "${PATCH_PATH}" > "${OUTPUT_PATH}"
