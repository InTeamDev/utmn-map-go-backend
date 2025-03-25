package utils

import (
	"encoding/json"
	"os"
	"strings"
	"unicode"

	"github.com/InTeamDev/utmn-map-go-backend/internal/domain/search/utils/translit"
)

// QueryProcessor обрабатывает поисковый запрос.
type QueryProcessor struct {
	synonymsPath string
}

// NewQueryProcessor создает новый экземпляр QueryProcessor.
func NewQueryProcessor(synonymsPath string) *QueryProcessor {
	return &QueryProcessor{
		synonymsPath: synonymsPath,
	}
}

// Process обрабатывает поисковый запрос.
func (p *QueryProcessor) Process(query string) string {
	query = p.handleTranslit(query)    // Транслитерация
	query = p.advancedNormalize(query) // Нормализация
	query = p.handleTypos(query)       // Исправление опечаток
	query = p.expandSynonyms(query)    // Расширение синонимов
	return query
}

// handleTranslit транслитерирует текст (например, "привет" -> "privet").
func (p *QueryProcessor) handleTranslit(query string) string {
	return translit.ToLatin(query) // Используем внешнюю библиотеку или свою реализацию
}

// advancedNormalize нормализует текст: приводит к нижнему регистру, удаляет лишние пробелы и спецсимволы.
func (p *QueryProcessor) advancedNormalize(query string) string {
	query = strings.ToLower(query)   // Приводим к нижнему регистру
	query = strings.TrimSpace(query) // Удаляем лишние пробелы

	// Удаляем спецсимволы
	var normalized strings.Builder
	for _, r := range query {
		if unicode.IsLetter(r) || unicode.IsDigit(r) || unicode.IsSpace(r) {
			normalized.WriteRune(r)
		}
	}
	return normalized.String()
}

// handleTypos исправляет опечатки в запросе.
func (p *QueryProcessor) handleTypos(query string) string {
	// Пример: замена часто встречающихся опечаток
	typos := map[string]string{
		"коридор":   "корридор",
		"лесница":   "лестница",
		"аудитория": "аудиторя",
	}

	for typo, correct := range typos {
		query = strings.ReplaceAll(query, typo, correct)
	}
	return query
}

// expandSynonyms расширяет запрос синонимами.
func (p *QueryProcessor) expandSynonyms(query string) string {
	synonyms := p.loadSynonyms()
	for word, synonymsList := range synonyms {
		if strings.Contains(query, word) {
			query = strings.Join(append([]string{query}, synonymsList...), " ")
		}
	}
	return query
}

// loadSynonyms загружает синонимы из JSON-файла.
func (p *QueryProcessor) loadSynonyms() map[string][]string {
	file, err := os.Open(p.synonymsPath)
	if err != nil {
		return map[string][]string{} // Возвращаем пустой словарь в случае ошибки
	}
	defer file.Close()

	var synonyms map[string][]string
	if err := json.NewDecoder(file).Decode(&synonyms); err != nil {
		return map[string][]string{} // Возвращаем пустой словарь в случае ошибки
	}
	return synonyms
}
