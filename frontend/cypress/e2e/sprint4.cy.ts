describe('Sprint 4 Features', () => {
  const mockNote = {
    ID: 1,
    title: 'Sprint 4 Note',
    content: 'Testing features',
    courseId: 1,
    userId: 1,
    authorName: 'Test Author',
    tags: ['tag1', 'tag2'],
    helpfulCount: 5,
    isHelpful: false,
    isSaved: false,
    CreatedAt: new Date().toISOString(),
    author: { id: 1, email: 'owner@test.com' }
  };

  beforeEach(() => {
    window.localStorage.setItem('token', 'fake-token');
    window.localStorage.setItem('currentUser', JSON.stringify({ id: 1, email: 'owner@test.com' }));
    
    cy.intercept('GET', 'https://courseshare.onrender.com/notes/1', { statusCode: 200, body: mockNote }).as('getNote');
    cy.intercept('GET', 'https://courseshare.onrender.com/courses/1/notes', { statusCode: 200, body: [mockNote] }).as('getNotes');
    cy.intercept('POST', 'https://courseshare.onrender.com/notes/1/helpful', { statusCode: 200, body: { helpfulCount: 6, isHelpful: true } }).as('markHelpful');
    cy.intercept('DELETE', 'https://courseshare.onrender.com/notes/1/helpful', { statusCode: 200, body: { helpfulCount: 5, isHelpful: false } }).as('unmarkHelpful');
    cy.intercept('POST', 'https://courseshare.onrender.com/notes/1/save', { statusCode: 200, body: { isSaved: true } }).as('saveNote');
    cy.intercept('DELETE', 'https://courseshare.onrender.com/notes/1/save', { statusCode: 200, body: { isSaved: false } }).as('unsaveNote');
    cy.intercept('GET', 'https://courseshare.onrender.com/my-notes', { statusCode: 200, body: [mockNote] }).as('getMyNotes');
    cy.intercept('GET', 'https://courseshare.onrender.com/saved-notes', { statusCode: 200, body: [mockNote] }).as('getSavedNotes');
  });

  it('displays author name if provided', () => {
    cy.visit('/notes/1');
    cy.wait('@getNote');
    cy.contains('Uploaded by Test Author');
  });

  it('toggles helpful and save status', () => {
    cy.visit('/notes/1');
    cy.wait('@getNote');

    // Toggle Helpful
    cy.get('#helpful-btn').click();
    cy.wait('@markHelpful');
    cy.get('.helpful-count').contains('6');

    // Toggle Save
    cy.get('#save-btn').click();
    cy.wait('@saveNote');
    cy.get('#save-btn').should('have.css', 'color', 'rgb(255, 64, 129)'); // accent color
  });

  it('filters notes in the list', () => {
    cy.visit('/courses/1/notes');
    cy.wait('@getNotes');

    // Search
    cy.get('input[placeholder="Keywords..."]').type('Sprint');
    cy.get('.note-card').should('have.length', 1);

    // Personal filters
    cy.get('#personal-filters-select').click({ force: true });
    cy.get('mat-option').contains('Written by me').click();
    cy.get('body').click(); // close select
    cy.get('.note-card').should('have.length', 1);
  });

  it('shows dashboards', () => {
    // My Notes
    cy.visit('/my-notes');
    cy.wait('@getMyNotes');
    cy.contains('Sprint 4 Note');

    // Saved Notes
    cy.visit('/saved-notes');
    cy.wait('@getSavedNotes');
    cy.contains('Sprint 4 Note');
  });
});
