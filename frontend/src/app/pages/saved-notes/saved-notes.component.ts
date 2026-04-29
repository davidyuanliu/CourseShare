import { Component, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { RouterLink, RouterModule } from '@angular/router';
import { CourseShareApiService } from '../../services/course-share-api.service';
import { Note } from '../../mock/mock-data.service';
import { MatCardModule } from '@angular/material/card';
import { MatButtonModule } from '@angular/material/button';
import { MatProgressSpinnerModule } from '@angular/material/progress-spinner';
import { MatIconModule } from '@angular/material/icon';
import { MatChipsModule } from '@angular/material/chips';
import { MatSnackBar, MatSnackBarModule } from '@angular/material/snack-bar';

@Component({
  selector: 'app-saved-notes',
  standalone: true,
  imports: [CommonModule, MatCardModule, MatButtonModule, MatProgressSpinnerModule, RouterLink, RouterModule, MatIconModule, MatChipsModule, MatSnackBarModule],
  templateUrl: './saved-notes.component.html',
  styleUrl: '../notes-list/notes-list.component.css'
})
export class SavedNotesComponent implements OnInit {
  notes: Note[] = [];
  loading = true;
  error: string | null = null;

  constructor(private apiService: CourseShareApiService, private snackBar: MatSnackBar) {}

  ngOnInit(): void {
    this.fetchSavedNotes();
  }

  fetchSavedNotes(): void {
    this.loading = true;
    this.error = null;
    this.apiService.getSavedNotes().subscribe({
      next: (data) => {
        this.notes = data;
        this.loading = false;
      },
      error: (err) => {
        this.error = err.error?.error || 'Failed to load saved notes.';
        this.loading = false;
      }
    });
  }

  unsaveNote(noteId: number, event: Event): void {
    event.stopPropagation();
    this.apiService.unsaveNote(noteId).subscribe({
      next: () => {
        this.notes = this.notes.filter(n => n.id !== noteId);
        this.snackBar.open('Note removed from saved', 'Close', { duration: 2000 });
      },
      error: () => this.snackBar.open('Failed to unsave note', 'Close', { duration: 3000 })
    });
  }
}
