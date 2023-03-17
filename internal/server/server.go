package main

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net"
	"net/http"

	infov1 "github.com/MatthewSerre/hyundai-bluelink-protobufs/gen/go/protos/information/v1"
	o "github.com/MatthewSerre/hyundai-bluelink-vehicle-information-service/internal/owner_info_service"
	"google.golang.org/grpc"
)

var addr string = "0.0.0.0:50052"

type Server struct {
	infov1.InformationServiceServer
}

type Auth struct {
	Username   string
	PIN        string
	JWTToken  string
	JWTExpiry int64
}

type Vehicle struct {
	RegistrationID string
	VIN string
	Generation string
	Mileage string
}

func main() {
	lis, err := net.Listen("tcp", addr)

	if err != nil {
		log.Fatalf("failed to listen on: %v\n", err)
	}

	log.Printf("vehicle information server listening on %s\n", addr)

	s := grpc.NewServer()
	infov1.RegisterInformationServiceServer(s, &Server{})

	if err = s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v\n", err)
	}
}

func (s *Server) GetVehicleInfo(context context.Context, request *infov1.VehicleInfoRequest) (*infov1.VehicleInfoResponse, error) {
	log.Println("processing vehicle information request")

	info, err := getVehicleInfo(Auth{Username: request.Username, JWTToken: request.JwtToken, PIN: request.Pin, JWTExpiry: request.JwtExpiry})

	if err != nil {
		log.Println("error obtaining vehicle information:", err)
		return &infov1.VehicleInfoResponse{}, err
	}

	log.Println("vehicle information request processed successfully")

	return &infov1.VehicleInfoResponse{
		RegistrationId: info.RegistrationID,
		Vin: info.VIN,
		Generation: info.Generation,
		Mileage: info.Mileage,
	}, nil
}

func getVehicleInfo(auth Auth) (Vehicle, error) {
	// Generate a request to obtain owner information

	req, err := http.NewRequest("POST", "https://owners.hyundaiusa.com/bin/common/MyAccountServlet", nil)

	if err != nil {
		log.Println("error getting owner info req:", err)
		return Vehicle{}, err
	}

	// Set the request headers using a helper method

	setReqHeaders(req, auth)

	// Add query parameters to the request

	q := req.URL.Query()
	q.Add("username", auth.Username)
	q.Add("token", auth.JWTToken)
	q.Add("service", "getOwnerInfoService")
	q.Add("url", "https://owners.hyundaiusa.com/us/en/page/dashboard.html")
	req.URL.RawQuery = q.Encode()

	// Check the response status

	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		log.Println("error obtaining vehicle information:", err)
		return Vehicle{}, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		log.Println("error obtaining vehicle information:", resp.Status)
		return Vehicle{}, err
	}

	// Print the response body as JSON

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		log.Println("error reading owner info response:", err)
		return Vehicle{}, err
	}

	var ownerInfo o.OwnerInfoService

	json.Unmarshal([]byte(body), &ownerInfo)

	vehicles := ownerInfo.ResponseString.OwnersVehiclesInfo

	vehicle := Vehicle{ RegistrationID: vehicles[0].RegistrationID, VIN: vehicles[0].VinNumber, Generation: vehicles[0].IsGen2, Mileage: vehicles[0].Mileage }

	return vehicle, nil
}

func setReqHeaders(req *http.Request, auth Auth) {
	// set request headers
	req.Header.Add("CSRF-Token", "undefined")
	req.Header.Add("accept-language", "en-US,en;q=0.9")
	req.Header.Add("X-Requested-With", "XMLHttpRequest")
	req.Header.Add("Referer", "https://owners.hyundaiusa.com/content/myhyundai/us/en/page/dashboard.html")
	req.Header.Add("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/16.2 Safari/605.1.15")
	req.Header.Add("Host", "owners.hyundaiusa.com")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	req.Header.Add("Origin", "https://owners.hyundaiusa.com")
	req.Header.Add("Cookie", "jwt_token=" + auth.JWTToken + "; s_name=" + auth.Username)
}