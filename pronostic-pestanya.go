package main

import (
	"bytes"
	"errors"
	"fmt"
	"image"
	"image/png"
	"io"
	"log"
	"os"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
)

func (myApp *Config) pronosticTab() *fyne.Container {
	grafic := myApp.obtenirGrafic()              //Invoquem el mètode obtenirGrafic() i el guardem en la variable grafic
	graficContainer := container.NewVBox(grafic) //Creem un nou contenidor Vertical a on afegim la variable grafic
	myApp.PronosticGraficContainer = graficContainer

	//Retornem el contenidor
	return graficContainer
}

func (myApp *Config) obtenirGrafic() *canvas.Image {
	var img *canvas.Image

	// Generem la gràfica multi-diaria directament
	errGen := GenerarGraficaMultiDiaria()
	if errGen != nil {
		log.Println("Error generant la gràfica de 4 dies:", errGen)
		img = canvas.NewImageFromResource(resourceNodisponiblePng)
	} else {
		img = canvas.NewImageFromFile("pronostic.png")
	}

	img.SetMinSize(fyne.Size{
		Width:  windowWidth * 0.95,
		Height: windowHeight * 0.95,
	})

	return img
}

func (myApp *Config) descarregarArxiu(URL string, nomArxiu string) error {
	//Obtenim la resposta en bytes desde la crida a una url
	response, err := myApp.HTTPClient.Get(URL)
	if err != nil {
		return err
	}

	if response.StatusCode != 200 {
		return errors.New("rebem un codi de resposta erronia quan descarreguem la imatge")
	}

	b, err := io.ReadAll(response.Body)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	//Decodifiquem la imatge en bytes per poder tractarla
	img, _, err := image.Decode(bytes.NewReader(b))
	if err != nil {
		return err
	}

	//Obtenim la sortida de l'arxiu
	out, err := os.Create(fmt.Sprintf("./%s", nomArxiu))
	if err != nil {
		return err
	}

	//Codifiquem a png transmetent els parametres de la ruta a on es creara l'arxiu i el contingut del binari
	err = png.Encode(out, img)
	if err != nil {
		return err
	}

	return nil
}
