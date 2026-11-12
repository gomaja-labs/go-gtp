#!/usr/bin/env bash
# End-to-end check of the gw-tester: start SGW, PGW, MME and eNB on the
# loopback interface and wait until every node reports all subscribers
# attached. The nodes configure addresses and routes on lo, so they run as
# root (through sudo when needed). Run from anywhere; Linux only.
set -euo pipefail

here="$(cd "$(dirname "$0")" && pwd)"
examples="$(cd "$here/.." && pwd)"
work="$(mktemp -d)"
sudo=""
if [ "$(id -u)" != 0 ]; then
	sudo="sudo"
fi

want="$(grep -cE '^[[:space:]]*- imsi:' "$here/enb/enb.yml")"
timeout_s="${E2E_TIMEOUT:-60}"

pids=()
cleanup() {
	for pid in "${pids[@]}"; do
		$sudo kill "$pid" 2>/dev/null || true
	done
	wait 2>/dev/null || true
	rm -rf "$work"
}
trap cleanup EXIT

(cd "$examples" && for node in sgw pgw mme enb; do
	go build -o "$work/$node" "./gw-tester/$node"
done)

for node in sgw pgw mme enb; do
	$sudo "$work/$node" -config "$here/$node/$node.yml" >"$work/$node.log" 2>&1 &
	pids+=("$!")
	sleep 1
done

# node name and metrics address, as configured by prom_addr in each yml.
nodes=(enb:127.0.10.1 mme:127.0.10.2 sgw:127.0.10.3 pgw:127.0.10.4)

# sessions prints the active session count of a node, or nothing while its
# metrics endpoint is not up yet.
sessions() {
	{ curl -fsS "http://$2:58080/metrics" 2>/dev/null || true; } |
		awk -v m="$1_active_sessions_total" '$1 == m { print $2 }'
}

deadline=$((SECONDS + timeout_s))
while :; do
	ok=1
	report=""
	for n in "${nodes[@]}"; do
		got="$(sessions "${n%%:*}" "${n#*:}")"
		report+="${n%%:*}=${got:-?} "
		[ "${got:-}" = "$want" ] || ok=0
	done
	if [ "$ok" = 1 ]; then
		echo "all nodes report $want active sessions: $report"
		exit 0
	fi
	if [ "$SECONDS" -ge "$deadline" ]; then
		echo "timed out after ${timeout_s}s waiting for $want sessions: $report" >&2
		for node in sgw pgw mme enb; do
			echo "----- $node log -----" >&2
			cat "$work/$node.log" >&2 || true
		done
		exit 1
	fi
	sleep 1
done
