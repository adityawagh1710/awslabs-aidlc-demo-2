// Tag management e2e tests:
//   1. Open tag manager panel.
//   2. Create a tag.
//   3. Delete a tag.
//   4. Create a todo with a tag attached.

const uniq = () => `cy_${Date.now()}_${Math.random().toString(36).slice(2, 8)}`

const email = `tag_${Date.now()}@example.com`
const password = 'Password1!'
let accessToken

describe('tag management', () => {
  before(() => {
    cy.request('POST', '/api/auth/register', { email, password }).then(res => {
      expect(res.status).to.eq(201)
      accessToken = res.body.access_token
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
  })

  it('opens the tag manager panel', () => {
    cy.contains('button', 'Tags').click()
    cy.contains('Manage Tags').should('be.visible')
    cy.get('input[placeholder="New tag name…"]').should('be.visible')
  })

  it('creates a new tag', () => {
    const tagName = `tag-${uniq()}`
    cy.contains('button', 'Tags').click()
    cy.contains('Manage Tags').should('be.visible')
    cy.get('input[placeholder="New tag name…"]').type(tagName)
    cy.contains('button', 'Add').click()
    cy.contains(tagName).should('be.visible')
  })

  it('deletes a tag', () => {
    const tagName = `del-${uniq()}`
    // Create via API
    cy.request({
      method: 'POST',
      url: '/api/tags/',
      headers: { Authorization: `Bearer ${accessToken}` },
      body: { name: tagName },
    })
    cy.visit('/todos')
    cy.contains('My Tasks').should('be.visible')
    cy.contains('button', 'Tags').click()
    cy.contains('Manage Tags').should('be.visible')
    cy.contains(tagName).parent().find('button').click()
    cy.contains(tagName).should('not.exist')
  })

  it('creates a todo with a tag attached', () => {
    const tagName = `att-${uniq()}`
    const todoTitle = `Tagged ${uniq()}`
    // Create tag via API
    cy.request({
      method: 'POST',
      url: '/api/tags/',
      headers: { Authorization: `Bearer ${accessToken}` },
      body: { name: tagName },
    })
    cy.visit('/todos')
    cy.contains('My Tasks').should('be.visible')

    cy.contains('button', 'New Task').click()
    cy.get('input[placeholder="What needs to be done?"]').type(todoTitle)
    // Select the tag badge in the modal
    cy.contains('button', tagName).click()
    cy.contains('button', 'Create Task').click()

    cy.contains(todoTitle).should('be.visible')
    cy.contains(todoTitle).closest('.card-hover').should('contain.text', tagName)
  })
})
