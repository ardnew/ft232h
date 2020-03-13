# matrix
This program generates the Travis CI configuration file (`.travis.yml`) containing jobs to compile and integration test the [`github.com/ardnew/ft232h`](https://github.com/ardnew/ft232h) Go module for all supported platforms. 

I was having trouble expressing the complex targets and dependencies using an implicit build matrix, and explicitly constructing all jobs seemed like a maintenance headache. 

### Usage
This program takes no arguments and prints the complete YAML content to stdout.

For example, to replace the build matrix from this directory:

```sh
go run . > ../.travis.yml 
```

### Generated Jobs
The following table lists all platforms that are automatically verified by integration test for every release:

||`x86` (32-bit)|`x86_64` (64-bit)|`ARM` (32-bit)|`ARMv8` (64-bit)|
|:---:|:---:|:---:|:---:|:---:|
|Linux|✔|✔|✔|✔|
|macOS||✔||||

Additionally, the following Go versions (`x` being most recent patch) are tested on each platform:
- Go `v1.11.x`
- Go `v1.12.x`
- Go `v1.13.x`
- Go `v1.14.x`
- Go `master`

The integration test for each platform (and each supported version of Go) performs the following:
1. Install new `gcc` C compiler
2. Install Go
3. Clone `git` repository
4. Compile native driver `libft232h.a`
5. Run module test suite via `go test`

In total, the module is tested against x5 versions of Go for x4 Linux and x1 macOS architectures (5x(4+1)=25 build target configurations). 
