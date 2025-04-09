package main

import (
	"database/sql"
	"fmt"
	"regexp"

	_ "github.com/lib/pq"
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

const (
	host     = "localhost"
	port     = 5432
	user     = "hreuter"
	password = "cavalo77"
	dbname   = "db_telefones"
)

func main() {
	psql := fmt.Sprintf("host=%s port=%d user=%s password=%s sslmode=disable", host, port, user, password)
	db, err := sql.Open("postgres", psql)
	if err != nil {
		panic(err)
	}
	err = createDB(db, dbname)
	if err != nil {
		panic(err)
	}
	db.Close()
}

func createDB(db *sql.DB, name string) error {
	_, err := db.Exec("CREATE DATABASE" + name)
	fmt.Printf("%s\n",name)
	if err != nil {
		return err
	}
	return nil
}

func normalize(telefone string) string {
	re := regexp.MustCompile("[^0-9]")
	return re.ReplaceAllString(telefone, "")

}
