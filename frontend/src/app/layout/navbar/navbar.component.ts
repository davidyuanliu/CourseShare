import { Component } from '@angular/core';
import { MatToolbarModule } from '@angular/material/toolbar';
import { MatButtonModule } from '@angular/material/button';
import { RouterLink, RouterLinkActive } from '@angular/router';
import { DemoControlService } from '../../services/demo-control.service';
import { MatSlideToggleModule } from '@angular/material/slide-toggle';
import { AuthService } from '../../services/auth.service';
import { CommonModule } from '@angular/common';

@Component({
  selector: 'app-navbar',
  standalone: true,
  imports: [
    MatToolbarModule,
    MatButtonModule,
    RouterLink,
    RouterLinkActive,
    MatSlideToggleModule,
    CommonModule
  ],
  templateUrl: './navbar.component.html',
  styleUrl: './navbar.component.css'
})
export class NavbarComponent {
  constructor(
      public demoControl: DemoControlService,
      public authService: AuthService
  ) { }

  logout(): void {
    this.authService.logout();
  }
}
