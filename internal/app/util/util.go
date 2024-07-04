package util

import (
	"crypto/rand"
	"encoding/base64"
)

func GenerateShortURL(length int) string {
	// Создаем байтовый массив нужного размера
	byteLength := (length*3 + 3) / 4
	bytes := make([]byte, byteLength)

	_, _ = rand.Read(bytes)

	// Кодируем в base64 и обрезаем до нужной длины
	randomString := base64.URLEncoding.EncodeToString(bytes)
	return randomString[:length]
}
