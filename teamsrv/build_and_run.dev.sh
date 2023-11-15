# type Config struct {
# 	App struct {
# 		Version string
# 		Source  struct {
# 			Url string
# 		}
# 	}
# 	OpenApi struct {
# 		AuthorizationURL string
# 		TokenURL         string
# 		Scopes           map[string]string
# 		OpenIdConnectUrl string
# 	}
# }

echo "Building the app"
env GOOS=linux CGO_ENABLED=0 go build -o /go/bin/teamsrv ./cmd

echo "Setting env variables"
#These are set in the container at build time
export App__Version="0.0.0.0"
export App__Source__Url="https://github.com/Miruken-Go/demo.microservice"

#These are set at deployment time
export Login__OAuth__0__Module="login.jwt"
export Login__OAuth__0__Options__Audience="07574dda-f3b0-4fed-aa9a-2e041b6ad3d1"
export Login__OAuth__0__Options__JWKS__Uri="https://teamsrvidentitydev.b2clogin.com/teamsrvidentitydev.onmicrosoft.com/discovery/v2.0/keys?p=b2c_1a_signup_signin"
export Login__Adb2c__0__Module="login.pwd"
export Login__Adb2c__0__Options__Credentials__0__Username="ooYymDzee5!V&v8gk7*s"
export Login__Adb2c__0__Options__Credentials__0__Password="i**72R#PLWbx8&#$I$ok"
export OpenApi__AuthorizationUrl="https://teamsrvidentitydev.b2clogin.com/teamsrvidentitydev.onmicrosoft.com/oauth2/v2.0/authorize?p=b2c_1a_signup_signin"
export OpenApi__TokenURL="https://teamsrvidentitydev.b2clogin.com/teamsrvidentitydev.onmicrosoft.com/oauth2/v2.0/token?p=b2c_1a_signup_signin"
export OpenApi__ClientId="3d8bd886-f1a7-42be-9319-acdf39a7673b"
export OpenApi__OpenIdConnectUrl="https://teamsrvidentitydev.b2clogin.com/teamsrvidentitydev.onmicrosoft.com/v2.0/.well-known/openid-configuration?p=B2C_1A_SIGNUP_SIGNIN"
export OpenApi__Scopes__0__Name="https://teamsrvidentitydev.onmicrosoft.com/teamsrv/Groups"
export OpenApi__Scopes__0__Description="Groups to which the user belongs."
export OpenApi__Scopes__1__Name="https://teamsrvidentitydev.onmicrosoft.com/teamsrv/Roles"
export OpenApi__Scopes__1__Description="Roles to which the user belongs."
export OpenApi__Scopes__2__Name="https://teamsrvidentitydev.onmicrosoft.com/teamsrv/Entitlements"
export OpenApi__Scopes__2__Description="Entitlements the user has."

echo "Starting the app: localhost:8080"
/go/bin/teamsrv