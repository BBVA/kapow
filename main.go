package main

import (
    "fmt"

    b "github.com/BBVA/kapow/pkg/banner"
)

func main() {
    ban := b.Banner("0.1.0")
    fmt.Println(ban)
}
