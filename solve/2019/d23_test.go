package solve2019

import "testing"

func TestNetworkTurnProcessesEmptyReadBeforeWaiting(t *testing.T) {
	nic := initNICs("3,100,104,255,104,17,3,101,104,43,3,102,1105,1,10", 1)[0]
	for _, want := range []int64{255, 17, 43} {
		got, waiting, activity, err := nic.turn()
		if err != nil || waiting || !activity || got == nil || *got != want {
			t.Fatalf("output: %v %t %t %v want %d", got, waiting, activity, err, want)
		}
	}
	if got, waiting, activity, err := nic.turn(); err != nil || got != nil || !waiting || activity {
		t.Fatalf("waiting: %v %t %t %v", got, waiting, activity, err)
	}
	if nic.vm.mem[101] != -1 {
		t.Fatal("empty input was not processed before output")
	}
}

func TestNetworkIdleRequiresCompletePackets(t *testing.T) {
	for _, count := range []int{1, 2} {
		if networkIdle([]networkNIC{{outputCount: count}}, true, false) {
			t.Fatalf("partial packet of %d values called idle", count)
		}
	}
	if networkIdle([]networkNIC{{input: []int64{17, 43}}}, true, false) {
		t.Fatal("queued packet called idle")
	}
	if networkIdle([]networkNIC{{}}, false, false) || networkIdle([]networkNIC{{}}, true, true) {
		t.Fatal("active network called idle")
	}
	if !networkIdle([]networkNIC{{}}, true, false) {
		t.Fatal("quiescent network not idle")
	}
}

func TestNetworkFirstPacketAndLatestNAT(t *testing.T) {
	program := "3,100,104,255,104,17,104,43,104,255,104,19,104,43,104,255,104,23,104,47,3,101,1105,1,20"
	for _, count := range []int{1, 3, 50} {
		for _, tc := range []struct {
			solve func([]networkNIC) (int64, error)
			want  int64
		}{{firstNetworkY, 43}, {repeatedNATY, 47}} {
			got, err := tc.solve(initNICs(program, count))
			if err != nil || got != tc.want {
				t.Fatalf("%d computers: %d %v want %d", count, got, err, tc.want)
			}
		}
	}
}

func TestNetworkRoutesPacketValuesInOrder(t *testing.T) {
	// The first computer sends a packet to the second. The receiver ignores -1,
	// then sends the received X and Y to the NAT.
	sender := parseIntcode("104,1,104,37,104,59,3,100,1105,1,6")
	receiver := parseIntcode("3,100,1008,100,-1,101,1005,101,0,3,102,104,255,4,100,4,102,3,103,1105,1,17")
	got, err := firstNetworkY([]networkNIC{{vm: sender}, {vm: receiver}})
	if err != nil || got != 59 {
		t.Fatalf("routed packet: %d %v", got, err)
	}
}

func TestNetworkProtocolErrors(t *testing.T) {
	for _, program := range []string{"99", "104,77,104,17,104,43,99", "104,255,99"} {
		if _, err := repeatedNATY(initNICs(program, 1)); err == nil {
			t.Fatalf("accepted broken protocol %q", program)
		}
	}
}

func TestNetworkRejectsMalformedWords(t *testing.T) {
	for _, input := range []string{"", "3,,99", "3,x,99", "99,", "9223372036854775808", "99, 0"} {
		for _, fn := range []func(string) (string, error){(Day23{}).Part1, (Day23{}).Part2} {
			if _, err := fn(input); err == nil {
				t.Fatalf("accepted malformed program %q", input)
			}
		}
	}
}

func TestNetworkEventResumesCurrentSweep(t *testing.T) {
	network := newPacketNetwork(initNICs("3,100,104,255,4,100,4,100,3,101,1105,1,8", 3))
	for _, want := range []int64{0, 1, 2} {
		packet, idle, err := network.nextEvent()
		if err != nil || idle || packet == nil || packet.x != want || packet.y != want {
			t.Fatalf("packet for %d: %v %t %v", want, packet, idle, err)
		}
	}
	packet, idle, err := network.nextEvent()
	if err != nil || packet != nil || !idle {
		t.Fatalf("idle after all packets: %v %t %v", packet, idle, err)
	}
}
