package gnmi_server

import (
	"context"
	gnoi_system_pb "github.com/openconfig/gnoi/system"
	log "github.com/golang/glog"
	"time"
	spb "proto/gnoi"
	transutil "transl_utils"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/codes"
	"encoding/json"
	"strings"
	"fmt"
	jwt "github.com/dgrijalva/jwt-go"
)

func (srv *Server) Reboot(ctx context.Context, req *gnoi_system_pb.RebootRequest) (*gnoi_system_pb.RebootResponse, error) {
	err := authenticate(srv.config.UserAuth, ctx, false)
	if err != nil {
		return nil, err
	}
	log.V(1).Info("gNOI: Reboot")
	return nil, nil
}
func (srv *Server) RebootStatus(ctx context.Context, req *gnoi_system_pb.RebootStatusRequest) (*gnoi_system_pb.RebootStatusResponse, error) {
	err := authenticate(srv.config.UserAuth, ctx, false)
	if err != nil {
		return nil, err
	}
	log.V(1).Info("gNOI: RebootStatus")
	return nil, nil
}
func (srv *Server) CancelReboot(ctx context.Context, req *gnoi_system_pb.CancelRebootRequest) (*gnoi_system_pb.CancelRebootResponse, error) {
	err := authenticate(srv.config.UserAuth, ctx, false)
	if err != nil {
		return nil, err
	}
	log.V(1).Info("gNOI: CancelReboot")
	return nil, nil
}
func (srv *Server) Ping(req *gnoi_system_pb.PingRequest, rs gnoi_system_pb.System_PingServer) error {
	ctx := rs.Context()
	err := authenticate(srv.config.UserAuth, ctx, false)
	if err != nil {
		return err
	}
	log.V(1).Info("gNOI: Ping")
	return nil
}
func (srv *Server) Traceroute(req *gnoi_system_pb.TracerouteRequest, rs gnoi_system_pb.System_TracerouteServer) error {
	ctx := rs.Context()
	err := authenticate(srv.config.UserAuth, ctx, false)
	if err != nil {
		return err
	}
	log.V(1).Info("gNOI: Traceroute")
	return nil
}
func (srv *Server) SetPackage(rs gnoi_system_pb.System_SetPackageServer) error {
	ctx := rs.Context()
	err := authenticate(srv.config.UserAuth, ctx, false)
	if err != nil {
		return err
	}
	log.V(1).Info("gNOI: SetPackage")
	return nil
}
func (srv *Server) SwitchControlProcessor(ctx context.Context, req *gnoi_system_pb.SwitchControlProcessorRequest) (*gnoi_system_pb.SwitchControlProcessorResponse, error) {
	err := authenticate(srv.config.UserAuth, ctx, false)
	if err != nil {
		return nil, err
	}
	log.V(1).Info("gNOI: SwitchControlProcessor")
	return nil, nil
}
func (srv *Server) Time(ctx context.Context, req *gnoi_system_pb.TimeRequest) (*gnoi_system_pb.TimeResponse, error) {
	err := authenticate(srv.config.UserAuth, ctx, false)
	if err != nil {
		return nil, err
	}
	log.V(1).Info("gNOI: Time")
	var tm gnoi_system_pb.TimeResponse
	tm.Time = uint64(time.Now().UnixNano())
	return &tm, nil
}


func (srv *Server) ShowTechsupport(ctx context.Context, req *spb.TechsupportRequest) (*spb.TechsupportResponse, error) {
	err := authenticate(srv.config.UserAuth, ctx, false)
	if err != nil {
		return nil, err
	}
	log.V(1).Info("gNOI: Sonic ShowTechsupport")
	var jobj map[string]map[string]string
	resp := &spb.TechsupportResponse{
		Output: &spb.TechsupportResponse_Output {

		},
	}
	reqstr := fmt.Sprintf("{\"sonic-show-techsupport-info:input\": {\"date\": \"%s\"}}", req.Input.Date)
	jsresp, err:= transutil.TranslProcessAction("/sonic-show-techsupport:sonic-show-techsupport-info", []byte(reqstr))
	if err != nil {
		return nil, status.Error(codes.Unknown, err.Error())
	}
	jsresp = []byte(strings.Replace(string(jsresp), "sonic-show-techsupport:output", "output", 1))
	err = json.Unmarshal(jsresp, &jobj)
	if err != nil {
		return nil, status.Error(codes.Unknown, err.Error())
	}
	resp.Output.OutputFilename = jobj["output"]["output-filename"]
	fmt.Println(resp)
	return resp, nil
}

func (srv *Server) Sum(ctx context.Context, req *spb.SumRequest) (*spb.SumResponse, error) {
	err := authenticate(srv.config.UserAuth, ctx, false)
	if err != nil {
		return nil, err
	}
	log.V(1).Info("gNOI: Sonic Sum")
	var resp spb.SumResponse
	reqstr := fmt.Sprintf("{\"sonic-tests:input\": {\"left\": %d, \"right\": %d}}", req.Input.Left, req.Input.Right)
	jsresp, err:= transutil.TranslProcessAction("/sonic-tests:sum", []byte(reqstr))
	
	if err != nil {
		return nil, status.Error(codes.Unknown, err.Error())
	}
	jsresp = []byte(strings.Replace(string(jsresp), "sonic-tests:output", "output", 1))
	err = json.Unmarshal(jsresp, &resp)
	if err != nil {
		return nil, status.Error(codes.Unknown, err.Error())
	}
	return &resp, nil
}

