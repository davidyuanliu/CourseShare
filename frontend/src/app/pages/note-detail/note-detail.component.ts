import { Component, OnInit } from '@angular/core';
import { CommonModule, Location } from '@angular/common';
import { ActivatedRoute, RouterLink } from '@angular/router';
import { CourseShareApiService } from '../../services/course-share-api.service';
import { Note } from '../../mock/mock-data.service';
import { MatCardModule } from '@angular/material/card';
import { MatButtonModule } from '@angular/material/button';
import { MatProgressSpinnerModule } from '@angular/material/progress-spinner';
import { MatIconModule } from '@angular/material/icon';
import { MatDividerModule } from '@angular/material/divider';

@Component({
  selector: 'app-note-detail',
  standalone: true,
  imports: [CommonModule, MatCardModule, MatButtonModule, MatProgressSpinnerModule, RouterLink, MatIconModule, MatDividerModule],
  templateUrl: './note-detail.component.html',
  styleUrl: './note-detail.component.css'
})
export class NoteDetailComponent implements OnInit {
  note: Note | null = null;
  loading = true;
  error: string | null = null;

  constructor(
    private route: ActivatedRoute,
    private apiService: CourseShareApiService,
    private location: Location
  ) { }

  ngOnInit(): void {
    const noteId = this.route.snapshot.paramMap.get('noteId');
    if (noteId) {
      this.fetchNote(noteId);
    }
  }

  fetchNote(noteId: string): void {
    this.loading = true;
    this.error = null;
    this.apiService.getNote(noteId).subscribe({
      next: (data) => {
        this.note = data;
        this.loading = false;
      },
      error: (err) => {
        this.error = err.message || 'Note not found or failed to load.';
        this.loading = false;
      }
    });
  }

  goBack(): void {
    this.location.back();
  }
}
