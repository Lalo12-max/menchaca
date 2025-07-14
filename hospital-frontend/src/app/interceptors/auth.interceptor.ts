import { HttpInterceptorFn } from '@angular/common/http';
import { inject } from '@angular/core';
import { Router } from '@angular/router';
import { catchError, throwError } from 'rxjs';

export const authInterceptor: HttpInterceptorFn = (req, next) => {
  const router = inject(Router);
  const token = localStorage.getItem('access_token');

  if (token) {
    req = req.clone({
      setHeaders: {
        Authorization: `Bearer ${token}`
      }
    });
  }

  return next(req).pipe(
    catchError(error => {
      if (error.status === 401) {
        localStorage.removeItem('access_token');
        localStorage.removeItem('refresh_token');
        router.navigate(['/login']);
      } else if (error.status === 403) {
        // Manejar errores de permisos
        console.error('Permisos insuficientes:', error.error);
        // Podrías mostrar un mensaje de error o redirigir a una página de "sin permisos"
      }
      return throwError(() => error);
    })
  );
};