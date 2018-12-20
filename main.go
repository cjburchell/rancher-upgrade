package main

import (
	"github.com/cjburchell/rancher-upgrade/Rancher"
)

func main() {

}

const UPGRADED = "Upgraded"
const ACTIVE = "Active"

func upgradeService(enviroment string, service string, confirm bool) error {

	client := rancher.Client{}.Environment(enviroment).Service(service)

	// checkServiceState(service, listener);
	dockerUUID := ""
	serviceUpgrade := rancher.ServiceUpgrade{
		InServiceStrategy: rancher.InServiceStrategy{
			LaunchConfig: rancher.LaunchConfig{
				ImageUuid: dockerUUID,
			},
		},
	}

	err := client.Upgrade(serviceUpgrade)
	if err != nil {
		return err
	}

	err = waitUntilServiceStateIs(client, service, UPGRADED)
	if err != nil {
		return err
	}

	if !confirm {
		return nil
	}

	err = client.FinishUpgrade()
	if err != nil {
		return err
	}

	return waitUntilServiceStateIs(client, service, ACTIVE)
}

func waitUntilServiceStateIs(client rancher.Client, serviceId string, state string) error {
	const timeout = 50
	timeoutMs := 1000 * timeout

	/*start := System.currentTimeMillis();
	current := System.currentTimeMillis();
	success := false
		for (current - start) < timeoutMs {
		checkService, err := client.Get()
		String state = checkService.get().getState()
		if (state.equalsIgnoreCase(targetState)) {
		listener.getLogger().println("current service state is " + targetState)
		success = true
		break
	}
		Thread.sleep(2000);
		current = System.currentTimeMillis();
	}
		if !success {
		return errors.New("timeout";
	}*/

	return nil
}
