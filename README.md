# boltstore  :convenience_store:

[![GoDoc](https://godoc.org/github.com/schollz/boltstore?status.svg)](https://godoc.org/github.com/schollz/boltstore)

*boltstore* is a Go-library for a simple thread-safe in-memory key-store with persistent BoltDB backend. 

## Usage

First, install the library using:

```
go get -u -v github.com/schollz/boltstore
```

Then you can add it to your program. Check out the examples, or see below for basic usage:

```golang
ks := boltstore.Open("mystore.db")

// set a key to any object you want
type Human struct {
  Name   string
  Height float64
}
err := ks.Set("human:1", Human{"Dante", 5.4})
if err != nil {
  panic(err)
}

// get the data back via an interface
var human Human
err = ks2.Get("human:1", &human)
if err != nil {
  panic(err) // returns error if key doesn't exist
}
fmt.Println(human.Name) // Prints 'Dante'
```

# Benchmark
```
$ go test -bench=.
goos: linux
goarch: amd64
pkg: github.com/schollz/boltstore
BenchmarkGet-4            500000              2474 ns/op
BenchmarkSet-4               300           7375338 ns/op
```


# License

MIT