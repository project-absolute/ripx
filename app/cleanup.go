package main

import (
	"context"
	"os"
	"path/filepath"
	"time"
)

// cleanupDuration определяет как долго изображения хранятся перед очисткой
const cleanupDuration = 24 * time.Hour

// cleanupInterval определяет как часто запускается очистка
const cleanupInterval = 1 * time.Hour

// startCleanupWorker запускает фоновый процесс очистки старых изображений
func startCleanupWorker(ctx context.Context) {
	ticker := time.NewTicker(cleanupInterval)
	defer ticker.Stop()
	
	for {
		select {
		case <-ctx.Done():
			// Graceful shutdown при отмене контекста
			return
		case <-ticker.C:
			// Периодическая очистка
			err := cleanupOldImages()
			if err != nil {
				// Ошибки игнорируем, продолжаем работу
				continue
			}
			
			err = removeEmptyDirectories()
			if err != nil {
				// Ошибки игнорируем, продолжаем работу
				continue
			}
		}
	}
}

// cleanupOldImages удаляет старые изображения
func cleanupOldImages() error {
	// Проверяем существование директории /data
	dataDir := "/data"
	if _, err := os.Stat(dataDir); os.IsNotExist(err) {
		return nil
	}
	
	// Читаем все пользовательские директории
	entries, err := os.ReadDir(dataDir)
	if err != nil {
		return err
	}
	
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		
		userDir := filepath.Join(dataDir, entry.Name())
		err := cleanupUserImages(userDir)
		if err != nil {
			// Продолжаем очистку других директорий при ошибке
			continue
		}
	}
	
	return nil
}

// cleanupUserImages очищает старые изображения в директории пользователя
func cleanupUserImages(userDir string) error {
	entries, err := os.ReadDir(userDir)
	if err != nil {
		return err
	}
	
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		
		filePath := filepath.Join(userDir, entry.Name())
		
		// Получаем информацию о файле
		info, err := entry.Info()
		if err != nil {
			// Пропускаем файл при ошибке
			continue
		}
		
		// Проверяем, является ли файл старым
		// Используем время модификации как fallback для времени доступа
		lastAccess := info.ModTime()
		if isImageOld(lastAccess) {
			// Удаляем старый файл
			os.Remove(filePath)
		}
	}
	
	return nil
}

// isImageOld проверяет, является ли изображение старым
func isImageOld(lastAccess time.Time) bool {
	// Изображение считается старым, если время последнего доступа
	// старше cleanupDuration (24 часа)
	return time.Since(lastAccess) > cleanupDuration
}

// removeEmptyDirectories удаляет пустые директории
func removeEmptyDirectories() error {
	dataDir := "/data"
	
	// Сначала удаляем пустые пользовательские директории
	entries, err := os.ReadDir(dataDir)
	if err != nil {
		return err
	}
	
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		
		userDir := filepath.Join(dataDir, entry.Name())
		
		// Проверяем, пуста ли директория
		isEmpty, err := isDirEmpty(userDir)
		if err != nil {
			continue
		}
		
		if isEmpty {
			// Удаляем пустую директорию пользователя
			os.Remove(userDir)
		}
	}
	
	return nil
}

// isDirEmpty проверяет, пуста ли директория
func isDirEmpty(dirPath string) (bool, error) {
	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return false, err
	}
	
	return len(entries) == 0, nil
}