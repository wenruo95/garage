package main

import (
	"context"

	"github.com/sirupsen/logrus"
	"github.com/wenruo95/garage/common"
	"github.com/wenruo95/garage/core/config"
)

func main() {
	logger := logrus.WithContext(context.Background()).WithField("cmd", "start_process")
	if err := config.Parse(); err != nil {
		logger.Errorf("load config error:%v", err)
		return
	}

	logger.Infof("config:%v", common.JsonStr(config.Get()))
}
