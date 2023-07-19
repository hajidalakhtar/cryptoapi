package rpc

import (
	"context"
	"cryptoapi/domain"
	utils "cryptoapi/helper"
	"cryptoapi/infrastructure/wallet/repository/helper"
	"log"
	"math/big"
	"sync"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	hdwallet "github.com/miguelmota/go-ethereum-hdwallet"
)

type RPCWalletRepository struct {
	client *ethclient.Client
}

func NewWalletRepository(client *ethclient.Client) domain.WalletRepository {
	return &RPCWalletRepository{client: client}
}

func (r *RPCWalletRepository) GetBalanceFromMnemonic(ctx context.Context, mnemonic string) ([]domain.Balance, string, error) {
	wallet, err := hdwallet.NewFromMnemonic(mnemonic)
	if err != nil {
		log.Fatal(err)
	}

	path := hdwallet.MustParseDerivationPath("m/44'/60'/0'/0/0")
	account, err := wallet.Derive(path, false)
	if err != nil {
		log.Fatal(err)
	}

	bal, err := r.GetBalance(ctx, account.Address.Hex())
	if err != nil {
		return nil, "", err
	}

	return bal, account.Address.Hex(), nil
}

func (r *RPCWalletRepository) GenerateNewWallet(ctx context.Context) (string, string, error) {
	mnemonic, err := hdwallet.NewMnemonic(256)
	if err != nil {
		log.Fatal(err)
	}

	wallet, err := hdwallet.NewFromMnemonic(mnemonic)
	if err != nil {
		log.Fatal(err)
	}

	path := hdwallet.MustParseDerivationPath("m/44'/60'/0'/0/0")
	account, err := wallet.Derive(path, false)
	if err != nil {
		log.Fatal(err)
	}

	return account.Address.Hex(), mnemonic, nil
}

