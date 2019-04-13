import queue
from binascii import hexlify

import serial
import threading
import logging
import time

class SerialInterface():
    def __init__(self, connection):
        self._logger = logging.getLogger('serial.interface')

        self._send_queue = queue.Queue()
        self._receive_queue = queue.Queue()
        self._connected = False

        # to use serial_for_url the port and the settings need to be separated out
        # from the connection dictionary and then removed to prevent a duplicate
        # 'port' error
        self._port = connection['port']
        self._connection = connection
        del self._connection['port']

    @property
    def receive_queue(self):
        return self._receive_queue

    @property
    def connected(self):
        return self._connected

    def connect(self):
        # create a thread to listen for commands
        self._connected = True
        self._listener = threading.Thread(target=self._send_receive_loop)

        # create the virtual ports for emulate MEMS
        # self._create_virtual_ports()

        # open port for listening
        self.serial = serial.serial_for_url(self._port, **self._connection)
        self.serial.flush()
        # start listening for commands and responses
        self._listener.start()

        self._logger.debug('connected')

    def disconnect(self):
        while self.has_data_to_send():
            pass

        # stop listening
        self._connected = False

        return True

    def send(self, data, expected_result_length = 0):
        # add data to the send buffer
        self._send_queue.put({'data':data, 'expected_result_length':expected_result_length})

    def read(self):
        # read data from the read buffer
        if self.has_received_data():
            data = self._receive_queue.get()

            self._logger.debug('read %s from receive queue', self.hex_to_string(data))
            return data
        else:
            return None

    def _send_receive_loop(self):
        while self._connected:
            expected_result_length = 0

            if self.has_data_to_send():
                send_data = self._send_queue.get()
                data = send_data['data']
                expected_result_length = send_data['expected_result_length']

                self.serial.write(data)
                self.serial.flush()
                self._logger.debug('sent %s', self.hex_to_string(data))

            while self.serial.in_waiting:
                # determine the expected length to be received
                if (expected_result_length > 0):
                    data = self.serial.read(expected_result_length)
                else:
                    data = self.serial.read()

                self.serial.flush()
                self._logger.debug('read %s', self.hex_to_string(data))
                self._receive_queue.put(data)

                # this allows a response to be generated based on the received
                # data. This function can be overridden for custom processing.
                self.receive_loop_hook(data)

        # sleep for 50ms to yield time to the cpu without interfering with
        # the serial communications
        time.sleep(0.5)

    def hex_to_string(self, code):
        request_code = hexlify(code)
        return str(request_code.decode("utf-8")).upper()

    def receive_loop_hook(self, data):
        pass

    def has_data_to_send(self):
        return not self._send_queue.empty()

    def has_received_data(self):
        return not self._receive_queue.empty()
