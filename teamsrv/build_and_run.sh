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

env GOOS=linux CGO_ENABLED=0 go build -o /go/bin/teamsrv ./cmd

export App__Version=0.0.0.0
export App__Source__Url=https://github.com/Miruken-Go/demo.microservice
export OpenApi__AuthorizationUrl=https://teamsrvidentitydev.b2clogin.com/teamsrvidentitydev.onmicrosoft.com/oauth2/v2.0/authorize?p=b2c_1a_signup_signin
export OpenApi__TokenURL=https://teamsrvidentitydev.b2clogin.com/teamsrvidentitydev.onmicrosoft.com/oauth2/v2.0/token?p=b2c_1a_signup_signin
export OpenApi__Scopes__0__Name=https://teamsrvidentitydev.onmicrosoft.com/teamsrv/Groups
export OpenApi__Scopes__0__Description="Groups to which the user belongs."
export OpenApi__Scopes__1__Name="https://teamsrvidentitydev.onmicrosoft.com/teamsrv/Roles"
export OpenApi__Scopes__1__Description="Roles to which the user belongs."
export OpenApi__Scopes__2__Name=https://teamsrvidentitydev.onmicrosoft.com/teamsrv/Entitlements 
export OpenApi__Scopes__2__Description="Entitlements the user has."

/go/bin/teamsrv