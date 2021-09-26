package ifile

// MimeType get File Mime Type name. eg "image/png"
func MimeType(path string) (mime string) {
	if path == "" {
		return
	}

	file, err := os.Open(path)
	if err != nil {
		return
	}

	return ReaderMimeType(file)
}

// ReaderMimeType get the io.Reader mimeType
// Usage:
// 	file, err := os.Open(filepath)
// 	if err != nil {
// 		return
// 	}
//	mime := ReaderMimeType(file)
func ReaderMimeType(r io.Reader) (mime string) {
	var buf [MimeSniffLen]byte
	n, _ := io.ReadFull(r, buf[:])
	if n == 0 {
		return ""
	}

	return http.DetectContentType(buf[:n])
}