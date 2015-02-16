R2D2, Mozilla OpSec's IRC bot
=============================

Get all the dependencies, build with `go build *.go -o r2d2` and run with `./r2d2 -c r2d2.cfg`

Contributing
------------

r2d2 is a fairly monolithic but reliable robot. Add your code in a separate
file, and reference it the `handleRequest` function of `r2d2.go`.
