# ConcealCli - conceal identifiable data from Aerospike log file

The `concealcli` tool can be used to predictably replace potentially identifiable data from Aerospike log files.

The data being replaced includes node IDs, IP addresses, ports, namespace names, and many others.

The tool uses an encryption key together with a deterministic algorithm, so that multiple logs from multiple machines can be processed in exactly the same way, with exactly the same replaces - as long as the key file used is the same.

The algorithm cannot be reversed without the key, and may not be reversable even with the key. The tool therefore generates a simple replacement map that can be used to match the replaced values to originals.

## TODO - THIS APPLICATIONS IS A WORK IN PROGRESS

Only handling NodeIDs from the ticker line and IP addresses (plus IP:PORT pairs) from logs so far. More will be added in due course.
