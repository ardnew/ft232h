/*
High-level driver for the FTDI FT232H USB to GPIO/SPI/I²C/JTAG/UART protocol
converter.

Developing Peripheral Device Drivers

  TBD(andrew@ardnew.com): add driver development details

API Design

The design of this module API was intended to marry the following principles,
ordered by importance, and the bulleted items underneath each are how they are
realized:

  1. Simple / Concise –- Clear, consistent, conventional behavior
     - Document all exported features
     - Transparent infrastructure, auto-configure as much as possible
     - Minimize subtleties and unapparent side-effects (or document
       those that exist)
     - Prefer Go-native and fixed-size composite types
     - Design flexible, high-level protocol methods
     - Integrate command-line flags
  2. Robust / General Purpose -– Maximize peripheral device support
     - Provide clear, descriptive error messages
     - Integration test every revision automatically
     - Verify correct behavior at edge cases
     - Isolate system-dependent functionality
     - Expose low-level device methods for granular control
  3. Performant / Efficient –- Utilize resources effectively
     - Minimize USB transactions and HID interframe delays
     - Maximize throughput of serial protocol transfers

FTDI Native Drivers

FTDI uses a custom (vendor-defined) protocol to communicate with their USB
devices and releases proprietary driver software `D2XX` (binary-only) that
application programmers (see: YOU) should use with FTDI USB devices.

The `D2XX` driver is thus only available for systems officially supported by
FTDI. These drivers are fairly low-level, containing only C source code headers
and a thin user guide as documentation.

These drivers have been modified and recompiled into a single static library for
each supported platform and distributed with this Go module, so the user isn't
required to compile them on his/her own.

For more details and instructions on compiling the library yourself, refer to
the godoc of package github.com/ardnew/ft232h/native.

*/
package ft232h
