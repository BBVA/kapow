
from setuptools import setup

setup(
    name='kapow',
    version='0.0.1',
    py_modules=['kapow'],
    include_package_data=True,
    install_requires=[
        'aiofiles==0.4.0',
        'aiohttp==3.5.4',
        'pyparsing==2.3.1',
        'Click==7.0'
    ],
    zip_safe=True,
    entry_points={
        'console_scripts': [
            'kapow = kapow:main',
        ]
    }
)
