package httpx

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"reflect"
	"strings"
)

type RequestBody interface {
	string | url.Values | map[string]string
}

var (
	ErrorUnSupportMethodType = errors.New("not support method type")
	ErrorUnSupportBodyType   = errors.New("not support body type")
	ErrorMarshalBody         = errors.New("marshal body failed")
)

func HttpContextGet[T RequestBody](ctx context.Context, httpUrl string, body T) ([]byte, error) {
	request, err := createRequest(ctx, "GET", httpUrl, body)
	if err != nil {
		return nil, err
	}
	return doAndRsp(request)
}
func HttpContextPost[T RequestBody](ctx context.Context, url string, body T) ([]byte, error) {
	request, err := createRequest(ctx, "POST", url, body)
	if err != nil {
		return nil, err
	}
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	return doAndRsp(request)
}
func HttpContextJson(ctx context.Context, url string, body interface{}) ([]byte, error) {
	request, err := createRequest(ctx, "JSON", url, body)
	if err != nil {
		return nil, err
	}
	request.Header.Set("Content-Type", "application/json")
	return doAndRsp(request)
}

func createRequest(ctx context.Context, method, httpUrl string, body interface{}) (*http.Request, error) {
	var bodyReader io.Reader
	if strings.ToUpper(method) == "JSON" {
		rValue := reflect.ValueOf(body)
		if rValue.Kind() == reflect.Pointer {
			rValue = rValue.Elem()
		}
		switch rValue.Kind() {
		case reflect.Map, reflect.Struct, reflect.Slice:
			if bodyBytes, ok := body.([]byte); ok {
				bodyReader = bytes.NewReader(bodyBytes)
			} else {
				bodyBytes, err := json.Marshal(body)
				if err != nil {
					return nil, ErrorMarshalBody
				}
				bodyReader = bytes.NewReader(bodyBytes)
			}
		case reflect.String:
			bodyReader = strings.NewReader(body.(string))
		default:
			return nil, ErrorUnSupportBodyType
		}

	} else {
		switch body.(type) {
		case map[string]string:
			urlValues := make(url.Values)

			for k, v := range body.(map[string]string) {
				urlValues[k] = []string{v}
			}
			bodyReader = strings.NewReader(urlValues.Encode())
		case url.Values:
			bodyReader = strings.NewReader(body.(url.Values).Encode())
		case string:
			bodyReader = strings.NewReader(body.(string))
		case []byte:
			bodyReader = bytes.NewReader(body.([]byte))
		default:
			return nil, ErrorUnSupportBodyType
		}
	}
	switch strings.ToUpper(method) {
	case "GET":
		bodyBytes, _ := io.ReadAll(bodyReader)
		if len(bodyBytes) > 0 {
			httpUrl += "?" + string(bodyBytes)
		}
		return http.NewRequestWithContext(ctx, "GET", httpUrl, nil)
	case "POST", "JSON":
		return http.NewRequestWithContext(ctx, "POST", httpUrl, bodyReader)
	default:
		return nil, ErrorUnSupportMethodType
	}
}

func doAndRsp(request *http.Request) ([]byte, error) {
	rsp, err := http.DefaultClient.Do(request)
	if err != nil {
		return nil, err
	}
	if rsp.Body != nil {
		defer rsp.Body.Close()
	}
	rspByte, err := io.ReadAll(rsp.Body)
	if err != nil {
		return nil, err
	}
	return rspByte, nil
}
