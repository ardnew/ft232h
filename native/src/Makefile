################################################################################
#                                                                              #
# file: Makefile                                                               #
# date: 18 Feb 2020                                                            #
# auth: ardnew (andrew@ardnew.com)                                             #
# desc: GNU Makefile to build and install native library files required to     #
#       interface with FTDI FT232H using Go module "github.com/ardnew/ft232h". #
#                                                                              #
#       This builds the libft232h library, an amalgamation of FTDI's driver    #
#       libraries libMPSSE and libftd2xx, modified for static linkage support, #
#       particularly convenient for the Go build system and package manager.   #
#                                                                              #
#       Run "make" without any arguments to see available options.             #
#                                                                              #
################################################################################

platforms := linux-amd64 linux-386 linux-arm64 linux-arm darwin-amd64

platform ?= linux-amd64

ifeq "" "$(strip $(filter $(platforms),$(platform)))"
$(error unsupported platform "$(platform)" (see: "make help"))
endif

# --------------------------------------------------- TOOLCHAIN CONFIGURATION --

# if you are cross-compiling for a non-native host, set the cross variable to
# the prefix of your cross-compiler. a few examples are commented out below.
cross   ?=
#cross   ?= i686-linux-gnu-
#cross   ?= arm-linux-gnueabihf-
#cross   ?= aarch64-linux-gnu-
objdump  = $(cross)objdump
CC       = $(cross)gcc
AR       = $(cross)ar
LD       = $(cross)ld


# ------------------------------------------------------ SOURCE CONFIGURATION --

# name of the shared library we are actually building
target = libft232h

# version numbers for the libftd2xx proprietary library
linux-amd64-ftd2xx-vers  = 1.4.8
linux-386-ftd2xx-vers    = 1.4.8
linux-arm64-ftd2xx-vers  = 1.4.8
linux-arm-ftd2xx-vers    = 1.4.8
darwin-amd64-ftd2xx-vers = 1.4.16


# ------------------------------------------------------- BUILD CONFIGURATION --

# parse the OS name ("linux", "darwin", ...) from the platform specification
os = $(word 1,$(subst -, ,$(platform)))

# where to install the resulting libft232h build output. you shouldn't need to
# change this for any reason. the Go code looks for the libft232h library in
# this directory at compile-time (via cgo preamble in native_bridge.go).
cgolib = ../lib/$(platform)

# add search paths for the libftd2xx header/lib and libMPSSE headers
includes = $(addprefix -I,. ../inc)
libpaths = $(addprefix -L,$(ftd2xx-root))
libnames = $(addprefix -l,$(ftd2xx-name))

# all sources files that constitute libMPSSE
i2c-src = ftdi_i2c.c
spi-src = ftdi_spi.c
sources = ftdi_infra.c ftdi_mid.c $(i2c-src) $(spi-src)
objname = obj
objpath = $(platform)/$(objname)
objects = $(patsubst %,$(objpath)/%,$(sources:.c=.o)) # expected object files
outputs = $(platform)/$(target).a

# paths and object files related to the libftd2xx proprietary library
ftd2xx-name        = ftd2xx
ftd2xx-root        = $(platform)/lib$(ftd2xx-name)/$($(platform)-ftd2xx-vers)
ftd2xx-objpath     = $(ftd2xx-root)/$(objname)
ftd2xx-lib         = $(ftd2xx-root)/lib$(ftd2xx-name).a
ftd2xx-content     = $(filter %.o,$(shell $(AR) -t "$(ftd2xx-lib)")) # must stay sorted!
ftd2xx-objects     = $(patsubst %,$(ftd2xx-objpath)/%,$(ftd2xx-content))

# command used to generate build metadata info
linux-objdump  = $(objdump) -fpt
darwin-objdump = $(objdump) -macho -dylibs-used -archive-headers -private-header -t -indirect-symbols -g

defines  = -DFTD2XX_STATIC
#defines  += -DINFRA_DEBUG_ENABLE
optimize = -Og
debug    = -g
warnings = -Wall
cflags   = -std=gnu11
ldflags  = -static
linux-cflags   = $(cflags) -fPIC
linux-ldflags  = $(ldflags)
darwin-cflags  = $(cflags)
darwin-ldflags = $(ldflags)

CFLAGS  = $(debug) $(optimize) $(warnings) $(defines) $(includes) $($(os)-cflags)
LDFLAGS = $($(os)-ldflags) $(libpaths) $(libnames)

buildtime := $(shell date)


# -------------------------------------------------------------- MAKE RECIPES --

.PHONY: default help available showconfig build clean

default: showconfig build

help: available
	@printf "to build for a specific target, set the 'platform' variable when calling make.\n"
	@printf "to use a specific cross-compiler, also set the 'cross' variable.\n\n"
	@printf "for example, to build for macOS using native toolchain:\n"
	@printf "	$$ make platform=\"darwin-amd64\"\n\n"
	@printf "or to build for Raspberry Pi model 4 using a GCC cross-compiler:\n"
	@printf "	$$ make platform=\"linux-arm64\" cross=\"aarch64-linux-gnu-\"\n\n"
	@printf "supported platforms:\n"
	@printf "	%s\n" $(platforms)

available:
	@printf "available make targets:\n"
	@printf "	clean		- removes all compiled objects (.o, .a, .so, .dylib) and metadata\n"
	@printf "	build		- compiles and installs the static library (default)\n\n"

showconfig:
	@printf "current build configuration:\n"
	@printf "	platform	= %s\n" $(platform)
	@printf "	build output	= %s\n" ./$(platform)/
	@printf "	install path	= %s\n" $(cgolib)
	@printf "	CC		= %s\n" "$(CC)"
	@printf "	CFLAGS		= %s\n" "$(CFLAGS)"
	@printf "	LDFLAGS		= %s\n\n" "$(LDFLAGS)"

build: $(showconfig) $(outputs) $(cgolib)
	@printf "\n\t** build complete: %s (%s)\n\n" "$(target).a" "$(buildtime)" | tee "$(platform)/$(target).a.info"
	@$($(os)-objdump) "$(platform)/$(target).a" >> "$(platform)/$(target).a.info"
	@cp $(outputs) "$(cgolib)"

clean:
	rm -rf "$(objpath)" "$(ftd2xx-objpath)" "$(platform)"/*.a "$(platform)"/*.info

$(cgolib):
	@test -d "$@" || mkdir -p "$@"

$(objpath):
	@test -d "$@" || mkdir -p "$@"

$(ftd2xx-objpath):
	@test -d "$@" || mkdir -p "$@"
	@# perform in subshell, because pushd/popd doesnt seem to have an effect
	( cd $@; $(AR) -x ../"lib$(ftd2xx-name).a" $(ftd2xx-content) )

$(platform)/$(target).a: $(objpath) $(objects) $(ftd2xx-objpath)
	$(AR) -rcs "$@" $(objects) $(ftd2xx-objects)

$(objects): $(objpath)/%.o: %.c
	$(CC) $(CFLAGS) -c -o $@ $<

