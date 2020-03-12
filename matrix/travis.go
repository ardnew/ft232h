package main

import (
	"fmt"
	"strings"
)

func indentation() *indent {
	const indentWidth = 4
	return indentBy(indentWidth, 0)
}

var (
	// versions to test against
	version = Version{
		travis:  "~> 2.1",
		lang:    "go",
		os:      "linux",
		langver: []string{"1.11.x", "1.12.x", "1.13.x", "1.14.x", "master"},
		distro:  "bionic",
		xcode:   "xcode11.3",
		fast:    true,
		target: []*Target{
			&Target{
				arch: "amd64",
				os:   "linux",
				env: []*Env{
					&Env{platform: "linux-amd64", compiler: "gcc"},
					&Env{platform: "linux-386", compiler: "gcc", pkgs: []string{"gcc-i686-linux-gnu", "libc-dev-i386-cross"},
						cross: "i686-linux-gnu-", mach: "i386", setarch: "setarch i386 --verbose --32bit"},
				},
			},
			&Target{
				arch: "arm64",
				os:   "linux",
				env: []*Env{
					&Env{platform: "linux-arm64", compiler: "gcc"},
					&Env{platform: "linux-arm", compiler: "gcc", pkgs: []string{"gcc-arm-linux-gnueabihf", "libc-dev-armhf-cross"},
						cross: "arm-linux-gnueabihf-", mach: "armhf", setarch: "setarch linux32 --verbose --32bit"},
				},
			},
			&Target{
				arch: "amd64",
				os:   "osx",
				env: []*Env{
					&Env{platform: "darwin-amd64", compiler: "gcc"},
				},
			},
		},
		script: []string{
			"pushd native/src",
			"${setarch} make platform=${platform} cross=${cross} clean build",
			"popd",
			"${setarch} go test -v -short -count=1 -args ./...",
		},
	}
)

func main() {

	for _, l := range *version.line() {
		fmt.Println(l)
	}
}

type Version struct {
	travis  string
	lang    string
	os      string
	langver []string
	distro  string
	xcode   string
	fast    bool
	target  []*Target
	script  []string
}

type Env struct {
	platform string
	compiler string
	pkgs     []string
	cross    string
	mach     string
	setarch  string
}

type Target struct {
	arch string
	os   string
	env  []*Env
}

func (v *Version) line() *line {

	ln := &line{}
	ind := indentation()

	*ln = append(*ln, *v.travisLine(ind).new()...)
	*ln = append(*ln, *v.languageLine(ind).new()...)
	*ln = append(*ln, *v.osLine(ind).new()...)
	*ln = append(*ln, *v.distLine(ind).new()...)
	*ln = append(*ln, *v.osximageLine(ind).new()...)
	*ln = append(*ln, *v.jobsLine(ind).new()...)
	*ln = append(*ln, *v.scriptLine(ind).new()...)

	return ln
}

func (v *Version) travisLine(ind *indent) *line {
	ln := &line{}
	ln.add(ind, "version: %s", v.travis)
	return ln
}

func (v *Version) languageLine(ind *indent) *line {
	ln := &line{}
	ln.add(ind, "language: %s", v.lang)
	return ln
}

func (v *Version) osLine(ind *indent) *line {
	ln := &line{}
	ln.add(ind, "os: %s", v.os)
	return ln
}

func (v *Version) distLine(ind *indent) *line {
	ln := &line{}
	ln.add(ind, "dist: %s", v.distro)
	return ln
}

func (v *Version) osximageLine(ind *indent) *line {
	ln := &line{}
	ln.add(ind, "osx_image: %s", v.xcode)
	return ln
}

func (v *Version) jobsLine(ind *indent) *line {
	ln := &line{}
	ln.add(ind, "jobs:")
	ln.add(ind.by(1), "fast_finish: %t", v.fast)
	ln.add(ind.by(1), "include:")
	for _, lang := range v.langver {
		for _, host := range v.target {
			for _, targ := range host.env {
				ln.add(ind.by(2), "- name: %q", fmt.Sprintf("%s (%s %s)", targ.platform, v.lang, lang))
				//ln.add(ind.by(3), "language: %s", v.lang)
				ln.add(ind.by(2), "%s: %q", v.lang, lang)
				ln.add(ind.by(2), "arch: %s", host.arch)
				ln.add(ind.by(2), "os: %s", host.os)
				ln.add(ind.by(2), "compiler: %s", targ.compiler)
				ln.add(ind.by(2), "env:")
				if "" != targ.platform {
					ln.add(ind.by(3), "- platform=%q", targ.platform)
				}
				if "" != targ.cross {
					ln.add(ind.by(3), "- cross=%q", targ.cross)
				}
				if "" != targ.mach {
					ln.add(ind.by(3), "- mach=%q", targ.mach)
				}
				if "" != targ.setarch {
					ln.add(ind.by(3), "- setarch=%q", targ.setarch)
				}
				if "" != targ.cross && "" != targ.mach && "" != targ.setarch {
					ln.add(ind.by(2), "before_install:")
					ln.add(ind.by(3), "- apt search gcc")
					ln.add(ind.by(3), "- sudo dpkg --add-architecture %s", targ.mach)
					ln.add(ind.by(3), "- sudo apt -yq update")
					if nil != targ.pkgs && len(targ.pkgs) > 0 {
						ln.add(ind.by(3), "- sudo apt -yq install %s", strings.Join(targ.pkgs, " "))
					}
				}
			}
		}
	}
	return ln
}

func (v *Version) scriptLine(ind *indent) *line {
	ln := &line{}
	ln.add(ind, "script:")
	for _, s := range v.script {
		ln.add(ind.by(1), "- %s", s)
	}
	return ln
}

type line []string

func (l *line) add(ind *indent, fs string, ar ...interface{}) *line {
	*l = append(*l, ind.ent(fmt.Sprintf(fs, ar...)))
	return l
}

func (l *line) new() *line {
	*l = append(*l, "") // add an empty line break
	return l
}

type indent struct {
	size  int
	level int
}

func indentBy(size int, level int) *indent { return &indent{size: size, level: level} }

func (ind *indent) inc()                 { ind.level++ }
func (ind *indent) dec()                 { ind.level-- }
func (ind *indent) set(level int)        { ind.level = level }
func (ind *indent) by(delta int) *indent { return indentBy(ind.size, ind.level+delta) }
func (ind *indent) ent(s string) string {
	if ind.level < 0 {
		ind.level = 0
	}
	if ind.size < 2 {
		ind.size = 2
	}
	pos := ind.size * ind.level
	if strings.HasPrefix(s, "- ") {
		pos = pos - 2
	}
	return fmt.Sprintf("%*s%s", pos, "", s)
}
