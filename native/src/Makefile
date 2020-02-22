################################################################################
#
# file: Makefile
# date: 18 Feb 2020
# auth: ardnew (andrew@ardnew.com)
# desc: GNU Makefile to build and install native library files required to
#       interface with FTDI FT232H using Go module "github.com/ardnew/ft232h".
#       This builds LibMPSSE, which depends on driver FTD2XX (proprietary binary
#       included with this package under relevant OS-arch subdir).
#
#       Run "make" without any arguments to see available options.
#
################################################################################


# -------------------------------------------------------- HOST CONFIGURATION --

# host OS and arch where you will be running your Go executables. the actual
# values used here are inherited from $GOOS and $GOARCH, and are what I used for
# directory names for any software that is system-specific (incl. proprietary)
os       = linux
#os       = darwin
arch     = amd64
#arch     = 386
#arch     = arm64

# where the resulting required libs will be installed. this should be somewhere
# fairly global so that your Go executables won't have any trouble locating them
install  = /usr/local/lib


# ------------------------------------------------------ SOURCE CONFIGURATION --

# name  of the shared library we are actually building
target = libMPSSE

# path to ftd2xx proprietary library and headers, constructed with system-
# specific path names and version numbers.
ftd2xx-name        = ftd2xx
linux-ftd2xx-vers  = 1.4.8
darwin-ftd2xx-vers = 1.4.16
ftd2xx-root        = $(os)-$(arch)/lib$(ftd2xx-name)-$(os)-$(arch)-$($(os)-ftd2xx-vers)


#       //-----------================####=================-----------\\        #
#      |||                                                           |||       #
#      ||| You should know what you're doing if you intend to make   |||       #
#      ||| changes to any of the following sections. Otherwise, the  |||       #
#      ||| basic configuration should suffice and you should now try |||       #
#      ||| to build the library with `make && make install`.         |||       #
#      |||                                                           |||       #
#       \\-----------================####=================-----------//        #


# ------------------------------------------------------- BUILD CONFIGURATION --

# add include and library search paths for the ftd2xx driver, and the include
# path to the libMPSSE API header files
includes = $(addprefix -I,$(ftd2xx-root) . ../inc)
libpaths = $(addprefix -L,$(ftd2xx-root))
libnames = $(addprefix -l,$(ftd2xx-name))

# all sources files that constitute libMPSSE
i2c-src = ftdi_i2c.c
spi-src = ftdi_spi.c
sources = ftdi_infra.c ftdi_mid.c $(i2c-src) $(spi-src)
objpath = $(os)-$(arch)/obj
objects = $(patsubst %,$(objpath)/%,$(sources:.c=.o)) # expected object files

# file naming convention for shared/dynamic libraries
linux-dynext      = so
darwin-dynext     = dylib
linux-ftd2xx-dynext  = $($(os)-dynext).$($(os)-ftd2xx-vers)
darwin-ftd2xx-dynext = $($(os)-ftd2xx-vers).$($(os)-dynext) # macOS reverses the extensions

# command used to generate metadata info about both compiled libraries (static
# and dynamic)
linux-objdump  = $(objdump) -fpt
darwin-objdump = $(objdump) -macho -dylibs-used -archive-headers -private-header -t -indirect-symbols -g

# all files required to use libMPSSE from a Go application are considered
# "outputs" from this Makefile. running a successful `make install` copies these
# files to the "fairly global" directory $(install) defined above.
outputs  = $(os)-$(arch)/$(target).a \
           $(os)-$(arch)/$(target).$($(os)-dynext) \
           $(ftd2xx-root)/lib$(ftd2xx-name).a \
           $(ftd2xx-root)/lib$(ftd2xx-name).$($(os)-ftd2xx-dynext) \
           $(ftd2xx-root)/lib$(ftd2xx-name).$($(os)-dynext)

defines  =
#defines  += -DINFRA_DEBUG_ENABLE
optimize = -Og
debug    = -g
warnings = -Wall
linux-cflags   = -fPIC
linux-ldflags  = -fPIC -shared
darwin-cflags  =
darwin-ldflags = -dynamiclib

cross   =
#cross   = i686-linux-gnu-
#cross   = aarch64-linux-gnu-
objdump = $(cross)objdump
CC      = $(cross)gcc
AR      = $(cross)ar
LD      = $(cross)ld
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
	@printf "	install		- installs both %s and %s to system library path\n" $(target) $(ftd2xx-name)
	@printf "	reinstall		- compiles and installs both %s and %s to system library path\n" $(target) $(ftd2xx-name)
	@printf "	build		- compiles both a static and shared/dynamic library %s\n" $(target)
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
	@test -d $^ || mkdir -p $^
	cp $(outputs) $^
	@echo "libraries installed to: $^"
	@printf "\t%s\n" $(outputs)

reinstall: $(outputs) install

build: $(outputs)

clean:
	rm -rf $(objpath) $(os)-$(arch)/*.a $(os)-$(arch)/*.$($(os)-dynext) $(os)-$(arch)/*.info

$(objpath):
	@test -d $@ || mkdir -p $@

$(os)-$(arch)/$(target).a: $(objpath) $(objects)
	$(AR) -rcs $@ $(objects)
	@printf "\t** build complete: %s (%s)\n\n" "$(target).a" "$(buildtime)" | tee "$@.info"
	@$($(os)-objdump) $@ >> "$@.info"

$(os)-$(arch)/$(target).$($(os)-dynext): $(objpath) $(objects)
	$(CC) -o $@ $(objects) $(LDFLAGS)
	@printf "\t** build complete: %s (%s)\n\n" "$(target).$($(os)-dynext)" "$(buildtime)" | tee "$@.info"
	@$($(os)-objdump) $@ >> "$@.info"

$(objects): $(objpath)/%.o: %.c
	$(CC) $(CFLAGS) -c -o $@ $<

