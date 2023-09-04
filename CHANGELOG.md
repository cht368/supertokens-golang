# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [unreleased]

## [0.8.4] - 2022-08-31

- Adds logic to retry network calls if the core returns status 429

## [0.8.3] - 2022-07-30
### Added
- Adds test to verify that session container uses overridden functions
- Adds with-go-zero example: https://github.com/supertokens/supertokens-golang/issues/157
- UserId Mapping functionality and compatibility with CDI 2.15
- Adds `CreateUserIdMapping`, `GetUserIdMapping`, `DeleteUserIdMapping`, `UpdateOrDeleteUserIdMappingInfo` functions to supertokens package


## [0.8.2] - 2022-07-18

### Fixes:
- Fixes JWKS Keyfunc call that resulted in a goroutine leak: https://github.com/supertokens/supertokens-golang/issues/155

## [0.8.1] - 2022-07-12

### Fixes:
- Fixes issue with 404 status being sent for apple redirect callback route.

## [0.8.0] - 2022-07-08

### Breaking change:
-   Changes session recipe interfaces to not return an `UNAUTHORISED` error when the input is a sessionHandle: https://github.com/supertokens/backend/issues/83
-   `GetSessionInformation` now returns `nil` is the session does not exist
-   `UpdateSessionData` now returns `nil` if the input `sessionHandle` does not exist.
-   `UpdateAccessTokenPayload` now returns `false` if the input `sessionHandle` does not exist.
-   `RegenerateAccessToken` now returns `nil` if the input access token's `sessionHandle` does not exist.
-   The session container functions have not changed in behaviour and return errors if `sessionHandle` does not exist. This works on the current session.

### Fixes
-   Clears cookies when RevokeSession is called using the session container, even if the session did not exist from before: https://github.com/supertokens/supertokens-node/issues/343

### Adds:
-   Adds default userContext for API calls that contains the request object. It can be used in APIs / functions override like so:

```golang
SignIn: func (..., userContext supertokens.UserContext) {
    if _default, ok := (*userContext)["_default"].(map[string]interface{}); ok {
        if req, ok := _default["request"].(*http.Request); ok {
            // do something here with the request object
        }
    }
}
```

## [0.7.2] - 2022-06-29
-   Adds unit tests for resend email & sms services for passwordless and thirdpartypasswordless recipes
-   Adds User Roles recipe and compatibility with CDI 2.14

## [0.7.1] - 2022-06-27
-   Fixes panic while returning empty result object with nil error in the API overrides. Related to https://github.com/supertokens/supertokens-golang/issues/107

## [0.7.0] - 2022-06-23
### Breaking change
-   Renamed `SMTPServiceConfig` to `SMTPSettings`
-   Changed type of `Secure` in `SMTPSettings` from `*bool` to `bool`
-   Renamed `SMTPServiceFromConfig` to `SMTPFrom`
-   Renamed `SMTPGetContentResult` to `EmailContent`
-   Renamed `SMTPTypeInput` to `SMTPServiceConfig`
-   Renamed field `SMTPSettings` to `Settings` in `SMTPServiceConfig`
-   Renamed `SMTPServiceInterface` to `SMTPInterface`
-   Renamed all instances of `MakeSmtpService` to `MakeSMTPService`
-   All instances of `MakeSMTPService` returns `*EmailDeliveryInterface` instead of `EmailDeliveryInterface`
-   Renamed `TwilioServiceConfig` to `TwilioSettings`
-   Renamed `TwilioGetContentResult` to `SMSContent`
-   Renamed `TwilioTypeInput` to `TwilioServiceConfig`
-   Renamed field `TwilioSettings` to `Settings` in `TwilioServiceConfig`
-   Changed types of fields `From` and `MessagingServiceSid` in `TwilioSettings` from `*string` to `string`
-   Renamed `MakeSupertokensService` to `MakeSupertokensSMSService`
-   All instances of `MakeSupertokensSMSService` and `MakeTwilioService` returns `*SmsDeliveryInterface` instead of `SmsDeliveryInterface`
-   Removed `SupertokensServiceConfig` and `MakeSupertokensSMSService` accepts `apiKey` directly instead of `SupertokensServiceConfig`
-   Renamed `TwilioServiceInterface` to `TwilioInterface`
- Removes support for FDIs that are < 1.14

