package memcache

import (
	"log"
	"io/ioutil"
)

type StdLogger interface {
	Print(v ...interface{})
	Printf(format string, v ...interface{})
	Println(v ...interface{})
}

var Logger StdLogger = log.New(ioutil.Discard, "[BeansDB-Client] ", log.LstdFlags)
