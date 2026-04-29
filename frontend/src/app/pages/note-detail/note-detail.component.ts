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
import { MatDialog, MatDialogModule } from '@angular/material/dialog';
import { MatSnackBar, MatSnackBarModule } from '@angular/material/snack-bar';
import { AuthService } from '../../services/auth.service';
import { Router } from '@angular/router';
import { ConfirmDialogComponent } from './confirm-dialog.component';
import { MatChipsModule } from '@angular/material/chips';
import { MatTooltipModule } from '@angular/material/tooltip';

@Component({
  selector: 'app-note-detail',
  standalone: true,
  imports: [CommonModule, MatCardModule, MatButtonModule, MatProgressSpinnerModule, RouterLink, MatIconModule, MatDividerModule, MatDialogModule, MatSnackBarModule, MatChipsModule, MatTooltipModule],
  templateUrl: './note-detail.component.html',
  styleUrl: './note-detail.component.css'
})
export class NoteDetailComponent implements OnInit {
  note: Note | null = null;
  loading = true;
  error: string | null = null;
  isOwner = false;

  constructor(
    private route: ActivatedRoute,
    private apiService: CourseShareApiService,
    private authService: AuthService,
    private location: Location,
    private router: Router,
    private dialog: MatDialog,
    private snackBar: MatSnackBar
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
        const currentUser = this.authService.currentUserValue;
        this.isOwner = !!currentUser && currentUser.id === this.note.userId;
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

  deleteNote(): void {
    const dialogRef = this.dialog.open(ConfirmDialogComponent, {
      width: '400px'
    });

    dialogRef.afterClosed().subscribe(result => {
      if (result) {
        this.apiService.deleteNote(this.note!.ID).subscribe({
          next: () => {
            this.snackBar.open('Note deleted successfully', 'Close', { duration: 3000 });
            this.router.navigate(['/courses', this.note!.courseId, 'notes']);
          },
          error: (err) => {
            this.snackBar.open(err.error?.error || 'Failed to delete note', 'Close', { duration: 3000 });
          }
        });
      }
    });
  }

  toggleHelpful(): void {
    if (!this.note || !this.authService.currentUserValue) return;

    if (this.note.isHelpful) {
      this.apiService.unmarkNoteHelpful(this.note.ID).subscribe({
        next: () => {
          this.note!.isHelpful = false;
          if (this.note!.helpfulCount !== undefined) {
            this.note!.helpfulCount--;
          }
        },
        error: () => this.snackBar.open('Failed to unmark helpful', 'Close', { duration: 3000 })
      });
    } else {
      this.apiService.markNoteHelpful(this.note.ID).subscribe({
        next: () => {
          this.note!.isHelpful = true;
          this.note!.helpfulCount = (this.note!.helpfulCount || 0) + 1;
        },
        error: () => this.snackBar.open('Failed to mark helpful', 'Close', { duration: 3000 })
      });
    }
  }

  toggleSave(): void {
    if (!this.note || !this.authService.currentUserValue) return;

    if (this.note.isSaved) {
      this.apiService.unsaveNote(this.note.ID).subscribe({
        next: () => {
          this.note!.isSaved = false;
          this.snackBar.open('Note removed from saved', 'Close', { duration: 2000 });
        },
        error: () => this.snackBar.open('Failed to unsave note', 'Close', { duration: 3000 })
      });
    } else {
      this.apiService.saveNote(this.note.ID).subscribe({
        next: () => {
          this.note!.isSaved = true;
          this.snackBar.open('Note saved successfully', 'Close', { duration: 2000 });
        },
        error: () => this.snackBar.open('Failed to save note', 'Close', { duration: 3000 })
      });
    }
  }

  get isLoggedIn(): boolean {
    return !!this.authService.currentUserValue;
  }
}
