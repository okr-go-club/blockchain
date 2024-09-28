package p2p

import (
	"context"
	"time"
)

func task(ctx context.Context, node *Node) {
	// запускаем бесконечный цикл
	for {
		select {
		// проверяем не завершён ли ещё контекст и выходим, если завершён
		case <-ctx.Done():
			return

		// выполняем нужный нам код
		default:
			println("Hello gophers!")
		}
		// делаем паузу перед следующей итерацией
		time.Sleep(time.Minute)
	}
}

func cronjob(node *Node) {
	// create a scheduler
	// создаём контекст с функцией завершения
	ctx, cancel := context.WithCancel(context.Background())
	go task(ctx, node)
	// делаем паузу, чтобы дать горутине поработать
	time.Sleep(10 * time.Minute)
	// завершаем контекст, чтобы завершить горутину
	cancel()
}
