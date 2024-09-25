package network

import (
	"net"
	"time"

	"go.uber.org/zap"
	"in-memory-db/internal/initialization"
	"in-memory-db/internal/synchronization"
	"in-memory-db/internal/utils"
)

type TcpServer struct {
	listener  net.Listener
	semaphore *synchronization.Semaphore
	init      *initialization.ServerInitializer
}

func NewServer() (*TcpServer, error) {
	init, err := initialization.InitializeServer()
	if err != nil {
		return &TcpServer{}, err
	}

	logger := init.Logger

	address := init.Network.Address
	l, err := net.Listen("tcp", address)
	if err != nil {
		logger.Error("Failed to listen on address", zap.String("address", address), zap.Error(err))
		return nil, err
	}

	server := &TcpServer{
		listener: l,
		init:     init,
	}

	server.semaphore = synchronization.NewSemaphore(init.Network.MaxConnections)

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

		timeout, err := utils.ParseTime(s.init.Network.IdleTimeout)
		if err != nil {
			s.init.Logger.Error("Failed to parse idle timeout", zap.Error(err))
			continue
		}

		err = conn.SetWriteDeadline(timeout)
		if err != nil {
			s.init.Logger.Error("Failed to set write deadline", zap.Error(err))
			continue
		}

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
		s.init.Logger.Error("Failed to set deadline", zap.Error(err))
		return
	}

	messageSize := s.init.Network.MaxMessageSize
	size, err := utils.ParseSize(messageSize)
	if err != nil {
		s.init.Logger.Error("Failed to parse message size", zap.Error(err))
		return
	}

	buf := make([]byte, size)
	for {
		request, err := conn.Read(buf)
		if err != nil {
			s.init.Logger.Error("Error reading from client", zap.Error(err))
			continue
		}

		query := string(buf[:request])
		res, err := s.init.DB.HandleQuery(query)
		if err != nil {
			s.init.Logger.Error("Failed to handle query", zap.String("query", query), zap.Error(err))
			_, _ = conn.Write([]byte("Error processing request\n"))
		}
		s.init.Logger.Info("Received message", zap.String("message", res))
		_, _ = conn.Write([]byte(res))
	}
}

func (s *TcpServer) Close() {
	s.init.Logger.Info("Closing server")
	_ = s.listener.Close()
}
