# Task Management API

## Introduction

This project is a RESTful API for a Task Management microservice built in Go. It provides endpoints for managing users, tasks, and projects, including operations such as creating, updating, deleting, and searching.

## Installation

1. Clone the repository:

```sh
git clone https://github.com/togzhanzhakhani/projects.git
cd projects
```

2. Build the project:

```sh
make build
```

3. Run the server:

```sh
make run
```

## Deployed on RENDER:

Base URL: https://projects-j02i.onrender.com

## API Endpoints
## Users
### URL: /users
#### GET /users: Get a list of all users.
#### POST /users: Create a new user.
### Request Body:

```sh
{
    "name": "John Doe",
    "email": "johndoe@example.com"
    "role": "admin"
}
```
#### GET /users/{id}: Get details of a specific user.
#### PUT /users/{id}: Update details of a specific user.
### Request Body:

```sh
{
    "name": "John Doe",
    "email": "johndoe@example.com"
    "role": "manager"
}
```

#### DELETE /users/{id}: Delete a specific user.
#### GET /users/{id}/tasks: Get a list of tasks for a specific user.
#### GET /users/search?name={name}: Find users by name.
#### GET /users/search?email={email}: Find users by email.

## Projects
### URL: /projects
#### GET /projects: Get a list of all projects.
#### POST /projects: Create a new project.
### Request Body:

```sh
{
    "name": "Project Alpha",
    "description": "A new innovative project",
    "start_date": "2024-07-01",
    "end_date": "2024-12-31",
    "manager_id": 1
}
```
#### GET /projects/{id}: Get details of a specific project.
#### PUT /projects/{id}: Update details of a specific project.
### Request Body:

```sh
{
    "name": "Project Beta",
    "description": "An updated innovative project",
    "start_date": "2024-07-01",
    "end_date": "2024-12-31",
    "manager_id": 1
}
```
#### DELETE /projects/{id}: Delete a specific project.
#### GET /projects/{id}/tasks: Get a list of tasks in a specific project.
#### GET /projects/search?title={title}: Find projects by title.
#### GET /projects/search?manager={userId}: Find projects by manager ID.

## Tasks
### URL: /tasks
#### GET /tasks: Get a list of all tasks.
#### POST /tasks: Create a new task.
### Request Body:

```sh
{
    "title": "Finish Report",
    "description": "Complete the quarterly financial report",
    "priority": "High",
    "status": "InProgress",
    "assignee_id": 3,
    "project_id": 5,
    "created_at": "2024-07-01",
    "completed_at": "2024-07-15"
}
```
#### GET /tasks/{id}: Get details of a specific task.
#### PUT /tasks/{id}: Update details of a specific task.
### Request Body:

```sh
{
    "title": "Finish Report",
    "description": "Complete the quarterly financial report and review",
    "priority": "High",
    "status": "Completed",
    "assignee_id": 3,
    "project_id": 5,
    "created_at": "2024-07-01",
    "completed_at": "2024-07-10"
}
```
#### DELETE /tasks/{id}: Delete a specific task.
#### GET /tasks/search?title={title}: Find tasks by title. specific project.
#### GET /tasks/search?status={status}: Find tasks by status.
#### GET /tasks/search?priority={priority}: Find tasks by priority.
#### GET /tasks/search?assignee={userId}: Find tasks by assignee ID.
#### GET /tasks/search?project={projectId}: Find tasks by project ID.