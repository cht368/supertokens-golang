package emailverification

import (
	"net/http"

	"github.com/supertokens/supertokens-golang/errors"
	"github.com/supertokens/supertokens-golang/recipe/emailverification/api"
	"github.com/supertokens/supertokens-golang/recipe/emailverification/models"
	"github.com/supertokens/supertokens-golang/supertokens"
)

const RECIPE_ID = "emailverification"

type Recipe struct {
	RecipeModule supertokens.RecipeModule
	Config       models.TypeNormalisedInput
	RecipeImpl   models.RecipeImplementation
	APIImpl      models.APIImplementation
}

var r *Recipe = nil

func MakeRecipe(recipeId string, appInfo supertokens.NormalisedAppinfo, config models.TypeInput) Recipe {
	querierInstance, _ := supertokens.GetNewQuerierInstanceOrThrowError(recipeId)
	recipeModuleInstance := supertokens.MakeRecipeModule(recipeId, appInfo, HandleAPIRequest, GetAllCORSHeaders, GetAPIsHandled)
	verifiedConfig := validateAndNormaliseUserInput(appInfo, config)
	recipeImplementation := MakeRecipeImplementation(*querierInstance)

	return Recipe{
		RecipeModule: recipeModuleInstance,
		Config:       verifiedConfig,
		RecipeImpl:   verifiedConfig.Override.Functions(recipeImplementation),
		APIImpl:      verifiedConfig.Override.APIs(api.MakeAPIImplementation()),
	}
}

func GetInstanceOrThrowError() (*Recipe, error) {
	if r != nil {
		return r, nil
	}
	return nil, errors.BadInputError{Msg: "Initialisation not done. Did you forget to call the SuperTokens.init function?"}
}

func RecipeInit(config models.TypeInput) supertokens.RecipeListFunction {
	return func(appInfo supertokens.NormalisedAppinfo) (*supertokens.RecipeModule, error) {
		if r == nil {
			recipe := MakeRecipe(RECIPE_ID, appInfo, config)
			r = &recipe
			return &r.RecipeModule, nil
		}
		return nil, errors.BadInputError{Msg: "Emailverification recipe has already been initialised. Please check your code for bugs."}
	}
}

func (r *Recipe) CreateEmailVerificationToken(userID, email string) (string, error) {
	response, err := r.RecipeImpl.CreateEmailVerificationToken(userID, email)
	if err != nil {
		return "", err
	}
	if response.Ok != nil {
		return response.Ok.Token, nil
	}
	return "", errors.BadInputError{Msg: "Email has already been verified"}
}

func (r *Recipe) VerifyEmailUsingToken(token string) (*models.User, error) {
	response, err := r.RecipeImpl.VerifyEmailUsingToken(token)
	if err != nil {
		return nil, err
	}
	if response.Ok != nil {
		return &response.Ok.User, nil
	}
	return nil, errors.BadInputError{Msg: "Invalid email verification token"}
}

// implement RecipeModule

func GetAPIsHandled() ([]supertokens.APIHandled, error) {
	generateEmailVerifyTokenAPI, err := supertokens.NewNormalisedURLPath(GenerateEmailVerifyTokenAPI)
	if err != nil {
		return nil, err
	}
	emailVerifyAPI, err := supertokens.NewNormalisedURLPath(EmailVerifyAPI)
	if err != nil {
		return nil, err
	}
	return []supertokens.APIHandled{{
		Method:                 "post",
		PathWithoutAPIBasePath: *generateEmailVerifyTokenAPI,
		ID:                     GenerateEmailVerifyTokenAPI,
		Disabled:               r.APIImpl.GenerateEmailVerifyTokenPOST == nil,
	}, {
		Method:                 "post",
		PathWithoutAPIBasePath: *emailVerifyAPI,
		ID:                     EmailVerifyAPI,
		Disabled:               r.APIImpl.VerifyEmailPOST == nil,
	}, {
		Method:                 "get",
		PathWithoutAPIBasePath: *emailVerifyAPI,
		ID:                     EmailVerifyAPI,
		Disabled:               r.APIImpl.IsEmailVerifiedGET == nil,
	}}, nil
}

func HandleAPIRequest(id string, req *http.Request, res http.ResponseWriter, theirHandler http.HandlerFunc, _ supertokens.NormalisedURLPath, _ string) error {
	options := models.APIOptions{
		Config:               r.Config,
		RecipeID:             r.RecipeModule.GetRecipeID(),
		RecipeImplementation: r.RecipeImpl,
		Req:                  req,
		Res:                  res,
		OtherHandler:         theirHandler,
	}
	var err error = nil
	if id == GenerateEmailVerifyTokenAPI {
		err = api.GenerateEmailVerifyToken(r.APIImpl, options)
	} else {
		err = api.EmailVerify(r.APIImpl, options)
	}
	return err
}

func GetAllCORSHeaders() []string {
	return []string{}
}
