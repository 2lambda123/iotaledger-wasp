./wasp-cli chain deploy --committee=0,1,2,3,4 --peers=0,1,2,3,4,5,6,7 --quorum=3 --chain=testchain --description="Test Chain"
./wasp-cli chain deposit IOTA:10000
./wasp-cli chain deploy-contract wasmtime fairroulette "fairroulette"  fairroulette_bg.wasm
#./wasp-cli --verbose chain post-request fairroulete placeBet string number int 2