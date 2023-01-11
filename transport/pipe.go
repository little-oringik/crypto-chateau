package transport

import (
	"net"
	"time"
)

type Pipe struct {
	tcpConn            net.Conn
	reservedData       []byte
	futurePacketLength uint16
	cfg                connCfg
}

func NewPipe(tcpConn net.Conn) *Pipe {
	return &Pipe{
		tcpConn:      tcpConn,
		reservedData: make([]byte, 0, 1024),
	}
}

func (cn *Pipe) Write(p []byte) (int, error) {
	dataWithLength := make([]byte, 0, len(p)+2)
	convertedLength := uint16(len(p))
	dataWithLength = append(dataWithLength, byte(convertedLength), byte(convertedLength>>8))
	dataWithLength = append(dataWithLength, p...)
	n, err := cn.tcpConn.Write(dataWithLength)
	return n, err
}

type PipeReadCfg struct {
	BufSize int
}

func (cn *Pipe) Read(cfg PipeReadCfg) ([]byte, error) {
	if cfg.BufSize == 0 {
		cfg.BufSize = 1024
	}

	fullMessage, err := GetFullMessage(cn.tcpConn, cfg.BufSize+2, cn.reservedData, cn.futurePacketLength)
	if err != nil {
		return nil, err
	}

	cn.futurePacketLength = fullMessage.gotFuturePacketLength
	cn.reservedData = fullMessage.gotReservedData

	return fullMessage.msg, nil
}

func (cn *Pipe) GetConn() net.Conn {
	return cn.tcpConn
}

func (cn *Pipe) SetConn(conn net.Conn) {
	cn.tcpConn = conn
}

func (cn *Pipe) CloseConn() error {
	return cn.tcpConn.Close()
}

func (cn *Pipe) LocalAddr() net.Addr {
	return cn.tcpConn.LocalAddr()
}

func (cn *Pipe) RemoteAddr() net.Addr {
	return cn.tcpConn.RemoteAddr()
}

func (cn *Pipe) SetDeadline(t time.Time) error {
	return cn.tcpConn.SetDeadline(t)
}

func (cn *Pipe) SetReadDeadline(t time.Time) error {
	return cn.tcpConn.SetReadDeadline(t)
}

func (cn *Pipe) SetWriteDeadline(t time.Time) error {
	return cn.tcpConn.SetWriteDeadline(t)
}
