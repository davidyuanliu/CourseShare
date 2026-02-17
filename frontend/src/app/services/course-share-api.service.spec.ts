import { TestBed } from '@angular/core/testing';

import { CourseShareApiService } from './course-share-api.service';

describe('CourseShareApiService', () => {
  let service: CourseShareApiService;

  beforeEach(() => {
    TestBed.configureTestingModule({});
    service = TestBed.inject(CourseShareApiService);
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });
});
