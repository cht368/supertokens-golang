package session

import (
	"net/http"

	"github.com/supertokens/supertokens-golang/recipe/session/api"
	"github.com/supertokens/supertokens-golang/recipe/session/schema"
	"github.com/supertokens/supertokens-golang/supertokens"
)

const RECIPE_ID = "session"

type SessionRecipe struct {
	RecipeModule supertokens.RecipeModule
	Config       schema.TypeNormalisedInput
	RecipeImpl   schema.RecipeImplementation
	APIImpl      schema.APIImplementation
}

var r *SessionRecipe = nil

func NewRecipe(recipeId string, appInfo supertokens.NormalisedAppinfo, config *schema.TypeInput) SessionRecipe {
	querierInstance, _ := supertokens.GetNewQuerierInstanceOrThrowError(recipeId)
	recipeModuleInstance := supertokens.MakeRecipeModule(recipeId, appInfo, HandleAPIRequest, GetAllCORSHeaders, GetAPIsHandled)
	verifiedConfig, _ := validateAndNormaliseUserInput(r, appInfo, config)

	recipeImplementation := MakeRecipeImplementation(*querierInstance, verifiedConfig)

	return SessionRecipe{
		RecipeModule: recipeModuleInstance,
		Config:       verifiedConfig,
		RecipeImpl:   verifiedConfig.Override.Functions(recipeImplementation),
		APIImpl:      verifiedConfig.Override.APIs(api.MakeAPIImplementation()),
	}
}

// Implement RecipeModule

func GetAPIsHandled() ([]supertokens.APIHandled, error) {
	refreshAPIPath, err := supertokens.NewNormalisedURLPath(RefreshAPIPath)
	if err != nil {
		return nil, err
	}
	signoutAPIPath, err := supertokens.NewNormalisedURLPath(SignoutAPIPath)
	if err != nil {
		return nil, err
	}
	return []supertokens.APIHandled{{
		Method:                 "post",
		PathWithoutAPIBasePath: *refreshAPIPath,
		ID:                     RefreshAPIPath,
		Disabled:               r.APIImpl.RefreshPOST == nil,
	}, {
		Method:                 "post",
		PathWithoutAPIBasePath: *signoutAPIPath,
		ID:                     SignoutAPIPath,
		Disabled:               r.APIImpl.SignOutPOST == nil,
	}}, nil
}

func HandleAPIRequest(id string, req *http.Request, res http.ResponseWriter, thierhandler http.HandlerFunc, _ supertokens.NormalisedURLPath, _ string) error {
	options := schema.APIOptions{
		Config:               r.Config,
		RecipeID:             r.RecipeModule.GetRecipeID(),
		RecipeImplementation: r.RecipeImpl,
		Req:                  req,
		Res:                  res,
		OtherHandler:         thierhandler,
	}
	if id == RefreshAPIPath {
		api.HandleRefreshAPI(r.APIImpl, options)
	} else {
		return api.SignOutAPI(r.APIImpl, options)
	}
	return nil
}

func GetAllCORSHeaders() []string {
	return []string{antiCsrfHeaderKey, ridHeaderKey}
}