func (r *RPCWalletRepository) GetBalance(ctx context.Context, addr string) ([]domain.Balance, error) {
	address := common.HexToAddress(addr)

	balance, err := r.client.BalanceAt(context.Background(), address, nil)
	if err != nil {
		return nil, err
	}

	nativeTokenBalance := utils.ConvertBalanceToFloat(balance, 18)

	balanceResp := []domain.Balance{{Name: "BNB", Amount: nativeTokenBalance}}

	tokenSymbols := map[string]string{
		//"CAKE":         "0x0E09FaBB73Bd3Ade0a17ECC321fD13a19e81cE82",
		//"ETH":          "0xA72Ff2B799324B042AE379809eE54dACE99db3A5",
		//"ADA":          "0x3EE2200Efb3400fAbB9AacF31297cBdD1d435D47",
		//"ALICE":        "0xAC51066d7bEC65Dc4589368da368b212745d63E8",
		//"ALIX":         "0xaF6Bd11A6F8f9c44b9D18f5FA116E403db599f8E",
		//"ALPHA":        "0xa1faa113cbE53436Df28FF0aEe54275c13B40975",
		//"ALU":          "0x8263CD1601FE73C066bf49cc09841f35348e3be0",
		//"ATA":          "0xA2120b9e674d3fC3875f415A7DF52e382F141225",
		//"ATOM":         "0x0Eb3a705fc54725037CC9e008bDede697f62F335",
		//"AXS":          "0x715D400F88C167884bbCc41C5FeA407ed4D2f8A0",
		//"BABYDOGE":     "0xc748673057861a797275CD8A068AbB95A902e8de",
		//"BEAR":         "0xc3EAE9b061Aa0e1B9BD3436080Dc57D2d63FEdc1",
		//"BEL":          "0x8443f091997f06a61670B735ED92734F5628692F",
		//"BELT":         "0xE0e514c71282b6f4e823703a39374Cf58dc3eA4f",
		//"BIN":          "0xe56842Ed550Ff2794F010738554db45E60730371",
		//"BMON":         "0x08ba0619b1e7A582E0BCe5BBE9843322C954C340",
		//"BNX":          "0x8C851d1a123Ff703BD1f9dabe631b69902Df5f97",
		//"BP":           "0xACB8f52DC63BB752a51186D1c55868ADbFfEe9C1",
		//"BSCPAD":       "0x5A3010d4d8D3B5fB49f8B6E57FB9E48063f16700",
		//"BTTOLD":       "0x8595F9dA7b868b1822194fAEd312235E43007b49",
		//"C98":          "0xaEC945e04baF28b135Fa7c640f624f8D90F1C3a6",
		//"CHESS":        "0x20de22029ab63cf9A7Cf5fEB2b737Ca1eE4c82A6",
		//"CHR":          "0xf9CeC8d50f6c8ad3Fb6dcCEC577e05aA32B224FE",
		//"CP":           "0x82C19905B036bf4E329740989DCF6aE441AE26c1",
		//"DERC":         "0x373E768f79c820aA441540d254dCA6d045c6d25b",
		//"DODO":         "0x67ee3Cb086F8a16f34beE3ca72FAD36F7Db929e2",
		//"DOGE":         "0xbA2aE424d960c26247Dd6c32edC70B295c744C43",
		//"DPET":         "0xfb62AE373acA027177D1c18Ee0862817f9080d08",
		//"DRACE":        "0xA6c897CaaCA3Db7fD6e2D2cE1a00744f40aB87Bb",
		//"DRS":          "0xc8E8ecB2A5B5d1eCFf007BF74d15A86434aA0c5C",
		//"DVI":          "0x758FB037A375F17c7e195CC634D77dA4F554255B",
		//"ECC":          "0x8D047F4F57A190C96c8b9704B39A1379E999D82B",
		//"EPS":          "0xA7f552078dcC247C2684336020c03648500C6d9F",
		//"FARA":         "0xF4Ed363144981D3A65f42e7D0DC54FF9EEf559A1",
		//"FLOKI":        "0x2B3F34e9D4b127797CE6244Ea341a83733ddd6E4",
		//"FORM":         "0x25A528af62e56512A19ce8c3cAB427807c28CC19",
		//"FRONT":        "0x928e55daB735aa8260AF3cEDadA18B5f70C72f1b",
		//"GOLD":         "0xb3a6381070B1a15169DEA646166EC0699fDAeA79",
		//"HERO":         "0xD40bEDb44C081D2935eebA6eF5a3c8A31A1bBE13",
		//"HONEY":        "0xFa363022816aBf82f18a9C2809dCd2BB393F6AC5",
		//"HUNNY":        "0x565b72163f17849832A692A3c5928cc502f46D69",
		//"INJ":          "0xa2B726B1145A4773F68593CF171187d8EBe4d495",
		//"IOTX":         "0x9678E42ceBEb63F23197D726B29b1CB20d0064E5",
		//"KABY":         "0x02A40C048eE2607B5f5606e445CFc3633Fb20b58",
		//"KMON":         "0xc732B6586A93b6B7CF5FeD3470808Bc74998224D",
		//"LINA":         "0x762539b45A1dCcE3D36d080F74d1AED37844b878",
		//"LINK":         "0xF8A0BF9cF54Bb92F17374d9e9A321E6a111a51bD",
		//"MASK":         "0x2eD9a5C8C13b93955103B9a7C167B67Ef4d568a3",
		//"MBOX":         "0x3203c9E46cA618C8C1cE5dC67e7e9D75f5da2377",
		//"MINIFOOTBALL": "0xD024Ac1195762F6F13f8CfDF3cdd2c97b33B248b",
		//"MIST":         "0x68E374F856bF25468D365E539b700b648Bf94B67",
		//"MND":          "0x4c97c901B5147F8C1C7Ce3c5cF3eB83B44F244fE",
		//"MONI":         "0x9573c88aE3e37508f87649f87c4dd5373C9F31e0",
		//"NAFT":         "0xD7730681B1DC8f6F969166B29D8A5EA8568616a3",
		//"NBL":          "0xA67a13c9283Da5AABB199Da54a9Cb4cD8B9b16bA",
		//"NFTB":         "0xde3dbBE30cfa9F437b293294d1fD64B26045C71A",
		//"NRV":          "0x42F6f551ae042cBe50C739158b4f0CAC0Edb9096",
		//"ONE":          "0x03fF0ff224f904be3118461335064bB48Df47938",
		//"PAID":         "0xAD86d0E9764ba90DDD68747D64BFfBd79879a238",
		//"PETG":         "0x09607078980CbB0665ABa9c6D1B84b8eAD246aA0",
		//"PINK":         "0x9133049Fb1FdDC110c92BF5b7Df635abB70C89DC",
		//"PMON":         "0x1796ae0b0fa4862485106a0de9b654eFE301D0b2",
		//"POCO":         "0x394bBA8F309f3462b31238B3fd04b83F71A98848",
		//"POTS":         "0x3Fcca8648651E5b974DD6d3e50F61567779772A8",
		//"PVU":          "0x31471E0791fCdbE82fbF4C44943255e923F1b794",
		//"PWT":          "0xf3eDD4f14a018df4b6f02Bf1b2cF17A8120519A2",
		//"QBT":          "0x17B7163cf1Dbd286E262ddc68b553D899B93f526",
		//"RACA":         "0x12BB890508c125661E03b09EC06E404bc9289040",
		//"RAMP":         "0x8519EA49c997f50cefFa444d240fB655e89248Aa",
		//"REEF":         "0xF21768cCBC73Ea5B6fd3C687208a7c2def2d966e",
		//"RUSD":         "0x07663837218A003e66310a01596af4bf4e44623D",
		//"SFP":          "0xD41FDb03Ba84762dD66a0af1a6C8540FF1ba5dfb",
		//"SFUND":        "0x477bC8d23c634C154061869478bce96BE6045D12",
		//"SHI":          "0x7269d98Af4aA705e0B1A5D8512FadB4d45817d5a",
		//"SKILL":        "0x154A9F9cbd3449AD22FDaE23044319D6eF2a1Fab",
		//"SMON":         "0xAB15B79880f11cFfb58dB25eC2bc39d28c4d80d2",
		//"SPS":          "0x1633b7157e7638C4d6593436111Bf125Ee74703F",
		//"SUSHI":        "0x947950BcC74888a40Ffa2593C5798F11Fc9124C4",
		//"SXP":          "0x47BEAd2563dCBf3bF2c9407fEa4dC236fAbA485A",
		//"TKO":          "0x9f589e3eabe42ebC94A44727b3f3531C0c877809",
		//"TLM":          "0x2222227E22102Fe3322098e4CBfE18cFebD57c95",
		//"TPT":          "0xECa41281c24451168a37211F0bc2b8645AF45092",
		//"TRONPAD":      "0x1Bf7AedeC439D6BFE38f8f9b20CF3dc99e3571C4",
		//"TRX":          "0x85EAC5Ac2F758618dFa09bDbe0cf174e7d574D5B",
		//"TSC":          "0xA2a26349448ddAfAe34949a6Cc2cEcF78c0497aC",
		//"TUSD":         "0x14016E85a25aeb13065688cAFB43044C2ef86784",
		//"TWT":          "0x4B0F1812e5Df2A09796481Ff14017e6005508003",
		//"UNCL":         "0x0E8D5504bF54D9E44260f8d153EcD5412130CaBb",
		//"UNCX":         "0x09a6c44c3947B69E2B45F4D51b67E6a39ACfB506",
		//"UNI":          "0xBf5140A22578168FD562DCcF235E5D43A02ce9B1",
		//"UST":          "0x23396cF899Ca06c4472205fC903bDB4de249D6fC",
		//"VAI":          "0x4BD17003473389A42DAF6a0a729f6Fdb328BbBd7",
		//"WANA":         "0x339C72829AB7DD45C3C52f965E7ABe358dd8761E",
		//"WEYU":         "0xFAfD4CB703B25CB22f43D017e7e0d75FEBc26743",
		//"WIN":          "0xaeF0d72a118ce24feE3cD1d43d383897D05B4e99",
		//"XRP":          "0x1D2F0da169ceB9fC7B3144628dB156f3F6c60dBE",
		//"XWG":          "0x6b23C89196DeB721e6Fd9726E6C76E4810a464bc",
		//"YAY":          "0x524dF384BFFB18C0C8f3f43d012011F8F9795579",
		//"ZIN":          "0xFbe0b4aE6E5a200c36A341299604D5f71A5F0a48",
	}
	balanceChan := make(chan domain.Balance, len(tokenSymbols))
	var wg sync.WaitGroup

	for symbol, tokenAddr := range tokenSymbols {
		wg.Add(1)
		go func(symbol, tokenAddr string) {
			defer wg.Done()

			tokenBalance, err := helper.GetTokenBalance(r.client, addr, tokenAddr)
			if err != nil {
				log.Printf("Failed to get token balance: %v", err)
				return
			}

			balanceChan <- domain.Balance{Name: symbol, Amount: tokenBalance}
		}(symbol, tokenAddr)
	}

	go func() {
		wg.Wait()
		close(balanceChan)
	}()

	for balance := range balanceChan {
		balanceResp = append(balanceResp, balance)
	}

	return balanceResp, nil
}

