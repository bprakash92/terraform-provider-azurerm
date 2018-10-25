package azurerm

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"

	"github.com/Azure/azure-sdk-for-go/services/containerinstance/mgmt/2018-04-01/containerinstance"
)

var emptyCheck = func(cg containerinstance.ContainerGroup) error { return nil }

func TestAccAzureRMContainerGroup_imageRegistryCredentials(t *testing.T) {
	resourceName := "azurerm_container_group.test"
	ri := acctest.RandInt()

	config := testAccAzureRMContainerGroup_imageRegistryCredentials(ri, testLocation())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMContainerGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMContainerGroupExists(resourceName, emptyCheck),
					resource.TestCheckResourceAttr(resourceName, "image_registry_credential.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "image_registry_credential.0.server", "hub.docker.com"),
					resource.TestCheckResourceAttr(resourceName, "image_registry_credential.0.username", "yourusername"),
					resource.TestCheckResourceAttr(resourceName, "image_registry_credential.0.password", "yourpassword"),
					resource.TestCheckResourceAttr(resourceName, "image_registry_credential.1.server", "mine.acr.io"),
					resource.TestCheckResourceAttr(resourceName, "image_registry_credential.1.username", "acrusername"),
					resource.TestCheckResourceAttr(resourceName, "image_registry_credential.1.password", "acrpassword"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"image_registry_credential.0.password",
					"image_registry_credential.1.password",
				},
			},
		},
	})
}

func TestAccAzureRMContainerGroup_imageRegistryCredentialsUpdate(t *testing.T) {
	resourceName := "azurerm_container_group.test"
	ri := acctest.RandInt()

	config := testAccAzureRMContainerGroup_imageRegistryCredentials(ri, testLocation())
	updated := testAccAzureRMContainerGroup_imageRegistryCredentialsUpdated(ri, testLocation())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMContainerGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMContainerGroupExists(resourceName, func(cg containerinstance.ContainerGroup) error {
						if cg.Containers == nil || len(*cg.Containers) != 1 {
							return fmt.Errorf("unexpected number of containers created")
						}
						containers := *cg.Containers
						if containers[0].Ports == nil || len(*containers[0].Ports) != 1 {
							return fmt.Errorf("unexpected number of ports created")
						}
						ports := *containers[0].Ports
						if *ports[0].Port != 5443 {
							return fmt.Errorf("expected port to be 80, instead got %d", ports[0].Port)
						}
						if string(ports[0].Protocol) != "UDP" {
							return fmt.Errorf("expected port to be UDP, instead got %s", string(ports[0].Protocol))
						}
						return nil
					}),
					resource.TestCheckResourceAttr(resourceName, "image_registry_credential.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "image_registry_credential.0.server", "hub.docker.com"),
					resource.TestCheckResourceAttr(resourceName, "image_registry_credential.0.username", "yourusername"),
					resource.TestCheckResourceAttr(resourceName, "image_registry_credential.0.password", "yourpassword"),
					resource.TestCheckResourceAttr(resourceName, "image_registry_credential.1.server", "mine.acr.io"),
					resource.TestCheckResourceAttr(resourceName, "image_registry_credential.1.username", "acrusername"),
					resource.TestCheckResourceAttr(resourceName, "image_registry_credential.1.password", "acrpassword"),
					resource.TestCheckResourceAttr(resourceName, "container.0.port", "5443"),
					resource.TestCheckResourceAttr(resourceName, "container.0.protocol", "UDP"),
				),
			},
			{
				Config: updated,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMContainerGroupExists(resourceName, func(cg containerinstance.ContainerGroup) error {
						if cg.Containers == nil || len(*cg.Containers) != 1 {
							return fmt.Errorf("unexpected number of containers created")
						}
						containers := *cg.Containers
						if containers[0].Ports == nil || len(*containers[0].Ports) != 1 {
							return fmt.Errorf("unexpected number of ports created")
						}
						ports := *containers[0].Ports
						if *ports[0].Port != 80 {
							return fmt.Errorf("expected port to be 80, instead got %d", ports[0].Port)
						}
						if string(ports[0].Protocol) != "TCP" {
							return fmt.Errorf("expected port to be TCP, instead got %s", string(ports[0].Protocol))
						}
						return nil
					}),
					resource.TestCheckResourceAttr(resourceName, "image_registry_credential.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "image_registry_credential.0.server", "hub.docker.com"),
					resource.TestCheckResourceAttr(resourceName, "image_registry_credential.0.username", "updatedusername"),
					resource.TestCheckResourceAttr(resourceName, "image_registry_credential.0.password", "updatedpassword"),
					resource.TestCheckResourceAttr(resourceName, "container.0.ports.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "container.0.ports.0.port", "80"),
					resource.TestCheckResourceAttr(resourceName, "container.0.ports.0.protocol", "TCP"),
				),
			},
		},
	})
}

