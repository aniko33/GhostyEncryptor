package main

import (
	"os"
	"path"
)

func init() {
    exePath, err := os.Executable()

    if err != nil {
        panic(err)
    }

    bannerBytes, err := os.ReadFile(
        path.Join(
            path.Dir(
                exePath,
            ),
            "banner.txt",
        ),
    )

    if err != nil {
        panic(err)
    }

    println(string(bannerBytes))
}
