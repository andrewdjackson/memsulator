# ROSCO - Rover Communication Protocol
import os
from binascii import hexlify


class Rosco(object):
    def __init__(self):
        self._version   = 'MNE101170'
        self._connected = False

        self._versions = {
                'MNE101070' : b'\xd0\x99\x00\x02\x03',
                'MNE101170' : b'\xd0\x99\x00\x03\x03'
            }

        #
        # The tables below describe the known fields in the data frames
        # that are sent by the ECU in responses to commands 0x7D and 0x80, respectively.
        # Multi-byte fields are sent in big-endian format.
        #

        self._dataframes = [
                {'command' : b'\x7d', 'fields' : ['command','dataframe_size','','throttle_angle','', '','lambda_voltage','lambda_frequency','lambda_duty_cycle','lambda_status','loop_indicator','long_term_trim','short_term_trim','carbon_canister_purge_valve_duty_cycle','','idle_base_position','','','','','idle_error','','','','','','','','','','','','']},
                {'command' : b'\x80', 'fields' : ['command','dataframe_size','engine_speed_low_byte','engine_speed_high_byte','coolant_temperature','ambient_temperature','intake_air_temperature','fuel_temperature','map_sensor','battery_voltage','throttle_pot_voltage','idle_switch','','park_neutral_switch','coolant_temp_inlet_air_temp_sensor_fault','fuel_pump_throttle_pot_circuit_fault','','','','idle_air_contol_position','idle_speed_deviation_low_byte','idle_speed_deviation_high_byte','','ignition_advance','coil_time_low_byte','coil_time_high_byte','','','']}
            ]

        self._commands = [
                {'open_fuel_pump_relay' : b'\x01'},
                {'open_ptc_relay'       : b'\x02'},
                {'open_aircond_relay'   : b'\x03'},
                {'close_purge_valve'    : b'\x08'},
                {'open_heater_relay'    : b'\x09'},
                {'close_fuel_pump_relay': b'\x11'},
                {'close_ptc_relay'      : b'\x12'},
                {'close_aircond_relay'  : b'\x13'},
                {'open_purge_valve'     : b'\x18'},
                {'close_heater_relay'   : b'\x19'},
                {'close_fan_relay'      : b'\x1e'},
                {'request_data_frame_a' : b'\x7d'},
                {'inc_fuel_trim'        : b'\x79'},
                {'dec_fuel_trim'        : b'\x7a'},
                {'inc_fuel_trim_alt'    : b'\x7b'},
                {'dec_fuel_trim_alt'    : b'\x7c'},
                {'request_data_frame_b' : b'\x80'},
                {'inc_idle_decay'       : b'\x89'},
                {'dec_idle_decay'       : b'\x8a'},
                {'inc_idle_speed'       : b'\x91'},
                {'dec_idle_speed'       : b'\x92'},
                {'inc_ignition_advance' : b'\x93'},
                {'dec_ignition_advance' : b'\x94'},
                {'clear_fault_codes'    : b'\xcc'},
                {'heartbeat'            : b'\xf4'},
                {'actuate_fuel_injector': b'\xf7'},
                {'fire_ignition_coil'   : b'\xf8'},
                {'current_iac_position' : b'\xfb'},
                {'open_iac_by_one_step' : b'\xfd'},
                {'close_iac_by_one_step': b'\xfe'},
            ]

    def hex_to_string(self, code):
        request_code = hexlify(code)
        return str(request_code.decode("utf-8")).upper()

    def get_expected_result(self, code):
        response = bytearray()
        request_code = self.hex_to_string(code)

        abspath = os.path.abspath(__file__)
        dname = os.path.dirname(abspath)
        responsefile = dname + '/responses/' + request_code + '.hex'

        if os.path.isfile(responsefile):
            with open(responsefile, mode='rb') as f:
                data = f.read()
                response.extend(data)
        else:
            response.extend(code)
            response.extend(b'00')

        return response

    def get_command_code(self, command_name):
         for c in self._commands:
            if command_name in c:
                return c[command_name]

         return command_name

    def get_version(self, response_code):
        for key, value in self._versions.items():
            if response_code == value:
                return key

    def get_dataframe(self, command_name):
        code = self.get_command_code(command_name)

        for dataframe in self._dataframes:
            if code == dataframe['command']:
                return dataframe['fields']

    @property
    def initialization_sequence(self):
        return [
                   {'tx': b'\xca', 'response': b'\xca'},
                   {'tx': b'\x75', 'response': b'\x75'},
                   {'tx': b'\xd0', 'response': self._version}
               ]

