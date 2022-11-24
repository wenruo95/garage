package config

import (
	"flag"
	"sync/atomic"
	"unsafe"

	"gopkg.in/ini.v1"
)

var gflag *Flag

type Flag struct {
	ConfigFile string
}

func init() {
	gflag = new(Flag)
	flag.StringVar(&gflag.ConfigFile, "conf", "./conf/app.ini", "-conf=./conf/app.ini")
}

func Parse() error {
	flag.Parse()

	if err := loadConfig(gflag.ConfigFile); err != nil {
		return err
	}

	return nil
}

func loadConfig(file string) error {
	config := &GarageConfig{}
	if err := ini.MapTo(config, file); err != nil {
		return err
	}

	setConfig(config)
	return nil
}
func setConfig(config *GarageConfig) {
	if gaconfig == nil {
		gaconfig = new(unsafe.Pointer)
	}
	atomic.StorePointer(gaconfig, unsafe.Pointer(config))
}
