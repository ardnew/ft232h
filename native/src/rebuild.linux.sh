#!/bin/bash
# ------------------------------------------------------------------------------
#
#  build fresh libft232h.a needed by github.com/ardnew/ft232h cgo interface for
#  all of the linux platforms.
#
#  i ~think~ only the following cross-compiler packages are needed on a host
#  running 64-bit Ubuntu (amd64):
#
#    gcc-i686-linux-gnu   gcc-aarch64-linux-gnu   gcc-arm-linux-gnueabihf
#    ------------------ | --------------------- | -----------------------
#        linux-386             linux-arm64               linux-arm
#
# ------------------------------------------------------------------------------

rebuild()
{
	targets=$1
	platform="platform=$2"
	[[ $# -gt 2 ]] && cross="cross=$3"

	banner=$(perl -le '$s=(shift); printf "%s [%s] %s", "="x(74-length($s)), $s, "="x2' "$2")
	printf -- "\n%s\n\n" "$banner"

	if ! make $platform $cross $targets; then
		printf -- "\n\t** BUILD FAILED | -- [%s] | %s **\n\n" "$platform" "$(date +'%Y-%b-%d %T %Z')"
	fi
}

# clean and rebuild by default
targets="clean build"
# replace with all targets given on command line (if any)
[[ $# -gt 0 ]] && targets=$@

rebuild  "$targets"  linux-amd64
rebuild  "$targets"  linux-386    i686-linux-gnu-
rebuild  "$targets"  linux-arm64  aarch64-linux-gnu-
rebuild  "$targets"  linux-arm    arm-linux-gnueabihf-
