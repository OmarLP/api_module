package domain

import "time"

type VerticalReport struct {
	IDCita                     *int64     `gorm:"column:id_cita" json:"id_cita,omitempty"`
	Anio                       *int       `gorm:"column:anio" json:"anio,omitempty"`
	Mes                        *int       `gorm:"column:mes" json:"mes,omitempty"`
	Dia                        *int       `gorm:"column:dia" json:"dia,omitempty"`
	Lote                       *string    `gorm:"column:lote" json:"lote,omitempty"`
	NumPag                     *int       `gorm:"column:num_pag" json:"num_pag,omitempty"`
	NumReg                     *int       `gorm:"column:num_reg" json:"num_reg,omitempty"`
	Red                        *string    `gorm:"column:red" json:"red,omitempty"`
	Microred                   *string    `gorm:"column:microred" json:"microred,omitempty"`
	PdDigitacion               *string    `gorm:"column:pd_digitacion" json:"pd_digitacion,omitempty"`
	CodigoUnico                *string    `gorm:"column:codigo_unico" json:"codigo_unico,omitempty"`
	NombreEstablecimiento      *string    `gorm:"column:nombre_establecimiento" json:"nombre_establecimiento,omitempty"`
	IDPaciente                 *string    `gorm:"column:id_paciente" json:"id_paciente,omitempty"`
	Documento                  *string    `gorm:"column:documento" json:"documento,omitempty"`
	NombrePaciente             *string    `gorm:"column:nombre_paciente" json:"nombre_paciente,omitempty"`
	Genero                     *string    `gorm:"column:genero" json:"genero,omitempty"`
	IDCondicionEstablecimiento *string    `gorm:"column:id_condicion_establecimiento" json:"id_condicion_establecimiento,omitempty"`
	IDCondicionServicio        *string    `gorm:"column:id_condicion_servicio" json:"id_condicion_servicio,omitempty"`
	CursoVida                  *string    `gorm:"column:curso_vida" json:"curso_vida,omitempty"`
	EdadReg                    *int       `gorm:"column:edad_reg" json:"edad_reg,omitempty"`
	TipoEdad                   *string    `gorm:"column:tipo_edad" json:"tipo_edad,omitempty"`
	AnioActualPaciente         *int       `gorm:"column:anio_actual_paciente" json:"anio_actual_paciente,omitempty"`
	MesActualPaciente          *int       `gorm:"column:mes_actual_paciente" json:"mes_actual_paciente,omitempty"`
	DiaActualPaciente          *int       `gorm:"column:dia_actual_paciente" json:"dia_actual_paciente,omitempty"`
	IDTurno                    *string    `gorm:"column:id_turno" json:"id_turno,omitempty"`
	CodigoItem                 *string    `gorm:"column:codigo_item" json:"codigo_item,omitempty"`
	DescripcionItem            *string    `gorm:"column:descripcion_item" json:"descripcion_item,omitempty"`
	TipoDiagnostico            *string    `gorm:"column:tipo_diagnostico" json:"tipo_diagnostico,omitempty"`
	ValorLab                   *string    `gorm:"column:valor_lab" json:"valor_lab,omitempty"`
	IDCorrelativo              *int       `gorm:"column:id_correlativo" json:"id_correlativo,omitempty"`
	IDCorrelativoLab           *int       `gorm:"column:id_correlativo_lab" json:"id_correlativo_lab,omitempty"`
	Peso                       *float64   `gorm:"column:peso;type:numeric(9,1)" json:"peso,omitempty"`
	Talla                      *float64   `gorm:"column:talla;type:numeric(9,1)" json:"talla,omitempty"`
	Hemoglobina                *float64   `gorm:"column:hemoglobina;type:numeric(9,1)" json:"hemoglobina,omitempty"`
	PerimetroAbdominal         *float64   `gorm:"column:perimetro_abdominal;type:numeric(9,1)" json:"perimetro_abdominal,omitempty"`
	PerimetroCefalico          *float64   `gorm:"column:perimetro_cefalico;type:numeric(9,1)" json:"perimetro_cefalico,omitempty"`
	DescripcionOtraCondicion   *string    `gorm:"column:descripcion_otra_condicion" json:"descripcion_otra_condicion,omitempty"`
	FechaUltimaRegla           *time.Time `gorm:"column:fecha_ultima_regla;type:date" json:"fecha_ultima_regla,omitempty"`
	FechaSolicitudHb           *time.Time `gorm:"column:fecha_solicitud_hb;type:date" json:"fecha_solicitud_hb,omitempty"`
	FechaResultadoHb           *time.Time `gorm:"column:fecha_resultado_hb" json:"fecha_resultado_hb,omitempty"`
	NombrePersonalSalud        *string    `gorm:"column:nombre_personal_salud" json:"nombre_personal_salud,omitempty"`
	NombreRegistrador          *string    `gorm:"column:nombre_registrador" json:"nombre_registrador,omitempty"`
	FechaRegistro              *time.Time `gorm:"column:fecha_registro" json:"fecha_registro,omitempty"`
	FechaModificacion          *time.Time `gorm:"column:fecha_modificacion;type:date" json:"fecha_modificacion,omitempty"`
}

func (VerticalReport) TableName() string {
	return "reporte_vertical"
}

// estructura para buscar por id_cita envuelta
type BuscarIdCita struct {
	Anio            int    `json:"anio_p" `
	MesInicio       int    `json:"mes_inicio"`
	MesFinal        int    `json:"mes_final"`
	Red             string `json:"red_p"`
	Microred        string `json:"microred_p"`
	PdDigitacion    string `json:"pd_digitacion_p"`
	Establecimiento string `json:"establecimiento_p"`
	CursoVida       string `json:"curso_vida_p"`
	CieCpm          string `json:"cie_cmp_p"`
	Documento       string `json:"documento_p"`
}
