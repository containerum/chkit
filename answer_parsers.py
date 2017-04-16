from datetime import datetime
from dateutil import parser
from prettytable import PrettyTable
import json

class GetParser:
    def __init__(self, row_answer):
        self.items = row_answer.get("results")[0].get("data").get("items")


    def show_human_readable_result(self):
        self.table = PrettyTable(["NAME", "READY", "STATUS", "RESTARTS", "AGE", "IP"])
        self.table.align = "l"
        for i in self.items:
            restarts = i.get("status").get("containerStatuses")
            restarts_sum = 0
            if restarts:
                for r in restarts:
                    restarts_sum += r.get("restartCount")
                else:
                    restarts_sum = None
            ip = i.get("status").get("podIP")
            status = i.get("status").get("phase")
            name = i.get("metadata").get("name")
            ready = "-/-"
            time = get_datetime_diff(i.get("metadata").get("creationTimestamp"))
            self.table.add_row([name, ready, status, restarts_sum, time, ip])
        print(self.table)


def get_datetime_diff(timestamp):
    created_date = parser.parse(timestamp)
    current_date = datetime.now()
    diff = ((current_date.day - created_date.day, "d"), (current_date.hour - created_date.hour, "h"),
            (current_date.minute - created_date.minute, "m"), (current_date.second - created_date.second,"s"))
    diff = tuple(filter(lambda x: x[0] > 0, diff))[0]
    return str(diff[0]) + diff[1]