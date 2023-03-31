package db

import (
	"testing"
	"time"
)

func Test_fixedConf(t *testing.T) {
	var c Conf
	c = fixedConf(c)
	d, _ := time.ParseDuration(c.MaxLife)
	t.Logf("conf %+v ,%v", c, d)
}
