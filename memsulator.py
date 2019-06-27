import logging
import serial
import time
from mems.MemsReader import MemsReader

logging.basicConfig(level=logging.DEBUG, format='%(asctime)s - %(name)s - %(levelname)s - %(message)s')

# create an emulator for MEMS, this will simulate the MEMS responses
emulator = MemsReader('ttycodereader', True)
emulator.connect()

try:
    while True:
        time.sleep(0.1)
except KeyboardInterrupt:
    pass

emulator.disconnect()
exit(0)