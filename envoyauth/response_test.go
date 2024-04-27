package envoyauth

import (
	"encoding/json"
	"reflect"
	"testing"

	ext_core_v3 "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	_structpb "github.com/golang/protobuf/ptypes/struct"
	"google.golang.org/protobuf/proto"
)

func TestIsAllowed(t *testing.T) {

	input := make(map[string]interface{})
	er := EvalResult{
		Decision: input,
	}
	var err error

	_, err = er.IsAllowed()
	if err == nil {
		t.Fatal("Expected error but got nil")
	}

	input["allowed"] = 1
	_, err = er.IsAllowed()
	if err == nil {
		t.Fatal("Expected error but got nil")
	}

	input["allowed"] = true
	var result bool
	result, err = er.IsAllowed()

	if err != nil {
		t.Fatalf("Expected no error but got %v", err)
	}

	if result != true {
		t.Fatalf("Expected value for IsAllowed %v but got %v", true, result)
	}
}

func TestGetRequestHTTPHeadersToRemove(t *testing.T) {
	tests := map[string]struct {
		decision interface{}
		exp      []string
		wantErr  bool
	}{
		"bool_eval_result": {
			true,
			[]string{},
			false,
		},
		"invalid_eval_result": {
			"hello",
			[]string{},
			true,
		},
		"empty_map_result": {
			map[string]interface{}{},
			[]string{},
			false,
		},
		"bad_header_value": {
			map[string]interface{}{"request_headers_to_remove": "test"},
			[]string{},
			true,
		},
		"string_array_header_value": {
			map[string]interface{}{"request_headers_to_remove": []string{"foo", "bar"}},
			[]string{"foo", "bar"},
			false,
		},
		"interface_array_header_value": {
			map[string]interface{}{"request_headers_to_remove": []interface{}{"foo", "bar", "fuz"}},
			[]string{"foo", "bar", "fuz"},
			false,
		},
		"interface_array_bad_header_value": {
			map[string]interface{}{"request_headers_to_remove": []interface{}{1}},
			[]string{},
			true,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			er := EvalResult{
				Decision: tc.decision,
			}

			result, err := er.GetRequestHTTPHeadersToRemove()

			if tc.wantErr {
				if err == nil {
					t.Fatal("Expected error but got nil")
				}
			} else {
				if err != nil {
					t.Fatalf("Unexpected error %v", err)
				}

				if !reflect.DeepEqual(tc.exp, result) {
					t.Fatalf("Expected result %v but got %v", tc.exp, result)
				}
			}
		})
	}
}

func TestGetResponseHTTPHeadersToAdd(t *testing.T) {
	input := make(map[string]interface{})
	er := EvalResult{
		Decision: input,
	}

	result, err := er.GetResponseHTTPHeadersToAdd()
	if err != nil {
		t.Fatalf("Expected no error but got %v", err)
	}

	if len(result) != 0 {
		t.Fatal("Expected no headers")
	}

	badHeader := "test"
	input["response_headers_to_add"] = badHeader

	_, err = er.GetResponseHTTPHeadersToAdd()
	if err == nil {
		t.Fatal("Expected error but got nil")
	}

	testHeaders := make(map[string]interface{})
	testHeaders["foo"] = "bar"
	input["response_headers_to_add"] = testHeaders

	result, err = er.GetResponseHTTPHeadersToAdd()
	if err != nil {
		t.Fatalf("Expected no error but got %v", err)
	}

	if len(result) != 1 {
		t.Fatalf("Expected one header but got %v", len(result))
	}

	testHeaders["baz"] = 1

	_, err = er.GetResponseHTTPHeadersToAdd()
	if err == nil {
		t.Fatal("Expected error but got nil")
	}

	input["response_headers_to_add"] = []interface{}{
		map[string]interface{}{
			"foo": "bar",
		},
		map[string]interface{}{
			"foo": "baz",
		},
	}

	result, err = er.GetResponseHTTPHeadersToAdd()
	if err != nil {
		t.Fatalf("Expected no error but got %v", err)
	}

	if len(result) != 2 {
		t.Fatalf("Expected two headers but got %v", len(result))
	}

	testAddHeaders := make(map[string]interface{})
	testAddHeaders["foo"] = []string{"bar", "baz"}
	input["response_headers_to_add"] = testAddHeaders

	result, err = er.GetResponseHTTPHeadersToAdd()

	if err != nil {
		t.Fatalf("Expected no error but got %v", err)
	}

	if len(result) != 2 {
		t.Fatalf("Expected two header but got %v", len(result))
	}
}

