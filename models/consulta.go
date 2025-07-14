package models

import "time"

type Consulta struct {
	IDConsulta    int       `json:"id_consulta"`
	ConsultorioID int       `json:"consultorio_id"`
	MedicoID      int       `json:"medico_id"`
	PacienteID    int       `json:"paciente_id"`
	Tipo          string    `json:"tipo"`
	Horario       time.Time `json:"horario"`
	Diagnostico   *string   `json:"diagnostico"`
	Costo         *float64  `json:"costo"`
}
