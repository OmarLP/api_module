package verticalreport

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/OmarLP/api_module/internal/domain"
	"github.com/OmarLP/api_module/pkg/meta"
	"github.com/xuri/excelize/v2"
)

type (
	Controller func(w http.ResponseWriter, r *http.Request)

	Endpoints struct {
		BuscarIdCita            Controller
		BuscaEspecifica         Controller
		GetRedes                Controller
		GetMicroredes           Controller
		GetPdDigitaciones       Controller
		GetEstablecimientos     Controller
		GetCursoVida            Controller
		ExportarBuscarIdCita    Controller
		ExportarBuscaEspecifica Controller
	}

	Response struct {
		Status int         `json:"status"`
		Data   interface{} `json:"data,omitempty"`
		Err    string      `json:"error,omitempty"`
		Meta   *meta.Meta  `json:"meta,omitempty"`
	}
)

func MakeEndpoints(s Service) Endpoints {
	return Endpoints{
		BuscarIdCita:            makeBuscarIdCitaEndpoint(s),
		BuscaEspecifica:         makeBuscaEspecificaEndpoint(s),
		GetRedes:                makeGetRedesEndpoint(s),
		GetMicroredes:           makeGetMicroredesEndpoint(s),
		GetPdDigitaciones:       makeGetPdDigitacionesEndpoint(s),
		GetEstablecimientos:     makeGetEstablecimientosEndpoint(s),
		GetCursoVida:            makeGetCursoVidaEndpoint(s),
		ExportarBuscarIdCita:    makeExportarBuscarIdCitaEndpoint(s),
		ExportarBuscaEspecifica: makeExportarBuscaEspecificaEndpoint(s),
	}
}

func makeGetCursoVidaEndpoint(s Service) Controller {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		cursoVida, err := s.GetCursoVida()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(&Response{Status: 500, Err: "internal server request"})
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(&Response{Status: 200, Data: cursoVida})
	}
}

func makeGetRedesEndpoint(s Service) Controller {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		redes, err := s.GetRedes()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(&Response{Status: 500, Err: "internal server request"})
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(&Response{Status: 200, Data: redes})
	}
}
func makeGetMicroredesEndpoint(s Service) Controller {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		red := r.URL.Query().Get("red")
		microredes, err := s.GetMicroredes(red)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(&Response{Status: 500, Err: "internal server request"})
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(&Response{Status: 200, Data: microredes})
	}
}
func makeGetPdDigitacionesEndpoint(s Service) Controller {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		red := r.URL.Query().Get("red")
		microred := r.URL.Query().Get("microred")
		pdDigitaciones, err := s.GetPdDigitaciones(red, microred)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(&Response{Status: 500, Err: "internal server request"})
			return
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(&Response{Status: 200, Data: pdDigitaciones})

	}
}
func makeGetEstablecimientosEndpoint(s Service) Controller {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		red := r.URL.Query().Get("red")
		microred := r.URL.Query().Get("microred")
		pdDigitacion := r.URL.Query().Get("pd_digitacion")
		establecimientos, err := s.GetEstablecimientos(red, microred, pdDigitacion)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(&Response{Status: 500, Err: "internal server request"})
			return
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(&Response{Status: 200, Data: establecimientos})
	}
}

func makeBuscarIdCitaEndpoint(s Service) Controller {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		var req domain.BuscarIdCita
		// depurar: ver como llega el body
		// validar request de acuerdo a la estructura solicitada
		decoder := json.NewDecoder(r.Body)
		decoder.DisallowUnknownFields()

		if err := decoder.Decode(&req); err != nil {
			// imprime el error de decode en la terminal
			fmt.Printf("error decode json: %v\n", err)
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(&Response{Status: 400, Err: "invalid request format"})
			return
		}

		// imprimir el struct recibido
		fmt.Printf("req recibido en GO: %+v\n", req)

		//validar anio
		currentYear := time.Now().Year()
		if req.Anio < 2000 || req.Anio > currentYear {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(&Response{Status: 400, Err: "year must be between 2000 and 2026"})
			return
		}

		// el mes de inicio no puede ser mayor que el mes final
		if req.MesInicio > req.MesFinal {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(&Response{Status: 400, Err: "start month cannot be greater than end month"})
			return
		}

		// leer parametros de la url (?page=2&limit=10)
		/*v := r.URL.Query()
		limit, _ := strconv.ParseInt(v.Get("limit"), 10, 64)
		page, _ := strconv.ParseInt(v.Get("page"), 10, 64)*/

		// leer parametros de la url con valores por defecto
		v := r.URL.Query()
		limit, err := strconv.ParseInt(v.Get("limit"), 10, 64)
		if err != nil || limit <= 0 {
			limit = 22 // Valor por defecto
		}

		page, err := strconv.ParseInt(v.Get("page"), 10, 64)
		if err != nil || page <= 0 {
			page = 1 // Valor por defecto
		}

		// obtener el total de registros
		count, err := s.Count(req)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(&Response{Status: 500, Err: err.Error()})
			return
		}

		// inicio de metadata de paginarcion
		meta, err := meta.New(int64(page), int64(limit), int64(count))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(&Response{Status: 500, Err: err.Error()})
			return
		}

		// respuesta con req, offset y limit
		data, err := s.BuscarIdCita(req, meta.Offset(), meta.Limit())
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(&Response{Status: 500, Err: err.Error()})
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(&Response{Status: 200, Data: data, Meta: meta})

	}
}

