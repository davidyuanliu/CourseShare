import { Injectable } from '@angular/core';
import { HttpClient, HttpHeaders } from '@angular/common/http';
import { Observable, throwError } from 'rxjs';
import { Course, Note } from '../mock/mock-data.service';
import { DemoControlService } from './demo-control.service';

@Injectable({
  providedIn: 'root'
})
export class CourseShareApiService {
  private apiUrl = 'https://unconstellated-ruthann-fiducially.ngrok-free.dev';
  private headers = new HttpHeaders().set('ngrok-skip-browser-warning', 'true');

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
    return this.http.get<Course[]>(`${this.apiUrl}/courses`, { headers: this.headers });
  }

  createCourse(course: { name: string }): Observable<Course> {
    const error = this.checkError('create course');
    if (error) return error;
    return this.http.post<Course>(`${this.apiUrl}/courses`, course, { headers: this.headers });
  }

  getNotesByCourse(courseId: string): Observable<Note[]> {
    const error = this.checkError('fetch notes');
    if (error) return error;
    return this.http.get<Note[]>(`${this.apiUrl}/courses/${courseId}/notes`, { headers: this.headers });
  }

  getNote(noteId: string): Observable<Note> {
    const error = this.checkError('fetch note detail');
    if (error) return error;
    return this.http.get<Note>(`${this.apiUrl}/notes/${noteId}`, { headers: this.headers });
  }

  createNote(note: Partial<Note>): Observable<Note> {
    const error = this.checkError('create note');
    if (error) return error;
    return this.http.post<Note>(`${this.apiUrl}/notes`, note, { headers: this.headers });
  }
}
