package istrings

import "fmt"

// rpad adds padding to the right of a string.
func Rpad(s string, padding int) string {
	template := fmt.Sprintf("%%-%ds", padding)
	return fmt.Sprintf(template, s)
}

func PadRight(str, padstr string, strlenght int) string {
	if len(str) >= strlenght {
		return str
	}
	for {
		str += padstr
		if len(str) >= strlenght {
			return str[0:strlenght]
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