func makeExportarBuscarIdCitaEndpoint(s Service) Controller {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		var req domain.BuscarIdCita
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// Llamamos al servicio para obtener todos los registros filtrados
		registros, err := s.ExportarBuscarIdCita(req)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		// Creamos el archivo Excel en memoria
		f := excelize.NewFile()
		sheetName := "Reporte Vertical"
		f.SetSheetName("Sheet1", sheetName)

		// Escribimos los encabezados de las columnas
		headers := []string{
			"ID Cita", "Año", "Mes", "Día", "Lote", "Pág", "Reg", "Código EESS",
			"Establecimiento", "ID Paciente", "Documento", "Paciente", "Género",
			"Cond. EESS", "Cond. Serv.", "Edad", "Tipo Edad", "Item", "Tipo Dx",
			"Valor Lab", "Correlativo", "Correlativo Lab", "Peso", "Talla",
			"Hemoglobina", "P. Abdominal", "P. Cefálico", "Condición", "FUR",
			"Personal Salud", "Registrador", "F. Registro", "F. Modificación",
		}

		for i, header := range headers {
			cell, _ := excelize.CoordinatesToCellName(i+1, 1)
			f.SetCellValue(sheetName, cell, header)
		}

		// Escribimos los datos fila por fila
		for rowIndex, item := range registros {
			row := rowIndex + 2 // Empezamos en la fila 2 (la 1 es de encabezados)

			f.SetCellValue(sheetName, fmt.Sprintf("A%d", row), safeInt64(item.IDCita))
			f.SetCellValue(sheetName, fmt.Sprintf("B%d", row), safeInt(item.Anio))
			f.SetCellValue(sheetName, fmt.Sprintf("C%d", row), safeInt(item.Mes))
			f.SetCellValue(sheetName, fmt.Sprintf("D%d", row), safeInt(item.Dia))
			f.SetCellValue(sheetName, fmt.Sprintf("E%d", row), safeString(item.Lote))
			f.SetCellValue(sheetName, fmt.Sprintf("F%d", row), safeInt(item.NumPag))
			f.SetCellValue(sheetName, fmt.Sprintf("G%d", row), safeInt(item.NumReg))
			f.SetCellValue(sheetName, fmt.Sprintf("H%d", row), safeString(item.CodigoUnico))
			f.SetCellValue(sheetName, fmt.Sprintf("I%d", row), safeString(item.NombreEstablecimiento))
			f.SetCellValue(sheetName, fmt.Sprintf("J%d", row), safeString(item.IDPaciente))
			f.SetCellValue(sheetName, fmt.Sprintf("K%d", row), safeString(item.Documento))
			f.SetCellValue(sheetName, fmt.Sprintf("L%d", row), safeString(item.NombrePaciente))
			f.SetCellValue(sheetName, fmt.Sprintf("M%d", row), safeString(item.Genero))
			f.SetCellValue(sheetName, fmt.Sprintf("N%d", row), safeString(item.IDCondicionEstablecimiento))
			f.SetCellValue(sheetName, fmt.Sprintf("O%d", row), safeString(item.IDCondicionServicio))
			f.SetCellValue(sheetName, fmt.Sprintf("P%d", row), safeInt(item.EdadReg))
			f.SetCellValue(sheetName, fmt.Sprintf("Q%d", row), safeString(item.TipoEdad))
			f.SetCellValue(sheetName, fmt.Sprintf("R%d", row), safeString(item.CodigoItem))
			f.SetCellValue(sheetName, fmt.Sprintf("S%d", row), safeString(item.TipoDiagnostico))
			f.SetCellValue(sheetName, fmt.Sprintf("T%d", row), safeString(item.ValorLab))
			f.SetCellValue(sheetName, fmt.Sprintf("U%d", row), safeInt(item.IDCorrelativo))
			f.SetCellValue(sheetName, fmt.Sprintf("V%d", row), safeInt(item.IDCorrelativoLab))
			f.SetCellValue(sheetName, fmt.Sprintf("W%d", row), safeFloat(item.Peso))
			f.SetCellValue(sheetName, fmt.Sprintf("X%d", row), safeFloat(item.Talla))
			f.SetCellValue(sheetName, fmt.Sprintf("Y%d", row), safeFloat(item.Hemoglobina))
			f.SetCellValue(sheetName, fmt.Sprintf("Z%d", row), safeFloat(item.PerimetroAbdominal))
			f.SetCellValue(sheetName, fmt.Sprintf("AA%d", row), safeFloat(item.PerimetroCefalico))
			f.SetCellValue(sheetName, fmt.Sprintf("AB%d", row), safeString(item.DescripcionOtraCondicion))
			f.SetCellValue(sheetName, fmt.Sprintf("AC%d", row), safeTime(item.FechaUltimaRegla))
			f.SetCellValue(sheetName, fmt.Sprintf("AD%d", row), safeString(item.NombrePersonalSalud))
			f.SetCellValue(sheetName, fmt.Sprintf("AE%d", row), safeString(item.NombreRegistrador))
			f.SetCellValue(sheetName, fmt.Sprintf("AF%d", row), safeTime(item.FechaRegistro))
			f.SetCellValue(sheetName, fmt.Sprintf("AG%d", row), safeTime(item.FechaModificacion))
		}

		// Configuramos los headers de la respuesta HTTP para que el navegador descargue el archivo
		w.Header().Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
		w.Header().Set("Content-Disposition", "attachment; filename=Reporte_Vertical.xlsx")
		w.WriteHeader(http.StatusOK)

		// Escribimos el buffer del excel directamente al stream de respuesta
		f.Write(w)
	}
}

