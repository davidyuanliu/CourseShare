import { TestBed } from '@angular/core/testing';
import { HttpClientTestingModule, HttpTestingController } from '@angular/common/http/testing';
import { AuthService } from './auth.service';

describe('AuthService', () => {
  let service: AuthService;
  let httpMock: HttpTestingController;

  beforeEach(() => {
    TestBed.configureTestingModule({
      imports: [HttpClientTestingModule]
    });
    service = TestBed.inject(AuthService);
    httpMock = TestBed.inject(HttpTestingController);
  });

  afterEach(() => {
    httpMock.verify();
  });

  it('should log in and save token', () => {
    const mockReponse = { token: 'mock-token', user: { id: 1, email: 'test@example.com' }, message: 'ok' };
    
    service.login({ email: 'test@example.com', password: '123' }).subscribe(res => {
      expect(res.token).toEqual('mock-token');
      expect(service.getToken()).toEqual('mock-token');
      expect(service.currentUserValue?.email).toEqual('test@example.com');
    });

    const req = httpMock.expectOne('https://courseshare.onrender.com/auth/login');
    expect(req.request.method).toBe('POST');
    req.flush(mockReponse);
  });
});
