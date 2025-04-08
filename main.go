package main

import "bytes"

func normalize(telefone string) string {
	// Queremos o telefone nesse modelo - 012345678

	var buf bytes.Buffer

	for _, ch := range telefone {
		//0 - 9 na tabela ascii estão entre os índices 0 e 9
		if ch >= '0' && ch <= '9' {
			buf.WriteRune(ch)
		}

	}
	return buf.String()

}
