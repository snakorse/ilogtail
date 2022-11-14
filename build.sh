#!/bin/bash
set -xue
set -o pipefail

function ramAvail () {
  local ramavail=$(cat /proc/meminfo | grep -i 'MemAvailable' | grep -o '[[:digit:]]*')
  echo $ramavail
}

nproc=$(nproc)
ram_size=$(ramAvail)
ram_limit_nproc=$((ram_size / 1024 / 768))
[[ $ram_limit_nproc -ge $nproc ]] || nproc=$ram_limit_nproc
[[ $nproc -gt 0 ]] || nproc=1

mkdir -p core/build && cd core/build && cmake -DCMAKE_BUILD_TYPE=Release -DLOGTAIL_VERSION=byted -DBUILD_LOGTAIL_UT=OFF -DENABLE_COMPATIBLE_MODE=OFF -DENABLE_STATIC_LINK_CRT=OFF .. && make -sj$nproc && cd - && ./scripts/upgrade_adapter_lib.sh && ./scripts/plugin_build.sh mod c-shared output
cp core/build/ilogtail output
cp core/build/plugin/libPluginAdapter.so output
echo -e "{\n}" > output/ilogtail_config.json
mkdir -p output/user_yaml_config.d