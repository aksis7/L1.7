package main

import (
	"fmt"
	"sync"
)

// SafeMap — структура, которая содержит обычную карту и мьютекс для защиты её от конкурентного доступа.
type SafeMap struct {
	mu   sync.Mutex             // Мьютекс для синхронизации доступа
	data map[string]interface{} // Карта для хранения данных
}

// NewSafeMap создает новый экземпляр SafeMap.
func NewSafeMap() *SafeMap {
	return &SafeMap{
		data: make(map[string]interface{}), // Инициализация карты
	}
}

// Set добавляет или обновляет элемент в карте, защищённой мьютексом.
func (sm *SafeMap) Set(key string, value interface{}) {
	sm.mu.Lock()         // Блокируем доступ для других горутин
	defer sm.mu.Unlock() // Разблокируем после завершения функции
	sm.data[key] = value
}

// Get возвращает значение из карты по ключу, если оно существует.
func (sm *SafeMap) Get(key string) (interface{}, bool) {
	sm.mu.Lock()         // Блокируем доступ для чтения
	defer sm.mu.Unlock() // Разблокируем после завершения
	val, ok := sm.data[key]
	return val, ok
}

func main() {
	safeMap := NewSafeMap() // Создаем потокобезопасную карту
	var wg sync.WaitGroup   // Для ожидания завершения всех горутин

	// Запускаем 10 горутин для конкурентной записи данных в карту
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done() // Уменьшаем счетчик wg после завершения горутины
			safeMap.Set(fmt.Sprintf("key%d", i), i)
		}(i)
	}

	wg.Wait() // Ждем завершения всех горутин

	// Считываем и выводим все значения из карты
	for i := 0; i < 10; i++ {
		val, ok := safeMap.Get(fmt.Sprintf("key%d", i))
		if ok {
			fmt.Printf("key%d: %v\n", i, val)
		}
	}
}
