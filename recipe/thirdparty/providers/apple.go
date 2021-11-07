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

package providers

import (
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/supertokens/supertokens-golang/recipe/thirdparty/api"
	"github.com/supertokens/supertokens-golang/recipe/thirdparty/tpmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
)

const appleID = "apple"

func Apple(config tpmodels.AppleConfig) tpmodels.TypeProvider {
	return tpmodels.TypeProvider{
		ID: appleID,
		Get: func(redirectURI, authCodeFromRequest *string) tpmodels.TypeProviderGetResponse {
			accessTokenAPIURL := "https://appleid.apple.com/auth/token"
			clientSecret, err := getClientSecret(config.ClientID, config.ClientSecret.KeyId, config.ClientSecret.TeamId, config.ClientSecret.PrivateKey)
			if err != nil {
				panic(err)
			}
			accessTokenAPIParams := map[string]string{
				"client_id":     config.ClientID,
				"client_secret": clientSecret,
				"grant_type":    "authorization_code",
			}
			if authCodeFromRequest != nil {
				accessTokenAPIParams["code"] = *authCodeFromRequest
			}
			if redirectURI != nil {
				accessTokenAPIParams["redirect_uri"] = *redirectURI
			}

			authorisationRedirectURL := "https://appleid.apple.com/auth/authorize"
			scopes := []string{"email"}
			if config.Scope != nil {
				scopes = config.Scope
			}

			var additionalParams map[string]interface{} = nil
			if config.AuthorisationRedirect != nil && config.AuthorisationRedirect.Params != nil {
				additionalParams = config.AuthorisationRedirect.Params
			}

			authorizationRedirectParams := map[string]interface{}{
				"scope":         strings.Join(scopes, " "),
				"response_mode": "form_post",
				"response_type": "code",
				"client_id":     config.ClientID,
			}
			for key, value := range additionalParams {
				authorizationRedirectParams[key] = value
			}

			return tpmodels.TypeProviderGetResponse{
				AccessTokenAPI: tpmodels.AccessTokenAPI{
					URL:    accessTokenAPIURL,
					Params: accessTokenAPIParams,
				},
				AuthorisationRedirect: tpmodels.AuthorisationRedirect{
					URL:    authorisationRedirectURL,
					Params: authorizationRedirectParams,
				},
				GetProfileInfo: func(authCodeResponse interface{}) (tpmodels.UserInfo, error) {
					claims := jwt.MapClaims{}
					_, _, err := new(jwt.Parser).ParseUnverified(authCodeResponse.(map[string]interface{})["id_token"].(string), claims)
					if err != nil {
						return tpmodels.UserInfo{}, err
					}

					var email string
					var isVerified bool
					var id string
					for key, val := range claims {
						if key == "sub" {
							id = val.(string)
						} else if key == "email" {
							email = val.(string)
						} else if key == "email_verified" {
							isVerified = val.(string) == "true"
						}
					}
					return tpmodels.UserInfo{
						ID: id,
						Email: &tpmodels.EmailStruct{
							ID:         email,
							IsVerified: isVerified,
						},
					}, nil
				},
				GetClientId: func() string {
					return config.ClientID
				},
				GetRedirectURI: func() (string, error) {
					supertokens, err := supertokens.GetInstanceOrThrowError()
					if err != nil {
						return "", err
					}
					return supertokens.AppInfo.APIDomain.GetAsStringDangerous() + "/callback/apple", nil
				},
			}
		},
		IsDefault: config.IsDefault,
	}
}

func getClientSecret(clientId, keyId, teamId, privateKey string) (string, error) {
	claims := jwt.StandardClaims{
		ExpiresAt: time.Now().Unix() + 86400*180,
		IssuedAt:  time.Now().Unix(),
		Audience:  "https://appleid.apple.com",
		Id:        keyId,
		Subject:   api.GetActualClientIdFromDevelopmentClientId(clientId),
		Issuer:    teamId,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodES256, claims)

	ecdsaPrivateKey, err := getECDSPrivateKey(privateKey)
	if err != nil {
		return "", err
	}

	// Finally sign the token with the value of type *ecdsa.PrivateKey
	return token.SignedString(ecdsaPrivateKey)
}

func getECDSPrivateKey(privateKey string) (*ecdsa.PrivateKey, error) {
	block, _ := pem.Decode([]byte(privateKey))
	// Check if it's a private key
	if block == nil || block.Type != "PRIVATE KEY" {
		return nil, errors.New("failed to decode PEM block containing private key")
	}
	// Get the encoded bytes
	x509Encoded := block.Bytes

	// Now you need an instance of *ecdsa.PrivateKey
	parsedKey, err := x509.ParsePKCS8PrivateKey(x509Encoded) // EDIT to x509Encoded from p8bytes
	if err != nil {
		return nil, err
	}

	ecdsaPrivateKey, ok := parsedKey.(*ecdsa.PrivateKey)
	if !ok {
		return nil, errors.New("not ecdsa private key")
	}
	return ecdsaPrivateKey, nil
}
