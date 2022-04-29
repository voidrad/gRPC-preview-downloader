package mapservice

import "sync"

//добавляем значения в мапу с использыванием мьютекса дабы избежать состояния гонки
func AddToMap(m map[string][]byte, youtubeLink string, preview []byte, mutex *sync.Mutex) map[string][]byte {
	mutex.Lock()
	defer mutex.Unlock()
	m[youtubeLink] = preview
	return m
}
