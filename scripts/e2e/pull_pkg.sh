#!/usr/bin/env bash

# set the kpm default registry and repository
export KPM_REG="ghcr.io"
export KPM_REPO="kcl-lang"
export OCI_REG_PLAIN_HTTP=off

current_dir=$(pwd)

mkdir -p ./scripts/e2e/pkg_in_reg/

cd ./scripts/e2e/pkg_in_reg/


# Check if file exists
if [ ! -d "./ghcr.io/kcl-lang/k8s/1.28" ]; then
  $current_dir/bin/kcl mod pull k8s:1.28
fi

if [ ! -d "./ghcr.io/kcl-lang/helloworld" ]; then
  $current_dir/bin/kcl mod pull helloworld
fi

cd "$current_dir"
