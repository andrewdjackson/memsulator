import logging
import serial
from binascii import hexlify
import mems.protocol.Rosco
from mems import MemsEmulator
from mems.interfaces.MemsSerialInterface import MemsSerialInterface

class MemsReader(object):
    def __init__(self, serial_port, is_emulator=False):
        self._logger = logging.getLogger('mems.reader')
        self._client = None

        mems_connection = {
            'port' : '/dev/null',
            'baudrate': 9600,
            'bytesize': serial.EIGHTBITS,
            'stopbits': serial.STOPBITS_ONE,
            'parity': serial.PARITY_NONE,
            'xonxoff': False,
            'rtscts': False,
            'dsrdtr': False,
            'write_timeout': 1,
            'timeout': 0
        }

        # create the serial interface to connect to MEMS ECU
        mems_connection['port'] = serial_port

        if is_emulator:
            self._client = mems.MemsEmulator(mems_connection)
        else:
            interface = MemsSerialInterface(mems_connection)
            self._client = mems.MemsClient(interface)

    @property
    def mems_interface(self):
        return self._client

    def connect(self):
        self._client.connect()

    def disconnect(self):
        self._client.disconnect()