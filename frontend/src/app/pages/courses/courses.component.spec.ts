import { ComponentFixture, TestBed } from '@angular/core/testing';
import { provideRouter } from '@angular/router';
import { of } from 'rxjs';
import { BrowserAnimationsModule } from '@angular/platform-browser/animations';

import { CoursesComponent } from './courses.component';
import { CourseShareApiService } from '../../services/course-share-api.service';
import { Course } from '../../mock/mock-data.service';

describe('CoursesComponent', () => {
  let component: CoursesComponent;
  let fixture: ComponentFixture<CoursesComponent>;
  let apiServiceSpy: jasmine.SpyObj<CourseShareApiService>;

  beforeEach(async () => {
    const spy = jasmine.createSpyObj('CourseShareApiService', ['getCourses', 'createCourse']);
    
    // Set up default mock return value before component initializes
    spy.getCourses.and.returnValue(of([
      { ID: 1, CreatedAt: '', UpdatedAt: '', name: 'Test Course 1' },
      { ID: 2, CreatedAt: '', UpdatedAt: '', name: 'Test Course 2' }
    ]));

    await TestBed.configureTestingModule({
      imports: [CoursesComponent, BrowserAnimationsModule],
      providers: [
        provideRouter([]),
        { provide: CourseShareApiService, useValue: spy }
      ]
    })
    .compileComponents();
    
    apiServiceSpy = TestBed.inject(CourseShareApiService) as jasmine.SpyObj<CourseShareApiService>;
    
    fixture = TestBed.createComponent(CoursesComponent);
    component = fixture.componentInstance;
    fixture.detectChanges(); // This triggers ngOnInit -> fetchCourses
  });

  it('should create', () => {
    (expect(component) as any).toBeTruthy();
  });

  it('should load courses on init', () => {
    (expect(apiServiceSpy.getCourses.calls.any()) as any).withContext('should call getCourses').toBe(true);
    (expect(component.courses.length) as any).toBe(2);
    (expect(component.courses[0].name) as any).toBe('Test Course 1');
    (expect(component.loading) as any).toBe(false);
  });
});
