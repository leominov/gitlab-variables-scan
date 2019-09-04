package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"regexp"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Insecure  bool             `yaml:"-"`
	Debug     bool             `yaml:"debug"`
	Endpoint  string           `yaml:"endpoint"`
	Token     string           `yaml:"token"`
	GroupIDs  []int            `yaml:"groupIDs"`
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
	for _, variable := range c.KeysRaw {
		re, err := regexp.Compile(variable)
		if err != nil {
			return fmt.Errorf("Failed to parse %s regexp", variable)
		}
		c.Keys = append(c.Keys, re)
	}
	for _, value := range c.ValuesRaw {
		re, err := regexp.Compile(value)
		if err != nil {
			return fmt.Errorf("Failed to parse %s regexp", value)
		}
		c.Values = append(c.Values, re)
	}
	for _, pair := range c.PairsRaw {
		re, err := regexp.Compile(pair)
		if err != nil {
			return fmt.Errorf("Failed to parse %s regexp", pair)
		}
		c.Pairs = append(c.Pairs, re)
	}
	return nil
}
