package controllers

import (
	"api-enjor/app/repository"
	"api-enjor/internal"
	"api-enjor/pkg/models"
	"api-enjor/pkg/utils"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthController struct {
	DB *gorm.DB
}

type PayloadsClaims struct {
	jwt.StandardClaims
	Sub *models.ModuleProfile `json:"sub"`
}

func NewAuthController(DB *gorm.DB) AuthController {
	return AuthController{DB}
}

func (ac *AuthController) SigninUserController(c *fiber.Ctx) error {
	user := new(repository.Users)

	err := c.BodyParser(user)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.NewBaseErrorResponse(err.Error(), fiber.StatusBadRequest))
	}

	user.ID = uuid.New()
	hash, _ := bcrypt.GenerateFromPassword([]byte(user.Password), 14)
	user.Password = string(hash)

	err = ac.DB.Create(&user).Error
	if err != nil {
		return c.Status(fiber.StatusNotImplemented).JSON(models.NewBaseErrorResponse(err.Error(), fiber.StatusNotImplemented))
	}
	response := map[string]interface{}{
		"id":           user.ID,
		"e_mail":       user.E_mail,
		"profile_name": user.First_name,
		"activate":     user.Activate,
	}
	return c.Status(fiber.StatusCreated).JSON(models.NewBaseResponse(response, fiber.StatusCreated))
}

func (ac *AuthController) UserOauthController(c *fiber.Ctx) error {
	state := c.Query("state")
	if state == "" {
		return c.Status(fiber.StatusBadRequest).JSON(models.NewBaseErrorResponse(fiber.Map{
			"message": "state is notfound",
		}, fiber.StatusBadRequest))
	}

	url := internal.Oauth().AuthCodeURL(state)

	return c.Redirect(url)
}

func (ac *AuthController) CallbackUserController(c *fiber.Ctx) error {
	request := new(models.OauthRequest)
	if err := c.BodyParser(request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.NewBaseErrorResponse(err.Error(), fiber.StatusBadRequest))
	}

	config := internal.Oauth()
	tok, err := config.Exchange(c.Context(), request.Code)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	resp, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + tok.AccessToken)
	if err != nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	defer resp.Body.Close()
	content, err := ioutil.ReadAll(resp.Body)

	var info models.ModuleProfileOauth
	json.Unmarshal(content, &info)

	model_user := &models.ModuleProfile{
		User_name:    info.Name,
		E_mail:       info.Email,
		Activate:     1,
		Profile_name: info.Name,
		Picture:      info.Picture,
		Is_oauth:     true,
	}

	items, err := utils.GenerateTokenJWT(model_user, true)
	model_user.Access_token = items
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.NewBaseErrorResponse(err.Error(), fiber.StatusBadRequest))
	}

	return c.JSON(models.NewBaseResponse(model_user, fiber.StatusOK))
}

func (ac *AuthController) LoginUserControlles(c *fiber.Ctx) error {
	payload := new(models.SignInInput)
	var user repository.Users

	if err := c.BodyParser(payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.NewBaseErrorResponse(err.Error(), fiber.StatusBadRequest))
	}

	err := ac.DB.Raw(`SELECT id, user_name, password, e_mail, first_name, activate FROM "users" WHERE e_mail = ?`, payload.E_mail).Scan(&user).Error
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.NewBaseErrorResponse(err.Error(), fiber.StatusBadRequest))
	}

	if user.E_mail == "" {
		return c.Status(fiber.StatusNotFound).JSON(models.NewBaseErrorResponse(map[string]interface{}{
			"message": "User not found",
		}, fiber.StatusNotFound))
	}
	if is_passwor_hash := utils.CheckPasswordHash(payload.Password, user.Password); !is_passwor_hash {
		return c.Status(fiber.StatusNotAcceptable).JSON(models.NewBaseErrorResponse(map[string]interface{}{
			"message": "Password don't matching",
			"match":   false,
		}, fiber.StatusNotAcceptable))
	}

	model_user := &models.ModuleProfile{
		ID:           user.ID,
		User_name:    user.User_name,
		E_mail:       user.E_mail,
		Activate:     user.Activate,
		Profile_name: user.First_name,
		Is_oauth:     false,
	}

	access_token, _ := utils.GenerateTokenJWT(model_user, true)
	refresh_token, _ := utils.GenerateTokenJWT(model_user, false)

	decode, _ := utils.Decode(os.Getenv("JWT_SECRET"))
	token, _err := jwt.ParseWithClaims(access_token, &utils.PayloadsClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(decode), nil
	})

	if _err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.NewBaseErrorResponse(_err.Error(), fiber.StatusBadRequest))
	}

	claims := token.Claims.(*utils.PayloadsClaims)

	return c.JSON(models.NewBaseResponse(utils.GenerateJWTOption{
		Access_token:  access_token,
		Refresh_token: refresh_token,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: claims.ExpiresAt,
			IssuedAt:  claims.IssuedAt,
		},
	}, fiber.StatusOK))
}

func (ac *AuthController) RefreshTokenControlles(c *fiber.Ctx) error {
	user := utils.ModuleUser(c)

	decode, _ := utils.Decode(os.Getenv("JWT_SECRET"))

	access_token, _ := utils.GenerateTokenJWT(user, true)

	token, _ := jwt.ParseWithClaims(access_token, &utils.PayloadsClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(decode), nil
	})

	claims := token.Claims.(*utils.PayloadsClaims)

	return c.JSON(models.NewBaseResponse(utils.GenerateJWTOption{
		Access_token: access_token,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: claims.ExpiresAt,
			IssuedAt:  claims.IssuedAt,
		},
	}, fiber.StatusOK))
}
