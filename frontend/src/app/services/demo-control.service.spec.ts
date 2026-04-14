import { TestBed } from '@angular/core/testing';

import { DemoControlService } from './demo-control.service';

describe('DemoControlService', () => {
  let service: DemoControlService;

  beforeEach(() => {
    TestBed.configureTestingModule({});
    service = TestBed.inject(DemoControlService);
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });
});
