package main

import (
	"fmt"
	"log"
	"net"
	"syscall"
)

func runChapter1() {
	var netDeviceList []netDevice
	events := make([]syscall.EpollEvents, 10)
	// epollの作成
	epdfd, err := syscall.EpollCreate1(0)
	if err != nil {
		log.Fatal("epoll create err : %s", err)
	}
	// NICの情報を取得する
	interfaces, _ := net.Interfaces()

	for _, netif := range interfaces {
		//無視するインターフェースか確認
		if !isIgnoreInterfaces(netif.Name) {
			sock, err := syscall.Socket(syscall.AF_PACKET, syscall.SOCK_RAW, int(htons(syscall.ETH_P_ALL)))
			if err != nil {
				log.Fatal("create socket err : %s", err)
			}
			addr := syscall.SockaddrLinklayer{
				Protocol: htons(syscall.ETH_P_ALL),
				Ifinfex:  netif.Index,
			}
			err = syscall.Bind(sock, &addr)
			if err != nil {
				log.Fatalf("bind err : %s", err)
			}
			fmt.Printf("Created device %s socket %d address %s\n", netif.Name, sock, netif.HardwareAddr.String())
			err = syscall.EpollCtl(epfd, syscall.EPOLL_CTL_ADD, sock,
				&syscall.EpollEvent{
					Events: syscall.EPOLLIN,
					Fd:     int32(sock),
				})
			if err != nil {
				log.Fatalf("epoll ctrl err : %s", err)
			}
			//netDevice構造体を作成
			//net_deviceの連結リストに連結
			netDeviceList = append(netDeviceList, netDevice{
				name:     netif.Name,
				macaddr:  setMacAddr(netif.HardwareAddr),
				socket:   sock,
				sockaddr: addr,
			})

		}
	}
}
