// intention is having this guy as a web server instead of a console app
package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strconv"

	pi "github.com/TomasCruz/pluginboss/plugininfo"
	"github.com/hashicorp/go-plugin"
)

// internal global variable holding all of the plugin metainfo
// not strictly neccessary to have it global, but I have better things to do :)
var pluginCatalogue pi.PluginInfo

func main() {
	// We don't want to see the plugin logs.
	log.SetOutput(ioutil.Discard)

	// load all the plugin info first
	pluginCatalogue = loadRegisteredPlugins()

	// Plugin maps
	pluginClients := map[string]*plugin.Client{} // need references to *plugin.Clients so I can kill 'em
	converterPlugins := map[string]pi.ConverterPlugin{}

	// Creating plugin clients
	err := converterPlugin("conversion_one", pluginClients, converterPlugins)
	if err != nil {
		fmt.Println("Couldn't start the plugin conversion_one, err: ", err.Error())
		os.Exit(1)
	}

	err = converterPlugin("conversion_two", pluginClients, converterPlugins)
	if err != nil {
		fmt.Println("Couldn't start the plugin conversion_two, err: ", err.Error())
		os.Exit(1)
	}

	err = converterPlugin("conversion_three", pluginClients, converterPlugins)
	if err != nil {
		fmt.Println("Couldn't start the plugin conversion_three, err: ", err.Error())
	}

	err = converterPlugin("conversion_four", pluginClients, converterPlugins)
	if err != nil {
		// plugin four dies straight away, so just keep on going
		fmt.Println("Couldn't start the plugin conversion_four, err: ", err.Error())
	}

	// kill 'em all on exit
	defer killAllClients(pluginClients)

	out, err := pluginCall("conversion_one", converterPlugins, os.Args[1])
	if err != nil {
		// if plugin didn't succesfully do it's work, so just keep on going
		fmt.Println("Error: ", err.Error())
	} else {
		fmt.Printf("Done conversion of F -> C, %s Fahrenheit == %f Celsius.\n Will Exit(0) now\n", os.Args[1], out)
	}

	out, err = pluginCall("conversion_two", converterPlugins, os.Args[2])
	if err != nil {
		fmt.Println("Error: ", err.Error())
	} else {
		fmt.Printf("Done conversion of C -> F, %s Celsius == %f Fahrenheit.\n Will Exit(0) now\n", os.Args[2], out)
	}

	out, err = pluginCall("conversion_three", converterPlugins, os.Args[3])
	if err != nil {
		// plugin three will die here, but pluginboss just keeps on going
		fmt.Println("Error: ", err.Error())
	} else {
		fmt.Printf("Done conversion %s == %f.\n Will Exit(0) now\n", os.Args[3], out)
	}

	out, err = pluginCall("conversion_four", converterPlugins, os.Args[4])
	if err != nil {
		fmt.Println("Error: ", err.Error())
	} else {
		fmt.Printf("Done conversion %s == %f.\n Will Exit(0) now\n", os.Args[4], out)
	}

	// proof that plugin engine and plugin one are still working normally, should a long-lived Client be needed
	out, err = pluginCall("conversion_one", converterPlugins, os.Args[5])
	if err != nil {
		fmt.Println("Error: ", err.Error())
	} else {
		fmt.Printf("Done conversion of F -> C, %s Fahrenheit == %f Celsius.\n Will Exit(0) now\n", os.Args[5], out)
	}

	os.Exit(0)
}

func killAllClients(pluginClients map[string]*plugin.Client) {
	for _, c := range pluginClients {
		c.Kill()
	}
}

func newCataloguedClient(pluginName string, pluginCatalogue pi.PluginInfo) *plugin.Client {
	plInfo := pluginCatalogue.PluginDataMap[pluginName]

	return plugin.NewClient(&plugin.ClientConfig{
		HandshakeConfig:  plInfo.Handshake,
		Plugins:          map[string]plugin.Plugin{pluginName: plInfo.Plugin},
		Cmd:              exec.Command("sh", "-c", plInfo.Cmd),
		AllowedProtocols: []plugin.Protocol{plugin.ProtocolGRPC},
	})
}

func loadRegisteredPlugins() pi.PluginInfo {
	catalogue := pi.PluginInfo{
		// initialize
		PluginDataMap: map[string]pi.PluginData{},
	}

	// hard-coded list of plugins available to pluginboss
	insertPlugin(&catalogue, "conversion_one", "bin/plugin_one", 1, "PLUGIN_ONE", "hello")
	insertPlugin(&catalogue, "conversion_two", "bin/plugin_two", 1, "PLUGIN_TWO", "hello")
	insertPlugin(&catalogue, "conversion_three", "bin/plugin_three", 1, "PLUGIN_THREE", "hello")
	insertPlugin(&catalogue, "conversion_four", "bin/plugin_four", 1, "PLUGIN_FOUR", "hello")

	return catalogue
}

func insertPlugin(pc *pi.PluginInfo, name, cmd string,
	protocolVersion uint, magicKey, magicValue string) {

	pc.PluginDataMap[name] = pi.PluginData{
		Cmd: cmd,
		Handshake: plugin.HandshakeConfig{
			ProtocolVersion:  protocolVersion,
			MagicCookieKey:   magicKey,
			MagicCookieValue: magicValue,
		},
		Plugin: &pi.ConverterGRPCPlugin{},
	}
}

func converterPlugin(pluginName string, pluginClients map[string]*plugin.Client, converterPlugins map[string]pi.ConverterPlugin) (err error) {
	_, ok := pluginClients[pluginName]
	if !ok {
		pluginClients[pluginName] = newCataloguedClient(pluginName, pluginCatalogue)
	}

	// Connect via RPC
	// Using plugin's client implies launching the plugin process
	// If it dies, we have no prob with it TODO
	rpcClient, err := pluginClients[pluginName].Client()
	if err != nil {
		return
	}

	// Request the plugin
	raw, err := rpcClient.Dispense(pluginName)
	if err != nil {
		return
	}

	// We should now have ConverterPlugin going. Commands issued to the plugin go there over a gRPC connection
	converterPlugins[pluginName], ok = raw.(pi.ConverterPlugin)
	if !ok {
		err = errors.New("Plugin not ConverterPlugin!")
		return
	}

	return
}

func pluginCall(pluginName string, converterPlugins map[string]pi.ConverterPlugin, arg string) (out float64, err error) {
	in, err := strconv.ParseFloat(arg, 64)
	if err != nil {
		fmt.Println("ParseFloat error:", err.Error())
		os.Exit(1)
	}

	cp, ok := converterPlugins[pluginName]
	if !ok {
		err = errors.New(fmt.Sprintf("Couldn't get plugin %s", pluginName))
		return
	}

	out, err = cp.Convert(in)
	return
}
