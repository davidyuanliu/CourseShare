describe('Note Ownership & UI Restrictions', () => {
  const dummyNote = {
    ID: 100,
    title: 'Testing Cypress Note',
    content: 'Integration content',
    courseId: 1,
    userId: 1,
    CreatedAt: new Date().toISOString(),
    author: {
      id: 1,
      email: 'student@example.com'
    }
  };

  beforeEach(() => {
    cy.intercept('GET', 'https://unconstellated-ruthann-fiducially.ngrok-free.dev/notes/100', {
      statusCode: 200,
      body: dummyNote
    }).as('getNote');

    cy.intercept('PUT', 'https://unconstellated-ruthann-fiducially.ngrok-free.dev/notes/100', {
      statusCode: 200,
      body: { message: 'Note updated successfully', note: { ...dummyNote, title: 'Updated Title' } }
    }).as('updateNote');

    cy.intercept('DELETE', 'https://unconstellated-ruthann-fiducially.ngrok-free.dev/notes/100', {
      statusCode: 200,
      body: { message: 'Note deleted successfully' }
    }).as('deleteNote');
  });

  it('hides edit/delete buttons for logged out users', () => {
    window.localStorage.removeItem('token');
    window.localStorage.removeItem('currentUser');
    
    cy.visit('/notes/100');
    cy.wait('@getNote');

    cy.contains('Testing Cypress Note');
    cy.get('[aria-label="Edit Note"]').should('not.exist');
    cy.get('[aria-label="Delete Note"]').should('not.exist');
  });

  it('hides edit/delete buttons for authenticated non-owners', () => {
    window.localStorage.setItem('token', 'fake-token-non-owner');
    window.localStorage.setItem('currentUser', JSON.stringify({ id: 2, email: 'other@example.com' })); // ID 2 is not owner ID 1
    
    cy.visit('/notes/100');
    cy.wait('@getNote');

    cy.contains('Testing Cypress Note');
    cy.get('[aria-label="Edit Note"]').should('not.exist');
    cy.get('[aria-label="Delete Note"]').should('not.exist');
  });

  it('shows edit/delete buttons and allows editing for the owner', () => {
    window.localStorage.setItem('token', 'fake-token');
    window.localStorage.setItem('currentUser', JSON.stringify({ id: 1, email: 'student@example.com' }));
    
    cy.visit('/notes/100');
    cy.wait('@getNote');

    cy.contains('Testing Cypress Note');
    cy.get('[aria-label="Edit Note"]').should('exist').click();

    cy.url().should('include', '/edit');
    cy.wait('@getNote'); // Edit component refetches the note to fill form

    // Form should be pre-filled
    cy.get('input[formControlName="title"]').should('have.value', 'Testing Cypress Note');
    
    cy.get('input[formControlName="title"]').clear().type('Updated Title');
    cy.get('button[type="submit"]').click();

    cy.wait('@updateNote');
    cy.contains('Note updated successfully');
  });

  it('allows deleting the note as an owner with confirmation', () => {
    window.localStorage.setItem('token', 'fake-token');
    window.localStorage.setItem('currentUser', JSON.stringify({ id: 1, email: 'student@example.com' }));
    
    cy.visit('/notes/100');
    cy.wait('@getNote');

    cy.get('[aria-label="Delete Note"]').click();
    
    // Expect dialog to appear
    cy.contains('Confirm Deletion').should('be.visible');
    
    // Click Delete in dialog
    cy.contains('button', 'Delete').click();

    cy.wait('@deleteNote');
    cy.contains('Note deleted successfully');
    cy.url().should('include', '/courses/1/notes');
  });
});
