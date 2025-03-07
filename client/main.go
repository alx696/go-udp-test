package main

import (
	"fmt"
	"log"
	"net"
	"sync"
	"time"
)

func run(addr *net.UDPAddr, waitGroup *sync.WaitGroup, limitChan chan bool, tag string) error {
	defer func() {
		<-limitChan
	}()
	defer waitGroup.Done()

	c, e := net.DialUDP("udp4", nil, addr)
	if e != nil {
		return e
	}
	defer c.Close()

	count := 0
	for count < 10000 {
		count++
		_, e = c.Write([]byte(fmt.Sprint("a", tag, "b", count)))
		if e != nil {
			return e
		}

		buf := make([]byte, 1024)
		n, e := c.Read(buf)
		if e != nil {
			return e
		}
		log.Println(string(buf[:n]))
	}

	return nil
}

func main() {
	var waitGroup sync.WaitGroup
	limitChan := make(chan bool, 100)

	addr, e := net.ResolveUDPAddr("udp4", "172.17.0.1:10000")
	if e != nil {
		log.Fatalln(e)
	}

	timeBegin := time.Now()
	count := 0
	// 数字比较大时不限制并发的话waitGroup容易卡住不能结束!!!
	for count < 1 {
		count++
		waitGroup.Add(1)
		limitChan <- true
		go run(addr, &waitGroup, limitChan, fmt.Sprint(count))
	}

	waitGroup.Wait()
	close(limitChan)

	log.Println("耗时", time.Now().UnixMilli()-timeBegin.UnixMilli())
}
