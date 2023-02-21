package shortUrl

// 计算
func Generate1(id int64) (tiny string) {
	num := id
	alpha := merge(getRange(48, 57), getRange(65, 90))
	alpha = merge(alpha, getRange(97, 122))
	if num < 62 {
		tiny = string(alpha[num])
		return tiny
	} else {
		var runes []rune
		runes = append(runes, alpha[num%62])
		num = num / 62
		for num >= 1 {
			if num < 62 {
				runes = append(runes, alpha[num-1])
			} else {
				runes = append(runes, alpha[num%62])
			}
			num = num / 62

		}
		for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
			runes[i], runes[j] = runes[j], runes[i]
		}
		tiny = string(runes)
		return tiny
	}
	return tiny
}

func getRange(start, end rune) (ran []rune) {
	for i := start; i <= end; i++ {
		ran = append(ran, i)
	}
	return ran
}

func merge(a, b []rune) []rune {
	c := make([]rune, len(a)+len(b))
	copy(c, a)
	copy(c[len(a):], b)
	return c
}
