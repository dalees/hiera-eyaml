package main

import (
	"fmt"

	"github.com/lyraproj/dgo/dgo"
	"github.com/lyraproj/hierasdk/hiera"
	"github.com/lyraproj/hierasdk/plugin"
	"github.com/lyraproj/hierasdk/register"

	"gopkg.in/yaml.v2"
)
/*
import (
	"github.com/lyraproj/hiera/hieraapi"
	"github.com/lyraproj/issue/issue"
	"github.com/lyraproj/pcore/px"
	"github.com/lyraproj/pcore/types"
	"github.com/lyraproj/pcore/yaml"
)*/

func main() {

	register.LookupKey(`lookup_eyaml`, myLookupKey)

	// Start RESTful service that makes all registered functions available
	plugin.ServeAndExit()
}

/*
func yamlData(ctx hieraapi.ServerContext) px.OrderedMap {
	pv := ctx.Option(`path`)
	if pv == nil {
		panic(px.Error(hieraapi.MissingRequiredOption, issue.H{`option`: `path`}))
	}
	path := pv.String()
	if bin, ok := types.BinaryFromFile2(path); ok {
		v := yaml.Unmarshal(ctx.(hieraapi.ServerContext).Invocation(), bin.Bytes())
		if data, ok := v.(px.OrderedMap); ok {
			return data
		}
		panic(px.Error(hieraapi.YamlNotHash, issue.H{`path`: path}))
	}
	return px.EmptyMap
}
*/

func myLookupKey(hc hiera.ProviderContext, key string) dgo.Value {
	// No need to support this yet.
	if key == `lookup_options` {
		return nil
	}

	path, ok := hc.StringOption(`path`)
	hardcodedKey := "%{hiera('profiles::api::host')}"
	if !ok {
		panic(fmt.Errorf(`missing required provider option 'path'`))
	}
	// open our yaml file.
	// find the referenced key (nil if not exist)
	// $ eyaml decrypt --pkcs7-private-key secure/keys/private_key.pkcs7.pem --pkcs7-public-key secure/keys/public_key.pkcs7.pem --eyaml data/cloudci.yaml
	// $ eyaml decrypt --string "ENC[PKCS7,Miixxx==]"
	// decrypt if required, stripping extra data.
	// optimisation: keep last X files in memory, along with crypt keys?

	/* privateKey, ok := hc.StringOption(`pkcs7_private_key`)
	if !ok {
		panic(fmt.Errorf(`missing required provider option 'pkcs7_private_key'`))
	}
	publicKey, ok := hc.StringOption(`pkcs7_public_key`)
	if !ok {
		panic(fmt.Errorf(`missing required provider option 'pkcs7_private_key'`))
	}*/

	if key != "profiles::api::host" {
		return hc.ToData(hardcodedKey)
	}
	return nil

	return hc.ToData(path)

	// TODO: The below is taken from https://github.com/lyraproj/hiera_azure/blob/master/vaultlookupkey.go
	/*	vaultName, ok := hc.StringOption(`vault_name`)
		if !ok {
			panic(fmt.Errorf(`missing required provider option 'vault_name'`))
		}
		var authorizer autorest.Authorizer
		var err error
		if os.Getenv("AZURE_TENANT_ID") != "" && os.Getenv("AZURE_CLIENT_ID") != "" && os.Getenv("AZURE_CLIENT_SECRET") != "" {
			authorizer, err = auth.NewAuthorizerFromEnvironment()
		} else {
			authorizer, err = auth.NewAuthorizerFromCLI()
		}
		if err != nil {
			panic(err)
		}
		client := keyvault.New()
		client.Authorizer = authorizer
		resp, err := client.GetSecret(context.Background(), "https://"+vaultName+".vault.azure.net", key, "")
		if err != nil {
			if ResponseWasStatusCode(resp.Response, http.StatusNotFound) {
				return nil
			}
			panic(err)
		}
		return hc.ToData(*resp.Value)*/
}
