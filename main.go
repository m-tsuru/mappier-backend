package main

import (
	"github.com/m-tsuru/mappier-backend/lib/db"
	"github.com/m-tsuru/mappier-backend/lib/handler"
)

func main() {
	db.Setup()
	handler.Setup()
}
