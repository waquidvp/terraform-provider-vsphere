package vsphere

import (
	"fmt"
	"github.com/hashicorp/terraform-provider-vsphere/vsphere/internal/helper/testhelper"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccDataSourceVSphereNetwork_dvsPortgroup(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			RunSweepers()
			testAccPreCheck(t)
			testAccDataSourceVSphereNetworkPreCheck(t)
			testAccSkipIfEsxi(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceVSphereNetworkConfigDVSPortgroup(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.vsphere_network.net", "type", "DistributedVirtualPortgroup"),
					resource.TestCheckResourceAttrPair(
						"data.vsphere_network.net", "id",
						"vsphere_distributed_port_group.pg", "id",
					),
				),
			},
		},
	})
}

func TestAccDataSourceVSphereNetwork_absolutePathNoDatacenter(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			RunSweepers()
			testAccPreCheck(t)
			testAccDataSourceVSphereNetworkPreCheck(t)
			testAccSkipIfEsxi(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceVSphereNetworkConfigDVSPortgroupAbsolute(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.vsphere_network.net", "type", "DistributedVirtualPortgroup"),
					resource.TestCheckResourceAttrPair(
						"data.vsphere_network.net", "id",
						"vsphere_distributed_port_group.pg", "id",
					),
				),
			},
		},
	})
}

func TestAccDataSourceVSphereNetwork_hostPortgroups(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			RunSweepers()
			testAccPreCheck(t)
			testAccDataSourceVSphereNetworkPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceVSphereNetworkConfigHostPortgroup(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.vsphere_network.net", "type", "Network"),
				),
			},
		},
	})
}

func testAccDataSourceVSphereNetworkPreCheck(t *testing.T) {
	if os.Getenv("TF_VAR_VSPHERE_PG_NAME") == "" {
		t.Skip("set TF_VAR_VSPHERE_PG_NAME to run vsphere_network acceptance tests")
	}
}

func testAccDataSourceVSphereNetworkConfigDVSPortgroup() string {
	return fmt.Sprintf(`
%s

resource "vsphere_distributed_virtual_switch" "dvs" {
  name          = "testacc-dvs"
  datacenter_id = "${data.vsphere_datacenter.rootdc1.id}"
}

resource "vsphere_distributed_port_group" "pg" {
  name                            = "terraform-test-pg"
  distributed_virtual_switch_uuid = "${vsphere_distributed_virtual_switch.dvs.id}"
}

data "vsphere_network" "net" {
  name          = "${vsphere_distributed_port_group.pg.name}"
  datacenter_id = "${data.vsphere_datacenter.rootdc1.id}"
  distributed_virtual_switch_uuid = "${vsphere_distributed_virtual_switch.dvs.id}"
}
`,
		testhelper.CombineConfigs(testhelper.ConfigDataRootDC1(), testhelper.ConfigDataRootPortGroup1()),
	)
}

func testAccDataSourceVSphereNetworkConfigDVSPortgroupAbsolute() string {
	return fmt.Sprintf(`
%s

resource "vsphere_distributed_virtual_switch" "dvs" {
  name          = "testacc-dvs"
  datacenter_id = "${data.vsphere_datacenter.rootdc1.id}"
}

resource "vsphere_distributed_port_group" "pg" {
  name                            = "terraform-test-pg"
  distributed_virtual_switch_uuid = "${vsphere_distributed_virtual_switch.dvs.id}"
}

data "vsphere_network" "net" {
  name          = "/${data.vsphere_datacenter.rootdc1.name}/network/${vsphere_distributed_port_group.pg.name}"
}
`,
		testhelper.CombineConfigs(testhelper.ConfigDataRootDC1(), testhelper.ConfigDataRootPortGroup1()),
	)
}

func testAccDataSourceVSphereNetworkConfigHostPortgroup() string {
	return fmt.Sprintf(`
%s

resource "vsphere_host_virtual_switch" "switch" {
  name           = "vSwitchTerraformTest"
  host_system_id = "${data.vsphere_host.esxi_host.id}"

  network_adapters = ["${var.host_nic0}", "${var.host_nic1}"]
  active_nics      = ["${var.host_nic0}", "${var.host_nic1}"]
  standby_nics     = []
}

resource "vsphere_host_port_group" "pg" {
  name                = "PGTerraformTest"
  host_system_id      = "${data.vsphere_host.esxi_host.id}"
  virtual_switch_name = "${vsphere_host_virtual_switch.switch.name}"
}


`,
		testhelper.CombineConfigs(testhelper.ConfigDataRootDC1(), testhelper.ConfigDataRootPortGroup1()),
	)
}
