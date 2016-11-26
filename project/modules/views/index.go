package helpers

import (
    "github.com/nogio/noggo"
)

func init() {

    noggo.View.Helper("test", func() string {
        return "test"
    })

}
