package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/pascaldekloe/jwt"
	age "github.com/theTardigrade/golang-age"

	"otus-homework/internal/database"
	"otus-homework/internal/helpers"
	"otus-homework/internal/models"
	"otus-homework/internal/password"
	"otus-homework/internal/request"
	"otus-homework/internal/response"
	"otus-homework/internal/validator"
	"otus-homework/internal/worker"
)

func (app *application) status(w http.ResponseWriter, r *http.Request) {
	data := map[string]string{
		"Status": "OK",
	}

	err := response.JSON(w, http.StatusOK, data)
	if err != nil {
		app.serverError(w, r, err)
	}
}

func (app *application) createUser(w http.ResponseWriter, r *http.Request) {
	input := models.InputUser{}

	err := request.DecodeJSON(w, r, &input)
	if err != nil {
		app.badRequest(w, r, err)
		return
	}

	input.Validator.CheckField(input.Password != "", "Password", "Password is required")
	input.Validator.CheckField(input.FirstName != "", "FirstName", "FirstName is required")
	input.Validator.CheckField(input.SecondName != "", "SecondName", "SecondName is required")
	input.Validator.CheckField(!time.Time(input.Birthdate).IsZero(), "Birthdate", "Birthdate is required")
	input.Validator.CheckField(input.Biography != "", "Biography", "Biography is required")
	input.Validator.CheckField(input.City != "", "City", "City is required")

	input.Validator.CheckField(len(input.Password) >= 4, "Password", "Password is too short")
	input.Validator.CheckField(len(input.Password) <= 72, "Password", "Password is too long")
	input.Validator.CheckField(len(input.FirstName) <= 50, "FirstName", "FirstName is too long")
	input.Validator.CheckField(len(input.SecondName) <= 50, "SecondName", "SecondName is too long")
	input.Validator.CheckField(len(input.City) <= 50, "City", "City is too long")

	input.Validator.CheckField(validator.NotIn(input.Password, password.CommonPasswords...), "Password", "Password is too common")

	if input.Validator.HasErrors() {
		app.failedValidation(w, r, input.Validator)
		return
	}

	hashedPassword, err := password.Hash(input.Password)
	if err != nil {
		app.serverError(w, r, err)
		return
	}
	var userIDToResp string
	err = app.db.Transactional(context.Background(), func(ctx context.Context) error {
		userID, err := app.db.InsertToUserData(ctx, database.UserData{
			FirstName:  input.FirstName,
			SecondName: input.SecondName,
			Birthdate:  time.Time(input.Birthdate),
			Sex:        input.Sex,
			Biography:  input.Biography,
			City:       input.City,
		})
		if err != nil {
			app.serverError(w, r, err)
			return err
		}
		err = app.db.InsertToUserAuth(ctx, userID, hashedPassword)
		if err != nil {
			app.serverError(w, r, err)
			return err
		}
		userIDToResp = userID
		return nil
	})
	if err != nil {
		app.serverError(w, r, fmt.Errorf("transaction err: %w", err))
	}

	data := map[string]string{
		"user_id": userIDToResp,
	}

	err = response.JSON(w, http.StatusOK, data)
	if err != nil {
		app.serverError(w, r, err)
	}
}
func (app *application) getUserByUserID(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())
	userID := params.ByName("id")
	if !helpers.IsValidUUID(userID) {
		app.badRequest(w, r, fmt.Errorf("wrong uuid format"))
	}

	data, err := app.db.GetUserDataByID(userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			app.notFound(w, r)
			return
		}
		app.serverError(w, r, err)
		return
	}

	respData := models.OutputUser{}
	respData.FirstName = data.FirstName
	respData.SecondName = data.SecondName
	respData.City = data.City
	respData.Biography = data.Biography
	respData.Sex = data.Sex
	respData.Birthdate = models.JsonBirthDate(data.Birthdate)
	respData.Age = age.CalculateToNow(data.Birthdate)

	err = response.JSON(w, http.StatusOK, respData)
	if err != nil {
		app.serverError(w, r, err)
	}
}

