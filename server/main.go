package main

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
)

func run(ctx context.Context, addr *net.UDPAddr) {
	c, e := net.ListenUDP("udp4", addr)
	if e != nil {
		log.Fatalln(e)
	}
	defer c.Close()

	buf := make([]byte, 1024)
	for {
		select {
		case <-ctx.Done():
			log.Println("上下文结束")
			return
		default:
			n, clientAddr, e := c.ReadFromUDP(buf)
			if e != nil {
				log.Println("读取错误", e)
				continue
			}
			text := buf[:n]
			log.Printf("收到:%s 请求: %s", clientAddr, text)

			_, e = c.WriteToUDP([]byte(text), clientAddr)
			if e != nil {
				log.Println("回复错误", e)
				continue
			}
		}
	}
}

func main() {
	ctx, ctxCancel := context.WithCancel(context.Background())

	addr, e := net.ResolveUDPAddr("udp4", "172.17.0.1:10000")
	if e != nil {
		log.Fatalln(e)
	}
	go run(ctx, addr)

	// 等待关闭信号
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGTERM, syscall.SIGINT)
	stopSignal := <-signalChan
	log.Println("收到关闭信号", stopSignal)

	ctxCancel()
}
