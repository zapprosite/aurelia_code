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
