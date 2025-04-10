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
	user     = "hqr777"
	password = "curryH"
	dbname   = "db_telefones"
)

func main() {
	psql := fmt.Sprintf("host=%s port=%d user=%s password=%s sslmode=disable", host, port, user, password)
	bd, err := sql.Open("postgres", psql)
	check_error(err)
	err = resetBD(bd, dbname)
	check_error(err)
	bd.Close()

	psql = fmt.Sprintf("%s dbname=%s", psql, dbname)
	bd, err = sql.Open("postgres", psql)
	check_error(err)
	defer bd.Close()
	check_error(criaTabelaTelefones(bd))

}

func criaTabelaTelefones(bd *sql.DB) error {
	declaracao := `
			CREATE TABLE IF NOT EXISTS num_telefones (
				id SERIAL,
				valor VARCHAR(255)
			)
		`
	_, err := bd.Exec(declaracao)
	return err

}

func check_error(err error) {
	if err != nil {
		panic(err)
	}
}

func createBD(db *sql.DB, name string) error {
	_, err := db.Exec("CREATE DATABASE " + name)
	if err != nil {
		return err
	}
	return nil
}

func resetBD(db *sql.DB, name string) error {
	_, err := db.Exec("DROP DATABASE IF EXISTS " + name)
	if err != nil {
		return err
	}
	return createBD(db, name)
}

func normalize(telefone string) string {
	re := regexp.MustCompile("[^0-9]")
	return re.ReplaceAllString(telefone, "")

}
