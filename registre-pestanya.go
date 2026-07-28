package main

import (
	"ecohortapp/repository"
	"fmt"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

// Realitzem una funció que retornara un contenidor de Fyne amb el contingut
func (myApp *Config) registresTab() *fyne.Container {
	//Invoquem la funcio anterior per carregar l'estructura de dedes amb la interficie de slice de slices
	myApp.Registres = myApp.getRegistresSlice()
	//També invoquem el mètode getRegistresTable() i l'asignem al item RegistresTable del struct
	myApp.RegistresTable = myApp.getRegistresTable()
	//Creem un contenidor amb una capça vertical i a on situem el widget que em general de la taula Registres

	registresContainer := container.NewBorder(
		nil,
		nil,
		nil,
		nil,
		//definirem un contenidor que ens permetra realizar graelles adaptatives i que indicarem amb dos parametres: el nombre de files/columnas i l’objecte que situarem.
		container.NewAdaptiveGrid(1, myApp.RegistresTable),
	)

	return registresContainer

}

// Realitzem una funcio adicional que ens retornara el punter a la widget en forma de taula i a on situarem les dades
func (myApp *Config) getRegistresTable() *widget.Table {
	// 1. Assegurem que myApp.Registres tingui dades inicialment
	myApp.Registres = myApp.getRegistresSlice()

	t := widget.NewTable(
		func() (int, int) {

			if len(myApp.Registres) == 0 {
				return 0, 0
			}
			return len(myApp.Registres), len(myApp.Registres[0])
		},
		func() fyne.CanvasObject {
			ctr := container.NewVBox(widget.NewLabel(""))
			return ctr
		},
		func(i widget.TableCellID, o fyne.CanvasObject) {

			if len(myApp.Registres) == 0 || i.Row >= len(myApp.Registres) {
				return
			}

			if i.Col == (len(myApp.Registres[0])-1) && i.Row != 0 {
				w := widget.NewButtonWithIcon("Borrar", theme.DeleteIcon(), func() {
					dialog.ShowConfirm("Borrar?", "", func(deleted bool) {
						if deleted {
							id, _ := strconv.Atoi(myApp.Registres[i.Row][0].(string))
							err := myApp.DB.BorrarRegistre(int64(id))
							if err != nil {
								myApp.ErrorLog.Println(err)
							}
						}
						// Forcem el refresc de la taula
						myApp.actualitzarRegistresTable()
					}, myApp.MainWindow)
				})
				w.Importance = widget.HighImportance
				o.(*fyne.Container).Objects = []fyne.CanvasObject{w}
			} else {
				o.(*fyne.Container).Objects = []fyne.CanvasObject{
					widget.NewLabel(myApp.Registres[i.Row][i.Col].(string)),
				}
			}
		},
	)

	// Establim l'ample de les diferents cel·les de Registres Climatològics
	colWidths := []float32{windowWidth * 0.07, windowWidth * 0.10, windowWidth * 0.26, windowWidth * 0.09, windowWidth * 0.11, windowWidth * 0.11, windowWidth * 0.1, windowWidth * 0.1}
	for i := 0; i < len(colWidths); i++ {
		t.SetColumnWidth(i, colWidths[i])
	}

	return t
}

// Realitzem una funció per obtenir tots els Registres en un Slice de Slices através d'una interficie que ens sera retornada
func (myApp *Config) getRegistresSlice() [][]interface{} {
	var slice [][]interface{}

	//Invoquem el métode inferior registresActuals()
	dades, err := myApp.registresActuals()
	if err != nil {
		myApp.ErrorLog.Println(err)
	}

	//Realitzem un append per incloure els registres obtinguts en forma de files i definint alhora les etiquetes de cada columna per la fila inicial.
	slice = append(slice, []interface{}{"ID", "Data", "Municipi", "Pluja", "Temp. Màx.", "Temp. Min.", "Humitat", "Opcions"})

	//Executem un for per elaborar tantes files com resultats ha obtingut de la BD
	for _, x := range dades {
		//Creem una interficie buida per la fila actual
		var filaActual []interface{}

		//anem afegint a la fila actual cada un dels valors que corresponen a cada columna definida al inici
		filaActual = append(filaActual, strconv.FormatInt(x.ID, 10)) //Transformem el valor numeric a String en base 10
		filaActual = append(filaActual, x.Data.Format("2006-01-02")) //Formategem la data al standard americà
		filaActual = append(filaActual, fmt.Sprintf("%s", x.Municipi))
		filaActual = append(filaActual, fmt.Sprintf("%d%%", x.Precipitacio))   //Formatagem la sortida a un valor decimal enter
		filaActual = append(filaActual, fmt.Sprintf("%d", x.TempMaxima))       //Formatagem la sortida a un valor decimal enter
		filaActual = append(filaActual, fmt.Sprintf("%d", x.TempMinima))       //Formatagem la sortida a un valor decimal enter
		filaActual = append(filaActual, fmt.Sprintf("%d%%", x.Humitat))        //Formatagem la sortida a un valor decimal enter
		filaActual = append(filaActual, widget.NewButton("Borrar", func() {})) //Definim el boto per eliminar i que invocarà una funció que ja definirem

		//Afegim aquesta fila a el slice de files
		slice = append(slice, filaActual)
	}

	return slice
}

// Realitzem una altre funció per obtenir tots els Registres amb un slice pero del nostre repositori en la DB
func (myApp *Config) registresActuals() ([]repository.Registres, error) {
	registres, err := myApp.DB.VeureTotsRegistres()
	if err != nil {
		//Capturem el possible error en el log d'errors
		myApp.ErrorLog.Println(err)
		return nil, err
	}

	return registres, nil
}
