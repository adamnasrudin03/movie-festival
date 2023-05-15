# Movie Festival App
## Overview
Your company wanted to create an app for a movie festival and you has been assigned for 
backend developer for the app.
Before the festival, an admin will manually collect movie files from invited participants. The
admin will then upload the movie files and add relevant movie attributes (title, etc) through
CMS, and the admin is also able to update the entry if necessary.
During the festival, everyone who installed the app can search and view the uploaded
movies. Every time someone sees a movie, it’s counted as 1 view for the movie.
Some users can register and login to the system. These authenticated users can vote for a
movie that they like, with limitation that 1 user can only vote for the same movie once.
However, an authenticated user can vote for multiple movies that they like, and they also can
unvote a movie if they changed their mind later.
After the festival finishes, the admin can see which movies are the most popular (have the
most viewership) and which movies are the most liked (have the most votes).


##  How to run this project

1. clone this project, create .env file according to .env.example file, and create db
2. run the project

```
go run main.go
```

##  Testing api with postman

```
import the file in Postman located at ./movie-festival.postman_collection.json
```


##  Tools used this project

1. Gin-gonic
2. Gorm
3. PosgreSQL
4. JWT
5. godotenv


##  Other information
default user admin account :

```
{
    "email": "admin@gmail.com",
    "password": "password123"
}
```
