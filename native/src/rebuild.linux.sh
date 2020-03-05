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

# remove and potentially use the debug flag if provided, i.e. (ft)debug(=1), but not (ft)debug=0
for (( i = 1; i <= $#; ++i )) do
	if [[ ${!i} =~ debug ]]; then
		[[ ${!i} =~ (^|[[:space:]])(ft)?debug(=[^[:space:]]*[1-9a-zA-Z][^[:space:]]*)?($|[[:space:]]) ]] &&
			debug=1
	else
		given=( ${given[@]} ${!i} )
	fi
done

# use the given make targets if any were provided
[[ ${#given} -gt 0 ]] && targets="${given[@]}"

# add the debug flag if provided
[[ -n ${debug} ]] && targets="ftdebug=1 $targets"

rebuild  "$targets"  linux-amd64
rebuild  "$targets"  linux-386    i686-linux-gnu-
rebuild  "$targets"  linux-arm64  aarch64-linux-gnu-
rebuild  "$targets"  linux-arm    arm-linux-gnueabihf-
