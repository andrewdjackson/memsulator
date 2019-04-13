import logging
from binascii import hexlify
import mems.protocol.Rosco

class MemsClient(object):
    def __init__(self, interface):
        self._logger = logging.getLogger('mems.client')

        self._rosco = mems.protocol.Rosco.Rosco()
        self._interface = interface
        self._initialized = False

    @property
    def responses(self):
        return self._interface.recieve_queue

    @property
    def connected(self):
        return self._interface.connected

    @property
    def has_data_to_send(self):
        return self._interface.has_data_to_send()

    def connect(self):
        if self._interface.connected == False:
            self._interface.connect()
            self._logger.info('connected')

    def disconnect(self):
        if self._interface.connected == True:
            self._interface.disconnect()
            self._logger.info('disconnected')

    def hex_to_string(self, code):
        request_code = hexlify(code)
        return str(request_code.decode("utf-8")).upper()

    def send_command(self, command, data=bytearray()):
        # get the command code for the command
        command_code = self._rosco.get_command_code(command)

        mems_command = bytearray()
        mems_command.extend(command_code)

        # add the data payload if required
        if data is not None:
            if len(data) > 0:
                mems_command.extend(data)

        return self._send_data(mems_command)

    def _send_data(self, mems_command):
        self._logger.debug('sending command %s to MEMS',  self.hex_to_string(mems_command))

        expected_result = self._rosco.get_expected_result(mems_command)
        expected_bytes = len(expected_result)

        # send the command
        self._interface.send(mems_command, expected_bytes)

        self._logger.debug('waiting for response from MEMS..')
        # read the response
        # or wait until the read timeout has expired
        while self._interface.waiting_for_mems_response:
            pass

        response = self._interface.read()
        self._logger.info('%s response from MEMS', self.hex_to_string(response))

        return response
