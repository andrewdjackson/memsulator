import logging
from mems.interfaces.SerialInterface import SerialInterface

class MemsSerialInterface(SerialInterface):
    def __init__(self, connection):
        super().__init__(connection)
        self._logger = logging.getLogger('mems.serial.interface')
        self._waiting_for_mems_response = False

    @property
    def waiting_for_mems_response(self):
        return self._waiting_for_mems_response

    def connect(self):
        super().connect()
        self._logger.debug('connected')

    def disconnect(self):
        super().disconnect()
        self._logger.debug('disconnected')

    def send(self, data, expected_bytes):
        self._waiting_for_mems_response = True
        super().send(data, expected_bytes)
        self._logger.debug('sent %s to MEMS', str(data))

    def receive_loop_hook(self, data):
        self._waiting_for_mems_response = False
        super().receive_loop_hook(data)

    def read(self):
        data = bytearray()

        while self._receive_queue.qsize() > 0:
            data.extend( super().read() )

        self._logger.debug('read %s from MEMS', str(data))
        return data
