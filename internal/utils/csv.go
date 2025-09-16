// Package utils provides CSV parsing utilities for bulk operations.
package utils

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strings"
)

// BulkUser represents a user for bulk creation.
type BulkUser struct {
	Username  string `csv:"username" json:"username"`
	Project   string `csv:"project" json:"project"`
	Blueprint string `csv:"blueprint" json:"blueprint"`
	Bundle    string `csv:"bundle" json:"bundle"`
	Groups    []string `csv:"groups" json:"groups"`
}

// BulkGroup represents a group for bulk creation.
type BulkGroup struct {
	Name        string   `csv:"name" json:"name"`
	Description string   `csv:"description" json:"description"`
	Policies    []string `csv:"policies" json:"policies"`
	Project     string   `csv:"project" json:"project"`
}

// ParseUsersCSV parses a CSV file containing user information.
func ParseUsersCSV(filename string) ([]BulkUser, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open CSV file: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.TrimLeadingSpace = true

	// Read header
	header, err := reader.Read()
	if err != nil {
		return nil, fmt.Errorf("failed to read CSV header: %w", err)
	}

	// Validate required columns
	requiredColumns := []string{"username", "project", "blueprint", "bundle"}
	columnMap := make(map[string]int)
	for i, col := range header {
		columnMap[strings.ToLower(col)] = i
	}

	for _, required := range requiredColumns {
		if _, exists := columnMap[required]; !exists {
			return nil, fmt.Errorf("required column '%s' not found in CSV", required)
		}
	}

	var users []BulkUser
	lineNum := 1

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("error reading CSV line %d: %w", lineNum+1, err)
		}

		lineNum++

		if len(record) != len(header) {
			return nil, fmt.Errorf("line %d: expected %d columns, got %d", lineNum, len(header), len(record))
		}

		user := BulkUser{
			Username:  strings.TrimSpace(record[columnMap["username"]]),
			Project:   strings.TrimSpace(record[columnMap["project"]]),
			Blueprint: strings.TrimSpace(record[columnMap["blueprint"]]),
			Bundle:    strings.TrimSpace(record[columnMap["bundle"]]),
		}

		// Parse groups if column exists
		if groupsIdx, exists := columnMap["groups"]; exists && groupsIdx < len(record) {
			groupsStr := strings.TrimSpace(record[groupsIdx])
			if groupsStr != "" {
				user.Groups = strings.Split(groupsStr, ";")
				for i := range user.Groups {
					user.Groups[i] = strings.TrimSpace(user.Groups[i])
				}
			}
		}

		// Validate required fields
		if user.Username == "" {
			return nil, fmt.Errorf("line %d: username cannot be empty", lineNum)
		}
		if user.Project == "" {
			return nil, fmt.Errorf("line %d: project cannot be empty", lineNum)
		}
		if user.Blueprint == "" {
			return nil, fmt.Errorf("line %d: blueprint cannot be empty", lineNum)
		}
		if user.Bundle == "" {
			return nil, fmt.Errorf("line %d: bundle cannot be empty", lineNum)
		}

		users = append(users, user)
	}

	return users, nil
}

// ParseGroupsCSV parses a CSV file containing group information.
func ParseGroupsCSV(filename string) ([]BulkGroup, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open CSV file: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.TrimLeadingSpace = true

	// Read header
	header, err := reader.Read()
	if err != nil {
		return nil, fmt.Errorf("failed to read CSV header: %w", err)
	}

	// Validate required columns
	requiredColumns := []string{"name", "policies"}
	columnMap := make(map[string]int)
	for i, col := range header {
		columnMap[strings.ToLower(col)] = i
	}

	for _, required := range requiredColumns {
		if _, exists := columnMap[required]; !exists {
			return nil, fmt.Errorf("required column '%s' not found in CSV", required)
		}
	}

	var groups []BulkGroup
	lineNum := 1

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("error reading CSV line %d: %w", lineNum+1, err)
		}

		lineNum++

		if len(record) != len(header) {
			return nil, fmt.Errorf("line %d: expected %d columns, got %d", lineNum, len(header), len(record))
		}

		group := BulkGroup{
			Name: strings.TrimSpace(record[columnMap["name"]]),
		}

		// Parse policies (semicolon separated)
		policiesStr := strings.TrimSpace(record[columnMap["policies"]])
		if policiesStr != "" {
			group.Policies = strings.Split(policiesStr, ";")
			for i := range group.Policies {
				group.Policies[i] = strings.TrimSpace(group.Policies[i])
			}
		}

		// Optional fields
		if descIdx, exists := columnMap["description"]; exists && descIdx < len(record) {
			group.Description = strings.TrimSpace(record[descIdx])
		}
		if projIdx, exists := columnMap["project"]; exists && projIdx < len(record) {
			group.Project = strings.TrimSpace(record[projIdx])
		}

		// Validate required fields
		if group.Name == "" {
			return nil, fmt.Errorf("line %d: name cannot be empty", lineNum)
		}
		if len(group.Policies) == 0 {
			return nil, fmt.Errorf("line %d: policies cannot be empty", lineNum)
		}

		groups = append(groups, group)
	}

	return groups, nil
}

// GenerateUsersCSVTemplate creates a sample CSV template for bulk user creation.
func GenerateUsersCSVTemplate(filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create CSV template: %w", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write header
	header := []string{"username", "project", "blueprint", "bundle", "groups"}
	if err := writer.Write(header); err != nil {
		return fmt.Errorf("failed to write CSV header: %w", err)
	}

	// Write sample data
	samples := [][]string{
		{"alice", "research-team", "ubuntu_22_04", "app_standard_xl_1_0", "researchers;data-analysts"},
		{"bob", "research-team", "ubuntu_22_04", "app_standard_2xl_1_0", "researchers"},
		{"charlie", "ml-team", "ubuntu_22_04", "gpu_nvidia_xl_1_0", "ml-engineers;gpu-users"},
	}

	for _, sample := range samples {
		if err := writer.Write(sample); err != nil {
			return fmt.Errorf("failed to write CSV sample: %w", err)
		}
	}

	return nil
}

// GenerateGroupsCSVTemplate creates a sample CSV template for bulk group creation.
func GenerateGroupsCSVTemplate(filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create CSV template: %w", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write header
	header := []string{"name", "description", "policies", "project"}
	if err := writer.Write(header); err != nil {
		return fmt.Errorf("failed to write CSV header: %w", err)
	}

	// Write sample data
	samples := [][]string{
		{"researchers", "Research team members", "arn:aws:iam::aws:policy/ReadOnlyAccess", "research-team"},
		{"data-analysts", "Data analysis team", "arn:aws:iam::aws:policy/ReadOnlyAccess;arn:aws:iam::aws:policy/IAMUserChangePassword", "research-team"},
		{"ml-engineers", "Machine learning engineers", "arn:aws:iam::aws:policy/PowerUserAccess", "ml-team"},
	}

	for _, sample := range samples {
		if err := writer.Write(sample); err != nil {
			return fmt.Errorf("failed to write CSV sample: %w", err)
		}
	}

	return nil
}