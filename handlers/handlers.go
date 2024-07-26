package handlers

import (

	//User-defined packages
	"blog/logs"
	"blog/middleware"
	"blog/models"
	"blog/repository"

	//Inbuild packages
	"fmt"
	"reflect"
	"regexp"
	"strconv"

	//Third-party packages
	"github.com/fatih/structs"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type Database struct {
	Db *gorm.DB
}

// This is for Signup
func (db Database) Signup(c *fiber.Ctx) error {
	var (
		data models.User
		role models.Roles
	)
	log := logs.Log()
	log.Info.Println("Message : 'signup-API called'")

	//Get user details from request body
	if err := c.BodyParser(&data); err != nil {
		log.Error.Println("Error : 'internal server error' Status : 500")
		return c.JSON(fiber.Map{
			"status": 500,
			"error":  "internal server error",
		})
	}

	//To check if any credential is missing or not
	fields := structs.Names(&models.SignupReq{})
	for _, field := range fields {
		if reflect.ValueOf(&data).Elem().FieldByName(field).Interface() == "" {
			stmt := fmt.Sprintf("missing %s", field)
			log.Error.Printf("Error : '%s' Status : 400\n", stmt)
			return c.JSON(fiber.Map{
				"status": 400,
				"error":  stmt,
			})
		}
	}

	//validate email format
	emailRegex := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	if !emailRegex.MatchString(data.Email) {
		log.Error.Println("Error : 'Invalid Email' Status : 400")
		return c.JSON(fiber.Map{
			"status": 400,
			"error":  "Invalid Email",
		})
	}

	//validate the password
	if len(data.Password) < 8 {
		log.Error.Println("Error : 'password must be greater than 8 characters' Status : 400")
		return c.JSON(fiber.Map{
			"status": 400,
			"error":  "password must be greater than 8 characters",
		})
	}

	//validate the role
	if data.Role != "admin" && data.Role != "user" {
		log.Error.Println("Error : 'Invalid role' Status : 400")
		return c.JSON(fiber.Map{
			"status": 400,
			"error":  "Invalid role",
		})
	}

	//To check if the user details already exist or not
	data, err := repository.ReadUserByEmail(db.Db, data)
	if err == nil {
		log.Error.Println("Error : 'user already exist' Status : 409")
		return c.JSON(fiber.Map{
			"status": 409,
			"error":  "user already exist",
		})
	}

	//To change the password into hashedPassword
	password, err := bcrypt.GenerateFromPassword([]byte(data.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Error.Printf("Error : '%s'\n", err)
		return nil
	}
	data.Password = string(password)

	//Select a role_id for specified role
	role, _ = repository.ReadRoleIdByRole(db.Db, data)
	data.RoleId = role.RoleId

	//Adding a user details into our database
	if err = repository.CreateUser(db.Db, data); err != nil {
		log.Error.Println("Error : 'email already exist' Status : 409")
		return c.JSON(fiber.Map{
			"status": 409,
			"error":  "email already exist",
		})
	}

	log.Info.Println("Message : 'signup successful!!!' Status : 200")
	return c.JSON(fiber.Map{
		"status":    200,
		"message":   "signup successful!!!",
		"user data": data,
	})
}

// This is for Login
func (db Database) Login(c *fiber.Ctx) error {
	var data models.User
	log := logs.Log()
	log.Info.Println("Message : 'login-API called'")
	//Get mail-id and password from request body
	if err := c.BodyParser(&data); err != nil {
		log.Error.Println("Error : 'internal server error' Status : 500")
		return c.JSON(fiber.Map{
			"status": 500,
			"error":  "internal server error",
		})
	}

	//To check if any credential is missing or not
	fields := structs.Names(&models.LoginReq{})
	for _, field := range fields {
		if reflect.ValueOf(&data).Elem().FieldByName(field).Interface() == "" {
			stmt := fmt.Sprintf("missing %s", field)
			log.Error.Printf("Error : '%s' Status : 400\n", stmt)
			return c.JSON(fiber.Map{
				"Status": 400,
				"error":  stmt,
			})
		}
	}

	//validates correct email format
	emailRegex := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	if !emailRegex.MatchString(data.Email) {
		log.Error.Println("Error : 'Invalid Email' Status : 400")
		return c.JSON(fiber.Map{
			"status": 400,
			"error":  "Invalid Email",
		})
	}

	//To verify if the user email is exist or not
	user, err := repository.ReadUserByEmail(db.Db, data)
	if err == nil {
		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(data.Password)); err == nil {
			// Fetch a JWT token
			auth, err := repository.ReadTokenByUserId(db.Db, user)
			if err == nil {
				log.Info.Println("Message : 'login successful!!!' Status : 200")
				return c.JSON(fiber.Map{
					"status":  200,
					"message": "Login Successful!!!",
					"token":   auth.Token,
				})
			}

			//Create a token
			token, err := middleware.CreateToken(user, c)
			if err != nil {
				return err
			}
			auth.UserId, auth.Token = user.UserId, token
			if err = repository.AddToken(db.Db, auth); err != nil {
				log.Error.Printf("Error : '%s' Status : 409\n", err)
				return c.JSON(fiber.Map{
					"status": 409,
					"error":  err.Error(),
				})
			}

			log.Info.Println("Message : 'login successful!!!' Status : 200")
			return c.JSON(fiber.Map{
				"status":  200,
				"message": "Login Successful!!!",
				"token":   token,
			})
		}
		log.Error.Println("Error : 'incorrect password' Status : 400")
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"status": 400,
			"error":  "incorrect password",
		})

	}
	log.Error.Println("Error : 'user not found' Status : 404")
	c.Status(fiber.StatusNotFound)
	return c.JSON(fiber.Map{
		"status": 404,
		"error":  "user not found",
	})
}

