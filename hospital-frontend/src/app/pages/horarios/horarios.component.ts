import { Component, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule, ReactiveFormsModule, FormBuilder, FormGroup, Validators } from '@angular/forms';

import { TableModule } from 'primeng/table';
import { ButtonModule } from 'primeng/button';
import { DialogModule } from 'primeng/dialog';
import { DropdownModule } from 'primeng/dropdown';
import { MessageModule } from 'primeng/message';
import { ToastModule } from 'primeng/toast';
import { TagModule } from 'primeng/tag';
import { TooltipModule } from 'primeng/tooltip';
import { ConfirmDialogModule } from 'primeng/confirmdialog';
import { MessageService, ConfirmationService } from 'primeng/api';

import { Horario, CreateHorarioRequest } from '../../models/horario.model';
import { HorarioService } from '../../services/horario.service';
import { UsuarioService } from '../../services/usuario.service';
import { ConsultorioService } from '../../services/consultorio.service';
import { ConsultaService } from '../../services/consulta.service';
import { Usuario, TipoUsuario } from '../../models/usuario.model';
import { Consultorio } from '../../models/consultorio.model';
import { Consulta } from '../../models/consulta.model';

interface DropdownOption {
  label: string;
  value: any;
}

@Component({
  selector: 'app-horarios',
  standalone: true,
  imports: [
    CommonModule, 
    FormsModule, 
    ReactiveFormsModule,
    TableModule,
    ButtonModule,
    DialogModule,
    DropdownModule,
    MessageModule,
    ToastModule,
    TagModule,
    TooltipModule,
    ConfirmDialogModule
  ],
  providers: [MessageService, ConfirmationService],
  templateUrl: './horarios.component.html',
  styleUrls: ['./horarios.component.css']
})
export class HorariosComponent implements OnInit {
  horarios: Horario[] = [];
  horarioForm: FormGroup;
  editingHorario: Horario | null = null;
  showForm = false;
  loading = false;
  medicos: Usuario[] = [];
  consultorios: Consultorio[] = [];
  consultas: Consulta[] = [];
  errorMessage = '';
  successMessage = '';

  consultorioOptions: DropdownOption[] = [];
  medicoOptions: DropdownOption[] = [];
  consultaOptions: DropdownOption[] = [];
  turnoOptions: DropdownOption[] = [
    { label: 'Mañana', value: 'mañana' },
    { label: 'Tarde', value: 'tarde' },
    { label: 'Noche', value: 'noche' }
  ];

  constructor(
    private horarioService: HorarioService,
    private usuarioService: UsuarioService,
    private consultorioService: ConsultorioService,
    private consultaService: ConsultaService,
    private fb: FormBuilder,
    private messageService: MessageService,
    private confirmationService: ConfirmationService
  ) {
    this.horarioForm = this.fb.group({
      consultorio_id: ['', [Validators.required]],
      medico_id: ['', [Validators.required]],
      turno: ['', [Validators.required]],
      consulta_id: ['']
    });
  }

  ngOnInit(): void {
    this.loadHorarios();
    this.loadMedicos();
    this.loadConsultorios();
    this.loadConsultas();
  }

  loadHorarios(): void {
    this.loading = true;
    this.horarioService.getHorarios().subscribe({
      next: (horarios) => {
        this.horarios = horarios;
        this.loading = false;
        if (horarios.length === 0) {
          this.messageService.add({
            severity: 'info',
            summary: 'Sin datos',
            detail: 'No se encontraron horarios en el sistema'
          });
        }
      },
      error: (error) => {
        console.error('Error loading horarios:', error);
        this.errorMessage = 'Error al cargar horarios';
        this.messageService.add({
          severity: 'error',
          summary: 'Error',
          detail: 'No se pudieron cargar los horarios'
        });
        this.loading = false;
      }
    });
  }

  loadMedicos(): void {
    this.usuarioService.getUsuarios().subscribe({
      next: (usuarios) => {
        this.medicos = usuarios.filter(u => u.tipo === TipoUsuario.MEDICO);
        this.medicoOptions = this.medicos.map(medico => ({
          label: medico.nombre,
          value: medico.id_usuario
        }));
      },
      error: (error) => {
        console.error('Error loading medicos:', error);
        this.messageService.add({
          severity: 'error',
          summary: 'Error',
          detail: 'No se pudieron cargar los médicos'
        });
      }
    });
  }

  loadConsultorios(): void {
    this.consultorioService.getConsultorios().subscribe({
      next: (consultorios) => {
        this.consultorios = consultorios;
        this.consultorioOptions = consultorios.map(consultorio => ({
          label: consultorio.nombre,
          value: consultorio.id_consultorio
        }));
      },
      error: (error) => {
        console.error('Error loading consultorios:', error);
        this.messageService.add({
          severity: 'error',
          summary: 'Error',
          detail: 'No se pudieron cargar los consultorios'
        });
      }
    });
  }

  loadConsultas(): void {
    this.consultaService.getConsultas().subscribe({
      next: (consultas) => {
        this.consultas = consultas;
        this.consultaOptions = consultas.map(consulta => ({
          label: `${consulta.tipo} - ${consulta.diagnostico || 'Sin diagnóstico'}`,
          value: consulta.id_consulta
        }));
      },
      error: (error) => {
        console.error('Error loading consultas:', error);
        this.messageService.add({
          severity: 'error',
          summary: 'Error',
          detail: 'No se pudieron cargar las consultas'
        });
      }
    });
  }