func makeBuscaEspecificaEndpoint(s Service) Controller {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		var req domain.BuscaEspecifica
		decoder := json.NewDecoder(r.Body)
		decoder.DisallowUnknownFields()

		if err := decoder.Decode(&req); err != nil {
			fmt.Printf("error decode json: %v\n", err)
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(&Response{Status: 400, Err: "invalid request format"})
			return
		}

		//imprimir el struct recibido
		fmt.Printf("req recibido en Go: %v\n", req)

		//validar anio
		currentYear := time.Now().Year()
		if req.Anio < 2000 || req.Anio > currentYear {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(&Response{Status: 400, Err: "year must be between 2000 and 2026"})
			return
		}

		// el mes de inicio no puede ser mayor que el mes final
		if req.MesInicio > req.MesFinal {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(&Response{Status: 400, Err: "start month cannot be greater than end month"})
			return
		}

		//leer parametros de la url con valores or defecto
		v := r.URL.Query()
		limit, err := strconv.ParseInt(v.Get("limit"), 10, 64)
		if err != nil || limit <= 0 {
			limit = 22
		}

		page, err := strconv.ParseInt(v.Get("page"), 10, 64)
		if err != nil || page <= 0 {
			page = 1 // Valor por defecto
		}

		// obtener el total de registros
		count, err := s.CountBE(req)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(&Response{Status: 500, Err: err.Error()})
			return
		}

		// inicio de metadata de paginarcion
		meta, err := meta.New(int64(page), int64(limit), int64(count))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(&Response{Status: 500, Err: err.Error()})
			return
		}

		// respuesta con req, offset y limit
		data, err := s.BuscaEspecifica(req, meta.Offset(), meta.Limit())
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(&Response{Status: 500, Err: err.Error()})
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(&Response{Status: 200, Data: data, Meta: meta})

	}
}

