package main

import (
	"fmt"
	"sync"
)

func main() {
	var sm sync.Map       // sync.Map — потокобезопасная структура для хранения данных
	var wg sync.WaitGroup // Для ожидания завершения всех горутин

	// Запускаем 10 горутин для конкурентной записи данных в sync.Map
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()                      // Уменьшаем счетчик wg после завершения горутины
			sm.Store(fmt.Sprintf("key%d", i), i) // Store — метод для записи данных в sync.Map
		}(i)
	}

	wg.Wait() // Ждем завершения всех горутин

	// Читаем и выводим значения из sync.Map
	for i := 0; i < 10; i++ {
		if val, ok := sm.Load(fmt.Sprintf("key%d", i)); ok { // Load — метод для чтения данных
			fmt.Printf("key%d: %v\n", i, val)
		}
	}
}
