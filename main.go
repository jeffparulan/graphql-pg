package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/graphql-go/graphql"
	"github.com/graphql-go/handler"
	_ "github.com/lib/pq"

	"github.com/jinzhu/gorm"
)

const (
	DB_HOST     = "localhost"
	DB_PORT     = "5432"
	DB_USER     = "postgres"
	DB_PASSWORD = "xxx"
	DB_NAME     = "postgres"
)

type Patient struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at`
}

type Post struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	PatientID int       `json:"patient_id"`
	CreatedAt time.Time `json:"created_at"`
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func ConnectDB() (*gorm.DB, error) {
	databaseURL := os.Getenv("postgresql://localhost:5432")
	databaseName := os.Getenv("postgres")

	db, err := gorm.Open("postgres", fmt.Sprintf("%v/%v", databaseURL, databaseName))
	if err != nil {
		log.Printf("Error connecting database.\n%v", err)
		return nil, err
	}

	return db, nil
}

func main() {

	dbinfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		DB_HOST, DB_PORT, DB_USER, DB_PASSWORD, DB_NAME)
	db, err := sql.Open("postgres", dbinfo)
	checkErr(err)

	defer db.Close()

	PatientType := graphql.NewObject(graphql.ObjectConfig{
		Name:        "Patient",
		Description: "A Patient who wrote the post",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type:        graphql.NewNonNull(graphql.Int),
				Description: "The identifier of the Patient.",
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					if Patient, ok := p.Source.(*Patient); ok {
						return Patient.ID, nil
					}

					return nil, nil
				},
			},
			"name": &graphql.Field{
				Type:        graphql.NewNonNull(graphql.String),
				Description: "The name of the Patient.",
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					if Patient, ok := p.Source.(*Patient); ok {
						return Patient.Name, nil
					}

					return nil, nil
				},
			},
			"email": &graphql.Field{
				Type:        graphql.NewNonNull(graphql.String),
				Description: "The email address of the Patient.",
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					if Patient, ok := p.Source.(*Patient); ok {
						return Patient.Email, nil
					}

					return nil, nil
				},
			},
			"created_at": &graphql.Field{
				Type:        graphql.NewNonNull(graphql.String),
				Description: "The created_at date of the Patient.",
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					if Patient, ok := p.Source.(*Patient); ok {
						return Patient.CreatedAt, nil
					}

					return nil, nil
				},
			},
		},
	})

	postType := graphql.NewObject(graphql.ObjectConfig{
		Name:        "Post",
		Description: "A Post made by a registered Patient",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type:        graphql.NewNonNull(graphql.Int),
				Description: "The identifier of the post.",
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					if post, ok := p.Source.(*Post); ok {
						return post.ID, nil
					}

					return nil, nil
				},
			},
			"title": &graphql.Field{
				Type:        graphql.NewNonNull(graphql.String),
				Description: "The title of the post.",
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					if post, ok := p.Source.(*Post); ok {
						return post.Title, nil
					}

					return nil, nil
				},
			},
			"content": &graphql.Field{
				Type:        graphql.NewNonNull(graphql.String),
				Description: "The content of the post.",
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					if post, ok := p.Source.(*Post); ok {
						return post.Content, nil
					}

					return nil, nil
				},
			},
			"created_at": &graphql.Field{
				Type:        graphql.NewNonNull(graphql.String),
				Description: "The created_at date of the post.",
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					if post, ok := p.Source.(*Post); ok {
						return post.CreatedAt, nil
					}

					return nil, nil
				},
			},
			"Patient": &graphql.Field{
				Type: PatientType,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					if post, ok := p.Source.(*Post); ok {
						Patient := &Patient{}
						err = db.QueryRow("select id, name, email from patient where id = $1", post.PatientID).Scan(&Patient.ID, &Patient.Name, &Patient.Email)
						checkErr(err)

						return Patient, nil
					}

					return nil, nil
				},
			},
		},
	})

	rootQuery := graphql.NewObject(graphql.ObjectConfig{
		Name: "RootQuery",
		Fields: graphql.Fields{
			"Patient": &graphql.Field{
				Type:        PatientType,
				Description: "Get a Patient.",
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{
						Type: graphql.Int,
					},
				},
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					id, _ := params.Args["id"].(int)

					Patient := &Patient{}
					err = db.QueryRow("select id, name, email from patient where id = $1", id).Scan(&Patient.ID, &Patient.Name, &Patient.Email)
					checkErr(err)

					return Patient, nil
				},
			},
			"Patients": &graphql.Field{
				Type:        graphql.NewList(PatientType),
				Description: "List of Patients.",
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					rows, err := db.Query("SELECT id, name, email FROM patient")
					checkErr(err)
					var Patients []*Patient

					for rows.Next() {
						Patient := &Patient{}

						err = rows.Scan(&Patient.ID, &Patient.Name, &Patient.Email)
						checkErr(err)
						Patients = append(Patients, Patient)
					}

					return Patients, nil
				},
			},
			"post": &graphql.Field{
				Type:        postType,
				Description: "Get a patient's post.",
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{
						Type: graphql.Int,
					},
				},
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					id, _ := params.Args["id"].(int)

					post := &Post{}
					err = db.QueryRow("select id, title, content, patient_id from posts where id = $1", id).Scan(&post.ID, &post.Title, &post.Content, &post.PatientID)
					checkErr(err)

					return post, nil
				},
			},
			"posts": &graphql.Field{
				Type:        graphql.NewList(postType),
				Description: "List of posts.",
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					rows, err := db.Query("SELECT id, title, content, patient_id FROM posts")
					checkErr(err)
					var posts []*Post

					for rows.Next() {
						post := &Post{}

						err = rows.Scan(&post.ID, &post.Title, &post.Content, &post.PatientID)
						checkErr(err)
						posts = append(posts, post)
					}

					return posts, nil
				},
			},
		},
	})

	rootMutation := graphql.NewObject(graphql.ObjectConfig{
		Name: "RootMutation",
		Fields: graphql.Fields{
			// Patient
			"createPatient": &graphql.Field{
				Type:        PatientType,
				Description: "Create new Patient",
				Args: graphql.FieldConfigArgument{
					"name": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
					"email": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
				},
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					name, _ := params.Args["name"].(string)
					email, _ := params.Args["email"].(string)
					createdAt := time.Now()

					var lastInsertId int
					err = db.QueryRow("INSERT INTO patient(name, email, created_at) VALUES($1, $2, $3) returning id;", name, email, createdAt).Scan(&lastInsertId)
					checkErr(err)

					newPatient := &Patient{
						ID:        lastInsertId,
						Name:      name,
						Email:     email,
						CreatedAt: createdAt,
					}

					return newPatient, nil
				},
			},
			"updatePatient": &graphql.Field{
				Type:        PatientType,
				Description: "Update an Patient",
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.Int),
					},
					"name": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
					"email": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
				},
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					id, _ := params.Args["id"].(int)
					name, _ := params.Args["name"].(string)
					email, _ := params.Args["email"].(string)

					stmt, err := db.Prepare("UPDATE patient SET name = $1, email = $2 WHERE id = $3")
					checkErr(err)

					_, err2 := stmt.Exec(name, email, id)
					checkErr(err2)

					newPatient := &Patient{
						ID:    id,
						Name:  name,
						Email: email,
					}

					return newPatient, nil
				},
			},
			"deletePatient": &graphql.Field{
				Type:        PatientType,
				Description: "Delete an Patient",
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{
						Type: graphql.Int,
					},
				},
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					id, _ := params.Args["id"].(int)

					stmt, err := db.Prepare("DELETE FROM patient WHERE id = $1")
					checkErr(err)

					_, err2 := stmt.Exec(id)
					checkErr(err2)

					return nil, nil
				},
			},
			// Post
			"createPost": &graphql.Field{
				Type:        postType,
				Description: "Create new post",
				Args: graphql.FieldConfigArgument{
					"title": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
					"content": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
					"Patient_id": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.Int),
					},
				},
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					title, _ := params.Args["title"].(string)
					content, _ := params.Args["content"].(string)
					PatientId, _ := params.Args["Patient_id"].(int)
					createdAt := time.Now()

					var lastInsertId int
					err = db.QueryRow("INSERT INTO posts(title, content, patient_id, created_at) VALUES($1, $2, $3, $4) returning id;", title, content, PatientId, createdAt).Scan(&lastInsertId)
					checkErr(err)

					newPost := &Post{
						ID:        lastInsertId,
						Title:     title,
						Content:   content,
						PatientID: PatientId,
						CreatedAt: createdAt,
					}

					return newPost, nil
				},
			},
			"updatePost": &graphql.Field{
				Type:        postType,
				Description: "Update a post",
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.Int),
					},
					"title": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
					"content": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
					"Patient_id": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.Int),
					},
				},
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					id, _ := params.Args["id"].(int)
					title, _ := params.Args["title"].(string)
					content, _ := params.Args["content"].(string)
					PatientId, _ := params.Args["patient_id"].(int)

					stmt, err := db.Prepare("UPDATE posts SET title = $1, content = $2, Patient_id = $3 WHERE id = $4")
					checkErr(err)

					_, err2 := stmt.Exec(title, content, PatientId, id)
					checkErr(err2)

					newPost := &Post{
						ID:        id,
						Title:     title,
						Content:   content,
						PatientID: PatientId,
					}

					return newPost, nil
				},
			},
			"deletePost": &graphql.Field{
				Type:        postType,
				Description: "Delete a post",
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{
						Type: graphql.Int,
					},
				},
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					id, _ := params.Args["id"].(int)

					stmt, err := db.Prepare("DELETE FROM posts WHERE id = $1")
					checkErr(err)

					_, err2 := stmt.Exec(id)
					checkErr(err2)

					return nil, nil
				},
			},
		},
	})

	schema, _ := graphql.NewSchema(graphql.SchemaConfig{
		Query:    rootQuery,
		Mutation: rootMutation,
	})

	h := handler.New(&handler.Config{
		Schema:   &schema,
		Pretty:   true,
		GraphiQL: true,
	})

	// serve HTTP
	http.Handle("/graphql", h)
	http.ListenAndServe(":8080", nil)
}
