# matrix
This program generates the Travis CI configuration file (`.travis.yml`) containing jobs to compile and integration test the [`github.com/ardnew/ft232h`](https://github.com/ardnew/ft232h) Go module for all supported platforms.

I was having trouble expressing the complex targets and dependencies using an implicit build matrix, and explicitly constructing all jobs seemed like a maintenance headache.

## Usage
This program takes no arguments and prints the complete YAML content to stdout.

For example, to replace the build matrix from this directory:

```sh
go run . > ../.travis.yml
```

## Generated Jobs
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
2. Install target Go version
3. Clone `git` repository
4. Compile native driver `libft232h.a`
5. Run module test suite via `go test -short`

> The `-short` flag is used because of the nature of this module. Being a device driver, the Travis CI build containers would need an actual FT232H device connected in order to run the full test suite. The reduced, `-short` suite acts as a smoke test that verifies linkage with the native driver `libft232h.a` as well as various utility types and methods. Omit the `-short` flag to perform full integration testing with a real FT232H device connected to the system.

In total, the module is tested against `×5` versions of Go on `×1` macOS and `×4` Linux architectures (`5×(4+1)` = `×25` build target configurations).

## YAML Verification
It can be a pain to spend so much time reconfiguring the generated `.travis.yml`, pushing upstream, wait for Travis to spawn jobs, and then find the builds all failing for a common configuration mistake (unrelated to what you're really trying to test). 

You can save a lot of time using [Travis CI's `travis-yml`](https://github.com/travis-ci/travis-yml) parser to validate the configuration before attempting to push the changes.

The easiest method I've found to use the package (i.e. to _not_ write any Ruby code) is to use the included Web API with `curl` as described in [their project README](https://github.com/travis-ci/travis-yml#web-api).

> At the time of this writing (`2020 Mar 13`), the `travis-yml` project was configured for Ruby 2.6.2. Be sure you have that version of Ruby installed, otherwise the software won't run. You can check required version in their project file `travis-yml/.ruby-version`. I was able to just modify this file's content to match my installation from `2.6.2` to `2.6.3` without issue, but you probably want to get their expected version for compatibility sake.

```sh
# download the project
git clone https://github.com/travis-ci/travis-yml.git
cd travis-yml
# fetch its dependencies
bundle install
# start the validation Web service
bundle exec rackup
```

This should start a local Web server listening on `tcp://127.0.0.1:9292`. You can then validate the generated `.travis.yml` with the following:
> You should also have the JSON command-line parser `jq` installed to make these outputs readable
```sh
# must be run from this directory
cd $GOPATH/src/github.com/ardnew/ft232h/matrix
# alias to parse the generated YAML as JSON
alias _parse='curl -sX POST --data-binary @<( go run . ) localhost:9292/v1/parse'
# alias to perform build matrix expansion on JSON parsed from generated YAML
alias _expand='curl -sX POST --data-binary @<( _parse ) localhost:9292/v1/expand'
# view the parsed JSON content
_parse | jq
# check if any errors exist (returns empty `[]` if no problems were found)
_parse | jq '.messages, .full_messages'
# view the expanded build matrix (all errors MUST be fixed before calling)
_expand | jq
```

Once verified, refer to [Usage](#usage) to install the configuration.

