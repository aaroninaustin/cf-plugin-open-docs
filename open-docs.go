package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"
	"io"

	"github.com/cloudfoundry/cli/plugin"
	"github.com/skratchdot/open-golang/open"
)

type serviceInstanceResponse struct {
	Entity serviceInstance `json:"entity"`
}

type serviceInstance struct {
	DashboardURL string `json:"dashboard_url"`
}

func main() {
	plugin.Start(&OpenDocsPlugin{})
}

// OpenPlugin empty struct for plugin
type OpenDocsPlugin struct{}

// Run of seeder plugin
func (plugin OpenDocsPlugin) Run(cliConnection plugin.CliConnection, args []string) {
	err := checkArgs(cliConnection, args)
	if err != nil {
		os.Exit(1)
	}
	if args[0] == "docs" {
		plugin.runAppOpen(cliConnection, args)
	}
}

func (OpenDocsPlugin) GetMetadata() plugin.PluginMetadata {
    return plugin.PluginMetadata{
        Name:    "docs",
        Version: plugin.VersionType{Major: 1, Minor: 1, Build: 0},
        Commands: []plugin.Command{
            {
                Name:     "docs",
                HelpText: "open app swagger url in browser",
                UsageDetails: plugin.Usage{
                    Usage: "docs <appname>",
                },
            },
        },
    }
}

func (plugin OpenDocsPlugin) runAppOpen(cliConnection plugin.CliConnection, args []string) {
	output, err := cliConnection.CliCommandWithoutTerminalOutput("app", args[1])
	if err != nil {
		fmt.Fprintln(os.Stdout, "error: app does not exist")
		os.Exit(1)
	}
	url, err := getUrlFromOutput(output)
	if err != nil {
		fmt.Fprintln(os.Stdout, "error: ", err)
		os.Exit(1)
	}

	open.Run(multiRoutesMenu(os.Stdin, url))

}

func getUrlFromOutput(output []string) ([]string, error) {
	urls := []string{}
	for _, line := range output {
		splitLine := strings.Split(strings.TrimSpace(line), " ")
		if splitLine[0] == "urls:" {
			if len(splitLine) > 1 {
				for p := 1; p < len(splitLine); p++ {
					url := "https://" + strings.Trim(splitLine[p] + "/v1/docs/", ",")
					url = strings.TrimSpace(url)
					urls = append(urls, url)
				}

			} else if len(splitLine) == 1 {
				return []string{""}, errors.New("App has no route")
			}
		}
	}
	return urls, nil
}

func multiRoutesMenu(input io.Reader, urls []string) string {
	if len(urls) == 1 {
		return urls[0]
	} else {
		var choice int
		fmt.Println("Multiple routes detected. Please choose one: ")
		for u := 0; u < len(urls); u++ {
			fmt.Printf("%d - %s\n", u+1, urls[u])
		}
		fmt.Print("Enter route to open: ")
		fmt.Fscanf(input, "%d", &choice)
		for !(choice >= 1 && choice <= len(urls)) {
			fmt.Printf("Please enter valid number(1 to %d): ", len(urls))
			fmt.Fscanf(input, "%d", &choice)
		}
		return urls[choice-1]
	}
}

func (plugin OpenDocsPlugin) runServiceOpen(cliConnection plugin.CliConnection, args []string) {
	output, err := cliConnection.CliCommandWithoutTerminalOutput("service", args[1], "--guid")
	if err != nil {
		fmt.Fprintln(os.Stdout, "error: service does not exist")
		os.Exit(1)
	}
	serviceInstanceGUID := strings.TrimSpace(output[0])

	output, err = cliConnection.CliCommandWithoutTerminalOutput("curl", fmt.Sprintf("/v2/service_instances/%s", serviceInstanceGUID))
	if err != nil {
		fmt.Fprintln(os.Stdout, "error: service does not exist")
		os.Exit(1)
	}
	jsonStr := ""
	for _, line := range output {
		jsonStr += line + "\n"
	}

	response := serviceInstanceResponse{}
	json.Unmarshal([]byte(jsonStr), &response)

	url := response.Entity.DashboardURL
	if url == "" {
		fmt.Println("No dashboard available")
	} else {
		open.Run(url)
	}
}

func checkArgs(cliConnection plugin.CliConnection, args []string) error {
	if len(args) < 2 {
		if args[0] == "docs" {
			cliConnection.CliCommand(args[0], "-h")
			return errors.New("Appname is needed")
		} else if args[0] == "service-open-docs" {
			cliConnection.CliCommand(args[0], "-h")
			return errors.New("Appname is needed")
		}
	}
	return nil
}
