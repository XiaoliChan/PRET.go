package main

import (
	"io"
	"net"
	"time"
)

func connectTCP(target string) net.Conn {
	conn, err := net.DialTimeout("tcp", target, time.Duration(5)*time.Second)
	if err != nil {
		defer conn.Close()
		return nil
	}
	err = conn.SetDeadline(time.Now().Add(time.Duration(5) * time.Second))
	if err != nil {
		defer conn.Close()
		return nil
	}
	return conn
}

func GetNetResponse(conn net.Conn) ([]byte, error) {
	conn.SetReadDeadline(time.Now().Add(time.Second))
	bytes, err := io.ReadAll(conn)
	if len(bytes) > 0 {
		err = nil
	}
	return bytes, err
}

func ReadBytes(conn net.Conn) (result []byte, err error) {
	buf := make([]byte, 4096)
	for {
		count, err := conn.Read(buf)
		if err != nil {
			break
		}
		result = append(result, buf[0:count]...)
		if count < 4096 {
			break
		}
	}
	return result, err
}
