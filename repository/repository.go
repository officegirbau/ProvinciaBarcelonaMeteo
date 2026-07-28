package repository

import (
	"errors"
	"time"
)

var (
	actualitzarErr = errors.New("l'actualització ha fallat")
	borrarErr      = errors.New("l'esborrat ha fallat")
)

type Repository interface {
	//Establim la funció migrate per crear totes les taules que necessitem en la nostra bd
	Migrate() error
	AgefirRegistre(nouRegistre Registres) (*Registres, error)
	VeureTotsRegistres() ([]Registres, error)
	VeureRegistre(id int64) (*Registres, error)
	ActualitzarRegistre(id int64, actualitzar Registres) error
	BorrarRegistre(id int64) error
}

// A continuació definim un struct amb els camps i el tipus de dades que emprarem
type Registres struct {
	ID           int64     `json:"id"`
	Municipi     string    `json:"municipi"`
	Data         time.Time `json:"data_registre"`
	Precipitacio int       `json:"precipitacio"`
	TempMaxima   int       `json:"temp_maxima"`
	TempMinima   int       `json:"temp_minima"`
	Humitat      int       `json:"humitat"`
}
