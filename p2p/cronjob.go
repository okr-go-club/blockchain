package p2p

import (
	"blockchain/chain"
	"context"
	"fmt"
	"net"
	"time"
)

func task(ctx context.Context, node *Node, conn net.Conn, blockchain *chain.Blockchain) {
	for {
		select {
		// проверяем не завершён ли ещё контекст и выходим, если завершён
		case <-ctx.Done():
			return

		default:
			// запросить длину блокчейна
			_, err := conn.Write([]byte("Give me length of your blockchain!\n"))
			if err != nil {
				fmt.Println("Error writing to connection:", err)
				node.RemoveConnection(conn.RemoteAddr().String())
				return
			}
		}
		// делаем паузу перед следующей итерацией
		time.Sleep(time.Minute)
	}
}

func cronjob(node *Node, conn net.Conn, blockchain *chain.Blockchain) {
	// create a scheduler
	// создаём контекст с функцией завершения
	ctx, cancel := context.WithCancel(context.Background())
	go task(ctx, node, conn, blockchain)
	// делаем паузу, чтобы дать горутине поработать
	time.Sleep(10 * time.Minute)
	// завершаем контекст, чтобы завершить горутину
	cancel()
}
