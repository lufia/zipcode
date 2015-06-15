package zipcode

func normalizeText(s string) string {
	tab := map[rune]rune{
		'（': '(',
		'）': ')',
		'０': '0',
		'１': '1',
		'２': '2',
		'３': '3',
		'４': '4',
		'５': '5',
		'６': '6',
		'７': '7',
		'８': '8',
		'９': '9',
	}
	r := []rune(s)
	buf := make([]rune, len(r))
	for i, c := range r {
		t, ok := tab[c]
		if ok {
			buf[i] = t
		} else {
			buf[i] = c
		}
	}
	return string(buf)
}
