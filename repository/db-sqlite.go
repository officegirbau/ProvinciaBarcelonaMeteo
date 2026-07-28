package repository

import (
	"database/sql"
	"errors"
	"time"
)

type SQLiteRepository struct {
	Conn *sql.DB
}

// Aquesta funció retornarà el struct poblat amb la connexió a la bd
func NewSQLiteRepository(db *sql.DB) *SQLiteRepository {
	return &SQLiteRepository{
		Conn: db,
	}
}

// Desenvolupem les funcions que hem mencionat en la interfície
// Hem de indicar en el receptor que farem servir un punter sobre el mètode NewSQLiteRepository per emprar la connexió establerta per aquestes accions
func (repo *SQLiteRepository) Migrate() error {
	sentencia := `
	create table if not exists registres(
		id integer primary key autoincrement,
		municipi text not null,
		data_registre integer not null,
		precipitacio integer not null,
		temp_maxima integer not null,
		temp_minima integer not null,
		humitat integer not null)
		`
	_, err := repo.Conn.Exec(sentencia)
	return err
}

func (repo *SQLiteRepository) AgefirRegistre(nouRegistre Registres) (*Registres, error) {
	//Preparem la instrucció per afegir un registre en la taula registres
	sentencia := "insert into registres (municipi, data_registre, precipitacio, temp_maxima, temp_minima, humitat) values (?,?,?,?,?,?)" //Executem la instrucció
	resp, err := repo.Conn.Exec(
		sentencia,
		nouRegistre.Municipi,
		nouRegistre.Data.Unix(),
		nouRegistre.Precipitacio,
		nouRegistre.TempMaxima,
		nouRegistre.TempMinima,
		nouRegistre.Humitat,
	)
	if err != nil {
		return nil, err
	}
	//Afegim una crida a la funció LastInsertId() de la resposta per obtenir la id que s'ha generat amb aquesta inserció
	id, err := resp.LastInsertId()
	if err != nil {
		return nil, err
	}
	nouRegistre.ID = id
	return &nouRegistre, nil //Preparem els retorn amb un objecte o nil en d'error
}

func (repo *SQLiteRepository) VeureTotsRegistres() ([]Registres, error) {
	sentencia := "select id, municipi, data_registre, precipitacio, temp_maxima, temp_minima, humitat from registres order by id" //Executem la consulta
	files, err := repo.Conn.Query(sentencia)
	if err != nil {
		return nil, err
	}
	defer files.Close()

	var resultats []Registres
	for files.Next() {
		var fila Registres
		var unixTime int64
		err := files.Scan(
			&fila.ID,
			&fila.Municipi,
			&unixTime,
			&fila.Precipitacio,
			&fila.TempMaxima,
			&fila.TempMinima,
			&fila.Humitat,
		)
		if err != nil {
			return nil, err
		}
		fila.Data = time.Unix(unixTime, 0)
		resultats = append(resultats, fila)
	}
	return resultats, nil
}

func (repo *SQLiteRepository) VeureRegistre(id int64) (*Registres, error) {
	resp := repo.Conn.QueryRow("select id, municipi, data_registre, precipitacio, temp_maxima, temp_minima, humitat from registres where id = ?", id)

	var fila Registres
	var unixTime int64

	err := resp.Scan(
		&fila.ID,
		&fila.Municipi,
		&unixTime,
		&fila.Precipitacio,
		&fila.TempMaxima,
		&fila.TempMinima,
		&fila.Humitat,
	)

	if err != nil {
		return nil, err
	}

	fila.Data = time.Unix(unixTime, 0)

	return &fila, nil
}

func (repo *SQLiteRepository) ActualitzarRegistre(id int64, actualitzar Registres) error {
	//filtre
	if id == 0 {
		return errors.New("el valor rebut com id no és correcte")
	}
	sentencia := "update registres set municipi = ?, data_registre = ?, precipitacio = ?, temp_maxima = ?, temp_minima = ?, humitat = ? where id = ?"

	resp, err := repo.Conn.Exec(
		sentencia,
		actualitzar.Municipi,
		actualitzar.Data.Unix(),
		actualitzar.Precipitacio,
		actualitzar.TempMaxima,
		actualitzar.TempMinima,
		actualitzar.Humitat,
		id,
	)
	if err != nil {
		return err
	}

	numero, err := resp.RowsAffected()
	if err != nil {
		return err
	}
	if numero == 0 {
		return actualitzarErr
	}
	return nil
}

func (repo *SQLiteRepository) BorrarRegistre(id int64) error {
	//filtre
	if id == 0 {
		return errors.New("el valor rebut com id no és correcte")
	}

	sentencia := "delete from registres where id = ?"
	resp, err := repo.Conn.Exec(sentencia, id)
	if err != nil {
		return err
	}
	numero, err := resp.RowsAffected()
	if err != nil {
		return err
	}
	if numero == 0 {
		return borrarErr
	}
	return nil
}
