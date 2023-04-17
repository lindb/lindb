#!/usr/bin/env bash

CURRENT_DIR=$(dirname $(readlink -f "$0"))

BIN_DIR="bin"
RELEASE_DIR="release"
PACKAGE_NAME="lindb"

build() {
	target=$1

	LD_FLAGS=("-s -w -X github.com/lindb/lindb/config.Version=${PACKAGE_VERSION}")
	LD_FLAGS+=("-X github.com/lindb/lindb/config.BuildTime=${BUILD_TIME}")


	GO_BUILD_ENV=("GOOS=${GOOS}" "GOARCH=${GOARCH}")
	if [ "${GOOS}" == "darwin" ]; then
		# macOS need enable cgo
		# https://github.com/shirou/gopsutil/issues/592
		GO_BUILD_ENV=("GOOS=${GOOS}" "GOARCH=${GOARCH}" "CGO_ENABLED=1")
	fi
	target_file="${BIN_DIR}/${target}/${PACKAGE_NAME}"

	# build go package
	env "${GO_BUILD_ENV[@]}" go build -o "${BIN_DIR}/${target}/lind-cli" "-ldflags=${LD_FLAGS[*]}" ./cmd/cli
	env "${GO_BUILD_ENV[@]}" go build -o "${BIN_DIR}/${target}/lind" "-ldflags=${LD_FLAGS[*]}" ./cmd/lind

	# if windows os rename binary package name
	if [ "${GOOS}" == "windows" ]; then
		mv "${BIN_DIR}/${target}/lind-cli" "${BIN_DIR}/${target}/lind-cli.exe"
		mv "${BIN_DIR}/${target}/lind" "${BIN_DIR}/${target}/lind.exe"
	fi
}

clean() {
	rm -rf "${BIN_DIR}"
	rm -rf "${RELEASE_DIR}"
}

function main() {
	# get version
	if [[ -z "${VERSION}" ]]; then 
		# if env not set, get git tag as version
    	ver="$(git describe --tags --exact-match --match "v*.*.*" \
  				|| git describe --match "v*.*.*" --tags \
  				|| git describe --tags \
  				|| git rev-parse HEAD)"
	else
		ver="${VERSION}"
	fi

	export PACKAGE_VERSION=${ver}
	export BUILD_TIME=$(date "+%Y-%m-%dT%H:%M:%S%z")

	echo "start build ${PACKAGE_NAME} release packages, version: ${PACKAGE_VERSION}, build time: ${BUILD_TIME}"	

	# clean old build result
	clean
	# mkdir release dir
	mkdir "${RELEASE_DIR}"

	# build package for supported os
	for os in darwin windows linux; do
		export GOOS=${os}
	
		# default os arch
		TARGET_ARCHS=("amd64")
	
		# if linux os build arm
		if [ ${GOOS} == "linux" ]; then
			TARGET_ARCHS+=("arm64")
		fi
	
		# build each os arch package
		for TARGET_ARCH in "${TARGET_ARCHS[@]}"; do 
			export GOARCH=${TARGET_ARCH}

			# target folder name
			TARGET="${PACKAGE_NAME}-${PACKAGE_VERSION}-${GOOS}-${GOARCH}"
	
			echo "start build ${TARGET} ....."
	
			build "${TARGET}"

			# compress binary package, build final release package
			cd "${CURRENT_DIR}/${BIN_DIR}"
			if [ ${GOOS} == "linux" ]; then
				tar -zcf "${TARGET}.tar.gz" ${TARGET}
				cp "${TARGET}.tar.gz" "${CURRENT_DIR}/${RELEASE_DIR}/"
			else
				zip -qr "${TARGET}.zip" ${TARGET}
				cp "${TARGET}.zip" "${CURRENT_DIR}/${RELEASE_DIR}"
			fi
			cd "${CURRENT_DIR}"

			echo "build ${TARGET} complete. ...."
		done
	done
}

main