func TestGetResponseHeaders(t *testing.T) {
	input := make(map[string]interface{})
	er := EvalResult{
		Decision: input,
	}

	result, err := er.GetResponseEnvoyHeaderValueOptions()
	if err != nil {
		t.Fatalf("Expected no error but got %v", err)
	}

	if len(result) != 0 {
		t.Fatal("Expected no headers")
	}

	badHeader := "test"
	input["headers"] = badHeader

	_, err = er.GetResponseEnvoyHeaderValueOptions()
	if err == nil {
		t.Fatal("Expected error but got nil")
	}

	testHeaders := make(map[string]interface{})
	testHeaders["foo"] = "bar"
	input["headers"] = testHeaders

	result, err = er.GetResponseEnvoyHeaderValueOptions()
	if err != nil {
		t.Fatalf("Expected no error but got %v", err)
	}

	if len(result) != 1 {
		t.Fatalf("Expected one header but got %v", len(result))
	}

	testHeaders["baz"] = 1

	_, err = er.GetResponseEnvoyHeaderValueOptions()
	if err == nil {
		t.Fatal("Expected error but got nil")
	}

	input["headers"] = []interface{}{
		map[string]interface{}{
			"foo": "bar",
		},
		map[string]interface{}{
			"foo": "baz",
		},
	}

	result, err = er.GetResponseEnvoyHeaderValueOptions()
	if err != nil {
		t.Fatalf("Expected no error but got %v", err)
	}

	if len(result) != 2 {
		t.Fatalf("Expected two header but got %v", len(result))
	}

	seen := map[string]int{}
	for _, hdr := range result {
		seen[hdr.Header.Value]++
	}
	if seen["bar"] != 1 || seen["baz"] != 1 {
		t.Errorf("expected 'bar' and 'baz', got %v", seen)
	}

	testAddHeaders := make(map[string]interface{})
	testAddHeaders["foo"] = []string{"bar", "baz"}
	input["headers"] = testAddHeaders

	result, err = er.GetResponseEnvoyHeaderValueOptions()

	if err != nil {
		t.Fatalf("Expected no error but got %v", err)
	}

	if len(result) != 2 {
		t.Fatalf("Expected two header but got %v", len(result))
	}
}

func TestGetResponseBody(t *testing.T) {
	input := make(map[string]interface{})
	er := EvalResult{
		Decision: input,
	}

	result, err := er.GetResponseBody()
	if err != nil {
		t.Fatalf("Expected no error but got %v", err)
	}

	if result != "" {
		t.Fatalf("Expected empty body but got %v", result)
	}

	input["body"] = "hello"
	result, err = er.GetResponseBody()
	if err != nil {
		t.Fatalf("Expected no error but got %v", err)
	}

	if result != "hello" {
		t.Fatalf("Expected result \"hello\" but got %v", result)
	}

	input["body"] = 123
	result, err = er.GetResponseBody()
	if err == nil {
		t.Fatal("Expected error but got nil", err)
	}
}

func TestGetResponseHttpStatus(t *testing.T) {
	input := make(map[string]interface{})
	er := EvalResult{
		Decision: input,
	}

	result, err := er.GetResponseEnvoyHTTPStatus()
	if err != nil {
		t.Fatalf("Expected no error but got %v", err)
	}

	if result.GetCode().String() != "Forbidden" {
		t.Fatalf("Expected http status code \"Forbidden\" but got %v", result.GetCode().String())
	}

	input["http_status"] = true
	result, err = er.GetResponseEnvoyHTTPStatus()
	if err == nil {
		t.Fatal("Expected error but got nil")
	}

	input["http_status"] = json.Number("1")
	result, err = er.GetResponseEnvoyHTTPStatus()
	if err == nil {
		t.Fatal("Expected error but got nil")
	}

	input["http_status"] = json.Number("9999")
	result, err = er.GetResponseEnvoyHTTPStatus()
	if err == nil {
		t.Fatal("Expected error but got nil")
	}

	input["http_status"] = json.Number("400")
	result, err = er.GetResponseEnvoyHTTPStatus()
	if err != nil {
		t.Fatalf("Expected no error but got %v", err)
	}

	if result.GetCode().String() != "BadRequest" {
		t.Fatalf("Expected http status code \"BadRequest\" but got %v", result.GetCode().String())
	}
}

