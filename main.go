package main

import (
    "fmt"

    b "github.com/BBVA/kapow/pkg/banner"
)

func main() {
    ban := b.Banner()
    fmt.Println(ban)
}
