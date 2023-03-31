package db

const DefaultPageSize = 20
const DefaultKey = "default"

type Conf struct {
	User     string
	Pass     string
	Server   string
	Database string
	Name     string
	MaxIdle  uint32
	MaxOpen  uint32
	MaxLife  string
}
