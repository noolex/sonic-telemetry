/* Copyright 2017 Google Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    https://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Binary gnmi_set performs a set request against a gNMI target with the specified config file.
package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
	"time"

	log "github.com/golang/glog"
	"golang.org/x/net/context"
	"google.golang.org/grpc"

	"github.com/google/gnxi/utils"
	"github.com/google/gnxi/utils/credentials"
	//"github.com/google/gnxi/utils/xpath"
	"github.com/jipanyang/gnxi/utils/xpath"
	"google.golang.org/grpc/metadata"
	"github.com/golang/protobuf/proto"
	pb "github.com/openconfig/gnmi/proto/gnmi"
	ext_pb "github.com/openconfig/gnmi/proto/gnmi_ext"
	spb "proto"
)

type arrayFlags []string

func (i *arrayFlags) String() string {
	return "my string representation"
}

func (i *arrayFlags) Set(value string) error {
	*i = append(*i, value)
	return nil
}

var (
	deleteOpt  arrayFlags
	replaceOpt arrayFlags
	updateOpt  arrayFlags
	targetAddr = flag.String("target_addr", "localhost:10161", "The target address in the format of host:port")
	targetName = flag.String("target_name", "hostname.com", "The target name use to verify the hostname returned by TLS handshake")
	timeOut    = flag.Duration("time_out", 10*time.Second, "Timeout for the Get request, 10 seconds by default")
	pathTarget = flag.String("xpath_target", "", "name of the target for which the path is a member")
	jwtToken   = flag.String("jwt_token", "", "JWT Token if required")
	bundleVersion    = flag.String("bundle_ver", "", "Optional version specifier for model bundle version.")
)

func buildPbUpdateList(pathValuePairs []string) []*pb.Update {
	var pbUpdateList []*pb.Update
	for _, item := range pathValuePairs {
		modName := strings.SplitN(item, "/", 3)
		
		pathValuePair := make([]string, 2)
		lc := strings.LastIndex(modName[2],":")
		
		if lc == -1 {
			log.Exitf("invalid path-value pair: %v", item)
		}
		pathValuePair_r := strings.SplitN(modName[2], ":", 2)
		// pathValuePair[0] = modName[2][0:lc]

		pathValuePair[0] = "/" + modName[1] + "/" + pathValuePair_r[0]

		pathValuePair[1] = modName[2][lc+1:]
		fmt.Println(pathValuePair[0])
		fmt.Println(pathValuePair[1])


		if len(pathValuePair) != 2 || len(pathValuePair[1]) == 0 {
			log.Exitf("invalid path-value pair: %v", item)
			log.Exitf("invalid path-value pair: %v", modName)
		}


		pbPath, err := xpath.ToGNMIPath(pathValuePair[0])
		if err != nil {
			log.Exitf("error in parsing xpath %q to gnmi path", pathValuePair[0])
		}
		var pbVal *pb.TypedValue
		if pathValuePair[1][0] == '@' {
			jsonFile := pathValuePair[1][1:]
			jsonConfig, err := ioutil.ReadFile(jsonFile)
			if err != nil {
				log.Exitf("cannot read data from file %v", jsonFile)
			}
			jsonConfig = bytes.Trim(jsonConfig, " \r\n\t")
			pbVal = &pb.TypedValue{
				Value: &pb.TypedValue_JsonIetfVal{
					JsonIetfVal: jsonConfig,
				},
			}
		} else {
			if strVal, err := strconv.Unquote(pathValuePair[1]); err == nil {
				pbVal = &pb.TypedValue{
					Value: &pb.TypedValue_StringVal{
						StringVal: strVal,
					},
				}
			} else {
				if intVal, err := strconv.ParseInt(pathValuePair[1], 10, 64); err == nil {
					pbVal = &pb.TypedValue{
						Value: &pb.TypedValue_IntVal{
							IntVal: intVal,
						},
					}
				} else if floatVal, err := strconv.ParseFloat(pathValuePair[1], 32); err == nil {
					pbVal = &pb.TypedValue{
						Value: &pb.TypedValue_FloatVal{
							FloatVal: float32(floatVal),
						},
					}
				} else if boolVal, err := strconv.ParseBool(pathValuePair[1]); err == nil {
					pbVal = &pb.TypedValue{
						Value: &pb.TypedValue_BoolVal{
							BoolVal: boolVal,
						},
					}
				} else {
					pbVal = &pb.TypedValue{
						Value: &pb.TypedValue_StringVal{
							StringVal: pathValuePair[1],
						},
					}
				}
			}
		}
		pbUpdateList = append(pbUpdateList, &pb.Update{Path: pbPath, Val: pbVal})
	}
	return pbUpdateList
}

func main() {
	flag.Var(&deleteOpt, "delete", "xpath to be deleted.")
	flag.Var(&replaceOpt, "replace", "xpath:value pair to be replaced. Value can be numeric, boolean, string, or IETF JSON file (. starts with '@').")
	flag.Var(&updateOpt, "update", "xpath:value pair to be updated. Value can be numeric, boolean, string, or IETF JSON file (. starts with '@').")
	flag.Parse()

	opts := credentials.ClientCredentials(*targetName)
	conn, err := grpc.Dial(*targetAddr, opts...)
	if err != nil {
		log.Exitf("Dialing to %q failed: %v", *targetAddr, err)
	}
	defer conn.Close()

	var deleteList []*pb.Path
	for _, xPath := range deleteOpt {
		pbPath, err := xpath.ToGNMIPath(xPath)
		if err != nil {
			log.Exitf("error in parsing xpath %q to gnmi path", xPath)
		}
		deleteList = append(deleteList, pbPath)
	}
	replaceList := buildPbUpdateList(replaceOpt)
	updateList := buildPbUpdateList(updateOpt)
	var prefix pb.Path
	prefix.Target = *pathTarget
	setRequest := &pb.SetRequest{
		Prefix:    &prefix,
		Delete:  deleteList,
		Replace: replaceList,
		Update:  updateList,
	}
	if len(*bundleVersion) > 0 {
		bv, err := proto.Marshal(&spb.BundleVersion{
			Version: *bundleVersion,
		})
		if err != nil {
			log.Exitf("%v", err)
		}

		setRequest.Extension = append(setRequest.Extension, &ext_pb.Extension{
			Ext: &ext_pb.Extension_RegisteredExt {
				RegisteredExt: &ext_pb.RegisteredExtension {
				Id: 999,
				Msg: bv,
			}}})
	}
	fmt.Println("== setRequest:")
	utils.PrintProto(setRequest)

	cli := pb.NewGNMIClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), *timeOut)
	defer cancel()

	if len(*jwtToken) > 0 {
		ctx = metadata.AppendToOutgoingContext(ctx, "access_token", *jwtToken)
	}
	setResponse, err := cli.Set(ctx, setRequest)
	if err != nil {
		log.Exitf("Set failed: %v", err)
	}

	fmt.Println("== setResponse:")
	utils.PrintProto(setResponse)
}
