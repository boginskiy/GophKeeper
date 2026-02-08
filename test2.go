package main

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
)

func generateRandom(size int) ([]byte, error) {
	// генерируем случайную последовательность байт
	b := make([]byte, size)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}

	return b, nil
}

func main() {
	// подписываемое сообщение
	src := []byte("Видишь гофера? Нет. И я нет. А он есть.")

	// создаём случайный ключ
	key, err := generateRandom(16)
	if err != nil {
		fmt.Printf("error: %v\n", err)
		return
	}

	// подписываем алгоритмом HMAC, используя SHA-256
	h := hmac.New(sha256.New, key)
	h.Write(src)
	dst := h.Sum(nil)

	fmt.Printf("%x", dst)
}

// Для проверки подлинности подписи надо:
// Заново пересчитать подпись с тем же ключом.
// Функцией hmac.Equal проверить, что подписи совпадают.
