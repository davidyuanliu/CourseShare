import { Injectable } from '@angular/core';

@Injectable({
  providedIn: 'root'
})
export class DemoControlService {
  errorMode = false;

  toggleErrorMode() {
    this.errorMode = !this.errorMode;
    console.log('Error mode is now:', this.errorMode ? 'ON' : 'OFF');
  }
}