func (app *application) searchUser(w http.ResponseWriter, r *http.Request) {
	input := models.InputUserSearch{}

	err := request.DecodeJSON(w, r, &input)
	if err != nil {
		app.badRequest(w, r, err)
		return
	}

	input.Validator.CheckField(input.FirstName != "", "FirstName", "FirstName is required")
	input.Validator.CheckField(input.SecondName != "", "SecondName", "SecondName is required")

	input.Validator.CheckField(len(input.FirstName) <= 50, "FirstName", "FirstName is too long")
	input.Validator.CheckField(len(input.SecondName) <= 50, "SecondName", "SecondName is too long")

	if input.Validator.HasErrors() {
		app.failedValidation(w, r, input.Validator)
		return
	}

	usersData, err := app.db.SearchUserData(input.FirstName, input.SecondName)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	if len(usersData) == 0 {
		app.notFound(w, r)
		return
	}

	var respData []models.OutputUser
	for _, userData := range usersData {
		respData = append(respData, models.OutputUser{
			FirstName:  userData.FirstName,
			SecondName: userData.SecondName,
			Sex:        userData.Sex,
			Biography:  userData.Biography,
			City:       userData.City,
			Birthdate:  models.JsonBirthDate(userData.Birthdate),
			Age:        age.CalculateToNow(userData.Birthdate),
		})
	}

	err = response.JSON(w, http.StatusOK, respData)
	if err != nil {
		app.serverError(w, r, err)
	}
}

func (app *application) createAuthenticationToken(w http.ResponseWriter, r *http.Request) {
	input := models.InputAuthToken{}

	err := request.DecodeJSON(w, r, &input)
	if err != nil {
		app.badRequest(w, r, err)
		return
	}
	input.Validator.CheckField(input.UserID != "", "ID", "ID is required")
	input.Validator.CheckField(input.Password != "", "Password", "Password is required")
	if !helpers.IsValidUUID(input.UserID) {
		app.badRequest(w, r, fmt.Errorf("wrong uuid format"))
		return
	}

	user, err := app.db.GetUserAuthByUserID(input.UserID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			app.notFound(w, r)
			return
		}
		app.serverError(w, r, err)
		return
	}

	input.Validator.CheckField(user != nil, "UserID", "UserID address could not be found")

	if user != nil {
		passwordMatches, err := password.Matches(input.Password, user.HashedPassword)
		if err != nil {
			app.serverError(w, r, err)
			return
		}

		input.Validator.CheckField(input.Password != "", "Password", "Password is required")
		input.Validator.CheckField(passwordMatches, "Password", "Password is incorrect")
	}

	if input.Validator.HasErrors() {
		app.failedValidation(w, r, input.Validator)
		return
	}

	var claims jwt.Claims
	claims.Subject = user.UserID

	expiry := time.Now().Add(24 * time.Hour)
	claims.Issued = jwt.NewNumericTime(time.Now())
	claims.NotBefore = jwt.NewNumericTime(time.Now())
	claims.Expires = jwt.NewNumericTime(expiry)

	claims.Issuer = app.config.baseURL
	claims.Audiences = []string{app.config.baseURL}

	jwtBytes, err := claims.HMACSign(jwt.HS256, []byte(app.config.jwt.secretKey))
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	data := map[string]string{
		"AuthenticationToken":       string(jwtBytes),
		"AuthenticationTokenExpiry": expiry.Format(time.RFC3339),
	}

	err = response.JSON(w, http.StatusOK, data)
	if err != nil {
		app.serverError(w, r, err)
	}
}

func (app *application) addFriend(w http.ResponseWriter, r *http.Request) {
	authenticatedUser := contextGetAuthenticatedUser(r)
	params := httprouter.ParamsFromContext(r.Context())
	friendID := params.ByName("id")
	if !helpers.IsValidUUID(friendID) {
		app.badRequest(w, r, fmt.Errorf("wrong uuid format"))
	}
	_, err := app.db.GetUserDataByID(friendID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			app.notFound(w, r)
			return
		}
		app.serverError(w, r, err)
		return
	}
	err = app.db.InsertToUserFriend(r.Context(), database.UserFriend{
		UserID:       authenticatedUser.UserID,
		FriendUserID: friendID,
	})
	if err != nil {
		app.serverError(w, r, err)
		return
	}
}

