import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable } from 'rxjs';

@Injectable({
  providedIn: 'root'
})
export class AuthService {
  private apiUrl = 'http://localhost:3001/api/v1/auth';

  constructor(private http: HttpClient) {}

  login(credentials: {email: string, password: string, totp_code?: string}): Observable<any> {
    return this.http.post(`${this.apiUrl}/login`, credentials);
  }

  register(userData: {nombre: string, email: string, tipo: string, password: string}): Observable<any> {
    return this.http.post(`${this.apiUrl}/register`, userData);
  }


  isAuthenticated(): boolean {
    const token = localStorage.getItem('access_token');
    return !!token;
  }


  getToken(): string | null {
    return localStorage.getItem('access_token');
  }


  logout(): void {
    localStorage.removeItem('access_token');
    localStorage.removeItem('refresh_token');
  }
}