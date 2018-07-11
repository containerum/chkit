Create configmap.
Configmap is a file storage, which can be mounted into a container. The most common usage of configmap is keeping config files, read-only DB, and secrets. Basically, you can think about it like about very simple key-value storage.

There are several ways to construct configmap:
- --item-string flag, formatted as KEY:VALUE pairs. The VALUE can be token, short init file, etc.
- --item-file flag, KEY:FILE_PATH or FILE_PATH (filename will be used as KEY)
- interactive wizard
- --import-file flag. Fields in imported file must be not base64 encoded!

Use the --force flag to skip wizard