func TestAccAzureRMContainerGroup_linuxBasic(t *testing.T) {
	resourceName := "azurerm_container_group.test"
	ri := acctest.RandInt()

	config := testAccAzureRMContainerGroup_linuxBasic(ri, testLocation())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMContainerGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMContainerGroupExists(resourceName, func(cg containerinstance.ContainerGroup) error {
						if cg.Containers == nil || len(*cg.Containers) != 1 {
							return fmt.Errorf("unexpected number of containers created")
						}
						containers := *cg.Containers
						if containers[0].Ports == nil || len(*containers[0].Ports) != 1 {
							return fmt.Errorf("unexpected number of ports created")
						}
						ports := *containers[0].Ports
						if *ports[0].Port != 80 {
							return fmt.Errorf("expected port to be 80, instead got %d", ports[0].Port)
						}
						if string(ports[0].Protocol) != "TCP" {
							return fmt.Errorf("expected port to be TCP, instead got %s", string(ports[0].Protocol))
						}
						return nil
					}),
					resource.TestCheckResourceAttr(resourceName, "container.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "os_type", "Linux"),
					resource.TestCheckResourceAttr(resourceName, "container.0.port", "80"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"image_registry_credential.0.password",
					"image_registry_credential.1.password",
				},
			},
		},
	})
}

func TestAccAzureRMContainerGroup_linuxBasicUpdate(t *testing.T) {
	resourceName := "azurerm_container_group.test"
	ri := acctest.RandInt()

	config := testAccAzureRMContainerGroup_linuxBasic(ri, testLocation())
	updatedConfig := testAccAzureRMContainerGroup_linuxBasicUpdated(ri, testLocation())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMContainerGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMContainerGroupExists(resourceName, emptyCheck),
					resource.TestCheckResourceAttr(resourceName, "container.#", "1"),
				),
			},
			{
				Config: updatedConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMContainerGroupExists(resourceName, func(cg containerinstance.ContainerGroup) error {
						if cg.Containers == nil || len(*cg.Containers) != 1 {
							return fmt.Errorf("unexpected number of containers created")
						}
						containers := *cg.Containers
						if containers[0].Ports == nil || len(*containers[0].Ports) != 2 {
							return fmt.Errorf("unexpected number of ports created")
						}
						ports := *containers[0].Ports
						if *ports[0].Port != 80 {
							return fmt.Errorf("expected port to be 80, instead got %d", ports[0].Port)
						}
						if string(ports[0].Protocol) != "TCP" {
							return fmt.Errorf("expected port to be TCP, instead got %s", string(ports[0].Protocol))
						}
						if *ports[1].Port != 5443 {
							return fmt.Errorf("expected port to be 80, instead got %d", ports[0].Port)
						}
						if string(ports[1].Protocol) != "UDP" {
							return fmt.Errorf("expected port to be UDP, instead got %s", string(ports[0].Protocol))
						}
						return nil
					}),
					resource.TestCheckResourceAttr(resourceName, "container.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "container.0.ports.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "container.0.ports.0.port", "80"),
					resource.TestCheckResourceAttr(resourceName, "container.0.ports.0.protocol", "TCP"),
					resource.TestCheckResourceAttr(resourceName, "container.0.ports.1.port", "5443"),
					resource.TestCheckResourceAttr(resourceName, "container.0.ports.1.protocol", "UDP"),
				),
			},
		},
	})
}

