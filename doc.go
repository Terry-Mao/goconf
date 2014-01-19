/*
Package goconf provides configuraton read and write implementations.      
Examples:
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
*/
package goconf
