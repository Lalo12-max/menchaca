export enum TipoUsuario {
  PACIENTE = 'paciente',
  MEDICO = 'medico',
  ENFERMERA = 'enfermera',
  ADMIN = 'admin'
}

export interface Usuario {
  id_usuario?: number;
  nombre: string;
  email?: string;
  tipo: TipoUsuario;
  role_id?: number;
  mfa_enabled?: boolean;
  created_at?: Date;
  updated_at?: Date;
}

export interface CreateUsuarioRequest {
  nombre: string;
  email?: string;
  password: string; 
  tipo: TipoUsuario;
  role_id?: number;
}

export interface LoginResponse {
  access_token?: string;
  refresh_token?: string;
  user?: Usuario;
  requires_mfa?: boolean;
  needs_mfa_setup?: boolean;
  user_id?: number;
  message?: string;
}

export interface MFASetupResponse {
  secret: string;
  qr_code_url: string;
  message: string;
  instructions: string;
}

export interface MFAVerifyRequest {
  secret: string;
  totp_code: string;
}