package gudev

import (
    "testing"
)

func TestClient(t *testing.T)  {
    subsystems := []string{"power_supply"}
    client := NewClient(subsystems)
    t.Log("client:", client)

    devices := client.QueryBySubsystem("power_supply")
    for _, device := range devices {
        t.Log("device",device)
        sysfsPath := device.GetSysfsPath()
        t.Log("sysfs path:", sysfsPath)
    }
    client.Unref()
}