  openCreateForm(): void {
    this.editingHorario = null;
    this.horarioForm.reset();
    this.showForm = true;
    this.clearMessages();
  }

  openEditForm(horario: Horario): void {
    this.editingHorario = horario;
    this.horarioForm.patchValue({
      consultorio_id: horario.consultorio_id,
      medico_id: horario.medico_id,
      turno: horario.turno,
      consulta_id: horario.consulta_id
    });
    this.showForm = true;
    this.clearMessages();
  }

  closeForm(): void {
    this.showForm = false;
    this.editingHorario = null;
    this.horarioForm.reset();
    this.clearMessages();
  }

  onSubmit(): void {
    if (this.horarioForm.valid) {
      const formValue = this.horarioForm.value;
      
      const horarioData: CreateHorarioRequest = {
        consultorio_id: parseInt(formValue.consultorio_id, 10),
        medico_id: parseInt(formValue.medico_id, 10),
        turno: formValue.turno,
        consulta_id: formValue.consulta_id ? parseInt(formValue.consulta_id, 10) : undefined
      };
      
      if (isNaN(horarioData.consultorio_id)) {
        this.messageService.add({
          severity: 'warn',
          summary: 'Validación',
          detail: 'Debe seleccionar un consultorio válido'
        });
        return;
      }
      
      if (isNaN(horarioData.medico_id)) {
        this.messageService.add({
          severity: 'warn',
          summary: 'Validación',
          detail: 'Debe seleccionar un médico válido'
        });
        return;
      }
      
      if (this.editingHorario) {
        this.updateHorario(this.editingHorario.id_horario!, horarioData);
      } else {
        this.createHorario(horarioData);
      }
    } else {
      this.messageService.add({
        severity: 'warn',
        summary: 'Formulario inválido',
        detail: 'Por favor complete todos los campos requeridos'
      });
    }
  }

  createHorario(horarioData: CreateHorarioRequest): void {
    this.loading = true;
    this.horarioService.createHorario(horarioData).subscribe({
      next: (horario) => {
        this.horarios.push(horario);
        this.messageService.add({
          severity: 'success',
          summary: 'Éxito',
          detail: 'Horario creado correctamente'
        });
        this.closeForm();
        this.loading = false;
      },
      error: (error) => {
        console.error('Error creating horario:', error);
        this.messageService.add({
          severity: 'error',
          summary: 'Error',
          detail: error.error?.error || 'No se pudo crear el horario'
        });
        this.loading = false;
      }
    });
  }

  updateHorario(id: number, horarioData: CreateHorarioRequest): void {
    this.loading = true;
    this.horarioService.updateHorario(id, horarioData).subscribe({
      next: (horario) => {
        const index = this.horarios.findIndex(h => h.id_horario === id);
        if (index !== -1) {
          this.horarios[index] = horario;
        }
        this.messageService.add({
          severity: 'success',
          summary: 'Éxito',
          detail: 'Horario actualizado correctamente'
        });
        this.closeForm();
        this.loading = false;
      },
      error: (error) => {
        console.error('Error updating horario:', error);
        this.messageService.add({
          severity: 'error',
          summary: 'Error',
          detail: error.error?.error || 'No se pudo actualizar el horario'
        });
        this.loading = false;
      }
    });
  }

  confirmDelete(id: number): void {
    this.confirmationService.confirm({
      message: '¿Está seguro de que desea eliminar este horario?',
      header: 'Confirmar eliminación',
      icon: 'pi pi-exclamation-triangle',
      acceptLabel: 'Sí, eliminar',
      rejectLabel: 'Cancelar',
      acceptButtonStyleClass: 'p-button-danger',
      rejectButtonStyleClass: 'p-button-text',
      accept: () => {
        this.deleteHorario(id);
      }
    });
  }

  deleteHorario(id: number): void {
    this.loading = true;
    this.horarioService.deleteHorario(id).subscribe({
      next: () => {
        this.horarios = this.horarios.filter(h => h.id_horario !== id);
        this.messageService.add({
          severity: 'success',
          summary: 'Éxito',
          detail: 'Horario eliminado correctamente'
        });
        this.loading = false;
      },
      error: (error) => {
        console.error('Error deleting horario:', error);
        this.messageService.add({
          severity: 'error',
          summary: 'Error',
          detail: error.error?.error || 'No se pudo eliminar el horario'
        });
        this.loading = false;
      }
    });
  }

  getMedicoNombre(medicoId: number): string {
    const medico = this.medicos.find(m => m.id_usuario === medicoId);
    return medico ? medico.nombre : 'N/A';
  }

  getConsultorioNombre(consultorioId: number): string {
    const consultorio = this.consultorios.find(c => c.id_consultorio === consultorioId);
    return consultorio ? consultorio.nombre : 'N/A';
  }

  getTurnoSeverity(turno: string): string {
    switch (turno.toLowerCase()) {
      case 'mañana':
        return 'success';
      case 'tarde':
        return 'warning';
      case 'noche':
        return 'info';
      default:
        return 'secondary';
    }
  }

  clearMessages(): void {
    this.errorMessage = '';
    this.successMessage = '';
  }
}