# magefile-example

Example magefile that I have used for 2 projects

To use it in your own project, copy `magefile.go` from this directory to your own project
and change `EXE` to whatever your project is called.

## Usage

Set up:

```sh
go install github.com/magefile/mage
go install .
```

Then run `mage` in the root directory. This should show:

```
Targets:
  build                Builds an executable for this computer
  buildAll             Builds an executable for all supported platforms
  buildLinuxAmd64      Builds an executable for Linux AMD64
  buildLinuxArm64      Builds an executable for Linux ARM64
  buildMacAmd64        Builds an executable for Mac AMD64
  buildMacArm64        Builds an executable for Mac ARM64
  buildWindowsAmd64    Builds an executable for Windows AMD64
  check                Runs go vet and go fmt, and checks that they don't say anything
  checkRepoClean       Checks that the repo is clean
  ci                   Runs everything that a CI system might want to do
  clean                Cleans the bin directory
  debug                Builds a debug executable for this computer
  run                  Runs the program in debug mode without arguments
  runRelease           Runs the program in release mode without arguments
  test                 Runs go test in verbose mode and prettifies the output
```

Then you can e.g. run `mage buildAll`, which should say something like:

```
Running 'go build -o bin/windows-amd64/example.debug.exe -tags debug .' with env GOOS="windows", GOARCH="amd64"
Running 'go build -o bin/darwin-arm64/example -tags release .' with env GOOS="darwin", GOARCH="arm64"
Running 'go build -o bin/windows-amd64/example.exe -tags release .' with env GOARCH="amd64", GOOS="windows"
Running 'go build -o bin/darwin-amd64/example -tags release .' with env GOOS="darwin", GOARCH="amd64"
Running 'go build -o bin/linux-arm64/example -tags release .' with env GOOS="linux", GOARCH="arm64"
Running 'go build -o bin/linux-amd64/example -tags release .' with env GOOS="linux", GOARCH="amd64"
```

There is a Github Actions script in .github/workflows/build.yml which builds the example project
by running `mage ci` (this runs some sanity checks, runs the tests and builds for all platforms).
It also uploads the built artifacts.
