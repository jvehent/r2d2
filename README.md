R2D2, yet another IRC bot
=========================

**(but this one's mine!)**

If you're on linux/amd64, just grab the r2d2 binary from the repo. Otherwise,
get all the dependencies, build with `go build *.go -o r2d2` and run with `./r2d2 -c r2d2.cfg`

Contributing
------------

r2d2 is a fairly monolithic but reliable robot. Add your code in a separate
file, and reference it the `handleRequest` function of `r2d2.go`.
