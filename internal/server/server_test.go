/*
	(c) Copyright 2015 Vlad Didenko

	Licensed under the Apache License, Version 2.0 (the "License");
	you may not use this file except in compliance with the License.
	You may obtain a copy of the License at

	    http://www.apache.org/licenses/LICENSE-2.0

	Unless required by applicable law or agreed to in writing, software
	distributed under the License is distributed on an "AS IS" BASIS,
	WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
	See the License for the specific language governing permissions and
	limitations under the License.
*/

package server

import (
	"bufio"
	"net/http"
	"strconv"
	"strings"
	"testing"
)

const (
	testPort uint16 = 8001
)

var (
	testUrl string = "http://localhost:" + strconv.Itoa(int(testPort))
)

func TestPaldHttp(t *testing.T) {

	testCases := []struct {
		request  string
		httpCode int
		respFore string
	}{
		{request: "/set?service=a0", httpCode: http.StatusOK, respFore: "49200"},
		{request: "/set?service=a1", httpCode: http.StatusOK, respFore: "49201"},
		{request: "/get?service=a2", httpCode: http.StatusNotFound, respFore: "Name \"a2\" not found in the port registry"},
		{request: "/set?service=a2", httpCode: http.StatusOK, respFore: "49202"},
		{request: "/get?service=a1", httpCode: http.StatusOK, respFore: "49201"},
		{request: "/set?service=a2", httpCode: http.StatusPreconditionFailed, respFore: "Name \"a2\" is already taken"},
		{request: "/set?service=a2&port=49000", httpCode: http.StatusPreconditionFailed, respFore: "Name \"a2\" is already taken"},
		{request: "/set?service=f0&port=49000", httpCode: http.StatusOK, respFore: "OK"},
		{request: "/get?service=f0", httpCode: http.StatusOK, respFore: "49000"},
		{request: "/del?port=49201", httpCode: http.StatusOK, respFore: "OK"},
		{request: "/set?service=a3", httpCode: http.StatusOK, respFore: "49201"},
		{request: "/set?service=a4", httpCode: http.StatusPreconditionFailed, respFore: "No ports available"},
		{request: "/set?svc=er", httpCode: http.StatusBadRequest, respFore: "Service name is missing"},
		{request: "/get?svc=er", httpCode: http.StatusBadRequest, respFore: "Service name is missing"},
		{request: "/set?service=f1&port=49001", httpCode: http.StatusOK, respFore: "OK"},
		{request: "/del?port=49O01", httpCode: http.StatusBadRequest, respFore: "strconv.ParseUint:"},
		{request: "/set?sevice=er", httpCode: http.StatusBadRequest, respFore: "Service name is missing"},
		{request: "/del?svc=er", httpCode: http.StatusBadRequest, respFore: "Port number is missing"},
	}

	go Run(testPort, 49200, 49202)

	for _, tc := range testCases {

		resp, err := http.Get(testUrl + tc.request)
		defer resp.Body.Close()
		if err != nil {
			t.Error(err)
		}

		if resp.StatusCode != tc.httpCode {
			t.Errorf("Received code %d instead of %d from %q request", resp.StatusCode, tc.httpCode, tc.request)
		}

		line, err := bufio.NewReader(resp.Body).ReadString('\n')
		if !strings.HasPrefix(line, tc.respFore) {
			t.Errorf("Wrong response body. Expected to start with %q, but it is %q", tc.respFore, line)
		}
	}
}
