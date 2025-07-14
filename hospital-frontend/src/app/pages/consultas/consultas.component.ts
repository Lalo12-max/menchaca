import { Component, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule, ReactiveFormsModule, FormBuilder, FormGroup, Validators } from '@angular/forms';

import { ButtonModule } from 'primeng/button';
import { TableModule } from 'primeng/table';
import { DialogModule } from 'primeng/dialog';
import { InputTextModule } from 'primeng/inputtext';
import { DropdownModule } from 'primeng/dropdown';
import { CardModule } from 'primeng/card';
import { TagModule } from 'primeng/tag';
import { TooltipModule } from 'primeng/tooltip';
import { MessageModule } from 'primeng/message';
import { ToastModule } from 'primeng/toast';
import { ConfirmDialogModule } from 'primeng/confirmdialog';
import { MessageService, ConfirmationService } from 'primeng/api';

import { Consulta, CreateConsultaRequest } from '../../models/consulta.model';
import { ConsultaService } from '../../services/consulta.service';
import { UsuarioService } from '../../services/usuario.service';
import { ConsultorioService } from '../../services/consultorio.service';
import { Usuario } from '../../models/usuario.model';
import { Consultorio } from '../../models/consultorio.model';

@Component({
  selector: 'app-consultas',
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
    CardModule,
    TagModule,
    TooltipModule,
    MessageModule,
    ToastModule,
    ConfirmDialogModule
  ],
  providers: [MessageService, ConfirmationService],
  templateUrl: './consultas.component.html',
  styleUrls: ['./consultas.component.css']
})
export class ConsultasComponent implements OnInit {
  consultas: Consulta[] = [];
  consultaForm: FormGroup;
  editingConsulta: Consulta | null = null;
  showForm = false;
  loading = false;
  medicos: Usuario[] = [];
  pacientes: Usuario[] = [];
  consultorios: Consultorio[] = [];

  consultorioOptions: any[] = [];
  medicoOptions: any[] = [];
  pacienteOptions: any[] = [];
  tipoOptions = [
    { label: 'Consulta General', value: 'consulta_general' },
    { label: 'Especialidad', value: 'especialidad' },
    { label: 'Urgencia', value: 'urgencia' },
    { label: 'Control', value: 'control' }
  ];

  constructor(
    private consultaService: ConsultaService,
    private usuarioService: UsuarioService,
    private consultorioService: ConsultorioService,
    private fb: FormBuilder,
    private messageService: MessageService,
    private confirmationService: ConfirmationService
  ) {
    this.consultaForm = this.fb.group({
      consultorio_id: ['', [Validators.required]],
      medico_id: ['', [Validators.required]],
      paciente_id: ['', [Validators.required]],
      tipo: ['', [Validators.required]],
      horario: ['', [Validators.required]],
      diagnostico: [''],
      costo: ['', [Validators.min(0)]]
    });
  }

  ngOnInit(): void {
    this.loadConsultas();
    this.loadUsuarios();
    this.loadConsultorios();
  }

  loadConsultas(): void {
    this.loading = true;
    this.consultaService.getConsultas().subscribe({
      next: (consultas) => {
        this.consultas = consultas;
        this.loading = false;
      },
      error: (error) => {
        console.error('Error loading consultas:', error);
        this.loading = false;
      }
    });
  }

  loadUsuarios(): void {
    this.usuarioService.getUsuarios().subscribe({
      next: (usuarios) => {
        this.medicos = usuarios.filter(u => u.tipo === 'medico');
        this.pacientes = usuarios.filter(u => u.tipo === 'paciente');
        
        this.medicoOptions = this.medicos.map(medico => ({
          label: medico.nombre,
          value: medico.id_usuario
        }));
        
        this.pacienteOptions = this.pacientes.map(paciente => ({
          label: paciente.nombre,
          value: paciente.id_usuario
        }));
      },
      error: (error) => {
        console.error('Error loading usuarios:', error);
        this.messageService.add({
          severity: 'error',
          summary: 'Error',
          detail: 'Error al cargar usuarios'
        });
      }
    });
  }

  loadConsultorios(): void {
    this.consultorioService.getConsultorios().subscribe({
      next: (consultorios) => {
        this.consultorios = consultorios;
        
        this.consultorioOptions = this.consultorios.map(consultorio => ({
          label: consultorio.nombre,
          value: consultorio.id_consultorio
        }));
      },
      error: (error) => {
        console.error('Error loading consultorios:', error);
        this.messageService.add({
          severity: 'error',
          summary: 'Error',
          detail: 'Error al cargar consultorios'
        });
      }
    });
  }

  openCreateForm(): void {
    this.editingConsulta = null;
    this.consultaForm.reset();
    this.showForm = true;
  }

