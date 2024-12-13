#!/usr/bin/env bash

set -e

COL_GREEN='[32m'
COL_RED='[91m'
COL_NC='[0m'

INFO="[i]"
TICK="[${COL_GREEN}✓${COL_NC}]"
CROSS="[${COL_RED}✗${COL_NC}]"

TMP_DIR="$(mktemp -d)"
TMP_RELEASE_TAR="$TMP_DIR/release.tar.gz"
OUTPUT_DIR="/usr/local/bin"
PROGRAM_NAME="adless"

check_system() {
  local uname_os="$(uname -s)"
  local uname_arch="$(uname -m)"

  case $uname_os in
    Linux*) OS="linux" ;;
    Darwin*) OS="darwin" ;;
    *) fail "${uname_os} operation system is unsupported" ;;
  esac

  case $uname_arch in
    x86_64 | amd64) ARCH="x86_64" ;;
    arm64 | arm) ARCH="arm64" ;;
    *) fail "${uname_arch} arch is unsupported" ;;
  esac

  echo "${INFO} Detected OS: ${OS}_${ARCH}"
}

check_dependencies() {
  which curl >/dev/null || fail "curl not installed"
  which grep >/dev/null || fail "grep not installed"
  which sed >/dev/null || fail "sed not installed"
}

find_latest_version() {
  LATEST_VERSION=$(curl -s "https://api.github.com/repos/wittyjudge/adless/releases/latest" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')
  echo -e "${INFO} Latest version: $LATEST_VERSION"
}

find_latest_release_tar_url() {
  URL="https://github.com//WIttyJudge/adless/releases/download/${LATEST_VERSION}/adless_${LATEST_VERSION}_${OS}_${ARCH}.tar.gz"
  echo -e "${INFO} Latest release tar: ${URL}"
}

install_release_tar() {
  echo -e "${INFO} Downloading release tar ${URL}.."
  curl --fail --progress-bar -L -o "$TMP_RELEASE_TAR" "$URL"
  echo -e "${TICK} Release tar downloaded"

  tar -xzf "$TMP_RELEASE_TAR" -C "$TMP_DIR"

  mv "$TMP_DIR/${PROGRAM_NAME}" "${OUTPUT_DIR}" || fail "Failed to make program executable, re-run the command using \"sudo bash\""
  chmod +x "${OUTPUT_DIR}/${PROGRAM_NAME}" ||  fail "Failed to move binary, re-run the command using \"sudo bash\""

  cleanup
}

finish() {
  echo ""
  echo "Adless ${LATEST_VERSION} successfully installed at ${OUTPUT_DIR}/${PROGRAM_NAME}"
  echo "run: adless --help"
}

cleanup() {
  rm -rf $TMP_DIR >/dev/null
}

fail() {
  cleanup
  msg=$1
  echo "${CROSS} $msg" 1>&2
  exit 1
}

# Execution

cat <<'EOF'
           _ _
  __ _  __| | | ___  ___ ___
 / _` |/ _` | |/ _ \/ __/ __|
| (_| | (_| | |  __/\__ \__ \
 \__,_|\__,_|_|\___||___/___/

EOF

check_system
check_dependencies

find_latest_version
find_latest_release_tar_url

install_release_tar

finish
