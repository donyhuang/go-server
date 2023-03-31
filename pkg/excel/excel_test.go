package excel

import (
	"os"
	"testing"
)

type excelData struct {
	Name    string `excel:"名字"`
	Age     int    `excel:"年龄"`
	Address string `excel:"地址"`
	test    string `excel:"测试"`
}

func TestSaveStructToBuff(t *testing.T) {
	var rows = []excelData{
		{
			Name: "dony",
			Age:  18,
			test: "t1",
		},
		{
			Name:    "pyf",
			Age:     17,
			Address: "西乡",
			test:    "t2",
		},
		{
			Name:    "bg",
			Age:     16,
			Address: "ct",
			test:    "t3",
		},
	}
	buf, _ := SaveStructToBuff(rows)
	f, _ := os.Create("test.xlsx")
	f.Write(buf)
}