// Handler for post a poster
func (db Database) PostPoster(c *fiber.Ctx) error {
	var (
		Post     models.Post
		Catagory models.Catagory
	)
	log := logs.Log()
	if err := middleware.AdminAuth(c); err != nil {
		log.Error.Println("Error : 'unauthorized entry' Status : 401")
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error":  "unauthorized entry",
			"status": 401,
		})
	}
	log.Info.Println("Message : 'poster-API called'")
	if err := c.BodyParser(&Post); err != nil {
		log.Error.Println("Error : 'internal server error' Status : 500")
		return c.JSON(fiber.Map{
			"status": 500,
			"error":  "internal server error",
		})
	}

	//To check if any credential is missing or not
	fields := structs.Names(&models.PostReq{})
	for _, field := range fields {
		if reflect.ValueOf(&Post).Elem().FieldByName(field).Interface() == "" {
			stmt := fmt.Sprintf("missing %s", field)
			log.Error.Printf("Error : '%s' Status : 400\n", stmt)
			return c.JSON(fiber.Map{
				"Status": 400,
				"error":  stmt,
			})
		}
	}
	claims := middleware.GetTokenClaims(c)
	userId, _ := strconv.Atoi(claims.Id)
	Post.UserId = uint(userId)
	Catagory, err := repository.ReadCatagoryIdByCatagory(db.Db, Post)
	if err != nil {
		log.Error.Println("Error : 'Invalid catagory' Status : 400")
		return c.JSON(fiber.Map{
			"status": 400,
			"error":  "Invalid catagory",
		})
	}
	Post.CatagoryId = Catagory.CatagoryId
	if err = repository.CreatePost(db.Db, Post); err != nil {
		log.Error.Printf("Error : '%s' Status : 400\n", err)
		return c.JSON(fiber.Map{
			"status": 400,
			"error":  err.Error(),
		})
	}
	log.Info.Println("Message : 'Post added successfully' Status : 200")
	return c.JSON(fiber.Map{
		"status":  200,
		"message": "Post added successfully",
	})
}

// Handler for get posters which were posted by a particular user
func (db Database) GetPostersByUserId(c *fiber.Ctx) error {
	log := logs.Log()
	if err := middleware.AdminAuth(c); err != nil {
		log.Error.Println("Error : 'unauthorized entry' Status : 401")
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error":  "unauthorized entry",
			"status": 401,
		})
	}
	log.Info.Println("Message : 'Getposters-API called'")
	claims := middleware.GetTokenClaims(c)
	Posts, err := repository.ReadPostersByUserId(db.Db, claims.Id)
	if err == nil {
		log.Info.Println("Message : 'Post retrieved successfully' Status : 200")
		return c.JSON(fiber.Map{
			"status": 200,
			"Posts":  Posts,
		})
	}
	log.Error.Println("Error : 'Post not found' Status : 404")
	c.Status(fiber.StatusNotFound)
	return c.JSON(fiber.Map{
		"status": 404,
		"error":  "Post not found",
	})
}

// Handler for get all posters
func (db Database) GetAllPosters(c *fiber.Ctx) error {
	log := logs.Log()
	log.Info.Println("Message : 'GetAllPosters-API called'")
	Posts, err := repository.ReadAllPosters(db.Db)
	if err == nil {
		log.Info.Println("Message : 'Post retrieved successfully' Status : 200")
		return c.JSON(fiber.Map{
			"status": 200,
			"Posts":  Posts,
		})
	}
	log.Error.Println("Error : 'Post not found' Status : 404")
	c.Status(fiber.StatusNotFound)
	return c.JSON(fiber.Map{
		"status": 404,
		"error":  "Post not found",
	})
}

