package app

import (
	"crypto/rand"
	"encoding/base64"
)

func GenerateRandomString(length int) (string, error) {
	// Создаем байтовый массив нужного размера
	byteLength := (length*3 + 3) / 4
	bytes := make([]byte, byteLength)

	// Заполняем байтовый массив случайными байтами
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}

	// Кодируем в base64 и обрезаем до нужной длины
	randomString := base64.URLEncoding.EncodeToString(bytes)
	return randomString[:length], nil
}
