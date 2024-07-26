# Blog Project
This is a simple blog project implemented in Go (Golang) using the Fiber web framework and GORM as the ORM library for interacting with the database. The project allows users to sign up, log in, post articles, add comments to articles, view articles, and manage user sessions using JWT (JSON Web Tokens) for authentication.

Clone the git repo - `git clone https://github.com/muthukumar89uk/gin-RESTAPI-postgres-gorm.git` - or [download it](https://github.com/muthukumar89uk/gin-RESTAPI-postgres-gorm/zipball/master).

## Technologies used
The project is built using the following technologies:
- **Golang**  : The backend is written in Go (Golang), a statically typed, compiled language.
- **Fiber**   : The Fiber web framework is used to create RESTful APIs and handle HTTP requests.
- **JWT**     : JSON Web Tokens are used for secure user authentication and authorization.
- **bcrypt**  : Passwords are stored securely in hashed form using the bcrypt hashing algorithm.
- **Postgres**: Here, users data and Post articles data are handled in Postgres SQL.

## Project Structure
The project is organized into several packages, each responsible for specific functionalities:
- `handlers`  : Contains the HTTP request handlers for different API endpoints.
- `logs`      : Custom package for logging.
- `middleware`: Custom middleware for handling authentication and authorization.
- `models`    : Defines the data models used in the application.
- `repository`: Contains functions for interacting with the database.
- `drivers`   : Contains functions for establish a connection to database.
- `helper`    : Custom package that contains all the constants.

## Project Explanations
**Signup**: Allows users to sign up by providing their signup credentials. Passwords are securely hashed using bcrypt before storing them in the database.

**Login**: Users can log in using their registered email and password. Successful login generates a JWT token that is sent back to the client for further authenticated requests.

**Post Article**: Authenticated users can create and post articles. The articles are associated with the logged-in user and a category.

**Get Articles**: Users, both authenticated and unauthenticated, can view all articles available in the system.

**Get Articles by User**: Authenticated users can retrieve articles posted by a particular user.

**Get Article by ID**: Allows retrieval of a specific article based on its ID.

**Update Article**: Authenticated users can update their own articles by providing new content.

**Delete Article**: Authenticated users can delete their own articles.

**Add Comment**: Authenticated users can add comments to articles.

**Get Comments by Article**: Allows retrieval of all comments associated with a specific article.

**Edit Comment**: Authenticated users can edit their own comments.

**Delete Comment**: Authenticated users can delete their own comments.

**Logout**: Authenticated users can log out, and their JWT tokens are invalidated.

## Note
- The project uses the Fiber web framework for handling HTTP requests.[here](https://github.com/gofiber/fiber).