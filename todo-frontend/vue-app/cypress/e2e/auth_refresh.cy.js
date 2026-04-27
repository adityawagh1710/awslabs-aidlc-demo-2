// Verifies the sessionStorage swap end-to-end:
//   1. Register lands tokens in sessionStorage (and NOT localStorage).
//   2. Reloading the page keeps the user signed in.
//   3. Authenticated API call (GET /todos) succeeds after reload.

const uniq = () => `cy_${Date.now()}_${Math.random().toString(36).slice(2, 8)}`

describe('auth + refresh', () => {
  it('keeps user logged in after a page reload (sessionStorage)', () => {
    const email = `${uniq()}@example.com`
    const password = 'Password1!'

    cy.visit('/register')
    cy.get('input[type="email"]').type(email)
    cy.get('#reg-pass').type(password)
    cy.get('#reg-confirm').type(password)
    cy.get('button[type="submit"]').click()

    // Successful register → app routes us to /todos.
    cy.url().should('include', '/todos')

    // Tokens should be in sessionStorage, not localStorage.
    cy.window().then(win => {
      expect(win.sessionStorage.getItem('access_token'), 'sessionStorage access_token').to.be.a('string').and.not.empty
      expect(win.sessionStorage.getItem('refresh_token'), 'sessionStorage refresh_token').to.be.a('string').and.not.empty
      expect(win.localStorage.getItem('access_token'), 'localStorage access_token').to.be.null
      expect(win.localStorage.getItem('refresh_token'), 'localStorage refresh_token').to.be.null
    })

    // Reload — same tab — should keep the session.
    cy.reload()
    cy.url().should('include', '/todos')
    cy.window().then(win => {
      expect(win.sessionStorage.getItem('access_token'), 'access_token survives reload').to.be.a('string').and.not.empty
    })

    // Authenticated API call after reload should succeed (not 401).
    cy.window().then(win => {
      const token = win.sessionStorage.getItem('access_token')
      cy.request({
        method: 'GET',
        url: '/api/todos/',
        headers: { Authorization: `Bearer ${token}` },
        failOnStatusCode: false,
      }).then(res => {
        expect(res.status, 'GET /todos/ after reload').to.eq(200)
      })
    })
  })

  it('logs out when the window is "closed" (sessionStorage cleared)', () => {
    const email = `${uniq()}@example.com`
    const password = 'Password1!'

    cy.visit('/register')
    cy.get('input[type="email"]').type(email)
    cy.get('#reg-pass').type(password)
    cy.get('#reg-confirm').type(password)
    cy.get('button[type="submit"]').click()
    cy.url().should('include', '/todos')

    // Simulate "window closed" by clearing sessionStorage and revisiting.
    cy.window().then(win => win.sessionStorage.clear())
    cy.visit('/')
    // App should bounce us to the login page when there's no session.
    cy.url().should('include', '/login')
  })
})