func (r *RPCWalletRepository) Transfer(ctx context.Context, mnemonic string, toAddr common.Address, amount *big.Int, gasPrice *big.Int, gasLimit uint64) (string, error) {
	chainID, err := r.client.NetworkID(context.Background())
	if err != nil {
		log.Fatalf("Failed to get network ID: %v", err)
	}

	wallet, err := hdwallet.NewFromMnemonic(mnemonic)
	if err != nil {
		log.Fatal(err)
	}

	path := hdwallet.MustParseDerivationPath("m/44'/60'/0'/0/0")
	account, err := wallet.Derive(path, true)
	if err != nil {
		log.Fatal(err)
	}

	nonce, err := r.client.NonceAt(context.Background(), account.Address, nil)
	if err != nil {
		log.Fatalf("Failed to get nonce: %v", err)
	}

	gasPrice, err = r.client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatalf("Failed to get gas price: %v", err)
	}

	gasLimit = uint64(21000)
	var data []byte

	tx := types.NewTransaction(nonce, toAddr, amount, gasLimit, gasPrice, data)
	signedTx, err := wallet.SignTx(account, tx, chainID)
	if err != nil {
		log.Fatal(err)
	}

	err = r.client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		log.Fatalf("Failed to send the transaction: %v", err)
	}

	return tx.Hash().Hex(), nil
}
