package main

import (
	"fmt"
	"net"
	"time"

	"github.com/adnpa/lotus/zconf"
	"github.com/adnpa/lotus/zpack"
)

func main() {
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", zconf.GGlobalObj.Host, zconf.GGlobalObj.TcpPort))
	if err != nil {
		fmt.Println("client conntect err: ", err)
		return
	}

	for {

		dp := zpack.NewDatapack()
		msg := zpack.NewMessage(1, []byte("你好"))

		sendDatapack, err := dp.Pack(msg)
		if err != nil {
			fmt.Println("pack msg err:", err)
			return
		}

		_, err = conn.Write(sendDatapack)
		if err != nil {
			fmt.Println("cnt send msg err:", err)
			return
		}

		buf := make([]byte, 512)
		n, err := conn.Read(buf)
		if err != nil {
			fmt.Println("read buf err")
			return
		}

		fmt.Printf("server call back: %s, cnt=%d \n", buf, n)
		time.Sleep(1 * time.Second)
	}

}
