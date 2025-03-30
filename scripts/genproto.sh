#!/usr/bin/env bash

set -euo pipefail

shopt -s globstar

if ! [[ "$0" =~ scripts/genproto.sh ]]; then
    echo "must be run from repo root"
    exit 255
fi


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
        rm -rf $go_out
    fi

    for d in $dir_list; do
        local service="${d:0:${#d}-2}"
        # echo $service
        local pb_file="${service}.proto"
        if [ -d "${go_out}/${d}" ]; then
            echo "found exist ${go_out}/${d}, after cleaning, run again"
            rm -rf "${go_out}/${d}"
        else
            mkdir -p "${go_out}/${d}"
            # paths=source_relative：输出文件与输入文件放在相同的相对目录中
            protoc \
                -I="/home/phr/protobuf/include" \
                -I="${API_ROOT}" \
                "--go_out=${go_out}" --go_opt=paths=source_relative \
                --go-grpc_opt=require_unimplemented_servers=false \
                "--go-grpc_out=${go_out}" --go-grpc_opt=paths=source_relative \
                "${API_ROOT}/${d}/${pb_file}"
        fi
    done 
    echo "genproto successfully"
}

gen_for_modules
