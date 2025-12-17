package main

import (
	"net/http"
)

func main() {
	// Настройка маршрутов
	mux := http.NewServeMux()
	
	// Регистрация обработчиков
	mux.HandleFunc("/", indexHandler)
	mux.HandleFunc("/upload", uploadHandler)
	
	// Запуск сервера
	http.ListenAndServe("0.0.0.0:8000", mux)
}