package main

import (
	"log"

	"github.com/alxand/nalo-workspace/internal/domain/models"
	"gorm.io/gorm"
)

func seedContinents(db *gorm.DB) error {
	continents := []models.Continent{
		{Name: "Africa", Code: "AF", Description: "Second-largest continent"},
		{Name: "Asia", Code: "AS", Description: "Largest continent by area and population"},
		{Name: "Europe", Code: "EU", Description: "Second-smallest continent"},
		{Name: "North America", Code: "NA", Description: "Third-largest continent"},
		{Name: "South America", Code: "SA", Description: "Fourth-largest continent"},
		{Name: "Australia", Code: "AU", Description: "Smallest continent"},
		{Name: "Antarctica", Code: "AN", Description: "Southernmost continent"},
	}

	for _, continent := range continents {
		if err := db.FirstOrCreate(&continent, models.Continent{Code: continent.Code}).Error; err != nil {
			return err
		}
		log.Printf("Seeded continent: %s", continent.Name)
	}
	return nil
}

func seedCountries(db *gorm.DB) error {
	// Get continents first
	var continents []models.Continent
	if err := db.Find(&continents).Error; err != nil {
		return err
	}

	// Create a map for easy lookup
	continentMap := make(map[string]int64)
	for _, c := range continents {
		continentMap[c.Code] = c.ID
	}

	countries := []models.Country{
		{Name: "United States", Code: "USA", ContinentID: continentMap["NA"], Description: "United States of America"},
		{Name: "Canada", Code: "CAN", ContinentID: continentMap["NA"], Description: "Canada"},
		{Name: "United Kingdom", Code: "GBR", ContinentID: continentMap["EU"], Description: "United Kingdom"},
		{Name: "Germany", Code: "DEU", ContinentID: continentMap["EU"], Description: "Germany"},
		{Name: "France", Code: "FRA", ContinentID: continentMap["EU"], Description: "France"},
		{Name: "Japan", Code: "JPN", ContinentID: continentMap["AS"], Description: "Japan"},
		{Name: "China", Code: "CHN", ContinentID: continentMap["AS"], Description: "China"},
		{Name: "India", Code: "IND", ContinentID: continentMap["AS"], Description: "India"},
		{Name: "Australia", Code: "AUS", ContinentID: continentMap["AU"], Description: "Australia"},
		{Name: "Brazil", Code: "BRA", ContinentID: continentMap["SA"], Description: "Brazil"},
		{Name: "South Africa", Code: "ZAF", ContinentID: continentMap["AF"], Description: "South Africa"},
		{Name: "Nigeria", Code: "NGA", ContinentID: continentMap["AF"], Description: "Nigeria"},
	}

	for _, country := range countries {
		if err := db.FirstOrCreate(&country, models.Country{Code: country.Code}).Error; err != nil {
			return err
		}
		log.Printf("Seeded country: %s", country.Name)
	}
	return nil
}

func seedCompanies(db *gorm.DB) error {
	// Get countries first
	var countries []models.Country
	if err := db.Find(&countries).Error; err != nil {
		return err
	}

	// Create a map for easy lookup
	countryMap := make(map[string]int64)
	for _, c := range countries {
		countryMap[c.Code] = c.ID
	}

	companies := []models.Company{
		{
			Name:      "Apple Inc.",
			Code:      "AAPL",
			CountryID: countryMap["USA"],
			Industry:  "Technology",
			Size:      "large",
			Founded:   1976,
			Website:   "https://www.apple.com",
		},
		{
			Name:      "Microsoft Corporation",
			Code:      "MSFT",
			CountryID: countryMap["USA"],
			Industry:  "Technology",
			Size:      "large",
			Founded:   1975,
			Website:   "https://www.microsoft.com",
		},
		{
			Name:      "Samsung Electronics",
			Code:      "SAMSUNG",
			CountryID: countryMap["JPN"],
			Industry:  "Technology",
			Size:      "large",
			Founded:   1938,
			Website:   "https://www.samsung.com",
		},
		{
			Name:      "Volkswagen Group",
			Code:      "VW",
			CountryID: countryMap["DEU"],
			Industry:  "Automotive",
			Size:      "large",
			Founded:   1937,
			Website:   "https://www.volkswagen.com",
		},
		{
			Name:      "Nestlé",
			Code:      "NESTLE",
			CountryID: countryMap["FRA"],
			Industry:  "Food & Beverage",
			Size:      "large",
			Founded:   1866,
			Website:   "https://www.nestle.com",
		},
		{
			Name:      "StartupXYZ",
			Code:      "STXYZ",
			CountryID: countryMap["USA"],
			Industry:  "Technology",
			Size:      "small",
			Founded:   2020,
			Website:   "https://www.startupxyz.com",
		},
	}

	for _, company := range companies {
		if err := db.FirstOrCreate(&company, models.Company{Code: company.Code}).Error; err != nil {
			return err
		}
		log.Printf("Seeded company: %s", company.Name)
	}
	return nil
}
