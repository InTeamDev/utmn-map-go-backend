package translit

import (
	"strings"
	"unicode"
)

var translitMap = map[rune]string{
	'а': "a", 'б': "b", 'в': "v", 'г': "g", 'д': "d", 'е': "e", 'ё': "yo", 'ж': "zh",
	'з': "z", 'и': "i", 'й': "y", 'к': "k", 'л': "l", 'м': "m", 'н': "n", 'о': "o",
	'п': "p", 'р': "r", 'с': "s", 'т': "t", 'у': "u", 'ф': "f", 'х': "kh", 'ц': "ts",
	'ч': "ch", 'ш': "sh", 'щ': "shch", 'ъ': "", 'ы': "y", 'ь': "", 'э': "e", 'ю': "yu",
	'я': "ya",
}

// ToLatin транслитерирует русский текст в латиницу.
func ToLatin(text string) string {
	var result strings.Builder
	for _, r := range text {
		if val, ok := translitMap[unicode.ToLower(r)]; ok {
			result.WriteString(val)
		} else {
			result.WriteRune(r)
		}
	}
	return result.String()
}
