package main

import (
	"database/sql"
	"ecohortapp/repository"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/widget"
	_ "github.com/glebarez/go-sqlite"
)

type Config struct {
	App                                fyne.App
	InfoLog                            *log.Logger
	ErrorLog                           *log.Logger
	MainWindow                         fyne.Window
	HTTPClient                         http.Client
	ClimaDadesContainer                *fyne.Container
	PronosticGraficContainer           *fyne.Container
	DB                                 repository.Repository
	Registres                          [][]interface{}
	RegistresTable                     *widget.Table
	RegistresList                      []repository.Registres
	AfegirRegistresDataRegistreEntrada *widget.Entry
	AfegirRegistresPrecipitacioEntrada *widget.Entry
	AfegirRegistresTempMaximaEntrada   *widget.Entry
	AfegirRegistresTempMinimaEntrada   *widget.Entry
	AfegirRegistresHumitatEntrada      *widget.Entry
	LlistaMunicipis                    []Municipi
}

var myApp Config

type Municipi struct {
	Nom  string
	Codi string
}

func carregarMunicipis(rutaArxiu string) ([]Municipi, error) {
	fitxer, err := os.Open(rutaArxiu)
	if err != nil {
		return nil, fmt.Errorf("no s'ha pogut obrir el fitxer: %v", err)
	}
	defer fitxer.Close()

	lector := csv.NewReader(fitxer)

	// 1. Llegim i descartem la primera línia (la capçalera: Municipi,Codi)
	_, err = lector.Read()
	if err != nil {
		if err == io.EOF {
			return nil, fmt.Errorf("el fitxer CSV està buit")
		}
		return nil, fmt.Errorf("error en llegir la capçalera: %v", err)
	}

	var llistaMunicipis []Municipi

	// 2. Iterem per la resta de files del CSV
	for {
		registre, err := lector.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("error en llegir una línia del CSV: %v", err)
		}

		if len(registre) >= 2 {
			municipi := Municipi{
				Nom:  registre[0], // Columna 1: Municipi
				Codi: registre[1], // Columna 2: Codi
			}
			llistaMunicipis = append(llistaMunicipis, municipi)
		}
	}

	return llistaMunicipis, nil
}

const (
	appID        = "com.ecohortapp.app"
	windowTitle  = "Previsó temps provincia de Barcelona propers 4 dies"
	windowWidth  = 800
	windowHeight = 500
)

func main() {
	//Iniciem l'app de fyne
	laMevaApp := app.NewWithID(appID)
	myApp.App = laMevaApp

	//Crear els logs
	myApp.InfoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	myApp.ErrorLog = log.New(os.Stdout, "ERROR\t", log.Ldate|log.Lshortfile)
	//inicialització variable municipi
	MunicipiD = "073"
	//Connexió amb la BBDD
	sqlDB, err := myApp.connectSQL()
	if err != nil {
		log.Panic(err)
	}

	//Repositori de la BBDD
	myApp.setupDB(sqlDB)
	myApp.MainWindow = laMevaApp.NewWindow(windowTitle)
	myApp.MainWindow.Resize(fyne.NewSize(windowWidth, windowHeight))
	myApp.MainWindow.SetFixedSize(true)
	myApp.MainWindow.SetMaster()

	// Carreguem els municipis des del fitxer CSV extern
	municipisCarregats, err := carregarMunicipis("codiMunicipis.csv")
	if err != nil {
		myApp.ErrorLog.Printf("Error al carregar els municipis del CSV: %v", err)
		// Pots decidir si fer un panic o assignar una llista buida
	} else {
		myApp.LlistaMunicipis = municipisCarregats
		myApp.InfoLog.Printf("S'han carregat %d municipis correctament des del CSV.", len(municipisCarregats))
	}
	//Invocació de la GUI
	myApp.makeUI()

	//Executar la APP
	myApp.MainWindow.ShowAndRun()
}

func (myApp *Config) connectSQL() (*sql.DB, error) {
	path := ""
	if os.Getenv("DB_PATH") != "" {
		//En cas de tenir valor, el recupera
		path = os.Getenv("DB_PATH")
	} else {
		path = myApp.App.Storage().RootURI().Path() + "/sql.db"
		myApp.InfoLog.Println("la base de datos está en...", path)
	}
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func (myApp *Config) setupDB(sqldb *sql.DB) {
	myApp.DB = repository.NewSQLiteRepository(sqldb)

	err := myApp.DB.Migrate()
	if err != nil {
		myApp.ErrorLog.Println(err)
		log.Panic()
	}

}
