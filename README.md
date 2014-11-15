# Goneric

Generics for go - at last!

## Status

This isn't really working yet

## Use

Do not.

## Licence

You may not execute this software whilst developing anything another human may ever have to maintain.
Doing so constitutes a violation of basic human decency.

## How?

Goneric mounts a FUSE filesystem in `$GOROOT/src/goneric`.

You can then import like so: `import fileList "goneric/tuple/___os/File___uint"`.

When this file is read, goneric will look for `*.goneric` files in
`$GOROOT/src/tuple`, and interpret them using `text/template`.

goneric files look like normal go code, except they can use
`text/template` interpolation - e.g.
`{{ range .Imports }} {{ .Path }} {{ .Type }} {{ end }}`

Note that Path could be empty (if a builtin type is used).

## Dependencies

You need to have goneric running whenever you build your code.
Goneric depends on FUSE, so you'll need FUSE drivers for your platform.
Unfortunately windows is not yet supported (although DokanX is promising).

