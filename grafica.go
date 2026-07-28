package main

import (
	"fmt"
	"math"
	"os"
	"time"

	"github.com/wcharczuk/go-chart/v2"
	"github.com/wcharczuk/go-chart/v2/drawing"
)

type TempBand struct {
	Min   float64
	Max   float64
	Color drawing.Color
}

type TempGradientSeries struct {
	Bands []TempBand
}

func (s TempGradientSeries) GetName() string {
	return "TempGradient"
}

func (s TempGradientSeries) GetStyle() chart.Style {
	return chart.Style{}
}

func (s TempGradientSeries) GetXAxis() chart.Axis {
	return nil
}

func (s TempGradientSeries) GetYAxis() chart.YAxisType {
	return chart.YAxisPrimary
}

func (s TempGradientSeries) PreRender(r chart.Renderer, canvasBox chart.Box, xrange, yrange chart.Range, defaults chart.Style) {
}

func (s TempGradientSeries) Validate() error {
	return nil
}

func (s TempGradientSeries) Render(r chart.Renderer, canvasBox chart.Box, xrange, yrange chart.Range, defaults chart.Style) {
	ymin := yrange.GetMin()
	ymax := yrange.GetMax()

	if ymax <= ymin {
		return
	}

	canvasHeight := float64(canvasBox.Height())

	for _, band := range s.Bands {
		// Ignorar bandes completament fora del rang visible actual
		if band.Max < ymin || band.Min > ymax {
			continue
		}

		// Mapejar proporcionalment respecte al mínim i màxim visible de l'eix Y
		fracTop := (ymax - band.Max) / (ymax - ymin)
		fracBottom := (ymax - band.Min) / (ymax - ymin)

		yTopFloat := float64(canvasBox.Top) + fracTop*canvasHeight
		yBottomFloat := float64(canvasBox.Top) + fracBottom*canvasHeight

		top := int(math.Min(yTopFloat, yBottomFloat))
		bottom := int(math.Max(yTopFloat, yBottomFloat))

		if top < canvasBox.Top {
			top = canvasBox.Top
		}
		if bottom > canvasBox.Bottom {
			bottom = canvasBox.Bottom
		}
		if top >= bottom {
			continue
		}

		r.SetFillColor(band.Color)
		r.SetStrokeColor(band.Color)
		r.SetStrokeWidth(0)

		r.MoveTo(canvasBox.Left, top)
		r.LineTo(canvasBox.Right, top)
		r.LineTo(canvasBox.Right, bottom)
		r.LineTo(canvasBox.Left, bottom)
		r.Close()
		r.Fill()
	}
}

