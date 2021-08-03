set -euo pipefail
script_dirpath="$(cd "$(dirname "${0}")" && pwd)"
root_dirpath="$(dirname "${script_dirpath}")"

# ==========================================================================================
#                                         Constants
# ==========================================================================================
KURTOSIS_DOCKERHUB_ORG="kurtosistech"
SCRIPTS_DIRNAME="scripts"
BUILD_AND_RUN_FILENAME="build-and-run.sh"
TESTSUITE_DIRNAME="testsuite"

ERROR_LOG_KEYWORD="ERRO"

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

# Building and running testsuite take a very long time, so we do some optimizations:
# 1) skip building/running testsuite if only docs changes
if git --no-pager diff --exit-code origin/develop...HEAD -- . ':!*.md' >/dev/null; then
  echo "Skipping building and running testsuite as the only changes are in Markdown files"
  exit 0
fi

# 2) skip building/running testsuite if testsuite folder doesn't have any changes
if git --no-pager diff --exit-code origin/develop...HEAD -- "${TESTSUITE_DIRNAME}" >/dev/null; then
  echo "Skipping building and running testsuite because testsuite folder doesn't have any changes"
  exit 0
fi

# 3) build and run testsuite
echo "Building and running testsuite..."

buildscript_filepath="${root_dirpath}/${SCRIPTS_DIRNAME}/${BUILD_AND_RUN_FILENAME}"
output_filepath="$(mktemp)"
if ! bash "${buildscript_filepath}" all --client-id "${kurtosis_client_id}" --client-secret "${kurtosis_client_secret_DO_NOT_LOG}" 2>&1 | tee "${output_filepath}"; then
  echo "Error: Building and running testsuite failed" >&2
  exit 1
fi
echo "Successfully built and run testsuite"

# This helps us catch errors that might show up in the testsuite logs but not get propagated to the actual exit codes
echo "Scanning output file for error log keyword '${ERROR_LOG_KEYWORD}'..."
if grep "${ERROR_LOG_KEYWORD}" "${output_filepath}"; then
  echo "Error: Detected error log pattern '${ERROR_LOG_KEYWORD}' in output file" >&2
  exit 1
fi
echo "No instances of error log keyword found"

echo "Successfully built and ran testsuite in need of validation"
