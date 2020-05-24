package services_test

import (
	"testing"

	"github.com/kaancfidan/bouncer/models"
	"github.com/kaancfidan/bouncer/services"
)

func Test_AuthorizerImpl_Authorize(t *testing.T) {
	type args struct {
		policyNames []string
		claims      map[string]interface{}
	}

	tests := []struct {
		name            string
		claimPolicies   map[string][]models.ClaimRequirement
		args            args
		wantFailedClaim string
		wantErr         bool
	}{
		{
			name:          "zero config - no claims",
			claimPolicies: map[string][]models.ClaimRequirement{},
			args: args{
				policyNames: make([]string, 0),
				claims:      map[string]interface{}{},
			},
			wantFailedClaim: "",
			wantErr:         false,
		},
		{
			name:          "zero config - irrelevant claims",
			claimPolicies: map[string][]models.ClaimRequirement{},
			args: args{
				policyNames: make([]string, 0),
				claims: map[string]interface{}{
					"claim": "value",
				},
			},
			wantFailedClaim: "",
			wantErr:         false,
		},
		{
			name:          "non-existing policy",
			claimPolicies: map[string][]models.ClaimRequirement{},
			args: args{
				policyNames: []string{"NonExistingPolicyName"},
				claims:      map[string]interface{}{},
			},
			wantFailedClaim: "",
			wantErr:         true,
		},
		{
			name: "claim exists",
			claimPolicies: map[string][]models.ClaimRequirement{
				"HasName": {
					models.ClaimRequirement{
						Claim: "name",
					},
				},
			},
			args: args{
				policyNames: []string{"HasName"},
				claims: map[string]interface{}{
					"name": "John",
				},
			},
			wantFailedClaim: "",
			wantErr:         false,
		},
		{
			name: "claim does not exist",
			claimPolicies: map[string][]models.ClaimRequirement{
				"HasName": {
					models.ClaimRequirement{
						Claim: "name",
					},
				},
			},
			args: args{
				policyNames: []string{"HasName"},
				claims: map[string]interface{}{
					"last_name": "Doe",
				},
			},
			wantFailedClaim: "name",
			wantErr:         false,
		},
		{
			name: "claim value matches",
			claimPolicies: map[string][]models.ClaimRequirement{
				"NamedJohn": {
					models.ClaimRequirement{
						Claim:  "name",
						Values: []string{"John"},
					},
				},
			},
			args: args{
				policyNames: []string{"NamedJohn"},
				claims: map[string]interface{}{
					"name": "John",
				},
			},
			wantFailedClaim: "",
			wantErr:         false,
		},
		{
			name: "claim value does not match",
			claimPolicies: map[string][]models.ClaimRequirement{
				"NamedJohn": {
					models.ClaimRequirement{
						Claim:  "name",
						Values: []string{"John"},
					},
				},
			},
			args: args{
				policyNames: []string{"NamedJohn"},
				claims: map[string]interface{}{
					"name": "Jane",
				},
			},
			wantFailedClaim: "name",
			wantErr:         false,
		},
		{
			name: "multiple claim values match",
			claimPolicies: map[string][]models.ClaimRequirement{
				"SpecificJohn": {
					models.ClaimRequirement{
						Claim:  "name",
						Values: []string{"John"},
					},
					models.ClaimRequirement{
						Claim:  "last_name",
						Values: []string{"Doe"},
					},
				},
			},
			args: args{
				policyNames: []string{"SpecificJohn"},
				claims: map[string]interface{}{
					"name":      "John",
					"last_name": "Doe",
				},
			},
			wantFailedClaim: "",
			wantErr:         false,
		},
		{
			name: "multiple claim values, one does not match",
			claimPolicies: map[string][]models.ClaimRequirement{
				"SpecificJohn": {
					models.ClaimRequirement{
						Claim:  "name",
						Values: []string{"John"},
					},
					models.ClaimRequirement{
						Claim:  "last_name",
						Values: []string{"Doe"},
					},
				},
			},
			args: args{
				policyNames: []string{"SpecificJohn"},
				claims: map[string]interface{}{
					"name":      "Jane",
					"last_name": "Doe",
				},
			},
			wantFailedClaim: "name",
			wantErr:         false,
		},
		{
			name: "array claim value matches",
			claimPolicies: map[string][]models.ClaimRequirement{
				"CanTest": {
					models.ClaimRequirement{
						Claim:  "permission",
						Values: []string{"Test"},
					},
				},
			},
			args: args{
				policyNames: []string{"CanTest"},
				claims: map[string]interface{}{
					"permission": []interface{}{"Test"},
				},
			},
			wantFailedClaim: "",
			wantErr:         false,
		},
		{
			name: "array claim does not match",
			claimPolicies: map[string][]models.ClaimRequirement{
				"CanDelete": {
					models.ClaimRequirement{
						Claim:  "permission",
						Values: []string{"Delete"},
					},
				},
			},
			args: args{
				policyNames: []string{"CanDelete"},
				claims: map[string]interface{}{
					"permission": []interface{}{"Test", "Add"},
				},
			},
			wantFailedClaim: "permission",
			wantErr:         false,
		},
		{
			name: "value claim match to array",
			claimPolicies: map[string][]models.ClaimRequirement{
				"CanAddOrDelete": {
					models.ClaimRequirement{
						Claim:  "permission",
						Values: []string{"Add", "Delete"},
					},
				},
			},
			args: args{
				policyNames: []string{"CanAddOrDelete"},
				claims: map[string]interface{}{
					"permission": "Add",
				},
			},
			wantFailedClaim: "",
			wantErr:         false,
		},
		{
			name: "array claim match to array",
			claimPolicies: map[string][]models.ClaimRequirement{
				"CanAddOrDelete": {
					models.ClaimRequirement{
						Claim:  "permission",
						Values: []string{"Add", "Delete"},
					},
				},
			},
			args: args{
				policyNames: []string{"CanAddOrDelete"},
				claims: map[string]interface{}{
					"permission": []interface{}{"Test", "Add"},
				},
			},
			wantFailedClaim: "",
			wantErr:         false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := services.NewAuthorizer(tt.claimPolicies)

			gotFailedPolicy, err := a.Authorize(tt.args.policyNames, tt.args.claims)

			if (err != nil) != tt.wantErr {
				t.Errorf("Authenticate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if gotFailedPolicy != tt.wantFailedClaim {
				t.Errorf("Authorize() = %v, want %v", gotFailedPolicy, tt.wantFailedClaim)
			}
		})
	}
}

