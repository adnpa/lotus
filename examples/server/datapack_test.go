package znet_test

import (
	"fmt"
	"io"
	"net"
	"testing"
	"time"

	"github.com/adnpa/lotus/zpack"
)

func TestDataPack(t *testing.T) {

	listener, err := net.Listen("tcp", "127.0.0.1:7777")
	if err != nil {
		fmt.Println("lsiten err", err)
		return
	}

	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				fmt.Println("accept err", err)
			}

			go func() {
				dp := zpack.NewDatapack()
				for {

					// 读两次 第一次已知头长度 第二次从头里知道包长度
					headData := make([]byte, dp.GetHeadLen())
					_, err := io.ReadFull(conn, headData)
					if err != nil {
						fmt.Println("read head err")
						return
					}

					// 第一次读
					msgHead, err := dp.Unpack(headData)
					if err != nil {
						fmt.Println("unpack data err", err)
						return
					}

					//第二次读
					if msgHead.GetMsgLen() > 0 {
						//使用断言 因为Read需要用data指针
						msg := msgHead.(*zpack.Message)
						msg.Data = make([]byte, msg.DataLen)
						_, err := io.ReadFull(conn, msg.Data)
						if err != nil {
							fmt.Println("unpack data err", err)
							return
						}

						t.Log(msg.DataLen)
						t.Log(msg.Id)
						t.Log(msg.Data)
						fmt.Println("recv data: ", msg)
						time.Sleep(2 * time.Second)

					}

				}
			}()
		}

	}()

	// 模拟客户端
	conn, err := net.Dial("tcp", "127.0.0.1:7777")
	if err != nil {
		fmt.Println("dial err:", err)
		return
	}

	dp := zpack.NewDatapack()
	msg1 := &zpack.Message{
		Id:      1,
		DataLen: 4,
		Data:    []byte{'a', 'b', 'c', 'd'},
	}
	msg2 := &zpack.Message{
		Id:      2,
		DataLen: 4,
		Data:    []byte{'z', 'i', 'n', 'x'},
	}
	sendData1, err := dp.Pack(msg1)
	if err != nil {
		fmt.Println("pack err:", err)
		return
	}

	sendData2, err := dp.Pack(msg2)
	if err != nil {
		fmt.Println("pack err:", err)
		return
	}

	// 把两个包粘在一起
	sendData3 := append(sendData1, sendData2...)
	conn.Write(sendData3)

	// 客户端阻塞 否则会出现没有打印的情况
	// 因为最内层的协程没有足够的时间来执行t.Log语句，而主协程很快就会退出，从而导致日志消息丢失。
	// 1 使用同步等待
	// var wg sync.WaitGroup
	// wg.Add(1)
	// defer wg.Done()
	// 2 阻塞 select {}
	// 3 主线程休眠
	time.Sleep(3 * time.Second)
}
