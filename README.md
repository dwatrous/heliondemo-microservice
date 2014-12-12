heliondemo-microservice
=======================

HelionDemo microservice example

# need some cloundfoundry env variables set up. not all the values are used.
export VCAP_SERVICES='{"user-provided":[{"name":"mongo-heliondemo","label":"user-provided","tags":[],"credentials":{"host":"mongo.heliondemo.com","port":"27000","database":"surveys"},"syslog_drain_url":""}]}'
export VCAP_APPLICATION='{"instance_id":"752711c063b544ddab0e9349e228e194","instance_index":0,"host":"0.0.0.0","port":53335,"started_at":"2014-12-06 19:50:41 -0800","started_at_timestamp":1417924241,"start":"2014-12-06 19:50:41 -0800","state_timestamp":1417924241,"limits":{"mem":128,"disk":2048,"fds":16384,"allow_sudo":false},"application_version":"d020c6b6-2128-43d4-864d-62de17770bb5","application_name":"go-env-dan","application_uris":["go-env-dan.15.126.246.248.xip.io"],"sso_enabled":false,"version":"d020c6b6-2128-43d4-864d-62de17770bb5","name":"go-env-dan","uris":["go-env-dan.15.126.246.248.xip.io"],"users":null}'
# set port for dev
export PORT=8080

go get github.com/dwatrous/heliondemo-microservice
cd $GOPATH/src/github.com/dwatrous/heliondemo-microservice
go build

./heliondemo-microservice

# now open http://localhost:8080/survey/dog which corresponds to survey/dogs.json in the repo
# submit form
# visit http://localhost:8080/result/dog for all results of dogs survey

