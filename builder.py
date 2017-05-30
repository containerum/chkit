from cx_Freeze import Executable,setup

setup(
    packages=['queue'],
    name="chkit",
    description="Containerum Hosting Client",
    version="1.2.2",
    executables=[Executable("chkit.py")]
)

