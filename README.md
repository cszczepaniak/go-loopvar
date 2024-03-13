# go-loopvar
A static analysis tool to find unnecessary loop variable captures.

### Motivation
As of Go 1.22, it's no longer necessary to capture loop variables for use in closures, to take their
address, etc. Now loop variables are scoped to single iterations rather than the entire loop.

This tool finds cases that are probably unnecessary captures of loop variables. If run with `-fix`,
it will also remove the unnecessary captures.
