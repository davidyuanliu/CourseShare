import { Injectable } from '@angular/core';
import { InMemoryDbService } from 'angular-in-memory-web-api';

export interface Course {
  id: string;
  code: string;
  title: string;
}

export interface Note {
  id: string;
  course_id: string;
  title: string;
  author_name: string;
  created_at: string;
  content: string;
}

@Injectable({
  providedIn: 'root'
})
export class MockDataService implements InMemoryDbService {
  createDb() {
    const courses: Course[] = [
      { id: '1', code: 'CEN5035', title: 'Agile Software Engineering' },
      { id: '2', code: 'COP4600', title: 'Operating Systems' },
      { id: '3', code: 'CIS4930', title: 'Special Topics: Web 2.0' }
    ];

    const notes: Note[] = [
      { id: '101', course_id: '1', title: 'Sprint Planning Basics', author_name: 'David Liu', created_at: '2024-02-10T10:00:00Z', content: 'Sprint planning is an event in scrum that kicks off the sprint...' },
      { id: '102', course_id: '1', title: 'User Stories and Acceptance Criteria', author_name: 'Jane Doe', created_at: '2024-02-12T14:30:00Z', content: 'User stories are short, simple descriptions of a feature...' },
      { id: '201', course_id: '2', title: 'CPU Scheduling Algorithms', author_name: 'Bob Smith', created_at: '2024-01-20T09:15:00Z', content: 'CPU scheduling deals with the problem of deciding which of the processes in the ready queue is to be allocated the CPU...' },
      { id: '301', course_id: '3', title: 'Angular Standalone Components', author_name: 'Alice Brown', created_at: '2024-02-15T11:45:00Z', content: 'Standalone components provide a simplified way to build Angular applications...' }
    ];

    return { courses, notes };
  }
}