### Added
-   Exposed `MakeSMTPService` from emailverification, emailpassword, passwordless, thirdparty, thirdpartyemailpassword and thirdpartypasswordless recipes
-   Exposed `MakeSupertokensSMSService` and `MakeTwilioService` from passwordless and thirdpartypasswordless recipes

### Fixes
- Fixes Cookie SameSite config validation.
- Changes `getEmailForUserIdForEmailVerification` function inside thirdpartypasswordless to take into account passwordless emails and return an empty string in case a passwordless email doesn't exist. This helps situations where the dev wants to customise the email verification functions in the thirdpartypasswordless recipe.

## [0.6.8] - 2022-06-17
### Added
- `EmailDelivery` user config for Emailpassword, Thirdparty, ThirdpartyEmailpassword, Passwordless and ThirdpartyPasswordless recipes.
- `SmsDelivery` user config for Passwordless and ThirdpartyPasswordless recipes.
- `Twilio` service integration for SmsDelivery ingredient.
- `SMTP` service integration for EmailDelivery ingredient.
- `Supertokens` service integration for SmsDelivery ingredient.

### Deprecated
- For Emailpassword recipe input config, `ResetPasswordUsingTokenFeature.CreateAndSendCustomEmail` and `EmailVerificationFeature.CreateAndSendCustomEmail` have been deprecated.
- For Thirdparty recipe input config, `EmailVerificationFeature.CreateAndSendCustomEmail` has been deprecated.
- For ThirdpartyEmailpassword recipe input config, `ResetPasswordUsingTokenFeature.CreateAndSendCustomEmail` and `EmailVerificationFeature.CreateAndSendCustomEmail` have been deprecated.
- For Passwordless recipe input config, `CreateAndSendCustomEmail` and `CreateAndSendCustomTextMessage` have been deprecated.
- For ThirdpartyPasswordless recipe input config, `CreateAndSendCustomEmail`, `CreateAndSendCustomTextMessage` and `EmailVerificationFeature.CreateAndSendCustomEmail` have been deprecated.

### Migration

Following is an example of ThirdpartyPasswordless recipe migration. If your existing code looks like

```go
func passwordlessLoginEmail(email string, userInputCode *string, urlWithLinkCode *string, codeLifetime uint64, preAuthSessionId string, userContext supertokens.UserContext) error {
	// some custom logic
}

func passwordlessLoginSms(phoneNumber string, userInputCode *string, urlWithLinkCode *string, codeLifetime uint64, preAuthSessionId string, userContext supertokens.UserContext) error {
	// some custom logic
}

func verifyEmail(user tplmodels.User, emailVerificationURLWithToken string, userContext supertokens.UserContext) {
	// some custom logic
}

supertokens.Init(supertokens.TypeInput{
    AppInfo: supertokens.AppInfo{
        AppName:       "...",
        APIDomain:     "...",
        WebsiteDomain: "...",
    },
    RecipeList: []supertokens.Recipe{
        thirdpartypasswordless.Init(tplmodels.TypeInput{
            FlowType: "...",
            ContactMethodEmailOrPhone: plessmodels.ContactMethodEmailOrPhoneConfig{
                Enabled: true,
                CreateAndSendCustomEmail: passwordlessLoginEmail,
                CreateAndSendCustomTextMessage: passwordlessLoginSms,
            },
            EmailVerificationFeature: &tplmodels.TypeInputEmailVerificationFeature{
                CreateAndSendCustomEmail: verifyEmail,
            },
        }),
    },
})
```

