import { ComponentFixture, TestBed } from '@angular/core/testing';
import { NoteDetailComponent } from './note-detail.component';
import { HttpClientTestingModule } from '@angular/common/http/testing';
import { RouterTestingModule } from '@angular/router/testing';
import { AuthService } from '../../services/auth.service';
import { BehaviorSubject } from 'rxjs';
import { User } from '../../models/auth';

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
      imports: [NoteDetailComponent, HttpClientTestingModule, RouterTestingModule],
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
    component.note = { id: 1, courseId: 1, title: 'T', content: 'C', CreatedAt: '', userId: 1, author: {id: 1, email: 'A'} };
    currentUserSubject.next(null);
    component.isOwner = !!component.note && !!authServiceSpy.currentUserValue && authServiceSpy.currentUserValue.id === component.note.userId;
    expect(component.isOwner).toBeFalsy();
  });

  it('should be owner if logged in as author', () => {
    component.note = { id: 1, courseId: 1, title: 'T', content: 'C', CreatedAt: '', userId: 1, author: {id: 1, email: 'A'} };
    currentUserSubject.next({ id: 1, email: 'A' });
    component.isOwner = !!component.note && !!authServiceSpy.currentUserValue && authServiceSpy.currentUserValue.id === component.note.userId;
    expect(component.isOwner).toBeTruthy();
  });
});
