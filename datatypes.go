package autocomplete

// AVSPlace represents response from aviasales autocomplete service
type AVSPlace struct {
	CityName     string            `json:"city_name"`
	Type         string            `json:"type"`
	StateCode    string            `json:"state_code"`
	Cases        map[string]string `json:"cases"`
	CityCases    map[string]string `json:"city_cases"`
	CountryCases map[string]string `json:"country_cases"`
	Name         string            `json:"name"`
	CountryCode  string            `json:"country_code"`
	Code         string            `json:"code"`
	CityCode     string            `json:"city_code"`
	CountryName  string            `json:"country_name"`
	Coordinates  *Coordinate       `json:"coordinates"`
	Weight       int64             `json:"weight"`
	IndexStrings []string          `json:"index_strings"`
}

// Coordinate contains geographic coordinates
type Coordinate struct {
	Lat float64 `json:"lat"`
	Lon float64 `json:"lon"`
}

// WidgetData is desired data format for stub service
type WidgetData struct {
	Slug     string `json:"slug"`
	Title    string `json:"title"`
	Subtitle string `json:"subtitle"`
}

// CreateWidgetData converts AVSPlace to WidgetData basing on place type
func CreateWidgetData(p *AVSPlace) *WidgetData {
	switch p.Type {
	case "country":
		return &WidgetData{
			Slug:     p.Code,
			Title:    p.Name,
			Subtitle: p.Name,
		}
	case "airport":
		return &WidgetData{
			Slug:     p.Code,
			Title:    p.Name,
			Subtitle: p.CityName,
		}
	case "city":
		return &WidgetData{
			Slug:     p.Code,
			Title:    p.Name,
			Subtitle: p.CountryName,
		}
	default:
		return &WidgetData{
			Slug: p.Code,
		}
	}
}
