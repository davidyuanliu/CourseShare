import { Injectable } from '@angular/core';
import { HttpClient, HttpHeaders } from '@angular/common/http';
import { Observable, throwError } from 'rxjs';
import { Course, Note } from '../mock/mock-data.service';
import { DemoControlService } from './demo-control.service';

@Injectable({
  providedIn: 'root'
})
export class CourseShareApiService {
  private apiUrl = 'https://courseshare.onrender.com';
  private headers = new HttpHeaders();

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

  updateNote(noteId: number | string, note: Partial<Note>): Observable<Note> {
    const error = this.checkError('update note');
    if (error) return error;
    return this.http.put<Note>(`${this.apiUrl}/notes/${noteId}`, note, { headers: this.headers });
  }

  deleteNote(noteId: number | string): Observable<any> {
    const error = this.checkError('delete note');
    if (error) return error;
    return this.http.delete(`${this.apiUrl}/notes/${noteId}`, { headers: this.headers });
  }

  getMyNotes(): Observable<Note[]> {
    const error = this.checkError('fetch my notes');
    if (error) return error;
    return this.http.get<Note[]>(`${this.apiUrl}/my-notes`, { headers: this.headers });
  }

  getSavedNotes(): Observable<Note[]> {
    const error = this.checkError('fetch saved notes');
    if (error) return error;
    return this.http.get<Note[]>(`${this.apiUrl}/saved-notes`, { headers: this.headers });
  }

  markNoteHelpful(noteId: number | string): Observable<any> {
    const error = this.checkError('mark note helpful');
    if (error) return error;
    return this.http.post(`${this.apiUrl}/notes/${noteId}/helpful`, {}, { headers: this.headers });
  }

  unmarkNoteHelpful(noteId: number | string): Observable<any> {
    const error = this.checkError('unmark note helpful');
    if (error) return error;
    return this.http.delete(`${this.apiUrl}/notes/${noteId}/helpful`, { headers: this.headers });
  }

  saveNote(noteId: number | string): Observable<any> {
    const error = this.checkError('save note');
    if (error) return error;
    return this.http.post(`${this.apiUrl}/notes/${noteId}/save`, {}, { headers: this.headers });
  }

  unsaveNote(noteId: number | string): Observable<any> {
    const error = this.checkError('unsave note');
    if (error) return error;
    return this.http.delete(`${this.apiUrl}/notes/${noteId}/save`, { headers: this.headers });
  }
}
