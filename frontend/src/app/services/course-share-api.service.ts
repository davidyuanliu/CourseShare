import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable, throwError } from 'rxjs';
import { Course, Note } from '../mock/mock-data.service';
import { DemoControlService } from './demo-control.service';

@Injectable({
  providedIn: 'root'
})
export class CourseShareApiService {
  private apiUrl = 'api';

  constructor(private http: HttpClient, private demoControl: DemoControlService) { }

  private checkError(operation: string): Observable<never> | null {
    if (this.demoControl.errorMode) {
      return throwError(() => new Error(`Demo Error: Failed to ${operation}. Please disable Error Mode to try again.`));
    }
    return null;
  }

  getCourses(): Observable<Course[]> {
    const error = this.checkError('fetch courses');
    if (error) return error;
    return this.http.get<Course[]>(`${this.apiUrl}/courses`);
  }

  getNotesByCourse(courseId: string): Observable<Note[]> {
    const error = this.checkError('fetch notes');
    if (error) return error;
    return this.http.get<Note[]>(`${this.apiUrl}/notes?course_id=${courseId}`);
  }

  getNote(noteId: string): Observable<Note> {
    const error = this.checkError('fetch note detail');
    if (error) return error;
    return this.http.get<Note>(`${this.apiUrl}/notes/${noteId}`);
  }
}
