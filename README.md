

# PFE-LMS &nbsp;&horbar;&nbsp; Plateforme de Gestion des Projets de Fin d'Etudes

<br>

<div align="center">

  <img src="https://go.dev/blog/go-brand/Go-Logo/SVG/Go-Logo_Blue.svg" height="60" alt="Go" />&nbsp;&nbsp;&nbsp;&nbsp;
  <img src="https://raw.githubusercontent.com/sveltejs/branding/refs/heads/master/svelte-logo.svg" height="60" alt="SvelteKit" />&nbsp;&nbsp;&nbsp;&nbsp;
  <img src="https://www.sqlite.org/images/sqlite370_banner.gif" height="50" alt="SQLite" />

  <br><br>

  <p align="center">
    A <strong>full-stack</strong> academic management platform for orchestrating the entire PFE <em>(Projet de Fin d'Etudes)</em> lifecycle &mdash; from subject proposal and validation, through team assignment and supervision, to defense scheduling, jury grading, and final transcript generation. Built with a <strong>Go + Fiber v3</strong> backend and a <strong>SvelteKit (Svelte 5)</strong> frontend.
  </p>

  <a href="https://github.com/rhpo/lms">
    <strong>
      Explore the repository &raquo;
    </strong>
  </a>

  <br><br>

  <a href="https://github.com/rhpo/lms/issues">Report Bug</a>
  &middot;
  <a href="https://github.com/rhpo/lms/issues">Request Feature</a>

  <br><br>

  ![Go](https://img.shields.io/badge/Go-1.23-00ADD8?style=for-the-badge&logo=go&logoColor=white)
  ![Svelte](https://img.shields.io/badge/Svelte-5-FF3E00?style=for-the-badge&logo=svelte&logoColor=white)
  ![Fiber](https://img.shields.io/badge/Fiber-v3-00ACD7?style=for-the-badge)
  ![SQLite](https://img.shields.io/badge/SQLite-3-003B57?style=for-the-badge&logo=sqlite&logoColor=white)
  ![License](https://img.shields.io/badge/License-MIT-green?style=for-the-badge)

</div>

<br>

---

<details>
  <summary><strong>Table of Contents</strong></summary>
  <ol>
    <li><a href="#about-the-project">About The Project</a></li>
    <li><a href="#architecture">Architecture</a></li>
    <li><a href="#backend-deep-dive">Backend Deep-Dive</a></li>
    <li><a href="#frontend-overview">Frontend Overview</a></li>
    <li><a href="#user-roles">User Roles</a></li>
    <li><a href="#features">Features</a></li>
    <li><a href="#getting-started">Getting Started</a></li>
    <li><a href="#project-structure">Project Structure</a></li>
    <li><a href="#api-endpoints">API Endpoints</a></li>
    <li><a href="#license">License</a></li>
    <li><a href="#contact">Contact</a></li>
  </ol>
</details>

<br>

<h2 id="about-the-project">&bull; About The Project</h2>

**PFE-LMS** is a comprehensive Learning Management System purpose-built for Algerian universities (and similar academic institutions) to manage the full lifecycle of final-year graduation projects.

The platform digitizes and streamlines a process that traditionally relies on paper forms, email chains, and Excel spreadsheets &mdash; replacing it with a role-aware, real-time web application that keeps students, teachers, administrators, and partner companies in sync.

### What it solves

- **Fragmented communication** between students, supervisors, and administration
- **Manual subject validation** workflows that are slow and error-prone
- **Opaque assignment processes** where students can't track their PFE status
- **Unstructured defense scheduling** and grade resolution across jury members
- **Zero traceability** on supervision meetings and progress reports

<br>

<h2 id="architecture">&bull; Architecture</h2>

```
┌──────────────────────────────────────────────────────────────────┐
│                         Client (Browser)                         │
│                    SvelteKit 2 · Svelte 5 Runes                  │
│                         SSR disabled (SPA)                       │
└────────────────────────────┬─────────────────────────────────────┘
                             │  JSON / REST
                             ▼
┌──────────────────────────────────────────────────────────────────┐
│                     Go Backend (Fiber v3)                         │
│                                                                   │
│  ┌─────────┐   ┌───────────┐   ┌──────────────┐                 │
│  │ Handler  │──▶│  Service   │──▶│  Repository   │                │
│  │  (HTTP)  │   │ (Business) │   │   (SQL/DB)    │                │
│  └─────────┘   └───────────┘   └──────┬───────┘                 │
│       │                                │                          │
│       │         ┌──────────────────────┘                          │
│       │         ▼                                                 │
│  ┌─────────────────────┐   ┌───────────────────┐                 │
│  │   SQLite Database    │   │  Shared Packages   │                │
│  │  (single-file, WAL) │   │  notify · auth ·   │                │
│  └─────────────────────┘   │  apperror · ...    │                │
│                             └───────────────────┘                 │
│                                      │                            │
│                          ┌───────────┴──────────┐                │
│                          │   Resend Email API    │                │
│                          └──────────────────────┘                │
└──────────────────────────────────────────────────────────────────┘
```

The system follows a **strict layered architecture**:

| Layer | Responsibility | Location |
|-------|---------------|----------|
| **Handler** | HTTP parsing, validation, response formatting | `backend/internal/handler/` |
| **Service** | Business logic, authorization, orchestration | `backend/internal/service/` |
| **Repository** | Raw SQL queries, data access | `backend/internal/repository/` |
| **Entity** | Domain models (1:1 with DB schema) | `backend/internal/entity/` |
| **Shared** | Cross-cutting: auth middleware, notifications, error types, validators | `backend/internal/shared/` |

<br>

<h2 id="backend-deep-dive">&bull; Backend Deep-Dive</h2>

### Tech Stack

| Component | Technology |
|-----------|-----------|
| Language | **Go 1.23** |
| HTTP Framework | **Fiber v3** (fasthttp-based) |
| Database | **SQLite** (WAL mode, single-file) |
| Auth | **JWT** (HS256, configurable expiry) |
| Email | **Resend API** (transactional emails) |
| File Storage | Local filesystem (`/uploads`) |

### Notification System

The backend features a **dual-channel notification system** (`shared/notify`):

```
notifier.Send(recipientProfileID, notifyType, message)
    ├──▶ DB Channel   → INSERT into notifications table (in-app)
    └──▶ Email Channel → Resend API (transactional email)
```

- `notifier.Send()` fans out to all registered channels
- `notifier.NotifyAdmins()` broadcasts to all admin profiles
- Email channel is conditionally registered (only when a valid Resend API key is configured)
- Notification types: `validation_requise`, `affectation`, `jury`, `disponibilite`, `sujet`

### Grading Workflow

The grading system implements a **multi-step jury evaluation process**:

1. **Member (Examinateur)** submits their individual evaluation (4 criteria, each /4) + archive decision
2. **Supervisor** submits their evaluation (criterion 5, /4)
3. **President** reviews both evaluations, then either adopts the member's grades or enters new ones
4. Final grade = sum of 5 criteria = **/20**
5. Student is automatically notified upon grade finalization

### Key Design Decisions

- **Profile ID vs Entity ID**: The system distinguishes between `profiles.id` (auth identity) and `teachers.id`/`students.id` (domain entity). All foreign keys in `defense_juries` reference teacher entity IDs, while auth middleware provides profile IDs — the service layer handles the mapping.
- **NullString / NullFloat64 wrappers**: Custom JSON-serializable nullable types for SQLite compat.
- **Auto-migration**: Schema migrations run on startup in `main.go` via `ALTER TABLE` statements.
- **Zero-ORM**: All SQL is hand-written for full control and performance.

<br>

<h2 id="frontend-overview">&bull; Frontend Overview</h2>

### Tech Stack

| Component | Technology |
|-----------|-----------|
| Framework | **SvelteKit 2** |
| UI Library | **Svelte 5** (Runes: `$state`, `$derived`, `$effect`, `$props`) |
| Icons | **Lucide Svelte** |
| Charts | **Chart.js** |
| Styling | **CSS custom properties** (design tokens), scoped `<style>` |
| Build | **Vite** |
| SSR | Disabled (SPA mode) |

### Component Library

A custom, purpose-built component library (`src/lib/components/ui/`):

`AppShell` · `Avatar` · `Badge` · `Button` · `Calendar` · `DateInput` · `FormField` · `Modal` · `Notification` · `Page` · `Pagination` · `SearchInput` · `Select` · `Switch` · `Table` · `Tabs` · `ThemeToggle` · `Tooltip`

### API Layer

A centralized API module (`src/lib/api.ts`) provides typed functions grouped by role:

```typescript
admin.listSubjects()       // → GET  /api/admin/subjects
teacher.listJuryDuties()   // → GET  /api/teacher/jury-duties
student.getSoutenance()     // → GET  /api/student/soutenance
company.submitSubject(...)  // → POST /api/company/subjects
```

<br>

<h2 id="user-roles">&bull; User Roles</h2>

| Role | Dashboard | Key Capabilities |
|------|-----------|-----------------|
| **Admin** | Full oversight | Manage users, validate subjects, assign PFEs, schedule defenses, resolve grades, audit logs, academic year settings |
| **Teacher** | Teaching portal | Propose subjects, validate peer subjects, supervise PFEs, track meetings, evaluate students, serve on juries |
| **Student** | Student portal | Browse catalogue, submit wishes, track PFE progress, log meetings, upload thesis, view defense & grades |
| **Company** | Partner portal | Propose industry subjects, co-supervise PFEs, track assignment status |

<br>

<h2 id="features">&bull; Features</h2>

### Subject Management
- Teachers and companies can propose PFE subjects
- Dual-validator review workflow (each validator independently approves/rejects)
- Subjects tagged with domains for smart jury recommendations
- Status tracking: `en_attente` → `valide` / `accepte_sous_reserve` / `refuse`

### Student Assignment
- Students browse a catalogue and submit ranked wishes
- Teachers accept candidates (supports monome/binome/trinome)
- Auto-generated PFE codes upon assignment
- Catalogue shows real-time availability ("Deja pris" for assigned subjects)

### Supervision & Progress
- Structured meeting logs (date, duration, type, topics, observations)
- Progress report status tracking (`a_faire` → `en_cours` → `termine`)
- Thesis (memoire) PDF upload and tracking

### Defense & Jury
- Admin schedules defenses with room, date, and jury composition
- Smart jury recommendations based on domain overlap
- President/member role distinction enforced in UI and backend
- Jury confirmation workflow with print preference

### Grading
- 4-criterion evaluation (/4 each) by jury + supervisor note (/4) = /20
- Archive decision per jury member (archivable / minor corrections / major corrections)
- President finalizes grade with choice to adopt or override
- Automatic student notification on grade publication

### Notifications
- Real-time in-app notifications per role
- Transactional email delivery via Resend
- Notification types: validation, assignment, jury, availability, subject updates

### Administration
- Full audit log of system actions
- Academic year management with submission windows
- Department, speciality, and promotion management
- User account management (activate/deactivate, role assignment)
- Bulk data export (Excel/CSV)

<br>

<h2 id="getting-started">&bull; Getting Started</h2>

### Prerequisites

- **Go** 1.23+
- **Node.js** 20+ & **pnpm**
- (Optional) **Resend** API key for email notifications

### Installation

```bash
# Clone the repository
git clone https://github.com/rhpo/lms.git
cd lms

# ── Backend ──
cd backend
go mod download
cp .env.example .env      # configure JWT_SECRET, RESEND_KEY, etc.
go run cmd/server/main.go  # starts on :8080

# ── Frontend ──
cd ..
pnpm install
pnpm dev                   # starts on :5173, proxies /api to :8080
```

### Environment Variables

| Variable | Description | Default |
|----------|------------|---------|
| `JWT_SECRET` | HMAC signing key for JWT tokens | *(required)* |
| `RESEND_API_KEY` | Resend API key for email delivery | *(optional, disables email if unset)* |
| `DB_PATH` | SQLite database file path | `./pfe.db` |
| `PORT` | Server listen port | `8080` |

<br>

<h2 id="project-structure">&bull; Project Structure</h2>

```
lms/
├── backend/
│   ├── cmd/server/main.go          # Entrypoint, routes, migrations
│   ├── internal/
│   │   ├── config/                  # App configuration
│   │   ├── entity/                  # Domain models (Go structs ↔ DB rows)
│   │   ├── handler/                 # HTTP handlers (admin, teacher, student, company)
│   │   ├── repository/              # SQL queries (one repo per entity)
│   │   ├── service/                 # Business logic layer
│   │   └── shared/
│   │       ├── apperror/            # Typed application errors
│   │       ├── middleware/           # JWT auth, role guards
│   │       ├── notify/              # Dual-channel notification system
│   │       ├── pfe_code/            # PFE code generator
│   │       ├── response/            # Standardized JSON responses
│   │       └── validator/           # Input validation helpers
│   ├── tests/                       # Integration tests + test helpers
│   └── uploads/                     # File storage (thesis PDFs, avatars)
│
├── src/
│   ├── lib/
│   │   ├── api.ts                   # Typed API client (admin/teacher/student/company)
│   │   ├── components/ui/           # Reusable UI component library
│   │   ├── constants/               # Branding, config constants
│   │   ├── stores/                  # Svelte stores (auth, theme)
│   │   └── types/                   # TypeScript domain types (1:1 with Go entities)
│   ├── routes/
│   │   ├── (app)/                   # Public routes (landing, login)
│   │   └── (dashboard)/
│   │       ├── admin/               # Admin panel (15+ pages)
│   │       ├── teacher/             # Teacher portal (8 sections)
│   │       ├── student/             # Student portal (9 sections)
│   │       └── company/             # Company portal
│   └── app.html                     # SPA shell
│
├── static/                          # Static assets (fonts, media)
├── svelte.config.js
├── vite.config.ts
└── package.json
```

<br>

<h2 id="api-endpoints">&bull; API Endpoints</h2>

<details>
<summary><strong>Admin</strong> &mdash; 30+ endpoints</summary>

| Method | Endpoint | Description |
|--------|----------|-------------|
| `GET` | `/api/admin/dashboard` | Dashboard statistics |
| `GET` | `/api/admin/subjects` | List all subjects |
| `POST` | `/api/admin/subjects/:id/assign-validators` | Assign peer reviewers |
| `GET` | `/api/admin/pfes` | List all PFE assignments |
| `POST` | `/api/admin/pfes/:id/co-supervisor` | Assign co-supervisor |
| `POST` | `/api/admin/defenses` | Schedule a defense |
| `GET` | `/api/admin/defenses/recommend-jury` | AI-assisted jury recommendations |
| `POST` | `/api/admin/defenses/:id/resolve-grade` | Resolve final grade |
| `GET` | `/api/admin/settings/deadlines` | Academic year settings |
| `GET` | `/api/admin/audit-log` | System audit trail |

</details>

<details>
<summary><strong>Teacher</strong> &mdash; 15+ endpoints</summary>

| Method | Endpoint | Description |
|--------|----------|-------------|
| `GET` | `/api/teacher/dashboard` | Teacher dashboard |
| `GET` | `/api/teacher/proposed-subjects` | My proposed subjects |
| `GET` | `/api/teacher/subjects-to-validate` | Subjects awaiting my review |
| `POST` | `/api/teacher/subjects-to-validate/:id/validate` | Submit validation decision |
| `GET` | `/api/teacher/supervised-pfes` | My supervised PFEs |
| `POST` | `/api/teacher/supervised-pfes/:id/evaluation` | Submit supervisor evaluation |
| `GET` | `/api/teacher/jury-duties` | My jury assignments |
| `GET` | `/api/teacher/jury-duties/:id/grade-context` | Get grading context (role, existing grades) |
| `POST` | `/api/teacher/jury-duties/:id/grade` | Submit member evaluation |
| `POST` | `/api/teacher/jury-duties/:id/final-grade` | Submit final grade (president only) |

</details>

<details>
<summary><strong>Student</strong> &mdash; 10+ endpoints</summary>

| Method | Endpoint | Description |
|--------|----------|-------------|
| `GET` | `/api/student/dashboard` | Student dashboard |
| `GET` | `/api/student/catalogue` | Browse available subjects |
| `POST` | `/api/student/voeux` | Submit a wish |
| `GET` | `/api/student/my-pfe` | My PFE assignment |
| `POST` | `/api/student/my-pfe/memoire` | Upload thesis |
| `GET` | `/api/student/soutenance` | My defense info & grades |

</details>

<details>
<summary><strong>Company</strong></summary>

| Method | Endpoint | Description |
|--------|----------|-------------|
| `GET` | `/api/company/dashboard` | Company dashboard |
| `POST` | `/api/company/subjects` | Propose a subject |
| `GET` | `/api/company/pfes` | Track assigned PFEs |

</details>

<br>

<h2 id="license">&bull; License</h2>

Distributed under the **MIT License**. See `LICENSE` for more information.

<br>

<h2 id="contact">&bull; Contact</h2>

**Ramy Hadid** &mdash; [ramy.hadid@esst-sup.com](mailto:ramy.hadid@esst-sup.com)

Project Link: [https://github.com/rhpo/lms](https://github.com/rhpo/lms)

<br>

---

<div align="center">
  <sub>Built with Go + SvelteKit for ESST University.</sub>
</div>
