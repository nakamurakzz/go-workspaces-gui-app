package main

import (
	"encoding/json"
	"os"
)

type Settings struct {
	activeProfile string
	activeFeature activeFeature
	Profiles      []string
}

type activeFeature struct {
	EC2       bool
	WorkSpace bool
}

func NewSrttings() *Settings {
	return &Settings{
		activeProfile: "",
		Profiles:      []string{},
	}
}

func LoadSettings() (*Settings, error) {
	file, err := os.Open("settings.json")
	if err != nil {
		if os.IsNotExist(err) {
			return &Settings{
				Profiles: []string{},
			}, nil
		}
		return NewSrttings(), nil
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	settings := &Settings{}
	err = decoder.Decode(settings)
	if err != nil {
		return nil, err
	}

	return settings, nil
}

func (s *Settings) Save() error {
	file, err := os.Create("settings.json")
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	err = encoder.Encode(s)
	if err != nil {
		return err
	}

	return nil
}

func (s *Settings) SetActiveProfile(profile string) {
	s.activeProfile = profile
}

func (s *Settings) GetActiveProfile() string {
	return s.activeProfile
}

func (s *Settings) changeEC2ListScreenVisible() {
	s.activeFeature.EC2 = !s.activeFeature.EC2
}

func (s *Settings) changeWorkSpaceListScreenVisible() {
	s.activeFeature.WorkSpace = !s.activeFeature.WorkSpace
}

func (s *Settings) GetActiveFeature() activeFeature {
	return s.activeFeature
}
