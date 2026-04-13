import { Component, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { ReactiveFormsModule, FormBuilder, FormGroup, Validators } from '@angular/forms';
import { ActivatedRoute, Router, RouterLink } from '@angular/router';
import { CourseShareApiService } from '../../services/course-share-api.service';
import { MatCardModule } from '@angular/material/card';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatInputModule } from '@angular/material/input';
import { MatButtonModule } from '@angular/material/button';
import { MatIconModule } from '@angular/material/icon';
import { MatProgressSpinnerModule } from '@angular/material/progress-spinner';
import { MatSnackBar, MatSnackBarModule } from '@angular/material/snack-bar';
import { HttpErrorResponse } from '@angular/common/http';

@Component({
  selector: 'app-edit-note',
  standalone: true,
  imports: [
    CommonModule,
    ReactiveFormsModule,
    RouterLink,
    MatCardModule,
    MatFormFieldModule,
    MatInputModule,
    MatButtonModule,
    MatIconModule,
    MatProgressSpinnerModule,
    MatSnackBarModule
  ],
  template: `
    <div class="header-row">
        <button mat-icon-button [routerLink]="['/notes', noteId]" aria-label="Back to note">
            <mat-icon>arrow_back</mat-icon>
        </button>
        <h2>Edit Note</h2>
    </div>

    <div *ngIf="loading" class="loading-container">
        <mat-spinner diameter="50"></mat-spinner>
    </div>

    <mat-card class="form-card" *ngIf="!loading">
        <mat-card-content>
            <div *ngIf="serverError" class="error-message server-error">
                <p>{{ serverError }}</p>
            </div>

            <form [formGroup]="noteForm" (ngSubmit)="onSubmit()">
                <mat-form-field appearance="outline" class="full-width">
                    <mat-label>Title</mat-label>
                    <input matInput formControlName="title" placeholder="Enter note title">
                    <mat-error *ngIf="noteForm.get('title')?.hasError('required')">
                        Title is required
                    </mat-error>
                </mat-form-field>

                <mat-form-field appearance="outline" class="full-width">
                    <mat-label>Content</mat-label>
                    <textarea matInput formControlName="content" placeholder="Enter note content"
                        rows="8"></textarea>
                    <mat-error *ngIf="noteForm.get('content')?.hasError('required')">
                        Content is required
                    </mat-error>
                </mat-form-field>

                <div class="form-actions" style="display:flex; justify-content:flex-end; gap:10px; margin-top:20px;">
                    <button mat-stroked-button type="button" [routerLink]="['/notes', noteId]">
                        Cancel
                    </button>
                    <button mat-raised-button color="primary" type="submit" [disabled]="submitting">
                        <mat-spinner *ngIf="submitting" diameter="20" class="btn-spinner" style="display:inline-block"></mat-spinner>
                        {{ submitting ? 'Submitting...' : 'Save Changes' }}
                    </button>
                </div>
            </form>
        </mat-card-content>
    </mat-card>
  `,
  styleUrl: '../create-note/create-note.component.css'
})
export class EditNoteComponent implements OnInit {
  noteForm!: FormGroup;
  noteId: string = '';
  submitting = false;
  loading = true;
  serverError: string | null = null;
  courseId: number | null = null;

  constructor(
    private fb: FormBuilder,
    private route: ActivatedRoute,
    private router: Router,
    private apiService: CourseShareApiService,
    private snackBar: MatSnackBar
  ) { }

  ngOnInit(): void {
    this.noteId = this.route.snapshot.paramMap.get('noteId') || '';
    
    this.noteForm = this.fb.group({
      title: ['', [Validators.required]],
      content: ['', [Validators.required]]
    });

    if (this.noteId) {
      this.fetchNote();
    }
  }

  fetchNote(): void {
    this.apiService.getNote(this.noteId).subscribe({
      next: (note) => {
        this.noteForm.patchValue({
          title: note.title,
          content: note.content
        });
        this.courseId = note.courseId;
        this.loading = false;
      },
      error: (err) => {
        this.serverError = 'Failed to load note';
        this.loading = false;
      }
    });
  }

  onSubmit(): void {
    if (this.noteForm.invalid) {
      this.noteForm.markAllAsTouched();
      return;
    }

    this.submitting = true;
    this.serverError = null;

    this.apiService.updateNote(this.noteId, this.noteForm.value).subscribe({
      next: () => {
        this.submitting = false;
        this.snackBar.open('Note updated successfully!', 'Close', {
          duration: 3000
        });
        this.router.navigate(['/notes', this.noteId]);
      },
      error: (err: HttpErrorResponse) => {
        this.submitting = false;
        if (err.error && typeof err.error === 'object' && err.error.error) {
          this.serverError = err.error.error;
        } else if (err.error && typeof err.error === 'string') {
          this.serverError = err.error;
        } else {
          this.serverError = err.message || 'Failed to update note. Please try again.';
        }
      }
    });
  }
}