func TestAccAzureRMContainerGroup_linuxComplete(t *testing.T) {
	resourceName := "azurerm_container_group.test"
	ri := acctest.RandInt()

	config := testAccAzureRMContainerGroup_linuxComplete(ri, testLocation())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMContainerGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMContainerGroupExists(resourceName, func(cg containerinstance.ContainerGroup) error {
						if cg.Containers == nil || len(*cg.Containers) != 1 {
							return fmt.Errorf("unexpected number of containers created")
						}
						containers := *cg.Containers
						if containers[0].Ports == nil || len(*containers[0].Ports) != 1 {
							return fmt.Errorf("unexpected number of ports created")
						}
						ports := *containers[0].Ports
						if *ports[0].Port != 80 {
							return fmt.Errorf("expected port to be 80, instead got %d", ports[0].Port)
						}
						if string(ports[0].Protocol) != "TCP" {
							return fmt.Errorf("expected port to be TCP, instead got %s", string(ports[0].Protocol))
						}
						return nil
					}),
					resource.TestCheckResourceAttr(resourceName, "container.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "container.0.ports.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "container.0.ports.0.port", "80"),
					resource.TestCheckResourceAttr(resourceName, "container.0.ports.0.protocol", "TCP"),
					resource.TestCheckResourceAttr(resourceName, "container.0.command", "/bin/bash -c ls"),
					resource.TestCheckResourceAttr(resourceName, "container.0.commands.#", "3"),
					resource.TestCheckResourceAttr(resourceName, "container.0.commands.0", "/bin/bash"),
					resource.TestCheckResourceAttr(resourceName, "container.0.commands.1", "-c"),
					resource.TestCheckResourceAttr(resourceName, "container.0.commands.2", "ls"),
					resource.TestCheckResourceAttr(resourceName, "container.0.environment_variables.%", "2"),
					resource.TestCheckResourceAttr(resourceName, "container.0.environment_variables.foo", "bar"),
					resource.TestCheckResourceAttr(resourceName, "container.0.environment_variables.foo1", "bar1"),
					resource.TestCheckResourceAttr(resourceName, "container.0.secure_environment_variables.%", "2"),
					resource.TestCheckResourceAttr(resourceName, "container.0.secure_environment_variables.secureFoo", "secureBar"),
					resource.TestCheckResourceAttr(resourceName, "container.0.secure_environment_variables.secureFoo1", "secureBar1"),
					resource.TestCheckResourceAttr(resourceName, "container.0.volume.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "container.0.volume.0.mount_path", "/aci/logs"),
					resource.TestCheckResourceAttr(resourceName, "container.0.volume.0.name", "logs"),
					resource.TestCheckResourceAttr(resourceName, "container.0.volume.0.read_only", "false"),
					resource.TestCheckResourceAttr(resourceName, "os_type", "Linux"),
					resource.TestCheckResourceAttr(resourceName, "restart_policy", "OnFailure"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"container.0.volume.0.storage_account_key",
					"container.0.secure_environment_variables.%",
					"container.0.secure_environment_variables.secureFoo",
					"container.0.secure_environment_variables.secureFoo1",
				},
			},
		},
	})
}

func TestAccAzureRMContainerGroup_windowsBasic(t *testing.T) {
	resourceName := "azurerm_container_group.test"
	ri := acctest.RandInt()

	config := testAccAzureRMContainerGroup_windowsBasic(ri, testLocation())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMContainerGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMContainerGroupExists(resourceName, func(cg containerinstance.ContainerGroup) error {
						if cg.Containers == nil || len(*cg.Containers) != 1 {
							return fmt.Errorf("unexpected number of containers created")
						}
						containers := *cg.Containers
						if containers[0].Ports == nil || len(*containers[0].Ports) != 2 {
							return fmt.Errorf("unexpected number of ports created")
						}
						ports := *containers[0].Ports
						if *ports[0].Port != 80 {
							return fmt.Errorf("expected port to be 80, instead got %d", ports[0].Port)
						}
						if string(ports[0].Protocol) != "TCP" {
							return fmt.Errorf("expected port to be TCP, instead got %s", string(ports[0].Protocol))
						}
						if *ports[1].Port != 443 {
							return fmt.Errorf("expected port to be 80, instead got %d", ports[0].Port)
						}
						if string(ports[1].Protocol) != "TCP" {
							return fmt.Errorf("expected port to be TCP, instead got %s", string(ports[0].Protocol))
						}
						return nil
					}),
					resource.TestCheckResourceAttr(resourceName, "container.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "os_type", "Windows"),
					resource.TestCheckResourceAttr(resourceName, "container.0.ports.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "container.0.ports.0.port", "80"),
					resource.TestCheckResourceAttr(resourceName, "container.0.ports.0.protocol", "TCP"),
					resource.TestCheckResourceAttr(resourceName, "container.0.ports.0.port", "443"),
					resource.TestCheckResourceAttr(resourceName, "container.0.ports.0.protocol", "TCP"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAzureRMContainerGroup_windowsComplete(t *testing.T) {
	resourceName := "azurerm_container_group.test"
	ri := acctest.RandInt()

	config := testAccAzureRMContainerGroup_windowsComplete(ri, testLocation())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMContainerGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMContainerGroupExists(resourceName, func(cg containerinstance.ContainerGroup) error {
						if cg.Containers == nil || len(*cg.Containers) != 1 {
							return fmt.Errorf("unexpected number of containers created")
						}
						containers := *cg.Containers
						if containers[0].Ports == nil || len(*containers[0].Ports) != 1 {
							return fmt.Errorf("unexpected number of ports created")
						}
						ports := *containers[0].Ports
						if *ports[0].Port != 80 {
							return fmt.Errorf("expected port to be 80, instead got %d", ports[0].Port)
						}
						if string(ports[0].Protocol) != "TCP" {
							return fmt.Errorf("expected port to be TCP, instead got %s", string(ports[0].Protocol))
						}
						return nil
					}),
					resource.TestCheckResourceAttr(resourceName, "container.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "container.0.ports.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "container.0.ports.0.port", "80"),
					resource.TestCheckResourceAttr(resourceName, "container.0.ports.0.protocol", "TCP"),
					resource.TestCheckResourceAttr(resourceName, "container.0.command", "cmd.exe echo hi"),
					resource.TestCheckResourceAttr(resourceName, "container.0.commands.#", "3"),
					resource.TestCheckResourceAttr(resourceName, "container.0.commands.0", "cmd.exe"),
					resource.TestCheckResourceAttr(resourceName, "container.0.commands.1", "echo"),
					resource.TestCheckResourceAttr(resourceName, "container.0.commands.2", "hi"),
					resource.TestCheckResourceAttr(resourceName, "container.0.environment_variables.%", "2"),
					resource.TestCheckResourceAttr(resourceName, "container.0.environment_variables.foo", "bar"),
					resource.TestCheckResourceAttr(resourceName, "container.0.environment_variables.foo1", "bar1"),
					resource.TestCheckResourceAttr(resourceName, "container.0.secure_environment_variables.%", "2"),
					resource.TestCheckResourceAttr(resourceName, "container.0.secure_environment_variables.secureFoo", "secureBar"),
					resource.TestCheckResourceAttr(resourceName, "container.0.secure_environment_variables.secureFoo1", "secureBar1"),
					resource.TestCheckResourceAttr(resourceName, "os_type", "Windows"),
					resource.TestCheckResourceAttr(resourceName, "restart_policy", "Never"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"container.0.secure_environment_variables.%",
					"container.0.secure_environment_variables.secureFoo",
					"container.0.secure_environment_variables.secureFoo1",
				},
			},
		},
	})
}

