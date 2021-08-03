set -euo pipefail
script_dirpath="$(cd "$(dirname "${0}")" && pwd)"
root_dirpath="$(dirname "${script_dirpath}")"

# ==========================================================================================
#                                         Constants
# ==========================================================================================
KURTOSIS_DOCKERHUB_ORG="kurtosistech"
LANG_SCRIPTS_DIRNAME="scripts"
BUILD_AND_RUN_FILENAME="build-and-run.sh"

ERROR_LOG_KEYWORD="ERRO"

SUPPORTED_LANGS_FILENAME="supported-languages.txt"

# ==========================================================================================
#                                        Arg-parsing
# ==========================================================================================
docker_username="${1:-}"
docker_password_DO_NOT_LOG="${2:-}" # WARNING: DO NOT EVER LOG THIS!!
kurtosis_client_id="${3:-}"
kurtosis_client_secret_DO_NOT_LOG="${4:-}" # WARNING: DO NOT EVER LOG THIS!!

# ==========================================================================================
#                                        Arg validation
# ==========================================================================================
if [ -z "${docker_username}" ]; then
    echo "Error: Docker username cannot be empty" >&2
    exit 1
fi
if [ -z "${docker_password_DO_NOT_LOG}" ]; then
    echo "Error: Docker password cannot be empty" >&2
    exit 1
fi
if [ -z "${kurtosis_client_id}" ]; then
    echo "Error: Kurtosis client ID cannot be empty" >&2
    exit 1
fi
if [ -z "${kurtosis_client_secret_DO_NOT_LOG}" ]; then
    echo "Error: Kurtosis client secret cannot be empty" >&2
    exit 1
fi

# ==========================================================================================
#                                           Main code
# ==========================================================================================
# Docker is restricting anonymous image pulls, so we log in before we do any pulling
if ! docker login -u "${docker_username}" -p "${docker_password_DO_NOT_LOG}"; then
    echo "Error: Logging in to Docker failed" >&2
    exit 1
fi

# Building and running testsuites take a very long time, so we do some optimizations:
# 1) skip building/running testsuites if only docs changes
if git --no-pager diff --exit-code origin/develop...HEAD -- . ':!*.md' > /dev/null; then
    echo "Skipping building and running testsuites as the only changes are in Markdown files"
    exit 0
fi
# 2) if there are changes in the code shared across all langs, we always need to build all testsuites
supported_langs_filepath="${root_dirpath}/${SUPPORTED_LANGS_FILENAME}"
not_lang_dirs_filters=""
for lang in $(cat "${supported_langs_filepath}"); do
    not_lang_dirs_filters="${not_lang_dirs_filters} :!${lang}"
done
if git --no-pager diff --exit-code origin/develop...HEAD -- . ':!*.md' ${not_lang_dirs_filters} > /dev/null; then
    has_shared_code_changes="false"
else
    has_shared_code_changes="true"
fi
# 3) if no shared code changes, then we only need to build the testsuites that had changes
lang_dirs_needing_building=()
for lang in $(cat "${supported_langs_filepath}"); do
    if ! "${has_shared_code_changes}" && git --no-pager diff --exit-code origin/develop...HEAD -- "${lang}" > /dev/null; then
        echo "Skipping adding ${lang} directory to list of testsuites to build as there are no shared code changes and the directory doesn't have any changes"
        continue
    fi
    lang_dirs_needing_building+=("${lang}")
done

echo "Building and running all example testsuites in need of validation..."
for lang in "${lang_dirs_needing_building[@]}"; do
    echo "Building and running ${lang} testsuite..."
    buildscript_filepath="${root_dirpath}/${lang}/${LANG_SCRIPTS_DIRNAME}/${BUILD_AND_RUN_FILENAME}"
    output_filepath="$(mktemp)"
    if ! bash "${buildscript_filepath}" all --client-id "${kurtosis_client_id}" --client-secret "${kurtosis_client_secret_DO_NOT_LOG}" 2>&1 | tee "${output_filepath}"; then
        echo "Error: Building and running ${lang} testsuite failed" >&2
        exit 1
    fi
    echo "Successfully built and run ${lang} testsuite"

    # This helps us catch errors that might show up in the testsuite logs but not get propagated to the actual exit codes
    echo "Scanning output file for error log keyword '${ERROR_LOG_KEYWORD}'..."
    if grep "${ERROR_LOG_KEYWORD}" "${output_filepath}"; then
        echo "Error: Detected error log pattern '${ERROR_LOG_KEYWORD}' in output file" >&2
        exit 1
    fi
    echo "No instances of error log keyword found"
done
echo "Successfully built and ran all example testsuites in need of validation"
