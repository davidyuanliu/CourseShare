import { Component, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { CourseShareApiService } from '../../services/course-share-api.service';
import { Course } from '../../mock/mock-data.service';
import { MatCardModule } from '@angular/material/card';
import { MatButtonModule } from '@angular/material/button';
import { MatProgressSpinnerModule } from '@angular/material/progress-spinner';
import { ActivatedRoute, Router, RouterLink } from '@angular/router';
import { FormsModule } from '@angular/forms';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatInputModule } from '@angular/material/input';
import { MatIconModule } from '@angular/material/icon';
import { DemoControlService } from '../../services/demo-control.service';

@Component({
  selector: 'app-courses',
  standalone: true,
  imports: [CommonModule, MatCardModule, MatButtonModule, MatProgressSpinnerModule, RouterLink, FormsModule, MatFormFieldModule, MatInputModule, MatIconModule],
  templateUrl: './courses.component.html',
  styleUrl: './courses.component.css'
})
export class CoursesComponent implements OnInit {
  courses: Course[] = [];
  loading = true;
  error: string | null = null;
  newCourseName = '';
  addingCourse = false;

  constructor(
    private apiService: CourseShareApiService,
    public demoControl: DemoControlService,
    private router: Router
  ) { }

  ngOnInit(): void {
    this.fetchCourses();
  }

  fetchCourses(): void {
    this.loading = true;
    this.error = null;
    this.apiService.getCourses().subscribe({
      next: (data) => {
        this.courses = data;
        this.loading = false;
      },
      error: (err) => {
        this.error = err.message || 'Failed to load courses. Please try again.';
        this.loading = false;
      }
    });
  }

  addCourse(): void {
    if (!this.newCourseName.trim()) return;

    this.addingCourse = true;
    this.apiService.createCourse({ name: this.newCourseName }).subscribe({
      next: (course) => {
        console.log('Course added successfully:', course);
        this.courses.push(course);
        this.newCourseName = '';
        this.addingCourse = false;
      },
      error: (err) => {
        console.error('Failed to add course:', err);
        this.addingCourse = false;
      }
    });
  }

  viewNotes(courseId: number): void {
    console.log('Navigating to notes for course ID:', courseId);
    if (courseId) {
      this.router.navigate(['/courses', courseId, 'notes']);
    } else {
      console.error('Cannot navigate: courseId is undefined or 0');
    }
  }
}