func testAccAzureRMContainerGroup_linuxBasic(ri int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_container_group" "test" {
  name                = "acctestcontainergroup-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  ip_address_type     = "public"
  os_type             = "Linux"

  container {
    name   = "hw"
    image  = "microsoft/aci-helloworld:latest"
    cpu    = "0.5"
    memory = "0.5"
    port   = 80
  }

  tags {
    environment = "Testing"
  }
}
`, ri, location, ri)
}

func testAccAzureRMContainerGroup_imageRegistryCredentials(ri int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_container_group" "test" {
  name                = "acctestcontainergroup-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  ip_address_type     = "public"
  os_type             = "Linux"

  container {
    name     = "hw"
    image    = "microsoft/aci-helloworld:latest"
    cpu      = "0.5"
    memory   = "0.5"
    port     = 5443
    protocol = "udp"
  }

  image_registry_credential {
    server   = "hub.docker.com"
    username = "yourusername"
    password = "yourpassword"
  }

  image_registry_credential {
    server   = "mine.acr.io"
    username = "acrusername"
    password = "acrpassword"
  }

  container {
    name   = "sidecar"
    image  = "microsoft/aci-tutorial-sidecar"
    cpu    = "0.5"
    memory = "0.5"
  }

  tags {
    environment = "Testing"
  }
}
`, ri, location, ri)
}

func testAccAzureRMContainerGroup_imageRegistryCredentialsUpdated(ri int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_container_group" "test" {
  name                = "acctestcontainergroup-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  ip_address_type     = "public"
  os_type             = "Linux"

  container {
    name   = "hw"
    image  = "microsoft/aci-helloworld:latest"
    cpu    = "0.5"
    memory = "0.5"
    ports  = {
      port = 80
    }
  }

  image_registry_credential {
    server   = "hub.docker.com"
    username = "updatedusername"
    password = "updatedpassword"
  }

  container {
    name   = "sidecar"
    image  = "microsoft/aci-tutorial-sidecar"
    cpu    = "0.5"
    memory = "0.5"
  }

  tags {
    environment = "Testing"
  }
}
`, ri, location, ri)
}

func testAccAzureRMContainerGroup_linuxBasicUpdated(ri int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_container_group" "test" {
  name                = "acctestcontainergroup-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  ip_address_type     = "public"
  os_type             = "Linux"

  container {
    name   = "hw"
    image  = "microsoft/aci-helloworld:latest"
    cpu    = "0.5"
    memory = "0.5"
    ports  = {
      port     = 80
    }
    ports  = {
      port     = 5443
      protocol = "udp"
    }
  }

  container {
    name   = "sidecar"
    image  = "microsoft/aci-tutorial-sidecar"
    cpu    = "0.5"
    memory = "0.5"
  }

  tags {
    environment = "Testing"
  }
}
`, ri, location, ri)
}

