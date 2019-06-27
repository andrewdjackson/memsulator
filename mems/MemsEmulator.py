import logging
import time
import subprocess
import threading
import mems.protocol.Rosco
from binascii import hexlify

from mems.interfaces.SerialInterface import SerialInterface

class MemsEmulator(SerialInterface):
    def __init__(self, connection):
        self._rosco = mems.protocol.Rosco.Rosco()

        self._serial_port_loopback = threading.Thread(target=self._create_virtual_ports)
        self._serial_port_loopback.start()
        time.sleep(1)

        super().__init__(connection)
        self._logger = logging.getLogger('mems.emulator.interface')

    def disconnect(self):
        yield super().disconnect()
        self._remove_virtal_ports()
        time.sleep(1)

    def receive_loop_hook(self, data):
        super().receive_loop_hook(data)
        response = self._get_response(data)
        self._logger.debug('mems response %s', self.hex_to_string(data))
        self.send(response)

    def hex_to_string(self, code):
        request_code = hexlify(code)
        return str(request_code.decode("utf-8")).upper()

    def _get_response(self, code):
        return self._rosco.get_response(code)

    def _create_virtual_ports(self):
        cmd = ['socat', '-d', '-d', 'pty,link=ttycodereader,raw,echo=0', 'pty,link=ttyecu,raw,echo=0']
        try:
            self.proc = subprocess.Popen(cmd)
        except Exception as ex:
            print (ex)

    def _remove_virtal_ports(self):
        # kill the virtual ports
        self.proc.kill()