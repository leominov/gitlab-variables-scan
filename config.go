package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"regexp"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Insecure bool     `yaml:"-"`
	Debug    bool     `yaml:"debug"`
	Endpoint string   `yaml:"endpoint"`
	Token    string   `yaml:"token"`
	GroupIDs []int    `yaml:"groupIDs"`
	KeysRaw  []string `yaml:"keys"`
	Include  Filters  `yaml:"include"`
	Exclude  Filters  `yaml:"exclude"`
}

type Filters struct {
	KeysRaw   []string         `yaml:"keys"`
	Keys      []*regexp.Regexp `yaml:"-"`
	ValuesRaw []string         `yaml:"values"`
	Values    []*regexp.Regexp `yaml:"-"`
	PairsRaw  []string         `yaml:"pairs"`
	Pairs     []*regexp.Regexp `yaml:"-"`
}

func LoadFromFile(filename string) (*Config, error) {
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	c := &Config{}
	err = yaml.Unmarshal(b, &c)
	if err != nil {
		return nil, err
	}
	err = c.fillFromEnv()
	if err != nil {
		return nil, err
	}
	err = c.parseRawData()
	if err != nil {
		return nil, err
	}
	return c, nil
}

func (c *Config) fillFromEnv() error {
	if len(c.Endpoint) == 0 {
		c.Endpoint = os.Getenv("GITLAB_ADDRESS")
	}
	if len(c.Token) == 0 {
		c.Token = os.Getenv("GITLAB_TOKEN")
	}
	return nil
}

func (c *Config) parseRawData() error {
	if err := c.Exclude.Parse(); err != nil {
		return err
	}
	if err := c.Include.Parse(); err != nil {
		return err
	}
	return nil
}

func (f *Filters) Parse() error {
	for _, variable := range f.KeysRaw {
		re, err := regexp.Compile(variable)
		if err != nil {
			return fmt.Errorf("Failed to parse %s regexp", variable)
		}
		f.Keys = append(f.Keys, re)
	}
	for _, value := range f.ValuesRaw {
		re, err := regexp.Compile(value)
		if err != nil {
			return fmt.Errorf("Failed to parse %s regexp", value)
		}
		f.Values = append(f.Values, re)
	}
	for _, pair := range f.PairsRaw {
		re, err := regexp.Compile(pair)
		if err != nil {
			return fmt.Errorf("Failed to parse %s regexp", pair)
		}
		f.Pairs = append(f.Pairs, re)
	}
	return nil
}
