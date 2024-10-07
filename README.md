# rye-gioui

Gio UI library with Rye language. This is the first working version ... more documentation, examples and better bindings are yet to come.

## Build

To build rye-gio binary use the following commands.

```sh
./build

# of

go build -o bin/rye-gio
```

Or you can use Makefile to regenerate the binding and build it (needs some tweaking at this point for second example):

```sh
# runs go generate and build into the local .bin folder
make all

## Runs an example
make run
```

## Usage


```sh
## Run the Hellp example
bin/rye-gio examples/hello_gio.rye

## Run the Click example
bin/rye-gio examples/click_counter.rye
```


## Examples

![example render](./docs/hello.png)

![example render](./docs/click-counter.png)
