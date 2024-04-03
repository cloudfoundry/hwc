#!/bin/bash

set -eu
set -o pipefail

THIS_FILE_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"
WORKSPACE_DIR="${THIS_FILE_DIR}/../.."
CI="${WORKSPACE_DIR}/wg-app-platform-runtime-ci"
. "$CI/shared/helpers/git-helpers.bash"
REPO_NAME=$(git_get_remote_name)
REPO_PATH="${THIS_FILE_DIR}/../"
unset THIS_FILE_DIR

package=${1:-.}
echo "Testing ${package}"

FLY_OS=windows \
DIR="${package}" \
"$CI/bin/fly-exec.bash" run-bin-test -i repo="${REPO_PATH}" -p
