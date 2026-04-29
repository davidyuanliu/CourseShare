import { Routes } from '@angular/router';
import { CoursesComponent } from './pages/courses/courses.component';
import { NotesListComponent } from './pages/notes-list/notes-list.component';
import { NoteDetailComponent } from './pages/note-detail/note-detail.component';
import { CreateNoteComponent } from './pages/create-note/create-note.component';
import { RegisterComponent } from './pages/register/register.component';
import { LoginComponent } from './pages/login/login.component';
import { EditNoteComponent } from './pages/edit-note/edit-note.component';
import { authGuard } from './services/auth.guard';
import { MyNotesComponent } from './pages/my-notes/my-notes.component';
import { SavedNotesComponent } from './pages/saved-notes/saved-notes.component';

export const routes: Routes = [
    { path: '', redirectTo: '/courses', pathMatch: 'full' },
    { path: 'courses', component: CoursesComponent },
    { path: 'login', component: LoginComponent },
    { path: 'register', component: RegisterComponent },
    { path: 'courses/:courseId/notes/new', component: CreateNoteComponent, canActivate: [authGuard] },
    { path: 'courses/:courseId/notes', component: NotesListComponent },
    { path: 'notes/:noteId', component: NoteDetailComponent },
    { path: 'notes/:noteId/edit', component: EditNoteComponent, canActivate: [authGuard] },
    { path: 'my-notes', component: MyNotesComponent, canActivate: [authGuard] },
    { path: 'saved-notes', component: SavedNotesComponent, canActivate: [authGuard] },
    { path: '**', redirectTo: '/courses' }
];
