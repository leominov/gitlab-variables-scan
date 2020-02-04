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
	entries := []*gitlab.Group{}
	options := &gitlab.ListSubgroupsOptions{
		AllAvailable: gitlab.Bool(true),
	}
	err := g.withPagination(func(opts gitlab.ListOptions) (*gitlab.Response, error) {
		options.ListOptions = opts
		fetchedEntries, r, err := g.cli.Groups.ListSubgroups(groupID, options, nil)
		if err != nil {
			return nil, err
		}
		entries = append(entries, fetchedEntries...)
		return r, nil
	})
	return entries, err
}

func (g *Git) GetProjects(groupID int) ([]*gitlab.Project, error) {
	entries := []*gitlab.Project{}
	options := &gitlab.ListGroupProjectsOptions{
		IncludeSubgroups: gitlab.Bool(true),
	}
	err := g.withPagination(func(opts gitlab.ListOptions) (*gitlab.Response, error) {
		options.ListOptions = opts
		fetchedEntries, r, err := g.cli.Groups.ListGroupProjects(groupID, options, nil)
		if err != nil {
			return nil, err
		}
		entries = append(entries, fetchedEntries...)
		return r, nil
	})
	return entries, err
}

func (g *Git) GetProjectVariables(projectID int) ([]*Variable, error) {
	entries := []*Variable{}
	err := g.withPagination(func(opts gitlab.ListOptions) (*gitlab.Response, error) {
		options := gitlab.ListVariablesOptions(opts)
		fetchedEntries, r, err := g.cli.ProjectVariables.ListVariables(projectID, &options, nil)
		if err != nil {
			return nil, err
		}
		for _, entry := range fetchedEntries {
			entries = append(entries, &Variable{
				Key:   entry.Key,
				Value: entry.Value,
			})
		}
		return r, nil
	})
	return entries, err
}

func (g *Git) GetGroupVariables(groupID int) ([]*Variable, error) {
	entries := []*Variable{}
	err := g.withPagination(func(opts gitlab.ListOptions) (*gitlab.Response, error) {
		options := gitlab.ListVariablesOptions(opts)
		fetchedEntries, r, err := g.cli.GroupVariables.ListVariables(groupID, &options, nil)
		if err != nil {
			return nil, err
		}
		for _, entry := range fetchedEntries {
			entries = append(entries, &Variable{
				Key:   entry.Key,
				Value: entry.Value,
			})
		}
		return r, nil
	})
	return entries, err
}

func (g *Git) withPagination(fetch func(opts gitlab.ListOptions) (*gitlab.Response, error)) error {
	page := 1
	for {
		opts := gitlab.ListOptions{
			Page: page,
		}
		r, err := fetch(opts)
		if err != nil {
			return err
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
	return nil
}

func (g *Git) GetGroup(groupID int) (*gitlab.Group, error) {
	group, _, err := g.cli.Groups.GetGroup(groupID)
	if err != nil {
		return nil, err
	}
	return group, nil
}
