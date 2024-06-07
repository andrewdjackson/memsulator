package responder

// MEMS Command List
/*
STOP Actuator Test

0x00	Coolant gauge	0x00, 0x00
0x01	Fuel pump relay	0x01, 0x00
0x02	PTC (inlet manifold heater) relay	0x02, 0x00
0x03	Air Conditioning relay	0x03, 0x00
0x04	Idle solenoid
0x05	ORFCO solenoid
0x06	Pulse air valve
0x07	EGR valve
0x08	Purge valve	0x08, 0x00
0x09	O2 heater relay	0x09, 0x00
0x0A	Emissions fail lamp
0x0B	Wastegate
0x0C	Fuel used
0x0D	Fan 1
0x0E	Fan 2
0x0F	Variable valve timing

START Actuator Test

0x10	Coolant gauge	0x10, 0x00
0x11	Fuel pump relay	0x11, 0x00
0x12	PTC (inlet manifold heater) relay	0x12, 0x00
0x13	Air Conditioning relay	0x13, 0x00
0x14	Idle solenoid
0x15	ORFCO solenoid
0x16	Pulse air valve
0x17	EGR valve
0x18	Purge valve	0x18, 0x00
0x19	O2 heater relay	0x19, 0x00
0x1A	Emissions fail lamp
0x1B	Wastegate
0x1C	Fuel used
0x1D	Fan 1 relay	0x1D
0x1E	Fan 2 relay	0x1E
0x1F	Variable valve timing

Commands

0x20	Engine bay temperature warning light off	20 00
0x21	Cruise control disable relay	21 00
0x30	Engine bay temperature warning light on	30 00
0x31	Cruise control disable relay off	31 00
0x60	Test RPM gauge stop / Exhaust backpressur valve test	60 00
0x61	Variable intake test	61 00
0x63	Test RPM gauge	63 00
0x64	Test Boost gauge	64 00
0x65	S/W Throt S.W	65 00
0x67	Fan 3 (engine bay) off	67 00
0x6B	Test RPM gauge start	6B 00
0x6D	?	6D 00
0x6F	Fan 3 (engine bay) on	6F 00
0x79	Increments fuel trim setting and returns the new value	0x79, [new value]
0x7A	Decrements fuel trim setting and returns the new value	0x7A, [new value]
0x7B	Increments fuel trim setting and returns the new value	0x7B, [new value]
0x7C	Decrements fuel trim setting and returns the new value	0x7C, [new value]
0x7D	Request data packet 0x7D	0x7D, [data packet]
0x7E	? Used as part of (auto?) idle adjustment	7E 08
0x7F	? Used as part of (auto?) ignition adjustment	7F 05
0x80	Request data frame 80, followed by 28-byte data frame
0x81	? Used at end of idle/ignition/clearing (auto?) adjustments	0x81, 0x00
0x82	?	82 09 9E 1D 00 00 60 05 FF FF
0x89	Increments idle decay setting and returns the new value	0x89, [new value]
0x8A	Decrements idle decay setting and returns the new value	0x8A, [new value]
0x91	Increments idle speed setting and returns the new value	0x91, [new value]
0x92	Decrements idle speed setting and returns the new value	0x92, [new value]
0x93	Increments ignition advance offset and returns the new value	0x93, [new value]
0x94	Decrements ignition advance offset and returns the new value	0x94, [new value]
0x9E	Alternate first byte of init (different diag mode/security level?)	9E
0xC4	Swap to diagnostic mode 4 (only from mode 3)	C4 xx
0xCA	First byte of "normal" init	CA
0xCB	?	CB 00
0xCC	Clear fault codes	CC 00
0xCD	Debug? Read RAM?	CD 01
0xCE	Alternate first byte of init (different diag mode/security level?)	CE
0xCF	Alternate first byte of init (different diag mode/security level?)	CF
0xD0	ECU/Software ID	D0 99 00 03 03
0xD1	ECU/Software IDs 1x integer, 1x Ascii	D1 41 42 4E 4D 50 30 30 33 99 00 03 03
        e.g. integer: 99 00 03 03
        e.g. string/Ascii: 41 42 4E 4D 50 30 30 33 = ABNMP003
0xD2	Read security status	D2, followed by 02 01, 00 01, or 01 01
0xD3	Recode ECU	D3, followed by 02 01, 00 02, or 01 01 (reply needs checking)
0xDA	Test injector 1 (mems 1.9)	DA, 01?
0xDB	Test injector 2 (mems 1.9)	DB, 01?
0xDE	Alternate first byte of init (different diag mode/security level?)	DE
0xE0	Alternate first byte of init (different diag mode/security level?)	E0
0xE5	Alternate first byte of init (different diag mode/security level?)	E5
0xE7	?	E7 02
0xE8	?	E8 05 26 01 00 01
0xE9	Clear faults 2nd? needs ignition cycle?
0xEA
0xEB
0xEC
0xED	?	ED 00
0xEE	?	EE 00
0xEF	Actuate fuel injectors? (MPi?)	EF 03
0xF0	Check current diagnostic mode	F0 14 - mode 3 (default), 1E - mode 4, 50 - mode 5 or 6 (you should know which)
0xF1
0xF2	Swap to diagnostic mode 6 (only from mode 4)	F2 xx
0xF3	Swap to diagnostic mode 4 (from mode 5,6)	F3 xx
0xF4	Swap to diagnostic mode 5 (from mode 3)	F4 xx
0xF5	Swap to diagnostic mode 3 (from mode 4,5,6)	F5 xx
0xF6	Disconnect/Reset diagnostic session
0xF7	Actuate fuel injector (SPi?)	F7 03
0xF8	Fire ignition coil	F8 02
0xF9	Adjust main map? 2 bytes input also?	0x00 on success
0xFA	Clear all adaptations	0xFA, 0x00
0xFB	Request current IAC position	0xFB [IAC position XX]
0xFC	?	FC 00
0xFD	Open IAC by one step and report current position	FD, [IAC position]
0xFE	Close IAC by one step and report current position	FE, [IAC position]
0xFF	Request current IAC position?
*/

