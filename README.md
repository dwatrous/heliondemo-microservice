HP Helion microservice demo app
=======================

This microservice implements a survey feature that can be embedded into a legacy application. This is an example app that can be pushed directly to Helion Development Platform. A MongoDB instance is required and can be created in a separate VM through Helion OpenStack. *This code is meant to be a starting point, not a production ready service.*

# Local Development

### Clone and build the repository
Golang, bzr and git are required to run the commands below. On Ubuntu they can be installed with the following commands.

```
#export http_proxy=http://proxy.company.com:8080
#export https_proxy=https://proxy.company.com:8080
sudo -E apt-get update
sudo -E apt-get install -y git golang bzr
```

The code can be cloned and built locally using the following commands. 

```
export GOPATH=~/go
go get github.com/dwatrous/heliondemo-microservice
cd $GOPATH/src/github.com/dwatrous/heliondemo-microservice
go build
```

### Set environment variables to simulate Development Platform instance environment.

Environment variables **VCAP_SERVICES**, **VCAP_APPLICATION** and **PORT** are used by the Go application to connect to MongoDB and expose the web service on the specified port. The following *export* commands can be modified to match your environment.

 * "**mongo-heliondemo**" is the name given to the MongoDB service when you add it to Helion
 * the "**credentials**" element should be updated to point to your mogno hsot and port

```
export VCAP_SERVICES='{"user-provided":[{"name":"mongo-heliondemo","label":"user-provided","tags":[],"credentials":{"host":"mongo.heliondemo.com","port":"27000","database":"surveys"},"syslog_drain_url":""}]}'
export VCAP_APPLICATION='{"instance_id":"752711c063b544ddab0e9349e228e194","instance_index":0,"host":"0.0.0.0","port":53335,"started_at":"2014-12-06 19:50:41 -0800","started_at_timestamp":1417924241,"start":"2014-12-06 19:50:41 -0800","state_timestamp":1417924241,"limits":{"mem":128,"disk":2048,"fds":16384,"allow_sudo":false},"application_version":"d020c6b6-2128-43d4-864d-62de17770bb5","application_name":"go-env-dan","application_uris":["go-env-dan.15.126.246.248.xip.io"],"sso_enabled":false,"version":"d020c6b6-2128-43d4-864d-62de17770bb5","name":"go-env-dan","uris":["go-env-dan.15.126.246.248.xip.io"],"users":null}'
# set port for dev
export PORT=8080
```

### Run the service
Run the service from the build directory

```
$GOPATH/src/github.com/dwatrous/heliondemo-microservice/heliondemo-microservice
```

### Verify the service is running

Render a survey http://localhost:8080/survey/dog which corresponds to survey/dog.json
Extract submitted surveys http://localhost:8080/result/dog for all results of dogs survey

# Deploy to Helion

# MongoDB Required