import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable } from 'rxjs';
import { Consultorio, CreateConsultorioRequest } from '../models/consultorio.model';

@Injectable({
  providedIn: 'root'
})
export class ConsultorioService {
  private apiUrl = 'http://localhost:3001/api/v1/consultorios'; 

  constructor(private http: HttpClient) {}

  getConsultorios(): Observable<Consultorio[]> {
    return this.http.get<Consultorio[]>(this.apiUrl);
  }

  getConsultorio(id: number): Observable<Consultorio> {
    return this.http.get<Consultorio>(`${this.apiUrl}/${id}`);
  }

  createConsultorio(consultorio: CreateConsultorioRequest): Observable<Consultorio> {
    return this.http.post<Consultorio>(this.apiUrl, consultorio);
  }

  updateConsultorio(id: number, consultorio: CreateConsultorioRequest): Observable<Consultorio> {
    return this.http.put<Consultorio>(`${this.apiUrl}/${id}`, consultorio);
  }

  deleteConsultorio(id: number): Observable<any> {
    return this.http.delete(`${this.apiUrl}/${id}`);
  }
}