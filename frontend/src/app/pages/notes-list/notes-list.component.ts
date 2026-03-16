import { Component, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { ActivatedRoute, RouterLink } from '@angular/router';
import { CourseShareApiService } from '../../services/course-share-api.service';
import { Note } from '../../mock/mock-data.service';
import { MatCardModule } from '@angular/material/card';
import { MatButtonModule } from '@angular/material/button';
import { MatProgressSpinnerModule } from '@angular/material/progress-spinner';
import { MatIconModule } from '@angular/material/icon';

@Component({
  selector: 'app-notes-list',
  standalone: true,
  imports: [CommonModule, MatCardModule, MatButtonModule, MatProgressSpinnerModule, RouterLink, MatIconModule],
  templateUrl: './notes-list.component.html',
  styleUrl: './notes-list.component.css'
})
export class NotesListComponent implements OnInit {
  notes: Note[] = [];
  courseId: string | null = null;
  loading = true;
  error: string | null = null;

  constructor(
    private route: ActivatedRoute,
    private apiService: CourseShareApiService
  ) { }

  ngOnInit(): void {
    this.route.paramMap.subscribe(params => {
      this.courseId = params.get('courseId');
      if (this.courseId) {
        this.fetchNotes(this.courseId);
      }
    });
  }

  fetchNotes(courseId: string): void {
    this.loading = true;
    this.error = null;
    this.apiService.getNotesByCourse(courseId).subscribe({
      next: (data) => {
        this.notes = data.sort((a, b) => new Date(b.created_at).getTime() - new Date(a.created_at).getTime());
        this.loading = false;
      },
      error: (err) => {
        this.error = err.message || 'Failed to load notes. Please try again.';
        this.loading = false;
      }
    });
  }
}
