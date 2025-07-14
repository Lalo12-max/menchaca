export interface Consulta {
  id_consulta?: number;
  consultorio_id: number;
  medico_id: number;
  paciente_id: number;
  tipo: string;
  horario: Date;
  diagnostico?: string;
  costo?: number;
}

export interface CreateConsultaRequest {
  consultorio_id: number;
  medico_id: number;
  paciente_id: number;
  tipo: string;
  horario: Date;
  diagnostico?: string;
  costo?: number;
}