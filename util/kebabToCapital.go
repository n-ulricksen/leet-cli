package util

func KebabToCapital(kebab string) string {
	var capital []rune
	beginWord := true

	for _, ch := range kebab {
		var newCh rune
		if beginWord {
			newCh = ch - 32
			beginWord = false
		} else if ch == '-' {
			newCh = ' '
			beginWord = true
		} else {
			newCh = ch
		}
		capital = append(capital, newCh)
	}

	return string(capital)
}
