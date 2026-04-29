import { ComponentFixture, TestBed } from '@angular/core/testing';
import { NotesListComponent } from './notes-list.component';
import { HttpClientTestingModule } from '@angular/common/http/testing';
import { ActivatedRoute } from '@angular/router';
import { RouterTestingModule } from '@angular/router/testing';
import { NoopAnimationsModule } from '@angular/platform-browser/animations';
import { AuthService } from '../../services/auth.service';
import { BehaviorSubject, of } from 'rxjs';
import { User } from '../../models/auth';
import { Note } from '../../mock/mock-data.service';
import { CourseShareApiService } from '../../services/course-share-api.service';

describe('NotesListComponent', () => {
  let component: NotesListComponent;
  let fixture: ComponentFixture<NotesListComponent>;
  let authServiceSpy: jasmine.SpyObj<AuthService>;
  let apiServiceSpy: jasmine.SpyObj<CourseShareApiService>;
  let currentUserSubject: BehaviorSubject<User | null>;

  const mockNotes: Note[] = [
    { 
        ID: 1, 
        title: 'Angular Basics', 
        content: 'Learn components', 
        courseId: 1, 
        userId: 1, 
        tags: ['angular', 'basics'],
        helpfulCount: 5,
        CreatedAt: '2023-01-01T00:00:00Z',
        isSaved: true
    },
    { 
        ID: 2, 
        title: 'React Introduction', 
        content: 'Hooks are cool', 
        courseId: 1, 
        userId: 2, 
        tags: ['react', 'web'],
        helpfulCount: 10,
        CreatedAt: '2023-02-01T00:00:00Z',
        isSaved: false
    },
    { 
        ID: 3, 
        title: 'Advanced Angular', 
        content: 'RxJS and more', 
        courseId: 1, 
        userId: 1, 
        tags: ['angular', 'rxjs'],
        helpfulCount: 2,
        CreatedAt: '2023-03-01T00:00:00Z',
        isSaved: false
    }
  ];

  beforeEach(async () => {
    currentUserSubject = new BehaviorSubject<User | null>(null);
    authServiceSpy = jasmine.createSpyObj('AuthService', ['logout']);
    Object.defineProperty(authServiceSpy, 'currentUserValue', {
      get: () => currentUserSubject.value
    });

    apiServiceSpy = jasmine.createSpyObj('CourseShareApiService', ['getNotesByCourse']);
    apiServiceSpy.getNotesByCourse.and.returnValue(of(mockNotes));

    const mockActivatedRoute = {
      paramMap: of({ get: (key: string) => key === 'courseId' ? '1' : null })
    };

    await TestBed.configureTestingModule({
      imports: [NotesListComponent, HttpClientTestingModule, RouterTestingModule, NoopAnimationsModule],
      providers: [
        { provide: AuthService, useValue: authServiceSpy },
        { provide: CourseShareApiService, useValue: apiServiceSpy },
        { provide: ActivatedRoute, useValue: mockActivatedRoute }
      ]
    })
    .compileComponents();
    
    fixture = TestBed.createComponent(NotesListComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });

  it('should filter notes by search query', () => {
    component.searchQuery = 'React';
    const filtered = component.filteredNotes;
    expect(filtered.length).toBe(1);
    expect(filtered[0].title).toBe('React Introduction');
  });

  it('should filter notes by multiple tags', () => {
    component.selectedTags = ['angular'];
    const filtered = component.filteredNotes;
    expect(filtered.length).toBe(2);
    expect(filtered.every(n => n.tags?.includes('angular'))).toBeTrue();
  });

  it('should filter by "Written by me"', () => {
    currentUserSubject.next({ id: 1, email: 'me@test.com' });
    component.showOnlyMine = true;
    const filtered = component.filteredNotes;
    expect(filtered.length).toBe(2);
    expect(filtered.every(n => n.userId === 1)).toBeTrue();
  });

  it('should filter by "Saved by me"', () => {
    component.showOnlySaved = true;
    const filtered = component.filteredNotes;
    expect(filtered.length).toBe(1);
    expect(filtered[0].ID).toBe(1);
  });

  it('should sort notes by likes', () => {
    component.sortBy = 'likes';
    const filtered = component.filteredNotes;
    expect(filtered[0].helpfulCount).toBe(10);
    expect(filtered[1].helpfulCount).toBe(5);
    expect(filtered[2].helpfulCount).toBe(2);
  });

  it('should sort notes by title', () => {
    component.sortBy = 'title';
    const filtered = component.filteredNotes;
    expect(filtered[0].title).toBe('Advanced Angular');
    expect(filtered[1].title).toBe('Angular Basics');
    expect(filtered[2].title).toBe('React Introduction');
  });

  it('should toggle all tags', () => {
    component.toggleAllTags(true);
    expect(component.selectedTags.length).toBe(component.allTags.length);
    component.toggleAllTags(false);
    expect(component.selectedTags.length).toBe(0);
  });
});
