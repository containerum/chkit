import os


def get_file_config_path():
    if os.name == 'nt':
        return os.path.join(os.getenv("HomePath"), "containerum", "CONFIG.json")
    else:
        return os.path.join(os.getenv("HOME"), ".containerum", "CONFIG.json")


def get_templates_path():
    if os.name == 'nt':
        return os.path.join(os.getenv("HomePath"), "containerum", "src", "json_templates")
    else:
        return os.path.join(os.getenv("HOME"), ".containerum", "src", "json_templates")


def create_folders():
    if os.name == 'nt':
        os.system("mkdir  %s" % os.path.join(os.getenv("HomePath"), "containerum", "src", "json_templates"))
    else:
        os.system("chmod 777 -R %s" % os.path.join(os.getenv("HOME"), ".containerum"))
        os.system("mkdir -p %s" % os.path.join(os.getenv("HOME"), ".containerum", "src", "json_templates"))