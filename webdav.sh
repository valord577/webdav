#!/usr/bin/env bash

# macro
# ------
BASE_DIR=$(cd "$(dirname ${BASH_SOURCE[0]})" && pwd)

BUILD_VERSION="v1.1"
BUILD_DATE=$(date '+%Y-%m-%dT%H:%M:%SZ%:z' -u)

GO_LDFLAGS_VERSION="'github.com/valord577/webdav/cmd.version=${BUILD_VERSION}'"
GO_LDFLAGS_DATE="'github.com/valord577/webdav/cmd.datetime=${BUILD_DATE}'"

PRESET_TARGETS=$(cat <<- 'EOF'
darwin/amd64
darwin/arm64
freebsd/amd64
freebsd/arm
freebsd/arm64
linux/amd64
linux/arm
linux/arm64
linux/riscv64
netbsd/amd64
netbsd/arm
netbsd/arm64
openbsd/amd64
openbsd/arm
openbsd/arm64
openbsd/mips64
windows/amd64
windows/arm
EOF
)

# default
# ------
# if "1", print help information
PRINT_HELP_INFO="1"
# if empty, build all target
# example: "linux/amd64,linux/arm64"
BUILD_TARGET=""
BUILD_OUTPUT="${BASE_DIR}/out"

# func
usage() {
  echo -e "WebDAV Build Script."
  echo -e ""
  echo -e "usage:"
  echo -e "  bash webdav.sh [-a | -t [build_target]] [-o [build_output]]"
  echo -e ""
  echo -e "Arguments:"
  echo -e "  -a   Build preset targets of the build files."
  echo -e "  -t   Declare the targets of the build files."
  echo -e "  -o   Declare the path of the output files."
}

# func
build() {
  echo -e "target: ${BUILD_TARGET}"
  echo -e "output: ${BUILD_OUTPUT}"

  if [ "${BUILD_TARGET}" = "" ]; then
    targets=(${PRESET_TARGETS})
  else
    targets=(${BUILD_TARGET//,/ })
  fi

  go mod download

  for t in ${targets[@]}; do
    build_one "${t}"
  done
}

build_one() {
  target=(${1//\// })

  os=${target[0]}
  arch=${target[1]}
  echo -e "build -> os: ${os} | arch: ${arch}"

  START_TIME=$(date '+%s.%N' -u)

  BIN_FILE="${BUILD_OUTPUT}/webdav"
  CGO_ENABLED=0 GOOS=${os} GOARCH=${arch} go build -o "${BIN_FILE}" -ldflags "-X ${GO_LDFLAGS_VERSION} -X ${GO_LDFLAGS_DATE}"

  if [ "${os}" = "windows" ]; then
    COMPRESSED_FILE="${BUILD_OUTPUT}/webdav_${BUILD_VERSION}_${os}_${arch}.zip"
    zip -q -j "${COMPRESSED_FILE}" "${BIN_FILE}"
  else
    COMPRESSED_FILE="${BUILD_OUTPUT}/webdav_${BUILD_VERSION}_${os}_${arch}.tar.gz"
    tar -zcf "${COMPRESSED_FILE}" -C "${BUILD_OUTPUT}" "webdav"
  fi

  openssl dgst "-sha256" "${COMPRESSED_FILE}" | sed "s/SHA256(${COMPRESSED_FILE//\//\\\/})= //" > "${COMPRESSED_FILE}.sha256sum.txt"

  rm -f "${BIN_FILE}"

  END_TIME=$(date '+%s.%N' -u)
  TOOK_TIME=$(printf "%.3f" $(echo "scale=3; ${END_TIME}-${START_TIME}" | bc))
  SHA256SUM=$(cat "${COMPRESSED_FILE}.sha256sum.txt")
  echo -e "      >> took ${TOOK_TIME}s | sha256sum: ${SHA256SUM}"
}

# -~*~- main -~*~-
while getopts ":at:o:" opt; do
  case ${opt} in
    a)
      BUILD_TARGET=""
      PRINT_HELP_INFO="0"
      ;;
    t)
      BUILD_TARGET="${OPTARG}"
      PRINT_HELP_INFO="0"
      ;;
    o)
      BUILD_OUTPUT="${OPTARG}"
      PRINT_HELP_INFO="0"
      ;;

    \?)
      # invalid(unknown) option
      ;;
    \:)
      # miss option argument
      ;;
  esac
done

if [ "${PRINT_HELP_INFO}" = "1" ]; then
  usage
else
  build
fi
