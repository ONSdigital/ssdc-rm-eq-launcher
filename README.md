# RM EQ Launcher

This project is based on https://github.com/ONSdigital/ssdc-rm-eq-launcher

### Building and Running
Install Go and ensure that your `GOPATH` env variable is set (usually it's `~/go`).

```
go get
go build
./ssdc-rm-eq-launcher

go run launch.go (Does both the build and run cmd above)
```

Open http://localhost:8000/

### Docker
The dockerfile is a multistage dockerfile which can be built using:

```
docker build -t ssdc-rm-eq-launcher:latest .
```

You can then run the image using `SURVEY_RUNNER_SCHEMA_URL` to point it at an instance of survey runner.

```
docker run -e SURVEY_RUNNER_SCHEMA_URL=http://localhost:5000 -it -p 8000:8000 onsdigital/ssdc-rm-eq-launcher:latest
```

The syntax for this will be slightly different on Mac

```
docker run -e SURVEY_RUNNER_SCHEMA_URL=http://host.docker.internal:5000 -it -p 8000:8000 onsdigital/ssdc-rm-eq-launcher:latest
```

You should then be able to access go launcher at `localhost:8000`

### Run

First, start the launcher running like this: `SURVEY_RUNNER_URL=https://test-runner.eq.gcp.onsdigital.uk CASE_API_URL=http://localhost:8161 ./ssdc-rm-eq-launcher`

Now navigate to http://localhost:8000/

### Settings
Environment Variable | Meaning | Default
---------------------|---------|--------
GO_LAUNCH_A_SURVEY_LISTEN_HOST|Host address  to listen on|0.0.0.0
GO_LAUNCH_A_SURVEY_LISTEN_PORT|Host port to listen on|8000
SURVEY_RUNNER_URL|URL of Survey Runner to re-direct to when launching a survey|http://localhost:5000
SURVEY_REGISTER_URL|URL of eq-survey-register to load schema list from |http://localhost:8080
JWT_ENCRYPTION_KEY_PATH|Path to the JWT Encryption Key (PEM format)|jwt-test-keys/sdc-user-authentication-encryption-sr-public-key.pem
JWT_SIGNING_KEY_PATH|Path to the JWT Signing Key (PEM format)|jwt-test-keys/sdc-user-authentication-signing-launcher-private-key.pem
