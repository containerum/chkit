from cx_Freeze import Executable,setup


exe = Executable(
    "chkit.py"
)

setup(
    packages=['queue'],
    name="chkit",
    description="Containerum Hosting Client",
    version="1.4.0",
    executables=[exe]
)