// MEMSHeartbeat command code for a communication heartbeat
var MEMSHeartbeat = []byte{0xf4}

// MEMSDataFrame request complete dataframes using 0x7d and 0x80 coomands
var MEMSDataFrame = []byte{0x80, 0x7d}

// MEMSReqData80 command code for requesting data frame 0x80
var MEMSReqData80 = []byte{0x80}

// MEMSReqData7D command code for requesting data frame 0x7D
var MEMSReqData7D = []byte{0x7d}

// MEMSInitCommandA command code to start initialisation sequence
var MEMSInitCommandA = []byte{0xca}

// MEMSInitCommandB command code forms second command as part of the initialisation sequence
var MEMSInitCommandB = []byte{0x75}

// MEMSInitECUID command code for retrieving the ECU ID as the final step in initialisation
var MEMSInitECUID = []byte{0xd0}

// MEMSGetECUID command code to retrieve the full ECU ID
// example: D1 41 42 4E 4D 50 30 30 33 99 00 03 03
// 41 42 4E 4D 50 30 30 33 = ABNMP003
// ID = 99 00 03 03
var MEMSGetECUSerial = []byte{0xd1}

// MEMSClearFaults command code to clear fault codes
var MEMSClearFaults = []byte{0xCC}

// MEMSResetAdj command code that instructs the ECU to clear all adjustments
var MEMSResetAdj = []byte{0x0F}

// MEMSResetECU command code that instructs the ECU to clear all computed and learnt settings
var MEMSResetECU = []byte{0xFA}

// MEMS Adjustment Settings
//
// | Setting                 | Decrement | Increment |
// | ----------------------- | --------- | --------- |
// | Fuel trim (Short Term?) |     7A    |     79    |
// | Fuel trim (Long Term?)  |     7C    |     7B    |
// | Idle decay              |     8A    |     89    |
// | Idle speed              |     92    |     91    |
// | Ignition advance offset |     94    |     93    |
// | IAC Position            |     FD    |     FE    |

//
// Fuel Trim Adjustments
//

// MEMSSTFTDecrement command
var MEMSSTFTDecrement = []byte{0x7a}

// MEMSSTFTIncrement command
var MEMSSTFTIncrement = []byte{0x79}

// MEMSLTFTDecrement command
var MEMSLTFTDecrement = []byte{0x7c}

// MEMSLTFTIncrement command
var MEMSLTFTIncrement = []byte{0x7b}
var MEMSFuelTrimMin = 0x00
var MEMSFuelTrimMax = 0xfe
var MEMSFuelTrimDefault = 0x8a

// MEMSGetIACPosition command code to retrieve the Idle Air Control Position
var MEMSGetIACPosition = []byte{0xfb}
var MEMSIACPositionDefault = 0x80
var MEMSIACIncrement = []byte{0xfd}
var MEMSIACDecrement = []byte{0xfe}

//
// Idle Decay Adjustments
//

// MEMSIdleDecayDecrement commad
var MEMSIdleDecayDecrement = []byte{0x8a}

// MEMSIdleDecayIncrement command
var MEMSIdleDecayIncrement = []byte{0x89}
var MEMSIdleDecayMin = 0x0a
var MEMSIdleDecayMax = 0x3c
var MEMSIdleDecayDefault = 0x23

//
// Idle Speed Adjustments
//

// MEMSIdleSpeedDecrement command
var MEMSIdleSpeedDecrement = []byte{0x92}

// MEMSIdleSpeedIncrement command
var MEMSIdleSpeedIncrement = []byte{0x91}
var MEMSIdleSpeedMin = 0x78
var MEMSIdleSpeedMax = 0x88
var MEMSIdleSpeedDefault = 0x80

//
// Ignition Advance Adjustments
//

// MEMSIgnitionAdvanceOffsetDecrement command
var MEMSIgnitionAdvanceOffsetDecrement = []byte{0x94}

// MEMSIgnitionAdvanceOffsetIncrement command
var MEMSIgnitionAdvanceOffsetIncrement = []byte{0x93}
var MEMSIgnitionAdvanceOffsetMin = 0x74
var MEMSIgnitionAdvanceOffsetMax = 0x8c
var MEMSIgnitionAdvanceOffsetDefault = 0x80

//
// Actuators
//

var MEMSFuelPumpOn = []byte{0x11}
var MEMSFuelPumpOff = []byte{0x01}
var MEMSPTCRelayOn = []byte{0x12}
var MEMSPTCRelayOff = []byte{0x02}
var MEMSACRelayOn = []byte{0x13}
var MEMSACRelayOff = []byte{0x03}
var MEMSPurgeValveOn = []byte{0x18}
var MEMSPurgeValveOff = []byte{0x08}
var MEMSO2HeaterOn = []byte{0x19}
var MEMSO2HeaterOff = []byte{0x09}
var MEMSBoostValveOn = []byte{0x1B}
var MEMSBoostValveOff = []byte{0x0B}
var MEMSFan1On = []byte{0x1D}
var MEMSFan1Off = []byte{0x0D}
var MEMSFan2On = []byte{0x1E}
var MEMSFan2Off = []byte{0x0E}
var MEMSTestInjectors = []byte{0xF7}
var MEMSFireCoil = []byte{0xF8}
