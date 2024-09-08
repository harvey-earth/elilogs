# elilogs

This Golang CLI application allows easy interaction with Elasticsearch.
It can be used in health checks for clusters, nodes, and indexes in a monitoring and alerting solution (such as Nagios or Sensu).
It can also be used for multi-index queries (which is not possible with Kibana).

## Installation
### From source
1. Clone the repository
2. Run `make`
3. Configure settings in __config.yml__

### From binary
1. Download release at https://github.com/harvey-earth/elilogs/releases.
2. Copy file from __default-config.yml__ to configuration directory and rename to config.yml.

## Configuration
Configuration can be achieved by settings in the config.yml file.
Copy the default-config.yml file to the same file as the binary or at __/etc/elilogs.config.yml__.
Elasticsearch connections will be attempted in the following order.

### Elasticsearch Cloud
Ensure the `cloud_id` and `cloud_api_key` variables are properly set.

### HTTPS certificate
Set `ca_cert_path` to the local path where the certificate is located.
Additionally the `username` and `password` values must be set to a user with proper access.

### HTTPS Fingerprint
Set `certificate_fingerprint` with the SHA256 fingerprint value of the CA certificate.
Additionally the `username` and `password` values must be set to a user with proper access.

### Basic authentication
If only the `username` and `password` values are set, the connection is attempted with Basic authentication.

## Usage
### check
Check is used to check the connection to an Elasticsearch server/cluster.
Its primary use is to confirm that initial configuration is working.
There is no need to call check before any other commands, as they all test and error handle for connection failures.

It will have an exit code 0 if the check is successful and 1 if the check is not successful.

`elilogs check` will test the connection to the Elasticsearch API.  
`elilogs test -c` will download a list of indices to the users ~/.cache/elilogs.txt directory to speed up future calls.

### list
List is used for listing information (namely health) for clusters, nodes, and indexes.
It can also output information on pending tasks and snapshots.

List will have an exit code of 0 if the health of all returned objects is green.
An exit code of 1 indicates an error running the check.
An exit code of 2 indicates that at least one of the returned objects is not in a healthy state.

`elilogs list index [index,...]` will list indexes along with the status and health of each.
`elilogs list cluster` is equivalent to `-a` and will list information about the cluster, nodes, pending tasks, and snapshots.

### search
Search is used to search multiple indexes for a query string.
The query string should be in Lucene query syntax.
It returns an exit status of 0 for a successful search with results, 1 if an error is encountered during the search, and 2 for a successful search with no results.

`elilogs search "query"` will search all indexes for matches to the query string.
`elilogs search -i [index,...] query` will search the listed index for matches to the query string.
