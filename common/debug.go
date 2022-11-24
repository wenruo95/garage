package common

import (
	"encoding/json"

	"github.com/sirupsen/logrus"
)

func JsonStr(i interface{}) string {
	if i == nil {
		return ""
	}
	buff, err := json.Marshal(i)
	if err != nil {
		logrus.Warnf("json marshal error:%v", err)
		return ""
	}
	return string(buff)
}
