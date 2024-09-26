package network

import (
	"net"
	"time"
)

type mockListener struct {
	acceptConn net.Conn
	acceptErr  error
	closeFunc  func() error
}

func (m *mockListener) Accept() (net.Conn, error) {
	return m.acceptConn, m.acceptErr
}

func (m *mockListener) Close() error {
	return m.closeFunc()
}

func (m *mockListener) Addr() net.Addr {
	return &net.TCPAddr{IP: net.ParseIP("127.0.0.1"), Port: 8080}
}

type mockConn struct {
	writeFunc func([]byte) (int, error)
	readFunc  func([]byte) (int, error)
	closeFunc func() error
}

func (m *mockConn) Write(b []byte) (int, error) {
	return m.writeFunc(b)
}

func (m *mockConn) Read(b []byte) (int, error) {
	return m.readFunc(b)
}

func (m *mockConn) Close() error                       { return m.closeFunc() }
func (m *mockConn) LocalAddr() net.Addr                { return nil }
func (m *mockConn) RemoteAddr() net.Addr               { return nil }
func (m *mockConn) SetDeadline(t time.Time) error      { return nil }
func (m *mockConn) SetReadDeadline(t time.Time) error  { return nil }
func (m *mockConn) SetWriteDeadline(t time.Time) error { return nil }