  openEditForm(consulta: Consulta): void {
    this.editingConsulta = consulta;
    const horarioFormatted = new Date(consulta.horario).toISOString().slice(0, 16);
    this.consultaForm.patchValue({
      consultorio_id: consulta.consultorio_id,
      medico_id: consulta.medico_id,
      paciente_id: consulta.paciente_id,
      tipo: consulta.tipo,
      horario: horarioFormatted,
      diagnostico: consulta.diagnostico,
      costo: consulta.costo
    });
    this.showForm = true;
  }

  closeForm(): void {
    this.showForm = false;
    this.editingConsulta = null;
    this.consultaForm.reset();
  }

  onSubmit(): void {
    if (this.consultaForm.valid) {
      const consultaData: CreateConsultaRequest = {
        consultorio_id: +this.consultaForm.value.consultorio_id,  
        medico_id: +this.consultaForm.value.medico_id,            
        paciente_id: +this.consultaForm.value.paciente_id,       
        tipo: this.consultaForm.value.tipo,
        horario: new Date(this.consultaForm.value.horario),
        diagnostico: this.consultaForm.value.diagnostico,
        costo: this.consultaForm.value.costo ? +this.consultaForm.value.costo : undefined  
      };
      
      console.log('🚀 Enviando datos de consulta:', consultaData);
      console.log('📅 Tipo de horario:', typeof consultaData.horario);
      console.log('📅 Valor de horario:', consultaData.horario);
      
      if (this.editingConsulta) {
        this.updateConsulta(this.editingConsulta.id_consulta!, consultaData);
      } else {
        this.createConsulta(consultaData);
      }
    }
  }

  createConsulta(consultaData: CreateConsultaRequest): void {
    console.log('🚀 Enviando datos de consulta:', consultaData);
    console.log('📅 Tipo de horario:', typeof consultaData.horario);
    console.log('📅 Valor de horario:', consultaData.horario);
    
    this.consultaService.createConsulta(consultaData).subscribe({
      next: (consulta) => {
        console.log('✅ Consulta creada exitosamente:', consulta);
        this.consultas.push(consulta);
        this.closeForm();
        this.messageService.add({
          severity: 'success',
          summary: 'Éxito',
          detail: 'Consulta creada exitosamente'
        });
      },
      error: (error) => {
        console.error('❌ Error creating consulta:', error);
        console.error('📄 Error details:', error.error);
        this.messageService.add({
          severity: 'error',
          summary: 'Error',
          detail: 'Error al crear la consulta'
        });
      }
    });
  }

  updateConsulta(id: number, consultaData: CreateConsultaRequest): void {
    this.consultaService.updateConsulta(id, consultaData).subscribe({
      next: (consulta) => {
        const index = this.consultas.findIndex(c => c.id_consulta === id);
        if (index !== -1) {
          this.consultas[index] = consulta;
        }
        this.closeForm();
        this.messageService.add({
          severity: 'success',
          summary: 'Éxito',
          detail: 'Consulta actualizada exitosamente'
        });
      },
      error: (error) => {
        console.error('Error updating consulta:', error);
        this.messageService.add({
          severity: 'error',
          summary: 'Error',
          detail: 'Error al actualizar la consulta'
        });
      }
    });
  }

  deleteConsulta(id: number): void {
    this.confirmationService.confirm({
      message: '¿Está seguro de que desea eliminar esta consulta?',
      header: 'Confirmar eliminación',
      icon: 'pi pi-exclamation-triangle',
      accept: () => {
        this.consultaService.deleteConsulta(id).subscribe({
          next: () => {
            this.consultas = this.consultas.filter(c => c.id_consulta !== id);
            this.messageService.add({
              severity: 'success',
              summary: 'Éxito',
              detail: 'Consulta eliminada exitosamente'
            });
          },
          error: (error) => {
            console.error('Error deleting consulta:', error);
            this.messageService.add({
              severity: 'error',
              summary: 'Error',
              detail: 'Error al eliminar la consulta'
            });
          }
        });
      }
    });
  }

  getTipoSeverity(tipo: string): string {
    switch (tipo) {
      case 'urgencia': return 'danger';
      case 'especialidad': return 'info';
      case 'control': return 'warning';
      case 'consulta_general': return 'success';
      default: return 'secondary';
    }
  }

  getMedicoNombre(medicoId: number): string {
    const medico = this.medicos.find(m => m.id_usuario === medicoId);
    return medico ? medico.nombre : 'N/A';
  }

  getPacienteNombre(pacienteId: number): string {
    const paciente = this.pacientes.find(p => p.id_usuario === pacienteId);
    return paciente ? paciente.nombre : 'N/A';
  }

  getConsultorioNombre(consultorioId: number): string {
    const consultorio = this.consultorios.find(c => c.id_consultorio === consultorioId);
    return consultorio ? consultorio.nombre : 'N/A';
  }
}