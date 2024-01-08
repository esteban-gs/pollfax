package main

import (
	"pollfax/internal/ingest"
	"pollfax/internal/util"
)

func main() {
	util.LoadAppEnv()
	ingest.Run()
}