func makeExportarBuscaEspecificaEndpoint(s Service) Controller {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		var req domain.BuscaEspecifica
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// Llamamos al servicio para obtener todos los registros filtrados
		registros, err := s.ExportarBuscaEspecifica(req)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		// Creamos el archivo Excel en memoria
		f := excelize.NewFile()
		sheetName := "Reporte Vertical"
		f.SetSheetName("Sheet1", sheetName)

		// Escribimos los encabezados de las columnas
		headers := []string{
			"ID Cita", "Año", "Mes", "Día", "Lote", "Pág", "Reg", "Código EESS",
			"Establecimiento", "ID Paciente", "Documento", "Paciente", "Género",
			"Cond. EESS", "Cond. Serv.", "Edad", "Tipo Edad", "Item", "Tipo Dx",
			"Valor Lab", "Correlativo", "Correlativo Lab", "Peso", "Talla",
			"Hemoglobina", "P. Abdominal", "P. Cefálico", "Condición", "FUR",
			"Personal Salud", "Registrador", "F. Registro", "F. Modificación",
		}

		for i, header := range headers {
			cell, _ := excelize.CoordinatesToCellName(i+1, 1)
			f.SetCellValue(sheetName, cell, header)
		}

		// Escribimos los datos fila por fila
		for rowIndex, item := range registros {
			row := rowIndex + 2 // Empezamos en la fila 2 (la 1 es de encabezados)

			f.SetCellValue(sheetName, fmt.Sprintf("A%d", row), safeInt64(item.IDCita))
			f.SetCellValue(sheetName, fmt.Sprintf("B%d", row), safeInt(item.Anio))
			f.SetCellValue(sheetName, fmt.Sprintf("C%d", row), safeInt(item.Mes))
			f.SetCellValue(sheetName, fmt.Sprintf("D%d", row), safeInt(item.Dia))
			f.SetCellValue(sheetName, fmt.Sprintf("E%d", row), safeString(item.Lote))
			f.SetCellValue(sheetName, fmt.Sprintf("F%d", row), safeInt(item.NumPag))
			f.SetCellValue(sheetName, fmt.Sprintf("G%d", row), safeInt(item.NumReg))
			f.SetCellValue(sheetName, fmt.Sprintf("H%d", row), safeString(item.CodigoUnico))
			f.SetCellValue(sheetName, fmt.Sprintf("I%d", row), safeString(item.NombreEstablecimiento))
			f.SetCellValue(sheetName, fmt.Sprintf("J%d", row), safeString(item.IDPaciente))
			f.SetCellValue(sheetName, fmt.Sprintf("K%d", row), safeString(item.Documento))
			f.SetCellValue(sheetName, fmt.Sprintf("L%d", row), safeString(item.NombrePaciente))
			f.SetCellValue(sheetName, fmt.Sprintf("M%d", row), safeString(item.Genero))
			f.SetCellValue(sheetName, fmt.Sprintf("N%d", row), safeString(item.IDCondicionEstablecimiento))
			f.SetCellValue(sheetName, fmt.Sprintf("O%d", row), safeString(item.IDCondicionServicio))
			f.SetCellValue(sheetName, fmt.Sprintf("P%d", row), safeInt(item.EdadReg))
			f.SetCellValue(sheetName, fmt.Sprintf("Q%d", row), safeString(item.TipoEdad))
			f.SetCellValue(sheetName, fmt.Sprintf("R%d", row), safeString(item.CodigoItem))
			f.SetCellValue(sheetName, fmt.Sprintf("S%d", row), safeString(item.TipoDiagnostico))
			f.SetCellValue(sheetName, fmt.Sprintf("T%d", row), safeString(item.ValorLab))
			f.SetCellValue(sheetName, fmt.Sprintf("U%d", row), safeInt(item.IDCorrelativo))
			f.SetCellValue(sheetName, fmt.Sprintf("V%d", row), safeInt(item.IDCorrelativoLab))
			f.SetCellValue(sheetName, fmt.Sprintf("W%d", row), safeFloat(item.Peso))
			f.SetCellValue(sheetName, fmt.Sprintf("X%d", row), safeFloat(item.Talla))
			f.SetCellValue(sheetName, fmt.Sprintf("Y%d", row), safeFloat(item.Hemoglobina))
			f.SetCellValue(sheetName, fmt.Sprintf("Z%d", row), safeFloat(item.PerimetroAbdominal))
			f.SetCellValue(sheetName, fmt.Sprintf("AA%d", row), safeFloat(item.PerimetroCefalico))
			f.SetCellValue(sheetName, fmt.Sprintf("AB%d", row), safeString(item.DescripcionOtraCondicion))
			f.SetCellValue(sheetName, fmt.Sprintf("AC%d", row), safeTime(item.FechaUltimaRegla))
			f.SetCellValue(sheetName, fmt.Sprintf("AD%d", row), safeString(item.NombrePersonalSalud))
			f.SetCellValue(sheetName, fmt.Sprintf("AE%d", row), safeString(item.NombreRegistrador))
			f.SetCellValue(sheetName, fmt.Sprintf("AF%d", row), safeTime(item.FechaRegistro))
			f.SetCellValue(sheetName, fmt.Sprintf("AG%d", row), safeTime(item.FechaModificacion))
		}

		// Configuramos los headers de la respuesta HTTP para que el navegador descargue el archivo
		w.Header().Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
		w.Header().Set("Content-Disposition", "attachment; filename=Reporte_Vertical.xlsx")
		w.WriteHeader(http.StatusOK)

		// Escribimos el buffer del excel directamente al stream de respuesta
		f.Write(w)
	}
}

// funciones auxiliares para extraer el valor de los punteros
func safeInt64(p *int64) int64 {
	if p == nil {
		return 0
	}
	return *p
}

func safeInt(p *int) int {
	if p == nil {
		return 0
	}
	return *p
}

func safeFloat(f *float64) float64 {
	if f == nil {
		return 0
	}
	return *f
}

func safeString(p *string) string {
	if p == nil {
		return ""
	}
	return *p
}

func safeTime(t *time.Time) string {
	if t == nil {
		return ""
	}
	return t.Format("2006-01-02")
}
