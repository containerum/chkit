from datetime import datetime
from dateutil import parser
from prettytable import PrettyTable
from keywords import EMPTY_NAMESPACE, NO_NAMESPACES




class TcpApiParser:
    def __init__(self, row_answer):
        if row_answer.get("results")[0].get("data").get("kind") == "PodList":
            self.items = row_answer.get("results")[0].get("data").get("items")
            self.show_human_readable_podlist()

        if row_answer.get("results")[0].get("data").get("kind") == "Pod":
            self.result = row_answer
            self.show_human_readable_pod()

    def show_human_readable_pod(self):
        metadata = self.result.get("results")[0].get("data").get("metadata")
        containers = self.result.get("results")[0].get("data").get("spec").get("containers")
        restartPolicy = self.result.get("results")[0].get("data").get("spec").get("restartPolicy")
        termination = self.result.get("results")[0].get("data").get("spec").get("terminationGracePeriodSeconds")
        system = self.result.get("results")[0].get("data").get("status")
        container_statuses = self.result.get("results")[0].get("data").get("status").get("containerStatuses")
        status = self.result.get("results")[0].get("data").get("status").get("conditions")

        print("Describe:")
        print("\t%-20s %s" % ("UserId:", self.result.get("UserId")))
        print("\t%-20s %s" % ("Channel:", self.result.get("channel")))
        print("\t%-20s %s" % ("CommandId:", self.result.get("id")))
        print("Pod:")
        print("\t%-20s %s" % ("CreationTime:", parser.parse(metadata.get("creationTimestamp"))))
        print("\tLabel:")
        print("\t\t%-20s %s" % ("App:", metadata.get("labels").get("app")))
        print("\t\t%-20s %s" % ("PodTemplateHash:", metadata.get("labels").get("pod-template-hash")))
        print("\t\t%-20s %s" % ("Role:", metadata.get("labels").get("role")))
        print("\t\t%-20s %s" % ("CommandId:", self.result.get("id")))
        print("Containers:")
        for c in containers:
            print("\t%s" % c.get("name"))
            if c.get("command"):
                print("\t\t%-20s %s" % ("Command:", "".join(c.get("command"))))
            print("\t\tPorts:")
            #print("\t\t\t%-20s %-10s %-10s" % ("Name", "Protocol", "ContPort"))
            ports = PrettyTable(["Name", "Protocol", "ContPort"])
            for p in c.get("ports"):
                #print("\t\t\t%-20s %-10s %-10s" % (p.get("name"), p.get("protocol"), p.get("containerPort")))
                ports.add_row([p.get("name"), p.get("protocol"), p.get("containerPort")])
            print(ports)
            if c.get("env"):
                env = PrettyTable(["Name","Value"])
                print("\t\tEnvironment:")
                for e in c.get("env"):
                    #print("\t\t\t%-20s %-10s %-10s" % (p.get("name"), p.get("protocol"), p.get("containerPort")))
                    env.add_row([e.get("name"), e.get("value")])
                print(env)
            print("\t\tResourceLimit:")
            print("\t\t\t%-10s %s" % ("CPU:", c.get("resources").get("limits").get("cpu")))
            print("\t\t\t%-10s %s" % ("Memory:", c.get("resources").get("limits").get("memory")))
            print("\t\t%-20s %s" % ("Image:", c.get("image")))
            print("\t\t%-20s %s" % ("ImagePullPolicy:", c.get("imagePullPolicy")))
            print("System:")
            print("\t%-30s %s" % ("PodIP:", system.get("podIP")))
            print("\t%-30s %s" % ("Phase:", system.get("phase")))
            print("\t%-30s %s" % ("StartTime:", parser.parse(system.get("startTime"))))
            print("\t%-30s %s" % ("TerminationGracePeriodSeconds:", termination))
            print("\t%-30s %s" % ("RestartPolicy:", restartPolicy))
            print("ContainerStatuses:")
            containerStatuses = PrettyTable(["Name","Ready","Restart Count"])
            for cs in container_statuses:
                containerStatuses.add_row([cs.get("name"), cs.get("ready"), cs.get("restartCount")])
            print(containerStatuses)
            print("Status:")
            StatusTable = PrettyTable(["Type:", "LastTransitionTime:", "Status:"])
            for s in status:
                StatusTable.add_row([s.get("type"), parser.parse(s.get("lastTransitionTime")), s.get("status")])
                # print("\t%-30s %s" % ("Type:", s.get("type")))
                # print("\t%-30s %s" % ("LastTransitionTime:", parser.parse(s.get("lastTransitionTime"))))
                # print("\t%-30s %s" % ("Status:", s.get("status")))
            print(StatusTable)

    def show_human_readable_podlist(self):
        if self.items:
            table = PrettyTable(["NAME", "READY", "STATUS", "RESTARTS", "AGE", "IP"])
            table.align = "l"
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
                table.add_row([name, ready, status, restarts_sum, time, ip])
            print(table)
        else:
            print(EMPTY_NAMESPACE)


