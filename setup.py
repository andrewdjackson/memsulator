from setuptools import setup

setup(name='memsulator',
      version='1.0.2',
      description='Rover MEMS 1.6 Serial Interface and Emulator',
      url='http://github.com/andrewdjackson/memsulator',
      author='Andrew Jackson',
      author_email='andrew.d.jackson@gmail.com',
      license='MIT',
      packages=['mems','mems.interfaces','mems.protocol'],
      zip_safe=False)