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


# -------------------------------------------------------- HOST CONFIGURATION --

# host OS and arch where you will be running your Go executables. if you are
# cross-compiling for a different host, be sure to set the `cross` variable
# further down below matching the prefix of your cross-compiler toolchain.
os       = linux
#os       = darwin
arch     = amd64
#arch     = 386
#arch     = arm
#arch     = arm64


# --------------------------------------------------- TOOLCHAIN CONFIGURATION --

cross   =
#cross   = i686-linux-gnu-
#cross   = arm-linux-gnueabihf-
#cross   = aarch64-linux-gnu-
objdump = $(cross)objdump
CC      = $(cross)gcc
AR      = $(cross)ar
LD      = $(cross)ld


# ------------------------------------------------------ SOURCE CONFIGURATION --

# name of the shared library we are actually building
target = libft232h

# version numbers for the libftd2xx proprietary library
linux-ftd2xx-vers  = 1.4.8
darwin-ftd2xx-vers = 1.4.16


################################################################################
#                                                                              #
#       //-----------================####=================-----------\\        #
#      |||                                                           |||       #
#      ||| You should know what you're doing if you intend to make   |||       #
#      ||| changes to any of the following sections. Otherwise, the  |||       #
#      ||| basic configuration should suffice and you should now try |||       #
#      ||| to build the library with `make build && make install`.   |||       #
#      |||                                                           |||       #
#       \\-----------================####=================-----------//        #
#                                                                              #
################################################################################


# ------------------------------------------------------- BUILD CONFIGURATION --

# where to install the resulting libft232h build output. you shouldn't need to
# change this for any reason. the Go code looks for the libft232h library in
# this directory at compile-time (via cgo preamble in native_bridge.go).
install  = ../lib/$(os)-$(arch)

# add library search paths for the libftd2xx driver, and include paths to the
# libMPSSE API header files
includes = $(addprefix -I,. ../inc)
libpaths = $(addprefix -L,$(ftd2xx-root))
libnames = $(addprefix -l,$(ftd2xx-name))

# all sources files that constitute libMPSSE
i2c-src = ftdi_i2c.c
spi-src = ftdi_spi.c
sources = ftdi_infra.c ftdi_mid.c $(i2c-src) $(spi-src)
objname = obj
objpath = $(os)-$(arch)/$(objname)
objects = $(patsubst %,$(objpath)/%,$(sources:.c=.o)) # expected object files
outputs = $(os)-$(arch)/$(target).a

# paths and object files related to the libftd2xx proprietary library
ftd2xx-name        = ftd2xx
ftd2xx-root        = $(os)-$(arch)/lib$(ftd2xx-name)-$(os)-$(arch)-$($(os)-ftd2xx-vers)
ftd2xx-objpath     = $(ftd2xx-root)/$(objname)
ftd2xx-lib         = $(ftd2xx-root)/lib$(ftd2xx-name).a
ftd2xx-content    := $(filter %.o,$(shell $(AR) -t "$(ftd2xx-lib)")) # must stay sorted!
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

.PHONY: default notarget available showconfig install reinstall build clean

default: notarget available showconfig

notarget:
	@printf "no target specified. usage:\n"
	@printf "	$$ make <target>\n\n"

available:
	@printf "available make targets:\n"
	@printf "	showconfig	- prints the current build/install configuration\n"
	@printf "	install		- copies %s to install path %s\n" "$(target)" "$(install)"
	@printf "	reinstall	- compiles and copies %s to install path %s\n" "$(target)" "$(install)"
	@printf "	build		- compiles the static library %s\n" "$(target).a"
	@printf "	clean		- removes all compiled objects (.o, .a, .so, .dylib) and metadata\n\n"

showconfig:
	@printf "current build configuration:\n"
	@printf "	OS		= %s\n" $(os)
	@printf "	arch		= %s\n" $(arch)
	@printf "	install path	= %s\n" $(install)
	@printf "	build output	= %s\n" ./$(os)-$(arch)/
	@printf "	CC		= %s\n" "$(CC)"
	@printf "	CFLAGS		= %s\n" "$(CFLAGS)"
	@printf "	LDFLAGS		= %s\n" "$(LDFLAGS)"

install: $(install)
	cp $(outputs) "$^"
	@printf "\n\t** install complete: %s/%s\n" "$^" $(outputs)

reinstall: $(outputs) install

build: $(outputs)

clean:
	rm -rf "$(objpath)" "$(ftd2xx-objpath)" "$(os)-$(arch)"/*.a "$(os)-$(arch)"/*.info

$(install):
	@test -d "$@" || mkdir -p "$@"

$(objpath):
	@test -d "$@" || mkdir -p "$@"

$(ftd2xx-objpath):
	@test -d "$@" || mkdir -p "$@"
	@# perform in subshell, because pushd/popd doesnt seem to have an effect
	( cd $@; $(AR) -x ../"lib$(ftd2xx-name).a" $(ftd2xx-content) )

$(os)-$(arch)/$(target).a: $(objpath) $(objects) $(ftd2xx-objpath)
	$(AR) -rcs "$@" $(objects) $(ftd2xx-objects)
	@printf "\n\t** build complete: %s (%s)\n\n" "$(target).a" "$(buildtime)" | tee "$@.info"
	@$($(os)-objdump) "$@" >> "$@.info"

$(objects): $(objpath)/%.o: %.c
	$(CC) $(CFLAGS) -c -o $@ $<

