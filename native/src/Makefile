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

# recognized/supported platforms
platforms := linux-amd64 linux-386 linux-arm64 linux-arm darwin-amd64

# default build target
platform ?= linux-amd64

# invalid build target provided
ifeq "" "$(strip $(filter $(platforms),$(platform)))"
$(error unsupported platform "$(platform)" (see: "make help"))
endif

# --------------------------------------------------- TOOLCHAIN CONFIGURATION --

# if you are cross-compiling for a non-native host, set the cross variable to
# the prefix of your cross-compiler. a few examples are commented out below.
# you can safely use an absolute/relative path in cross as well.
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
os   = $(word 1,$(subst -, ,$(platform)))
arch = $(word 2,$(subst -, ,$(platform)))

# where to install the resulting libft232h build output. you shouldn't need to
# change this for any reason. the Go code looks for the libft232h library in
# this directory at compile-time (via cgo preamble in native_bridge.go).
install = ../lib/$(platform)

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

# other data included in cgo preamble
cgo-includes      = libMPSSE_spi.h libMPSSE_i2c.h ftd2xx.h stdlib.h
darwin-frameworks = CoreFoundation IOKit

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
ifeq "" "$(strip $(ftdebug))"
ftdebug-yesno  = NO
else
ftdebug-yesno  = YES
defines       += -DINFRA_DEBUG_ENABLE
endif
optimize       = -Og
debug          = -g
warnings       = -Wall
cflags         = -std=gnu11
linux-cflags   = $(cflags) -fPIC
darwin-cflags  = $(cflags)

CFLAGS = $(debug) $(optimize) $(warnings) $(defines) $(includes) $($(os)-cflags)

buildtime := $(shell date +"%Y-%b-%d %T %Z")
marg  := 4
cols  := 12

# define the macro $(,) so that we can use commas in e.g. $(subst x,$(,),1x2x3)
,     := ,


# -------------------------------------------------------------- MAKE RECIPES --

.PHONY: default help available showconfig preamble build clean

default: build

help: available showconfig preamble
	@printf -- "\n"
	@printf -- "--------------------------------------------------------------------------------\n"
	@printf -- "\n"
	@printf -- "to build for a specific target, set the 'platform' variable when calling make.\n"
	@printf -- "\n"
	@printf -- "to use a specific cross-compiler, also set the 'cross' variable.\n"
	@printf -- "\n"
	@printf -- "  - for example, to build for macOS using native toolchain:\n"
	@printf -- "\n"
	@printf -- "	$$ make platform=\"darwin-amd64\"\n"
	@printf -- "\n"
	@printf -- "  - or to build for Raspberry Pi model 4 using a GCC cross-compiler:\n"
	@printf -- "\n"
	@printf -- "	$$ make platform=\"linux-arm64\" cross=\"aarch64-linux-gnu-\"\n"
	@printf -- "\n"
	@printf -- "  - to force a complete fresh rebuild and include driver debug printing:\n"
	@printf -- "\n"
	@printf -- "	$$ make ftdebug=1 clean build\n"
	@printf -- "\n"
	@printf -- "supported platforms:\n"
	@printf -- "\n"
	@printf -- "%*s%s\n" $(patsubst %,$(marg) " " %,$(platforms))
	@printf -- "\n"

available:
	@printf -- "\n"
	@printf -- "--------------------------------------------------------------------------------\n"
	@printf -- "\n"
	@printf -- "available make targets:\n"
	@printf -- "\n"
	@printf -- "%*s %-*s - prints available targets, usage examples, and other stuff\n"  $(marg) " " $(cols) "help"
	@printf -- "%*s %-*s - prints the build configuration for a given platform\n"        $(marg) " " $(cols) "showconfig"
	@printf -- "%*s %-*s - prints the suggested cgo preamble to use libft232h.a\n"       $(marg) " " $(cols) "preamble"
	@printf -- "%*s %-*s - removes all compiled objects (.o, .a) and metadata\n"         $(marg) " " $(cols) "clean"
	@printf -- "%*s %-*s - compiles and installs the static library (default)\n"         $(marg) " " $(cols) "build"
	@printf -- "\n"

showconfig:
	@printf -- "\n"
	@printf -- "--------------------------------------------------------------------------------\n"
	@printf -- "\n"
	@printf -- "current build configuration:\n"
	@printf -- "\n"
	@printf -- "%*s %-*s = %s\n" $(marg) " " $(cols) "platform" "$(platform)"
	@printf -- "%*s %-*s = %s\n" $(marg) " " $(cols) "install" "$(install)/"
	@printf -- "%*s %-*s = %s\n" $(marg) " " $(cols) "build path" "./$(platform)/"
	@printf -- "%*s %-*s = %s\n" $(marg) " " $(cols) "ftdebug" "$(ftdebug-yesno)"
	@printf -- "%*s %-*s = %s\n" $(marg) " " $(cols) "CC" "$(CC)"
	@printf -- "%*s %-*s = %s\n" $(marg) " " $(cols) "CFLAGS" "$(CFLAGS)"
	@printf -- "\n"

preamble:
	@printf -- "\n"
	@printf -- "--------------------------------------------------------------------------------\n"
	@printf -- "\n"
	@printf -- "cgo preamble to use with current target:\n"
	@printf -- "\n"
	@printf -- "%*s // #cgo               CFLAGS: -I%s\n" $(marg) " " \$${SRCDIR}/native/inc
	@test "darwin" = "$(os)" && (\
	 printf -- "%*s // #cgo %-12s LDFLAGS: %s\n"          $(marg) " " $(subst -,$(,),$(platform)) "$(patsubst %,-framework %,$(darwin-frameworks))" ) || true
	@printf -- "%*s // #cgo %-12s LDFLAGS: -L%s\n"        $(marg) " " $(subst -,$(,),$(platform)) \$${SRCDIR}/native/lib/$(platform)
	@printf -- "%*s // #cgo              LDFLAGS: -l%s\n" $(marg) " " "$(subst lib,,$(target))"
	@printf -- "%*s // #include \"%s\"\n"    $(patsubst %,$(marg) " " %,$(cgo-includes))
	@printf -- "%*s import \"C\"\n"                       $(marg) " "
	@printf -- "\n"

build: showconfig $(outputs) $(install) preamble
	@printf -- "\n"
	@printf -- "\t** BUILD SUCCESS | %s [%s] | %s **\n\n" "$(target).a" "$(platform)" "$(buildtime)" | tee "$(platform)/$(target).a.info"
	@$($(os)-objdump) "$(platform)/$(target).a" >> "$(platform)/$(target).a.info"
	@cp $(outputs) "$(install)"

clean:
	rm -rf "$(objpath)" "$(ftd2xx-objpath)" "$(platform)"/*.a "$(platform)"/*.info

$(install):
	@test -d "$@" || mkdir -p "$@"

$(objpath):
	@test -d "$@" || mkdir -p "$@"

$(ftd2xx-objpath):
	@printf -- "\n"
	@test -d "$@" || mkdir -p "$@"
	@# perform in subshell, because pushd/popd doesnt seem to have an effect
	( cd $@; $(AR) -x ../"lib$(ftd2xx-name).a" $(ftd2xx-content) )

$(platform)/$(target).a: $(objpath) $(objects) $(ftd2xx-objpath)
	@printf -- "\n"
	$(AR) -rcs "$@" $(objects) $(ftd2xx-objects)

$(objects): $(objpath)/%.o: %.c
	$(CC) $(CFLAGS) -c -o $@ $<

