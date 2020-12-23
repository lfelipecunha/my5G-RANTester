UPFNS="UPFns"
EXEC_UPFNS="sudo -E ip netns exec ${UPFNS}"

sudo ip netns add ${UPFNS}

sudo ip link add veth0 type veth peer name veth1
sudo ip link set veth0 up
sudo ip addr add 60.60.0.1 dev lo
sudo ip addr add 60.60.0.2 dev lo
sudo ip addr add 60.60.0.3 dev lo
sudo ip addr add 60.60.0.4 dev lo
sudo ip addr add 60.60.0.5 dev lo
sudo ip addr add 60.60.0.6 dev lo
sudo ip addr add 60.60.0.7 dev lo
sudo ip addr add 60.60.0.8 dev lo
sudo ip addr add 60.60.0.9 dev lo
sudo ip addr add 60.60.0.10 dev lo
sudo ip addr add 60.60.0.11 dev lo
sudo ip addr add 60.60.0.12 dev lo
sudo ip addr add 60.60.0.13 dev lo
sudo ip addr add 60.60.0.14 dev lo
sudo ip addr add 60.60.0.15 dev lo
sudo ip addr add 60.60.0.16 dev lo
sudo ip addr add 60.60.0.17 dev lo
sudo ip addr add 60.60.0.18 dev lo
sudo ip addr add 60.60.0.19 dev lo
sudo ip addr add 60.60.0.20 dev lo
sudo ip addr add 60.60.0.21 dev lo
sudo ip addr add 60.60.0.22 dev lo
sudo ip addr add 60.60.0.23 dev lo
sudo ip addr add 60.60.0.24 dev lo
sudo ip addr add 60.60.0.25 dev lo
sudo ip addr add 60.60.0.26 dev lo
sudo ip addr add 60.60.0.27 dev lo
sudo ip addr add 60.60.0.28 dev lo
sudo ip addr add 60.60.0.29 dev lo
sudo ip addr add 60.60.0.30 dev lo
sudo ip addr add 60.60.0.31 dev lo
sudo ip addr add 60.60.0.32 dev lo
sudo ip addr add 60.60.0.33 dev lo
sudo ip addr add 60.60.0.34 dev lo
sudo ip addr add 60.60.0.35 dev lo
sudo ip addr add 60.60.0.36 dev lo
sudo ip addr add 60.60.0.37 dev lo
sudo ip addr add 60.60.0.38 dev lo
sudo ip addr add 60.60.0.39 dev lo
sudo ip addr add 60.60.0.40 dev lo
sudo ip addr add 60.60.0.41 dev lo
sudo ip addr add 60.60.0.42 dev lo
sudo ip addr add 60.60.0.43 dev lo
sudo ip addr add 60.60.0.44 dev lo
sudo ip addr add 60.60.0.45 dev lo
sudo ip addr add 60.60.0.46 dev lo
sudo ip addr add 60.60.0.47 dev lo
sudo ip addr add 60.60.0.48 dev lo
sudo ip addr add 60.60.0.49 dev lo
sudo ip addr add 60.60.0.50 dev lo
sudo ip addr add 60.60.0.51 dev lo
sudo ip addr add 60.60.0.52 dev lo
sudo ip addr add 10.200.200.1/24 dev veth0
sudo ip addr add 10.200.200.2/24 dev veth0

sudo ip link set veth1 netns ${UPFNS}

${EXEC_UPFNS} ip link set lo up
${EXEC_UPFNS} ip link set veth1 up
${EXEC_UPFNS} ip addr add 60.60.0.101 dev lo
${EXEC_UPFNS} ip addr add 10.200.200.101/24 dev veth1
${EXEC_UPFNS} ip addr add 10.200.200.102/24 dev veth1