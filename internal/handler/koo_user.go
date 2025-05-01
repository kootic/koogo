package handler

import (
	"github.com/gofiber/fiber/v2"

	"github.com/kootic/koogo/internal/dto"
	"github.com/kootic/koogo/internal/service"
	"github.com/kootic/koogo/pkg/koohttp"
)

type KooUserHandler interface {
	KooCreateUser(c *fiber.Ctx) error
	KooGetUserByID(c *fiber.Ctx) error
	KooGetUserPet(c *fiber.Ctx) error
}

type kooUserHandler struct {
	kooUserService service.KooUserService
}

var _ KooUserHandler = (*kooUserHandler)(nil)

func NewKooUserHandler(services *service.Services) KooUserHandler {
	return &kooUserHandler{
		kooUserService: services.KooUserService,
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
func (h *kooUserHandler) KooCreateUser(c *fiber.Ctx) error {
	kooCreateUserRequest, err := koohttp.GetBodyAndValidate[dto.KooCreateUserRequest](c)
	if err != nil {
		return err
	}

	kooUser, err := h.kooUserService.KooCreateUser(c.Context(), kooCreateUserRequest)
	if err != nil {
		return err
	}

	return koohttp.Success(c, kooUser)
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
func (h *kooUserHandler) KooGetUserByID(c *fiber.Ctx) error {
	userID, err := koohttp.GetParamUUID(c, "userId")
	if err != nil {
		return err
	}

	user, err := h.kooUserService.KooGetUserByID(c.Context(), userID)
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
func (h *kooUserHandler) KooGetUserPet(c *fiber.Ctx) error {
	userID, err := koohttp.GetParamUUID(c, "userId")
	if err != nil {
		return err
	}

	userPet, err := h.kooUserService.KooGetPetByOwnerID(c.Context(), userID)
	if err != nil {
		return err
	}

	return koohttp.Success(c, userPet)
}
