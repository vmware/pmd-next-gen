package configfile

import (
	"errors"
	"strings"

	"github.com/go-ini/ini"
)

type Meta struct {
	Path string
	cfg  *ini.File
}

func Load(path string) (*Meta, error) {
	cfg, err := ini.LoadSources(ini.LoadOptions{AllowNonUniqueSections: true, AllowDuplicateShadowValues: true}, path)
	if err != nil {
		return nil, err
	}

	return &Meta{
		Path: path,
		cfg:  cfg,
	}, nil
}

func (m *Meta) Save() error {
	return m.cfg.SaveTo(m.Path)
}

func ParseKeyFromSectionString(path string, section string, key string) (string, error) {
	c, err := Load(path)
	if err != nil {
		return "", err
	}

	v := c.cfg.Section(section).Key(key).String()
	if v == "" {
		return "", errors.New("not found")
	}

	return v, nil
}

func (m *Meta) SetKeySectionString(section string, key string, value string) {
	m.cfg.Section(section).Key(key).SetValue(strings.ToLower(value))
}

func MapTo(cfg *ini.File, section string, v interface{}) error {
	if err := cfg.Section(section).MapTo(v); err != nil {
		return err
	}

	return nil
}
