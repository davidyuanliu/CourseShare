import { TestBed } from '@angular/core/testing';
import { HttpTestingController, provideHttpClientTesting } from '@angular/common/http/testing';
import { provideHttpClient } from '@angular/common/http';

import { CourseShareApiService } from './course-share-api.service';
import { Course } from '../mock/mock-data.service';

describe('CourseShareApiService', () => {
  let service: CourseShareApiService;
  let httpMock: HttpTestingController;

  beforeEach(() => {
    TestBed.configureTestingModule({
      providers: [
        CourseShareApiService,
        provideHttpClient(),
        provideHttpClientTesting()
      ]
    });
    service = TestBed.inject(CourseShareApiService);
    httpMock = TestBed.inject(HttpTestingController);
  });

  afterEach(() => {
    httpMock.verify();
  });

  it('should be created', () => {
    (expect(service) as any).toBeTruthy();
  });

  it('should fetch courses successfully', () => {
    const mockCourses: Course[] = [
      { ID: 1, CreatedAt: '', UpdatedAt: '', name: 'Software Engineering' },
      { ID: 2, CreatedAt: '', UpdatedAt: '', name: 'Database Systems' }
    ];

    service.getCourses().subscribe(courses => {
      (expect(courses.length) as any).toBe(2);
      (expect(courses) as any).toEqual(mockCourses);
    });

    const req = httpMock.expectOne('http://localhost:8080/courses');
    (expect(req.request.method) as any).toBe('GET');
    req.flush(mockCourses);
  });
});
