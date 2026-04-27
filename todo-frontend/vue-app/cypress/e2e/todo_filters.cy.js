// Todo filter & search e2e tests:
//   1. Filter by status dropdown.
//   2. Filter by priority dropdown.
//   3. Search by text.
//   4. Clear filters.
//   5. Empty state when no matches.

const uniq = () => `cy_${Date.now()}_${Math.random().toString(36).slice(2, 8)}`

const email = `filter_${Date.now()}@example.com`
const password = 'Password1!'
let accessToken

const prefix = `F${Date.now().toString(36)}`

describe('todo filters & search', () => {
  before(() => {
    cy.request('POST', '/api/auth/register', { email, password }).then(res => {
      expect(res.status).to.eq(201)
      accessToken = res.body.access_token

      // Seed 3 todos with different status/priority combos.
      const todos = [
        { title: `${prefix} Alpha`, priority: 'high' },
        { title: `${prefix} Beta`,  priority: 'medium' },
        { title: `${prefix} Gamma`, priority: 'low' },
      ]
      todos.forEach(t => {
        cy.request({
          method: 'POST',
          url: '/api/todos/',
          headers: { Authorization: `Bearer ${accessToken}` },
          body: t,
        }).its('status').should('eq', 201)
      })

      // Update Beta → in_progress, Gamma → done (must go through in_progress first)
      cy.request({
        method: 'GET',
        url: '/api/todos/',
        headers: { Authorization: `Bearer ${accessToken}` },
      }).then(listRes => {
        const beta  = listRes.body.find(t => t.title.includes('Beta'))
        const gamma = listRes.body.find(t => t.title.includes('Gamma'))
        if (beta) {
          cy.request({
            method: 'PUT',
            url: `/api/todos/${beta.id}`,
            headers: { Authorization: `Bearer ${accessToken}` },
            body: { status: 'in_progress' },
          })
        }
        if (gamma) {
          cy.request({
            method: 'PUT',
            url: `/api/todos/${gamma.id}`,
            headers: { Authorization: `Bearer ${accessToken}` },
            body: { status: 'in_progress' },
          }).then(() => {
            cy.request({
              method: 'PUT',
              url: `/api/todos/${gamma.id}`,
              headers: { Authorization: `Bearer ${accessToken}` },
              body: { status: 'done' },
            })
          })
        }
      })
    })
  })

  beforeEach(() => {
    cy.session(email, () => {
      cy.visit('/login')
      cy.get('#email').type(email)
      cy.get('#password').type(password)
      cy.get('button[type="submit"]').click()
      cy.url().should('include', '/todos')
    })
    cy.visit('/todos')
    cy.contains('My Tasks').should('be.visible')
    // Ensure all 3 seeded todos are loaded
    cy.contains(`${prefix} Alpha`).should('be.visible')
  })

  it('filters by status', () => {
    cy.get('select').eq(0).select('in_progress')
    cy.contains(`${prefix} Beta`).should('be.visible')
    cy.contains(`${prefix} Alpha`).should('not.exist')
    cy.contains(`${prefix} Gamma`).should('not.exist')
  })

  it('filters by priority', () => {
    cy.get('select').eq(1).select('high')
    cy.contains(`${prefix} Alpha`).should('be.visible')
    cy.contains(`${prefix} Beta`).should('not.exist')
    cy.contains(`${prefix} Gamma`).should('not.exist')
  })

  it('searches by text', () => {
    cy.get('input[placeholder="Search tasks…"]').type('Gamma')
    cy.contains(`${prefix} Gamma`).should('be.visible')
    cy.contains(`${prefix} Alpha`).should('not.exist')
    cy.contains(`${prefix} Beta`).should('not.exist')
  })

  it('clears all filters', () => {
    cy.get('select').eq(0).select('done')
    cy.contains(`${prefix} Alpha`).should('not.exist')
    cy.contains('button', 'Clear').click()
    cy.contains(`${prefix} Alpha`).should('be.visible')
    cy.contains(`${prefix} Beta`).should('be.visible')
    cy.contains(`${prefix} Gamma`).should('be.visible')
  })

  it('shows empty state when no matches', () => {
    cy.get('input[placeholder="Search tasks…"]').type('zzz_no_match_ever')
    cy.contains('No matching tasks').should('be.visible')
    cy.contains('button', 'Clear filters').should('be.visible')
  })
})
