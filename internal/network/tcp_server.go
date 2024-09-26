package network

import (
	"fmt"
	"net"
	"time"

	"in-memory-db/internal/initialization"
	"in-memory-db/internal/synchronization"

	"go.uber.org/zap"
)

type TcpServer struct {
	listener  net.Listener
	semaphore *synchronization.Semaphore
	app       *initialization.App
}

func NewServer() (*TcpServer, error) {
	app, err := initialization.NewApp()
	if err != nil {
		return nil, fmt.Errorf("failed to initialize app: %v", err)
	}

	logger := app.Logger

	address := app.Config.Network.Address
	l, err := net.Listen("tcp", address)
	if err != nil {
		logger.Error("Failed to listen on address", zap.String("address", address), zap.Error(err))
		return nil, err
	}

	server := &TcpServer{
		listener: l,
		app:      app,
	}

	server.semaphore = synchronization.NewSemaphore(app.Config.Network.MaxConnections)

	return server, nil
}

func (s *TcpServer) AcceptConnections() {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			s.app.Logger.Error("Failed to accept connection", zap.Error(err))
			continue
		}

		s.app.Logger.Info("New client connected", zap.String("remoteAddr", conn.RemoteAddr().String()))

		duration := s.app.Config.Network.IdleTimeout
		timeout := time.Now().Add(duration)

		err = conn.SetWriteDeadline(timeout)
		if err != nil {
			s.app.Logger.Error("Failed to set write deadline", zap.Error(err))
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
		s.app.Logger.Error("Failed to set deadline", zap.Error(err))
		return
	}

	buf := make([]byte, s.app.Config.Network.MaxMessageSize)
	for {
		request, err := conn.Read(buf)
		if err != nil {
			s.app.Logger.Error("Error reading from client", zap.Error(err))
			continue
		}

		query := string(buf[:request])
		res, err := s.app.DB.HandleQuery(query)
		if err != nil {
			s.app.Logger.Error("Failed to handle query", zap.String("query", query), zap.Error(err))
			_, _ = conn.Write([]byte("Error processing request\n"))
		}
		s.app.Logger.Info("Received message", zap.String("message", res))
		_, _ = conn.Write([]byte(res))
	}
}

func (s *TcpServer) Close() error {
	s.app.Logger.Info("Closing server")
	return s.listener.Close()
}
