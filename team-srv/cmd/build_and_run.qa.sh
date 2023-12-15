echo 'Building the app'
env GOOS=linux CGO_ENABLED=0 go build -o /go/bin/team-srv ./cmd

echo 'Setting env variables'
#These are set in the container at build time
export App__Version='0.0.0.0'
export App__Source__Url='https://github.com/Miruken-Go/demo.microservice'

#These are set at deployment time
export Login__OAuth__0__Module='login.jwt'
export Login__OAuth__0__Options__JWKS__Uri='https://teamsrvidentityqa.b2clogin.com/teamsrvidentityqa.onmicrosoft.com/discovery/v2.0/keys?p=b2c_1a_signup_signin'
export Login__Basic__0__Module='login.pwd'
export Login__Basic__0__Options__Credentials__0__Username='ooYymDzee5!V&v8gk7*s'
export Login__Basic__0__Options__Credentials__0__Password='i**72R#PLWbx8&#$I$ok'
export OpenApi__AuthorizationUrl='https://teamsrvidentityqa.b2clogin.com/teamsrvidentityqa.onmicrosoft.com/oauth2/v2.0/authorize?p=b2c_1a_signup_signin'
export OpenApi__TokenURL='https://teamsrvidentityqa.b2clogin.com/teamsrvidentityqa.onmicrosoft.com/oauth2/v2.0/token?p=b2c_1a_signup_signin'
export OpenApi__ClientId='c60feecc-e6fe-43cf-aab3-f33a3ab987de'
export OpenApi__OpenIdConnectUrl='https://teamsrvidentityqa.b2clogin.com/teamsrvidentityqa.onmicrosoft.com/v2.0/.well-known/openid-configuration?p=B2C_1A_SIGNUP_SIGNIN'
export OpenApi__Scopes__0__Name='https://teamsrvidentityqa.onmicrosoft.com/teamsrv/Group'
export OpenApi__Scopes__0__Description='Groups to which the user belongs.'
export OpenApi__Scopes__1__Name='https://teamsrvidentityqa.onmicrosoft.com/teamsrv/Role'
export OpenApi__Scopes__1__Description='Roles to which the user belongs.'
export OpenApi__Scopes__2__Name='https://teamsrvidentityqa.onmicrosoft.com/teamsrv/Entitlement'
export OpenApi__Scopes__2__Description='Entitlements the user has.'

echo 'Starting the app: localhost:8080'
/go/bin/team-srv