// Handler for get a poster by post-id
func (db Database) GetPosterById(c *fiber.Ctx) error {
	log := logs.Log()
	log.Info.Println("Message : 'Getposter-API called'")
	Post, err := repository.ReadPostByPostId(db.Db, c.Params("post_id", ""))
	if err == nil {
		log.Info.Println("Message : 'Post retrieved successfully' Status : 200")
		return c.JSON(fiber.Map{
			"status":    200,
			"Post data": Post,
		})
	}
	log.Error.Println("Error : 'Post not found' Status : 404")
	c.Status(fiber.StatusNotFound)
	return c.JSON(fiber.Map{
		"status": 404,
		"error":  "Post not found",
	})
}

// Handler for update a poster by post-id
func (db Database) UpdatePosterById(c *fiber.Ctx) error {
	var check int
	log := logs.Log()

	if err := middleware.AdminAuth(c); err != nil {
		log.Error.Println("Error : 'unauthorized entry' Status : 401")
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error":  "unauthorized entry",
			"status": 401,
		})
	}
	log.Info.Println("Message : 'Updateposter-API called'")
	Post, err := repository.ReadPostByPostId(db.Db, c.Params("post_id", ""))
	if err == nil {
		if err := c.BodyParser(&Post); err != nil {
			log.Error.Println("Error : 'internal server error' Status : 500")
			return c.JSON(fiber.Map{
				"status": 500,
				"error":  "internal server error",
			})
		}

		fields := structs.Names(models.PostReq{})
		for _, field := range fields {
			if reflect.ValueOf(&Post).Elem().FieldByName(field).Interface() == "" {
				check++
			}
		}
		if check == 3 {
			log.Error.Println("Error : 'no data found to do update' Status : 404")
			return c.JSON(fiber.Map{
				"status": 404,
				"error":  "no data found to do update",
			})
		}
		if Post.Catagory != "" {
			Catagory, err := repository.ReadCatagoryIdByCatagory(db.Db, Post)
			if err != nil {
				log.Error.Println("Error : 'invalid catagory' Status : 404")
				c.Status(fiber.StatusNotFound)
				return c.JSON(fiber.Map{
					"status": 404,
					"error":  "invalid catagory",
				})
			}
			Post.CatagoryId = Catagory.CatagoryId
		}
		if err := repository.UpdatePostByPostId(db.Db, c.Params("post_id", ""), Post); err == nil {
			log.Info.Println("Message : 'Post updated successfully' Status : 200")
			return c.JSON(fiber.Map{
				"status":  200,
				"message": "post updated Successfully!!!",
			})
		}
	}
	log.Error.Println("Error : 'Post not found' Status : 404")
	c.Status(fiber.StatusNotFound)
	return c.JSON(fiber.Map{
		"status": 404,
		"error":  "Post not found",
	})
}

// Handler for delete a poster by post-id
func (db Database) DeletePosterById(c *fiber.Ctx) error {
	log := logs.Log()
	if err := middleware.AdminAuth(c); err != nil {
		log.Error.Println("Error : 'unauthorized entry' Status : 401")
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error":  "unauthorized entry",
			"status": 401,
		})
	}
	log.Info.Println("Message : 'Deleteposter-API called'")
	if _, err := repository.ReadPostByPostId(db.Db, c.Params("post_id", "")); err == nil {
		repository.DeletePostByPostId(db.Db, c.Params("post_id", ""))
		log.Info.Println("Message : 'Post deleted successfully' Status : 200")
		return c.JSON(fiber.Map{
			"status":  200,
			"message": "post deleted Successfully!!!",
		})
	}

	log.Error.Println("Error : 'Post not found' Status : 404")
	c.Status(fiber.StatusNotFound)
	return c.JSON(fiber.Map{
		"status": 404,
		"error":  "Post not found",
	})
}

