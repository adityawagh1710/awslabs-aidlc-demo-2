# Requirements Clarification Questions — Todo App

Please answer each question by filling in the letter choice after the `[Answer]:` tag.
If none of the options match, choose the last option and describe your preference.
Let me know when you're done.

---

## Question 1
What type of application should this be?

A) Web application (browser-based)
B) Mobile application (iOS/Android)
C) Desktop application
D) Command-line interface (CLI)
E) Other (please describe after [Answer]: tag below)

[Answer]: A

---

## Question 2
Who are the users of this application?

A) Single user only (personal/local use, no login needed)
B) Multiple users with individual accounts (authentication required)
C) Team/collaborative (multiple users sharing tasks)
D) Other (please describe after [Answer]: tag below)

[Answer]: B

---

## Question 3
What are the core features needed?

A) Basic CRUD only — create, view, complete, delete todos
B) Basic CRUD + due dates and priorities
C) Basic CRUD + due dates, priorities, categories/tags, and search
D) Full-featured — all of the above + reminders, recurring tasks, file attachments
E) Other (please describe after [Answer]: tag below)

[Answer]: D

---

## Question 4
Where should todo data be stored?

A) In-memory only (data lost on restart — suitable for demo/prototype)
B) Local file or embedded database (e.g., SQLite — single user, no server)
C) Backend database with a REST API (e.g., PostgreSQL + Node/Python/Java)
D) Cloud-managed database (e.g., AWS DynamoDB, Firebase)
E) Other (please describe after [Answer]: tag below)

[Answer]: C

---

## Question 5
What technology stack do you prefer?

A) JavaScript/TypeScript (React frontend, Node.js backend)
B) Python (Flask or FastAPI backend, optional frontend)
C) Java or Kotlin (Spring Boot backend)
D) No preference — let AI-DLC recommend based on requirements
E) Other (please describe after [Answer]: tag below)

[Answer]: GO lang with fiber framework for backend and latest vue for frontend 

---

## Question 6
What is the deployment target?

A) Local development only (no deployment needed)
B) Self-hosted server or VM
C) AWS (e.g., Lambda, ECS, Amplify, S3+CloudFront)
D) Other cloud provider (GCP, Azure, Vercel, Netlify, etc.)
E) Other (please describe after [Answer]: tag below)

[Answer]: B

---

## Question 7
What is the expected scale and performance requirement?

A) Personal use — very low traffic, no performance constraints
B) Small team — up to 50 users, moderate performance
C) Production-grade — hundreds or more users, performance and reliability matter
D) Other (please describe after [Answer]: tag below)

[Answer]: C

---

## Question: Security Extensions
Should security extension rules be enforced for this project?

A) Yes — enforce all SECURITY rules as blocking constraints (recommended for production-grade applications)
B) No — skip all SECURITY rules (suitable for PoCs, prototypes, and experimental projects)
X) Other (please describe after [Answer]: tag below)

[Answer]: A

---

## Question: Property-Based Testing Extension
Should property-based testing (PBT) rules be enforced for this project?

A) Yes — enforce all PBT rules as blocking constraints (recommended for projects with business logic, data transformations, serialization, or stateful components)
B) Partial — enforce PBT rules only for pure functions and serialization round-trips
C) No — skip all PBT rules (suitable for simple CRUD applications, UI-only projects, or thin integration layers)
X) Other (please describe after [Answer]: tag below)

[Answer]: B
