import { Component, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule, ReactiveFormsModule, FormBuilder, FormGroup, Validators } from '@angular/forms';

import { ButtonModule } from 'primeng/button';
import { TableModule } from 'primeng/table';
import { DialogModule } from 'primeng/dialog';
import { InputTextModule } from 'primeng/inputtext';
import { DropdownModule } from 'primeng/dropdown';
import { ToastModule } from 'primeng/toast';
import { MessageModule } from 'primeng/message';
import { TagModule } from 'primeng/tag';
import { TooltipModule } from 'primeng/tooltip';
import { ConfirmDialogModule } from 'primeng/confirmdialog';
import { MessageService, ConfirmationService } from 'primeng/api';

import { Consultorio, CreateConsultorioRequest } from '../../models/consultorio.model';
import { ConsultorioService } from '../../services/consultorio.service';
import { UsuarioService } from '../../services/usuario.service';
import { Usuario, TipoUsuario } from '../../models/usuario.model';

@Component({
  selector: 'app-consultorios',
  standalone: true,
  imports: [
    CommonModule, 
    FormsModule, 
    ReactiveFormsModule,
    ButtonModule,
    TableModule,
    DialogModule,
    InputTextModule,
    DropdownModule,
    ToastModule,
    MessageModule,
    TagModule,
    TooltipModule,
    ConfirmDialogModule
  ],
  providers: [MessageService, ConfirmationService],
  templateUrl: './consultorios.component.html',
  styleUrls: ['./consultorios.component.css']
})
export class ConsultoriosComponent implements OnInit {
  consultorios: Consultorio[] = [];
  consultorioForm: FormGroup;
  editingConsultorio: Consultorio | null = null;
  showForm = false;
  loading = false;
  medicos: Usuario[] = [];
  
  tiposConsultorio = [
    { label: 'Cardiología', value: 'Cardiologia' },
    { label: 'Neurología', value: 'Neurologia' },
    { label: 'Pediatría', value: 'Pediatria' },
    { label: 'Ginecología', value: 'Ginecologia' },
    { label: 'Medicina General', value: 'Medicina General' }
  ];

  constructor(
    private consultorioService: ConsultorioService,
    private usuarioService: UsuarioService,
    private fb: FormBuilder,
    private messageService: MessageService,
    private confirmationService: ConfirmationService
  ) {
    this.consultorioForm = this.fb.group({
      nombre: ['', [Validators.required, Validators.minLength(2)]],
      tipo: ['', [Validators.required]],
      ubicacion: [''],
      medico_id: ['', [Validators.min(1)]]
    });
  }

  ngOnInit(): void {
    this.loadConsultorios();
    this.loadMedicos();
  }

  loadConsultorios(): void {
    this.loading = true;
    this.consultorioService.getConsultorios().subscribe({
      next: (consultorios) => {
        this.consultorios = consultorios;
        this.loading = false;
      },
      error: (error) => {
        console.error('Error loading consultorios:', error);
        this.messageService.add({
          severity: 'error',
          summary: 'Error',
          detail: 'Error al cargar consultorios'
        });
        this.loading = false;
      }
    });
  }

  openCreateForm(): void {
    this.editingConsultorio = null;
    this.consultorioForm.reset();
    this.showForm = true;
  }

  openEditForm(consultorio: Consultorio): void {
    this.editingConsultorio = consultorio;
    this.consultorioForm.patchValue({
      nombre: consultorio.nombre,
      tipo: consultorio.tipo,
      ubicacion: consultorio.ubicacion,
      medico_id: consultorio.medico_id
    });
    this.showForm = true;
  }

  closeForm(): void {
    this.showForm = false;
    this.editingConsultorio = null;
    this.consultorioForm.reset();
  }

