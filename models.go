package main

import "gopkg.in/yaml.v3"

type object struct {
	APIVersion string                 `yaml:"apiVersion"`
	Kind       string                 `yaml:"kind"`
	Metadata   map[string]interface{} `yaml:"metadata"`
	Spec       yaml.Node              `yaml:"spec"`
}

const (
	kindDashboard string = "Dashboard"
	kindLabel     string = "Label"
	kindTask      string = "Task"
)

type dashboard struct {
	APIVersion string                 `yaml:"apiVersion"`
	Kind       string                 `yaml:"kind"`
	Metadata   map[string]interface{} `yaml:"metadata"`
	Spec       dashboardSpec          `yaml:"spec"`
}

type dashboardSpec struct {
	Name         string        `yaml:"name"`
	Description  string        `yaml:"description"`
	Associations []association `yaml:"associations"`
	Charts       []chart       `yaml:"charts"`
}

type association struct {
	Kind string `yaml:"kind"`
	Name string `yaml:"name"`
}

type chart struct {
	Kind                       string        `yaml:"kind,omitempty"`
	Name                       string        `yaml:"name,omitempty"`
	Prefix                     string        `yaml:"prefix,omitempty"`
	TickPrefix                 string        `yaml:"tickPrefix,omitempty"`
	Suffix                     string        `yaml:"suffix,omitempty"`
	TickSuffix                 string        `yaml:"tickSuffix,omitempty"`
	Note                       string        `yaml:"note,omitempty"`
	NoteOnEmpty                bool          `yaml:"noteOnEmpty,omitempty"`
	DecimalPlaces              int           `yaml:"decimalPlaces,omitempty"`
	EnforceDecimals            bool          `yaml:"enforceDecimals,omitempty"`
	Shade                      bool          `yaml:"shade,omitempty"`
	HoverDimension             string        `yaml:"hoverDimension,omitempty"`
	Legend                     legend        `yaml:"legend,omitempty"`
	Colors                     []color       `yaml:"colors,omitempty"`
	Queries                    []query       `yaml:"queries,omitempty"`
	Axes                       []axis        `yaml:"axes,omitempty"`
	Geom                       string        `yaml:"geom,omitempty"`
	YSeriesColumns             []string      `yaml:"ySeriesColumns,omitempty"`
	XCol                       string        `yaml:"xCol,omitempty"`
	YCol                       string        `yaml:"yCol,omitempty"`
	GenerateXAxisTicks         []string      `yaml:"generateXAxisTicks,omitempty"`
	GenerateYAxisTicks         []string      `yaml:"generateYAxisTicks,omitempty"`
	XTotalTicks                int           `yaml:"xTotalTicks,omitempty"`
	YTotalTicks                int           `yaml:"yTotalTicks,omitempty"`
	XTickStart                 float64       `yaml:"xTickStart,omitempty"`
	YTickStart                 float64       `yaml:"yTickStart,omitempty"`
	XTickStep                  float64       `yaml:"xTickStep,omitempty"`
	YTickStep                  float64       `yaml:"yTickStep,omitempty"`
	UpperColumn                string        `yaml:"upperColumn,omitempty"`
	MainColumn                 string        `yaml:"mainColumn,omitempty"`
	LowerColumn                string        `yaml:"lowerColumn,omitempty"`
	XPos                       int           `yaml:"xPos,omitempty"`
	YPos                       int           `yaml:"yPos,omitempty"`
	Height                     int           `yaml:"height,omitempty"`
	Width                      int           `yaml:"width,omitempty"`
	BinSize                    int           `yaml:"binSize,omitempty"`
	BinCount                   int           `yaml:"binCount,omitempty"`
	Position                   string        `yaml:"position,omitempty"`
	FieldOptions               []fieldOption `yaml:"fieldOptions,omitempty"`
	FillColumns                []string      `yaml:"fillColumns,omitempty"`
	TableOptions               tableOptions  `yaml:"tableOptions,omitempty"`
	TimeFormat                 string        `yaml:"timeFormat,omitempty"`
	LegendColorizeRows         bool          `yaml:"legendColorize_rows,omitempty"`
	LegendOpacity              float64       `yaml:"legendOpacity,omitempty"`
	LegendOrientationThreshold int           `yaml:"legendOrientationThreshold,omitempty"`
	Zoom                       float64
	Center                     center
	MapStyle                   string
	AllowPanAndZoom            bool
	DetectCoordinateFields     bool
	GeoLayers                  []geoLayer
}

