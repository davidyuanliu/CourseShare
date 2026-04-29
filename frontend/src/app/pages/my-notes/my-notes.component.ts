import { Component, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { RouterLink } from '@angular/router';
import { CourseShareApiService } from '../../services/course-share-api.service';
import { Note } from '../../mock/mock-data.service';
import { MatCardModule } from '@angular/material/card';
import { MatButtonModule } from '@angular/material/button';
import { MatProgressSpinnerModule } from '@angular/material/progress-spinner';
import { MatIconModule } from '@angular/material/icon';
import { MatChipsModule } from '@angular/material/chips';

@Component({
  selector: 'app-my-notes',
  standalone: true,
  imports: [CommonModule, MatCardModule, MatButtonModule, MatProgressSpinnerModule, RouterLink, MatIconModule, MatChipsModule],
  templateUrl: './my-notes.component.html',
  styleUrl: '../notes-list/notes-list.component.css'
})
export class MyNotesComponent implements OnInit {
  notes: Note[] = [];
  loading = true;
  error: string | null = null;

  constructor(private apiService: CourseShareApiService) {}

  ngOnInit(): void {
    this.fetchMyNotes();
  }

  fetchMyNotes(): void {
    this.loading = true;
    this.error = null;
    this.apiService.getMyNotes().subscribe({
      next: (data) => {
        this.notes = data.sort((a, b) => new Date(b.CreatedAt).getTime() - new Date(a.CreatedAt).getTime());
        this.loading = false;
      },
      error: (err) => {
        this.error = err.error?.error || 'Failed to load your notes.';
        this.loading = false;
      }
    });
  }
}