func (srv *Server) SaveConfig(ctx context.Context, req *spb.SaveConfigRequest) (*spb.SaveConfigResponse, error) {
	err := authenticate(srv.config.UserAuth, ctx, false)
	if err != nil {
		return nil, err
	}
	log.V(1).Info("gNOI: Sonic SaveConfig")
	var resp spb.SaveConfigResponse
	reqstr := fmt.Sprintf("{\"save_config:input\": {\"file_path\": \"%s\"}", req.Input.FilePath)
	jsresp, err:= transutil.TranslProcessAction("/sonic-config-mgmt:save_config", []byte(reqstr))
	if err != nil {
		return nil, status.Error(codes.Unknown, err.Error())
	}

	jsresp = []byte(strings.Replace(string(jsresp), "save_config:output", "output", 1))
	err = json.Unmarshal(jsresp, &resp)
	if err != nil {
		return nil, status.Error(codes.Unknown, err.Error())
	}
	return &resp, nil
}

func (srv *Server) ReloadConfig(ctx context.Context, req *spb.ReloadConfigRequest) (*spb.ReloadConfigResponse, error) {
	err := authenticate(srv.config.UserAuth, ctx, false)
	if err != nil {
		return nil, err
	}
	log.V(1).Info("gNOI: Sonic ReloadConfig")
	var resp spb.ReloadConfigResponse
	reqstr := fmt.Sprintf("{\"reload_config:input\": {\"file_path\": \"%s\"}", req.Input.FilePath)
	jsresp, err:= transutil.TranslProcessAction("/sonic-config-mgmt:reload_config", []byte(reqstr))
	if err != nil {
		return nil, status.Error(codes.Unknown, err.Error())
	}

	jsresp = []byte(strings.Replace(string(jsresp), "reload_config:output", "output", 1))
	err = json.Unmarshal(jsresp, &resp)
	if err != nil {
		return nil, status.Error(codes.Unknown, err.Error())
	}
	return &resp, nil
}

func (srv *Server) LoadMgmtConfig(ctx context.Context, req *spb.LoadMgmtConfigRequest) (*spb.LoadMgmtConfigResponse, error) {
	err := authenticate(srv.config.UserAuth, ctx, false)
	if err != nil {
		return nil, err
	}
	log.V(1).Info("gNOI: Sonic LoadMgmtConfig")
	var resp spb.LoadMgmtConfigResponse
	reqstr := fmt.Sprintf("{\"load_mgmt_config:input\": {\"file_path\": \"%s\"}", req.Input.FilePath)
	jsresp, err:= transutil.TranslProcessAction("/sonic-config-mgmt:load_mgmt_config", []byte(reqstr))
	if err != nil {
		return nil, status.Error(codes.Unknown, err.Error())
	}

	jsresp = []byte(strings.Replace(string(jsresp), "load_mgmt_config:output", "output", 1))
	err = json.Unmarshal(jsresp, &resp)
	if err != nil {
		return nil, status.Error(codes.Unknown, err.Error())
	}
	return &resp, nil
}

func (srv *Server) LoadMinigraph(ctx context.Context, req *spb.LoadMinigraphRequest) (*spb.LoadMinigraphResponse, error) {
	err := authenticate(srv.config.UserAuth, ctx, false)
	if err != nil {
		return nil, err
	}
	log.V(1).Info("gNOI: Sonic LoadMinigraph")
	var resp spb.LoadMinigraphResponse
	reqstr := fmt.Sprintf("{\"load_minigraph:input\": {\"file_path\": \"%s\"}", req.Input.FilePath)
	jsresp, err:= transutil.TranslProcessAction("/sonic-config-mgmt:load_minigraph", []byte(reqstr))
	if err != nil {
		return nil, status.Error(codes.Unknown, err.Error())
	}

	jsresp = []byte(strings.Replace(string(jsresp), "load_minigraph:output", "output", 1))
	err = json.Unmarshal(jsresp, &resp)
	if err != nil {
		return nil, status.Error(codes.Unknown, err.Error())
	}
	return &resp, nil
}

func (srv *Server) Authenticate(ctx context.Context, req *spb.AuthenticateRequest) (*spb.AuthenticateResponse, error) {
	// Can't enforce normal authentication here.. maybe only enforce client cert auth if enabled?
	// err := authenticate(srv.config.UserAuth, ctx, false)
	// if err != nil {
	// 	return nil, err
	// }
	log.V(1).Info("gNOI: Sonic Authenticate")



	auth_success, _ := UserPwAuth(req.Username, req.Password)
	if  auth_success {
		return &spb.AuthenticateResponse{Token: tokenResp(req.Username)}, nil
	}
	return nil, status.Errorf(codes.PermissionDenied, "Invalid Username or Password")	

}
func (srv *Server) Refresh(ctx context.Context, req *spb.RefreshRequest) (*spb.RefreshResponse, error) {
	err := authenticate(srv.config.UserAuth, ctx, false)
	if err != nil {
		return nil, err
	}
	log.V(1).Info("gNOI: Sonic Refresh")
	
	token, err := JwtAuthenAndAuthor(ctx, false)
	if err != nil {
		return nil, err
	}

	claims := &Claims{}
	jwt.ParseWithClaims(token.AccessToken, claims, func(token *jwt.Token) (interface{}, error) {
		return hmacSampleSecret, nil
	})
	if time.Unix(claims.ExpiresAt, 0).Sub(time.Now()) > JwtRefreshInt {
		return nil, status.Errorf(codes.InvalidArgument, "Invalid JWT Token")
	}
	
	return &spb.RefreshResponse{Token: tokenResp(claims.Username)}, nil

}