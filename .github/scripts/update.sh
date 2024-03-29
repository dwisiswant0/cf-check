#!/bin/bash

function getCfPrefs() {
  curl -s "https://api.bgpview.io/asn/${1}/prefixes" | \
    jq -r '.data.ipv4_prefixes[].parent.prefix' | \
    grep / | sort -u
}

DB="$(echo -e "$(getCfPrefs "13335")\n$(getCfPrefs "395747")" | sort -u)"

if [[ "$(echo "${DB}" | wc -l)" != "$(curl -skL "${1}" | wc -l)" ]]; then
  echo "${DB}" > db/prefixes.txt
fi
