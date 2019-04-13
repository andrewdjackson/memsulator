import logging
import serial
import time
import os
from mems.MemsEmulator import MemsEmulator

logging.basicConfig(level=logging.INFO, format='%(asctime)s - %(name)s - %(levelname)s - %(message)s')

port_settings = {
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

abspath = os.path.abspath(__file__)
dname = os.path.dirname(abspath)
print(dname)

# create an emulator for MEMS, this will simulate the MEMS responses
# when commands are sent to the MEMS Interface
emulator_connection =  {'port': 'ttycodereader'}
emulator_connection.update(port_settings)
emulator = MemsEmulator(emulator_connection)
emulator.connect()

try:
    while True:
        time.sleep(0.1)
except KeyboardInterrupt:
    pass

emulator.disconnect()
exit(0)