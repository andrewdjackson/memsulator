cd # memsreader
Rover MEMS 1.6 Code Reader in Python


# Protocol

The ECU's diagnostic port uses a serial protocol, with 8 data bits, no parity, and 1 stop bit, running at 9600 bps.

When a diagnotic tool (or a PC running diagnostic software) is first connected to the ECU, it must send an initialization sequence to enable further communication. I haven't looked inside the SPi MEMS firmware and I've only run experiments with a single car, so I can't provide any further insight on the meaning of this sequence. I'd certainly like to hear from anyone who has more information. Note that the first byte returned by the ECU appears to always be an echo of the command byte.

## Table 1: Initialization sequence

| Sent by tool | Response from ECU (MNE101170) | Response from ECU (MNE101070) |
| --- | --- | --- |
| CA | CA | CA |
| 75 | 75 | 75 |
| D0 | D0 99 00 03 03 | D0 99 00 02 03 |

Once the initialization sequence is complete, any of several single-byte commands may be sent to the ECU. They are described in Table 3.

Some of the commands appear to increment or decrement settings that have a predefined range. The settings that I've found are described in Table 2.

## Table 2: Software counters/settings used in the firmware on MNE101170\. All values in hex.

| Setting | Default value | Minimum | Maximum | Decremented by command(s) | Incremented by command(s) |
| --- | --- | --- | --- | --- | --- |
| Fuel trim | 8A | 00 | FE | 7A and 7C | 79 and 7B |
| Idle decay | 23 | 0A | 3C | 8A | 89 |
| Idle speed | 80 | 78 | 88 | 92 | 91 |
| Ignition advance offset | 80 | 74 | 8C | 94 | 93 |

Although some of the actuators have pairs of on/off commands to drive them, I've found that the system fitted to the Mini SPi will actually turn off the AC relay, PTC relay, and fuel pump relay automatically after a short time (i.e. without requiring the 'off' command). The 'off' command is acknowledged by the ECU, but apparently has no effect.

## Table 3: Protocol command bytes (values in hex)

| Command byte (sent to ECU) | Description | Response from ECU (MNE101170) |
| --- | --- | --- |
| 01 | Open fuel pump relay (stop fuel pump) | 01 00 |
| 02 | Open PTC relay (inlet manifold heater) | 02 00 |
| 03 | Open air conditioning relay | 03 00 |
| 08 | Close purge valve? | 08 00 |
| 09 | Open O2 heater relay? | 09 00 |
| 11 | Close fuel pump relay (run fuel pump) | 11 00 |
| 12 | Close PTC relay (inlet manifold heater) | 12 00 |
| 13 | Close air conditioning relay | 13 00 |
| 18 | Open purge valve? | 18 00 |
| 19 | Close O2 heater relay? | 19 00 |
| 1D | Close Fan 1 relay? Together with command 1E, I found through protocol analysis that this has something to do with testing an electric fan, but my ECUs didn't respond beyond the echo byte. | 1D |
| 1E | Close Fan 2 relay? | 1E |
| 65 | ? | 65 00 |
| 6D | ? | 6D 00 |
| 79 | Increments fuel trim setting and returns the current value. | 79 xx, where the second byte is the counter value |
| 7A | Decrements fuel trim setting and returns the current value. | 7A xx, where the second byte is the counter value |
| 7B | Increments fuel trim setting and returns the current value. | 7B xx, where the second byte is the counter value |
| 7C | Decrements fuel trim setting and returns the current value. | 7C xx, where the second byte is the counter value |
| 7D | Request data frame | 7D, followed by 32-byte data frame; see Table 4 |
| 7E | ? | 7E 08 |
| 7F | ? | 7F 05 |
| 80 | Request data frame | 80, followed by 28-byte data frame; see Table 5 |
| 82 | ? | 82 09 9E 1D 00 00 60 05 FF FF |
| 89 | Increments idle decay setting and returns the current value | 89 xx, where the second byte is the counter value |
| 8A | Decrements idle decay setting and returns the current value | 8A xx, where the second byte is the counter value |
| 91 | Increments idle speed setting and returns the current value | 91 xx, where the second byte is the counter value |
| 92 | Decrements idle speed setting and returns the current value | 92 xx, where the second byte is the counter value |
| 93 | Increments ignition advance offset and returns the current value | 93 xx, where the second byte is the counter value |
| 94 | Decrements ignition advance offset and returns the current value | 94 xx, where the second byte is the counter value |
| CB | ? | CB 00 |
| CC | Clear fault codes | CC 00 |
| CD | ? | CD 01 |
| D0 | ? | D0 99 00 03 03 |
| D1 | ? | D1 41 42 4E 4D 50 30 30 33 99 00 03 03 41 42 4E 4D 50 30 30 33 99 00 03 03 41 42 4E 4D 50 30 30 33 99 00 03 03 |
| D2 | ? | D2, followed by 02 01, 00 01, or 01 01 |
| D3 | ? | D3, followed by 02 01, 00 02, or 01 01 |
| E7 | ? | E7 02 |
| E8 | ? | E8 05 26 01 00 01 |
| ED | ? | ED 00 |
| EE | ? | EE 00 |
| EF | Actuate fuel injectors? (MPi?) | EF 03 |
| F0 | ? | F0 50 |
| F3 | ? | F3 00 |
| F4 | NOP / heartbeat? Sent continuously by handheld diagnostic tools to verify serial link. | F4 00 |
| F5 | ? | F5 00 |
| F6 | Unknown. Disconnect or sleep command? During a bench test, the Mini SPi ECU changed its behavior after being commanded with F6; it would only echo single bytes that it received, but nothing else. For example, sending 80 would cause it to echo 80 but no data frame followed. Re-sending the init sequence of CA, 75, and D0 restored normal behavior. | F6 00 |
| F7 | Actuate fuel injector (SPi?) | F7 03 |
| F8 | Fire ignition coil | F8 02 |
| FB | Request current IAC position | FB xx, where the second byte represents the IAC position |
| FC | ? | FC 00 |
| FD | Open IAC by one step and report current position | FD xx, where the second byte represents the IAC position |
| FE | Close IAC by one step and report current position | FE xx, where the second byte represents the IAC position |
| FF | Request current IAC position? |

