module github.com/Azure/sonic-telemetry

go 1.12

require (
	github.com/Azure/sonic-mgmt-common v0.0.0-00010101000000-000000000000
	github.com/Workiva/go-datastructures v1.0.52
	github.com/bgentry/speakeasy v0.1.0 // indirect
	github.com/c9s/goprocinfo v0.0.0-20191125144613-4acdd056c72d
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/go-redis/redis v6.15.7+incompatible
	github.com/go-redis/redis/v7 v7.2.0
	github.com/gogo/protobuf v1.3.1
	github.com/golang/glog v0.0.0-20160126235308-23def4e6c14b
	github.com/golang/protobuf v1.4.0
	github.com/google/gnxi v0.0.0-20191016182648-6697a080bc2d
	github.com/grpc-ecosystem/go-grpc-middleware v1.2.0
	github.com/jipanyang/gnmi v0.0.0-20180820232453-cb4d464fa018
	github.com/jipanyang/gnxi v0.0.0-20181221084354-f0a90cca6fd0
	github.com/kylelemons/godebug v1.1.0
	github.com/msteinert/pam v0.0.0-20190215180659-f29b9f28d6f9
	github.com/openconfig/gnmi v0.0.0-20200307010808-e7106f7f5493
	github.com/openconfig/gnoi v0.0.0-20191206155121-b4d663a26026
	github.com/openconfig/ygot v0.7.1
	github.com/xeipuuv/gojsonschema v1.2.0
	golang.org/x/crypto v0.0.0-20200622213623-75b288015ac9
	golang.org/x/net v0.0.0-20200813134508-3edf25e44fcc
	google.golang.org/grpc v1.28.0
	gopkg.in/yaml.v2 v2.2.4
)

replace github.com/Azure/sonic-mgmt-common => ../sonic-mgmt-common