func testAccAzureRMContainerGroup_windowsBasic(ri int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_container_group" "test" {
  name                = "acctestcontainergroup-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  ip_address_type     = "public"
  os_type             = "windows"

  container {
    name   = "windowsservercore"
    image  = "microsoft/windowsservercore:latest"
    cpu    = "2.0"
    memory = "3.5"
    ports  = {
      port     = 80
      protocol = "tcp"
    }
    ports  = {
      port     = 443
      protocol = "tcp"
    }
  }

  tags {
    environment = "Testing"
  }
}
`, ri, location, ri)
}

func testAccAzureRMContainerGroup_windowsComplete(ri int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_container_group" "test" {
  name                = "acctestcontainergroup-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  ip_address_type     = "public"
  dns_name_label      = "acctestcontainergroup-%d"
  os_type             = "windows"
  restart_policy      = "Never"

  container {
    name   = "windowsservercore"
    image  = "microsoft/windowsservercore:latest"
    cpu    = "2.0"
    memory = "3.5"
    ports  = {
      port     = 80
      protocol = "tcp"
    }

    environment_variables {
      "foo"  = "bar"
      "foo1" = "bar1"
    }

    secure_environment_variables {
      "secureFoo"  = "secureBar"
      "secureFoo1" = "secureBar1"
    }

    commands = ["cmd.exe", "echo", "hi"]
  }

  tags {
    environment = "Testing"
  }
}
`, ri, location, ri, ri)
}

func testAccAzureRMContainerGroup_linuxComplete(ri int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_storage_account" "test" {
  name                     = "accsa%d"
  resource_group_name      = "${azurerm_resource_group.test.name}"
  location                 = "${azurerm_resource_group.test.location}"
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_storage_share" "test" {
  name = "acctestss-%d"

  resource_group_name  = "${azurerm_resource_group.test.name}"
  storage_account_name = "${azurerm_storage_account.test.name}"

  quota = 50
}

resource "azurerm_container_group" "test" {
  name                = "acctestcontainergroup-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  ip_address_type     = "public"
  dns_name_label      = "acctestcontainergroup-%d"
  os_type             = "linux"
  restart_policy      = "OnFailure"

  container {
    name   = "hf"
    image  = "seanmckenna/aci-hellofiles"
    cpu    = "1"
    memory = "1.5"

    ports  = {
      port     = 80
      protocol = "tcp"
    }

    volume {
      name       = "logs"
      mount_path = "/aci/logs"
      read_only  = false
      share_name = "${azurerm_storage_share.test.name}"

      storage_account_name = "${azurerm_storage_account.test.name}"
      storage_account_key = "${azurerm_storage_account.test.primary_access_key}"
    }

    environment_variables {
      "foo" = "bar"
      "foo1" = "bar1"
    }

    secure_environment_variables {
      "secureFoo"  = "secureBar"
      "secureFoo1" = "secureBar1"
    }

    commands = ["/bin/bash", "-c", "ls"]
  }

  tags {
    environment = "Testing"
  }
}
`, ri, location, ri, ri, ri, ri)
}

func testCheckAzureRMContainerGroupExists(name string, checkResp func(containerinstance.ContainerGroup) error) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for Container Registry: %s", name)
		}

		conn := testAccProvider.Meta().(*ArmClient).containerGroupsClient
		ctx := testAccProvider.Meta().(*ArmClient).StopContext

		resp, err := conn.Get(ctx, resourceGroup, name)
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: Container Group %q (resource group: %q) does not exist", name, resourceGroup)
			}
			return fmt.Errorf("Bad: Get on containerGroupsClient: %+v", err)
		}
		return checkResp(resp)
	}
}

func testCheckAzureRMContainerGroupDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*ArmClient).containerGroupsClient
	ctx := testAccProvider.Meta().(*ArmClient).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_container_group" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		resp, err := conn.Get(ctx, resourceGroup, name)

		if err != nil {
			if resp.StatusCode != http.StatusNotFound {
				return fmt.Errorf("Container Group still exists:\n%#v", resp)
			}

			return nil
		}

	}

	return nil
}
