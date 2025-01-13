package util

import (
	"log"
	"os"
)

var Log = log.New(os.Stdout, "[simulator]", log.Lshortfile|log.Ldate|log.Ltime)
