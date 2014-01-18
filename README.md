## Terry-Mao/goconf

`Terry-Mao/goconf` is an configuration file parse module.

## Requeriments
* Go 1.1 or higher

## Installation

Just pull `Terry-Mao/goconf` from github using `go get`:

```sh
# download the code
$ go get -u github.com/Terry-Mao/goconf
```

## Usage

```go
package main

import (
    "fmt"
    "github.com/Terry-Mao/goconf"
)

func main() {
    conf, err := goconf.New("./examples/conf_test.txt")
    if err != nil {
        panic(err)
    }

    id, err := conf.Int("id")
    if err != nil {
        panic(err)
    }
    
    fmt.Println(id)
}
```

## Documentation

Read the `Terry-Mao/goconf` documentation from a terminal

```go
$ go doc github.com/Terry-Mao/goconf
```

Alternatively, you can [goconf](http://go.pkgdoc.org/github.com/Terry-Mao/goconf) online.
