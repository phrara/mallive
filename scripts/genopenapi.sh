#!/usr/bin/env bash

set -euo pipefail

shopt -s globstar

if ! [[ "$0" =~ scripts/genopenapi.sh ]]; then
    echo "must be run from repo root"
    exit 255
fi


OPENAPI_ROOT="./api/openapi"

GEN_SERVER=(
    # "chi-server"
    # "echo-server"
    # "fiber-server"
    "gin-server"
)

if [ ${#GEN_SERVER[@]} != 1 ]; then
    echo "Multiple server frameworks have been selected!"
    exit 255
fi

echo "Server framework <${GEN_SERVER[0]}> has been selected."

function get_openapi_files() {
    openapi_files=$(find -type f -path "${OPENAPI_ROOT}" -name "*.yml" -o -name "*.yaml")
    echo "${openapi_files[@]}"
}

get_openapi_files

function prepare_dir() {
    local dir="$1"
    if [ -d "$dir" ]; then
        echo "found exist ${dir}, clean"
        find "${dir}" -mindepth 1 -delete
    else
        mkdir -p "${dir}"
    fi
}

function gen_openapi() {
    local output_dir=$1
    local package=$2
    local service=$3

    mkdir -p "${output_dir}"
    find "${output_dir}" -type f -name "*.gen.go" -delete

    client_path="internal/common/client/${service}"
    prepare_dir ${client_path}

    # Server
    oapi-codegen -generate types -o "${output_dir}/openapi_types.gen.go" -package "${package}" "${OPENAPI_ROOT}/${service}.yml"
    oapi-codegen -generate "${GEN_SERVER}" -o "${output_dir}/openapi_api.gen.go" -package "${package}" "${OPENAPI_ROOT}/${service}.yml"

    # Client
    oapi-codegen -generate types -o "${client_path}/openapi_types.gen.go" -package "${service}" "${OPENAPI_ROOT}/${service}.yml"
    oapi-codegen -generate client -o "${client_path}/openapi_client.gen.go" -package "${service}" "${OPENAPI_ROOT}/${service}.yml"
    
    echo "genopenapi successfully"
}


gen_openapi internal/order/ports ports order
