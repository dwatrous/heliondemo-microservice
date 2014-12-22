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
Extract submitted surveys http://localhost:8080/result/dog for all results of the dog survey

# MongoDB Required
The MongoDB requirement can be satisfied by creating another VM using the Ubunutu 14.04 image and running the following commands.

```
#export http_proxy=http://proxy.company.com:8080
#export https_proxy=https://proxy.company.com:8080
sudo -E apt-key adv --keyserver hkp://keyserver.ubuntu.com:80 --recv 7F0CEB10
echo 'deb http://downloads-distro.mongodb.org/repo/ubuntu-upstart dist 10gen' | sudo tee /etc/apt/sources.list.d/mongodb.list
sudo -E apt-get update
sudo -E apt-get -y install mongodb-org
```

At this point MongoDB is running, but by default it will only accept connections from localhost. To change this, it's necessary to edit the file **/etc/mongod.conf**. Update the line which declares **bind_ip** by either commenting it out or changing the value to **0.0.0.0**. After making the change, restart the service.

```
sudo service mongod restart
```

For convenience, allocate a floating IP address and associate it with the MongoDB VM so that this can be injected into the deployed instance.

# Deploy to Helion

### Setup proxy
If your HP Helion Development Platform installation requires a proxy to access the internet, the following kato command can be run from a DEA or the Cloud Controller:

```
helion@watroushdp-dea-2:~$ kato op upstream_proxy set proxy.company.com:8080

Upstream proxy set to: proxy.houston.hp.com:8080

Restart your DEA nodes to effect the changes

Polipo must be restarted for it to take effect:
        sudo /etc/init.d/polipo restart

helion@watroushdp-dea-2:~$ sudo /etc/init.d/polipo restart
Restarting polipo: polipo.
```

### Add MongoDB as a user provided service
Use the **create-service** command to add MnogoDB credentials and connection details as a service with the name **mongo-heliondemo**

```
ubuntu@hdp-install:~$ helion create-service user-provided mongo-heliondemo
Which credentials to use for connections [hostname, port, password]: host, port, database
host: 15.50.137.94
port: 27017
database: survey
Creating new service ... OK
```

### Push the code
Change into the cloned repository run *helion push*. When prompted to Bind existing services, choose "Y". Select the user provided service from the list to bind it to this application. 

```
ubuntu@hdp-install:~/heliondemo-microservice$ helion push
Would you like to deploy from the current directory ? [Yn]:
Using manifest file "manifest.yml"
Application Deployed URL [survey.15.50.137.82.xip.io]:
Application Url:   http://survey.15.50.137.82.xip.io
Enter GOVERSION [1.2]:
  Adding Environment Variable [GOVERSION=1.2]
Creating Application [survey] as [https://api.15.50.137.82.xip.io -> default -> default -> survey] ... OK
  Map http://survey.15.50.137.82.xip.io ... OK
Bind existing services to 'survey' ? [yN]: y
Which one ?
1. mongo-heliondemo
Choose:? 1
  Binding mongo-heliondemo to survey ... OK
Bind another ? [yN]:
Create services to bind to 'survey' ? [yN]:
Uploading Application [survey] ...
  Checking for bad links ... 18 OK
  Copying to temp space ... 17 OK
  Checking for available resources ...  OK
  Processing resources ... OK
  Packing application ... OK
  Uploading (1K) ... 100% OK
Push Status: OK
Starting Application [survey] ...
stackato[dea_ng]: Staging application
staging: -----> Downloaded app package (9.4M)
staging: -----> Installing Go 1.2...
staging: -----> Running: go get -tags heroku ./...
staging: -----> Uploading droplet (12M)
stackato[dea_ng]: Completed staging application
stackato[dea_ng.0]: Launching web process: bin/survey
app[stdout.0]: Listening on :46829
OK
http://survey.15.50.137.82.xip.io/ deployed
```

### Verify the service is running

Assuming the above output, the service can be verified as follows.

Render a survey http://survey.15.50.137.82.xip.io/survey/dog which corresponds to survey/dog.json

Extract submitted surveys http://survey.15.50.137.82.xip.io/result/dog for all results of the dog survey
