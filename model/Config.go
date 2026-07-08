package model

import "github.com/ethereum/go-ethereum/common"

type ApplicationConfig struct {
	NodeURL                 string          `mapstructure:"nodeURL"`
	WSURL                   string          `mapstructure:"wsURL"`
	ContractAddress         string          `mapstructure:"contractAddress"`
	RelayHubContractAddress *common.Address `mapstructure:"relayHubContractAddress"`
	NodeKeyPath             string          `mapstructure:"nodeKeyPath"`
	NodeAddressPath         string          `mapstructure:"nodeAddressPath"`
	Key                     string          `mapstructure:"key"`
	Port                    string          `mapstructure:"port"`
	// NonceCacheTTL: vida máxima (segundos) de una entrada del caché de nonces por sender.
	// 0 o ausente = default (300s).
	NonceCacheTTL int64 `mapstructure:"nonceCacheTTL"`
}

type KeyStoreConfig struct {
	Agent string `mapstructure:"agent"`
}

type PassphraseConfig struct {
	Agent string `mapstructure:"agent"`
}

type SecurityConfig struct {
	PermissionsEnabled     bool   `mapstructure:"permissionsEnabled"`
	AccountContractAddress string `mapstructure:"accountContractAddress"`
}

type Config struct {
	Application ApplicationConfig `mapstructure:"application"`
	KeyStore    KeyStoreConfig    `mapstructure:"keystore"`
	Passphrase  PassphraseConfig  `mapstructure:"passphrase"`
	Security    SecurityConfig    `mapstructure:"security"`
}
