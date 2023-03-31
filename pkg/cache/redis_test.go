package cache

import (
	"context"
	"testing"
	"time"
)

func TestHSetWithExpireContext(t *testing.T) {
	NewPools([]RedisConf{
		{
			Server: "192.168.88.253:6379",
		},
	})
	val, err := SetNXExpireContext(context.Background(), "advertisement_2023030717053207000827", `{"spreadId":1}`, time.Minute*15)
	t.Logf("val %v ,err %v", val, err)
}