func GenerarGraficaMultiDiaria() error {
	dies, ciutat, err := ObtenirPrevisionsMultiDiaria()
	if err != nil || len(dies) == 0 {
		return err
	}

	avuiStr := time.Now().Format("2006-01-02")
	startIndex := 0
	for idx, dItem := range dies {
		if len(dItem.Fecha) >= 10 && dItem.Fecha[:10] == avuiStr {
			startIndex = idx
			break
		}
	}

	var xValues []float64
	var yValues []float64

	var xAnnotations []float64
	var yAnnotations []float64
	var annotations []string

	limitDies := 4
	if len(dies)-startIndex < limitDies {
		limitDies = len(dies) - startIndex
	}

	for i := 0; i < limitDies; i++ {
		diaActual := dies[startIndex+i]
		offsetDia := float64(i * 24)

		pVal := 0
		if len(diaActual.ProbPrecipitacion) > 0 {
			pVal = diaActual.ProbPrecipitacion[0].Value
		}
		hVal := diaActual.HumedadRelativa.Maxima

		dataText := diaActual.Fecha
		if len(dataText) >= 10 {
			dataText = dataText[8:10] + "-" + dataText[5:7]
		}

		minT := float64(diaActual.Temperatura.Minima)
		maxT := float64(diaActual.Temperatura.Maxima)

		xAnnotations = append(xAnnotations, offsetDia+2.0)
		yAnnotations = append(yAnnotations, minT+1.5)

		labelText := fmt.Sprintf("Dia %s | P:%d%% | H:%d%%", dataText, pVal, hVal)
		annotations = append(annotations, labelText)

		for h := 0.0; h <= 24.0; h += 1.0 {
			factor := math.Sin((h / 24.0) * math.Pi)
			valY := minT + (maxT-minT)*factor

			xValues = append(xValues, offsetDia+h)
			yValues = append(yValues, valY)
		}
	}

	tempSeries := chart.ContinuousSeries{
		Name: "Temperatura (°C)",
		Style: chart.Style{
			StrokeColor: chart.ColorRed,
			StrokeWidth: 3,
		},
		XValues: xValues,
		YValues: yValues,
	}

	var chartAnnotations []chart.Value2
	for j := range xAnnotations {
		chartAnnotations = append(chartAnnotations, chart.Value2{
			XValue: xAnnotations[j],
			YValue: yAnnotations[j],
			Label:  annotations[j],
		})
	}

	annotationSeries := chart.AnnotationSeries{
		Annotations: chartAnnotations,
	}

	bands := []TempBand{
		{Min: 42, Max: 45, Color: drawing.Color{R: 139, G: 0, B: 0, A: 160}},     // Vermell molt fosc
		{Min: 39, Max: 42, Color: drawing.Color{R: 180, G: 30, B: 20, A: 160}},   // Vermell fosc
		{Min: 36, Max: 39, Color: drawing.Color{R: 217, G: 83, B: 79, A: 160}},   // Vermell
		{Min: 33, Max: 36, Color: drawing.Color{R: 240, G: 120, B: 30, A: 160}},  // Taronja intens
		{Min: 30, Max: 33, Color: drawing.Color{R: 240, G: 140, B: 40, A: 160}},  // Taronja
		{Min: 27, Max: 30, Color: drawing.Color{R: 245, G: 165, B: 50, A: 160}},  // Ambre fosc
		{Min: 24, Max: 27, Color: drawing.Color{R: 245, G: 180, B: 60, A: 160}},  // Ambre
		{Min: 21, Max: 24, Color: drawing.Color{R: 250, G: 200, B: 75, A: 160}},  // Groc intens
		{Min: 18, Max: 21, Color: drawing.Color{R: 250, G: 220, B: 90, A: 160}},  // Groc
		{Min: 15, Max: 18, Color: drawing.Color{R: 235, G: 235, B: 105, A: 160}}, // Groc verdós
		{Min: 12, Max: 15, Color: drawing.Color{R: 200, G: 230, B: 115, A: 160}}, // Verd llima clar
		{Min: 9, Max: 12, Color: drawing.Color{R: 160, G: 210, B: 140, A: 160}},  // Verd
		{Min: 6, Max: 9, Color: drawing.Color{R: 120, G: 195, B: 180, A: 160}},   // Verd blavós
		{Min: 3, Max: 6, Color: drawing.Color{R: 100, G: 180, B: 220, A: 160}},   // Blau clar
		{Min: 0, Max: 3, Color: drawing.Color{R: 70, G: 135, B: 200, A: 160}},    // Blau mitjà
		{Min: -3, Max: 0, Color: drawing.Color{R: 50, G: 110, B: 180, A: 160}},   // Blau fosc
		{Min: -5, Max: -3, Color: drawing.Color{R: 40, G: 90, B: 160, A: 160}},   // Blau molt fosc
	}

	graph := chart.Chart{
		Title: "Previsió Meteorològica - " + ciutat,
		Background: chart.Style{
			Padding: chart.Box{
				Top:    30,
				Left:   40,
				Right:  40,
				Bottom: 30,
			},
		},
		XAxis: chart.XAxis{
			Name: "Hores / Dies",
			Ticks: []chart.Tick{
				{Value: 0.0, Label: "00h"},
				{Value: 6.0, Label: "06h"},
				{Value: 12.0, Label: "12h"},
				{Value: 18.0, Label: "18h"},
				{Value: 24.0, Label: "00h"},
				{Value: 36.0, Label: "12h"},
				{Value: 48.0, Label: "00h"},
				{Value: 60.0, Label: "12h"},
				{Value: 72.0, Label: "00h"},
				{Value: 84.0, Label: "12h"},
				{Value: 96.0, Label: "24h"},
			},
		},
		Series: []chart.Series{
			TempGradientSeries{Bands: bands},
			tempSeries,
			annotationSeries,
		},
	}

	file, err := os.Create("pronostic.png")
	if err != nil {
		return err
	}
	defer file.Close()

	return graph.Render(chart.PNG, file)
}
