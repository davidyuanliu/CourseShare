import { ComponentFixture, TestBed } from '@angular/core/testing';
import { SavedNotesComponent } from './saved-notes.component';
import { HttpClientTestingModule } from '@angular/common/http/testing';
import { RouterTestingModule } from '@angular/router/testing';
import { NoopAnimationsModule } from '@angular/platform-browser/animations';
import { CourseShareApiService } from '../../services/course-share-api.service';
import { of } from 'rxjs';

describe('SavedNotesComponent', () => {
  let component: SavedNotesComponent;
  let fixture: ComponentFixture<SavedNotesComponent>;
  let apiService: CourseShareApiService;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [SavedNotesComponent, HttpClientTestingModule, RouterTestingModule, NoopAnimationsModule]
    })
    .compileComponents();
    
    fixture = TestBed.createComponent(SavedNotesComponent);
    component = fixture.componentInstance;
    apiService = TestBed.inject(CourseShareApiService);
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });

  it('should load saved notes on init', () => {
    const mockNotes = [{ ID: 1, title: 'Saved Note', content: 'Content', courseId: 1, userId: 2, CreatedAt: '' }];
    spyOn(apiService, 'getSavedNotes').and.returnValue(of(mockNotes));
    
    component.fetchSavedNotes();
    expect(component.notes.length).toBe(1);
    expect(component.notes[0].title).toBe('Saved Note');
  });

  it('should unsave a note', () => {
    const mockNoteId = 1;
    spyOn(apiService, 'unsaveNote').and.returnValue(of({ message: 'Unsaved', isSaved: false }));
    component.notes = [{ ID: 1, title: 'T', content: 'C', courseId: 1, userId: 2, CreatedAt: '' }];
    
    component.unsaveNote(mockNoteId, new MouseEvent('click'));
    expect(apiService.unsaveNote).toHaveBeenCalledWith(1);
    expect(component.notes.length).toBe(0);
  });
});
