package znet

import (
	"fmt"
	"gozinx/utils"
	"gozinx/ziface"
	"net"
)

type Server struct {
	Name        string // 服务器名称
	IPVersion   string // 服务器绑定ip版本
	IP          string // 服务器绑定的ip
	Port        int    // 服务器监听的端口
	MsgHandler  ziface.IMsgHandler
	ConnMgr     ziface.IConnManager
	OnConnStart func(conn ziface.IConnection)
	OnConnStop  func(conn ziface.IConnection)
}

// 启动服务器
func (s *Server) Start() {
	fmt.Printf("[%s] Server Listener at IP: %s, Port %d\n", s.Name, s.IP, s.Port)

	go func() {
		// 开始消息队列
		s.MsgHandler.StartWorkerPool()

		// 1.获取一个tcp的addr
		addr, err := net.ResolveTCPAddr(s.IPVersion, fmt.Sprintf("%s:%d", s.IP, s.Port))
		if err != nil {
			fmt.Println("resolve tcp addr error: ", err)
			return
		}
		// 2.监听这个服务地址
		listenner, err := net.ListenTCP(s.IPVersion, addr)
		if err != nil {
			fmt.Printf("listen %s error,%s\n", s.IPVersion, err)
			return
		}
		// fmt.Printf("start zinx server success %s ,Listenning...\n", s.Name)
		var cid uint32 = 0

		// 3.阻塞等待客户端连接，处理客户端连接业务（读写）
		for {
			// 如果有客户端连接过来，阻塞会返回
			conn, err := listenner.AcceptTCP()
			if err != nil {
				fmt.Println("Accept err", err)
				continue
			}

			if s.ConnMgr.Len() >= utils.GlobalObject.MaxConn {
				// todo 给客户端相应超出最大一个错误包
				conn.Close()
				fmt.Println("too many connection")
				continue
			}

			// 将处理新链接业务方法和conn进行绑定
			dealConn := NewConnection(s, conn, cid, s.MsgHandler)
			cid++

			// 启动当前的链接业务处理
			go dealConn.Start()
		}
	}()
}

// 停止服务器
func (s *Server) Stop() {
	// 将一些服务器资源、状态或者一些已经开辟的链接信息 进行停止或回收
	s.ConnMgr.ClearConn()
}

// 运行服务器
func (s *Server) Serve() {
	s.Start() // 启动server的服务器功能

	// 做一些服务器启动后额外的业务

	// 阻塞状态
	select {}
}

func (s *Server) AddRouter(msgId uint32, router ziface.IRouter) {
	s.MsgHandler.AddRouter(msgId, router)
}

func (s *Server) GetConnMgr() ziface.IConnManager {
	return s.ConnMgr
}

// 初始化Server模块方法
func NewServer() ziface.IServer {
	s := &Server{
		Name:       utils.GlobalObject.Name,
		IPVersion:  "tcp4",
		IP:         utils.GlobalObject.Host,
		Port:       utils.GlobalObject.TcpPort,
		MsgHandler: NewMsgHandler(),
		ConnMgr:    NewConnManager(),
	}
	return s
}

func (s *Server) SetOnConnStart(hookFunc func(ziface.IConnection)) {
	s.OnConnStart = hookFunc
}

func (s *Server) SetOnConnStop(hookFun func(ziface.IConnection)) {
	s.OnConnStop = hookFun
}

func (s *Server) CallOnConnStart(conn ziface.IConnection) {
	if s.OnConnStart != nil {
		s.OnConnStart(conn)
	}
}

func (s *Server) CallOnConnStop(conn ziface.IConnection) {
	if s.OnConnStop != nil {
		s.OnConnStop(conn)
	}
}
