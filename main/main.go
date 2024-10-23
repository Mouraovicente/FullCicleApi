package main

import (
	"fmt"
	"time"
)

func main2() {
	//cria um canal de inteiro com o nome da variavel canal
	canal := make(chan int)
	//t2
	go func() {
		for i := range 10 {
			canal <- i + 10
			time.Sleep(time.Second)
		}
	}()

	for v := range canal {
		fmt.Println(v)
	}

	// multithread
	//go contador(5)
	//go contador(8)

}

func contador(count int) {
	for i := range count {
		fmt.Println(i)
		time.Sleep(time.Second)
	}
}
