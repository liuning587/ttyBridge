package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"sync"
	"time"
	"ttyBridge/tty"
)

var (
	h         bool
	enableLog bool
	ttyName   string
	port      int
	baudrate  int
	lock      sync.Mutex
	connID    int
	connMap   map[int]net.Conn
)

func init() {
	flag.BoolVar(&h, "h", false, "this help")
	flag.BoolVar(&enableLog, "log", false, "enable log")

	flag.StringVar(&ttyName, "tty", "/tmp/dev/usb_exm3_1", "tty device name")
	flag.IntVar(&baudrate, "baudrate", 9600, "baudrate")
	flag.IntVar(&port, "port", 40084, "tcp server port")
}

func main() {
	log.Println("ttyBridge v1.0.0")
	flag.Parse()
	if h {
		flag.Usage()
		return
	}

	connMap = make(map[int]net.Conn)
	listener, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", port))
	if err != nil {
		log.Println("listen err:", err)
		return
	}
	defer listener.Close()
	log.Printf("listen %d OK!\n", port)

	t, err := tty.NewTtyConn(ttyName, baudrate)
	if err != nil {
		log.Println("TTY open "+ttyName+" err:", err)
		return
	}
	defer t.Disconnect()

	go func() {
		for {
			buf, err := t.RecvData()
			if err != nil {
				time.Sleep(time.Second)
				continue
			}
			if len(buf) == 0 {
				continue
			}
			log.Printf("TCP TX: % X", buf)
			lock.Lock()
			for _, conn := range connMap {
				_, err = conn.Write(buf)
				if err != nil {
					log.Println("send to tcp err:", err)
				}
			}
			lock.Unlock()
		}
	}()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("accept err:", err)
			return
		}
		lock.Lock()
		connMap[connID] = conn
		lock.Unlock()
		go process(conn, t, connID)
		connID++
	}
}

func process(conn net.Conn, t *tty.TtyConn, connID int) {
	defer conn.Close()

	for {
		var buf [2048]byte
		n, err := conn.Read(buf[:])

		if err != nil {
			log.Printf("read from connect failed, err: %v\n", err)
			break
		}

		err = t.SendData(buf[:n])
		if err != nil {
			log.Println("send to ble err:", err)
		}
	}
	lock.Lock()
	defer lock.Unlock()
	delete(connMap, connID)
}
