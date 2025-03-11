package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v3"
)

type NodeLabelMatch struct {
	MatchType  string `yaml:"matchType"`
	MatchValue string `yaml:"matchValue"`
}

type GPhotosLocale struct {
	SelectAllPhotosLabel            NodeLabelMatch `yaml:"selectAllPhotosLabel"`
	FileNameLabel                   NodeLabelMatch `yaml:"fileNameLabel"`
	DateLabel                       NodeLabelMatch `yaml:"dateLabel"`
	Today                           string         `yaml:"today"`
	Yesterday                       string         `yaml:"yesterday"`
	TimeLabel                       NodeLabelMatch `yaml:"timeLabel"`
	TzLabel                         NodeLabelMatch `yaml:"tzLabel"`
	ViewPreviousPhotoMatch          NodeLabelMatch `yaml:"viewPreviousPhotoMatch"`
	MoreOptionsLabel                NodeLabelMatch `yaml:"moreOptionsLabel"`
	DownloadLabel                   NodeLabelMatch `yaml:"downloadLabel"`
	DownloadOriginalLabel           NodeLabelMatch `yaml:"downloadOriginalLabel"`
	OpenInfoMatch                   NodeLabelMatch `yaml:"openInfoMatch"`
	VideoStillProcessingDialogLabel NodeLabelMatch `yaml:"videoStillProcessingDialogLabel"`
	VideoStillProcessingStatusText  string         `yaml:"videoStillProcessingStatusText"`
	NoWebpageFoundText              string         `yaml:"noWebpageFoundText"`
	ShortDayNames                   []string       `yaml:"shortDayNames"`
	LongDayNames                    []string       `yaml:"longDayNames"`
	ShortMonthNames                 []string       `yaml:"shortMonthNames"`
}

var locales map[string]GPhotosLocale = make(map[string]GPhotosLocale)

func readLocalesFromYAML() error {
	// If no filename is provided, use a default
	filename := os.Getenv("GPHOTOS_LOCALE_FILE")
	if filename == "" {
		filename = filepath.Join(filepath.Dir(os.Args[0]), "locales.yaml")
	}

	// Read the YAML file
	data, err := os.ReadFile(filename)
	if os.IsNotExist(err) && os.Getenv("GPHOTOS_LOCALE_FILE") == "" {
		return nil
	} else if err != nil {
		return fmt.Errorf("error reading locales YAML file: %w", err)
	}

	// Parse the YAML
	var parsedLocales map[string]GPhotosLocale
	err = yaml.Unmarshal(data, &parsedLocales)
	if err != nil {
		return fmt.Errorf("error parsing locales YAML: %w", err)
	}

	localeNames := make([]string, 0, len(parsedLocales))
	for name, locale := range parsedLocales {
		locales[name] = locale
		localeNames = append(localeNames, name)
		log.Debug().Msgf("Loaded locale %s with values\n%v", name, locale)
	}
	log.Info().Msgf("Loaded locales: %s", strings.Join(localeNames, ", "))

	return nil
}

func initLocales() error {
	// Get locale file path from env var
	if err := readLocalesFromYAML(); err != nil {
		return fmt.Errorf("error reading locales YAML file: %w", err)
	}

	// If no English locale exists, set the default
	if _, exists := locales["en"]; !exists {
		locales["en"] = GPhotosLocale{
			SelectAllPhotosLabel:            NodeLabelMatch{"startsWith", "Select all photos from"},
			FileNameLabel:                   NodeLabelMatch{"startsWith", "Filename:"},
			DateLabel:                       NodeLabelMatch{"startsWith", "Date taken:"},
			Today:                           "Today",
			Yesterday:                       "Yesterday",
			TimeLabel:                       NodeLabelMatch{"startsWith", "Time taken:"},
			TzLabel:                         NodeLabelMatch{"startsWith", "GMT"},
			ViewPreviousPhotoMatch:          NodeLabelMatch{"equals", "View previous photo"},
			MoreOptionsLabel:                NodeLabelMatch{"equals", "More options"},
			DownloadLabel:                   NodeLabelMatch{"equals", "Download - Shift+D"},
			DownloadOriginalLabel:           NodeLabelMatch{"equals", "Download original"},
			OpenInfoMatch:                   NodeLabelMatch{"equals", "Open info"},
			VideoStillProcessingDialogLabel: NodeLabelMatch{"startsWith", "Video still is processing"},
			VideoStillProcessingStatusText:  "Video is still processing &amp; can be downloaded later",
			NoWebpageFoundText:              "No webpage was found for the web address:",
			ShortDayNames:                   []string{"Sun", "Mon", "Tue", "Wed", "Thu", "Fri", "Sat"},
			LongDayNames:                    []string{"Sunday", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday"},
			ShortMonthNames:                 []string{"Jan", "Feb", "Mar", "Apr", "May", "Jun", "Jul", "Aug", "Sep", "Oct", "Nov", "Dec"},
		}
	}

	if _, exists := locales["nl"]; !exists {
		locales["nl"] = GPhotosLocale{
			SelectAllPhotosLabel:            NodeLabelMatch{"startsWith", "Alle foto's van"},
			FileNameLabel:                   NodeLabelMatch{"startsWith", "Bestandsnaam:"},
			DateLabel:                       NodeLabelMatch{"startsWith", "Fotodatum:"},
			Today:                           "Vandaag",
			Yesterday:                       "Gisteren",
			TimeLabel:                       NodeLabelMatch{"startsWith", "Tijdsduur:"},
			TzLabel:                         NodeLabelMatch{"startsWith", "GMT"},
			ViewPreviousPhotoMatch:          NodeLabelMatch{"equals", "Vorige foto bekijken"},
			MoreOptionsLabel:                NodeLabelMatch{"equals", "Meer opties"},
			DownloadLabel:                   NodeLabelMatch{"equals", "Downloaden - Shift+D"},
			DownloadOriginalLabel:           NodeLabelMatch{"equals", "Download original"},
			OpenInfoMatch:                   NodeLabelMatch{"equals", "Info openen"},
			VideoStillProcessingDialogLabel: NodeLabelMatch{"startsWith", "Video still is processing"},
			VideoStillProcessingStatusText:  "Video is still processing &amp; can be downloaded later",
			NoWebpageFoundText:              "No webpage was found for the web address:",
			ShortDayNames:                   []string{"zo", "ma", "di", "wo", "do", "vr", "za"},
			LongDayNames:                    []string{"zondag", "maandag", "dinsdag", "woensdag", "donderdag", "vrijdag", "zaterdag"},
			ShortMonthNames:                 []string{"jan", "feb", "mrt", "apr", "mei", "jun", "jul", "aug", "sep", "okt", "nov", "dec"},
		}
	}

	return nil
}

func getAriaLabelSelector(matcher NodeLabelMatch) string {
	eq := "="
	if matcher.MatchType == "equals" {
		eq = "="
	} else if matcher.MatchType == "startsWith" {
		eq = "^="
	} else if matcher.MatchType == "contains" {
		eq = "*="
	} else if matcher.MatchType == "endsWith" {
		eq = "$="
	}
	return fmt.Sprintf("[aria-label%s\"%s\"]", eq, matcher.MatchValue)
}
