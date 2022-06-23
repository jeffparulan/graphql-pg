
# GraphQL with Golang and pg as backend RDBMS
> POC: GraphQL complete example using Golang & PostgreSQL

## Prerequisite (Installations)
Install the following dependencies:
```
go get github.com/graphql-go/graphql
go get github.com/graphql-go/handler
go get github.com/lib/pq
```

Install & create postgres database
```
brew install postgres
createuser graphql --createdb
createdb graphql -U graphql
psql graphql -U graphql
```
Note: Windows users, please install postgreSQL (https://www.postgresql.org/download/windows/)

Create the tables
```sql
CREATE TABLE IF NOT EXISTS patient
(
    id serial PRIMARY KEY,
    name varchar(100) NOT NULL,
    email varchar(150) NOT NULL,
    created_at date
);

CREATE TABLE IF NOT EXISTS posts
(
    id serial PRIMARY KEY,
    title varchar(100) NOT NULL,
    content text NOT NULL,
    patient_id int,
    created_at date
);
```

## How to run the api
go run main.go

Invoke http://localhost:8080/graphql

Follow the different types of query within the GraphQL Usage

 ## GraphQL Usage
 Query to get the all patients
```
query {
  Patients {
    id,
    name,
    email
  }
}
```

Query to get a specific patient
```
query {
 Patient(id: 2) {
  id,
  name,
  email
}
}
```

Create new patient using mutation
```
mutation {
  createPatient(name: "Kailo Ben", email: "kailo.ben@gmail.com") {
    id
    name
    email
  }
}
```

Update an patient using mutation
```
mutation {
  updatePatient(id: 2, name: "Sohel Amin Shah", email: "kailo.ben.1@gmail.com") {
    id
    name
    email
  }
}
```

Delete an patient using mutation
```
mutation {
  deletePatient(id: 2) {
    id
  }
}
```

Query to get the posts with its relation patient
```
query {
  posts {
    id
    title
    content
    Patient {
      id
      name
      email
    }
  }
}
```