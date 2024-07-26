package repository

import (
	//User defined packages
	"blog/models"

	//Third-party packages
	"gorm.io/gorm"
)

// Adding a specified roles to the roles table
func AddRoles(Db *gorm.DB) {
	Roles := []models.Roles{
		{RoleId: 1, Role: "admin"},
		{RoleId: 2, Role: "user"},
	}
	Db.Create(&Roles)
}

// Adding a specified catagories to the catagories table
func Addcatagory(Db *gorm.DB) {
	Catagory := []models.Catagory{
		{CatagoryId: 1, Catagory: "Java"},
		{CatagoryId: 2, Catagory: "Python"},
		{CatagoryId: 3, Catagory: "Go"},
		{CatagoryId: 4, Catagory: "JavaScript"},
		{CatagoryId: 5, Catagory: "ColdFusion"},
	}
	Db.Create(&Catagory)
}

// Table creation
func TableCreation(Db *gorm.DB) {
	Db.AutoMigrate(&models.Roles{}, &models.Catagory{}, &models.User{}, &models.Authentication{}, &models.Post{}, &models.Comments{})
	AddRoles(Db)
	Addcatagory(Db)
}

// Retrieve the User's role by role-id
func ReadRoleIdByRole(Db *gorm.DB, data models.User) (models.Roles, error) {
	role := models.Roles{}
	err := Db.Select("role_id").Where("role=?", data.Role).First(&role).Error
	return role, err
}

// Adding a user details into users table
func CreateUser(Db *gorm.DB, data models.User) (err error) {
	err = Db.Create(&data).Error
	return
}

// Retrieve the User details by Email
func ReadUserByEmail(Db *gorm.DB, data models.User) (models.User, error) {
	err := Db.Where("email = ?", data.Email).First(&data).Error
	return data, err
}

// Retrieve a token by user-id
func ReadTokenByUserId(Db *gorm.DB, user models.User) (auth models.Authentication, err error) {
	err = Db.Where("user_id=?", user.UserId).First(&auth).Error
	return auth, err
}

// Adding a token into authorizations table
func AddToken(Db *gorm.DB, auth models.Authentication) error {
	err := Db.Create(&auth).Error
	return err
}

// Retrieve a catagory-id by post's catagory
func ReadCatagoryIdByCatagory(Db *gorm.DB, Post models.Post) (Catagory models.Catagory, err error) {
	err = Db.Select("catagory_id").Where("catagory=?", Post.Catagory).First(&Catagory).Error
	return
}

// Adding a post into posts table
func CreatePost(Db *gorm.DB, Post models.Post) error {
	err := Db.Create(&Post).Error
	return err
}

// Retrieve all posts which were posted by a user
func ReadPostersByUserId(Db *gorm.DB, userId string) (Posts []models.Post, err error) {
	err = Db.Where("user_id=?", userId).Find(&Posts).Error
	return
}

// Retrieve all posts from comment
func ReadAllPosters(Db *gorm.DB) (Posts []models.Post, err error) {
	err = Db.Find(&Posts).Error
	return
}

// Retrieve a post by post-id
func ReadPostByPostId(Db *gorm.DB, postId string) (Post models.Post, err error) {
	err = Db.Where("post_id=?", postId).First(&Post).Error
	return
}

// Update a post by post-id
func UpdatePostByPostId(Db *gorm.DB, postId string, Post models.Post) (err error) {
	err = Db.Where("post_id=?", postId).Save(&Post).Error
	return
}

// Delete a post by post-id
func DeletePostByPostId(Db *gorm.DB, postId string) (err error) {
	var Post models.Post
	err = Db.Where("post_id=?", postId).Delete(&Post).Error
	return
}

// Retrieve a post-id by post_title
func ReadPostIdbyPostTitle(Db *gorm.DB, postTitle string) (Post models.Post, err error) {
	err = Db.Where("post_title=?", postTitle).First(&Post).Error
	return
}

// Adding a comment into comments table
func CreateComment(Db *gorm.DB, comment models.Comments) (err error) {
	err = Db.Create(&comment).Error
	return
}

// Retrieve a comment by post-id
func ReadCommentsByPostId(Db *gorm.DB, postId string) (Comments []models.Comments, err error) {
	err = Db.Where("post_id=?", postId).Find(&Comments).Error
	return
}

// Retrieve a comment by comment-id
func ReadCommentByCommentId(Db *gorm.DB, commentId string) (Comment models.Comments, err error) {
	err = Db.Where("comment_id=?", commentId).First(&Comment).Error
	return
}

// Edit a comment by comment-id
func EditCommentByCommentId(Db *gorm.DB, commentId string, Comment models.Comments) (err error) {
	err = Db.Where("comment_id=?", commentId).Save(&Comment).Error
	return
}

// Delete a comment by comment-id
func DeleteComment(Db *gorm.DB, commentId string) (err error) {
	var comment models.Comments
	err = Db.Where("comment_id=?", commentId).Delete(&comment).Error
	return
}

// Delete a token by user-id
func DeleteToken(Db *gorm.DB, userId string) (err error) {
	var token models.Authentication
	err = Db.Where("user_id=?", userId).Delete(&token).Error
	return

}