  onSubmit(): void {
    if (this.consultorioForm.valid) {
      const consultorioData: CreateConsultorioRequest = {
        nombre: this.consultorioForm.value.nombre,
        tipo: this.consultorioForm.value.tipo,
        ubicacion: this.consultorioForm.value.ubicacion || undefined,
        medico_id: this.consultorioForm.value.medico_id || undefined
      };
      
      if (this.editingConsultorio) {
        this.updateConsultorio(this.editingConsultorio.id_consultorio!, consultorioData);
      } else {
        this.createConsultorio(consultorioData);
      }
    } else {
      this.messageService.add({
        severity: 'warn',
        summary: 'Formulario incompleto',
        detail: 'Por favor complete todos los campos requeridos'
      });
    }
  }

  createConsultorio(consultorioData: CreateConsultorioRequest): void {
    this.consultorioService.createConsultorio(consultorioData).subscribe({
      next: (consultorio) => {
        this.consultorios.push(consultorio);
        this.messageService.add({
          severity: 'success',
          summary: 'Éxito',
          detail: 'Consultorio creado exitosamente'
        });
        this.closeForm();
      },
      error: (error) => {
        console.error('Error creating consultorio:', error);
        this.messageService.add({
          severity: 'error',
          summary: 'Error',
          detail: error.error?.error || 'Error al crear consultorio'
        });
      }
    });
  }

  updateConsultorio(id: number, consultorioData: CreateConsultorioRequest): void {
    this.consultorioService.updateConsultorio(id, consultorioData).subscribe({
      next: (consultorio) => {
        const index = this.consultorios.findIndex(c => c.id_consultorio === id);
        if (index !== -1) {
          this.consultorios[index] = consultorio;
        }
        this.messageService.add({
          severity: 'success',
          summary: 'Éxito',
          detail: 'Consultorio actualizado exitosamente'
        });
        this.closeForm();
      },
      error: (error) => {
        console.error('Error updating consultorio:', error);
        this.messageService.add({
          severity: 'error',
          summary: 'Error',
          detail: error.error?.error || 'Error al actualizar consultorio'
        });
      }
    });
  }

  deleteConsultorio(id: number): void {
    this.confirmationService.confirm({
      message: '¿Está seguro de que desea eliminar este consultorio?',
      header: 'Confirmar eliminación',
      icon: 'pi pi-exclamation-triangle',
      acceptLabel: 'Sí, eliminar',
      rejectLabel: 'Cancelar',
      accept: () => {
        this.consultorioService.deleteConsultorio(id).subscribe({
          next: () => {
            this.consultorios = this.consultorios.filter(c => c.id_consultorio !== id);
            this.messageService.add({
              severity: 'success',
              summary: 'Éxito',
              detail: 'Consultorio eliminado exitosamente'
            });
          },
          error: (error) => {
            console.error('Error deleting consultorio:', error);
            this.messageService.add({
              severity: 'error',
              summary: 'Error',
              detail: error.error?.error || 'Error al eliminar consultorio'
            });
          }
        });
      }
    });
  }

  loadMedicos(): void {
    this.usuarioService.getUsuarios().subscribe({
      next: (usuarios) => {
        this.medicos = usuarios.filter(u => u.tipo === TipoUsuario.MEDICO);
      },
      error: (error) => {
        console.error('Error loading medicos:', error);
        this.messageService.add({
          severity: 'error',
          summary: 'Error',
          detail: 'Error al cargar médicos'
        });
      }
    });
  }

  getMedicoNombre(medicoId: number): string {
    const medico = this.medicos.find(m => m.id_usuario === medicoId);
    return medico ? medico.nombre : 'N/A';
  }

  getTipoSeverity(tipo: string): string {
    switch (tipo) {
      case 'Cardiologia':
        return 'danger';
      case 'Neurologia':
        return 'info';
      case 'Pediatria':
        return 'success';
      case 'Ginecologia':
        return 'warning';
      case 'Medicina General':
        return 'secondary';
      default:
        return 'info';
    }
  }
}