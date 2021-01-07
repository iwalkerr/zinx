package znet

import (
	"fmt"
	"io"
	"net"
	"testing"
)

func TestDataPack(t *testing.T) {
	listenner, err := net.Listen("tcp", "127.0.0.1:8999")
	if err != nil {
		fmt.Println("server listen err: ", err)
		return
	}

	go func() {
		for {
			conn, err := listenner.Accept()
			if err != nil {
				fmt.Println("server accept error", err)
			}
			go func(conn net.Conn) {
				pd := NewDataPack()
				for {
					headData := make([]byte, int(pd.GetHeadLen()))
					_, err := io.ReadFull(conn, headData)
					if err != nil {
						fmt.Println("read head error")
						break
					}

					msgHead, err := pd.Unpack(headData)
					if err != nil {
						fmt.Println("server unpacke err ", err)
						return
					}

					if msgHead.GetMsgLen() > 0 {
						msg := msgHead.(*Message)
						msg.Data = make([]byte, msg.GetMsgLen())

						_, err := io.ReadFull(conn, msg.Data)
						if err != nil {
							fmt.Println("server unpack data err ", err)
							return
						}

						fmt.Printf("-----> Recv data: %s\n", string(msg.Data))
					}
				}

			}(conn)
		}
	}()

	// 模拟客户端
	conn, err := net.Dial("tcp", "127.0.0.1:8999")
	if err != nil {
		fmt.Println("client dial err: ", err)
		return
	}

	dp := NewDataPack()

	msg1 := &Message{
		Id:      1,
		DataLen: 4,
		Data:    []byte{'z', 'i', 'n', 'x'},
	}
	sendData1, err := dp.Pack(msg1)
	if err != nil {
		fmt.Println("client pack msg1 error", err)
		return
	}

	msg2 := &Message{
		Id:      2,
		DataLen: 7,
		Data:    []byte{'n', 'i', 'h', 'a', 'o', '!', '!'},
	}
	sendData2, err := dp.Pack(msg2)
	if err != nil {
		fmt.Println("client pack msg2 error", err)
		return
	}

	// 两个包拼接在一起
	sendData1 = append(sendData1, sendData2...)

	_, _ = conn.Write(sendData1)

	select {}
}
