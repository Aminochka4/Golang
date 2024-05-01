# Student
22B030305 Amanzholova Amina

# Social media as Questionnaires

This is a simple app where everyone can create profile with their own questionnaires.

## API structure

### Endpoints:

### Users
+ ```POST /api/v1/users/register:``` Register a new user.
+ ```POST /api/v1/users/activated:```To activate an account
+ ```POST /api/v1/users/login:```To login into an account
+ ```GET /api/v1/users:``` Get all users.
+ ```GET /api/v1/users/{userId}:``` Get a user by ID.
+ ```PUT /api/v1/users/{userId}:``` Update a user by ID.

### Questionnaires
+ ```POST /api/v1/questionnaires:``` Create a new questionnaire.
+ ```GET /api/v1/questionnaires:``` Get all questionnaires.
+ ```GET /api/v1/questionnaires/{questionnairesId}:``` Get a questionnaire by ID.
+ ```PUT /api/v1/questionnaires/{questionnairesId}:``` Update a questionnaire by ID.
+ ```DELETE /api/v1/questionnaires/{questionnairesId}:``` Delete a questionnaire by ID.

### Answer
+ ```GET /api/v1/answer:``` Get all answers
+ ```GET /api/v1/answer/{answerId}:``` Get an answer by ID
+ ```PUT /api/v1/answer/{answerId}:``` Update an answer by ID
+ ```DELETE /api/v1/answer/{answerId}:``` Delete an answer by ID

## DB structure

![img.png](img.png)

```
Table user {
  id bigserial [primary key]
  createdAt timestamp
  name text
  surname text
  username text
  email text
  password text
}

Table questionnaire {
  id bigserial [pk]
  createdAt timestamp
  updatedAt timestamp
  topic text
  questions text
  userId bigserial
}

Table answer{
  id biserial [pk]
  createdAt timestamp
  updatedAt timestamp
  questionnaireId bigint
  answer text
  userId bigserial
}

Ref: questionnaire.userId < user.id
```


## Installation