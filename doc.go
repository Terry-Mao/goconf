/*
Package goconf provides configuraton read and write implementations.

Examples:
    package main

    import (
        "fmt"
        "github.com/Terry-Mao/goconf"
    )

    func main() {
        conf := goconf.New()
        if err := conf.Parse("./examples/conf_test.txt"); err != nil {
            panic(err)
        }
        core := conf.Get("core")
        if core == nil {
            panic("no core section")
        }
        id, err := core.Int("id")
        if err != nil {
            panic(err)
        }
        fmt.Println(id)
    }
*/
package goconf
