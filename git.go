package main

import (
	"strconv"

	"github.com/xanzy/go-gitlab"
)

type Git struct {
	cli *gitlab.Client
}

type Variable struct {
	Key   string
	Value string
}

func NewGit(endpoint, token string) (*Git, error) {
	cli := gitlab.NewClient(nil, token)
	err := cli.SetBaseURL(endpoint)
	if err != nil {
		return nil, err
	}
	return &Git{
		cli: cli,
	}, err
}

func (g *Git) GetSubGroups(groupID int) ([]*gitlab.Group, error) {
	groups, err := g.getSubGroups(groupID)
	if err != nil {
		return nil, err
	}
	if len(groups) == 0 {
		return groups, nil
	}
	for _, group := range groups {
		grps, err := g.GetSubGroups(group.ID)
		if err != nil {
			return nil, err
		}
		groups = append(groups, grps...)
	}
	return groups, nil
}

func (g *Git) getSubGroups(groupID int) ([]*gitlab.Group, error) {
	page := 1
	entries := []*gitlab.Group{}
	for {
		options := &gitlab.ListSubgroupsOptions{
			ListOptions: gitlab.ListOptions{
				Page: page,
			},
			AllAvailable: gitlab.Bool(true),
		}
		fetchedEntries, r, err := g.cli.Groups.ListSubgroups(groupID, options, nil)
		if err != nil {
			return nil, err
		}
		for _, entry := range fetchedEntries {
			entries = append(entries, entry)
		}
		nextPageRaw := r.Header.Get("X-Next-Page")
		if len(nextPageRaw) == 0 {
			break
		}
		nextPage, err := strconv.Atoi(nextPageRaw)
		if err != nil {
			break
		}
		page = nextPage
	}
	return entries, nil
}

func (g *Git) GetProjects(groupID int) ([]*gitlab.Project, error) {
	page := 1
	entries := []*gitlab.Project{}
	for {
		options := &gitlab.ListGroupProjectsOptions{
			ListOptions: gitlab.ListOptions{
				Page: page,
			},
			IncludeSubgroups: gitlab.Bool(true),
		}
		fetchedEntries, r, err := g.cli.Groups.ListGroupProjects(groupID, options, nil)
		if err != nil {
			return nil, err
		}
		for _, entry := range fetchedEntries {
			entries = append(entries, entry)
		}
		nextPageRaw := r.Header.Get("X-Next-Page")
		if len(nextPageRaw) == 0 {
			break
		}
		nextPage, err := strconv.Atoi(nextPageRaw)
		if err != nil {
			break
		}
		page = nextPage
	}
	return entries, nil
}

func (g *Git) GetProjectVariables(projectID int) ([]*Variable, error) {
	page := 1
	entries := []*Variable{}
	for {
		options := &gitlab.ListVariablesOptions{
			ListOptions: gitlab.ListOptions{
				Page: page,
			},
		}
		fetchedEntries, r, err := g.cli.ProjectVariables.ListVariables(projectID, options, nil)
		if err != nil {
			return nil, err
		}
		for _, entry := range fetchedEntries {
			entries = append(entries, &Variable{
				Key:   entry.Key,
				Value: entry.Value,
			})
		}
		nextPageRaw := r.Header.Get("X-Next-Page")
		if len(nextPageRaw) == 0 {
			break
		}
		nextPage, err := strconv.Atoi(nextPageRaw)
		if err != nil {
			break
		}
		page = nextPage
	}
	return entries, nil
}

func (g *Git) GetGroupVariables(groupID int) ([]*Variable, error) {
	page := 1
	entries := []*Variable{}
	for {
		options := &gitlab.ListVariablesOptions{
			ListOptions: gitlab.ListOptions{
				Page: page,
			},
		}
		fetchedEntries, r, err := g.cli.GroupVariables.ListVariables(groupID, options, nil)
		if err != nil {
			return nil, err
		}
		for _, entry := range fetchedEntries {
			entries = append(entries, &Variable{
				Key:   entry.Key,
				Value: entry.Value,
			})
		}
		nextPageRaw := r.Header.Get("X-Next-Page")
		if len(nextPageRaw) == 0 {
			break
		}
		nextPage, err := strconv.Atoi(nextPageRaw)
		if err != nil {
			break
		}
		page = nextPage
	}
	return entries, nil
}

func (g *Git) GetGroup(groupID int) (*gitlab.Group, error) {
	group, _, err := g.cli.Groups.GetGroup(groupID)
	if err != nil {
		return nil, err
	}
	return group, nil
}
