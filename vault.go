package newrelic

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"emperror.dev/errors"
	"github.com/hashicorp/vault/api"
)

type VaultClient struct {
	Client *api.Client
}

type VaultSettings struct {
	VaultAddr  string
	VaultToken string
}

const (
	// CredentialsPath https://github.com/Boostport/kubernetes-vault#init-container-configuration
	CredentialsPath       string = "/var/run/secrets/boostport.com"
	VaultSettingsFileName string = "vault-token"
)

type VaultFile struct {
	ClientToken string `json:"clientToken"`
	VaultAddr   string `json:"vaultAddr"`
}

func getVaultSettings() *VaultSettings {
	var vaultFile VaultFile
	jsonFile, err := os.Open(CredentialsPath + VaultSettingsFileName)
	if err != nil {
		fmt.Println(err)
	}
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)
	json.Unmarshal([]byte(byteValue), &vaultFile)
	if vaultFile.ClientToken == "" {
		fmt.Println("token not found or empty")
	}
	return &VaultSettings{
		VaultAddr:  vaultFile.VaultAddr,
		VaultToken: vaultFile.ClientToken,
	}
}

// NewVault vault client
func initClient(vaultSettings *VaultSettings) (*VaultClient, error) {
	client, err := api.NewClient(&api.Config{Address: vaultSettings.VaultAddr})
	if err != nil {
		return nil, errors.Wrap(err, "could not create vault client")
	}
	client.SetToken(vaultSettings.VaultToken)

	v := &VaultClient{
		Client: client,
	}

	return v, nil
}

func (vc *VaultClient) readSecret(secretName string, secretPath string) string {
	secret, err := vc.Client.Logical().Read(secretPath)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	m, ok := secret.Data["data"].(map[string]interface{})
	if !ok {
		return ""
	}
	str, ok := m[secretName].(string)
	if !ok {
		return "Key not found"
	}
	return str
}
