package services

import (
	"fmt"
	"io"

	"github.com/go-yaml/yaml"

	"github.com/kaancfidan/bouncer/models"
)

// ConfigParser is the config parsing interface
type ConfigParser interface {
	ParseConfig(reader io.Reader) (*models.Config, error)
}

// YamlConfigParser is the YAML deserialization implementation of ConfigParser
type YamlConfigParser struct{}

// ParseConfig implements config parsing from YAML files
func (_ YamlConfigParser) ParseConfig(reader io.Reader) (*models.Config, error) {
	decoder := yaml.NewDecoder(reader)

	cfg := models.Config{}
	err := decoder.Decode(&cfg)
	if err != nil {
		return nil, fmt.Errorf("could not parse config yaml: %v", err)
	}

	return &cfg, nil
}

// ValidateConfig validates a parsed Config struct against following constraints:
// - Both claim policies and route policies must not be nil. Empty map/slices are allowed.
// - All ClaimRequirement instances must have a claim named.
// - All RoutePolicy instances must have a path configured.
// - If a RoutePolicy is flagged with AllowAnonymous, it cannot name any claim policies
// - If a RoutePolicy has a claim policy named, that claim policy should be defined in the ClaimPolicies section.
func ValidateConfig(config *models.Config) error {
	if config.ClaimPolicies == nil {
		return fmt.Errorf("claim policies nil")
	}
	if config.RoutePolicies == nil {
		return fmt.Errorf("route policies nil")
	}

	// check claim policies
	for policyName, policy := range config.ClaimPolicies {
		for _, requirement := range policy {
			// Claim field is mandatory
			if requirement.Claim == "" {
				return fmt.Errorf("found claim policy (%s) with unnamed claim requirement: %v",
					policyName, policy)
			}
		}
	}

	// existing claim policy names
	existingPolicies := make(map[string]bool)
	for k := range config.ClaimPolicies {
		existingPolicies[k] = true
	}

	// check route policies
	for _, p := range config.RoutePolicies {
		if p.Path == "" {
			return fmt.Errorf("found route policy without a path denition: %v", p)
		}

		// anonymous routes cannot name claim policies
		if p.AllowAnonymous && (p.PolicyName != "") {
			return fmt.Errorf("found route policy with ambiguous claim policy config: %v", p)
		}

		// non-existing policy check (~foreign key constraint)
		if p.PolicyName != "" && !existingPolicies[p.PolicyName] {
			return fmt.Errorf("non-existing policy name found in route policy: %v", p)
		}
	}

	return nil
}