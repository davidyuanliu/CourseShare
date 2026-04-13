import { ComponentFixture, TestBed, fakeAsync, tick, flush } from '@angular/core/testing';
import { ReactiveFormsModule } from '@angular/forms';
import { RouterTestingModule } from '@angular/router/testing';
import { NoopAnimationsModule, BrowserAnimationsModule } from '@angular/platform-browser/animations';
import { Router, provideRouter } from '@angular/router';
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
    mockApiService.createNote.and.returnValue(of({ ID: 1, title: 'Test', content: 'Test', courseId: 1, CreatedAt: '', UpdatedAt: '' } as any));

    await TestBed.configureTestingModule({
      imports: [CreateNoteComponent, BrowserAnimationsModule],
      providers: [
        provideRouter([]),
        { provide: CourseShareApiService, useValue: mockApiService }
      ]
    })
    .overrideProvider(MatSnackBar, { useValue: mockSnackBar })
    .compileComponents();

    router = TestBed.inject(Router);
    spyOn(router, 'navigate');

    fixture = TestBed.createComponent(CreateNoteComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    (expect(component) as any).toBeTruthy();
  });

  it('should mark title as invalid when empty', () => {
    const titleControl = component.noteForm.get('title')!;
    titleControl.markAsTouched();
    (expect(titleControl.hasError('required')) as any).toBeTrue();
  });

  it('should mark content as invalid when empty', () => {
    const contentControl = component.noteForm.get('content')!;
    contentControl.markAsTouched();
    (expect(contentControl.hasError('required')) as any).toBeTrue();
  });

  it('should have invalid form when required fields are empty', () => {
    (expect(component.noteForm.valid) as any).toBeFalse();
  });

  it('should have valid form when title and content are filled', () => {
    component.noteForm.patchValue({ title: 'Test', content: 'Test content' });
    (expect(component.noteForm.valid) as any).toBeTrue();
  });

  it('should not call API when form is invalid', () => {
    component.onSubmit();
    (expect(mockApiService.createNote) as any).not.toHaveBeenCalled();
  });

  it('should call API and navigate on successful submit', fakeAsync(() => {
    const mockNote = { ID: 999, title: 'Test', content: 'Content', courseId: 1, CreatedAt: '' };
    mockApiService.createNote.and.returnValue(of(mockNote));
    component.courseId = '1';

    component.noteForm.patchValue({ title: 'Test', content: 'Content' });
    component.onSubmit();
    flush();

    (expect(mockApiService.createNote) as any).toHaveBeenCalled();
    (expect(mockSnackBar.open) as any).toHaveBeenCalledWith('Note created successfully!', 'Close', jasmine.any(Object));
    (expect(router.navigate) as any).toHaveBeenCalledWith(['/courses', '1', 'notes']);
    (expect(component.submitting) as any).toBeFalse();
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
    flush();

    (expect(component.serverError) as any).toBe('Title is required');
    (expect(router.navigate) as any).not.toHaveBeenCalled();
    (expect(component.submitting) as any).toBeFalse();
  }));

  it('should disable submit button while submitting', () => {
    component.submitting = true;
    fixture.detectChanges();
    const submitBtn = fixture.nativeElement.querySelector('#note-submit-btn');
    (expect(submitBtn.disabled) as any).toBeTrue();
  });
});
