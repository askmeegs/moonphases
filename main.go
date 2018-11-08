// Copyright 2018 Google Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"os"
	"time"

	_ "github.com/jnewmano/grpc-json-proxy/codec"
	pb "github.com/m-okeefe/moonphases/proto"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
)

var (
	addr     = flag.String("addr", ":8001", "[host]:port to listen")
	logLevel = flag.String("log-level", "info", "info, debug, warn, error")
	city     = flag.String("city", "San Francisco, CA", "<City>, <State-Abbrev>")
	log      *logrus.Entry
)

type MoonPhasesServer struct {
	City string
}

// USNO API Struct
type RawPhases struct {
	Error        bool        `json:"error"`
	ApiVersion   string      `json:"apiversion"`
	Year         int         `json:"year"`
	Month        int         `json:"month"`
	Day          int         `json:"day"`
	DayOfWeek    string      `json:"dayofweek"`
	DateChanged  bool        `json:"datechanged"`
	IsDst        string      `json:"isdst"`
	County       string      `json:"county"`
	Tz           int         `json:"tz"`
	State        string      `json:"state"`
	City         string      `json:"city"`
	Lon          float32     `json:"lon"`
	Lat          float32     `json:"lat"`
	MoonData     []Phenomena `json:"moondata"`
	SunData      []Phenomena `json:"sundata"`
	ClosestPhase Phase       `json:"closestphase"`
}
type Phenomena struct {
	Phen string `json:"phen"`
	Time string `json:"time"`
}

type Phase struct {
	Phase string `json:"phase"`
	Date  string `json:"date"`
	Time  string `json:"time"`
}

func init() {
	flag.Parse()
	host, err := os.Hostname()
	if err != nil {
		log.Fatal(errors.Wrap(err, "cannot get hostname"))
	}
	switch *logLevel {
	case "error":
		logrus.SetLevel(logrus.ErrorLevel)
	case "warn":
		logrus.SetLevel(logrus.WarnLevel)
	case "debug":
		logrus.SetLevel(logrus.DebugLevel)
	default:
		logrus.SetLevel(logrus.InfoLevel)
	}

	logrus.SetFormatter(&logrus.JSONFormatter{FieldMap: logrus.FieldMap{logrus.FieldKeyLevel: "severity"}})
	log = logrus.WithFields(logrus.Fields{
		"service": "moonphases",
		"host":    host,
	})
	grpclog.SetLogger(log.WithField("facility", "grpc"))
}

func main() {
	lis, err := net.Listen("tcp", *addr)
	if err != nil {
		log.Fatal(err)
	}
	grpcServer := grpc.NewServer()

	// Initialize new phases server
	s := &MoonPhasesServer{*city}
	pb.RegisterMoonPhasesServer(grpcServer, s)
	log.Info(" 🌒 moon phases grpc server 🌔   ")
	log.WithField("addr", *addr).Info("starting to listen on grpc")
	log.Fatal(grpcServer.Serve(lis))
}

// GetPhases is the sole endpoint for the PhasesServer
// Calls the Naval Observatory API to get moon phase info for a <City> on Today's <Date>
func (s *MoonPhasesServer) GetPhases(ctx context.Context, req *pb.GetPhasesRequest) (*pb.GetPhasesResponse, error) {

	log.Info("GET")

	// Construct query
	currentTime := time.Now().Local()
	f := currentTime.Format("01/7/2006")
	Url, _ := url.Parse("http://api.usno.navy.mil")
	Url.Path += "/rstt/oneday"
	parameters := url.Values{}
	parameters.Add("date", f)
	parameters.Add("loc", s.City)
	Url.RawQuery = parameters.Encode()
	fmt.Printf("Encoded URL is %q\n", Url.String())

	// Call USNO API
	resp, err := http.Get(Url.String())
	if err != nil {
		log.Error(err)
		return &pb.GetPhasesResponse{}, err
	}
	defer resp.Body.Close()

	// Unmarshal raw json into RawPhases struct
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error(err)
		return &pb.GetPhasesResponse{}, err
	}
	var r RawPhases
	err = json.Unmarshal(body, &r)
	if err != nil {
		log.Error(err)
		return &pb.GetPhasesResponse{}, err
	}

	// Ensure valid data
	m := r.MoonData
	if len(m) < 2 {
		log.Error("could not get adequate moondata")
		return &pb.GetPhasesResponse{}, fmt.Errorf("missing moondata: %#v", m)
	}
	rise, upperTransit, set := m[0], m[1], m[2]
	closest := r.ClosestPhase

	// Parse RawPhases into PhaseInfo proto format
	p := &pb.PhaseInfo{
		City:         fmt.Sprintf("%s, %s", r.City, r.State),
		Lat:          fmt.Sprintf("%f", r.Lat),
		Lon:          fmt.Sprintf("%f", r.Lon),
		ClosestPhase: fmt.Sprintf("%s: %s %s", closest.Phase, closest.Date, closest.Time),
		Rise:         fmt.Sprintf("%s - %s", rise.Phen, rise.Time),
		UpperTransit: fmt.Sprintf("%s - %s", upperTransit.Phen, upperTransit.Time),
		Set:          fmt.Sprintf("%s - %s", set.Phen, set.Time),
	}
	return &pb.GetPhasesResponse{PhaseInfo: p}, nil
}
