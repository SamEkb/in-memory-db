package network

import (
	"net"

	"go.uber.org/zap"
	"in-memory-db/internal/database"
	"in-memory-db/internal/synchronization"
)

type TcpServer struct {
	listener  net.Listener
	logger    *zap.Logger
	semaphore *synchronization.Semaphore
	db        *database.Database
}

func NewServer(address string, logger *zap.Logger, db *database.Database) (*TcpServer, error) {
	l, err := net.Listen("tcp", address)
	if err != nil {
		logger.Error("Failed to listen on address", zap.String("address", address), zap.Error(err))
		return nil, err
	}

	server := &TcpServer{
		listener: l,
		logger:   logger,
		db:       db,
	}

	server.semaphore = synchronization.NewSemaphore(5)

	return server, nil
}

func (s *TcpServer) AcceptConnections() {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			s.logger.Error("Failed to accept connection", zap.Error(err))
			continue
		}

		s.logger.Info("New client connected", zap.String("remoteAddr", conn.RemoteAddr().String()))

		s.semaphore.Acquire()
		go func() {
			s.handleClient(conn)
			defer s.semaphore.Release()
		}()
	}
}

func (s *TcpServer) handleClient(conn net.Conn) {
	defer conn.Close()

	buf := make([]byte, 1024)
	for {
		request, err := conn.Read(buf)
		if err != nil {
			s.logger.Error("Error reading from client", zap.Error(err))
			return
		}

		query := string(buf[:request])
		res, err := s.db.HandleQuery(query)
		if err != nil {
			s.logger.Error("Failed to handle query", zap.String("query", query), zap.Error(err))
			continue
		}
		s.logger.Info("Received message", zap.String("message", res))
		_, _ = conn.Write([]byte(res))
	}
}

func (s *TcpServer) Close() {
	s.logger.Info("Closing server")
	_ = s.listener.Close()
}
