#!/bin/bash

file=${1:-TMP_CONSUL_SERVICES.json}

# Get GoShimmer API address
goshimer_api=`curl -s localhost:8500/v1/catalog/service/goshimmer-testnet-leader | jq -r '.[] | select(.ServiceID | endswith("leader-api")) | ([.ServiceAddress, .ServicePort] | join(":"))'`

# Get service ports with the Node ports being the first set and the access ports being the second set
node_ports=`curl localhost:8500/v1/catalog/service/iscp-evm-node-wasp | jq '[ .[] | .ServiceID |= (split("-") | .[:-1] | join("-")) | .ServiceTags[] |= select(. != "wasp") | {ServiceID, Services: {(.ServiceTags[0]): ([.ServiceAddress, .ServicePort] | join(":"))}} ] | group_by(.ServiceID) | [ foreach .[] as $item ({}; reduce $item[] as $elem ({Services: {}}; .Services? += $elem.Services); .) ]'`
echo $node_ports | jq '.[]' > $file
access_ports=`curl localhost:8500/v1/catalog/service/iscp-evm-access-wasp | jq '[ .[] | .ServiceID |= (split("-") | .[:-1] | join("-")) | .ServiceTags[] |= select(. != "wasp") | {ServiceID, Services: {(.ServiceTags[0]): ([.ServiceAddress, .ServicePort] | join(":"))}} ] | group_by(.ServiceID) | [ foreach .[] as $item ({}; reduce $item[] as $elem ({Services: {}}; .Services? += $elem.Services); .) ]'`
echo $access_ports | jq '.[]' >> $file
all_ports=`jq --slurp . $file`
rm -rf $file

api_addresses=`echo $all_ports | jq -rc '.[].Services.api'`
declare -A pubKeys

# Get all public keys
for n in ${api_addresses[@]}
do
    echo "get peering for $n"
    peering=`curl -s -u wasp:wasp $n/adm/peering/self`
    pubKeys[$n]=$peering
done

# trust all public keys that are not the current node
for n in ${api_addresses[@]}
do
    for k in ${!pubKeys[@]}
    do
        # Nodes must trust themselves as well as all others
        echo "trusting peer $k ${pubKeys[$k]} for $n"
        pubKey=`echo ${pubKeys[$k]} | jq -r '.pubKey'`
        curl -u wasp:wasp -X PUT --header 'Content-Type: application/json' --data-raw ${pubKeys[$k]} $n/adm/peering/trusted/$pubKey
    done
done

# Init wasp cli
./wasp-cli init
# Set GoShimmer address
./wasp-cli set goshimmer.api $goshimer_api
# Set wasp addresses
length=`echo $all_ports | jq -c '. | length'`
for (( n=0; n<$length; n++ ))
do
    node=`echo $all_ports | jq -c --argjson index $n '.[$index].Services'`
    ./wasp-cli set "wasp.$n.api" `echo $node | jq -rc '.api'`
    ./wasp-cli set "wasp.$n.nanomsg" `echo $node | jq -rc '.nanomsg'`
    ./wasp-cli set "wasp.$n.peering" `echo $node | jq -rc '.peering'`
    echo "set node $n to $node"
done
