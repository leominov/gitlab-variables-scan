package main

import (
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/xanzy/go-gitlab"
)

const (
	secureValueMask = "***"
)

type Scanner struct {
	c      *Config
	git    *Git
	failed bool
}

func NewScanner(c *Config) (*Scanner, error) {
	git, err := NewGit(c.Endpoint, c.Token)
	if err != nil {
		return nil, err
	}
	return &Scanner{
		c:   c,
		git: git,
	}, nil
}

func (s *Scanner) Scan() error {
	groups := []*gitlab.Group{}
	for _, groupID := range s.c.GroupIDs {
		group, err := s.git.GetGroup(groupID)
		if err != nil {
			return err
		}
		groups = append(groups, group)
	}
	subgroups, err := s.fetchGroups()
	if err != nil {
		return err
	}
	groups = append(groups, subgroups...)
	err = s.checkGroupsVariables(groups)
	if err != nil {
		return err
	}

	projects, err := s.fetchProjects()
	if err != nil {
		return err
	}
	err = s.checkProjectsVariables(projects)
	if err != nil {
		return err
	}

	if s.failed {
		return errors.New("Failed. Found sensitive data")
	}
	return nil
}

func (s *Scanner) fetchGroups() ([]*gitlab.Group, error) {
	log.Println("Fetching groups...")
	groupNames := []string{}
	groups := []*gitlab.Group{}
	for _, groupID := range s.c.GroupIDs {
		grs, err := s.git.GetSubGroups(groupID)
		if err != nil {
			return nil, err
		}
		for _, group := range grs {
			groupNames = append(groupNames, group.Name)
		}
		groups = append(groups, grs...)
	}
	if s.c.Debug {
		log.Printf(
			"Found %d group(s): %s",
			len(groups),
			strings.Join(groupNames, ", "),
		)
	}
	return groups, nil
}

func (s *Scanner) fetchProjects() ([]*gitlab.Project, error) {
	log.Println("Fetching projects...")
	projectNames := []string{}
	projects := []*gitlab.Project{}
	for _, groupID := range s.c.GroupIDs {
		prs, err := s.git.GetProjects(groupID)
		if err != nil {
			return nil, err
		}
		for _, project := range prs {
			projectNames = append(projectNames, project.Name)
		}
		projects = append(projects, prs...)
	}
	if s.c.Debug {
		log.Printf(
			"Found %d project(s): %s",
			len(projects),
			strings.Join(projectNames, ", "),
		)
	}
	return projects, nil
}

func (s *Scanner) checkProjectsVariables(projects []*gitlab.Project) error {
	for _, project := range projects {
		log.Printf("Checking %s project...", project.NameWithNamespace)
		vars, err := s.git.GetProjectVariables(project.ID)
		if err != nil {
			return err
		}
		if s.c.Debug {
			log.Printf("Found %d variable(s)", len(vars))
		}
		isContains := s.IsVariablesContainsSensitiveData(vars)
		if isContains {
			s.failed = true
		}
	}
	return nil
}

func (s *Scanner) checkGroupsVariables(groups []*gitlab.Group) error {
	for _, group := range groups {
		log.Printf("Checking %s group...", group.FullName)
		vars, err := s.git.GetGroupVariables(group.ID)
		if err != nil {
			return err
		}
		if s.c.Debug {
			log.Printf("Found %d variable(s)", len(vars))
		}
		isContains := s.IsVariablesContainsSensitiveData(vars)
		if isContains {
			s.failed = true
		}
	}
	return nil
}

func (s *Scanner) IsVariablesContainsSensitiveData(vars []*Variable) bool {
	contains := false
	for _, variable := range vars {
		value := secureValueMask
		if s.c.Insecure {
			value = strings.Replace(variable.Value, "\n", "", -1)
		}
		if re, yes := s.IsVariableContainsSensitiveData(variable); yes {
			log.Printf("  * %s=%s [%s]", variable.Key, value, re)
			contains = true
		}
	}
	return contains
}

func (s *Scanner) IsVariableContainsSensitiveData(variable *Variable) (string, bool) {
	for _, rule := range s.c.Keys {
		if rule.MatchString(variable.Key) {
			return rule.String(), true
		}
	}
	for _, rule := range s.c.Values {
		if rule.MatchString(variable.Value) {
			return rule.String(), true
		}
	}
	for _, rule := range s.c.Pairs {
		pair := fmt.Sprintf("%s=%s", variable.Key, variable.Value)
		if rule.MatchString(pair) {
			return rule.String(), true
		}
	}
	return "", false
}