After migration to using new `EmailDelivery` and `SmsDelivery` config, your code would look like:
```go
func passwordlessLoginEmail(email string, userInputCode *string, urlWithLinkCode *string, codeLifetime uint64, preAuthSessionId string, userContext supertokens.UserContext) error {
	// some custom logic
	return nil
}

func passwordlessLoginSms(phoneNumber string, userInputCode *string, urlWithLinkCode *string, codeLifetime uint64, preAuthSessionId string, userContext supertokens.UserContext) error {
	// some custom logic
	return nil
}

func verifyEmail(user tplmodels.User, emailVerificationURLWithToken string, userContext supertokens.UserContext) {
	// some custom logic
}

var sendEmail = func(input emaildelivery.EmailType, userContext supertokens.UserContext) error {
	if input.EmailVerification != nil {
		verifyEmail(tplmodels.User{ID: input.EmailVerification.User.ID, Email: &input.EmailVerification.User.Email}, input.EmailVerification.EmailVerifyLink, userContext)
	} else if input.PasswordlessLogin != nil {
		return passwordlessLoginEmail(input.PasswordlessLogin.Email, input.PasswordlessLogin.UserInputCode, input.PasswordlessLogin.UrlWithLinkCode, input.PasswordlessLogin.CodeLifetime, input.PasswordlessLogin.PreAuthSessionId, userContext)
	}
	return nil
}

var sendSms = func(input smsdelivery.SmsType, userContext supertokens.UserContext) error {
	if input.PasswordlessLogin != nil {
		return passwordlessLoginSms(input.PasswordlessLogin.PhoneNumber, input.PasswordlessLogin.UserInputCode, input.PasswordlessLogin.UrlWithLinkCode, input.PasswordlessLogin.CodeLifetime, input.PasswordlessLogin.PreAuthSessionId, userContext)
	}
	return nil
}

supertokens.Init(supertokens.TypeInput{
    AppInfo: supertokens.AppInfo{
        AppName:       "...",
        APIDomain:     "...",
        WebsiteDomain: "...",
    },
    RecipeList: []supertokens.Recipe{
        thirdpartypasswordless.Init(tplmodels.TypeInput{
            FlowType: "...",
            ContactMethodEmailOrPhone: plessmodels.ContactMethodEmailOrPhoneConfig{
                Enabled: true,
            },
            EmailDelivery: &emaildelivery.TypeInput{
                Service: &emaildelivery.EmailDeliveryInterface{
                    SendEmail: &sendEmail,
                },
            },
            SmsDelivery: &smsdelivery.TypeInput{
                Service: &smsdelivery.SmsDeliveryInterface{
                    SendSms: &sendSms,
                },
            },
        }),
    },
})
```

## [0.6.7]
- Fixes panic when call to thirdparty provider API returns a non 2xx status.

### Breaking change
-   https://github.com/supertokens/supertokens-node/issues/220
    -   Adds `{status: "GENERAL_ERROR", message: string}` as a possible output to all the APIs.
    -   Changes `FIELD_ERROR` output status in third party recipe API to be `GENERAL_ERROR`.
    -   Replaced `FIELD_ERROR` status type in third party signinup API with `GENERAL_ERROR`.
    -   Removed `FIELD_ERROR` status type from third party signinup recipe function.
