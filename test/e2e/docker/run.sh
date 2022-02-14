#!/usr/bin/env bash

docker-compose down

docker volume create simnet_lnd_alice
docker volume create simnet_lnd_bob

docker-compose run -d --name alice --volume simnet_lnd_alice:/root/.lnd lnd

# find a way to check if Alice is started
sleep 5

export MINING_ADDRESS=$(docker exec -it alice lncli --network=simnet newaddress np2wkh | grep address | sed -r 's/^[^:]*:(.*)$/\1/' | tr -d '\r' | xargs echo -n)

docker-compose up -d btcd

# find a way to check if Alice is started
sleep 5

docker exec -it btcd /start-btcctl.sh generate 400

docker exec -it btcd /start-btcctl.sh getblockchaininfo | grep -A 1 segwit

sleep 5

docker exec -it alice lncli --network=simnet walletbalance

docker-compose run -d --name bob --volume simnet_lnd_bob:/root/.lnd lnd

sleep 5

# docker exec -it bob lncli --network=simnet getinfo

export BOB_PUBKEY=$(docker exec -it bob lncli --network=simnet getinfo | grep identity | sed -r 's/^[^:]*:(.*)$/\1/' | tr -d '\r' | tr -d ',' | xargs echo -n)

echo $BOB_PUBKEY

export BOB_IP=$(docker inspect bob | grep IPAddress | tail -n1 | sed -r 's/^[^:]*:(.*)$/\1/' | tr -d '\r' | tr -d ',' | xargs echo -n)

echo $BOB_PUBKEY@$BOB_IP

docker exec -it alice lncli --network=simnet connect $BOB_PUBKEY@$BOB_IP

docker exec -it alice lncli --network=simnet listpeers

sleep 2

docker exec -it alice lncli --network=simnet openchannel --node_key=$BOB_PUBKEY --local_amt=1000000

docker exec -it btcd /start-btcctl.sh generate 3

sleep 2

docker exec -it alice lncli --network=simnet listchannels

export PAY_HASH=$(docker exec -it bob lncli --network=simnet addinvoice --amt=10000 | grep payment_request | sed -r 's/^[^:]*:(.*)$/\1/' | tr -d '\r' | tr -d ',' | xargs echo -n)

echo $PAY_HASH

docker exec -it alice lncli --network=simnet sendpayment --force --pay_req=$PAY_HASH

docker exec -it alice lncli --network=simnet channelbalance

docker exec -it bob lncli --network=simnet channelbalance
