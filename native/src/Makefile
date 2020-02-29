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
#cross   ?= aarch64-linux-gnu-
#cross   ?= arm-linux-gnueabihf-
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

# add include paths for the libMPSSE and D2XX headers
includes = $(addprefix -I,. ../inc)

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

defines        = -DFTD2XX_STATIC
#defines        += -DINFRA_DEBUG_ENABLE
optimize       = -Og
debug          = -g
warnings       = -Wall
cflags         = -std=gnu11
linux-cflags   = $(cflags) -fPIC
darwin-cflags  = $(cflags)

CFLAGS = $(debug) $(optimize) $(warnings) $(defines) $(includes) $($(os)-cflags)

buildtime := $(shell date +"%Y-%b-%d %T %Z")
marg := 4
cols := 12


# -------------------------------------------------------------- MAKE RECIPES --

.PHONY: default help available showconfig build clean

default: build

help: showconfig available
	@printf -- "\n"
	@printf -- "--------------------------------------------------------------------------------\n"
	@printf -- "\n"
	@printf -- "to build for a specific target, set the 'platform' variable when calling make.\n"
	@printf -- "to use a specific cross-compiler, also set the 'cross' variable.\n"
	@printf -- "  - for example, to build for macOS using native toolchain:\n"
	@printf -- "\n"
	@printf -- "	$$ make platform=\"darwin-amd64\"\n"
	@printf -- "\n"
	@printf -- "  - or to build for Raspberry Pi model 4 using a GCC cross-compiler:\n"
	@printf -- "	$$ make platform=\"linux-arm64\" cross=\"aarch64-linux-gnu-\"\n"
	@printf -- "\n"
	@printf -- "  - to force a fresh rebuild, use both targets 'clean' & 'build' explicitly:\n"
	@printf -- "	$$ make clean build\n"
	@printf -- "\n"
	@printf -- "supported platforms:\n"
	@printf -- "%*s%-*s\n" $(patsubst %,$(marg) " " $(cols) %,$(platforms))
	@printf -- "\n"

available:
	@printf -- "\n"
	@printf -- "--------------------------------------------------------------------------------\n"
	@printf -- "\n"
	@printf -- "available make targets:\n"
	@printf -- "%*s %-*s - removes all compiled objects (.o, .a, .so, .dylib) and metadata\n" $(marg) " " $(cols) "clean"
	@printf -- "%*s %-*s - compiles and installs the static library (default)\n"              $(marg) " " $(cols) "build"
	@printf -- "\n"

showconfig:
	@printf -- "\n"
	@printf -- "--------------------------------------------------------------------------------\n"
	@printf -- "\n"
	@printf -- "current build configuration:\n"
	@printf -- "%*s %-*s = %s\n" $(marg) " " $(cols) "platform" $(platform)
	@printf -- "%*s %-*s = %s\n" $(marg) " " $(cols) "output" ./$(platform)/
	@printf -- "%*s %-*s = %s\n" $(marg) " " $(cols) "cgo path" $(cgolib)
	@printf -- "%*s %-*s = %s\n" $(marg) " " $(cols) "CC" "$(CC)"
	@printf -- "%*s %-*s = %s\n" $(marg) " " $(cols) "CFLAGS" "$(CFLAGS)"
	@printf -- "\n"

build: showconfig $(outputs) $(cgolib)
	@printf -- "\n\t** BUILD SUCCESS | %s [%s] | %s **\n\n" "$(target).a" "$(platform)" "$(buildtime)" | tee "$(platform)/$(target).a.info"
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

