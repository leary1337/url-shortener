package util

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
)

func GenerateShortURL(addr string, length int) string {
	// Создаем байтовый массив нужного размера
	byteLength := (length*3 + 3) / 4
	bytes := make([]byte, byteLength)

	_, _ = rand.Read(bytes)

	// Кодируем в base64 и обрезаем до нужной длины
	randomString := base64.URLEncoding.EncodeToString(bytes)
	return fmt.Sprintf("%s/%s", addr, randomString[:length])
}
