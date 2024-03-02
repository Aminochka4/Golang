# Social media as Questionnaires

This is simple app where everyone can create profile
with own questionnaires.

## API structure

## DB structure

![img.png](img.png)

```
Table user {
  id bigserial [primary key]
  createdAt timestamp
  updatedAt timestamp
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

Ref: questionnaire.userId < user.id
```


## Installation