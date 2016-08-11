package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
)

type proxyServer struct {
	dumpTo     *os.File
	localAddr  string
	remoteAddr string
}

func (ps *proxyServer) start() error {
	listener, err := net.Listen("tcp", ps.localAddr)
	if err != nil {
		return err
	}
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println(err)
		}
		go ps.handleClient(conn)
	}
}

func (ps *proxyServer) healthcheck() error {
	conn, err := net.Dial("tcp", ps.remoteAddr)
	if err != nil {
		return err
	}
	return conn.Close()
}

func (ps *proxyServer) handleClient(conn net.Conn) {
	remoteConn, err := net.Dial("tcp", ps.remoteAddr)
	if err != nil {
		log.Println(err)
		return
	}
	req := &dumper{w: remoteConn, label: requestFrom(conn), dumpTo: ps.dumpTo}
	resp := &dumper{w: conn, label: responseTo(conn), dumpTo: ps.dumpTo}
	go io.Copy(resp, remoteConn)
	_, err = io.Copy(req, conn)
	if err != nil {
		log.Println(err)
		return
	}
	err = ps.dumpTo.Sync()
	if err != nil {
		log.Println(err)
		return
	}
}

func requestFrom(conn net.Conn) string {
	return fmt.Sprintf("[==>>] Request from %s", conn.RemoteAddr().String())
}

func responseTo(conn net.Conn) string {
	return fmt.Sprintf("[<<==] Response to %s", conn.RemoteAddr().String())
}
