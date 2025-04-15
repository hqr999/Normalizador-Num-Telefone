package main

import (
	"fmt"
	tel_db "github.com/hqr999/Normalizador-Num-Telefone/db"
	_ "github.com/lib/pq"
	"regexp"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "hqr777"
	password = "curryH"
	dbname   = "db_telefones"
)

func main() {
	psql := fmt.Sprintf("host=%s port=%d user=%s password=%s sslmode=disable", host, port, user, password)
	check_error(tel_db.Reset("postgres", psql, dbname))

	psql = fmt.Sprintf("%s dbname=%s", psql, dbname)
	check_error(tel_db.Migrate("postgres", psql))

	db, err := tel_db.Open("postgres", psql)

	check_error(err)
	defer db.Close()

	err = db.Seed()
	check_error(err)

	tels, err := db.TodosTelefones()

	for _, p := range tels {
		fmt.Printf("Trabalhando em...%+v\n", p)
		num := normalize(p.Numero)
		if num != p.Numero {
			fmt.Println("Atualizando ou removendo...", num)
			existente, err := db.AchaTel(num)
			check_error(err)
			if existente != nil {
				check_error(db.DeletaTel(p.ID))
			} else {
				p.Numero = num
				check_error(db.AtualizaTel(p))
			}
		} else {
			fmt.Println("Sem mudanças necessárias")
		}
	}
}

/*
	func getTel(db *sql.DB, id int) (string, error) {
		var num string
		linha := db.QueryRow("SELECT * FROM num_telefones WHERE id=$1", id)
		err := linha.Scan(&id, &num)
		if err != nil {
			return "", err
		}
		return num, nil
	}
*/
func check_error(err error) {
	if err != nil {
		panic(err)
	}
}

func normalize(telefone string) string {
	re := regexp.MustCompile("[^0-9]")
	return re.ReplaceAllString(telefone, "")

}
