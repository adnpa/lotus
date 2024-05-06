package zconf

import (
	"encoding/json"
	"os"

	"github.com/adnpa/lotus/ziface"
)

type GlobalObj struct {
	// server
	TcpServer ziface.IServer
	Host      string
	TcpPort   int
	Name      string

	Version          string
	MaxConn          int
	MaxPackageSize   uint32
	WorkerPoolSize   uint32
	MaxWorkerTaskLen uint32
}

var GGlobalObj *GlobalObj

func (g *GlobalObj) Reload() {
	data, err := os.ReadFile("conf/lotus.json")
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(data, &GGlobalObj)
	if err != nil {
		panic(err)
	}
}

func init() {
	GGlobalObj = &GlobalObj{
		Name:             "github.com/adnpa/lotusServerApp",
		Version:          "v0.1",
		TcpPort:          8999,
		Host:             "0.0.0.0",
		MaxConn:          1000,
		MaxPackageSize:   4096,
		WorkerPoolSize:   10,
		MaxWorkerTaskLen: 1024,
	}
	GGlobalObj.Reload()
}

// func GetGlobalObj() *GlobalObj {
// 	return globalObj
// }
