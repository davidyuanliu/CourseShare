import { Component, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { ActivatedRoute, RouterLink } from '@angular/router';
import { CourseShareApiService } from '../../services/course-share-api.service';
import { Note } from '../../mock/mock-data.service';
import { MatCardModule } from '@angular/material/card';
import { MatButtonModule } from '@angular/material/button';
import { MatProgressSpinnerModule } from '@angular/material/progress-spinner';
import { MatIconModule } from '@angular/material/icon';
import { FormsModule } from '@angular/forms';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatInputModule } from '@angular/material/input';
import { MatSelectModule } from '@angular/material/select';
import { MatCheckboxModule } from '@angular/material/checkbox';
import { AuthService } from '../../services/auth.service';
import { MatDividerModule } from '@angular/material/divider';

@Component({
  selector: 'app-notes-list',
  standalone: true,
  imports: [CommonModule, MatCardModule, MatButtonModule, MatProgressSpinnerModule, RouterLink, MatIconModule, FormsModule, MatFormFieldModule, MatInputModule, MatSelectModule, MatCheckboxModule, MatDividerModule],
  templateUrl: './notes-list.component.html',
  styleUrl: './notes-list.component.css'
})
export class NotesListComponent implements OnInit {
  notes: Note[] = [];
  courseId: string | null = null;
  loading = true;
  error: string | null = null;

  searchQuery = '';
  sortBy = 'newest';
  selectedTags: string[] = [];
  showOnlyMine = false;
  showOnlySaved = false;

  constructor(
    private route: ActivatedRoute,
    private apiService: CourseShareApiService,
    private authService: AuthService
  ) { }

  ngOnInit(): void {
    this.route.paramMap.subscribe(params => {
      this.courseId = params.get('courseId');
      if (this.courseId) {
        this.fetchNotes(this.courseId);
      }
    });
  }

  get allTags(): string[] {
    const tags = new Set<string>();
    this.notes.forEach(note => {
      if (note.tags) {
        note.tags.forEach(tag => tags.add(tag));
      }
    });
    return Array.from(tags).sort();
  }

  get filteredNotes(): Note[] {
    let result = this.notes;

    if (this.searchQuery) {
      const q = this.searchQuery.toLowerCase();
      result = result.filter(n => n.title.toLowerCase().includes(q) || n.content.toLowerCase().includes(q));
    }

    if (this.selectedTags.length > 0) {
      result = result.filter(n => n.tags && n.tags.some(tag => this.selectedTags.includes(tag)));
    }

    if (this.showOnlyMine) {
      const currentUserId = this.authService.currentUserValue?.id;
      if (currentUserId) {
        result = result.filter(n => n.userId === currentUserId);
      }
    }

    if (this.showOnlySaved) {
      result = result.filter(n => n.isSaved);
    }

    result = result.slice().sort((a, b) => {
      if (this.sortBy === 'newest') {
        return new Date(b.CreatedAt).getTime() - new Date(a.CreatedAt).getTime();
      } else if (this.sortBy === 'oldest') {
        return new Date(a.CreatedAt).getTime() - new Date(b.CreatedAt).getTime();
      } else if (this.sortBy === 'title') {
        return a.title.localeCompare(b.title);
      } else if (this.sortBy === 'likes') {
        return (b.helpfulCount || 0) - (a.helpfulCount || 0);
      }
      return 0;
    });

    return result;
  }

  fetchNotes(courseId: string): void {
    this.loading = true;
    this.error = null;
    this.apiService.getNotesByCourse(courseId).subscribe({
      next: (data) => {
        this.notes = data;
        this.loading = false;
      },
      error: (err) => {
        this.error = err.message || 'Failed to load notes. Please try again.';
        this.loading = false;
      }
    });
  }

  get allSelected(): boolean {
    return this.allTags.length > 0 && this.selectedTags.length === this.allTags.length;
  }

  get partiallySelected(): boolean {
    return this.selectedTags.length > 0 && this.selectedTags.length < this.allTags.length;
  }

  toggleAllTags(checked: boolean): void {
    if (checked) {
      this.selectedTags = [...this.allTags];
    } else {
      this.selectedTags = [];
    }
  }

  onTagChange(tag: string, checked: boolean): void {
    if (checked) {
      if (!this.selectedTags.includes(tag)) {
        this.selectedTags = [...this.selectedTags, tag];
      }
    } else {
      this.selectedTags = this.selectedTags.filter(t => t !== tag);
    }
  }

  get isLoggedIn(): boolean {
    return !!this.authService.currentUserValue;
  }
}
