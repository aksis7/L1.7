package main

import (
	"fmt"
	"hash/fnv"
	"sync"
)

// Количество шардов для разделения карты
const shardCount = 32

// ShardedMap — структура, разделяющая карту на шардированные сегменты для повышения производительности.
type ShardedMap struct {
	shards []map[string]interface{} // Список карт (шарды)
	mu     []sync.Mutex             // Список мьютексов для каждого шарда
}

// NewShardedMap создает новый ShardedMap с заданным количеством шардов.
func NewShardedMap() *ShardedMap {
	shards := make([]map[string]interface{}, shardCount)
	mu := make([]sync.Mutex, shardCount)
	for i := 0; i < shardCount; i++ {
		shards[i] = make(map[string]interface{}) // Инициализация каждого шарда
	}
	return &ShardedMap{shards: shards, mu: mu}
}

// getShard определяет, к какому шару относится ключ, используя хэш-функцию.
func (sm *ShardedMap) getShard(key string) int {
	hasher := fnv.New32a() // Создаем хэш-функцию
	hasher.Write([]byte(key))
	return int(hasher.Sum32()) % shardCount // Определяем индекс шарда
}

// Set записывает данные в соответствующий шард.
func (sm *ShardedMap) Set(key string, value interface{}) {
	shardIndex := sm.getShard(key)   // Находим шард по ключу
	sm.mu[shardIndex].Lock()         // Блокируем доступ к этому шару
	defer sm.mu[shardIndex].Unlock() // Разблокируем после завершения
	sm.shards[shardIndex][key] = value
}

// Get возвращает данные из соответствующего шарда.
func (sm *ShardedMap) Get(key string) (interface{}, bool) {
	shardIndex := sm.getShard(key)   // Находим шард по ключу
	sm.mu[shardIndex].Lock()         // Блокируем доступ к этому шару
	defer sm.mu[shardIndex].Unlock() // Разблокируем после завершения
	val, ok := sm.shards[shardIndex][key]
	return val, ok
}

func main() {
	shardedMap := NewShardedMap() // Создаем шардированную карту
	var wg sync.WaitGroup         // Для ожидания завершения всех горутин

	// Запускаем 10 горутин для конкурентной записи данных в карту
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done() // Уменьшаем счетчик wg после завершения горутины
			shardedMap.Set(fmt.Sprintf("key%d", i), i)
		}(i)
	}

	wg.Wait() // Ждем завершения всех горутин

	// Считываем и выводим значения из карты
	for i := 0; i < 10; i++ {
		val, ok := shardedMap.Get(fmt.Sprintf("key%d", i))
		if ok {
			fmt.Printf("key%d: %v\n", i, val)
		}
	}
}
