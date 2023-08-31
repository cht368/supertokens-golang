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

package sessionwithjwt

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/supertokens/supertokens-golang/recipe/openid/openidmodels"
	"github.com/supertokens/supertokens-golang/recipe/session/sessmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
)

func newSessionWithJWTContainer(originalSessionClass sessmodels.SessionContainer, openidRecipeImplementation openidmodels.RecipeInterface) sessmodels.SessionContainer {

	originalUpdateAccessTokenPayloadWithContext := originalSessionClass.UpdateAccessTokenPayloadWithContext
	originalSessionClass.UpdateAccessTokenPayloadWithContext = func(newAccessTokenPayload map[string]interface{}, userContext supertokens.UserContext) error {
		if newAccessTokenPayload == nil {
			newAccessTokenPayload = map[string]interface{}{}
		}
		accessTokenPayload := originalSessionClass.GetAccessTokenPayloadWithContext(userContext)
		jwtPropertyName, ok := accessTokenPayload[ACCESS_TOKEN_PAYLOAD_JWT_PROPERTY_NAME_KEY]

		if !ok {
			return originalUpdateAccessTokenPayloadWithContext(newAccessTokenPayload, userContext)
		}

		existingJWT := accessTokenPayload[jwtPropertyName.(string)].(string)

		currentTimeInSeconds := uint64(time.Now().UnixNano() / 1000000000) // time in seconds

		claims := jwt.MapClaims{}
		decodedPayload := map[string]interface{}{}
		_, _, err := new(jwt.Parser).ParseUnverified(existingJWT, claims)
		if err != nil {
			return err
		}
		for key, val := range claims {
			decodedPayload[key] = val
		}

		jwtExpiry := uint64(decodedPayload["exp"].(float64)) - currentTimeInSeconds

		if jwtExpiry < 1 {
			// it can come here if someone calls this function well after
			// the access token and the jwt payload have expired (which can happen if an API takes a VERY long time). In this case, we still want the jwt payload to update, but the resulting JWT should
			// not be alive for too long (since it's expired already). So we set it to
			// 1 second lifetime.
			jwtExpiry = 1
		}

		newAccessTokenPayload, err = addJWTToAccessTokenPayload(newAccessTokenPayload, jwtExpiry, originalSessionClass.GetUserIDWithContext(userContext), jwtPropertyName.(string), openidRecipeImplementation, userContext)
		if err != nil {
			return err
		}

		return originalUpdateAccessTokenPayloadWithContext(newAccessTokenPayload, userContext)
	}
	return originalSessionClass
}
