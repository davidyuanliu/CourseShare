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

@Component({
  selector: 'app-note-detail',
  standalone: true,
  imports: [CommonModule, MatCardModule, MatButtonModule, MatProgressSpinnerModule, RouterLink, MatIconModule, MatDividerModule, MatDialogModule, MatSnackBarModule],
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
}
