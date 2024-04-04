package dockerhelper

import (
	"fmt"
	"net/mail"
	"strings"

	"github.com/bridgex-eu/drs/internal/command"
)

var traefikImage = "traefik:2.11"

func validEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

func getCertEmail(cli *command.Cli) (string, error) {
	email := cli.Config.Email()

	if email != "" {
		return email, nil
	}

	email, err := command.Prompt(cli.In, cli.Out, "Enter an email that will be used for SSL certificate", "")
	if err != nil {
		return "", err
	}

	if !validEmail(email) {
		fmt.Printf("error!")
		return "", fmt.Errorf("Please, enter a valid email address")
	}

	cli.Config.SetEmail(email)

	return email, nil
}

func getTraefikArgs(cli *command.Cli) ([]string, error) {
	email, err := getCertEmail(cli)
	if err != nil {
		return []string{}, err
	}

	return []string{
		"--providers.docker",
		"--providers.docker.exposedByDefault=false",
		"--log.level=DEBUG",
		"--entrypoints.web.address=:80",
		"--entrypoints.web.http.redirections.entryPoint.to=websecure",
		"--entrypoints.web.http.redirections.entryPoint.scheme=https",
		"--entrypoints.websecure.address=:443",
		fmt.Sprintf("--certificatesresolvers.myresolver.acme.email=%s", email),
		"--certificatesresolvers.myresolver.acme.storage=./acme.json",
		"--certificatesresolvers.myresolver.acme.tlschallenge=true",
		"--accesslog=true",
	}, nil
}

var traefikVolumes = []string{
	"/var/run/docker.sock:/var/run/docker.sock",
}

var traefikPorts = []string{
	"80:80",
	"443:443",
}

var traefikName = "drs-traefik"

func runTraefik(client *command.MachineClient, cli *command.Cli) error {
	traefikArgs, err := getTraefikArgs(cli)
	if err != nil {
		return err
	}

	return RunContainer(client, traefikImage, traefikName, traefikPorts, traefikVolumes, nil, nil, traefikArgs...)
}

func GetTraefikRouteLabel(domain, containerPort string) []string {
	domainLabel := strings.ReplaceAll(domain, ".", "-")
	return []string{
		"traefik.docker.network=" + drsNetwork,
		"traefik.enable=true",
		fmt.Sprintf("traefik.http.routers.%s.rule=Host(`%s`)", domainLabel, domain),
		fmt.Sprintf("traefik.http.routers.%s.entrypoints=websecure", domainLabel),
		fmt.Sprintf("traefik.http.routers.%s.tls=true", domainLabel),
		fmt.Sprintf("traefik.http.routers.%s.tls.certresolver=myresolver", domainLabel),
		fmt.Sprintf("traefik.http.services.%s.loadbalancer.server.port=%s", domainLabel, containerPort),
	}
}
