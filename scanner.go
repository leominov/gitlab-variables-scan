package main

import (
	"errors"
	"log"
	"strings"

	"github.com/xanzy/go-gitlab"
)

type Scanner struct {
	debug  bool
	c      *Config
	git    *Git
	failed bool
}

func NewScanner(c *Config, debug bool) (*Scanner, error) {
	git, err := NewGit(c.Endpoint, c.Token)
	if err != nil {
		return nil, err
	}
	return &Scanner{
		debug: debug,
		c:     c,
		git:   git,
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
	err = s.checkGroupsVaribles(groups)
	if err != nil {
		return err
	}

	projects, err := s.fetchProjects()
	if err != nil {
		return err
	}
	err = s.checkProjectsVaribles(projects)
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
	if s.debug {
		log.Println(strings.Join(groupNames, ", "))
		log.Printf("Found %d group(s)", len(groups))
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
	if s.debug {
		log.Println(strings.Join(projectNames, ", "))
		log.Printf("Found %d project(s)", len(projects))
	}
	return projects, nil
}

func (s *Scanner) checkProjectsVaribles(projects []*gitlab.Project) error {
	for _, project := range projects {
		log.Printf("Checking %s project...", project.NameWithNamespace)
		vars, err := s.git.GetProjectVariables(project.ID)
		if err != nil {
			return err
		}
		if s.debug {
			log.Printf("Found %d variable(s)", len(vars))
		}
		isContains := s.IsVariablesContainsSensitiveData(vars)
		if isContains {
			s.failed = true
		}
	}
	return nil
}

func (s *Scanner) checkGroupsVaribles(groups []*gitlab.Group) error {
	for _, group := range groups {
		log.Printf("Checking %s group...", group.FullName)
		vars, err := s.git.GetGroupVariables(group.ID)
		if err != nil {
			return err
		}
		if s.debug {
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
		value := strings.Replace(variable.Value, "\n", "", -1)
		match := false
		for _, rule := range s.c.VariablesRE {
			if rule.MatchString(variable.Key) {
				match = true
				log.Printf("  * %s=%s [by name]", variable.Key, value)
				break
			}
		}
		if match {
			contains = true
			continue
		}
		for _, rule := range s.c.ValuesRE {
			if rule.MatchString(variable.Value) {
				match = true
				log.Printf("  * %s=%s [by value]", variable.Key, value)
				break
			}
		}
		if match {
			contains = true
			continue
		}
	}
	return contains
}
