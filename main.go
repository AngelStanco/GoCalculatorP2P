package main

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"bufio"
	"math"
	"fmt"
	"io"
	"os"
)

// Funciones

// ##################
// ### Porcentaje ###
// ##################

func Percent(amount float64, percent float64) float64{
	return (amount * percent) / 100
}


// ################
// ### Redondeo ###
// ################

func Round(roundParamt string, num float64, decimal int8) float64 {
	factor := math.Pow(10, float64(decimal))
	switch roundParamt {
	case "floor":
		return math.Floor(num*factor) / factor
	case "ceil":
		return math.Ceil(num*factor) / factor
	case "Trunc":
		return math.Trunc(num*factor) / factor
	default:
		return math.Round(num*factor) / factor
	}
}


// ################
// ### Get APIs ###
// ################

func GetApi(url string) ([]byte, error){
  // Realizar la solicitud HTTP GET
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	// Cierre del cuerpo de la respuesta
	defer resp.Body.Close()

	// Leemos el cuerpo de la respuesta
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}


// ###################
// ### Input Float ###
// ###################

func InputFloat(initMsg string) float64{
	reader := bufio.NewReader(os.Stdin) // Crear un lector de entrada

	for {
		fmt.Print(initMsg)
		// Leer la entrada del usuario
		input, _ := reader.ReadString('\n')
		// Eliminar espacios y saltos de línea
		input = strings.TrimSpace(input)

		num, err := strconv.ParseFloat(input, 64)
		if err != nil {
			// Si hay un error, mostrar un mensaje y volver a pedir los datos
			fmt.Println("\n ¡Error! Debes ingresar un número válido.\n")
			continue
		}
		return num
	}
}

// Estructuras


// ######################
// ### JSON Bolívares ###
// ######################

type DataTime struct {
	Date string	`json:"date"`
	Time string	`json:"time"`
}

type Monitor struct {
	Symbol string `json:"symbol"`
	Price float64 `json:"price"`
	PriceOld float32 `json:"price_old"`
	Percent float32 `json:"percent"`
	LastUp string `json:"last_update"`
}

type Monitors struct {
	BCV Monitor `json:"bcv"`
	Prl Monitor `json:"enparalelovzla"`
}

type RespVES struct {
	DataTime DataTime `json:"datetime"`
	Monitors Monitors `json:"monitors"`
}

// Main

func main() {
	var usd bool
	var exg float64
	var mnt float64
	var cms1 float64
	var cms2 float64
	var symbol string

	// UI

	fmt.Println("\nDólares (USD) -> 0\n\nBolívares (VES) -> 1\n")
	exg = InputFloat("")

	switch exg {
		case 1:
			symbol = "Bs"

			respVES, err := GetApi("https://pydolarve.org/api/v1/dollar")
			if err != nil {
				// Error de Conexion
				exg = InputFloat("Error de Conexion,  Ingresa Manualmente el Valor Actual del Bolivar: ")
			}else {
				var jsonVES RespVES

				err = json.Unmarshal(respVES, &jsonVES)
				if err != nil {
					// Error al Parsear los Datos
					exg = InputFloat("Error al Parsear los Datos,  Ingresa Manualmente el Valor Actual del Bolivar: ")
				}else {
					// Dolar BCV
					var PriceBCV float64 = jsonVES.Monitors.BCV.Price
					var pctBCV string = fmt.Sprintf("%s %g%%", jsonVES.Monitors.BCV.Symbol, jsonVES.Monitors.BCV.Percent)

					fmt.Println("\n---- BCV ---- (X No Usado X)\n\nÚltima Actualización:", jsonVES.Monitors.BCV.LastUp, "\n\nValor actual del Bolívar:", (1 / PriceBCV), "\n\nPrecio actual:", PriceBCV, pctBCV, "\n\nPrecio anterior:", jsonVES.Monitors.BCV.PriceOld)

					// Dolar Paralelo
					var PricePrl float64 = jsonVES.Monitors.Prl.Price
					var pctPrl string = fmt.Sprintf("%s %g%%", jsonVES.Monitors.Prl.Symbol, jsonVES.Monitors.Prl.Percent)

					fmt.Println("\n\n---- EnParaleloVzla ----\n\nÚltima Actualización:", jsonVES.Monitors.Prl.LastUp, "\n\nValor actual del Bolívar:", (1 / PricePrl), "\n\nPrecio actual:", PricePrl, pctPrl, "\n\nPrecio anterior:", jsonVES.Monitors.Prl.PriceOld)

					exg = PricePrl
				}
			}
		default:
			symbol = "$"
			usd = true
			exg = 1
	}


		// INPUTS


	cms1 = InputFloat("\nPorcentaje de Comisión de la Plataforma P2P: ")

	cms2 = InputFloat("\n\nPorcentaje de Comisión de la Cuenta Bancaria: ")

	mnt = InputFloat("\n\nCantidad de Dólares con la que Deseas Trabajar: ")


		// Calculos


	// Conversión de dólares a moneda seleccionada
	var mntCvs float64 = Round("", (exg * mnt), 2)

	// Porcentaje total de las Comisiones
	var pctCms float64 = cms1 + cms2

	// Monto total de las comisiones
	var mntCms float64 = Round("", Percent(mnt, pctCms), 2) * exg

	// Monto de las Comisiones en Decimal
	var dcmCms float64 = Round("", mntCms / mnt, 4)

	// Precios Minimos
	var priceBuy float64 = Round("floor", (exg - dcmCms), 3)
	var priceSale float64 = Round("ceil", (exg + dcmCms), 3)

	// Inputs Precio
	priceBuy = InputFloat(fmt.Sprintf("\n\nTu Precio de Compra - (Mínimo %g%s): ", priceBuy, symbol))
	priceSale = InputFloat(fmt.Sprintf("\n\nTu Precio de Venta - (Mínimo %g%s): ", priceSale, symbol))

	// Montos finales
	var resultSale  float64 = Round("", (priceSale * mnt) - mntCms, 2)
	var resultBuy  float64 = Round("", (priceBuy * mnt) + mntCms, 2)
	var resultFinal float64 = Round("", (mntCvs - resultBuy) + (resultSale - mntCvs), 2)


 		// RESULTADOS


	fmt.Println("\n\n------------ Resultado de las Comisiones ------------")


	fmt.Printf("\n\nComisión por cada Transacción: %g%%  ->  Perdida: %g%s", pctCms, Round("", mntCms, 2), symbol)

	fmt.Printf("\n\n\nComisión Total de Compra y Venta: %g%%  ->  Perdida Total: %g%s", (pctCms * 2), Round("", (mntCms * 2), 2), symbol)


	fmt.Println("\n\n\n------------ Resultados ------------")

	if !usd {
		fmt.Printf("\n\nConversión de 1$ a %s = %g%s  |  %g$ a %s = %g%s", symbol, exg, symbol, mnt, symbol, mntCvs, symbol)
	}

	fmt.Printf("\n\n\nTu Precio de Compra: %g%s  ->  Pagaras: %g%s y Recibirás %g USDT", priceBuy, symbol, resultBuy, symbol, mnt)

	fmt.Printf("\n\n\nTu Precio de Venta: %g%s  ->  Pagaras: %g USDT y Recibirás %g%s", priceSale, symbol, mnt, resultSale, symbol)

	fmt.Printf("\n\n\nResultado Final: %g%s", resultFinal, symbol)


		// FIN


  fmt.Println("\n\n\nPresiona doble Enter para salir...")
  fmt.Scanln()
  fmt.Scanln()
}
