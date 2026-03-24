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
  selector: 'app-create-note',
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
  templateUrl: './create-note.component.html',
  styleUrl: './create-note.component.css'
})
export class CreateNoteComponent implements OnInit {
  noteForm!: FormGroup;
  courseId: string = '';
  submitting = false;
  serverError: string | null = null;

  constructor(
    private fb: FormBuilder,
    private route: ActivatedRoute,
    private router: Router,
    private apiService: CourseShareApiService,
    private snackBar: MatSnackBar
  ) { }

  ngOnInit(): void {
    this.courseId = this.route.snapshot.paramMap.get('courseId') || '';

    this.noteForm = this.fb.group({
      title: ['', [Validators.required]],
      content: ['', [Validators.required]],
      author_name: ['']
    });
  }

  onSubmit(): void {
    if (this.noteForm.invalid) {
      this.noteForm.markAllAsTouched();
      return;
    }

    this.submitting = true;
    this.serverError = null;

    const notePayload = {
      id: Date.now().toString(),
      ...this.noteForm.value,
      course_id: this.courseId,
      created_at: new Date().toISOString()
    };

    this.apiService.createNote(notePayload).subscribe({
      next: () => {
        this.submitting = false;
        this.snackBar.open('Note created successfully!', 'Close', {
          duration: 3000,
          horizontalPosition: 'center',
          verticalPosition: 'bottom'
        });
        this.router.navigate(['/courses', this.courseId, 'notes']);
      },
      error: (err: HttpErrorResponse) => {
        this.submitting = false;
        if (err.error && typeof err.error === 'object' && err.error.error) {
          this.serverError = err.error.error;
        } else if (err.error && typeof err.error === 'string') {
          this.serverError = err.error;
        } else {
          this.serverError = err.message || 'Failed to create note. Please try again.';
        }
      }
    });
  }
}
