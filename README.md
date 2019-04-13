# Rover MEMS 1.6 Emulator

Emulates the Rover MEMS 1.6 ECU. The emulator will create 2 virtual serial ports when run using socat.
The ports are shown as output on start up. The second port listed is the ECU. 

The reponses are held in binary files, so the response meets the protocol requirements but doesn't reflect the simulation of a real engine.

## Protocol

The ECU's diagnostic port uses a serial protocol, with 8 data bits, no parity, and 1 stop bit, running at 9600 bps.

## Readmems

I have included an example readmems.cfg. This config is used by my fork of librosco/readmems available in my readmems repo.
This version of readmems writes the output of read-raw to a log file and a binary file that can be used by the memsviewer diagnostic analyser.
