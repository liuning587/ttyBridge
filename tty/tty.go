package tty

import (
	"io"
	"log"
	"sync"

	"github.com/jacobsa/go-serial/serial"
)

// TtyConn tty connect
type TtyConn struct {
	sync.Mutex
	recvLock sync.Mutex
	serName  string             //串口名
	f        io.ReadWriteCloser //串口读写句柄
}

// NewTtyConn new tty connect
func NewTtyConn(serName string, baudrate int) (*TtyConn, error) {
	options := serial.OpenOptions{
		PortName:              serName,
		BaudRate:              115200,
		DataBits:              8,
		StopBits:              1,
		MinimumReadSize:       0,
		InterCharacterTimeout: 1000,
		ParityMode:            serial.PARITY_NONE,
	}
	f, err := serial.Open(options)
	if err != nil {
		return nil, err
	}

	conn := &TtyConn{
		serName: serName,
		f:       f,
	}

	return conn, nil
}

// SendData send data
// fixme 0D 0A
func (t *TtyConn) SendData(buf []byte) error {
	log.Printf("TTY TX: % X", buf)
	t.Lock()
	defer t.Unlock()

	_, err := t.f.Write(buf)
	if err != nil {
		log.Println("Write err:", err)
		return err
	}
	return nil
}

// RecvData send data
func (t *TtyConn) RecvData() ([]byte, error) {
	t.recvLock.Lock()
	defer t.recvLock.Unlock()

	buf := make([]byte, 2048)
	n, err := t.f.Read(buf)
	if err != nil {
		if err != io.EOF {
			log.Println("Error reading from serial port: ", err)
			return nil, err
		}
		return []byte{}, nil
	}

	if n > 0 {
		log.Printf("TTY RX: % X", buf[:n])
	}

	return buf[:n], nil
}

// Disconnect close conna nd serial
func (t *TtyConn) Disconnect() error {
	if t.f != nil {
		t.f.Close()
	}
	return nil
}
