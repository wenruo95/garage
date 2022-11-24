package config

import (
	"sync/atomic"
	"unsafe"
)

var gaconfig *unsafe.Pointer

func Get() *GarageConfig {
	if gaconfig != nil {
		if p := atomic.LoadPointer(gaconfig); p != nil {
			return (*GarageConfig)(p)
		}
	}
	return nil
}

type GarageConfig struct {
	Local *LocalConfig
	TCP   *TCPConfig
	Log   *LogConfig
}

type LocalConfig struct {
	PodIP   string
	PodName string
	Cluster string
}

type TCPConfig struct {
	Port string
}

type LogConfig struct {
	File string
}
