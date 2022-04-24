#!/usr/bin/env bash

# Download 国土数値情報/医療機関データ
# https://nlftp.mlit.go.jp/ksj/gml/datalist/KsjTmplt-P04-v3_0.html

set -euo pipefail
IFS=$'\n\t'

print_help() {
  cat <<EOD
Usage:
  ${_ME} [<arguments>]
  ${_ME} -h | --help

Options:
  -h --help  Show this screen.
EOD
}

download() {
  pushd data
  curl -O -s 'https://nlftp.mlit.go.jp/ksj/gml/data/P04/P04-20/P04-20_GML.zip'
  unzip -j P04-20_GML.zip '*.geojson'
  popd
}

main() {
  if [[ "${1:-}" =~ ^-h|--help$  ]]
  then
    print_help
  else
    download "$@"
  fi
}

main "$@"
