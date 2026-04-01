package homelab

import "testing"

func TestParseContainerStatuses(t *testing.T) {
	t.Parallel()

	raw := "aurelia\tUp 10 minutes\t127.0.0.1:3334->3334/tcp\nsupabase-storage\tRestarting (1) 20 seconds ago\t"
	items := parseContainerStatuses(raw)

	if len(items) != 2 {
		t.Fatalf("len(items) = %d", len(items))
	}
	if !items[0].Up || items[1].Up {
		t.Fatalf("unexpected up flags: %#v", items)
	}
	if items[1].Ports != "-" {
		t.Fatalf("expected fallback port '-', got %q", items[1].Ports)
	}
}

func TestParseZFSSnapshotsKeepsMostRecentTail(t *testing.T) {
	t.Parallel()

	raw := ""
	for i := 0; i < 20; i++ {
		raw += "tank/data@snap-" + string(rune('a'+i)) + "\tWed Mar 26 10:00 2026\n"
	}

	items := parseZFSSnapshots(raw)
	if len(items) != maxSnapshots {
		t.Fatalf("len(items) = %d", len(items))
	}
	if items[0].Name != "tank/data@snap-f" {
		t.Fatalf("unexpected first snapshot: %#v", items[0])
	}
}

func TestParseDiskUsage(t *testing.T) {
	t.Parallel()

	raw := "Filesystem      Size  Used Avail Use% Mounted on\n/dev/nvme0n1p1  916G  483G  387G  56% /srv"
	disk := parseDiskUsage(raw)

	if disk.Filesystem != "/dev/nvme0n1p1" || disk.UsedPercent != 56 {
		t.Fatalf("disk = %#v", disk)
	}
}

func TestTailLinesDropsEmptyRows(t *testing.T) {
	t.Parallel()

	got := tailLines([]string{"", "a", " ", "b", "c"}, 2)
	if got != "b\nc" {
		t.Fatalf("tailLines() = %q", got)
	}
}

func TestBuildHomelabServices_Healthy(t *testing.T) {
	t.Parallel()

	services, summary, counts := buildHomelabServices([]ContainerStatus{
		{Name: "supabase-db", Status: "Up 10 minutes", Up: true},
		{Name: "realtime-dev.supabase-realtime", Status: "Up 8 minutes", Up: true},
		{Name: "qdrant", Status: "Up 3 minutes", Up: true},
		{Name: "captain-api", Status: "Up 1 minute", Up: true},
	})

	assertCanonicalServices(t, services)
	if summary.Status != homelabStatusHealthy || summary.Healthy != 2 || summary.Degraded != 0 || summary.Offline != 0 {
		t.Fatalf("summary = %#v", summary)
	}
	if counts.Services != 2 || counts.Containers != 4 || counts.UpContainers != 4 {
		t.Fatalf("counts = %#v", counts)
	}
	for _, service := range services {
		if service.Status != homelabStatusHealthy {
			t.Fatalf("service %q expected healthy, got %#v", service.Name, service)
		}
	}
}

func TestBuildHomelabServices_Degraded(t *testing.T) {
	t.Parallel()

	services, summary, counts := buildHomelabServices([]ContainerStatus{
		{Name: "supabase-db", Status: "Up 10 minutes", Up: true},
		{Name: "supabase-realtime", Status: "Restarting (1) 20 seconds ago"},
		{Name: "qdrant", Status: "Up 3 minutes", Up: true},
	})

	assertCanonicalServices(t, services)
	if summary.Status != homelabStatusDegraded || summary.Healthy != 1 || summary.Degraded != 0 || summary.Offline != 1 {
		t.Fatalf("summary = %#v", summary)
	}
	if counts.Services != 2 || counts.Containers != 3 || counts.UpContainers != 2 || counts.RestartingContainers != 1 {
		t.Fatalf("counts = %#v", counts)
	}
	if services[0].Status != homelabStatusHealthy {
		t.Fatalf("qdrant expected healthy, got %#v", services[0])
	}
	if services[1].Status != homelabStatusOffline {
		t.Fatalf("caprover expected offline, got %#v", services[1])
	}
}

func TestBuildHomelabServices_Offline(t *testing.T) {
	t.Parallel()

	services, summary, counts := buildHomelabServices([]ContainerStatus{
		{Name: "supabase-db", Status: "Exited (0) 2 minutes ago"},
		{Name: "qdrant", Status: "Dead"},
		{Name: "captain-worker", Status: "Restarting (1) 10 seconds ago"},
	})

	assertCanonicalServices(t, services)
	if summary.Status != homelabStatusOffline || summary.Healthy != 0 || summary.Degraded != 0 || summary.Offline != 2 {
		t.Fatalf("summary = %#v", summary)
	}
	if counts.Services != 2 || counts.Containers != 3 || counts.UpContainers != 0 || counts.RestartingContainers != 1 || counts.ExitedContainers != 1 || counts.DeadContainers != 1 {
		t.Fatalf("counts = %#v", counts)
	}
	for _, service := range services {
		if service.Status != homelabStatusOffline {
			t.Fatalf("service %q expected offline, got %#v", service.Name, service)
		}
	}
}

func TestBuildHomelabServices_EmitsExactlyTwoCanonicalServices(t *testing.T) {
	t.Parallel()

	services, summary, counts := buildHomelabServices(nil)

	assertCanonicalServices(t, services)
	if summary.Status != homelabStatusOffline || summary.Healthy != 0 || summary.Degraded != 0 || summary.Offline != 2 {
		t.Fatalf("summary = %#v", summary)
	}
	if counts.Services != 2 || counts.Containers != 0 {
		t.Fatalf("counts = %#v", counts)
	}
}

func assertCanonicalServices(t *testing.T, services []HomelabService) {
	t.Helper()

	if len(services) != 2 {
		t.Fatalf("len(services) = %d", len(services))
	}
	want := []string{homelabServiceQdrant, homelabServiceCaprover}
	for i, service := range services {
		if service.Name != want[i] {
			t.Fatalf("service[%d].Name = %q, want %q", i, service.Name, want[i])
		}
	}
}
