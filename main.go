package main

import (
	"context"

	"github.com/sirupsen/logrus"
)

func main() {
	ctx := context.Background()
	entry := logrus.WithContext(ctx).WithField("hello", "world")
	entry.Infof("hello")

	logrus.WithContext(ctx).Infof("hello2")
	logrus.Infof("hello")
}
