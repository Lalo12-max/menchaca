import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable } from 'rxjs';
import { Consulta, CreateConsultaRequest } from '../models/consulta.model';

@Injectable({
  providedIn: 'root'
})
export class ConsultaService {
  private apiUrl = 'http://localhost:3001/api/v1/consultas';

  constructor(private http: HttpClient) {}

  getConsultas(): Observable<Consulta[]> {
    return this.http.get<Consulta[]>(this.apiUrl);
  }

  getConsulta(id: number): Observable<Consulta> {
    return this.http.get<Consulta>(`${this.apiUrl}/${id}`);
  }

  createConsulta(consulta: CreateConsultaRequest): Observable<Consulta> {
    console.log('🌐 Enviando POST a:', this.apiUrl);
    console.log('📦 Datos a enviar:', consulta);
    return this.http.post<Consulta>(this.apiUrl, consulta);
  }

  updateConsulta(id: number, consulta: CreateConsultaRequest): Observable<Consulta> {
    return this.http.put<Consulta>(`${this.apiUrl}/${id}`, consulta);
  }

  deleteConsulta(id: number): Observable<any> {
    return this.http.delete(`${this.apiUrl}/${id}`);
  }
}