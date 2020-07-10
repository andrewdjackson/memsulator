package scenarios

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strings"
	"time"

	"github.com/andrewdjackson/memsulator/utils"
	"github.com/gocarina/gocsv"
)

// MemsData is the mems information computed from dataframes 0x80 and 0x7d
type MemsData struct {
	Time                     string  `csv:"#time"`
	EngineRPM                uint16  `csv:"80x01-02_engine-rpm"`
	CoolantTemp              int     `csv:"80x03_coolant_temp"`
	AmbientTemp              int     `csv:"80x04_ambient_temp"`
	IntakeAirTemp            int     `csv:"80x05_intake_air_temp"`
	FuelTemp                 int     `csv:"80x06_fuel_temp"`
	ManifoldAbsolutePressure float32 `csv:"80x07_map_kpa"`
	BatteryVoltage           float32 `csv:"80x08_battery_voltage"`
	ThrottlePotSensor        float32 `csv:"80x09_throttle_pot"`
	IdleSwitch               bool    `csv:"80x0A_idle_switch"`
	AirconSwitch             bool    `csv:"80x0B_uk1"`
	ParkNeutralSwitch        bool    `csv:"80x0C_park_neutral_switch"`
	DTC0                     int     `csv:"80x0D-0E_fault_codes"`
	DTC1                     int     `csv:"-"`
	IdleSetPoint             int     `csv:"80x0F_idle_set_point"`
	IdleHot                  int     `csv:"80x10_idle_hot"`
	Uk8011                   int     `csv:"80x11_uk2"`
	IACPosition              int     `csv:"80x12_iac_position"`
	IdleSpeedDeviation       uint16  `csv:"80x13-14_idle_error"`
	IgnitionAdvanceOffset80  int     `csv:"80x15_ignition_advance_offset"`
	IgnitionAdvance          float32 `csv:"80x16_ignition_advance"`
	CoilTime                 float32 `csv:"80x17-18_coil_time"`
	CrankshaftPositionSensor bool    `csv:"80x19_crankshaft_position_sensor"`
	Uk801a                   int     `csv:"80x1A_uk4"`
	Uk801b                   int     `csv:"80x1B_uk5"`
	IgnitionSwitch           bool    `csv:"7dx01_ignition_switch"`
	ThrottleAngle            int     `csv:"7dx02_throttle_angle"`
	Uk7d03                   int     `csv:"7dx03_uk6"`
	AirFuelRatio             float32 `csv:"7dx04_air_fuel_ratio"`
	DTC2                     int     `csv:"7dx05_dtc2"`
	LambdaVoltage            int     `csv:"7dx06_lambda_voltage"`
	LambdaFrequency          int     `csv:"7dx07_lambda_sensor_frequency"`
	LambdaDutycycle          int     `csv:"7dx08_lambda_sensor_dutycycle"`
	LambdaStatus             int     `csv:"7dx09_lambda_sensor_status"`
	ClosedLoop               bool    `csv:"7dx0A_closed_loop"`
	LongTermFuelTrim         int     `csv:"7dx0B_long_term_fuel_trim"`
	ShortTermFuelTrim        int     `csv:"7dx0C_short_term_fuel_trim"`
	CarbonCanisterPurgeValve int     `csv:"7dx0D_carbon_canister_dutycycle"`
	DTC3                     int     `csv:"7dx0E_dtc3"`
	IdleBasePosition         int     `csv:"7dx0F_idle_base_pos"`
	Uk7d10                   int     `csv:"7dx10_uk7"`
	DTC4                     int     `csv:"7dx11_dtc4"`
	IgnitionAdvanceOffset7d  int     `csv:"7dx12_ignition_advance2"`
	IdleSpeedOffset          int     `csv:"7dx13_idle_speed_offset"`
	Uk7d14                   int     `csv:"7dx14_idle_error2"`
	Uk7d15                   int     `csv:"7dx14-15_uk10"`
	DTC5                     int     `csv:"7dx16_dtc5"`
	Uk7d17                   int     `csv:"7dx17_uk11"`
	Uk7d18                   int     `csv:"7dx18_uk12"`
	Uk7d19                   int     `csv:"7dx19_uk13"`
	Uk7d1a                   int     `csv:"7dx1A_uk14"`
	Uk7d1b                   int     `csv:"7dx1B_uk15"`
	Uk7d1c                   int     `csv:"7dx1C_uk16"`
	Uk7d1d                   int     `csv:"7dx1D_uk17"`
	Uk7d1e                   int     `csv:"7dx1E_uk18"`
	JackCount                int     `csv:"7dx1F_uk19"`
	Dataframe7d              string  `csv:"0x7d_raw"`
	Dataframe80              string  `csv:"0x80_raw"`
}

