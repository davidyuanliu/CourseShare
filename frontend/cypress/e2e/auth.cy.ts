describe('Authentication Flows', () => {
  beforeEach(() => {
    // Intercept to avoid real backend logic bleeding over if the DB state is dirty
    cy.intercept('POST', '**/auth/register', {
      statusCode: 201,
      body: {
        message: 'User registered successfully',
        token: 'fake-jwt-token',
        user: { id: 1, email: 'student@example.com' }
      }
    }).as('registerRequest');

    cy.intercept('POST', '**/auth/login', {
      statusCode: 200,
      body: {
        message: 'Login successful',
        token: 'fake-jwt-token',
        user: { id: 1, email: 'student@example.com' }
      }
    }).as('loginRequest');
  });

  it('allows a user to navigate to register and submit form', () => {
    cy.visit('/');
    cy.contains('Register').click();
    cy.url().should('include', '/register');

    cy.get('input[formControlName="email"]').type('student@example.com');
    cy.get('input[formControlName="password"]').type('password123');
    cy.get('button[type="submit"]').click();

    cy.wait('@registerRequest');
    cy.contains('Registration successful');
    cy.url().should('include', '/courses');
    
    // Auth state changes navbar
    cy.contains('Logout');
  });

  it('allows a user to login and correctly updates auth state', () => {
    cy.visit('/login');

    cy.get('input[formControlName="email"]').type('student@example.com');
    cy.get('input[formControlName="password"]').type('password123');
    cy.get('button[type="submit"]').click();

    cy.wait('@loginRequest');
    cy.contains('Login successful');
    cy.url().should('include', '/courses');
    
    // Auth state changes navbar
    cy.contains('Logout');
  });

  it('navigates to login when logging out', () => {
    // Seed localstorage to simulate logged in
    window.localStorage.setItem('token', 'fake-jwt-token');
    window.localStorage.setItem('currentUser', JSON.stringify({ id: 1, email: 'student@example.com' }));
    
    cy.visit('/');
    cy.contains('Logout').click();
    
    cy.contains('Login');
    cy.contains('Register');
  });
});
