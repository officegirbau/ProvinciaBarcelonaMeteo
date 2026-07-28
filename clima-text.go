package main

import (
	"fmt"
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
)

// Fem una funció que retornara quatre elements de texte amb fyne i que es realitzarà una inferencia a l'estructura Config
func (myApp *Config) getClimaText() (*canvas.Text, *canvas.Text, *canvas.Text, *canvas.Text, *canvas.Text) {
	//Definim variables
	var diaria Diaria
	var precipitacio, tempMax, tempMin, humitat, municipiT *canvas.Text

	parte, err := diaria.GetPrevisions()
	if err != nil {
		gris := color.NRGBA{R: 155, G: 155, B: 155, A: 255}
		precipitacio = canvas.NewText("Precipitació: Sense INFO", gris)
		tempMax = canvas.NewText("Temp. Max: Sense INFO", gris)
		tempMin = canvas.NewText("Temp. Min: Sense INFO", gris)
		humitat = canvas.NewText("Humitat: Sense INFO", gris)
		municipiT = canvas.NewText("Municipi: Sense INFO", gris)
	} else {
		displayColor := color.NRGBA{R: 0, G: 180, B: 0, A: 255}

		if parte.ProbPrecipitacio < 50 {
			displayColor = color.NRGBA{R: 180, G: 0, B: 0, A: 255}
		}
		municipiTxt := fmt.Sprintf("Municipi: %s", parte.Ciutat)
		precipitacioTxt := fmt.Sprintf("Precipitació: %d%%", parte.ProbPrecipitacio)
		tempMaxTxt := fmt.Sprintf("Temp. Max: %d", parte.TemperaturaMax)
		tempMinTxt := fmt.Sprintf("Temp. Min: %d", parte.TemperaturaMin)
		humitatTxt := fmt.Sprintf("Humitat: %d%%", parte.HumitatRelativa)
		//aplicar color
		municipiT = canvas.NewText(municipiTxt, nil)
		precipitacio = canvas.NewText(precipitacioTxt, displayColor)
		tempMax = canvas.NewText(tempMaxTxt, nil)
		tempMin = canvas.NewText(tempMinTxt, nil)
		humitat = canvas.NewText(humitatTxt, displayColor)

	}
	municipiT.Alignment = fyne.TextAlignLeading
	precipitacio.Alignment = fyne.TextAlignCenter
	tempMax.Alignment = fyne.TextAlignCenter
	tempMin.Alignment = fyne.TextAlignCenter
	humitat.Alignment = fyne.TextAlignTrailing

	//Retornem les cinc variables
	return municipiT, precipitacio, tempMax, tempMin, humitat
}
