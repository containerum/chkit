from cx_Freeze import Executable,setup

setup(
    packages=['queue'],
    name="client",
    description="Containerum client",
    version="1.0",
    executables=[Executable("client.py")]
)

