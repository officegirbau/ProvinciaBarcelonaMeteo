package main

import (
	"ecohortapp/repository"
	"sort"
	"strconv"
	"time"

	"fyne.io/fyne/v2"
	_ "fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func (myApp *Config) getToolBar(window fyne.Window) *widget.Toolbar {
	toolBar := widget.NewToolbar(
		widget.NewToolbarSpacer(),
		widget.NewToolbarAction(theme.DocumentCreateIcon(), func() {
			myApp.addRegistresDialog()
		}),
		widget.NewToolbarAction(theme.ViewRefreshIcon(), func() {
			myApp.actualitzarClimaDadesContent()
		}),
		widget.NewToolbarAction(theme.SettingsIcon(), func() {
			w := myApp.mostrarPreferencies()
			w.Resize(fyne.NewSize(300, 200))
			w.Show()
		}),
	)
	return toolBar
}

func (myApp *Config) addRegistresDialog() dialog.Dialog {
	var nomsMunicipis []string
	for _, m := range myApp.LlistaMunicipis {
		nomsMunicipis = append(nomsMunicipis, m.Nom)
	}
	sort.Strings(nomsMunicipis) // Ordenem alfabèticament

	municipiSelect := widget.NewSelect(nomsMunicipis, func(s string) {
		// Actualitzem el codi quan l'usuari selecciona un municipi pel seu nom
		for _, m := range myApp.LlistaMunicipis {
			if m.Nom == s {
				MunicipiD = m.Codi
				break
			}
		}
	})

	if len(nomsMunicipis) > 0 {
		municipiSelect.SetSelected(nomsMunicipis[0])
		// Assignem el codi del primer element per defecte
		for _, m := range myApp.LlistaMunicipis {
			if m.Nom == nomsMunicipis[0] {
				MunicipiD = m.Codi
				break
			}
		}
	}

	dataRegistreEntrada := widget.NewEntry()
	precipitacioEntrada := widget.NewEntry()
	tempMaximaEntrada := widget.NewEntry()
	tempMinimaEntrada := widget.NewEntry()
	humitatEntrada := widget.NewEntry()

	// Validacions...
	validacioData := func(s string) error {
		if _, err := time.Parse("2006-01-02", s); err != nil {
			return err
		}
		return nil
	}
	dataRegistreEntrada.Validator = validacioData

	esIntValidador := func(s string) error {
		_, err := strconv.Atoi(s)
		return err
	}
	precipitacioEntrada.Validator = esIntValidador
	tempMaximaEntrada.Validator = esIntValidador
	tempMinimaEntrada.Validator = esIntValidador
	humitatEntrada.Validator = esIntValidador

	dataRegistreEntrada.PlaceHolder = "YYYY-MM-DD"

	addForm := dialog.NewForm(
		"Afegir Registre",
		"Afegir",
		"Cancelar",
		[]*widget.FormItem{
			{Text: "Municipi", Widget: municipiSelect},
			{Text: "Data Registre", Widget: dataRegistreEntrada},
			{Text: "Probabilitat de precipitació", Widget: precipitacioEntrada},
			{Text: "Temperatura màxima", Widget: tempMaximaEntrada},
			{Text: "Temperatura mínima", Widget: tempMinimaEntrada},
			{Text: "Humitat", Widget: humitatEntrada},
		},
		func(valid bool) {
			if valid {

				municipi := municipiSelect.Selected
				dataRegistre, _ := time.Parse("2006-01-02", dataRegistreEntrada.Text)
				precipitacio, _ := strconv.Atoi(precipitacioEntrada.Text)
				tempMaxima, _ := strconv.Atoi(tempMaximaEntrada.Text)
				tempMinima, _ := strconv.Atoi(tempMinimaEntrada.Text)
				humitat, _ := strconv.Atoi(humitatEntrada.Text)

				_, err := myApp.DB.AgefirRegistre(repository.Registres{
					Municipi:     municipi,
					Data:         dataRegistre,
					Precipitacio: precipitacio,
					TempMaxima:   tempMaxima,
					TempMinima:   tempMinima,
					Humitat:      humitat,
				})
				if err != nil {
					myApp.ErrorLog.Println(err)
				}

				myApp.actualitzarRegistresTable()
			}
		},
		myApp.MainWindow,
	)

	addForm.Resize(fyne.Size{Width: 400})
	addForm.Show()

	return addForm
}

func (myApp *Config) mostrarPreferencies() fyne.Window {
	// 1. Creem una finestra nova en lloc d'un diàleg
	win := myApp.App.NewWindow("Configurar ajustaments")

	// Obtenim els noms dels municipis de la llista carregada del CSV
	var nomsMunicipis []string
	for _, m := range myApp.LlistaMunicipis {
		nomsMunicipis = append(nomsMunicipis, m.Nom)
	}
	sort.Strings(nomsMunicipis)

	// Creem el quadre desplegable (Select)
	selectMunicipi := widget.NewSelect(nomsMunicipis, func(s string) {})

	// Seleccionem per defecte el municipi actual si coincideix amb el codi
	for _, m := range myApp.LlistaMunicipis {
		if m.Codi == MunicipiD {
			selectMunicipi.SetSelected(m.Nom)
			break
		}
	}

	// Creem el formulari amb el desplegable
	form := widget.NewForm(
		widget.NewFormItem("Municipi", selectMunicipi),
	)

	// Definim l'acció quan es clica Guardar (Submit)
	form.OnSubmit = func() {
		if nomSeleccionat := selectMunicipi.Selected; nomSeleccionat != "" {
			// Busquem el codi corresponent al nom seleccionat
			for _, m := range myApp.LlistaMunicipis {
				if m.Nom == nomSeleccionat {
					MunicipiD = m.Codi
					break
				}
			}

			// Guardem a les preferències
			myApp.App.Preferences().SetString("municipi", MunicipiD)

			// Refresquem les dades i tanquem la finestra
			myApp.actualitzarClimaDadesContent()
			win.Close()
		}
	}

	// Definim l'acció quan es clica Cancel·lar
	form.OnCancel = func() {
		win.Close()
	}

	// Assignem el formulari com a contingut de la finestra i la dimensionem
	win.SetContent(form)
	win.Resize(fyne.NewSize(400, 200))

	return win
}
