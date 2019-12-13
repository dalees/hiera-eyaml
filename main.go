package main

import (
	"fmt"
	"io/ioutil"
	"path/filepath"

	"github.com/lyraproj/dgo/dgo"
	"github.com/lyraproj/hierasdk/hiera"
	"github.com/lyraproj/hierasdk/plugin"
	"github.com/lyraproj/hierasdk/register"
	"gopkg.in/yaml.v3"
)

func main() {

	register.LookupKey(`lookup_eyaml`, myLookupKey)

	// Start RESTful service that makes all registered functions available
	plugin.ServeAndExit()
}

/*// A fake main, used to simulate being a real plugin to allow debugging.
// Remove this and enable above main as required.
func main() {
	// make a fake context
	options := map[string]string{
		"path":              "/home/dale/dev/openstack-puppet-catalyst/hieradata/data/cloudci.yaml",
		"pkcs7_private_key": "/home/dale/dev/openstack-puppet-catalyst/hieradata/secure/keys/private_key.pkcs7.pem",
		"pkcs7_public_key":  "/home/dale/dev/openstack-puppet-catalyst/hieradata/secure/keys/public_key.pkcs7.pem",
	}
	optstr, ok := json.Marshal(options)
	if ok != nil {
		panic("Could not convert options to JSON")
	}
	v := url.Values{}
	v.Add("options", string(optstr))
	ctx := hiera.NewProviderContext(v)
	// call myLookupKey with some test data
	key := "keystone::admin_token"
	result := myLookupKey(ctx, key)

	// print output values
	fmt.Printf("Lookup of '%s' returned value '%s'\n", key, result)
}*/

/*func decryptValue(hc hiera.ProviderContext, value dgo.Value) {
	// decrypt temporary notes
	//   https://github.com/fullsailor/pkcs7
	//   https://github.com/mozilla-services/pkcs7

	// $ eyaml decrypt --pkcs7-private-key secure/keys/private_key.pkcs7.pem --pkcs7-public-key secure/keys/public_key.pkcs7.pem --eyaml data/cloudci.yaml
	// $ eyaml decrypt --string "ENC[PKCS7,Miixxx==]"
	// decrypt if required, stripping extra data.

	privateKey, ok := hc.StringOption(`pkcs7_private_key`)
	if !ok {
		panic(fmt.Errorf(`missing required provider option 'pkcs7_private_key'`))
	}
	publicKey, ok := hc.StringOption(`pkcs7_public_key`)
	if !ok {
		panic(fmt.Errorf(`missing required provider option 'pkcs7_private_key'`))
	}

}*/

func myLookupKey(hc hiera.ProviderContext, key string) dgo.Value {
	if key == `lookup_options` {
		return nil
	}

	path, ok := hc.StringOption(`path`)
	if !ok {
		panic(fmt.Errorf(`missing required provider option 'path'`))
	}
	// Parse the yaml file.
	filename, _ := filepath.Abs(path)
	yamlFile, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	// define an interface to store the data. Works okay for single level nesting.
	var config map[string]interface{}

	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		panic(err)
	}

	// Note: This currently only works with strings not complex data types.
	//       It will need to work with all, so we may need to use github.com/lyraproj/pcore/yaml
	// eg. profiles::icinga::users
	// eg. adjutant::plugin_settings
	if val, ok := config[key]; ok {
		// TODO: Perform decryption of val if required
		return hc.ToData(val)
	}

	return nil
}