type legend struct {
	Orientation string `yaml:"orientation,omitempty"`
	Type        string `yaml:"type"`
}

type color struct {
	ID    string   `yaml:"id,omitempty"`
	Name  string   `yaml:"name,omitempty"`
	Type  string   `yaml:"type,omitempty"`
	Hex   string   `yaml:"hex,omitempty"`
	Value *float64 `yaml:"value,omitempty"`
}

type query struct {
	Query string `yaml:"query"`
}

type axis struct {
	Base   string    `yaml:"base,omitempty"`
	Label  string    `yaml:"label,omitempty"`
	Name   string    `yaml:"name,omitempty"`
	Prefix string    `yaml:"prefix,omitempty"`
	Scale  string    `yaml:"scale,omitempty"`
	Suffix string    `yaml:"suffix,omitempty"`
	Domain []float64 `yaml:"domain,omitempty"`
}

type fieldOption struct {
	FieldName   string `yaml:"field_name,omitempty"`
	DisplayName string `yaml:"display_name,omitempty"`
	Visible     bool   `yaml:"visible,omitempty"`
}

type tableOptions struct {
	VerticalTimeAxis bool   `yaml:"verticalTimeAxis,omitempty"`
	SortByField      string `yaml:"sortByField,omitempty"`
	Wrapping         string `yaml:"wrapping,omitempty"`
	FixFirstColumn   bool   `yaml:"fixFirstColumn,omitempty"`
}

type center struct {
	Lat float64 `yaml:"lat,omitempty"`
	Lon float64 `yaml:"lon,omitempty"`
}

type geoLayer struct {
	Type               string  `yaml:"type,omitempty"`
	RadiusField        string  `yaml:"radiusField,omitempty"`
	ColorField         string  `yaml:"colorField,omitempty"`
	IntensityField     string  `yaml:"intensityField,omitempty"`
	ViewColors         []color `yaml:"viewColors,omitempty"`
	Radius             int32   `yaml:"radius,omitempty"`
	Blur               int32   `yaml:"blur,omitempty"`
	RadiusDimension    *axis   `yaml:"radiusDimension,omitempty"`
	ColorDimension     *axis   `yaml:"colorDimension,omitempty"`
	IntensityDimension *axis   `yaml:"intensityDimension,omitempty"`
	InterpolateColors  bool    `yaml:"interpolateColors,omitempty"`
	TrackWidth         int32   `yaml:"trackWidth,omitempty"`
	Speed              int32   `yaml:"speed,omitempty"`
	RandomColors       bool    `yaml:"randomColors,omitempty"`
	IsClustered        bool    `yaml:"isClustered,omitempty"`
}

type task struct {
	APIVersion string                 `yaml:"apiVersion"`
	Kind       string                 `yaml:"kind"`
	Metadata   map[string]interface{} `yaml:"metadata"`
	Spec       taskSpec               `yaml:"spec"`
}

type taskSpec struct {
	Name         string        `yaml:"name,omitempty"`
	Associations []association `yaml:"associations"`
	Cron         string        `yaml:"cron,omitempty"`
	Every        string        `yaml:"every,omitempty"`
	Offset       string        `yaml:"offset,omitempty"`
	Query        string        `yaml:"query,omitempty"`
	Status       string        `yaml:"status,omitempty"`
}
