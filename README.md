# VMmanager 6 API for GoLang

## Install

> go get -u github.com/mayerdev/go-vm6-api

## Use

```go
package main

import "github.com/mayerdev/go-vm6-api"

func main() {
    api := utils.NewVm6("https://vm.local", "admin@vm.local", "my_password")
    err := api.Login()

    if err != nil {
        panic(err)
    }

    res, err := vm_api.Send("GET", "v3", "vm", "host", nil)

    if err != nil {
        panic(err)
    }

    fmt.Println(res)
}
```