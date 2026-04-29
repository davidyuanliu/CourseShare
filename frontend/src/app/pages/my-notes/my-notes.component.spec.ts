import { ComponentFixture, TestBed } from '@angular/core/testing';
import { MyNotesComponent } from './my-notes.component';
import { HttpClientTestingModule } from '@angular/common/http/testing';
import { RouterTestingModule } from '@angular/router/testing';
import { NoopAnimationsModule } from '@angular/platform-browser/animations';
import { CourseShareApiService } from '../../services/course-share-api.service';
import { of } from 'rxjs';

describe('MyNotesComponent', () => {
  let component: MyNotesComponent;
  let fixture: ComponentFixture<MyNotesComponent>;
  let apiService: CourseShareApiService;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [MyNotesComponent, HttpClientTestingModule, RouterTestingModule, NoopAnimationsModule]
    })
    .compileComponents();
    
    fixture = TestBed.createComponent(MyNotesComponent);
    component = fixture.componentInstance;
    apiService = TestBed.inject(CourseShareApiService);
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });

  it('should load my notes on init', () => {
    const mockNotes = [{ ID: 1, title: 'My Note', content: 'Content', courseId: 1, userId: 1, CreatedAt: '' }];
    spyOn(apiService, 'getMyNotes').and.returnValue(of(mockNotes));
    
    component.fetchMyNotes();
    expect(component.notes.length).toBe(1);
    expect(component.notes[0].title).toBe('My Note');
  });
});
