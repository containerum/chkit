from datetime import datetime

from colorama import Fore
from requests import get
import os
from handlers.config_json_handler import check_last_update, save_checking_time

GITHUB_ADDR = "https://api.github.com/repos/containerum/chkit/releases/latest"
MAX_DIFF_TIME = 0.001


class Version:
    def __init__(self, current_version):
        self.current_version = current_version
        self.current_version_str = current_version

    def compare_current_version(self):
        last_checked = check_last_update()
        if last_checked:
            now = datetime.now()
            diff = now - datetime.strptime(last_checked, "%Y-%m-%d %H:%M:%S.%f")
            human_diff = (diff.seconds / 60)
            if human_diff > MAX_DIFF_TIME:
                result = get(url=GITHUB_ADDR)
                if result.status_code == 200:
                    latest_version = result.json().get("tag_name")
                    latest_version = latest_version.split(".")
                    self.latest_version = [int(i) for i in latest_version]
                    current_version = self.current_version.split(".")
                    self.current_version = [int(i) for i in current_version]
                    if latest_version[0] > current_version[0] or latest_version[1] > current_version[1]:
                        self.print_error()
                        return
                    if latest_version[0] == current_version[0] and latest_version[1] == current_version[1] and latest_version[2] > current_version[2]:
                        self.print_warning()
                    save_checking_time()
        return True

    def print_error(self):
        print('{}{}{}{}'.format(
            Fore.RED,
            "Your Version ",
            self.current_version_str,
            " is too old! Please update itself!",
        ))

    def print_warning(self):
        print('{}{}'.format(
            Fore.YELLOW,
            "We recommend you update your ChKit",
        ))

    def print_latest_version(self):
        print('{}{}'.format(
            Fore.GREEN,
            "You have the latest version of ChKit"

        ))

    def check_last_version(self):
        result = get(url=GITHUB_ADDR)
        if result.status_code == 200:
            latest_version = result.json().get("tag_name")
            latest_version = latest_version.split(".")
            self.latest_version = [int(i) for i in latest_version]
            current_version = self.current_version.split(".")
            self.current_version = [int(i) for i in current_version]
            if latest_version[2] >= current_version[2] and latest_version[1] >= current_version[1] and latest_version[0] >= current_version[0]:
                assets = result.json().get("assets")
                arch_urls = [i.get("browser_download_url") for i in assets]
                if os.name == 'posix':
                    for i in arch_urls:
                        if "linux" in i:
                            return i

            else:
                self.print_latest_version()
                return