func (app *application) deleteFriend(w http.ResponseWriter, r *http.Request) {
	authenticatedUser := contextGetAuthenticatedUser(r)
	params := httprouter.ParamsFromContext(r.Context())
	friendID := params.ByName("id")
	if !helpers.IsValidUUID(friendID) {
		app.badRequest(w, r, fmt.Errorf("wrong uuid format"))
	}
	err := app.db.DeleteUserFriendByUserID(authenticatedUser.UserID, friendID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			app.notFound(w, r)
			return
		}
		app.serverError(w, r, err)
		return
	}
}

func (app *application) createPost(w http.ResponseWriter, r *http.Request) {
	authenticatedUser := contextGetAuthenticatedUser(r)
	input := models.InputCreatePost{}

	err := request.DecodeJSON(w, r, &input)
	if err != nil {
		app.badRequest(w, r, err)
		return
	}

	input.Validator.CheckField(input.Text != "", "text", "Text is required")
	err = app.db.InsertToUserPosts(r.Context(), database.UserPost{
		UserID: authenticatedUser.UserID,
		Post:   input.Text,
	})
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	payload := database.Payload{
		UserID: authenticatedUser.UserID,
		Post:   input.Text,
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	err = app.db.CreateTask(r.Context(), worker.AddPost, string(jsonPayload))
	if err != nil {
		app.serverError(w, r, err)
		return
	}

}

func (app *application) updatePost(w http.ResponseWriter, r *http.Request) {
	authenticatedUser := contextGetAuthenticatedUser(r)
	input := models.InputUpdatePost{}

	err := request.DecodeJSON(w, r, &input)
	if err != nil {
		app.badRequest(w, r, err)
		return
	}

	input.Validator.CheckField(input.PostID != "", "id", "ID is required")
	input.Validator.CheckField(input.Text != "", "text", "Text is required")
	err = app.db.UpdateUserPosts(r.Context(), database.UserPost{
		ID:     input.PostID,
		UserID: authenticatedUser.UserID,
		Post:   input.Text,
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			app.notFound(w, r)
			return
		}
		app.serverError(w, r, err)
		return
	}
}

func (app *application) deletePost(w http.ResponseWriter, r *http.Request) {
	authenticatedUser := contextGetAuthenticatedUser(r)
	params := httprouter.ParamsFromContext(r.Context())
	postID := params.ByName("id")
	if !helpers.IsValidUUID(postID) {
		app.badRequest(w, r, fmt.Errorf("wrong uuid format"))
	}

	err := app.db.DeleteUserPosts(r.Context(), database.UserPost{
		Post:   postID,
		UserID: authenticatedUser.UserID,
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			app.notFound(w, r)
			return
		}
		app.serverError(w, r, err)
		return
	}
}

func (app *application) getPost(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())
	postID := params.ByName("id")
	if !helpers.IsValidUUID(postID) {
		app.badRequest(w, r, fmt.Errorf("wrong uuid format"))
	}

	post, err := app.db.GetUserPostByPostID(postID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			app.notFound(w, r)
			return
		}
		app.serverError(w, r, err)
		return
	}

	respData := models.OutputPost{}
	respData.Text = post.Post

	err = response.JSON(w, http.StatusOK, respData)
	if err != nil {
		app.serverError(w, r, err)
	}
}

func (app *application) getFeed(w http.ResponseWriter, r *http.Request) {
	authenticatedUser := contextGetAuthenticatedUser(r)

	posts, err := app.cache.Get(authenticatedUser.UserID)
	if err != nil {
		app.serverError(w, r, err)
	}

	respData := []models.OutputPost{}
	for _, post := range posts {
		respData = append(respData, models.OutputPost{Text: post})
	}

	err = response.JSON(w, http.StatusOK, respData)
	if err != nil {
		app.serverError(w, r, err)
	}
}
func (app *application) protected(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("This is a protected handler"))
}
