import { Routes } from '@angular/router';
import { CoursesComponent } from './pages/courses/courses.component';
import { NotesListComponent } from './pages/notes-list/notes-list.component';
import { NoteDetailComponent } from './pages/note-detail/note-detail.component';

export const routes: Routes = [
    { path: '', redirectTo: '/courses', pathMatch: 'full' },
    { path: 'courses', component: CoursesComponent },
    { path: 'courses/:courseId/notes', component: NotesListComponent },
    { path: 'notes/:noteId', component: NoteDetailComponent },
    { path: '**', redirectTo: '/courses' }
];
