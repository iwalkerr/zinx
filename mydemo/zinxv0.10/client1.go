package main

import (
	"fmt"
	"gozinx/znet"
	"io"
	"net"
	"time"
)

// 模拟客户端
func main() {
	time.Sleep(1 * time.Second)

	conn, err := net.Dial("tcp", "127.0.0.1:8999")
	if err != nil {
		fmt.Println("client start err, exit!")
		return
	}

	for {
		dp := znet.NewDataPack()
		binaryMsg, err := dp.Pack(znet.NewMsgPackege(1, []byte("zonx v0.8 client1 test")))
		if err != nil {
			fmt.Println("pack error ", err)
			return
		}

		if _, err = conn.Write(binaryMsg); err != nil {
			fmt.Println("write error ", err)
			return
		}

		binaryHead := make([]byte, dp.GetHeadLen())
		if _, err := io.ReadFull(conn, binaryHead); err != nil {
			fmt.Println("read head error", err)
			break
		}

		msgHead, err := dp.Unpack(binaryHead)
		if err != nil {
			fmt.Println("client unpack msgHead error", err)
			break
		}

		if msgHead.GetMsgLen() > 0 {
			msg := msgHead.(*znet.Message)
			msg.Data = make([]byte, msg.GetMsgLen())

			if _, err := io.ReadFull(conn, msg.Data); err != nil {
				fmt.Println("read msg data error", err)
				return
			}

			fmt.Printf("Recv Server data %s\n", string(msg.Data))
		}

		// cpu阻塞
		time.Sleep(1 * time.Second)
	}
}
