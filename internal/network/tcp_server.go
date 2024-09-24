package network

import (
	"fmt"
	"net"
	"time"

	"go.uber.org/zap"
	"in-memory-db/internal/initialization"
	"in-memory-db/internal/synchronization"
)

type TcpServer struct {
	listener  net.Listener
	semaphore *synchronization.Semaphore
	init      *initialization.ServerInitializer
}

func NewServer() (*TcpServer, error) {
	init, err := initialization.InitializeServer()
	if err != nil {
		panic(fmt.Sprintf("Initialization error: %v", err))
	}

	address := init.Config.Network.Address
	logger := init.Logger
	maxConnections := init.Config.Network.MaxConnection

	l, err := net.Listen("tcp", address)
	if err != nil {
		logger.Error("Failed to listen on address", zap.String("address", address), zap.Error(err))
		return nil, err
	}

	server := &TcpServer{
		listener: l,
		init:     init,
	}

	server.semaphore = synchronization.NewSemaphore(maxConnections)

	return server, nil
}

func (s *TcpServer) AcceptConnections() {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			s.init.Logger.Error("Failed to accept connection", zap.Error(err))
			continue
		}

		s.init.Logger.Info("New client connected", zap.String("remoteAddr", conn.RemoteAddr().String()))

		//TODO create time parser
		//conn.SetWriteDeadline(s.init.Config.Network.IdleTimeout)
		s.semaphore.Acquire()
		go func() {
			s.handleClient(conn)
			defer s.semaphore.Release()
		}()
	}
}

func (s *TcpServer) handleClient(conn net.Conn) {
	defer conn.Close()
	err := conn.SetReadDeadline(time.Now().Add(5 * time.Minute))
	if err != nil {

	}

	//TODO create message size parser
	//messageSize := s.init.Config.Network.MaxMessageSize
	buf := make([]byte, 1024)
	for {
		request, err := conn.Read(buf)
		if err != nil {
			s.init.Logger.Error("Error reading from client", zap.Error(err))
			return
		}

		query := string(buf[:request])
		res, err := s.init.DB.HandleQuery(query)
		if err != nil {
			s.init.Logger.Error("Failed to handle query", zap.String("query", query), zap.Error(err))
			continue
		}
		s.init.Logger.Info("Received message", zap.String("message", res))
		_, _ = conn.Write([]byte(res))
	}
}

func (s *TcpServer) Close() {
	s.init.Logger.Info("Closing server")
	_ = s.listener.Close()
}
