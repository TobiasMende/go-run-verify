package logging

import (
	"fmt"
	"net"
)

type SocketConnector struct {
	Port int
	udp  *net.UDPConn
}

func (conn *SocketConnector) Init() error {
	addr, err := net.ResolveUDPAddr("udp", fmt.Sprint("localhost:", conn.Port))
	if err != nil {
		Warning.Println("Unable to resolve UDP addr for port ", conn.Port)
		return err
	}
	udp, err := net.ListenUDP("udp", addr)
	if err != nil {
		Warning.Println("Unable to listen on UDP port ", conn.Port)
	}
	conn.udp = udp
	return err
}

func (conn *SocketConnector) Receive() (msg interface{}, err error) {
	var b, oob []byte
	_, _, _, _, err = conn.udp.ReadMsgUDP(b, oob)

	return b, err
}

func (conn *SocketConnector) Close() error {
	return conn.udp.Close()
}
