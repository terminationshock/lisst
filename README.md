[![pipeline status](https://gitlab.mpcdf.mpg.de/tmelson/lisst/badges/master/pipeline.svg)](https://gitlab.mpcdf.mpg.de/tmelson/lisst/-/pipelines)
[![download](https://img.shields.io/badge/download-executable-success)](https://gitlab.mpcdf.mpg.de/tmelson/lisst/-/jobs/artifacts/master/raw/lisst?job=build)
[![license](https://img.shields.io/badge/license-MIT-informational)](LICENSE)
[![Go](https://img.shields.io/badge/Go-1.19-informational)](https://go.dev/dl/)

# lisst

This tool displays the output of a program as an interactive list allowing you to launch another command for a selected line of interest.

![Screenshot](lisst.gif)

## Usage

The usage of *lisst* can best be explained with an example:

```bash
git log --oneline | lisst "[0-9a-f]{7,40}" git show
```

This command will pipe the output of `git log` into *lisst*. It will open an interactive list which you can browse with the arrow keys.
Each commit hash matching the regular expression `[0-9a-f]{7,40}` will be highlighted. When you select a line containing a commit hash
and hit the enter key, the command `git show <commit hash>` will be executed. When this command returns, the list will be shown again allowing you
to select the next line of interest.

*lisst* accepts all non-empty output piped into it and splits it on line breaks. Each line is matched against the given regular expression.
The first match within a line is highlighted. The enter key triggers the upstream command only if the selected line contains a match.
An arbitrary number of command line arguments can be added to the command. The highlighted match in the selected line is appended
to this list of arguments. The status bar at the bottom displays the command that is about to be executed when the enter key is pressed.

The following example demonstrates how to efficiently edit many files:

```bash
grep -r func | lisst "^(.*?):" vi
```

All occurrences of `func` in all files in the current directory and all sub-directories will be displayed as a list with each file name being highlighted.
When you select a certain line and press the enter key, the editor `vi` will be launched and you can edit the file as usual. When you close the editor,
the list will be visible again allowing you to edit the next file.

You can use human-readable keywords for frequently used patterns. In the screencast shown above, for example,
the keyword `--git-commit-hash` is used for convenience instead of the actual regular expression.
See [`lisst --help`](https://gitlab.mpcdf.mpg.de/tmelson/lisst/-/jobs/artifacts/master/raw/lisst-help.txt?job=build)
for a complete list of supported keywords and more useful examples.

## Building

This program requires Go version 1.19 or later for building. You can build the self-contained executable yourself with `./build.sh`
or you can download the artifact built from the latest commit of the `master` branch [here](https://gitlab.mpcdf.mpg.de/tmelson/lisst/-/jobs/artifacts/master/raw/lisst?job=build).

## Testing

Unit tests are executed during the build process. The integration test suite can be executed with `./test.sh`.

## License

*lisst* is licensed under the [MIT License](LICENSE).

## Licenses of dependencies

*lisst* depends on several other Go packages. Their licenses are listed [here](LICENSES_DEPENDENCIES.md).
