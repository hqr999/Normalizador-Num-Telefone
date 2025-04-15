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
	_, err = insertPhone(bd, "1234567890")
	check_error(err)
	_, err = insertPhone(bd, "123 456 7891")
	check_error(err)
	_, err = insertPhone(bd, "(123) 456 7892")
	check_error(err)
	_, err = insertPhone(bd, "(123) 456-7893")
	check_error(err)
	_, err = insertPhone(bd, "123-456-7894")
	check_error(err)
	_, err = insertPhone(bd, "123-456-7890")
	check_error(err)
	id, err := insertPhone(bd, "1234567892")
	check_error(err)
	_, err = insertPhone(bd, "(123)456-7892")
	check_error(err)

	num, err := getTel(bd, id)
	check_error(err)
	fmt.Println("Número é ... ", num)
	tels, err := todosTelefones(bd)
	check_error(err)

	for _, p := range tels {
		fmt.Printf("Trabalhando em...%+v\n", p)
		num := normalize(p.numero)
		if num != p.numero {
			fmt.Println("Atualizando ou removendo...", num)
			existente, err := achaTel(bd, num)
			check_error(err)
			if existente != nil {
				check_error(deletaTel(bd, p.id))
			} else {
				p.numero = num
				check_error(atualizaTel(bd, p))
			}
		} else {
			fmt.Println("Sem mudanças necessárias")
		}
	}
}

func atualizaTel(db *sql.DB, p telefone) error {
	statment := `UPDATE num_telefones SET valor=$2 WHERE id=$1`
	_, err := db.Exec(statment, p.id, p.numero)
	return err

}

func deletaTel(db *sql.DB, id int) error {
	statment := `DELETE FROM num_telefones WHERE id=$1`
	_, err := db.Exec(statment, id)
	return err
}

type telefone struct {
	id     int
	numero string
}

func todosTelefones(db *sql.DB) ([]telefone, error) {
	linhas, err := db.Query("SELECT id, valor FROM num_telefones")
	if err != nil {
		return nil, err
	}

	defer linhas.Close()

	var retorno []telefone
	for linhas.Next() {
		var p telefone
		if err := linhas.Scan(&p.id, &p.numero); err != nil {
			return nil, err
		}
		retorno = append(retorno, p)
	}
	if err := linhas.Err(); err != nil {
		return nil, err
	}
	return retorno, nil
}

func getTel(db *sql.DB, id int) (string, error) {
	var num string
	linha := db.QueryRow("SELECT * FROM num_telefones WHERE id=$1", id)
	err := linha.Scan(&id, &num)
	if err != nil {
		return "", err
	}
	return num, nil
}

func achaTel(db *sql.DB, number string) (*telefone, error) {
	var p telefone
	linha := db.QueryRow("SELECT * FROM num_telefones WHERE valor=$1", number)
	err := linha.Scan(&p.id, &p.numero)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		} else {
			return nil, err
		}
	}
	return &p, nil
}
func insertPhone(db *sql.DB, tel string) (int, error) {
	statement := `INSERT INTO num_telefones(valor) VALUES($1) RETURNING id`
	var id int
	err := db.QueryRow(statement, tel).Scan(&id)
	if err != nil {
		return -1, err
	}
	return id, nil
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
