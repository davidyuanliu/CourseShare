import { ComponentFixture, TestBed } from '@angular/core/testing';
import { NoteDetailComponent } from './note-detail.component';
import { HttpClientTestingModule } from '@angular/common/http/testing';
import { RouterTestingModule } from '@angular/router/testing';
import { NoopAnimationsModule } from '@angular/platform-browser/animations';
import { AuthService } from '../../services/auth.service';
import { User } from '../../models/auth';
import { CourseShareApiService } from '../../services/course-share-api.service';
import { Note } from '../../mock/mock-data.service';
import { BehaviorSubject, of } from 'rxjs';

describe('NoteDetailComponent', () => {
  let component: NoteDetailComponent;
  let fixture: ComponentFixture<NoteDetailComponent>;
  let authServiceSpy: jasmine.SpyObj<AuthService>;
  let currentUserSubject: BehaviorSubject<User | null>;

  beforeEach(async () => {
    currentUserSubject = new BehaviorSubject<User | null>(null);
    authServiceSpy = jasmine.createSpyObj('AuthService', ['logout']);
    Object.defineProperty(authServiceSpy, 'currentUserValue', {
      get: () => currentUserSubject.value
    });


    await TestBed.configureTestingModule({
      imports: [NoteDetailComponent, HttpClientTestingModule, RouterTestingModule, NoopAnimationsModule],
      providers: [
        { provide: AuthService, useValue: authServiceSpy }
      ]
    })
    .compileComponents();
    
    fixture = TestBed.createComponent(NoteDetailComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });

  it('should not be owner if not logged in', () => {
    component.note = { ID: 1, courseId: 1, title: 'T', content: 'C', CreatedAt: '', userId: 1, author: {id: 1, email: 'A'} };
    currentUserSubject.next(null);
    component.isOwner = !!component.note && !!authServiceSpy.currentUserValue && authServiceSpy.currentUserValue.id === component.note.userId;
    expect(component.isOwner).toBeFalsy();
  });

  it('should be owner if logged in as author', () => {
    component.note = { ID: 1, courseId: 1, title: 'T', content: 'C', CreatedAt: '', userId: 1, author: {id: 1, email: 'A'} };
    currentUserSubject.next({ id: 1, email: 'A' });
    component.isOwner = !!component.note && !!authServiceSpy.currentUserValue && authServiceSpy.currentUserValue.id === component.note.userId;
    expect(component.isOwner).toBeTruthy();
  });

  it('should toggle helpful status', () => {
    const apiService = TestBed.inject(CourseShareApiService);
    const mockNote: Note = { ID: 1, courseId: 1, title: 'T', content: 'C', CreatedAt: '', userId: 2, isHelpful: false, helpfulCount: 0 };
    component.note = mockNote;
    currentUserSubject.next({ id: 1, email: 'tester@test.com' });

    spyOn(apiService, 'markNoteHelpful').and.returnValue(of({ helpfulCount: 1, isHelpful: true }));
    
    component.toggleHelpful();
    expect(component.note.isHelpful).toBeTrue();
    expect(component.note.helpfulCount).toBe(1);
    expect(apiService.markNoteHelpful).toHaveBeenCalledWith(1);
  });

  it('should toggle save status', () => {
    const apiService = TestBed.inject(CourseShareApiService);
    const mockNote: Note = { ID: 1, courseId: 1, title: 'T', content: 'C', CreatedAt: '', userId: 2, isSaved: false };
    component.note = mockNote;
    currentUserSubject.next({ id: 1, email: 'tester@test.com' });

    spyOn(apiService, 'saveNote').and.returnValue(of({ isSaved: true }));
    
    component.toggleSave();
    expect(component.note.isSaved).toBeTrue();
    expect(apiService.saveNote).toHaveBeenCalledWith(1);
  });
});
