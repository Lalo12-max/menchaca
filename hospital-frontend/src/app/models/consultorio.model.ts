export interface Consultorio {
  id_consultorio?: number;
  nombre: string;
  tipo: string;
  ubicacion?: string;
  medico_id?: number;
}

export interface CreateConsultorioRequest {
  nombre: string;
  tipo: string;
  ubicacion?: string;
  medico_id?: number;
}

export interface CreateConsultorioRequest {
  nombre: string;
  ubicacion?: string;
  capacidad?: number;
  equipamiento?: string;
}