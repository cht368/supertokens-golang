/* Copyright (c) 2021, VRAI Labs and/or its affiliates. All rights reserved.
 *
 * This software is licensed under the Apache License, Version 2.0 (the
 * "License") as published by the Apache Software Foundation.
 *
 * You may not use this file except in compliance with the License. You may
 * obtain a copy of the License at http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
 * WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
 * License for the specific language governing permissions and limitations
 * under the License.
 */

package thirdpartyemailpassword

import (
	"errors"

	"github.com/supertokens/supertokens-golang/ingredients/emaildelivery"
	"github.com/supertokens/supertokens-golang/recipe/emailpassword"
	"github.com/supertokens/supertokens-golang/recipe/emailpassword/epmodels"
	"github.com/supertokens/supertokens-golang/recipe/emailverification"
	"github.com/supertokens/supertokens-golang/recipe/emailverification/evmodels"
	"github.com/supertokens/supertokens-golang/recipe/thirdpartyemailpassword/emaildelivery/backwardCompatibilityService"
	"github.com/supertokens/supertokens-golang/recipe/thirdpartyemailpassword/tpepmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
)

func validateAndNormaliseUserInput(recipeInstance *Recipe, appInfo supertokens.NormalisedAppinfo, config *tpepmodels.TypeInput) (tpepmodels.TypeNormalisedInput, error) {
	typeNormalisedInput := makeTypeNormalisedInput(recipeInstance)

	if config != nil && config.SignUpFeature != nil {
		typeNormalisedInput.SignUpFeature = config.SignUpFeature
	}

	if config != nil && config.Providers != nil {
		typeNormalisedInput.Providers = config.Providers
	}

	typeNormalisedInput.EmailVerificationFeature = validateAndNormaliseEmailVerificationConfig(recipeInstance, config)

	if config != nil && config.ResetPasswordUsingTokenFeature != nil {
		typeNormalisedInput.ResetPasswordUsingTokenFeature = config.ResetPasswordUsingTokenFeature
	}

	typeNormalisedInput.GetEmailDeliveryConfig = func(recipeImpl tpepmodels.RecipeInterface, epRecipeImpl epmodels.RecipeInterface) emaildelivery.TypeInputWithService {
		sendPasswordResetEmail := emailpassword.DefaultCreateAndSendCustomPasswordResetEmail(appInfo)
		if config != nil && config.ResetPasswordUsingTokenFeature != nil && config.ResetPasswordUsingTokenFeature.CreateAndSendCustomEmail != nil {
			sendPasswordResetEmail = config.ResetPasswordUsingTokenFeature.CreateAndSendCustomEmail
		}

		sendEmailVerificationEmail := emailverification.DefaultCreateAndSendCustomEmail(appInfo)
		if typeNormalisedInput.EmailVerificationFeature.CreateAndSendCustomEmail != nil {
			sendEmailVerificationEmail = typeNormalisedInput.EmailVerificationFeature.CreateAndSendCustomEmail
		}

		emailService := backwardCompatibilityService.MakeBackwardCompatibilityService(recipeImpl, epRecipeImpl, appInfo, sendEmailVerificationEmail, sendPasswordResetEmail)
		if config != nil && config.EmailDelivery != nil && config.EmailDelivery.Service != nil {
			emailService = *config.EmailDelivery.Service
		}
		result := emaildelivery.TypeInputWithService{
			Service: emailService,
		}
		if config != nil && config.EmailDelivery != nil && config.EmailDelivery.Override != nil {
			result.Override = config.EmailDelivery.Override
		}

		return result
	}

	if config != nil && config.Override != nil {
		if config.Override.Functions != nil {
			typeNormalisedInput.Override.Functions = config.Override.Functions
		}
		if config.Override.APIs != nil {
			typeNormalisedInput.Override.APIs = config.Override.APIs
		}
		if config.Override.EmailVerificationFeature != nil {
			typeNormalisedInput.Override.EmailVerificationFeature = config.Override.EmailVerificationFeature
		}
	}

	return typeNormalisedInput, nil
}

func makeTypeNormalisedInput(recipeInstance *Recipe) tpepmodels.TypeNormalisedInput {
	return tpepmodels.TypeNormalisedInput{
		SignUpFeature:                  nil,
		Providers:                      nil,
		ResetPasswordUsingTokenFeature: nil,
		EmailVerificationFeature:       validateAndNormaliseEmailVerificationConfig(recipeInstance, nil),
		Override: tpepmodels.OverrideStruct{
			Functions: func(originalImplementation tpepmodels.RecipeInterface) tpepmodels.RecipeInterface {
				return originalImplementation
			},
			APIs: func(originalImplementation tpepmodels.APIInterface) tpepmodels.APIInterface {
				return originalImplementation
			},
			EmailVerificationFeature: nil,
		},
	}
}

func validateAndNormaliseEmailVerificationConfig(recipeInstance *Recipe, config *tpepmodels.TypeInput) evmodels.TypeInput {
	emailverificationTypeInput := evmodels.TypeInput{
		GetEmailForUserID: recipeInstance.getEmailForUserId,
		Override:          nil,
	}

	if config != nil {
		if config.Override != nil {
			emailverificationTypeInput.Override = config.Override.EmailVerificationFeature
		}
		if config.EmailVerificationFeature != nil {
			if config.EmailVerificationFeature.GetEmailVerificationURL != nil {
				emailverificationTypeInput.GetEmailVerificationURL = func(user evmodels.User, userContext supertokens.UserContext) (string, error) {
					userInfo, err := (*recipeInstance.RecipeImpl.GetUserByID)(user.ID, userContext)
					if err != nil {
						return "", err
					}
					if userInfo == nil {
						return "", errors.New("Unknown User ID provided")
					}
					return config.EmailVerificationFeature.GetEmailVerificationURL(*userInfo, userContext)
				}
			}
		}
	}

	return emailverificationTypeInput
}
