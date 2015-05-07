# trout

ping a host every $configurable seconds and perform a traceroute if it
is unreachable

## prerequisites

trout depends on the following programs being present on the system:

* fping
* traceroute

## installation

with a relatively new version of the go language tools that understand
custom repositories, you can just do

    go install hubs.net.uk/sw/trout

otherwise, clone this repository and do

    go build .

from within the checked-out version.
