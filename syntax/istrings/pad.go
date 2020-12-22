package istrings

func PadRight(str, pad string, lenght int) string {
	if len(str) >= lenght {
		return str
	}
	for {
		str += pad
		if len(str) >= lenght {
			return str[0:lenght]
		}
	}
}

func PadLeft(str, pad string, lenght int) string {
	if len(str) >= lenght {
		return str
	}
	for {
		str = pad + str
		if len(str) >= lenght {
			return str[0:lenght]
		}
	}
}
