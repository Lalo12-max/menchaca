import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable } from 'rxjs';
import { Usuario, CreateUsuarioRequest, MFASetupResponse, MFAVerifyRequest } from '../models/usuario.model';

@Injectable({
  providedIn: 'root'
})
export class UsuarioService {
  private apiUrl = 'http://localhost:3001/api/v1/usuarios';
  private authUrl = 'http://localhost:3001/api/v1/auth';

  constructor(private http: HttpClient) {}

  getUsuarios(): Observable<Usuario[]> {
    return this.http.get<Usuario[]>(this.apiUrl);
  }

  getUsuario(id: number): Observable<Usuario> {
    return this.http.get<Usuario>(`${this.apiUrl}/${id}`);
  }

  createUsuario(usuario: CreateUsuarioRequest): Observable<Usuario> {
    return this.http.post<Usuario>(this.apiUrl, usuario);
  }

  updateUsuario(id: number, usuario: CreateUsuarioRequest): Observable<Usuario> {
    return this.http.put<Usuario>(`${this.apiUrl}/${id}`, usuario);
  }

  deleteUsuario(id: number): Observable<any> {
    return this.http.delete(`${this.apiUrl}/${id}`);
  }


  getMFASetup(userId: number): Observable<MFASetupResponse> {
    return this.http.get<MFASetupResponse>(`${this.authUrl}/mfa/setup/${userId}`);
  }

  verifyMFASetup(userId: number, request: MFAVerifyRequest): Observable<any> {
    return this.http.post(`${this.authUrl}/mfa/setup/${userId}/verify`, request);
  }
}