- Changes output of `VerifyEmailPOST` to `VerifyEmailPOSTResponse`
- Changes output of `PasswordResetPOST` to `ResetPasswordPOSTResponse`
- `SignInUp` recipe function doesn't return `FIELD_ERROR` anymore in thirdparty, thirdpartypasswordless and thirdpartyemailpassword recipe.
- `SignInUpPOST` api function returns `GENERAL_ERROR` instead of `FIELD_ERROR` in thirdparty, thirdpartypasswordless and thirdpartyemailpassword recipe.
- If there is an error in sending SMS or email in passwordless based recipes, then we no longer return a GENERAL_ERROR, but instead, we return a regular golang error.
- Changes `GetJWKSGET` in JWT recipe to return `GetJWKSAPIResponse` (that also contains a General Error response)
- Changes `GetOpenIdDiscoveryConfigurationGET` in Open ID recipe to return `GetOpenIdDiscoveryConfigurationAPIResponse` (that also contains a General Error response)
- Renames `OnGeneralError` callback (that's in user input) to `OnSuperTokensAPIError`
- If there is an error in the `errorHandler`, we no longer call `OnSuperTokensAPIError` in that, but instead, we return an error back.

## [0.6.6]
- Fixes facebook login

## [0.6.5]
- Fixes issue in reading request body in API override: https://github.com/supertokens/supertokens-golang/issues/116

## [0.6.4]
- Fixes issue in writing custom response in API override with general error
### Added
- Adds unit tests to thirdpartypasswordless recipe

## [0.6.3] - 2022-05-19
### Fixes
- Fixes the function signature of the `GetUserByThirdPartyInfo` function in the `thirdpartypasswordless` recipe.

## [0.6.2] - 2022-05-18
### Fixes
- Fixes issue in writing custom response in API Override

## [0.6.1] - 2022-05-17
### Fixes
- https://github.com/supertokens/supertokens-golang/issues/102. Sending `preAuthSessionID` instead of `preAuthSessionId` to the core.
- Fixes the error message in AuthorizationUrlAPI function in the `api` module of the thirdparty recipe in case when providers is nil

## [0.6.0] - 2022-05-13
### Breaking Change

- Adds both with context and without context functions to thirdparty passwordless recipe, Like all other recipes. Where we expose both WithContext functions and without context functions, which are basically the same as WithContext ones with an emtpy map[string]interface{} passed as context

### Added
- Adds unit tests to passwordless recipe 

### Fixes
- Fixes existing action to run go mod tidy in the examples folder
- Fixes stopSt function in testing utils

## [0.5.9] - 2022-05-10
### Fixes
- Fixes bug in the revokeCode function of the recipeimplementation in passwordless recipe 

## [0.5.8] - 2022-05-05
### Added
- Adds Github Actions for testing and pre-commit hooks.
- Adds more unit tests for thirdpary email password recipe
- Adds test to jwt recipe
- Adds test to opendID recipe


### Fixes
- Third party sign in up API response correction.

## [0.5.7] - 2022-04-23
- Adds functions to delete passwordless user info in recipes that have passwordless users.
- Fixes bug in signinup helper function exposed by passwordless recipe

## [0.5.6] - 2022-04-18

- Adds UserMetadata recipe

## [0.5.5] - 2022-04-11
### Added 
-   Adds functions for debug logging

## [0.5.4] - 2022-03-30

### Added
 - workflow to enforce go mod tidy is run when issuing a PR. 

## [0.5.3] - 2022-03-24

### Fixes
- Checks if discord returned email before setting it in the profile info obj.

## [0.5.2] - 2022-03-17
- Adds thirdpartypasswordless recipe: https://github.com/supertokens/supertokens-core/issues/331

## [0.5.1] - 2022-02-07

-   Adds testing framework along with unit tests for the recipes
-   Adds unit tests for thirdparty recipe and thirdpartyemailpassword recipe
-   Adds example implementation with go fiber

## [0.5.0] - 2022-02-20
### Breaking Change

-   Adds user context to all functions exposed to the user, and to API and Recipe interface functions. This is a non breaking change for User exposed function calls, but a breaking change if you are using the Recipe or APIs override feature
-   Returns session from API interface functions that create a session
-   Renames functions in ThirdPartyEmailPassword recipe (https://github.com/supertokens/supertokens-node/issues/219):
    -   Recipe Interface:
        -   `SignInUp` -> `ThirdPartySignInUp`
        -   `SignUp` -> `EmailPasswordSignUp`
        -   `SignIn` -> `EmailPasswordSignIn`
    -   API Interface:
        -   `EmailExistsGET` -> `EmailPasswordEmailExistsGET`
    -   User exposed functions (in `recipe/thirdpartyemailpassword/main.go`)
        -   `SignInUp` -> `ThirdPartySignInUp`
        -   `SignUp` -> `EmailPasswordSignUp`
        -   `SignIn` -> `EmailPasswordSignIn`

### Change:

-   Uses recipe interface inside session class so that any modification to those get reflected in the session class functions too.

## [0.4.2] - 2022-01-31
- Adds ability to give a path for each of the hostnames in the connectionURI: https://github.com/supertokens/supertokens-node/issues/252
- Adds workflow to verify if pr title follows conventional commits
- Added userId as an optional property to the response of `recipe/user/password/reset` (Compatibility with CDI 2.12).

### Added

-   Added `regenerateAccessToken` as a new recipe function for the session recipe.
-   Added a bunch of new functions inside the session container which gives user the ability to either call a       function with userContext or just call the function without it (for example: `RevokeSession` and `RevokeSessionWithContext`)
 
### Breaking changes:

-   Allows passing of custom user context everywhere: https://github.com/supertokens/supertokens-golang/issues/64


## [0.4.1] - 2022-01-27
-   Fixes https://github.com/supertokens/supertokens-node/issues/244 - throws an error if a user tries to update email / password of a third party login user.
-   Adds check to see if user has provided empty connectionInfo
-   Adds fixes to solve casting of data in session-functions

## [0.4.0] - 2022-01-14

-   Adds passwordless recipe
-   Adds compatibility with FDI 1.11 and CDI 2.11

## [0.3.5] - 2022-01-08

### Fixes
- Fixes issue of methods getting hidden due to DoneWriter wrapper around ResponseWriter: https://github.com/supertokens/supertokens-golang/issues/55

## [0.3.4] - 2022-01-06

### Fixes
- Sends application/json content-type in `SendNon200Response` function: https://github.com/supertokens/supertokens-golang/issues/53

## [0.3.3] - 2021-12-20

### Added
- Add DeleteUser function

## [0.3.2] - 2021-12-06
### Added
-   The ability to enable JWT creation with session management, this allows easier integration with services that require JWT based authentication: https://github.com/supertokens/supertokens-core/issues/250

## [0.3.1] - 2021-12-06
### Changes
- Upgrade `keyfunc` dependency to stable version.

### Fixes
- Removes use of apiGatewayPath from apple's redirect URI since that is already there in the apiBasePath


## [0.3.0] - 2021-11-23

### Breaking changes:
- Changes `FIELD_ERROR` type in sign in up response from `Error` to `ErrorMsg`

### Addition
- Sign in with google workspaces and discord

### Changes
- If getting profile info from third party provider throws an error, that is propagated a `FIELD_ERROR` to the client.

## [0.2.2] - 2021-11-15

### Changes
- Does not send a response if the user has already sent the response: https://github.com/supertokens/supertokens-node/issues/197

## [0.2.1] - 2021-11-08

### Changes
-   When routing, ignores `rid` value `"anti-csrf"`: https://github.com/supertokens/supertokens-node/issues/202

## [0.2.0] - 2021-10-21

### Breaking changes:
- Makes recipe and API interface have pointers to functions to fix https://github.com/supertokens/supertokens-node/issues/199
-   Support for FDI 1.10:
    -   Allow thirdparty `/signinup POST` API to take `authCodeResponse` XOR `code` so that it can supprt OAuth via PKCE

### Added:
- Makes recipe and API interface have pointers to functions to fix https://github.com/supertokens/supertokens-node/issues/199
-   Support for FDI 1.10:
    -   Adds apple sign in callback API
-   Optional `getRedirectURI` function added to social providers in case we set the `redirect_uri` on the backend.
-   Adds optional `IsDefault` param to auth providers so that they can be reused with different credentials.
- Adds sign in with apple support: https://github.com/supertokens/supertokens-golang/issues/20

## [0.1.0] - 2021-10-21

### Breaking change:

- Removes `SignInUpPost` from thirdpartyemailpassword API interface and replaces it with three APIs: `EmailPasswordSignInPOST`, `EmailPasswordSignUpPOST` and `ThirdPartySignInUpPOST`: https://github.com/supertokens/supertokens-node/issues/192
- Renames all JWT function names to use AccessToken instead for clarity

## [0.0.6] - 2021-10-18

### Changed

-  Changes implementation such that actual client IDs are not in the SDK, removes imports for OAuth dev related code.

## [0.0.5] - 2021-10-18

### Fixed

- URL protocol is being taken into account when determining the value of cookie same site: https://github.com/supertokens/supertokens-golang/issues/36

## [0.0.4] - 2021-10-12

### Added

- Adds OAuth development keys for Google and Github for faster recipe implementation.

## [0.0.3] - 2021-09-25

### Added

- Support for FDI 1.9
- JWT Recipe

### Fixed
- Sets response content-type as JSON

## [0.0.2] - 2021-09-22

### Added

-   Support for multiple access token signing keys: https://github.com/supertokens/supertokens-core/issues/305
-   Supporting CDI 2.9

## [0.0.1] - 2021-09-18

### Added
- Initial version of the repo