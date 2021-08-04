#!/usr/bin/env bash
# ^^^^^^^^^^^^^^^^^ this is the most platform-agnostic way to guarantee this script runs with Bash
# 2021-07-08 WATERMARK, DO NOT REMOVE - This script was generated from the Kurtosis Bash script template

set -euo pipefail   # Bash "strict mode"
script_dirpath="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
root_dirpath="$(dirname "${script_dirpath}")"
kurtosis_core_dirpath="${root_dirpath}/.kurtosis"


# ==================================================================================================
#                                       Arg Parsing & Validation
# ==================================================================================================
show_help_and_exit() {
    echo ""
    echo "Usage: $(basename "${0}") action [kurtosis.sh_arg1] [kurtosis.sh_arg2]..."
    echo ""
    echo "  action              The action that should be passed to the underlying build-and-run-core.sh script to tell it which action should be taken (call"
    echo "                          'bash ${kurtosis_core_dirpath}/build-and-run-core.sh help' directly for all available actions)"
    echo "  kurtosis.sh_args    Optional, supplemental args that should be passed to the kurtosis.sh script to modify testsuite execution behaviour (call"
    echo "                          'bash ${kurtosis_core_dirpath}/kurtosis.sh --help' directly for all available args)"
    echo ""
    exit 1  # Exit with error so CI will fail if it accidentally calls this
}

if [ "${#}" -eq 0 ]; then
    show_help_and_exit
fi
action="${1:-}"
shift 1
if [ "${action}" == "help" ]; then
    show_help_and_exit
fi


# ==================================================================================================
#                                             Main Logic
# ==================================================================================================
custom_params_json='{}'
# >>>>>>>> Add custom testsuite parameters here <<<<<<<<<<<<<

bash "${kurtosis_core_dirpath}/build-and-run-core.sh" \
    "${action}" \
    "kurtosis-onboarding-test" \
    "${root_dirpath}" \
    "${root_dirpath}/testsuite/Dockerfile" \
    "${kurtosis_core_dirpath}/kurtosis.sh" \
    --custom-params "${custom_params_json}" \
    ${1+"${@}"}
