import { ComponentFixture, TestBed } from '@angular/core/testing';
import { EditNoteComponent } from './edit-note.component';
import { HttpClientTestingModule } from '@angular/common/http/testing';
import { RouterTestingModule } from '@angular/router/testing';
import { provideAnimations } from '@angular/platform-browser/animations';

describe('EditNoteComponent', () => {
  let component: EditNoteComponent;
  let fixture: ComponentFixture<EditNoteComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [EditNoteComponent, HttpClientTestingModule, RouterTestingModule],
      providers: [provideAnimations()]
    })
    .compileComponents();
    
    fixture = TestBed.createComponent(EditNoteComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });

  it('should initialize form', () => {
    expect(component.noteForm).toBeDefined();
    expect(component.noteForm.controls['title']).toBeDefined();
    expect(component.noteForm.controls['content']).toBeDefined();
  });

  it('should validate form fields', () => {
    component.noteForm.controls['title'].setValue('');
    expect(component.noteForm.controls['title'].valid).toBeFalsy();
    
    component.noteForm.controls['title'].setValue('Test Title');
    expect(component.noteForm.controls['title'].valid).toBeTruthy();
  });
});
