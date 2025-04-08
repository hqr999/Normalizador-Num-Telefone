package main

import (
	"regexp"
)

//func normalize(telefone string) string {
// Queremos o telefone nesse modelo - 012345678

//var buf bytes.Buffer

//for _, ch := range telefone {
//0 - 9 na tabela ascii estão entre os índices 0 e 9
//if ch >= '0' && ch <= '9' {
//buf.WriteRune(ch)
//}

//}
//return buf.String()

//}

// Também é uma função que normaliza um número de telefone, mas usa regex
func normalize(telefone string) string {
	re := regexp.MustCompile("[^0-9]")
	return re.ReplaceAllString(telefone, "")

}
