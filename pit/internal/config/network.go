package config

import (
	"fmt"
	"os"
	"strings"
)

// Network is the product environment. One workspace uses exactly one network.
type Network string

const (
	Mainnet Network = "mainnet"
	Testnet Network = "testnet"
)

type Chain struct {
	Network          Network
	ChainID          int64
	RPC              string
	Explorer         string
	Serving          string
	Ledger           string
	StorageIndexer   string
	StorageFlow      string
	Identity8004     string
	Reputation8004   string
	DeskID           string
	HLInfo           string
	HLExchange       string
	RouterCatalog    string // catalog only — never used for the private book
	PCUI             string
}

func MainnetChain() Chain {
	return Chain{
		Network:        Mainnet,
		ChainID:        16661,
		RPC:            envOr("PIT_CHAIN_RPC_URL", "https://evmrpc.0g.ai"),
		Explorer:       envOr("PIT_EXPLORER_BASE", "https://chainscan.0g.ai"),
		Serving:        envOr("PIT_SERVING_CONTRACT", "0x47340d900bdFec2BD393c626E12ea0656F938d84"),
		Ledger:         envOr("PIT_LEDGER_CONTRACT", "0x2dE54c845Cd948B72D2e32e39586fe89607074E3"),
		StorageIndexer: envOr("PIT_STORAGE_INDEXER", "https://indexer-storage-turbo.0g.ai"),
		StorageFlow:    envOr("PIT_STORAGE_FLOW", "0x62D4144dB0F0a6fBBaeb6296c785C71B3D57C526"),
		Identity8004:   envOr("PIT_8004_IDENTITY", "0x8004A169FB4a3325136EB29fA0ceB6D2e539a432"),
		Reputation8004: envOr("PIT_8004_REPUTATION", "0x8004BAa17C55a88189AE136b182e5fdA19dE9b63"),
		DeskID:         envOr("PIT_DESK_ID_CONTRACT", "0xfdB3a8D39F1E2b77a8261b359eABaaa2F08f8c35"),
		HLInfo:         envOr("PIT_HL_INFO_URL", "https://api.hyperliquid.xyz/info"),
		HLExchange:     envOr("PIT_HL_EXCHANGE_URL", "https://api.hyperliquid.xyz/exchange"),
		RouterCatalog:  envOr("PIT_ROUTER_URL", "https://router-api.0g.ai"),
		PCUI:           "https://pc.0g.ai",
	}
}

func TestnetChain() Chain {
	return Chain{
		Network:        Testnet,
		ChainID:        16602,
		RPC:            envOr("PIT_TESTNET_RPC_URL", "https://evmrpc-testnet.0g.ai"),
		Explorer:       envOr("PIT_TESTNET_EXPLORER", "https://chainscan-galileo.0g.ai"),
		Serving:        envOr("PIT_TESTNET_SERVING_CONTRACT", "0xa79F4c8311FF93C06b8CfB403690cc987c93F91E"),
		Ledger:         envOr("PIT_TESTNET_LEDGER_CONTRACT", "0xE70830508dAc0A97e6c087c75f402f9Be669E406"),
		StorageIndexer: envOr("PIT_TESTNET_STORAGE_INDEXER", "https://indexer-storage-testnet-turbo.0g.ai"),
		StorageFlow:    envOr("PIT_TESTNET_STORAGE_FLOW", "0x22E03a6A89B950F1c82ec5e74F8eCa321a105296"),
		Identity8004:   envOr("PIT_TESTNET_8004_IDENTITY", "0x8004A818BFB912233c491871b3d84c89A494BD9e"),
		Reputation8004: envOr("PIT_TESTNET_8004_REPUTATION", "0x8004B663056A597Dffe9eCcC1965A193B7388713"),
		DeskID:         os.Getenv("PIT_TESTNET_DESK_ID_CONTRACT"),
		HLInfo:         envOr("PIT_HL_TESTNET_INFO", "https://api.hyperliquid-testnet.xyz/info"),
		HLExchange:     envOr("PIT_HL_TESTNET_EXCHANGE", "https://api.hyperliquid-testnet.xyz/exchange"),
		RouterCatalog:  envOr("PIT_TESTNET_ROUTER_URL", "https://router-api-testnet.integratenetwork.work"),
		PCUI:           "https://pc.testnet.0g.ai",
	}
}

func ParseNetwork(s string) (Network, error) {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "mainnet", "aristotle", "16661":
		return Mainnet, nil
	case "testnet", "galileo", "16602":
		return Testnet, nil
	default:
		return "", fmt.Errorf("unknown network %q", s)
	}
}

func For(n Network) Chain {
	if n == Testnet {
		return TestnetChain()
	}
	return MainnetChain()
}

func envOr(k, def string) string {
	if v := strings.TrimSpace(os.Getenv(k)); v != "" {
		return v
	}
	return def
}

// GuardFallbacks exits the process if mocks are enabled.
func GuardFallbacks() error {
	v := strings.ToLower(strings.TrimSpace(os.Getenv("PIT_ALLOW_FALLBACKS")))
	if v != "" && v != "false" && v != "0" {
		return fmt.Errorf("PIT_ALLOW_FALLBACKS must be false")
	}
	return nil
}

// RejectGlobalUser is a hard kill: the product never treats a fixture master as every user.
func RejectGlobalUser(workspaceOwner string) error {
	master := strings.TrimSpace(os.Getenv("PIT_MASTER_ADDRESS"))
	if master == "" {
		return nil
	}
	if strings.EqualFold(workspaceOwner, master) && os.Getenv("PIT_PRODUCT_MODE") == "true" {
		return fmt.Errorf("PIT_MASTER_ADDRESS cannot be the product user identity")
	}
	return nil
}