func Test_AuthorizerImpl_IsAnonymousAllowed(t *testing.T) {
	tests := []struct {
		name            string
		matchedPolicies []models.RoutePolicy
		want            bool
	}{
		{
			name:            "empty config",
			matchedPolicies: []models.RoutePolicy{},
			want:            false,
		},
		{
			name:            "single allow",
			matchedPolicies: []models.RoutePolicy{{AllowAnonymous: true}},
			want:            true,
		},
		{
			name:            "single disallow",
			matchedPolicies: []models.RoutePolicy{{AllowAnonymous: false}},
			want:            false,
		},
		{
			name: "one allow one disallow",
			matchedPolicies: []models.RoutePolicy{
				{AllowAnonymous: true},
				{AllowAnonymous: false},
			},
			want: false,
		},
		{
			name: "both allow",
			matchedPolicies: []models.RoutePolicy{
				{AllowAnonymous: true},
				{AllowAnonymous: true},
			},
			want: true,
		},
		{
			name: "both disallow",
			matchedPolicies: []models.RoutePolicy{
				{AllowAnonymous: false},
				{AllowAnonymous: false},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := services.NewAuthorizer(map[string][]models.ClaimRequirement{})

			if got := a.IsAnonymousAllowed(tt.matchedPolicies); got != tt.want {
				t.Errorf("IsAnonymousAllowed() = %v, want %v", got, tt.want)
			}
		})
	}
}