The tables below describe the known fields in the data frames that are sent by the ECU in response to commands 0x7D and 0x80, respectively. Multi-byte fields are sent in big-endian format.

## Table 4: Data frame (Command 0x7D)

| Byte position | Field |
| --- | --- |
| 0x00 | Size of data frame, including this byte |
| 0x01 | Unknown |
| 0x02 | Throttle angle? |
| 0x03 | Unknown |
| 0x04 | Sometimes documented to be air/fuel ratio, but often observed to never change from 0xFF |
| 0x05 | Unknown |
| 0x06 | Lambda sensor voltage, 0.5mv per LSB |
| 0x07 | Lambda sensor frequency? |
| 0x08 | Lambda sensor duty cycle? |
| 0x09 | Lambda sensor status? 0x01 for good, any other value for no good |
| 0x0A | Loop indicator, 0 for open loop and nonzero for closed loop |
| 0x0B | Long term trim? |
| 0x0C | Short term trim, 1% per LSB |
| 0x0D | Carbon canister purge valve duty cycle? |
| 0x0E | Unknown |
| 0x0F | Idle base position |
| 0x10 | Unknown |
| 0x11 | Unknown |
| 0x12 | Unknown |
| 0x13 | Unknown |
| 0x14 | Idle error? |
| 0x15 | Unknown |
| 0x16 | Unknown |
| 0x17 | Unknown |
| 0x18 | Unknown |
| 0x19 | Unknown |
| 0x1A | Unknown |
| 0x1B | Unknown |
| 0x1C | Unknown |
| 0x1D | Unknown |
| 0x1E | Unknown |
| 0x1F | Unknown |

## Table 5: Data frame (Command 0x80)

| Byte position | Field |
| --- | --- |
| 0x00 | Size of data frame, including this byte. This should be 0x1C (28 bytes) for the frame described here. |
| 0x01-2 | Engine speed in RPM (16 bits) |
| 0x03 | Coolant temperature in degrees C with +55 offset and 8-bit wrap |
| 0x04 | Computed ambient temperature in degrees C with +55 offset and 8-bit wrap |
| 0x05 | Intake air temperature in degrees C with +55 offset and 8-bit wrap |
| 0x06 | Fuel temperature in degrees C with +55 offset and 8-bit wrap. This is not supported on the Mini SPi, and always appears as 0xFF. |
| 0x07 | MAP sensor value in kilopascals |
| 0x08 | Battery voltage, 0.1V per LSB (e.g. 0x7B == 12.3V) |
| 0x09 | Throttle pot voltage, 0.02V per LSB. WOT should probably be close to 0xFA or 5.0V. |
| 0x0A | Idle switch. Bit 4 will be set if the throttle is closed, and it will be clear otherwise. |
| 0x0B | Unknown. Probably a bitfield. Observed as 0x24 with engine off, and 0x20 with engine running. A single sample during a fifteen minute test drive showed a value of 0x30. |
| 0x0C | Park/neutral switch. Zero is closed, nonzero is open. |
| 0x0D | Fault codes. On the Mini SPi, only two bits in this location are checked: |
 Bit 0: Coolant temp sensor fault (Code 1) |
 Bit 1: Inlet air temp sensor fault (Code 2) |
| 0x0E | Fault codes. On the Mini SPi, only two bits in this location are checked: |
 Bit 1: Fuel pump circuit fault (Code 10) |
 Bit 7: Throttle pot circuit fault (Code 16) |
| 0x0F | Unknown |
| 0x10 | Unknown |
| 0x11 | Unknown |
| 0x12 | Idle air control motor position. On the Mini SPi's A-series engine, 0 is closed, and 180 is wide open. |
| 0x13-14 | Idle speed deviation (16 bits) |
| 0x15 | Unknown |
| 0x16 | Ignition advance, 0.5 degrees per LSB with range of -24 deg (0x00) to 103.5 deg (0xFF) |
| 0x17-18 | Coil time, 0.002 milliseconds per LSB (16 bits) |
| 0x19 | Unknown |
| 0x1A | Unknown |
| 0x1B | Unknown |

## MNE ECU Codes
Single Point Injection except Cooper upto (V)059586
serial number: MNE10097

Single Point Injection except Cooper (V)BD059587 upto (V)BD103112
serial number: MNE10089 & MNE101040

Single Point Injection except Cooper (V)BD103113 upto (V)BD134454 with anti thieft system
serial number: MNE101150

Single Point Injection Cooper upto (V)BD060487
serial number: MNE10027

Single Point Injection Cooper (V)BD060488 upto (V)BD069492
serial number: MNE10092, alternative MNE101070 could be used

Single Point Injection Cooper (V)BD069493 upto (V)BD103112
serial number: MNE101070

Single Point Injection Cooper (V)BD103113 upto (V)BD134454
serial number: MNE101170

Single Point Injection with automatic gearbox upto (V)059844
serial number: MNE10026

Single Point Injection with automatic gearbox (V)BD059845 upto (V)BD134454
serial number: MNE10090 & MNE101060

Multi Point Injection (V)BD134455 upto (V)BD137649
serial number: MKC104290

Multi Point Injection (V)BD137650 on
serial number: MKC104291 & MKC104292
