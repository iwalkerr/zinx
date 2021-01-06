package znet

import (
	"fmt"
	"gozinx/ziface"
	"net"
)

type Server struct {
	Name      string // 服务器名称
	IPVersion string // 服务器绑定ip版本
	IP        string // 服务器绑定的ip
	Port      int    // 服务器监听的端口
}

// 启动服务器
func (s *Server) Start() {
	fmt.Printf("[start] Server Listener at IP:%s, Port %d,is starting\n", s.IP, s.Port)

	go func() {
		// 1.获取一个tcp的addr
		addr, err := net.ResolveTCPAddr(s.IPVersion, fmt.Sprintf("%s:%d", s.IP, s.Port))
		if err != nil {
			fmt.Println("resolve tcp addr error: ", err)
			return
		}
		// 2.监听这个服务地址
		listenner, err := net.ListenTCP(s.IPVersion, addr)
		if err != nil {
			fmt.Println("listen ", s.IPVersion, " err ", err)
			return
		}
		fmt.Println("start zinx server success, ", s.Name, " Listenning...")

		// 3.阻塞等待客户端连接，处理客户端连接业务（读写）
		for {
			// 如果有客户端连接过来，阻塞会返回
			conn, err := listenner.AcceptTCP()
			if err != nil {
				fmt.Println("Accept err", err)
				continue
			}

			// 已经与客户端建立连接，做一些业务,最大512字节
			go func() {
				for {
					buf := make([]byte, 512)
					cnt, err := conn.Read(buf)
					if err != nil {
						fmt.Println("recv buf err ", err)
						continue
					}

					fmt.Printf("recv client buf %s, cnt %d\n", buf, cnt)

					// 回显功能
					if _, err := conn.Write(buf[:cnt]); err != nil {
						fmt.Println("write back buf err ", err)
						continue
					}
				}
			}()
		}
	}()
}

// 停止服务器
func (s *Server) Stop() {
	// TODO 将一些服务器资源、状态或者一些已经开辟的链接信息 进行停止或回收

}

// 运行服务器
func (s *Server) Serve() {
	s.Start() // 启动server的服务器功能

	// 做一些服务器启动后额外的业务

	// 阻塞状态
	select {}
}

// 初始化Server模块方法
func NewServer(name string) ziface.IServer {
	s := &Server{
		Name:      name,
		IPVersion: "tcp4",
		IP:        "0.0.0.0",
		Port:      8999,
	}
	return s
}
