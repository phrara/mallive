#!/usr/bin/env bash

set -euo pipefail

shopt -s globstar

if ! [["$0" =~ scripts/genproto.sh]]; then
    