func TestGetDynamicMetadata(t *testing.T) {
	input := make(map[string]interface{})
	er := EvalResult{
		Decision: input,
	}

	result, err := er.GetDynamicMetadata()
	if err != nil {
		t.Fatalf("Expected no error but got %v", err)
	}

	if result != nil {
		t.Fatalf("Expected no dynamic metadata but got %v", result)
	}

	input["dynamic_metadata"] = map[string]interface{}{
		"foo": "bar",
	}
	result, err = er.GetDynamicMetadata()
	if err != nil {
		t.Fatalf("Expected no error but got %v", err)
	}

	expectedDynamicMetadata := &_structpb.Struct{
		Fields: map[string]*_structpb.Value{
			"foo": {
				Kind: &_structpb.Value_StringValue{
					StringValue: "bar",
				},
			},
		},
	}
	if !proto.Equal(result, expectedDynamicMetadata) {
		t.Fatalf("Expected result %v but got %v", expectedDynamicMetadata, result)
	}
}

func TestGetDynamicMetadataWithBooleanDecision(t *testing.T) {
	er := EvalResult{
		Decision: true,
	}

	result, err := er.GetDynamicMetadata()
	if err == nil {
		t.Fatal("Expected error error but got none")
	}

	if result != nil {
		t.Fatalf("Expected no result but got %v", result)
	}
}

func TestNewEvalResultWithDecisionID(t *testing.T) {
	type Opt func(*EvalResult)

	withDecisionID := func(decisionID string) Opt {
		return func(result *EvalResult) {
			result.DecisionID = decisionID
		}
	}

	expectedDecisionID := "some-decision-id"

	er, _, err := NewEvalResult(withDecisionID(expectedDecisionID))
	if err != nil {
		t.Fatalf("NewEvalResult() error = %v, wantErr %v", err, false)
	}
	if er.DecisionID != expectedDecisionID {
		t.Errorf("Expected DecisionID to be '%v', got '%v'", expectedDecisionID, er.DecisionID)
	}
}

func TestGetQueryParametersToSet(t *testing.T) {
	tests := map[string]struct {
		decision interface{}
		exp      []*ext_core_v3.QueryParameter
		wantErr  bool
	}{
		"bool_eval_result": {
			true,
			[]*ext_core_v3.QueryParameter{},
			false,
		},
		"invalid_eval_result": {
			"hello",
			[]*ext_core_v3.QueryParameter{},
			true,
		},
		"empty_map_result": {
			map[string]interface{}{},
			[]*ext_core_v3.QueryParameter{},
			false,
		},
		"bad_header_value": {
			map[string]interface{}{"query_parameters_to_set": "test"},
			[]*ext_core_v3.QueryParameter{},
			true,
		},
		"bad_string_array_header_value": {
			map[string]interface{}{"query_parameters_to_set": []string{"foo", "bar"}},
			[]*ext_core_v3.QueryParameter{},
			true,
		},
		"ok1": {
			map[string]interface{}{"query_parameters_to_set": []map[string]interface{}{
				{"key": "abc", "value": "123"},
			}},
			[]*ext_core_v3.QueryParameter{
				{
					Key:   "abc",
					Value: "123",
				},
			},
			false,
		},
		"ok2": {
			map[string]interface{}{"query_parameters_to_set": []map[string]interface{}{
				{"key": "abc", "value": "123"},
				{"key": "xyz", "value": "987"},
			}},
			[]*ext_core_v3.QueryParameter{
				{
					Key:   "abc",
					Value: "123",
				},
				{
					Key:   "xyz",
					Value: "987",
				},
			},
			false,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			er := EvalResult{
				Decision: tc.decision,
			}

			result, err := er.GetQueryParametersToSet()

			if tc.wantErr {
				if err == nil {
					t.Fatal("Expected error but got nil")
				}
			} else {
				if err != nil {
					t.Fatalf("Unexpected error %v", err)
				}

				if !reflect.DeepEqual(tc.exp, result) {
					t.Fatalf("Expected result %v but got %v", tc.exp, result)
				}
			}
		})
	}
}