class WebClientApiParser:
    def __init__(self, row_answer):
        self.items = row_answer
        self.result = row_answer

    def show_human_readable_namespace_list(self):
        if self.items:
            table = PrettyTable(["ID",  "IS ACTIVE",  "AGE", "CPU", "CPU LIMIT", "MEMORY", "MEMORY LIMIT"])
            table.align = "l"
            for i in self.items:
                status = i.get("active")
                name = i.get("id")
                cpu = i.get("cpu")
                cpu_limit = i.get("cpu_limit")
                memory = i.get("memory")
                memory_limit = i.get("memory_limit")
                time = get_datetime_diff(i.get("created"))
                table.add_row([name,  status,  time, cpu, cpu_limit, memory, memory_limit])
            print(table)
        else:
            print(NO_NAMESPACES)

    def show_human_readable_deployment_list(self):
        if self.items:
            table = PrettyTable(["NAME",  "PODS ACTIVE",  "CPU",  "RAM", "AGE"])
            table.align = "l"
            for i in self.items:
                name = i.get("name")
                cpu = i.get("cpu")
                pods_active = i.get("pods_active")
                memory = i.get("ram")
                time = get_datetime_diff(i.get("created_at"))
                table.add_row([name,  pods_active,  cpu,  memory, time])
            print(table)
        else:
            print(NO_NAMESPACES)

    def show_human_readable_deployment(self, namespace):
        status = self.result.get("status")
        conditions = self.result.get("conditions")
        containers = self.result.get("containers")
        print(self.result)
        if self.result:
            print("%-20s %s" % ("Name:", self.result.get("name")))
            print("%-20s %s" % ("Namespace:", namespace))
            print("%-20s %s" % ("CreationTimeStamp:", parser.parse(self.result.get("created_at"))))
            print("Labels:")
            for key,value in self.result.get("labels").items():
                print("\t%s=%s" % (key, value))
            status_tuple = ("Status:", status.get("updated"), "updated", status.get("total"), "total",
                            status.get("available"), "available", status.get("unavailable"), "unavailable")
            print("%-20s %s %s | %s %s | %s %s | %s %s" % status_tuple)
            print("Conditions:")
            conditions_table = PrettyTable(["TYPE", "STATUS", "REASON"])
            for c in conditions:
                conditions_table.add_row([c.get("type"), c.get("status"), c.get("reason")])
            print(conditions_table)
            print("Containers:")
            for c in containers:
                print("\t%s" % c.get("name"))


def get_datetime_diff(timestamp):
    created_date = parser.parse(timestamp)
    current_date = datetime.now()
    diff = ((current_date.year - created_date.year, "Y"),(current_date.month - created_date.month, "M"),
            (current_date.day - created_date.day, "d"), (current_date.hour - created_date.hour, "h"),
            (current_date.minute - created_date.minute, "m"), (current_date.second - created_date.second,"s"))
    diff = tuple(filter(lambda x: x[0] > 0, diff))[0]
    return str(diff[0]) + diff[1]