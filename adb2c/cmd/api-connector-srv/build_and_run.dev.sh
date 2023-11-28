echo 'Building the app'
env GOOS=linux CGO_ENABLED=0 go build -o /go/bin/api-connector-srv /go/src/app/adb2c/cmd/api-connector-srv

echo 'Setting env variables'
#These are set in the container at build time
export App__Version='0.0.0.0'
export App__Source__Url='https://github.com/Miruken-Go/demo.microservice'

#These are set at deployment time
export Login__Adb2c__0__Module='login.pwd'
export Login__Adb2c__0__Options__Credentials__0__Username='ooYymDzee5!V&v8gk7*s'
export Login__Adb2c__0__Options__Credentials__0__Password='i**72R#PLWbx8&#$I$ok'

env

echo 'Starting the app: localhost:8080'
/go/bin/api-connector-srv