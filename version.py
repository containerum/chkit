from requests import get
from config_json_handler import check_last_update
from datetime import datetime, timedelta
import json
from bcolors import BColors

GITHUB_ADDR = "https://api.github.com/repos/containerum/chkit/releases/latest"
MAX_DIFF_TIME = 2


class Version:
    def compare_current_version(self, current_version):
        last_checked = check_last_update()
        if last_checked:
            now = datetime.now()
            diff = now - datetime.strptime(last_checked, "%Y-%m-%d %H:%M:%S.%f")
            human_diff = (diff.seconds / 60)
            if human_diff > 2:
                result = get(url=GITHUB_ADDR)
                if result.status_code == 200:
                    latest_version = result.json().get("tag_name")
                    latest_version = latest_version.split(".")
                    self.latest_version = [int(i) for i in latest_version]
                    current_version = current_version.split(".")
                    self.current_version = [int(i) for i in current_version]
                    if latest_version[0] > current_version[0] or latest_version[1] > current_version[1]:
                        self.print_error()
                        return
                    if latest_version[0] == current_version[0] and latest_version[1] == current_version[1] and latest_version[2] > current_version[2]:
                        self.print_warning()
        return True

    def print_error(self):
        print('{}{}{}{}{}'.format(
            BColors.FAIL,
            "Your Version ",
            self.current_version,
            "is too old! Please update itself!",
            BColors.ENDC
        ))

    def print_warning(self):
        print('{}{}{}'.format(
            BColors.WARNING,
            "We recommend you update your ChKit",
            BColors.ENDC
        ))