// Handler for add comment to a post
func (db Database) AddComment(c *fiber.Ctx) error {
	var commentData models.Comments
	log := logs.Log()
	if err := middleware.UserAuth(c); err != nil {
		log.Error.Println("Error : 'unauthorized entry' Status : 401")
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error":  "unauthorized entry",
			"status": 401,
		})

	}
	log.Info.Println("Message : 'AddComment-API called'")
	if err := c.BodyParser(&commentData); err != nil {
		log.Error.Println("Error : 'internal server error' Status : 500")
		return c.JSON(fiber.Map{
			"status": 500,
			"error":  "internal server error",
		})
	}
	fields := structs.Names(&models.CommentReq{})
	for _, field := range fields {
		if reflect.ValueOf(&commentData).Elem().FieldByName(field).Interface() == "" {
			stmt := fmt.Sprintf("missing %s", field)
			log.Error.Printf("Error : '%s' Status : 400\n", stmt)
			return c.JSON(fiber.Map{
				"status": 400,
				"error":  stmt,
			})
		}
	}
	claims := middleware.GetTokenClaims(c)
	userId, _ := strconv.Atoi(claims.Id)
	commentData.UserId = uint(userId)
	Post, err := repository.ReadPostIdbyPostTitle(db.Db, commentData.PostTitle)
	if err != nil {
		log.Error.Println("Error : 'Post not found' Status : 404")
		return c.JSON(fiber.Map{
			"status": 404,
			"error":  "post not found",
		})
	}
	commentData.PostId = Post.PostId
	if err := repository.CreateComment(db.Db, commentData); err != nil {
		log.Error.Printf("Error : '%s' Status : 400\n", err)
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"status": 400,
			"error":  err.Error(),
		})
	}
	log.Info.Println("Message : 'Comment added successfully' Status : 200")
	c.Status(fiber.StatusAccepted)
	return c.JSON(fiber.Map{
		"status":  200,
		"message": "comment added successfully",
	})
}

// Handler for Edit a comment by comment-id
func (db Database) EditCommentByCommentId(c *fiber.Ctx) error {
	log := logs.Log()
	if err := middleware.UserAuth(c); err != nil {
		log.Error.Println("Error : 'unauthorized entry' Status : 401")
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error":  "unauthorized entry",
			"status": 401,
		})
	}
	log.Info.Println("Message : 'Updateposter-API called'")
	comment, err := repository.ReadCommentByCommentId(db.Db, c.Params("comment_id", ""))
	if err == nil {
		if err := c.BodyParser(&comment); err != nil {
			log.Error.Println("Error : 'internal server error' Status : 500")
			return c.JSON(fiber.Map{
				"status": 500,
				"error":  "internal server error",
			})
		}

		if comment.Comment == "" {
			log.Error.Println("Error : 'comment field is required' Status : 404")
			return c.JSON(fiber.Map{
				"status": 404,
				"error":  "comment field is required",
			})
		}

		if err := repository.EditCommentByCommentId(db.Db, c.Params("comment_id", ""), comment); err == nil {
			log.Info.Println("Message : 'comment edited  successfully' Status : 200")
			return c.JSON(fiber.Map{
				"status":  200,
				"message": "comment edited Successfully!!!",
			})
		}
	}
	log.Error.Println("Error : 'comment not found' Status : 404")
	c.Status(fiber.StatusNotFound)
	return c.JSON(fiber.Map{
		"status": 404,
		"error":  "comment not found",
	})
}

// Handler for get a comment by post-id
func (db Database) GetCommentByPostId(c *fiber.Ctx) error {
	log := logs.Log()
	log.Info.Println("Message : 'GetCommentById-API called'")
	commentData, err := repository.ReadCommentsByPostId(db.Db, c.Params("post_id", ""))
	if err == nil && commentData != nil {
		log.Info.Println("Message : 'comment(s) retrieved successfully' Status : 200")
		return c.JSON(fiber.Map{
			"status":   200,
			"Comments": commentData,
		})
	}
	log.Error.Println("Error : 'Comment not found for this post' Status : 404")
	c.Status(fiber.StatusNotFound)
	return c.JSON(fiber.Map{
		"status": 404,
		"error":  "Comment not found for this post",
	})
}

// Handler for delete a comment
func (db Database) DeleteCommentById(c *fiber.Ctx) error {
	log := logs.Log()
	log.Info.Println("Message : 'DeleteCommentById-API called'")
	if _, err := repository.ReadCommentByCommentId(db.Db, c.Params("comment_id", "")); err == nil {
		repository.DeleteComment(db.Db, c.Params("comment_id", ""))
		log.Info.Println("Message : 'comment deleted successfully' Status : 200")
		return c.JSON(fiber.Map{
			"status":  200,
			"message": "Comment deleted Successfully!!!",
		})
	}
	log.Error.Println("Error : 'Comment not found' Status : 404")
	c.Status(fiber.StatusNotFound)
	return c.JSON(fiber.Map{
		"status": 404,
		"error":  "Comment not found",
	})
}

// Handler for Logout
func (db Database) Logout(c *fiber.Ctx) error {
	log := logs.Log()
	log.Info.Println("Message : 'Logout-API called'")
	claims := middleware.GetTokenClaims(c)
	repository.DeleteToken(db.Db, claims.Id)
	log.Info.Println("Message : 'Logout Successful' Status : 200")
	return c.JSON(fiber.Map{
		"status":  200,
		"message": "Logout Successful",
	})
}
