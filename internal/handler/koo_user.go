package handler

// BOILERPLATE: This file demonstrates the handler pattern.
// Delete this file when bootstrapping a new project.
// See docs/BOOTSTRAPPING.md for details.

import (
	"github.com/gofiber/fiber/v2"

	"github.com/kootic/koogo/internal/dto"
	"github.com/kootic/koogo/internal/service"
	"github.com/kootic/koogo/pkg/koohttp"
)

type KooUserHandler interface {
	CreateUser(c *fiber.Ctx) error
	GetUserByID(c *fiber.Ctx) error
	GetUserPet(c *fiber.Ctx) error
}

type kooUserHandler struct {
	userService service.KooUserService
}

var _ KooUserHandler = (*kooUserHandler)(nil)

func NewKooUserHandler(userService service.KooUserService) KooUserHandler {
	return &kooUserHandler{
		userService: userService,
	}
}

// KooCreateUser godoc
//
//	@tags			Users
//	@Summary		Create a new user
//	@Description	Create a new user
//	@Accept			json
//	@Produce		json
//	@Param			kooCreateUserRequest	body		dto.KooCreateUserRequest	true	"Create user request"
//	@Success		200						{object}	dto.KooUserResponse
//	@Failure		400						{object}	koohttp.APIResponseError
//	@Failure		500						{object}	koohttp.APIResponseError
//	@Router			/v1/koo/users [post]
func (h *kooUserHandler) CreateUser(c *fiber.Ctx) error {
	req, err := koohttp.GetBodyAndValidate[dto.KooCreateUserRequest](c)
	if err != nil {
		return err
	}

	user, err := h.userService.KooCreateUser(c.Context(), req)
	if err != nil {
		return err
	}

	return koohttp.Success(c, user)
}

// KooGetUserByID godoc
//
//	@tags			Users
//	@Summary		Get a user by ID
//	@Description	Get a user by ID
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string	true	"User ID"
//	@Success		200	{object}	dto.KooUserResponse
//	@Failure		400	{object}	koohttp.APIResponseError
//	@Failure		404	{object}	koohttp.APIResponseError
//	@Failure		500	{object}	koohttp.APIResponseError
//	@Router			/v1/koo/users/{id} [get]
func (h *kooUserHandler) GetUserByID(c *fiber.Ctx) error {
	userID, err := koohttp.GetParamUUID(c, "userId")
	if err != nil {
		return err
	}

	user, err := h.userService.KooGetUserByID(c.Context(), userID)
	if err != nil {
		return err
	}

	return koohttp.Success(c, user)
}

// KooGetUserPet godoc
//
//	@tags			Users
//	@Summary		Get a user's pet
//	@Description	Get a user's pet
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string	true	"User ID"
//	@Success		200	{object}	dto.KooPetResponse
//	@Failure		400	{object}	koohttp.APIResponseError
//	@Failure		403	{object}	koohttp.APIResponseError
//	@Failure		404	{object}	koohttp.APIResponseError
//	@Failure		500	{object}	koohttp.APIResponseError
//	@Router			/v1/koo/users/{id}/pet [get]
func (h *kooUserHandler) GetUserPet(c *fiber.Ctx) error {
	userID, err := koohttp.GetParamUUID(c, "userId")
	if err != nil {
		return err
	}

	userPet, err := h.userService.KooGetPetByOwnerID(c.Context(), userID)
	if err != nil {
		return err
	}

	return koohttp.Success(c, userPet)
}
