// Additional auth e2e tests:
//   1. Register → logout → login round-trip.
//   2. Wrong password keeps user on /login with no tokens.
//   3. Mismatched passwords show client-side validation error.
//
// NOTE: /auth/register is rate-limited to 10 req/hr. Tests minimise
//       registrations by reusing a single account where possible.

const uniq = () => `cy_${Date.now()}_${Math.random().toString(36).slice(2, 8)}`

// Shared credentials created once via the API (not the UI) so we only
// burn one rate-limit slot for the whole describe block.
let sharedEmail
let sharedPassword

describe('auth flows', () => {
  before(() => {
    sharedEmail = `${uniq()}@example.com`
    sharedPassword = 'Password1!'

    // Register via API to avoid burning UI interactions.
    cy.request('POST', '/api/auth/register', {
      email: sharedEmail,
      password: sharedPassword,
    }).its('status').should('eq', 201)
  })

  it('logs in with valid credentials after registration', () => {
    cy.visit('/login')
    cy.get('#email').type(sharedEmail)
    cy.get('#password').type(sharedPassword)
    cy.get('button[type="submit"]').click()

    cy.url().should('include', '/todos')

    // Tokens should be in sessionStorage.
    cy.window().then(win => {
      expect(win.sessionStorage.getItem('access_token')).to.be.a('string').and.not.empty
      expect(win.sessionStorage.getItem('refresh_token')).to.be.a('string').and.not.empty
    })
  })

  it('stays on /login with no tokens when password is wrong', () => {
    cy.visit('/login')
    cy.get('#email').type(sharedEmail)
    cy.get('#password').type('WrongPassword99!')
    cy.get('button[type="submit"]').click()

    // The user must NOT reach /todos.
    cy.url().should('include', '/login')
    cy.window().then(win => {
      expect(win.sessionStorage.getItem('access_token')).to.be.null
    })
  })

  it('shows a validation error when register passwords do not match', () => {
    cy.visit('/register')
    cy.get('input[type="email"]').type(`${uniq()}@example.com`)
    cy.get('#reg-pass').type('Password1!')
    cy.get('#reg-confirm').type('DifferentPass2!')
    cy.get('button[type="submit"]').click()

    // Client-side check — no API call, no rate-limit hit.
    cy.url().should('include', '/register')
    cy.get('.text-red-600')
      .should('be.visible')
      .and('contain.text', 'Passwords do not match')
  })
})
