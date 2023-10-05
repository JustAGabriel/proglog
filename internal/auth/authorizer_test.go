package auth

import (
	"testing"

	"github.com/justagabriel/proglog/internal/config"
	"github.com/stretchr/testify/require"
)

const (
	validSubject string = "root"
)

func TestACLAuthorization(t *testing.T) {
	authorizer, err := New(config.ACLModelFile, config.ACLPolicyFile)
	require.NoError(t, err)

	scenarios := map[string]func(*testing.T, *Authorizer){
		"valid 'get' credentials returns 'true'":    testValidGetCreds,
		"valid 'create' credentials returns 'true'": testValidCreateCreds,
		"invalid subject returns 'false'":           testInvalidSubject,
		"invalid action returns 'false'":            testInvalidAction,
	}

	for scenario, test := range scenarios {
		t.Run(scenario, func(t *testing.T) {
			test(t, authorizer)
		})
	}
}

func testValidGetCreds(t *testing.T, authorizer *Authorizer) {
	// arrange
	getAction := "get"

	// act
	authError := authorizer.Authorize(validSubject, getAction)

	// assert
	require.NoError(t, authError, "credentials are valid - should work")
}

func testValidCreateCreds(t *testing.T, authorizer *Authorizer) {
	// arrange
	getAction := "create"

	// act
	authError := authorizer.Authorize(validSubject, getAction)

	// assert
	require.NoError(t, authError, "credentials are valid - should work")
}

func testInvalidSubject(t *testing.T, authorizer *Authorizer) {
	// arrange
	invalidSubject := "randomUser"
	validAction := "create"

	// act
	authError := authorizer.Authorize(invalidSubject, validAction)

	// assert
	require.Error(t, authError, "subject is not defined - should fail")
}

func testInvalidAction(t *testing.T, authorizer *Authorizer) {
	// arrange
	invalidAction := "destroy"

	// act
	authError := authorizer.Authorize(validSubject, invalidAction)

	// assert
	require.Error(t, authError, "action is not defined - should fail")
}
