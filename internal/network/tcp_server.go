package network

import (
	"net"
	"time"

	"in-memory-db/internal/initialization"
	"in-memory-db/internal/synchronization"

	"go.uber.org/zap"
)

type TcpServer struct {
	listener  net.Listener
	semaphore synchronization.ISemaphore
	app       initialization.IApp
}

func NewServer(semaphore synchronization.ISemaphore, app initialization.IApp) (*TcpServer, error) {
	logger := app.GetLogger()

	address := app.GetConfig().Network.Address
	l, err := net.Listen("tcp", address)
	if err != nil {
		logger.Error("Failed to listen on address", zap.String("address", address), zap.Error(err))
		return nil, err
	}

	server := &TcpServer{
		listener:  l,
		semaphore: semaphore,
		app:       app,
	}

	return server, nil
}

func (s *TcpServer) AcceptConnections() {
	logger := s.app.GetLogger()
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			logger.Error("Failed to accept connection", zap.Error(err))
			continue
		}

		logger.Info("New client connected", zap.String("remoteAddr", conn.RemoteAddr().String()))

		duration := s.app.GetConfig().Network.IdleTimeout
		timeout := time.Now().Add(duration)

		err = conn.SetWriteDeadline(timeout)
		if err != nil {
			logger.Error("Failed to set write deadline", zap.Error(err))
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
	logger := s.app.GetLogger()
	defer func(conn net.Conn) {
		err := conn.Close()
		if err != nil {
			logger.Error("Failed to close connection", zap.Error(err))
		}
	}(conn)
	err := conn.SetReadDeadline(time.Now().Add(5 * time.Minute))
	if err != nil {
		logger.Error("Failed to set deadline", zap.Error(err))
		return
	}

	buf := make([]byte, s.app.GetConfig().Network.MaxMessageSize)
	for {
		request, err := conn.Read(buf)
		if err != nil {
			logger.Error("Error reading from client", zap.Error(err))
			continue
		}

		query := string(buf[:request])
		res, err := s.app.GetDataBase().HandleQuery(query)
		if err != nil {
			logger.Error("Failed to handle query", zap.String("query", query), zap.Error(err))
			_, _ = conn.Write([]byte("Error processing request\n"))
		}
		logger.Info("Received message", zap.String("message", res))
		_, _ = conn.Write([]byte(res))
	}
}

func (s *TcpServer) Close() error {
	s.app.GetLogger().Info("Closing server")
	return s.listener.Close()
}