// Scenario represents the scenario data
type Scenario struct {
	file *os.File
	// Memsdata log
	Memsdata []*MemsData
	// Position in the log
	Position int
	// Count of items in the log
	Count int
}

// NewScenario creates a new scenario
func NewScenario() *Scenario {
	scenario := &Scenario{}
	// initialise the log
	scenario.Memsdata = []*MemsData{}
	// start at the beginning
	scenario.Position = 0
	// no items in the log
	scenario.Count = 0

	return scenario
}

// Open the CSV scenario file
func (scenario *Scenario) openFile(filepath string) {
	var err error

	scenario.file, err = os.OpenFile(filepath, os.O_RDWR|os.O_CREATE, os.ModePerm)

	if err != nil {
		utils.LogE.Printf("unable to open %s", err)
	}
}

// Load the scenario
func (scenario *Scenario) Load(filepath string) {
	scenario.openFile(filepath)

	if err := gocsv.Unmarshal(scenario.file, &scenario.Memsdata); err != nil {
		utils.LogE.Printf("unable to parse file %s", err)
	} else {
		scenario.Count = len(scenario.Memsdata)
		utils.LogI.Printf("loaded scenario %s (%d dataframes)", filepath, scenario.Count)
	}
}

// Next provides the next item in the log
func (scenario *Scenario) Next() *MemsData {
	item := scenario.Memsdata[scenario.Position]
	scenario.Position = scenario.Position + 1

	// if we pass the end, loop back to the start
	if scenario.Position >= scenario.Count {
		utils.LogW.Printf("reached end of scenario, restarting from beginning")
		scenario.Position = 0
	}

	return item
}

// ConvertReadMemsLogToMemsFCR takes Readmems Log files and converts them into MemsFCR format
func (scenario *Scenario) ConvertReadMemsLogToMemsFCR(filepath string) {
	var m *MemsData
	lines := scenario.readResponseFile(filepath)
	startTime := time.Now()

	// convert to a compress byte string line by line
	// from 80: 00 00...
	// to 800000...
	// skip header lines
	for i := 3; i < len(lines); i++ {
		line := lines[i]

		if scenario.isACommandResponse(line) {
			line = scenario.cleanCommandResponse(line)

			if strings.HasPrefix(line, "80") {
				m = &MemsData{}
				m.Dataframe80 = line
			}
			if strings.HasPrefix(line, "7D") {
				m.Dataframe7d = line
			}

			// convert to a struct, response order is 80, 7d so
			// make sure we don't increment until we have both dataframes
			if m.Dataframe80 != "" && m.Dataframe7d != "" {
				startTime = startTime.Add(1 * time.Second)
				m.Time = startTime.Format("15:04:05")
				scenario.Memsdata = append(scenario.Memsdata, m)
			}
		}
	}

}

// SaveCSVFile saves the Memdata to a CSV file
func (scenario *Scenario) SaveCSVFile(filepath string) {
	file, _ := os.OpenFile(filepath, os.O_RDWR|os.O_CREATE, os.ModePerm)

	err := gocsv.MarshalFile(&scenario.Memsdata, file)
	if err != nil {
		utils.LogI.Printf("error saving csv file %s", err)
	}
}

// returns true if the line starts with an 80 or 75 which
// indicate command responses
func (scenario *Scenario) isACommandResponse(line string) bool {
	if strings.HasPrefix(line, "80") {
		return true
	}

	if strings.HasPrefix(line, "7D") {
		return true
	}

	return false
}

// returns compacted and cleaned response
func (scenario *Scenario) cleanCommandResponse(line string) string {
	// remove the : character
	line = strings.ReplaceAll(line, ":", "")
	// remove all the spaces
	line = strings.ReplaceAll(line, " ", "")
	// remove all the LF
	line = strings.ReplaceAll(line, "\n", "")
	line = strings.ReplaceAll(line, "\r", "")

	return line
}

