package loader

import (
	"encoding/hex"
	"github.com/gocarina/gocsv"
	log "github.com/sirupsen/logrus"
	"os"
)

type CSVFields struct {
	Dataframe7d string `csv:"0x7d_raw"`
	Dataframe80 string `csv:"0x80_raw"`
}

type Dataframes struct {
	Dataframe7d []byte
	Dataframe80 []byte
}

type Scenario struct {
	file       *os.File
	csvFields  []*CSVFields
	dataframes []*Dataframes
	Count      int
}

func NewScenario() *Scenario {
	return &Scenario{
		csvFields:  []*CSVFields{},
		dataframes: []*Dataframes{},
		Count:      0,
	}
}

func (scenario *Scenario) Load(filepath string) {
	if _, err := os.Stat(filepath); err != nil {
		log.Errorf("unable to find file %s (%e)", filepath, err)
		return
	}

	scenario.openFile(filepath)

	if err := gocsv.Unmarshal(scenario.file, &scenario.csvFields); err != nil {
		log.Errorf("unable to parse file %s", err)
	} else {
		scenario.Count = len(scenario.csvFields)
		if err = scenario.convertFieldsToDataframes(); err == nil {
			log.Infof("loaded scenario %s (%d dataframes)", filepath, scenario.Count)
		}
	}
}

func (scenario *Scenario) openFile(filepath string) {
	var err error

	if scenario.file, err = os.OpenFile(filepath, os.O_RDWR|os.O_CREATE, os.ModePerm); err != nil {
		log.Errorf("unable to open %s", err)
	}
}

func (scenario *Scenario) convertFieldsToDataframes() error {
	var err error
	var d7d, d80 []byte

	for i := 0; i < len(scenario.csvFields); i++ {
		if d7d, err = convertHexStringToByteArray(scenario.csvFields[i].Dataframe7d); err == nil {
			if d80, err = convertHexStringToByteArray(scenario.csvFields[i].Dataframe80); err == nil {

				dataframe := &Dataframes{
					Dataframe7d: d7d,
					Dataframe80: d80,
				}

				scenario.dataframes = append(scenario.dataframes, dataframe)
			}
		}
	}

	return err
}

func convertHexStringToByteArray(byteString string) ([]byte, error) {
	return hex.DecodeString(byteString)
}
