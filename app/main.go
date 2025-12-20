package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
)

func main() {
	// Создание директории для хранения данных
	err := ensureDataDir()
	if err != nil {
		os.Exit(1)
	}
	
	// Создаем контекст для graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	
	// Настройка маршрутов
	mux := http.NewServeMux()
	
	// Статические файлы с правильным MIME типом
	mux.HandleFunc("/static/", func(w http.ResponseWriter, r *http.Request) {
		// Определяем MIME тип на основе расширения файла
		ext := filepath.Ext(r.URL.Path)
		switch ext {
		case ".css":
			w.Header().Set("Content-Type", "text/css")
		case ".js":
			w.Header().Set("Content-Type", "application/javascript")
		case ".png":
			w.Header().Set("Content-Type", "image/png")
		case ".jpg", ".jpeg":
			w.Header().Set("Content-Type", "image/jpeg")
		case ".gif":
			w.Header().Set("Content-Type", "image/gif")
		case ".svg":
			w.Header().Set("Content-Type", "image/svg+xml")
		}
		
		// Убираем префикс /static/ из пути
		filePath := r.URL.Path[1:] // убираем первый /
		http.ServeFile(w, r, filePath)
	})
	
	// Регистрация обработчиков
	mux.HandleFunc("/", indexHandler)
	mux.HandleFunc("/upload", uploadHandler)
	mux.HandleFunc("/delete-image", deleteImageHandler)
	mux.HandleFunc("/delete-album", deleteAlbumHandler)
	mux.HandleFunc("/delete-user", deleteUserHandler)
	
	// Запуск cleanup worker в отдельной goroutine
	go startCleanupWorker(ctx)
	
	// Настройка graceful shutdown через сигналы
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	
	// Запуск сервера в отдельной goroutine
	serverErr := make(chan error, 1)
	go func() {
		serverErr <- http.ListenAndServe("0.0.0.0:8000", mux)
	}()
	
	// Ожидание сигнала или ошибки сервера
	select {
	case <-sigChan:
		// Получен сигнал завершения
		cancel() // Останавливаем cleanup worker
	case <-serverErr:
		// Ошибка сервера
		cancel() // Останавливаем cleanup worker
	}
}