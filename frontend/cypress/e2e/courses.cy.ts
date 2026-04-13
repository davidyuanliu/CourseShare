describe('Courses Display and Navigation', () => {
  it('Visits the initial project page and shows courses', () => {
    cy.visit('/');
    
    // Check that we are redirected to /courses or we are on it
    cy.url().should('include', '/courses');
    
    // Check that there is a heading or the app title is present
    cy.contains('CourseShare').should('be.visible');

    // Make sure courses load (there might be a spinner first)
    // We check for mat-card which is used in courses.component.html
    cy.get('mat-card').should('have.length.greaterThan', 0);
    
    // Click on the first course card's 'View Notes' button to navigate
    cy.get('mat-card').first().contains('VIEW NOTES').click();

    // Verify it navigates to the notes list
    cy.url().should('match', /\/courses\/\d+\/notes/);
  });
});
