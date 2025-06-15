curl http://localhost:8080/generate-account | jq 

curl -X POST "http://localhost:8080/generate-stealth" -H "Content-Type: application/json" -d '{"pub_key": "0x04a5834a4f6ec6c8a4953225ff126760eb9480af64dc65b4af7cb2df158755a5a79a4ba7ec5b043c4ad3d5a353d3f64ad938dfa26c778201298b16b49368bc8b6a"}' | jq 

curl -X POST http://localhost:8080/recover-stealth-priv-key -H "Content-Type: application/json" -d '{ "recipient_privkey": "0xc97126efe6a9834ae551f70959c7777a2b57de6a190cfdb30aa3b9d3354347ae", "ephemeral_pubkey": "0x04786a24825903b0cb294785d2a13a5b0063364feab771f0a3a9746b4ad27ea6b299c87ac8d422edf1b9770567157ea7fc3d7cd0317014558a9c1c415a8c1fa883"}' | jq


SANCTION Curl:
curl -X POST http://localhost:8080/sanctions/check -H "Content-Type: application/json" -d '{"address": "0xAbc123"}' | jq
curl -X POST http://localhost:8080/sanctions/check -H "Content-Type: application/json" -d '{"address": "0xAbc123"}' | jq

curl -X POST http://localhost:8080/sanctions/remove -H "Content-Type: application/json" -d '{"address": "0xAbc123"}' | jq
curl -X POST http://localhost:8080/sanctions/check -H "Content-Type: application/json" -d '{"address": "0xAbc123"}' | jq

curl -X POST http://localhost:8080/sanctions/add -H "Content-Type: application/json" -d '{"address": "0xAbc123"}' | jq
curl -X POST http://localhost:8080/sanctions/check -H "Content-Type: application/json" -d '{"address": "0xAbc123"}' | jq 