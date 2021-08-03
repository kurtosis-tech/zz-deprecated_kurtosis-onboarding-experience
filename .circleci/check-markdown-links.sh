#!/usr/bin/env bash
# ^^^^^^^^^^^^^^^^^ this is the most platform-agnostic way to guarantee this script runs with Bash
# 2021-07-08 WATERMARK, DO NOT REMOVE - This script was generated from the Kurtosis Bash script template

set -euo pipefail   # Bash "strict mode"
script_dirpath="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
root_dirpath="$(dirname "${script_dirpath}")"



# ==================================================================================================
#                                             Main Logic
# ==================================================================================================
# We can use relative syntax to specify kurtosis-core/kurtosis-libs links in the Markdown (e.g. "./kurtosis-core/architecture")
# because everything gets published to the same docs.kurtosistech.com domain, but we need to expand to the full URL for the Markdown
# link checker
config_filepath="$(mktemp)"
cat << EOF > "${config_filepath}"
{
    "replacementPatterns": [
        {
            "pattern": "^../kurtosis-core",
            "replacement": "https://docs.kurtosistech.com/kurtosis-core"
        },
        {
            "pattern": "^../kurtosis-libs",
            "replacement": "https://docs.kurtosistech.com/kurtosis-libs"
        }
    ]
}
EOF

# Inspired by https://github.com/open-telemetry/opentelemetry-collector/pull/1156/files/2244e61f4dd0378deffc00d939edf6f800687dcf
exit_code=0
for filepath in $(find "${root_dirpath}" -iname '*.md' | sort); do
    markdown-link-check --config "${config_filepath}" -qv "${filepath}" || exit_code=1
    # Wait to scan files so that we don't overload github with requests which may result in 429 responses
    sleep 2
done

exit "${exit_code}"
