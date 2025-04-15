package db

import (
	"database/sql"

	_ "github.com/lib/pq"
)

type Telefone struct {
	ID     int
	Numero string
}

func Open(driverName, dataSource string) (*DB, error) {
	db, err := sql.Open(driverName, dataSource)
	if err != nil {
		return nil, err
	}
	return &DB{db}, nil

}

type DB struct {
	db *sql.DB
}

func (x *DB) Close() error {
	return x.db.Close()
}

func (x *DB) Seed() error {
	data := []string{
		"1234567890",
		"123 456 7891",
		"(123) 456 7892",
		"(123) 456-7893",
		"123-456-7894",
		"123-456-7890",
		"1234567892",
		"(123)456-7892",
	}

	for _, n := range data {
		if _, err := insertPhone(x.db, n); err != nil {
			return err
		}
	}
	return nil

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

func (x *DB) TodosTelefones() ([]Telefone, error) {
	return todosTelefones(x.db)
}

func todosTelefones(db *sql.DB) ([]Telefone, error) {
	linhas, err := db.Query("SELECT id, valor FROM num_telefones")
	if err != nil {
		return nil, err
	}

	defer linhas.Close()

	var retorno []Telefone
	for linhas.Next() {
		var p Telefone
		if err := linhas.Scan(&p.ID, &p.Numero); err != nil {
			return nil, err
		}
		retorno = append(retorno, p)
	}
	if err := linhas.Err(); err != nil {
		return nil, err
	}
	return retorno, nil
}

func (x *DB) AchaTel(number string) (*Telefone, error) {
	return achaTel(x.db, number)
}

func achaTel(db *sql.DB, number string) (*Telefone, error) {
	var p Telefone
	linha := db.QueryRow("SELECT * FROM num_telefones WHERE valor=$1", number)
	err := linha.Scan(&p.ID, &p.Numero)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		} else {
			return nil, err
		}
	}
	return &p, nil
}

func (x *DB) AtualizaTel(p Telefone) error {

	statment := `UPDATE num_telefones SET valor=$2 WHERE id=$1`
	_, err := x.db.Exec(statment, p.ID, p.Numero)
	return err

}

func (x *DB) DeletaTel(id int) error {
	return deletaTel(x.db, id)
}

func deletaTel(db *sql.DB, id int) error {
	statment := `DELETE FROM num_telefones WHERE id=$1`
	_, err := db.Exec(statment, id)
	return err
}

func Migrate(driverName, dataSource string) error {
	db, err := sql.Open(driverName, dataSource)
	if err != nil {
		return err
	}

	err = criaTabelaTelefones(db)
	if err != nil {
		return err
	}

	return db.Close()
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

func Reset(driverName, dataSource, dbName string) error {
	db, err := sql.Open(driverName, dataSource)
	if err != nil {
		return err
	}
	err = resetBD(db, dbName)
	if err != nil {
		return err
	}
	return db.Close()

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
