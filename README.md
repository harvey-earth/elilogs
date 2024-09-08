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

## Usage
### check
Check is used to check the connection to an Elasticsearch server/cluster.
Its primary use is to confirm that initial configuration is working.
There is no need to call check before any other commands, as they all test and error handle for connection failures.

It will have an exit code 0 if the check is successful and 2 if the check is not successful.

`elilogs check` will test the connection to the Elasticsearch API.  
`elilogs test -c` will download a list of indices to the users ~/.cache/elilogs.txt directory to speed up future calls.

### list
List is used for listing information (namely health) for clusters, nodes, and indexes.
It can also output information on pending tasks and snapshots.

List will have an exit code of 0 if the health of all returned objects is green.
An exit code of 1 indicates that at least one of the returned objects is not in a healthy state.
An exit code of 2 indicates an error running the check.

`elilogs list index [index,...]` will list indexes along with the status and health of each.
`elilogs list cluster` is equivalent to `-a` and will list information about the cluster, nodes, pending tasks, and snapshots.

### search
Search is used to search multiple indexes for a query string.
It returns an exit status of 0 for a successful search with results, 1 for a successful search with no results, and 2 if an error is encountered during the search.

`elilogs search "query"` will search all indexes for matches to the query string.
`elilogs search -i [index,...] query` will search the listed index for matches to the query string.