// ReadMems logs are piped output from the readmems console and
// in the format:
// 7D: 00 00 ...
// 80: 00 00 ...
// this function reads the response file into an array of strings for
// processing
func (scenario *Scenario) readResponseFile(path string) []string {
	file, err := os.Open(path)

	if err != nil {
		utils.LogE.Printf("unable to open %s", err)
	}
	defer file.Close()

	var lines []string

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	return lines
}

// ConvertCSVToMemsFCR takes Readmems CSV files and converts them into MemsFCR format
func (scenario *Scenario) ConvertCSVToMemsFCR(filepath string) {
	// load the Readmems CSV
	scenario.Load(filepath)

	// recreate the Dataframes from the CSV values
	for _, m := range scenario.Memsdata {
		scenario.recreateDataframes(m)
	}
}

func convertBooltoInt(b bool) uint8 {
	if b {
		return 1
	}

	return 0
}

// Recreate the Dataframe HEX data from the parameters
// The CSV data fields are calculated from the raw data, we need to undo
// those computations
func (scenario *Scenario) recreateDataframes(data *MemsData) {
	// undo all the computations and put all data back into integer/hex format
	df80 := fmt.Sprintf("801C"+
		"%04x%02x%02x%02x%02x%02x%02x%02x%02x%02x"+
		"%02x%02x%02x%02x%02x%02x%02x%04x%02x%02x%04x"+
		"%02x%02x%02x",
		uint16(data.EngineRPM),
		uint8(data.CoolantTemp+55),
		uint8(data.AmbientTemp+55),
		uint8(data.IntakeAirTemp+55),
		uint8(data.FuelTemp+55),
		uint8(data.ManifoldAbsolutePressure),
		uint8(data.BatteryVoltage*10),
		uint8(data.ThrottlePotSensor/0.02),
		convertBooltoInt(data.IdleSwitch),
		convertBooltoInt(data.AirconSwitch),
		convertBooltoInt(data.ParkNeutralSwitch),
		uint8(data.DTC0),
		uint8(data.DTC1),
		uint8(data.IdleSetPoint),
		uint8(data.IdleHot),
		uint8(data.Uk8011),
		uint8(math.Round(float64(data.IACPosition)*1.8)),
		uint16(data.IdleSpeedDeviation),
		uint8(data.IgnitionAdvanceOffset80),
		uint8((data.IgnitionAdvance*2)+24),
		uint16(data.CoilTime/0.002),
		convertBooltoInt(data.CrankshaftPositionSensor),
		uint8(data.Uk801a),
		uint8(data.Uk801b),
	)

	df7d := fmt.Sprintf("7D20"+
		"%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x"+
		"%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x"+
		"%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x",
		convertBooltoInt(data.IgnitionSwitch),
		uint8(data.ThrottleAngle/6*10),
		uint8(data.Uk7d03),
		uint8(data.AirFuelRatio*10),
		uint8(data.DTC2),
		uint8(data.LambdaVoltage/5),
		uint8(data.LambdaFrequency),
		uint8(data.LambdaDutycycle),
		uint8(data.LambdaStatus),
		convertBooltoInt(data.ClosedLoop),
		uint8(data.LongTermFuelTrim),
		uint8(data.ShortTermFuelTrim),
		uint8(data.CarbonCanisterPurgeValve),
		uint8(data.DTC3),
		uint8(data.IdleBasePosition),
		uint8(data.Uk7d10),
		uint8(data.DTC4),
		uint8(data.IgnitionAdvanceOffset7d+48),
		uint8((data.IdleSpeedOffset+128)/25),
		uint8(data.Uk7d14),
		uint8(data.Uk7d15),
		uint8(data.DTC5),
		uint8(data.Uk7d17),
		uint8(data.Uk7d18),
		uint8(data.Uk7d19),
		uint8(data.Uk7d1a),
		uint8(data.Uk7d1b),
		uint8(data.Uk7d1c),
		uint8(data.Uk7d1d),
		uint8(data.Uk7d1e),
		uint8(data.JackCount),
	)

	data.Dataframe7d = strings.ToUpper(df7d)
	data.Dataframe80 = strings.ToUpper(df80)

	fmt.Printf("0x80: %s\n", data.Dataframe80)
	fmt.Printf("0x7d: %s\n", data.Dataframe7d)
}
