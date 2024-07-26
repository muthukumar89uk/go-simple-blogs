package router

import (
	//User-defined packages
	"blog/handlers"
	"blog/logs"
	"blog/middleware"

	//Third-party package
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func Router(Db *gorm.DB) {
	log := logs.Log()
	control := handlers.Database{Db: Db}
	app := fiber.New()

	//Public
	app.Post("/signup", control.Signup)
	app.Post("/login", control.Login)
	app.Get("/getAllPosters", middleware.AuthMiddleware(control.Db), control.GetAllPosters)
	app.Get("/getPoster/:post_id", middleware.AuthMiddleware(control.Db), control.GetPosterById)
	app.Get("/getComments/:post_id", middleware.AuthMiddleware(control.Db), control.GetCommentByPostId)
	app.Delete("/deleteComment/:comment_id", middleware.AuthMiddleware(control.Db), control.DeleteCommentById)
	app.Get("/logout", middleware.AuthMiddleware(control.Db), control.Logout)

	//Only for user
	app.Post("/user/addComment", control.AddComment)
	app.Put("/user/editComment/:comment_id", middleware.AuthMiddleware(control.Db), control.EditCommentByCommentId)

	//Only for admin
	app.Post("/admin/postPoster", middleware.AuthMiddleware(control.Db), control.PostPoster)
	app.Get("/admin/getPostersByAdmin", middleware.AuthMiddleware(control.Db), control.GetPostersByUserId)
	app.Put("/admin/updatePoster/:post_id", middleware.AuthMiddleware(control.Db), control.UpdatePosterById)
	app.Delete("/admin/deletePoster/:post_id", middleware.AuthMiddleware(control.Db), control.DeletePosterById)

	//start a server
	log.Info.Println("Message : 'Server starts in port 8000.....'")
	app.Listen(":8000")
}
