import { ComponentFixture, TestBed, fakeAsync, tick } from '@angular/core/testing';
import { ReactiveFormsModule } from '@angular/forms';
import { RouterTestingModule } from '@angular/router/testing';
import { NoopAnimationsModule } from '@angular/platform-browser/animations';
import { Router } from '@angular/router';
import { of, throwError } from 'rxjs';
import { HttpErrorResponse } from '@angular/common/http';

import { CreateNoteComponent } from './create-note.component';
import { CourseShareApiService } from '../../services/course-share-api.service';
import { MatSnackBar } from '@angular/material/snack-bar';

describe('CreateNoteComponent', () => {
  let component: CreateNoteComponent;
  let fixture: ComponentFixture<CreateNoteComponent>;
  let mockApiService: jasmine.SpyObj<CourseShareApiService>;
  let mockSnackBar: jasmine.SpyObj<MatSnackBar>;
  let router: Router;

  beforeEach(async () => {
    mockApiService = jasmine.createSpyObj('CourseShareApiService', ['createNote']);
    mockSnackBar = jasmine.createSpyObj('MatSnackBar', ['open']);

    await TestBed.configureTestingModule({
      imports: [
        CreateNoteComponent,
        ReactiveFormsModule,
        RouterTestingModule.withRoutes([]),
        NoopAnimationsModule
      ],
      providers: [
        { provide: CourseShareApiService, useValue: mockApiService },
        { provide: MatSnackBar, useValue: mockSnackBar }
      ]
    }).compileComponents();

    router = TestBed.inject(Router);
    spyOn(router, 'navigate');

    fixture = TestBed.createComponent(CreateNoteComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });

  it('should initialize form with empty fields', () => {
    expect(component.noteForm.get('title')?.value).toBe('');
    expect(component.noteForm.get('content')?.value).toBe('');
    expect(component.noteForm.get('author_name')?.value).toBe('');
  });

  it('should mark title as invalid when empty', () => {
    const titleControl = component.noteForm.get('title')!;
    titleControl.markAsTouched();
    expect(titleControl.hasError('required')).toBeTrue();
  });

  it('should mark content as invalid when empty', () => {
    const contentControl = component.noteForm.get('content')!;
    contentControl.markAsTouched();
    expect(contentControl.hasError('required')).toBeTrue();
  });

  it('should have invalid form when required fields are empty', () => {
    expect(component.noteForm.valid).toBeFalse();
  });

  it('should have valid form when title and content are filled', () => {
    component.noteForm.patchValue({ title: 'Test', content: 'Test content' });
    expect(component.noteForm.valid).toBeTrue();
  });

  it('should not call API when form is invalid', () => {
    component.onSubmit();
    expect(mockApiService.createNote).not.toHaveBeenCalled();
  });

  it('should call API and navigate on successful submit', fakeAsync(() => {
    const mockNote = { id: '999', title: 'Test', content: 'Content', course_id: '1', author_name: '', created_at: '' };
    mockApiService.createNote.and.returnValue(of(mockNote));
    component.courseId = '1';

    component.noteForm.patchValue({ title: 'Test', content: 'Content' });
    component.onSubmit();
    tick();

    expect(mockApiService.createNote).toHaveBeenCalled();
    expect(mockSnackBar.open).toHaveBeenCalledWith('Note created successfully!', 'Close', jasmine.any(Object));
    expect(router.navigate).toHaveBeenCalledWith(['/courses', '1', 'notes']);
    expect(component.submitting).toBeFalse();
  }));

  it('should display server error on API failure', fakeAsync(() => {
    const errorResponse = new HttpErrorResponse({
      error: { error: 'Title is required' },
      status: 400
    });
    mockApiService.createNote.and.returnValue(throwError(() => errorResponse));
    component.courseId = '1';

    component.noteForm.patchValue({ title: 'Test', content: 'Content' });
    component.onSubmit();
    tick();

    expect(component.serverError).toBe('Title is required');
    expect(router.navigate).not.toHaveBeenCalled();
    expect(component.submitting).toBeFalse();
  }));

  it('should disable submit button while submitting', () => {
    component.submitting = true;
    fixture.detectChanges();
    const submitBtn = fixture.nativeElement.querySelector('#note-submit-btn');
    expect(submitBtn.disabled).toBeTrue();
  });
});
