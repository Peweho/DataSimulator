package util

import (
	"log"
	"os"
)

var Log = log.New(os.Stdout, "[FarmeEnvData]", log.Lshortfile|log.Ldate|log.Ltime)
