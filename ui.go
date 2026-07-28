package main

import (
	"fmt"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
)

func (myApp *Config) makeUI() {
	// Cridarem la API de la AEMET
	municipiT, precipitacio, tempMax, tempMin, humitat := myApp.getClimaText()

	// 1. Fila superior: Només el municipi (alineat a l'esquerra)
	liniaMunicipi := container.NewHBox(municipiT)

	// 2. Fila inferior: Les 4 dades meteorològiques en una graella de 4 columnes
	liniaDades := container.NewGridWithColumns(4,
		precipitacio,
		tempMax,
		tempMin,
		humitat,
	)

	// 3. Contenidor principal que ajunta les dues files en vertical (VBox)
	climaDadesContent := container.NewVBox(
		liniaMunicipi,
		liniaDades,
	)

	myApp.ClimaDadesContainer = climaDadesContent

	// Obtenir la barra d'eines
	toolBar := myApp.getToolBar(myApp.MainWindow)
	grafic := myApp.pronosticTab()
	registresTabContent := myApp.registresTab()

	// Crearem les pestanyes
	pestanes := container.NewAppTabs(
		container.NewTabItemWithIcon("Previsió", theme.HomeIcon(), grafic),
		container.NewTabItemWithIcon("Registre Meteorològic", theme.InfoIcon(), registresTabContent),
	)
	pestanes.SetTabLocation(container.TabLocationTop)

	contingutFinal := container.NewVBox(climaDadesContent, toolBar, pestanes)

	myApp.MainWindow.SetContent(contingutFinal)

	// Goroutine corregida: Interval de 30 minuts i segura per a Fyne
	go func() {
		ticker := time.NewTicker(30 * time.Minute)
		defer ticker.Stop()

		for range ticker.C {
			// Fyne exigeix que qualsevol actualització visual es faci dins de fyne.Do
			fyne.Do(func() {
				myApp.actualitzarClimaDadesContent()
			})
		}
	}()
}

func (myApp *Config) actualitzarClimaDadesContent() {
	fmt.Println("S'ha actualitzat les dades")

	municipiT, precipitacio, tempMax, tempMin, humitat := myApp.getClimaText()

	// 1. Fila superior: Només el municipi
	liniaMunicipi := container.NewHBox(municipiT)

	// 2. Fila inferior: Les 4 dades meteorològiques
	liniaDades := container.NewGridWithColumns(4,
		precipitacio,
		tempMax,
		tempMin,
		humitat,
	)

	// 3. Assignem les dues files al contenidor principal (que ja és un VBox o un contenidor flexible)
	myApp.ClimaDadesContainer.Objects = []fyne.CanvasObject{liniaMunicipi, liniaDades}
	myApp.ClimaDadesContainer.Refresh()

	// Cridarem el grafic de la primera pestanya
	grafic := myApp.obtenirGrafic()
	if grafic != nil {
		myApp.PronosticGraficContainer.Objects = []fyne.CanvasObject{grafic}
		myApp.PronosticGraficContainer.Refresh()
	}
}

func (app *Config) actualitzarRegistresTable() {
	//Invoquem el mètode contenidor dels slices i ho assignem a l'atribut 	Registres del struct Config
	myApp.Registres = myApp.getRegistresSlice()

	if myApp.RegistresTable != nil {
		myApp.RegistresTable.Refresh()
		fmt.Println("3. Taula de registres refrescada gràficament!")
	}
}
