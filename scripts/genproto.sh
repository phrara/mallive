#!/usr/bin/env bash

set -euo pipefail

shopt -s globstar

if ! [[ "$0" =~ scripts/genproto.sh ]]; then
    echo "must be run from repo root"
    exit 255
fi

# source ./scripts/lib.sh

API_ROOT="./api"

function get_dirs() {
    dir_list=()
    while IFS= read -r dir; do
        dir_list+=("$dir")
    done < <(find . -type f -name "*.proto" -exec dirname {} \; | xargs -n1 basename | sort -u)
    echo "${dir_list[@]}"
}

get_dirs

function get_pb_files() {
    pb_files=$(find . -type f -name "*.proto")
    echo "${pb_files[@]}"
}

get_pb_files

function gen_for_modules() {
    local go_out="internal/common/genproto"
    if [ -d "$go_out" ]; then
        echo "found exist ${go_out}, clean"
        run rm -rf $go_out
    fi

    for d in $dir_list; do
        local service="${d:0:${#d}-2}"
        # echo $service
        local pb_file="${service}.proto"
        

    done
}

gen_for_modules

echo ""