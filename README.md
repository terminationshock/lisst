[![pipeline status](https://gitlab.mpcdf.mpg.de/tmelson/lisst/badges/master/pipeline.svg)](https://gitlab.mpcdf.mpg.de/tmelson/lisst/-/commits/master)

# lisst

This tool displays the output of a program as an interactive list allowing you to launch another command with input from lines of interest.

## Usage

The usage of *lisst* can best be explained with an example:

```bash
git log | lisst "[0-9a-f]{40}" git show
```

This command will pipe the output of `git log` into *lisst*. It will open an interactive list in which you can browse with the arrow keys.
Each commit hash matching the regular expression `[0-9a-f]{40}` will be highlighted in red. When you hit the enter key on a line containing
a commit hash, the command `git show <commit hash>` will be executed. When this command returns, the list will be shown again allowing you
to select the next line of interest.

*lisst* will accept all output piped into it and split it on line breaks. Each line is matched against the regular expression.
The first match in a line is highlighted. The enter key will trigger the upstream command only if the selected line contains a match.
An arbitrary number of command line arguments can be added to the command. The highlighted match in the selected line will be appended
to this list of arguments.

Further examples can be found in the output of `lisst --help`.

## Building

This program requires Go version 1.19 or later for building. You can build the self-contained executable yourself with `./build.sh`
or you can download the artifact of the latest commit from the `master` branch [here](https://gitlab.mpcdf.mpg.de/tmelson/lisst/-/jobs/artifacts/master/raw/lisst?job=build).

## Testing

Unit tests are executed during the build process. However, you can also run the integration test suite with `./test.sh`.

## Licenses of dependencies

The provided executable depends on several modules. Their licenses are listed [here](LICENSES_DEPENDENCIES.md).
