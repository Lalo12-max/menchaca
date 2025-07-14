import { Component, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { Router, RouterModule } from '@angular/router';

import { ButtonModule } from 'primeng/button';
import { CardModule } from 'primeng/card';
import { ToastModule } from 'primeng/toast';
import { TooltipModule } from 'primeng/tooltip';
import { AvatarModule } from 'primeng/avatar';
import { DividerModule } from 'primeng/divider';
import { PanelModule } from 'primeng/panel';
import { MessageService } from 'primeng/api';

import { AuthService } from '../../services/auth.service';

@Component({
  selector: 'app-dashboard',
  standalone: true,
  imports: [
    CommonModule, 
    RouterModule,
    ButtonModule,
    CardModule,
    ToastModule,
    TooltipModule,
    AvatarModule,
    DividerModule,
    PanelModule
  ],
  providers: [MessageService],
  templateUrl: './dashboard.component.html',
  styleUrls: ['./dashboard.component.css']
})
export class DashboardComponent implements OnInit {
  
  constructor(
    private router: Router,
    private authService: AuthService,
    private messageService: MessageService
  ) {}

  ngOnInit(): void {
    const token = localStorage.getItem('access_token');
    console.log('Token en dashboard:', token ? token.substring(0, 20) + '...' : 'No hay token');
    
    if (!token) {
      console.log('No hay token, redirigiendo a login');
      this.messageService.add({
        severity: 'warn',
        summary: 'Sesión Expirada',
        detail: 'Por favor, inicia sesión nuevamente.',
        life: 3000
      });
      localStorage.removeItem('access_token');
      localStorage.removeItem('refresh_token');
      this.router.navigate(['/login']);
      return;
    }
    
    if (token.length < 10) {
      console.log('Token inválido, redirigiendo a login');
      this.messageService.add({
        severity: 'error',
        summary: 'Token Inválido',
        detail: 'El token de sesión es inválido.',
        life: 3000
      });
      localStorage.removeItem('token');
      this.router.navigate(['/login']);
      return;
    }
    
    console.log('Dashboard cargado correctamente con token válido');
    this.messageService.add({
      severity: 'success',
      summary: 'Bienvenido',
      detail: 'Dashboard cargado correctamente',
      life: 2000
    });
  }

  navigateTo(route: string): void {
    this.messageService.add({
      severity: 'info',
      summary: 'Navegando',
      detail: `Accediendo a ${route.replace('/', '')}...`,
      life: 1500
    });
    
    setTimeout(() => {
      this.router.navigate([route]);
    }, 500);
  }

  logout(): void {
    this.messageService.add({
      severity: 'info',
      summary: 'Cerrando Sesión',
      detail: 'Hasta pronto...',
      life: 2000
    });
    
    setTimeout(() => {
      this.authService.logout();
      this.router.navigate(['/login']);
    }, 1000);
  }
}